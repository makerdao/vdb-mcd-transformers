package test_data

import "github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"

type MockConfigParser struct {
	ConfigFilePathPassedIn string
	ConfigFileNamePassedIn string
	ParseErr               error
	ConfigToReturn         config.TransformersConfig
}

func (cp *MockConfigParser) ParseCurrentConfig(configFilePath, configFileName string) (config.TransformersConfig, error) {
	cp.ConfigFilePathPassedIn = configFilePath
	cp.ConfigFileNamePassedIn = configFileName
	return cp.ConfigToReturn, cp.ParseErr
}
