package generator

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"
)


type NewCollateralGenerator struct {
	ConfigFileName string
	ConfigFilePath string
	ConfigParser config.IParse
	ConfigUpdater config.IUpdate
}

func (g NewCollateralGenerator) AddToConfig() error{
	initialConfig, parseConfigErr := g.ConfigParser.ParseCurrentConfig(g.ConfigFilePath, g.ConfigFileName)
	if parseConfigErr != nil {
		return parseConfigErr
	}

	g.ConfigUpdater.SetInitialConfig(initialConfig)

	updateErr := g.ConfigUpdater.AddNewCollateralToConfig()
	if updateErr != nil {
		return updateErr
	}

	file, fileOpenErr := os.Create(g.ConfigFilePath + g.ConfigFileName + ".toml")
	if fileOpenErr != nil {
		return fileOpenErr
	}

	configToWrite := g.ConfigUpdater.GetUpdatedConfig()

	encodingErr := toml.NewEncoder(file).Encode(configToWrite)
	if encodingErr !=  nil {
		return encodingErr
	}

	return file.Close()

	////TODO: get this from args
	//collateral := config.NewCollateral("ETH-B", "1.2.3")
	//flipContract := config.Contract{
	//	Address:  "0xabc123",
	//	Abi:      "flip abi",
	//	Deployed: 1,
	//}
	//
	//medianContract := config.Contract{
	//	Address:  "0xdef456",
	//	Abi:      "median abi",
	//	Deployed: 2,
	//}
	//
	//osmContract := config.Contract{
	//	Address:  "0xghi789",
	//	Abi:      "osm abi",
	//	Deployed: 3,
	//}
	//
	//contracts := make(map[string]config.Contract)
	//contracts["flip"] = flipContract
	//contracts["median"] = medianContract
	//contracts["osm"] = osmContract
	//
}

