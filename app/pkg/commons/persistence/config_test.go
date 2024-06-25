package persistence_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
)

//go:embed testdata/mappers/*
var mappers embed.FS

func TestNewStorageConfig_Success(t *testing.T) {
	storageConfigOptions := persistence.StorageConfigOptions{
		ReadDirAndFileFS: mappers,
		MappersPath:      "testdata/mappers/initial",
		Driver:           "postgres",
	}
	storageConfig, err := persistence.NewStorageConfig(storageConfigOptions)
	require.Nil(t, err)

	_, ok := storageConfig.GetStatement("io.adhara.persistencefw.itest.insertUser")
	require.True(t, ok)

	_, ok = storageConfig.GetStatement("io.adhara.persistencefw.itest.additional.insertUser")
	require.False(t, ok)

	additionalStorageConfigOptions := persistence.StorageConfigOptions{
		ReadDirAndFileFS: mappers,
		MappersPath:      "testdata/mappers/additional",
		Driver:           "postgres",
	}
	additionalConfig, err := persistence.NewStorageConfig(additionalStorageConfigOptions)
	require.Nil(t, err)
	err = storageConfig.AddConfig(*additionalConfig)
	require.Nil(t, err)

	_, ok = storageConfig.GetStatement("io.adhara.persistencefw.itest.additional.insertUser")
	require.True(t, ok)
}
