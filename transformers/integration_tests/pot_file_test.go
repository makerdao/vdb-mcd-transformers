package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/dsr"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotFile EventTransformers", func() {
	BeforeEach(func() {
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
			blockNumber = int64(8928300)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			potFileDSRConfig := event.TransformerConfig{
				TransformerName:     constants.PotFileDSRTable,
				ContractAddresses:   []string{test_data.PotAddress()},
				ContractAbi:         constants.PotABI(),
				Topic:               constants.PotFileDSRSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = event.HexStringsToAddresses(potFileDSRConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(potFileDSRConfig.Topic)}

			initializer := event.ConfiguredTransformer{
				Config:      potFileDSRConfig,
				Transformer: dsr.Transformer{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			tr := initializer.NewTransformer(db)
			executeErr := tr.Execute(eventLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Pot.file dsr event", func() {
			var dbResult potFileModel
			getFileErr := db.Get(&dbResult, `SELECT msg_sender, what, data FROM maker.pot_file_dsr`)
			Expect(getFileErr).NotTo(HaveOccurred())

			msgSender := shared.GetChecksumAddressString("0x000000000000000000000000be8e3e3618f7474f8cb1d074a26affef007e98fb")
			msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
			Expect(msgSenderErr).NotTo(HaveOccurred())

			Expect(dbResult.MsgSender).To(Equal(msgSenderID))
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
			blockNumber = int64(8928163)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			potFileVowConfig := event.TransformerConfig{
				TransformerName:     constants.PotFileVowTable,
				ContractAddresses:   []string{test_data.PotAddress()},
				ContractAbi:         constants.PotABI(),
				Topic:               constants.PotFileVowSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = event.HexStringsToAddresses(potFileVowConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(potFileVowConfig.Topic)}

			initializer := event.ConfiguredTransformer{
				Config:      potFileVowConfig,
				Transformer: vow.Transformer{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			tr := initializer.NewTransformer(db)
			executeErr := tr.Execute(eventLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Pot.file vow event", func() {
			var dbResult potFileModel
			getFileErr := db.Get(&dbResult, `SELECT msg_sender, what, data FROM maker.pot_file_vow`)
			Expect(getFileErr).NotTo(HaveOccurred())

			msgSender := shared.GetChecksumAddressString("0x000000000000000000000000baa65281c2fa2baacb2cb550ba051525a480d3f4")
			msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
			Expect(msgSenderErr).NotTo(HaveOccurred())

			Expect(dbResult.MsgSender).To(Equal(msgSenderID))
			Expect(dbResult.What).To(Equal("vow"))
			Expect(dbResult.Data).To(Equal(test_data.VowAddress()))
		})
	})
})

type potFileModel struct {
	MsgSender int64 `db:"msg_sender"`
	What      string
	Data      string
}
