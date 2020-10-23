package new_collateral

import (
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/initializer"
	pluginConfig "github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/plugin/writer"
)

type NewCollateralGenerator struct {
	ConfigFileName       string
	ConfigFilePath       string
	ConfigParser         config.IParse
	ConfigUpdater        config.IUpdate
	InitializerGenerator initializer.IGenerate
}

func (g NewCollateralGenerator) UpdateConfig() error {
	initialConfig, parseConfigErr := g.ConfigParser.ParseCurrentConfig(g.ConfigFilePath, g.ConfigFileName)
	if parseConfigErr != nil {
		return parseConfigErr
	}

	g.ConfigUpdater.SetInitialConfig(initialConfig)

	updateErr := g.ConfigUpdater.AddNewCollateralToConfig()
	if updateErr != nil {
		return updateErr
	}

	file, fileOpenErr := os.Create(helpers.GetFullConfigFilePath(g.ConfigFilePath, g.ConfigFileName))
	if fileOpenErr != nil {
		return fileOpenErr
	}

	updatedConfig, updatedConfigErr := g.ConfigUpdater.GetUpdatedConfigForToml()
	if updatedConfigErr != nil {
		return updatedConfigErr
	}

	encodingErr := toml.NewEncoder(file).Encode(updatedConfig)
	if encodingErr != nil {
		return encodingErr
	}

	return file.Close()
}

func (g *NewCollateralGenerator) UpdatePluginExporter() error {
	pluginConfig, pluginErr := g.PreparePluginConfig()
	if pluginErr != nil {
		return pluginErr
	}

	pluginWriter := writer.NewPluginWriter(pluginConfig)
	return pluginWriter.WritePlugin()
}

func (g *NewCollateralGenerator) PreparePluginConfig() (pluginConfig.Plugin, error) {
	updatedConfig := g.ConfigUpdater.GetUpdatedConfig()
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
func (g *NewCollateralGenerator) WriteInitializers() error {
	flipInitializersErr := g.InitializerGenerator.GenerateFlipInitializer()
	if flipInitializersErr != nil {
		return flipInitializersErr
	}
	return g.InitializerGenerator.GenerateMedianInitializer()
}
