package transformer_exporter

import (
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	pluginConfig "github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/plugin/writer"
)

type TransformerExporterUpdater interface {
	WritePlugin(pluginConfig pluginConfig.Plugin) error
	PreparePluginConfig(updatedConfig types.TransformersConfig) (pluginConfig.Plugin, error)
}

type updater struct{}

func NewTransformerExporterUpdater() updater {
	return updater{}
}

func (p updater) WritePlugin(pluginConfig pluginConfig.Plugin) error {
	pluginWriter := writer.NewPluginWriter(pluginConfig)
	return pluginWriter.WritePlugin()
}

func (p updater) PreparePluginConfig(updatedConfig types.TransformersConfig) (pluginConfig.Plugin, error) {
	transformers := make(map[string]pluginConfig.Transformer)
	for k, v := range updatedConfig.TransformerExporters {
		rank, rankErr := strconv.Atoi(v.Rank)
		if rankErr != nil {
			return pluginConfig.Plugin{}, rankErr
		}
		transformers[k] = pluginConfig.Transformer{
			Path:           v.Path,
			Type:           getTransformerType(v.Type),
			MigrationPath:  v.Migrations,
			MigrationRank:  uint64(rank),
			RepositoryPath: v.Repository,
		}
	}

	return pluginConfig.Plugin{
		Transformers: transformers,
		FilePath:     helpers.GetExecutePluginsPath(),
		FileName:     updatedConfig.ExporterMetadata.Name,
		Save:         updatedConfig.ExporterMetadata.Save,
		Home:         updatedConfig.ExporterMetadata.Home,
		Schema:       updatedConfig.ExporterMetadata.Schema,
	}, nil
}

func getTransformerType(typeString string) pluginConfig.TransformerType {
	switch typeString {
	case "eth_event":
		return pluginConfig.EthEvent
	case "eth_storage":
		return pluginConfig.EthStorage
	default:
		return pluginConfig.UnknownTransformerType
	}
}
