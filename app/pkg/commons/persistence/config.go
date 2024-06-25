// Package persistence defines the framework to interact with the database storage.
package persistence

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

const defaultMappersPath = "include/config/mappers"

// StatementConfig identifies a statement configuration
type StatementConfig struct {
	ID      string `xml:"id,attr"`
	Content string `xml:",chardata"`
}

type StorageConfigOptions struct {
	ReadDirAndFileFS ReadDirAndFileFS
	MappersPath      string
	Driver           string
}

// StorageConfig configuration for StatementConfig
type StorageConfig struct {
	dataMap map[string]StatementConfig
}

// MapperConfig identifies a mapper configuration for a set of StatementConfig
type MapperConfig struct {
	ID         string            `xml:"id,attr"`
	Statements []StatementConfig `xml:"statement"`
}

type ReadDirAndFileFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

// NewStorageConfig returns a new StorageConfig with the provided configuration options.
func NewStorageConfig(options StorageConfigOptions) (*StorageConfig, error) {
	if options.ReadDirAndFileFS == nil {
		return nil, errors.New("ReadDirAndFileFS cannot be nil to create a storage config")
	}
	mappersPath := defaultMappersPath

	if options.MappersPath != "" {
		mappersPath = options.MappersPath
	}

	mappersPathForDriver := path.Join(mappersPath, options.Driver)

	dirEntries, err := options.ReadDirAndFileFS.ReadDir(mappersPathForDriver)
	if err != nil {
		return nil, err
	}

	configFiles := make([]io.ReadCloser, 0)

	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".xml") {
			file, readFileErr := options.ReadDirAndFileFS.ReadFile(path.Join(mappersPathForDriver, entry.Name()))
			if readFileErr != nil {
				return nil, readFileErr
			}

			configFiles = append(configFiles, io.NopCloser(strings.NewReader(string(file))))
		}
	}

	persistenceConfig := newPersistenceFwConfig()
	for _, currentConfigFile := range configFiles {
		dataBytes, readAllErr := io.ReadAll(currentConfigFile)
		if readAllErr != nil {
			return nil, fmt.Errorf("can not read mapper file. %w", readAllErr)
		}
		extraConfig, readConfigXMLFileErr := readConfigXMLFileBytes(dataBytes)
		if readConfigXMLFileErr != nil {
			return nil, fmt.Errorf("%w", readConfigXMLFileErr)
		}
		persistenceConfig.addConfig(extraConfig)
	}

	return &persistenceConfig, nil
}

// NewEmptyStorageConfig returns an empty StorageConfig
func NewEmptyStorageConfig() StorageConfig {
	return newPersistenceFwConfig()
}

// AddConfig adds a new configuration to the configurations data map.
func (config *StorageConfig) AddConfig(newFwConfig StorageConfig) error {
	for k := range newFwConfig.dataMap {
		_, exists := config.dataMap[k]
		if exists {
			return fmt.Errorf("a statement config already exists with id [%s]", k)
		}
	}

	config.addConfig(&newFwConfig)

	return nil
}

// GetStatement obtain the StatementConfig with the AdminID provided
func (config *StorageConfig) GetStatement(id string) (StatementConfig, bool) {
	stmtConfig, ok := config.dataMap[id]
	return stmtConfig, ok
}

func (config *StorageConfig) addConfig(extraConfig *StorageConfig) {
	for key, value := range extraConfig.dataMap {
		config.dataMap[key] = value
	}
}

func (config *StorageConfig) addStatement(mapperId string, stmt StatementConfig) {
	config.dataMap[mapperId+"."+stmt.ID] = stmt
}

func newPersistenceFwConfig() StorageConfig {
	return StorageConfig{
		dataMap: make(map[string]StatementConfig),
	}
}

func readConfigXMLFileBytes(dataBytes []byte) (*StorageConfig, error) {
	var mapperConfig MapperConfig
	err := xml.Unmarshal(dataBytes, &mapperConfig)
	if err != nil {
		return nil, fmt.Errorf("File cand not be parsed as xml,\noriginal error: %w", err)
	}

	config := newPersistenceFwConfig()

	for counter := range mapperConfig.Statements {
		config.addStatement(mapperConfig.ID, mapperConfig.Statements[counter])
	}

	return &config, nil
}
