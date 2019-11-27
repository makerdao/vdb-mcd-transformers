package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/dsr"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotFile EventTransformers", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
	})

	Describe("Pot file dsr", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			logs        []types.Log
			topics      []common.Hash
		)

		BeforeEach(func() {
			blockNumber = int64(14764699)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			potFileDSRConfig := transformer.EventTransformerConfig{
				TransformerName:     constants.PotFileDSRLabel,
				ContractAddresses:   []string{test_data.PotAddress()},
				ContractAbi:         constants.PotABI(),
				Topic:               constants.PotFileDSRSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(potFileDSRConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(potFileDSRConfig.Topic)}

			initializer := event.Transformer{
				Config:     potFileDSRConfig,
				Converter:  &dsr.Converter{},
				Repository: &dsr.Repository{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

			tr := initializer.NewTransformer(db)
			executeErr := tr.Execute(headerSyncLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Pot.file dsr event from Kovan", func() {
			var dbResult potFileDSRModel
			getFileErr := db.Get(&dbResult, `SELECT what, data FROM maker.pot_file_dsr`)
			Expect(getFileErr).NotTo(HaveOccurred())

			Expect(dbResult.What).To(Equal("dsr"))
			Expect(dbResult.Data).To(Equal("1000000000627937192491029810"))
		})
	})

	Describe("Pot file vow", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			logs        []types.Log
			topics      []common.Hash
		)

		BeforeEach(func() {
			blockNumber = int64(14764543)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			potFileVowConfig := transformer.EventTransformerConfig{
				TransformerName:     constants.PotFileVowLabel,
				ContractAddresses:   []string{test_data.PotAddress()},
				ContractAbi:         constants.PotABI(),
				Topic:               constants.PotFileVowSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(potFileVowConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(potFileVowConfig.Topic)}

			initializer := event.Transformer{
				Config:     potFileVowConfig,
				Converter:  &vow.Converter{},
				Repository: &vow.Repository{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

			tr := initializer.NewTransformer(db)
			executeErr := tr.Execute(headerSyncLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Pot.file vow event from Kovan", func() {
			var dbResult potFileDSRModel
			getFileErr := db.Get(&dbResult, `SELECT what, data FROM maker.pot_file_vow`)
			Expect(getFileErr).NotTo(HaveOccurred())

			Expect(dbResult.What).To(Equal("vow"))
			Expect(dbResult.Data).To(Equal(test_data.VowAddress()))
		})
	})
})

type potFileDSRModel struct {
	What string
	Data string
}
