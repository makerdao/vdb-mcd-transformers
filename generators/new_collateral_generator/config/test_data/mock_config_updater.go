package test_data

import "github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"

type MockConfigUpdater struct {
	SetCurrentConfigCalled bool
	InitialConfigPassedIn  config.TransformersConfig
	AddNewCollateralCalled bool
	AddNewCollateralErr    error
	GetUpdatedConfigCalled bool
	UpdatedConfig          config.TransformersConfigForTomlEncoding
}

func (cu *MockConfigUpdater) SetInitialConfig(initialConfig config.TransformersConfig) {
	cu.InitialConfigPassedIn = initialConfig
	cu.SetCurrentConfigCalled = true
}

func (cu *MockConfigUpdater) AddNewCollateralToConfig() error {
	cu.AddNewCollateralCalled = true
	return cu.AddNewCollateralErr
}

func (cu *MockConfigUpdater) GetUpdatedConfig() config.TransformersConfigForTomlEncoding {
	cu.GetUpdatedConfigCalled = true
	return cu.UpdatedConfig
}
