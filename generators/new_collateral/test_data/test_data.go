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
		"path":       Cat110Exporter.Path,
		"type":       Cat110Exporter.Type,
		"repository": Cat110Exporter.Repository,
		"migrations": Cat110Exporter.Migrations,
		"rank":       Cat110Exporter.Rank,
		"contracts":  nil,
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
		"path":       DenyExporter.Path,
		"type":       DenyExporter.Type,
		"repository": DenyExporter.Repository,
		"migrations": DenyExporter.Migrations,
		"contracts": []interface{}{
			FlipBatAContractName,
			FlipEthAContractName,
			MedianBatContractName,
			OsmBatContractName,
			FlipEthBContractName,
			MedianEthBContractName,
			OsmEthBContractName,
		},
		"rank": "0",
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
		"path":       LogMedianPriceExporter.Path,
		"type":       LogMedianPriceExporter.Type,
		"repository": LogMedianPriceExporter.Repository,
		"migrations": LogMedianPriceExporter.Migrations,
		"contracts":  []interface{}{MedianBatContractName, MedianEthBContractName},
		"rank":       "0",
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
		"path":       LogValueExporter.Path,
		"type":       LogValueExporter.Type,
		"repository": LogValueExporter.Repository,
		"migrations": LogValueExporter.Migrations,
		"contracts":  []interface{}{OsmBatContractName, OsmEthBContractName},
		"rank":       LogValueExporter.Rank,
	}

	CatFileVowExporter = types.TransformerExporter{
		Path:       "transformers/events/cat_file/vow/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{"MCD_CAT_1_0_0", "MCD_CAT_1_1_0"},
		Rank:       "0",
	}

	CatFileVowExporterMap = map[string]interface{}{
		"path":       CatFileVowExporter.Path,
		"type":       CatFileVowExporter.Type,
		"repository": CatFileVowExporter.Repository,
		"migrations": CatFileVowExporter.Migrations,
		"rank":       CatFileVowExporter.Rank,
		"contracts": []interface{}{
			"MCD_CAT_1_0_0", "MCD_CAT_1_1_0",
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
		"path":       FlipEthBStorageExporter.Path,
		"type":       FlipEthBStorageExporter.Type,
		"repository": FlipEthBStorageExporter.Repository,
		"migrations": FlipEthBStorageExporter.Migrations,
		"rank":       FlipEthBStorageExporter.Rank,
		"contracts":  nil,
	}

	MedianEthBStorageExporter = types.TransformerExporter{
		Path:       "transformers/storage/median/initializers/median_eth_b/v1_1_3",
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	MedianEthBStorageExporterMap = map[string]interface{}{
		"path":       MedianEthBStorageExporter.Path,
		"type":       MedianEthBStorageExporter.Type,
		"repository": MedianEthBStorageExporter.Repository,
		"migrations": MedianEthBStorageExporter.Migrations,
		"rank":       MedianEthBStorageExporter.Rank,
		"contracts":  nil,
	}

	FlipBatAContractName  = "MCD_FLIP_BAT_A_1_0_0"
	FlipEthAContractName  = "MCD_FLIP_ETH_A_1_0_0"
	MedianBatContractName = "MEDIAN_BAT_1_0_0"
	OsmBatContractName    = "OSM_BAT"

	Cat100ContractName = "MCD_CAT_1_0_0"
	Cat100Contract     = types.Contract{
		Address:  "0x78f2c2af65126834c51822f56be0d7469d7a523e",
		Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}]}]`,
		Deployed: 8928165,
	}

	Cat110ContractName = "MCD_CAT_1_1_0"
	Cat110Contract     = types.Contract{
		Address:  "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea",
		Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}]}]`,
		Deployed: 10742907,
	}

	FlipEthBContractName = "MCD_FLIP_ETH_B_1_1_3"
	FlipEthBContract     = types.Contract{
		Address:  "testFlipContractAddress",
		Abi:      "testFlipContractAbi",
		Deployed: 123,
	}

	MedianEthBContractName = "MEDIAN_ETH_B_1_1_3"
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
				Abi:      `'[{"inputs":[{"internalType":"address","name":"src_","type":"address"}]}]'`,
				Deployed: 9974841,
			},
			"test-contract-2": types.Contract{
				Address:  "0x7382c066801E7Acb2299aC8562847B9883f5CD3c",
				Abi:      `'[{"inputs":[{"internalType":"address","name":"src_","type":"address"}]}]'`,
				Deployed: 10322969,
			},
		},
	}
	UpdatedConfigForToml = types.TransformersConfigForToml{
		Exporter: map[string]interface{}{
			"home":   "github.com/makerdao/vulcanizedb",
			"name":   "transformerExporter",
			"save":   false,
			"schema": "maker",
			"transformerNames": []interface{}{
				"cat_v1_1_0",
				"cat_file_vow",
				"flip_eth_b_v1_1_3",   // new storage flip transformer
				"median_eth_b_v1_1_3", // new median eth transformer
			},
			"median_eth_b_v1_1_3": MedianEthBStorageExporterMap,
			"flip_eth_b_v1_1_3":   FlipEthBStorageExporterMap,
			"cat_file_vow":        CatFileVowExporterMap,
			"cat_v1_1_0":          Cat110ExporterMap,
			"log_median_price":    UpdatedLogMedianPriceExporterMap,
			"log_value":           UpdatedLogValueExporterMap,
			"deny":                UpdatedDenyExporterMap,
		},
		Contracts: types.Contracts{
			"MCD_CAT_1_0_0":        Cat100Contract,
			"MCD_CAT_1_1_0":        Cat110Contract,
			"MCD_FLIP_ETH_B_1_1_3": FlipEthBContract,
			"MEDIAN_ETH_B_1_1_3":   MedianEthBContract,
			"OSM_ETH_B":            OsmEthBContract,
		},
	}

	TestConfigFileContent = `[exporter]
  home = "github.com/makerdao/vulcanizedb"
  name = "transformerExporter"
  save = false
  schema = "maker"
  transformerNames = ["cat_v1_1_0", "cat_file_vow", "flip_eth_b_v1_1_3", "median_eth_b_v1_1_3"]
  [exporter.cat_file_vow]
    contracts = ["MCD_CAT_1_0_0", "MCD_CAT_1_1_0"]
    migrations = "db/migrations"
    path = "transformers/events/cat_file/vow/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"
  [exporter.cat_v1_1_0]
    migrations = "db/migrations"
    path = "transformers/storage/cat/v1_1_0/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_storage"
  [exporter.deny]
    contracts = ["MCD_FLIP_BAT_A_1_0_0", "MCD_FLIP_ETH_A_1_0_0", "MEDIAN_BAT_1_0_0", "OSM_BAT", "MCD_FLIP_ETH_B_1_1_3", "MEDIAN_ETH_B_1_1_3", "OSM_ETH_B"]
    migrations = "db/migrations"
    path = "transformers/events/auth/deny/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"
  [exporter.flip_eth_b_v1_1_3]
    migrations = "db/migrations"
    path = "transformers/storage/flip/initializers/eth_b/v1_1_3"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_storage"
  [exporter.log_median_price]
    contracts = ["MEDIAN_BAT_1_0_0", "MEDIAN_ETH_B_1_1_3"]
    migrations = "db/migrations"
    path = "transformers/events/log_median_price/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"
  [exporter.log_value]
    contracts = ["OSM_BAT", "OSM_ETH_B"]
    migrations = "db/migrations"
    path = "transformers/events/log_value/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"
  [exporter.median_eth_b_v1_1_3]
    migrations = "db/migrations"
    path = "transformers/storage/median/initializers/median_eth_b/v1_1_3"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_storage"

[contract]
  [contract.MCD_CAT_1_0_0]
    address = "0x78f2c2af65126834c51822f56be0d7469d7a523e"
    abi = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vat_\",\"type\":\"address\"}]}]"
    deployed = 8928165
  [contract.MCD_CAT_1_1_0]
    address = "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea"
    abi = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vat_\",\"type\":\"address\"}]}]"
    deployed = 10742907
  [contract.MCD_FLIP_ETH_B_1_1_3]
    address = "testFlipContractAddress"
    abi = "testFlipContractAbi"
    deployed = 123
  [contract.MEDIAN_ETH_B_1_1_3]
    address = "testMedianContractAddress"
    abi = "testMedianContractAbi"
    deployed = 456
  [contract.OSM_ETH_B]
    address = "testOsmContractAddress"
    abi = "testOsmContractAbi"
    deployed = 789
`
)
