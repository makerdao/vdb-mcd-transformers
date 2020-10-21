package test_data

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
)

type MockConfigParser struct {
	ConfigFilePathPassedIn string
	ConfigFileNamePassedIn string
	ParseErr               error
	ConfigToReturn         types.TransformersConfig
}

func (cp *MockConfigParser) ParseCurrentConfig(configFilePath, configFileName string) (types.TransformersConfig, error) {
	cp.ConfigFilePathPassedIn = configFilePath
	cp.ConfigFileNamePassedIn = configFileName
	return cp.ConfigToReturn, cp.ParseErr
}
