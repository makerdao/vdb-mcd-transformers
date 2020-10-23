package test_data

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
)

type MockConfigUpdater struct {
	SetCurrentConfigCalled bool
	InitialConfigPassedIn  types.TransformersConfig
	AddNewCollateralCalled bool
	AddNewCollateralErr    error
	GetUpdatedConfigCalled bool
	UpdatedConfig          types.TransformersConfig
}

func (cu *MockConfigUpdater) SetInitialConfig(initialConfig types.TransformersConfig) {
	cu.InitialConfigPassedIn = initialConfig
	cu.SetCurrentConfigCalled = true
}

func (cu *MockConfigUpdater) AddNewCollateralToConfig() error {
	cu.AddNewCollateralCalled = true
	return cu.AddNewCollateralErr
}

func (cu *MockConfigUpdater) GetUpdatedConfig() types.TransformersConfig {
	cu.GetUpdatedConfigCalled = true
	return cu.UpdatedConfig
}
