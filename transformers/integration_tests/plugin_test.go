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
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/logs"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	p2 "github.com/makerdao/vulcanizedb/pkg/plugin"
	"github.com/makerdao/vulcanizedb/pkg/plugin/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var eventConfig = config.Plugin{
	Home: "github.com/makerdao/vdb-mcd-transformers",
	Transformers: map[string]config.Transformer{
		"bite": {
			Path:           "transformers/events/bite/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
		"cat_file": {
			Path:           "transformers/events/cat_file/flip/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
		"deal": {
			Path:           "transformers/events/deal/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
	},
	FileName: "testEventTransformerSet",
	FilePath: "$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var storageConfig = config.Plugin{
	Home: "github.com/makerdao/vdb-mcd-transformers",
	Transformers: map[string]config.Transformer{
		"jug": {
			Path:           "transformers/storage/jug/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
		"vat": {
			Path:           "transformers/storage/vat/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
	},
	FileName: "testStorageTransformerSet",
	FilePath: "$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var combinedConfig = config.Plugin{
	Home: "github.com/makerdao/vdb-mcd-transformers",
	Transformers: map[string]config.Transformer{
		"bite": {
			Path:           "transformers/events/bite/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
		"cat_file": {
			Path:           "transformers/events/cat_file/flip/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
		"deal": {
			Path:           "transformers/events/deal/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
		"jug": {
			Path:           "transformers/storage/jug/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
		"vat": {
			Path:           "transformers/storage/vat/initializer",
			Type:           config.EthStorage,
			MigrationPath:  "db/migrations",
			RepositoryPath: "github.com/makerdao/vdb-mcd-transformers",
		},
	},
	FileName: "testComboTransformerSet",
	FilePath: "$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var dbConfig = config.Database{
	Hostname: "localhost",
	Port:     5432,
	Name:     "vulcanize_testing",
}

type Exporter interface {
	Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []transformer.ContractTransformerInitializer)
}

var _ = Describe("Plugin test", func() {
	viper.SetConfigName("testing")
	viper.AddConfigPath("$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/environments/")

	var (
		g                            p2.Generator
		goPath, soPath               string
		hr                           datastore.HeaderRepository
		headerID                     int64
		ilk                          = "0x4554482d41000000000000000000000000000000000000000000000000000000"
		blockNumber                  = int64(8928180) //needs a mainnet block with a cat file flip
		maxConsecutiveUnexpectedErrs = 0
		retryInterval                = 2 * time.Second
		delegator                    logs.ILogDelegator
		extractor                    logs.ILogExtractor
		statusWriter                 fakes.MockStatusWriter
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	Describe("Event Transformers only", func() {
		BeforeEach(func() {
			var pathErr, initErr, generateErr error
			goPath, soPath, pathErr = eventConfig.GetPluginPaths()
			Expect(pathErr).ToNot(HaveOccurred())
			g, initErr = p2.NewGenerator(eventConfig, dbConfig)
			Expect(initErr).ToNot(HaveOccurred())
			generateErr = g.GenerateExporterPlugin()
			Expect(generateErr).ToNot(HaveOccurred())
			extractor = logs.NewLogExtractor(db, blockChain)
			delegator = logs.NewLogDelegator(db)
			statusWriter = fakes.MockStatusWriter{}
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
				hr = repositories.NewHeaderRepository(db)
				header1, err := blockChain.GetHeaderByNumber(blockNumber)
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

				w := watcher.NewEventWatcher(db, blockChain, extractor, delegator, maxConsecutiveUnexpectedErrs, retryInterval, &statusWriter)
				addErr := w.AddTransformers(eventTransformerInitializers)
				Expect(addErr).NotTo(HaveOccurred())

				var executeErr error
				go func() {
					executeErr = w.Execute(constants.HeaderUnchecked)
				}()

				Consistently(func() error {
					return executeErr
				}).Should(BeNil())
				expectedIlkID, err := shared.GetOrCreateIlk(ilk, db)
				Expect(err).NotTo(HaveOccurred())
				// including longer timeout because this test takes awhile to populate the db
				Eventually(func() int64 {
					var ilkID int64
					_ = db.Get(&ilkID, `SELECT ilk_id FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
					return ilkID
				}, time.Second*30).Should(Equal(expectedIlkID))
				Eventually(func() string {
					var what string
					_ = db.Get(&what, `SELECT what FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
					return what
				}).Should(Equal("flip"))
				Eventually(func() string {
					var flip string
					_ = db.Get(&flip, `SELECT flip FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
					return flip
				}).Should(Equal(test_data.FlipEthV100Address()))
			})

			It("rechecks checked headers for event logs", func() {
				hr = repositories.NewHeaderRepository(db)
				header1, err := blockChain.GetHeaderByNumber(blockNumber)
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

				w := watcher.NewEventWatcher(db, blockChain, extractor, delegator, maxConsecutiveUnexpectedErrs, retryInterval, &statusWriter)
				addErr := w.AddTransformers(eventTransformerInitializers)
				Expect(addErr).NotTo(HaveOccurred())
				var executeErrOne, executeErrTwo error

				go func() {
					executeErrOne = w.Execute(constants.HeaderUnchecked)
					executeErrTwo = w.Execute(constants.HeaderUnchecked)
				}()

				Consistently(func() error {
					return executeErrOne
				}).Should(BeNil())
				Consistently(func() error {
					return executeErrTwo
				}).Should(BeNil())
			})
		})
	})

	Describe("Storage Transformers only", func() {
		BeforeEach(func() {
			var pathErr, initErr, generateErr error
			goPath, soPath, pathErr = storageConfig.GetPluginPaths()
			Expect(pathErr).ToNot(HaveOccurred())
			g, initErr = p2.NewGenerator(storageConfig, dbConfig)
			Expect(initErr).ToNot(HaveOccurred())
			generateErr = g.GenerateExporterPlugin()
			Expect(generateErr).ToNot(HaveOccurred())
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
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				_, storageTransformerInitializers, _ := exporter.Export()

				w := watcher.NewStorageWatcher(db, -1, &statusWriter)
				w.AddTransformers(storageTransformerInitializers)
				// This blocks right now, need to make test file to read from
				//err = w.Execute()
				//Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Event and Storage Transformers in same instance", func() {
		BeforeEach(func() {
			var pathErr, initErr, generateErr error
			goPath, soPath, pathErr = combinedConfig.GetPluginPaths()
			Expect(pathErr).ToNot(HaveOccurred())
			g, initErr = p2.NewGenerator(combinedConfig, dbConfig)
			Expect(initErr).ToNot(HaveOccurred())
			generateErr = g.GenerateExporterPlugin()
			Expect(generateErr).ToNot(HaveOccurred())
			extractor = logs.NewLogExtractor(db, blockChain)
			delegator = logs.NewLogDelegator(db)
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
				hr = repositories.NewHeaderRepository(db)
				header1, err := blockChain.GetHeaderByNumber(blockNumber)
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

				ew := watcher.NewEventWatcher(db, blockChain, extractor, delegator, maxConsecutiveUnexpectedErrs, retryInterval, &statusWriter)
				addTransformersErr := ew.AddTransformers(eventInitializers)
				Expect(addTransformersErr).NotTo(HaveOccurred())

				var executeErr error
				go func() {
					executeErr = ew.Execute(constants.HeaderUnchecked)
				}()

				Consistently(func() error {
					return executeErr
				}).Should(BeNil())
				expectedIlkID, err := shared.GetOrCreateIlk(ilk, db)
				Expect(err).NotTo(HaveOccurred())
				// including longer timeout because this test takes awhile to populate the db
				Eventually(func() int64 {
					var ilkID int64
					_ = db.Get(&ilkID, `SELECT ilk_id FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
					return ilkID
				}, time.Second*30).Should(Equal(expectedIlkID))
				Eventually(func() string {
					var what string
					_ = db.Get(&what, `SELECT what FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
					return what
				}).Should(Equal("flip"))
				Eventually(func() string {
					var flip string
					_ = db.Get(&flip, `SELECT flip FROM maker.cat_file_flip WHERE header_id = $1`, headerID)
					return flip
				}).Should(Equal(test_data.FlipEthV100Address()))

				sw := watcher.NewStorageWatcher(db, -1, &statusWriter)
				sw.AddTransformers(storageInitializers)
				// This blocks right now, need to make test file to read from
				//err = w.Execute()
				//Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
