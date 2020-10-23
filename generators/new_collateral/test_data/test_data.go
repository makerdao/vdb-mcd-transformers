package test_data

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
)

var (
	InitialConfig = types.TransformersConfig{
		ExporterMetadata: types.ExporterMetaData{
			Home:             "github.com/makerdao/vulcanizedb",
			Name:             "transformerExporter",
			Save:             false,
			Schema:           "maker",
			TransformerNames: []string{"cat_v1_1_0", "cat_file_vow"},
		},
		TransformerExporters: types.TransformerExporters{
			"cat_v1_1_0":       Cat110Exporter,
			"cat_file_vow":     CatFileVowExporter,
			"deny":             DenyExporter,
			"log_value":        LogValueExporter,
			"log_median_price": LogMedianPriceExporter,
		},
		Contracts: types.Contracts{
			"MCD_CAT_1_0_0": Cat100Contract,
			"MCD_CAT_1_1_0": Cat110Contract,
		},
	}

	EthBCollateral = types.NewCollateral("ETH-B", "1.1.3")
	PaxgCollateral = types.NewCollateral("PAXG", "1.1.4")

	Cat110Exporter = types.TransformerExporter{
		Path:       "transformers/storage/cat/v1_1_0/initializer",
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	Cat110ExporterMap = map[string]interface{}{
		"Path":       Cat110Exporter.Path,
		"Type":       Cat110Exporter.Type,
		"Repository": Cat110Exporter.Repository,
		"Migrations": Cat110Exporter.Migrations,
		"Rank":       Cat110Exporter.Rank,
		"Contracts":  nil,
	}

	DenyExporter = types.TransformerExporter{
		Path:       "transformers/events/auth/deny/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{FlipBatAContractName, FlipEthAContractName, MedianBatContractName, OsmBatContractName},
		Rank:       "0",
	}

	UpdatedDenyExporterMap = map[string]interface{}{
		"Path":       DenyExporter.Path,
		"Type":       DenyExporter.Type,
		"Repository": DenyExporter.Repository,
		"Migrations": DenyExporter.Migrations,
		"Contracts": []interface{}{
			FlipBatAContractName,
			FlipEthAContractName,
			MedianBatContractName,
			OsmBatContractName,
			FlipEthBContractName,
			MedianEthBContractName,
			OsmEthBContractName,
		},
		"Rank": "0",
	}

	LogMedianPriceExporter = types.TransformerExporter{
		Path:       "transformers/events/log_median_price/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{MedianBatContractName},
		Rank:       "0",
	}

	UpdatedLogMedianPriceExporterMap = map[string]interface{}{
		"Path":       LogMedianPriceExporter.Path,
		"Type":       LogMedianPriceExporter.Type,
		"Repository": LogMedianPriceExporter.Repository,
		"Migrations": LogMedianPriceExporter.Migrations,
		"Contracts":  []interface{}{MedianBatContractName, MedianEthBContractName},
		"Rank":       "0",
	}

	LogValueExporter = types.TransformerExporter{
		Path:       "transformers/events/log_value/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{OsmBatContractName},
		Rank:       "0",
	}

	UpdatedLogValueExporterMap = map[string]interface{}{
		"Path":       LogValueExporter.Path,
		"Type":       LogValueExporter.Type,
		"Repository": LogValueExporter.Repository,
		"Migrations": LogValueExporter.Migrations,
		"Contracts":  []interface{}{OsmBatContractName, OsmEthBContractName},
		"Rank":       LogValueExporter.Rank,
	}

	CatFileVowExporter = types.TransformerExporter{
		Path:       "transformers/events/cat_file/vow/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{"MCD_CAT_1.0.0", "MCD_CAT_1.1.0"},
		Rank:       "0",
	}

	CatFileVowExporterMap = map[string]interface{}{
		"Path":       CatFileVowExporter.Path,
		"Type":       CatFileVowExporter.Type,
		"Repository": CatFileVowExporter.Repository,
		"Migrations": CatFileVowExporter.Migrations,
		"Rank":       CatFileVowExporter.Rank,
		"Contracts": []interface{}{
			"MCD_CAT_1.0.0", "MCD_CAT_1.1.0",
		},
	}

	FlipEthBStorageExporter = types.TransformerExporter{
		Path:       "transformers/storage/flip/initializers/eth_b/v1_1_3",
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	FlipEthBStorageExporterMap = map[string]interface{}{
		"Path":       FlipEthBStorageExporter.Path,
		"Type":       FlipEthBStorageExporter.Type,
		"Repository": FlipEthBStorageExporter.Repository,
		"Migrations": FlipEthBStorageExporter.Migrations,
		"Rank":       FlipEthBStorageExporter.Rank,
		"Contracts":  nil,
	}

	MedianEthBStorageExporter = types.TransformerExporter{
		Path:       "transformers/storage/median/initializers/median_eth_b",
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	MedianEthBStorageExporterMap = map[string]interface{}{
		"Path":       MedianEthBStorageExporter.Path,
		"Type":       MedianEthBStorageExporter.Type,
		"Repository": MedianEthBStorageExporter.Repository,
		"Migrations": MedianEthBStorageExporter.Migrations,
		"Rank":       MedianEthBStorageExporter.Rank,
		"Contracts":  nil,
	}

	FlipBatAContractName  = "MCD_FLIP_BAT_A_1_0_0"
	FlipEthAContractName  = "MCD_FLIP_ETH_A_1_0_0"
	MedianBatContractName = "MEDIAN_BAT"
	OsmBatContractName    = "OSM_BAT"

	Cat100ContractName = "MCD_CAT_1_0_0"
	Cat100Contract     = types.Contract{
		Address:  "0x78f2c2af65126834c51822f56be0d7469d7a523e",
		Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"lump","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
		Deployed: 8928165,
	}

	Cat110ContractName = "MCD_CAT_1_1_0"
	Cat110Contract     = types.Contract{
		Address:  "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea",
		Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"box","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"rad","type":"uint256"}],"name":"claw","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"dunk","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"litter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
		Deployed: 10742907,
	}

	FlipEthBContractName = "MCD_FLIP_ETH_B_1_1_3"
	FlipEthBContract     = types.Contract{
		Address:  "testFlipContractAddress",
		Abi:      "testFlipContractAbi",
		Deployed: 123,
	}

	MedianEthBContractName = "MEDIAN_ETH_B"
	MedianEthBContract     = types.Contract{
		Address:  "testMedianContractAddress",
		Abi:      "testMedianContractAbi",
		Deployed: 456,
	}

	OsmEthBContractName = "OSM_ETH_B"
	OsmEthBContract     = types.Contract{
		Address:  "testOsmContractAddress",
		Abi:      "testOsmContractAbi",
		Deployed: 789,
	}

	EthBContracts = types.Contracts{
		"flip":   FlipEthBContract,
		"median": MedianEthBContract,
		"osm":    OsmEthBContract,
	}

	UpdatedConfig = types.TransformersConfig{
		ExporterMetadata: types.ExporterMetaData{
			Home:             "test-home",
			Name:             "test-config",
			Save:             false,
			Schema:           "test-schema",
			TransformerNames: []string{"transformer1", "transformer2"},
		},
		TransformerExporters: types.TransformerExporters{
			"test-1": types.TransformerExporter{
				Path:       "path-test-1",
				Type:       "eth_storage",
				Repository: "repo-1",
				Migrations: "test-migrations",
				Contracts:  nil,
				Rank:       "0",
			},
			"test-2": types.TransformerExporter{
				Path:       "path-test-2",
				Type:       "eth_event",
				Repository: "repo-2",
				Migrations: "test-migrations",
				Contracts: []string{
					"test-contract-1",
					"test-contract-2",
				},
				Rank: "0",
			},
		},
		Contracts: types.Contracts{
			"test-contract-1": types.Contract{
				Address:  "0xf185d0682d50819263941e5f4EacC763CC5C6C42",
				Abi:      `'[{"inputs":[{"internalType":"address","name":"src_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"val","type":"bytes32"}],"name":"LogValue","type":"event"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"bud","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"src_","type":"address"}],"name":"change","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address[]","name":"a","type":"address[]"}],"name":"diss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"a","type":"address"}],"name":"diss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"hop","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address[]","name":"a","type":"address[]"}],"name":"kiss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"a","type":"address"}],"name":"kiss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"pass","outputs":[{"internalType":"bool","name":"ok","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"peek","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"},{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"peep","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"},{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"poke","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"read","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"src","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"start","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint16","name":"ts","type":"uint16"}],"name":"step","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"stop","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"stopped","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"void","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"zzz","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]'`,
				Deployed: 9974841,
			},
			"test-contract-2": types.Contract{
				Address:  "0x7382c066801E7Acb2299aC8562847B9883f5CD3c",
				Abi:      `'[{"inputs":[{"internalType":"address","name":"src_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"val","type":"bytes32"}],"name":"LogValue","type":"event"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"bud","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"src_","type":"address"}],"name":"change","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address[]","name":"a","type":"address[]"}],"name":"diss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"a","type":"address"}],"name":"diss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"hop","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address[]","name":"a","type":"address[]"}],"name":"kiss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"a","type":"address"}],"name":"kiss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"pass","outputs":[{"internalType":"bool","name":"ok","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"peek","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"},{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"peep","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"},{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"poke","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"read","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"src","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"start","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint16","name":"ts","type":"uint16"}],"name":"step","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"stop","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"stopped","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"void","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"zzz","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]'`,
				Deployed: 10322969,
			},
		},
	}
	UpdatedConfigForToml = types.TransformersConfigForToml{
		Exporter: map[string]interface{}{
			"Home":   "github.com/makerdao/vulcanizedb",
			"Name":   "transformerExporter",
			"Save":   false,
			"Schema": "maker",
			"TransformerNames": []interface{}{
				"cat_v1_1_0",
				"cat_file_vow",
				"flip_eth_b_v1_1_3", // new storage flip transformer
				"median_eth_b",      // new median eth transformer
			},
			"median_eth_b":      MedianEthBStorageExporterMap,
			"flip_eth_b_v1_1_3": FlipEthBStorageExporterMap,
			"cat_file_vow":      CatFileVowExporterMap,
			"cat_v1_1_0":        Cat110ExporterMap,
			"log_median_price":  UpdatedLogMedianPriceExporterMap,
			"log_value":         UpdatedLogValueExporterMap,
			"deny":              UpdatedDenyExporterMap,
		},
		Contracts: types.Contracts{
			"MCD_CAT_1_0_0":        Cat100Contract,
			"MCD_CAT_1_1_0":        Cat110Contract,
			"MCD_FLIP_ETH_B_1_1_3": FlipEthBContract,
			"MEDIAN_ETH_B":         MedianEthBContract,
			"OSM_ETH_B":            OsmEthBContract,
		},
	}

	TestConfigFileContent = `[Exporter]
  Home = "github.com/makerdao/vulcanizedb"
  Name = "transformerExporter"
  Save = false
  Schema = "maker"
  TransformerNames = ["cat_v1_1_0", "cat_file_vow", "flip_eth_b_v1_1_3", "median_eth_b"]
  [Exporter.cat_file_vow]
    Contracts = ["MCD_CAT_1.0.0", "MCD_CAT_1.1.0"]
    Migrations = "db/migrations"
    Path = "transformers/events/cat_file/vow/initializer"
    Rank = "0"
    Repository = "github.com/makerdao/vdb-mcd-transformers"
    Type = "eth_event"
  [Exporter.cat_v1_1_0]
    Migrations = "db/migrations"
    Path = "transformers/storage/cat/v1_1_0/initializer"
    Rank = "0"
    Repository = "github.com/makerdao/vdb-mcd-transformers"
    Type = "eth_storage"
  [Exporter.deny]
    Contracts = ["MCD_FLIP_BAT_A_1_0_0", "MCD_FLIP_ETH_A_1_0_0", "MEDIAN_BAT", "OSM_BAT", "MCD_FLIP_ETH_B_1_1_3", "MEDIAN_ETH_B", "OSM_ETH_B"]
    Migrations = "db/migrations"
    Path = "transformers/events/auth/deny/initializer"
    Rank = "0"
    Repository = "github.com/makerdao/vdb-mcd-transformers"
    Type = "eth_event"
  [Exporter.flip_eth_b_v1_1_3]
    Migrations = "db/migrations"
    Path = "transformers/storage/flip/initializers/eth_b/v1_1_3"
    Rank = "0"
    Repository = "github.com/makerdao/vdb-mcd-transformers"
    Type = "eth_storage"
  [Exporter.log_median_price]
    Contracts = ["MEDIAN_BAT", "MEDIAN_ETH_B"]
    Migrations = "db/migrations"
    Path = "transformers/events/log_median_price/initializer"
    Rank = "0"
    Repository = "github.com/makerdao/vdb-mcd-transformers"
    Type = "eth_event"
  [Exporter.log_value]
    Contracts = ["OSM_BAT", "OSM_ETH_B"]
    Migrations = "db/migrations"
    Path = "transformers/events/log_value/initializer"
    Rank = "0"
    Repository = "github.com/makerdao/vdb-mcd-transformers"
    Type = "eth_event"
  [Exporter.median_eth_b]
    Migrations = "db/migrations"
    Path = "transformers/storage/median/initializers/median_eth_b"
    Rank = "0"
    Repository = "github.com/makerdao/vdb-mcd-transformers"
    Type = "eth_storage"

[contract]
  [contract.MCD_CAT_1_0_0]
    address = "0x78f2c2af65126834c51822f56be0d7469d7a523e"
    abi = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vat_\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"urn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"art\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tab\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"flip\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"Bite\",\"type\":\"event\"},{\"anonymous\":true,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes4\",\"name\":\"sig\",\"type\":\"bytes4\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"usr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"arg1\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"arg2\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"LogNote\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"urn\",\"type\":\"address\"}],\"name\":\"bite\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"cage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"usr\",\"type\":\"address\"}],\"name\":\"deny\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"what\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"file\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"what\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"data\",\"type\":\"address\"}],\"name\":\"file\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"what\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"flip\",\"type\":\"address\"}],\"name\":\"file\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ilks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"flip\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chop\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lump\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"live\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"usr\",\"type\":\"address\"}],\"name\":\"rely\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"vat\",\"outputs\":[{\"internalType\":\"contract VatLike\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"vow\",\"outputs\":[{\"internalType\":\"contract VowLike\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"wards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
    deployed = 8928165
  [contract.MCD_CAT_1_1_0]
    address = "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea"
    abi = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vat_\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"urn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"art\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tab\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"flip\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"Bite\",\"type\":\"event\"},{\"anonymous\":true,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes4\",\"name\":\"sig\",\"type\":\"bytes4\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"usr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"arg1\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"arg2\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"LogNote\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"urn\",\"type\":\"address\"}],\"name\":\"bite\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"box\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"cage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rad\",\"type\":\"uint256\"}],\"name\":\"claw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"usr\",\"type\":\"address\"}],\"name\":\"deny\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"what\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"file\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"what\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"file\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"what\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"data\",\"type\":\"address\"}],\"name\":\"file\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ilk\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"what\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"flip\",\"type\":\"address\"}],\"name\":\"file\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ilks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"flip\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chop\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dunk\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"litter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"live\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"usr\",\"type\":\"address\"}],\"name\":\"rely\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"vat\",\"outputs\":[{\"internalType\":\"contract VatLike\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"vow\",\"outputs\":[{\"internalType\":\"contract VowLike\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"wards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
    deployed = 10742907
  [contract.MCD_FLIP_ETH_B_1_1_3]
    address = "testFlipContractAddress"
    abi = "testFlipContractAbi"
    deployed = 123
  [contract.MEDIAN_ETH_B]
    address = "testMedianContractAddress"
    abi = "testMedianContractAbi"
    deployed = 456
  [contract.OSM_ETH_B]
    address = "testOsmContractAddress"
    abi = "testOsmContractAbi"
    deployed = 789
`
)
