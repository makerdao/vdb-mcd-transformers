// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package integration_tests

import (
	"plugin"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	p2 "github.com/vulcanize/vulcanizedb/pkg/plugin"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
)

var eventConfig = config.Plugin{
	Home: "github.com/vulcanize/mcd_transformers",
	Transformers: map[string]config.Transformer{
		"bite": {
			Path:           "transformers/events/bite/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
		"cat_file": {
			Path:           "transformers/events/cat_file/flip/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
		"deal": {
			Path:           "transformers/events/deal/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
	},
	FileName: "testEventTransformerSet",
	FilePath: "$GOPATH/src/github.com/vulcanize/mcd_transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var storageConfig = config.Plugin{
	Home: "github.com/vulcanize/mcd_transformers",
	Transformers: map[string]config.Transformer{
		"jug": {
			Path:           "transformers/storage/jug/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
		"vat": {
			Path:           "transformers/storage/vat/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
	},
	FileName: "testStorageTransformerSet",
	FilePath: "$GOPATH/src/github.com/vulcanize/mcd_transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var combinedConfig = config.Plugin{
	Home: "github.com/vulcanize/mcd_transformers",
	Transformers: map[string]config.Transformer{
		"bite": {
			Path:           "transformers/events/bite/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
		"cat_file": {
			Path:           "transformers/events/cat_file/flip/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
		"deal": {
			Path:           "transformers/events/deal/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
		"jug": {
			Path:           "transformers/storage/jug/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
		"vat": {
			Path:           "transformers/storage/vat/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/vulcanize/mcd_transformers",
		},
	},
	FileName: "testComboTransformerSet",
	FilePath: "$GOPATH/src/github.com/vulcanize/mcd_transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var dbConfig = config.Database{
	Hostname: "localhost",
	Port:     5432,
	Name:     "vulcanize_testing",
}

type Exporter interface {
	Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer, []transformer.ContractTransformerInitializer)
}

func SetupDBandBC() (*postgres.DB, core.BlockChain) {
	rpcClient, ethClient, err := getClients(ipc)
	Expect(err).NotTo(HaveOccurred())
	bc, err := getBlockChain(rpcClient, ethClient)
	Expect(err).NotTo(HaveOccurred())
	db := test_config.NewTestDB(bc.Node())
	test_config.CleanTestDB(db)
	return db, bc
}

var _ = Describe("Plugin test", func() {
	var g p2.Generator
	var goPath, soPath string
	var err error
	var db *postgres.DB
	var hr repositories.HeaderRepository
	var headerID int64
	viper.SetConfigName("testing")
	viper.AddConfigPath("$GOPATH/src/github.com/vulcanize/mcd_transformers/environments/")

	Describe("Event Transformers only", func() {
		BeforeEach(func() {
			goPath, soPath, err = eventConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(eventConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified  TransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventTransformerInitializers, storageTransformerInitializers, _ := exporter.Export()
				Expect(len(eventTransformerInitializers)).To(Equal(3))
				Expect(len(storageTransformerInitializers)).To(Equal(0))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of TransformerInitializers that we can execute over", func() {
				db, bc := SetupDBandBC()
				hr = repositories.NewHeaderRepository(db)
				header1, err := bc.GetHeaderByNumber(12742101)
				Expect(err).ToNot(HaveOccurred())
				headerID, err = hr.CreateOrUpdateHeader(header1)
				Expect(err).ToNot(HaveOccurred())

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventTransformerInitializers, _, _ := exporter.Export()

				w := watcher.NewEventWatcher(db, bc)
				w.AddTransformers(eventTransformerInitializers)
				err = w.Execute(constants.HeaderMissing)
				Expect(err).ToNot(HaveOccurred())

				type model struct {
					Ilk              string `db:"ilk_id"`
					What             string
					Flip             string
					LogIndex         uint   `db:"log_idx"`
					TransactionIndex uint   `db:"tx_idx"`
					Raw              []byte `db:"raw_log"`
					Id               int64  `db:"id"`
					HeaderId         int64  `db:"header_id"`
				}

				returned := model{}

				err = db.Get(&returned, `SELECT * FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
				Expect(err).ToNot(HaveOccurred())

				ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
				Expect(err).NotTo(HaveOccurred())
				Expect(returned.Ilk).To(Equal(strconv.Itoa(ilkID)))
				Expect(returned.What).To(Equal("flip"))
				Expect(returned.Flip).To(Equal("0xe8b1439bbCc5a4Da4A8e938257AEa3a8D18F78ad"))
				Expect(returned.TransactionIndex).To(Equal(uint(0)))
				Expect(returned.LogIndex).To(Equal(uint(3)))
			})

			It("rechecks checked headers for event logs", func() {
				db, bc := SetupDBandBC()
				hr = repositories.NewHeaderRepository(db)
				header1, err := bc.GetHeaderByNumber(12742101)
				Expect(err).ToNot(HaveOccurred())
				headerID, err = hr.CreateOrUpdateHeader(header1)
				Expect(err).ToNot(HaveOccurred())

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventTransformerInitializers, _, _ := exporter.Export()

				w := watcher.NewEventWatcher(db, bc)
				w.AddTransformers(eventTransformerInitializers)
				err = w.Execute(constants.HeaderMissing)
				Expect(err).ToNot(HaveOccurred())

				err = w.Execute(constants.HeaderMissing)
				Expect(err).ToNot(HaveOccurred())

				var catFileFlipChecked int
				err = db.Get(&catFileFlipChecked, `SELECT cat_file_flip FROM public.checked_headers WHERE header_id = $1`, headerID)
				Expect(err).ToNot(HaveOccurred())
				Expect(catFileFlipChecked).To(Equal(2))
			})
		})
	})

	Describe("Storage Transformers only", func() {
		BeforeEach(func() {
			goPath, soPath, err = storageConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(storageConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified StorageTransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventTransformerInitializers, storageTransformerInitializers, _ := exporter.Export()
				Expect(len(storageTransformerInitializers)).To(Equal(2))
				Expect(len(eventTransformerInitializers)).To(Equal(0))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of StorageTransformerInitializers that we can execute over", func() {
				db, _ = SetupDBandBC()
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				_, storageTransformerInitializers, _ := exporter.Export()

				tailer := fs.FileTailer{Path: viper.GetString("filesystem.storageDiffsPath")}
				storageFetcher := fetcher.NewCsvTailStorageFetcher(tailer)
				w := watcher.NewStorageWatcher(storageFetcher, db)
				w.AddTransformers(storageTransformerInitializers)
				// This blocks right now, need to make test file to read from
				//err = w.Execute()
				//Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Event and Storage Transformers in same instance", func() {
		BeforeEach(func() {
			goPath, soPath, err = combinedConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(combinedConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified TransformerInitializers and StorageTransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventInitializers, storageInitializers, _ := exporter.Export()
				Expect(len(eventInitializers)).To(Equal(3))
				Expect(len(storageInitializers)).To(Equal(2))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of TransformerInitializers and StorageTransformerInitializers that we can execute over", func() {
				db, bc := SetupDBandBC()
				hr = repositories.NewHeaderRepository(db)
				header1, err := bc.GetHeaderByNumber(12742101)
				Expect(err).ToNot(HaveOccurred())
				headerID, err = hr.CreateOrUpdateHeader(header1)
				Expect(err).ToNot(HaveOccurred())

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventInitializers, storageInitializers, _ := exporter.Export()

				ew := watcher.NewEventWatcher(db, bc)
				ew.AddTransformers(eventInitializers)
				err = ew.Execute(constants.HeaderMissing)
				Expect(err).ToNot(HaveOccurred())

				type model struct {
					Ilk              string `db:"ilk_id"`
					What             string
					Flip             string
					LogIndex         uint   `db:"log_idx"`
					TransactionIndex uint   `db:"tx_idx"`
					Raw              []byte `db:"raw_log"`
					Id               int64  `db:"id"`
					HeaderId         int64  `db:"header_id"`
				}

				returned := model{}

				err = db.Get(&returned, `SELECT * FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
				Expect(err).ToNot(HaveOccurred())

				ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
				Expect(err).NotTo(HaveOccurred())
				Expect(returned.Ilk).To(Equal(strconv.Itoa(ilkID)))
				Expect(returned.What).To(Equal("flip"))
				Expect(returned.Flip).To(Equal("0xe8b1439bbCc5a4Da4A8e938257AEa3a8D18F78ad"))
				Expect(returned.TransactionIndex).To(Equal(uint(0)))
				Expect(returned.LogIndex).To(Equal(uint(3)))

				tailer := fs.FileTailer{Path: viper.GetString("filesystem.storageDiffsPath")}
				storageFetcher := fetcher.NewCsvTailStorageFetcher(tailer)
				sw := watcher.NewStorageWatcher(storageFetcher, db)
				sw.AddTransformers(storageInitializers)
				// This blocks right now, need to make test file to read from
				//err = w.Execute()
				//Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
