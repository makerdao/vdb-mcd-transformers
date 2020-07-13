package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_address"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VowFileAuctionAddress LogNoteTransformer", func() {
	var (
		initializer event.ConfiguredTransformer
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		config := event.TransformerConfig{
			TransformerName:   constants.VowFileAuctionAddressTable,
			ContractAddresses: []string{test_data.VowAddress()},
			ContractAbi:       constants.VowABI(),
			Topic:             constants.VowFileAuctionAddressSignature(),
		}

		addresses = event.HexStringsToAddresses(config.ContractAddresses)
		topics = []common.Hash{common.HexToHash(config.Topic)}

		initializer = event.ConfiguredTransformer{
			Config:      config,
			Transformer: auction_address.Transformer{},
		}
	})

	It("fetches and transforms a Vow.file auction address event", func() {
		blockNumber := int64(9017707)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := initializer.NewTransformer(db)
		executeErr := tr.Execute(eventLogs)
		Expect(executeErr).NotTo(HaveOccurred())

		var dbResult []vowFileAuctionAddressModel
		getVowFileErr := db.Select(&dbResult, `SELECT msg_sender, what, data from maker.vow_file_auction_address`)
		Expect(getVowFileErr).NotTo(HaveOccurred())

		var msgSenderId int64
		msgSenderAddress := common.HexToAddress("0xbe8e3e3618f7474f8cb1d074a26affef007e98fb").Hex()
		getMsgSenderIdErr := db.Get(&msgSenderId, `SELECT id FROM public.addresses WHERE address = $1`, msgSenderAddress)
		Expect(getMsgSenderIdErr).NotTo(HaveOccurred())

		var dataAddressId int64
		dataAddress := common.HexToAddress("0x4d95a049d5b0b7d32058cd3f2163015747522e99").Hex()
		getDataAddressIdErr := db.Get(&dataAddressId, `SELECT id FROM public.addresses WHERE address = $1`, dataAddress)
		Expect(getDataAddressIdErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].MsgSender).To(Equal(msgSenderId))
		Expect(dbResult[0].What).To(Equal("flopper"))
		Expect(dbResult[0].Data).To(Equal(dataAddressId))
	})
})

type vowFileAuctionAddressModel struct {
	MsgSender int64 `db:"msg_sender"`
	What      string
	Data      int64
}
