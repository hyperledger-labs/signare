// Package dbmigrator defines the functionalities for database migrations.
package dbmigrator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql"
)

const (
	dbSchemas         = "include/dbschemas"
	migrationFileName = "db_migration_steps.yml"
)

var ErrMigrationInProgress = errors.New("there is already a migration in progress")

// MigrationStep defines the step of a migration process.
type MigrationStep struct {
	UpFile             io.ReadCloser
	DownFile           io.ReadCloser
	VersionDescription string
}

// Migration defines a database migration.
type Migration struct {
	Steps                 []MigrationStep
	TargetVersion         int
	MigrationsTablePrefix *string
}

type migrationStep struct {
	file        io.ReadCloser
	version     int
	description string
}

// DbMigrator handles database migrations.
type DbMigrator struct {
	connection sql.Connection
}

// DbMigratorOptions defines the options to configure a DbMigrator.
type DbMigratorOptions struct {
	Connection sql.Connection
}

// MigrateFromFilesInput defines the input data to migrate database from files.
type MigrateFromFilesInput struct {
	FS                    fs.ReadDirFS
	TargetVersion         *int
	MigrationsTablePrefix *string
}

type migrationFileConfig struct {
	Steps []stepConfig `yaml:"migration_steps"`
}

type stepConfig struct {
	Up                 string `yaml:"up"`
	Down               string `yaml:"down"`
	VersionDescription string `yaml:"version_description"`
}

var lock = &sync.Mutex{}
var migratorInstance *DbMigrator // singleton instance

func NewDbMigrator(options DbMigratorOptions) (*DbMigrator, error) {
	if migratorInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if migratorInstance == nil {
			if options.Connection == nil {
				return nil, errors.New("connection cannot be nil")
			}
			migratorInstance = &DbMigrator{
				connection: options.Connection,
			}
		}
	}

	return migratorInstance, nil
}

// MigrateFromFiles migrates the database using migration files.
// It reads migration steps from files specified in the input MigrateFromFilesInput
// and executes them sequentially.
func (d DbMigrator) MigrateFromFiles(ctx context.Context, input MigrateFromFilesInput) error {
	if input.FS == nil {
		return errors.New("FS cannot be nil")
	}

	migrationFile := path.Join(dbSchemas, d.connection.GetDialectName(), migrationFileName)

	migrationStepsFile, err := input.FS.Open(migrationFile)
	if err != nil {
		panic(err)
	}
	migrationStepsFileData, err := io.ReadAll(migrationStepsFile)
	if err != nil {
		panic(err)
	}

	var config migrationFileConfig
	err = yaml.Unmarshal(migrationStepsFileData, &config)
	if err != nil {
		panic(err)
	}
	dbMigration := Migration{
		Steps:         make([]MigrationStep, 0),
		TargetVersion: len(config.Steps),
	}

	if input.TargetVersion != nil {
		dbMigration.TargetVersion = *input.TargetVersion
	}
	for _, currentMigrationStepConfig := range config.Steps {
		upFileReader, innerErr := input.FS.Open(removeLeadingDashFromFilename(currentMigrationStepConfig.Up))
		if innerErr != nil {
			panic(innerErr)
		}
		downFileReader, innerErr := input.FS.Open(removeLeadingDashFromFilename(currentMigrationStepConfig.Down))
		if innerErr != nil {
			panic(innerErr)
		}
		newMigrationStep := MigrationStep{
			UpFile:             upFileReader,
			DownFile:           downFileReader,
			VersionDescription: currentMigrationStepConfig.VersionDescription,
		}
		dbMigration.Steps = append(dbMigration.Steps, newMigrationStep)
	}
	dbMigration.MigrationsTablePrefix = input.MigrationsTablePrefix

	return d.Migrate(ctx, dbMigration)
}

// Migrate executes the database migration process.
// It initializes the migration process, obtains the current schema version,
// and migrates the database to the target version.
//
//nolint:staticcheck
func (d DbMigrator) Migrate(ctx context.Context, migration Migration) error {
	lock.Lock()
	defer lock.Unlock()
	logger.LogEntry(ctx).Info("migration started")

	migrator := d.connection.GetMigrator()
	defer func() {
		logger.LogEntry(ctx).Info("closing connection to database")
		err := migrator.CloseConnection(ctx)
		if err != nil {
			fmt.Printf("error closing connection [%s]", err.Error())
		}
	}()

	logger.LogEntry(ctx).Info("opening connection to database")
	err := migrator.OpenConnection(ctx)
	if err != nil {
		return err
	}

	logger.LogEntry(ctx).Info("initializating migration process")
	err = migrator.InitMigration(ctx, migration.MigrationsTablePrefix)
	if err != nil {
		return err
	}

	logger.LogEntry(ctx).Info("obtain current schemas version")
	version, err := migrator.GetMigrationVersion(ctx)
	if err != nil {
		return err
	}
	if version.Dirty {
		return fmt.Errorf("migration table is dirty in version [%d]", version.Version)
	}

	var stepsToMigrate []migrationStep
	if migration.TargetVersion > version.Version {
		logger.LogEntry(ctx).Infof("upgrading to target version [%d] from version [%d]", migration.TargetVersion, version.Version)
		stepsToMigrate, err = mapMigrationsUp(version.Version, migration.TargetVersion, migration.Steps)
		if err != nil {
			return err
		}
	} else if migration.TargetVersion < version.Version {
		logger.LogEntry(ctx).Infof("downgrading to target version [%d] from version [%d]", migration.TargetVersion, version.Version)
		stepsToMigrate, err = mapMigrationsDown(version.Version, migration.TargetVersion, migration.Steps)
		if err != nil {
			return err
		}
	} else {
		logger.LogEntry(ctx).Info("nothing to migrate")
		return nil
	}

	for _, step := range stepsToMigrate {
		errMigrating := d.migrateStep(ctx, migrator, step)
		if errMigrating != nil {
			return errMigrating
		}
	}
	logger.LogEntry(ctx).Info("migration finished")
	return nil
}

func (d DbMigrator) migrateStep(ctx context.Context, migrator sql.Migrator, step migrationStep) error {
	migrationVersion := sql.SetVersionInput{
		MigrationVersion: sql.MigrationVersion{
			Version:     step.version,
			Dirty:       true,
			Description: step.description,
		},
	}

	logger.LogEntry(ctx).Infof("set migration step [%d] as dirty", migrationVersion.Version)
	err := migrator.SetMigrationVersion(ctx, migrationVersion)
	if err != nil {
		return err
	}

	defer step.file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(step.file)
	if err != nil {
		return err
	}
	logger.LogEntry(ctx).Infof("executing migration step [%d] with description [%s]", migrationVersion.Version, migrationVersion.Description)
	err = migrator.RunMigration(ctx, buf.String())
	if err != nil {
		return err
	}

	migrationVersion.Dirty = false
	logger.LogEntry(ctx).Infof("set migration step [%d] as not dirty", migrationVersion.Version)
	err = migrator.SetMigrationVersion(ctx, migrationVersion)
	if err != nil {
		return err
	}
	return nil
}

func mapMigrationsUp(currentVersion int, targetVersion int, steps []MigrationStep) ([]migrationStep, error) {
	stepsToMigrate := make([]migrationStep, 0)

	if targetVersion > len(steps) {
		return nil, fmt.Errorf("target version is higher [%d] than the available number of steps [%d]", targetVersion, len(steps))
	}

	for c := currentVersion; c < targetVersion; c++ {
		currentStep := steps[c]
		newMigration := migrationStep{
			file:        currentStep.UpFile,
			version:     c + 1,
			description: currentStep.VersionDescription,
		}
		stepsToMigrate = append(stepsToMigrate, newMigration)

	}

	return stepsToMigrate, nil
}

func mapMigrationsDown(currentVersion int, targetVersion int, steps []MigrationStep) ([]migrationStep, error) {
	stepsToMigrate := make([]migrationStep, 0)

	if targetVersion < 0 {
		return nil, fmt.Errorf("target version is lower [%d] than zero", targetVersion)
	}

	if targetVersion > len(steps) {
		return nil, fmt.Errorf("target version is higher [%d] than the available number of steps [%d]", targetVersion, len(steps))
	}

	for c := currentVersion; c > targetVersion; c-- {
		currentStep := steps[c-1]
		newMigration := migrationStep{
			file:        currentStep.DownFile,
			version:     c - 1,
			description: currentStep.VersionDescription,
		}
		stepsToMigrate = append(stepsToMigrate, newMigration)

	}

	return stepsToMigrate, nil
}

func removeLeadingDashFromFilename(filename string) string {
	if strings.HasPrefix(filename, "/") {
		return filename[1:]
	}
	return filename
}
