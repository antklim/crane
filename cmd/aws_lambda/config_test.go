package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var envVars = [...]string{
	"CRANE_ARCHIVE_BUCKET",
	"CRANE_ARCHIVE_FOLDER",
	"CRANE_DEPLOY_BUCKET",
	"CRANE_STAGE_BUCKET",
	"CRANE_PRODUCTION_BUCKET",
}

var envSet = make(map[string]string, len(envVars))

func setup(t *testing.T) {
	for _, envVar := range envVars {
		envSet[envVar] = os.Getenv(envVar)
		err := os.Setenv(envVar, envVar)
		require.NoError(t, err)
	}
}

func teardown(t *testing.T) {
	for key, value := range envSet {
		var err error
		if value == "" {
			err = os.Unsetenv(key)
		} else {
			err = os.Setenv(key, value)
		}
		if err != nil {
			t.Logf("Error occurred during test teardown: %+v", err)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	setup(t)

	expectedConfig := &config{
		ArchiveBucket:    "CRANE_ARCHIVE_BUCKET",
		ArchiveFolder:    "CRANE_ARCHIVE_FOLDER",
		DeployBucket:     "CRANE_DEPLOY_BUCKET",
		StageBucket:      "CRANE_STAGE_BUCKET",
		ProductionBucket: "CRANE_PRODUCTION_BUCKET",
		Region:           "ap-southeast-2",
	}

	config, err := loadConfig()
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, config)

	teardown(t)
}
