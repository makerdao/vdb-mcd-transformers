package test_data

import "github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"

var (
	TestConfigString = `
[exporter]
    home     = "github.com/makerdao/vulcanizedb"
    name     = "transformerExporter"
    save     = false
    schema   = "maker"
    transformerNames = [
        "cat_v1_1_0",
		"cat_file_vow"
	]
    [exporter.cat_v1_1_0]
        path = "transformers/storage/cat/v1_1_0/initializer"
        type = "eth_storage"
        repository = "github.com/makerdao/vdb-mcd-transformers"
        migrations = "db/migrations"
        rank = "0"
    [exporter.cat_file_vow]
        path = "transformers/events/cat_file/vow/initializer"
        type = "eth_event"
        repository = "github.com/makerdao/vdb-mcd-transformers"
        migrations = "db/migrations"
        contracts = [
            "MCD_CAT_1.0.0",
            "MCD_CAT_1.1.0"
        ]
        rank = "0"
[contract]
    [contract.MCD_CAT_1_0_0]
        address  = "0x78f2c2af65126834c51822f56be0d7469d7a523e"
        abi      = '[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"lump","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]'
        deployed = 8928165
    [contract.MCD_CAT_1_1_0]
        address  = "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea"
        abi      = '[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"box","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"rad","type":"uint256"}],"name":"claw","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"dunk","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"litter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]'
        deployed = 10742907
`

	InitialConfig = config.TransformersConfig{
		ExporterMetadata:     ExporterMetadata,
		TransformerExporters: config.TransformerExporters{
			"exporter.cat_v1_1_0":   Cat110Exporter,
			"exporter.cat_file_vow": CatFileVowExporter,
			"exporter.deny": DenyExporter,
			"exporter.log_value": LogValueExporter,
			"exporter.log_median_price": LogMedianPriceExporter,
		},
		Contracts: config.Contracts{
			"MCD_CAT_1_0_0": Cat100Contract,
			"MCD_CAT_1_1_0": Cat110Contract,
		},
	}

	EthBCollateral = config.NewCollateral("ETH-B", "1.1.3")
	PaxgCollateral = config.NewCollateral("PAXG","1.1.4")

	ExporterMetadata = config.ExporterMetaData{
		Home:             "github.com/makerdao/vulcanizedb",
		Name:             "transformerExporter",
		Save:             false,
		Schema:           "maker",
		TransformerNames: []string{"cat_v1_1_0", "cat_file_vow"},
	}

	Cat110Exporter = config.TransformerExporter{
		Path:       "transformers/storage/cat/v1_1_0/initializer",
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	DenyExporter = config.TransformerExporter{
		Path:       "transformers/events/auth/deny_initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{"MCD_FLIP_BAT_A_1_0_0", "MEDIAN_BAT", "OSM_BAT"},
		Rank:       "0",
	}

	LogMedianPriceExporter = config.TransformerExporter{
		Path:       "transformers/events/log_median_price/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{"MEDIAN_BAT"},
		Rank:       "0",
	}

	LogValueExporter = config.TransformerExporter{
		Path:       "transformers/events/log_value/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{"OSM_BAT"},
		Rank:       "0",
	}

	CatFileVowExporter = config.TransformerExporter{
		Path:       "transformers/events/cat_file/vow/initializer",
		Type:       "eth_event",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Contracts:  []string{"MCD_CAT_1.0.0", "MCD_CAT_1.1.0"},
		Rank:       "0",
	}

	FlipEthBStorageExporter = config.TransformerExporter{
		Path:       "transformers/storage/flip/initializers/eth_b/v1_1_3",
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	MedianEthBStorageExporter = config.TransformerExporter{
		Path:       "transformers/storage/median/initializers/median_eth_b",
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	Cat100Contract = config.Contract{
		Address:  "0x78f2c2af65126834c51822f56be0d7469d7a523e",
		Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"lump","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
		Deployed: 8928165,
	}

	Cat110Contract = config.Contract{
		Address:  "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea",
		Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"box","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"rad","type":"uint256"}],"name":"claw","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"dunk","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"litter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
		Deployed: 10742907,
	}

	FlipEthBContractName = "MCD_FLIP_ETH_B_1_1_3"
	FlipEthBContract = config.Contract{
		Address:  "testFlipContractAddress",
		Abi:      "testFlipContractAbi",
		Deployed: 123,
	}

	MedianEthBContractName = "MEDIAN_ETH_B"
	MedianEthBContract = config.Contract{
		Address:  "testMedianContractAddress",
		Abi:      "testMedianContractAbi",
		Deployed: 456,
	}

	OsmEthBContractName = "OSM_ETH_B"
	OsmEthBContract = config.Contract{
		Address:  "testOsmContractAddress",
		Abi:      "testOsmContractAbi",
		Deployed: 789,
	}

	EthBContracts = config.Contracts{
		"flip": FlipEthBContract,
		"median": MedianEthBContract,
		"osm": OsmEthBContract,
	}
)
