package test_data

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	"github.com/makerdao/vulcanizedb/pkg/config"
)

type MockTransformerExporterUpdater struct {
	PreparePluginConfigCalled bool
	PluginConfigToReturn      config.Plugin
	PreparePluginErr          error
	WritePluginCalled         bool
	PluginConfigPassedIn      config.Plugin
	WritePluginErr            error
}

func (p *MockTransformerExporterUpdater) PreparePluginConfig(updatedConfig types.TransformersConfig) (config.Plugin, error) {
	p.PreparePluginConfigCalled = true
	return p.PluginConfigToReturn, p.PreparePluginErr
}

func (p *MockTransformerExporterUpdater) WritePlugin(pluginConfig config.Plugin) error {
	p.WritePluginCalled = true
	p.PluginConfigPassedIn = pluginConfig
	return p.WritePluginErr
}
