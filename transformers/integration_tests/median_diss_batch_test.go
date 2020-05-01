package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/batch"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MedianDissBatch EventTransformer", func() {
	medianDissConfig := event.TransformerConfig{
		TransformerName:   constants.MedianDissBatchTable,
		ContractAddresses: test_data.MedianAddresses(),
		ContractAbi:       constants.MedianABI(),
		Topic:             constants.MedianDissBatchSignature(),
	}

	It("transforms Median diss batch log events", func() {
		blockNumber := int64(8936530)
		medianDissConfig.StartingBlockNumber = blockNumber
		medianDissConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		// TODO: fetch event from blockchain once one exists
		// logFetcher := fetcher.NewLogFetcher(blockChain)
		// logs, logErr := logFetcher.FetchLogs(
		// 	event.HexStringsToAddresses(medianDissConfig.ContractAddresses),
		// 	[]common.Hash{common.HexToHash(medianDissConfig.Topic)},
		// 	header)
		// Expect(logErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, []types.Log{test_data.RawMedianDissBatchLogFiveAddresses}, db)

		transformer := event.ConfiguredTransformer{
			Config:      medianDissConfig,
			Transformer: batch.Transformer{},
		}.NewTransformer(db)

		transformErr := transformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResults []MedianDissBatchModel
		err := db.Select(&dbResults, `SELECT address_id, msg_sender, a_length FROM maker.median_diss_batch ORDER BY id`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(test_data.MedianEthAddress(), db)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))

		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress("0xe87F55Af91068a1DA44095138F3d37C45894Eb21", db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(msgSenderAddressID, 10)))

		Expect(dbResult.ALength).To(Equal("5"))

		var addresses []string
		addressesError := db.Get(pq.Array(&addresses), `SELECT a FROM maker.median_diss_batch ORDER BY id`)
		Expect(addressesError).NotTo(HaveOccurred())
		Expect(addresses).To(ConsistOf(
			"0xA52F23A651d1FA7c2610753C768103ee8C498f22",
			"0xce91dB32ad1C91278A56CBb2d8f24f9315043DE9",
			"0x3482f7a06Db71F8EcAc04F882546a66081311667",
			"0x702F365E1E559D9dC7b1af6aB9be64feb9c4D822",
		))
	})
})

type MedianDissBatchModel struct {
	AddressID string `db:"address_id"`
	MsgSender string `db:"msg_sender"`
	ALength   string `db:"a_length"`
}
