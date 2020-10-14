package generator_test

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("New Collateral Generator", func() {
	Context("ParseCurrentConfigFile", func() {
		var testConfig = `
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
[servers]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [servers.alpha]
	  ip = "10.0.0.1"
	  dc = "eqdc10"

  [servers.beta_1_0]
	  ip = "10.0.0.2"
	  dc = "eqdc10"
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

		It("reads in the exporter metadata", func() {
			expectedExporterMetadata := generator.ExporterMetaData{
				Home:             "github.com/makerdao/vulcanizedb",
				Name:             "transformerExporter",
				Save:             false,
				Schema:           "maker",
				TransformerNames: []string{"cat_v1_1_0", "cat_file_vow"},
			}

			configStruct, parseErr := generator.ParseCurrentConfig(testConfig)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(configStruct.ExporterMetadata).To(Equal(expectedExporterMetadata))
		})

		It("reads in the exporterTransformers", func() {
			expectedTransformerExporters := generator.TransformerExporters{
				"exporter.cat_v1_1_0": generator.TransformerExporter{
					Path:       "transformers/storage/cat/v1_1_0/initializer",
					Kind:       "eth_storage",
					Repository: "github.com/makerdao/vdb-mcd-transformers",
					Migrations: "db/migrations",
					Rank:       "0",
				},
				"exporter.cat_file_vow": generator.TransformerExporter{
					Path:       "transformers/events/cat_file/vow/initializer",
					Kind:       "eth_event",
					Repository: "github.com/makerdao/vdb-mcd-transformers",
					Migrations: "db/migrations",
					Contracts:  []string{"MCD_CAT_1.0.0", "MCD_CAT_1.1.0"},
					Rank:       "0",
				},
			}
			configStruct, parseErr := generator.ParseCurrentConfig(testConfig)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(configStruct.TransformerExporters).To(Equal(expectedTransformerExporters))
		})

		It("reads in the contracts", func() {
			expectedContracts := generator.Contracts{
				"MCD_CAT_1_0_0": generator.Contract{
					Address:  "0x78f2c2af65126834c51822f56be0d7469d7a523e",
					Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"lump","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
					Deployed: 8928165,
				},
				"MCD_CAT_1_1_0": generator.Contract{
					Address:  "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea",
					Abi:      `[{"inputs":[{"internalType":"address","name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"ilk","type":"bytes32"},{"indexed":true,"internalType":"address","name":"urn","type":"address"},{"indexed":false,"internalType":"uint256","name":"ink","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"art","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":false,"internalType":"address","name":"flip","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"}],"name":"Bite","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"address","name":"urn","type":"address"}],"name":"bite","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"box","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"cage","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"rad","type":"uint256"}],"name":"claw","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"ilk","type":"bytes32"},{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"address","name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"ilks","outputs":[{"internalType":"address","name":"flip","type":"address"},{"internalType":"uint256","name":"chop","type":"uint256"},{"internalType":"uint256","name":"dunk","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"litter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"live","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"internalType":"contract VowLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
					Deployed: 10742907,
				},
			}
			configStruct, parseErr := generator.ParseCurrentConfig(testConfig)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(configStruct.Contracts).To(Equal(expectedContracts))
		})
	})
})
