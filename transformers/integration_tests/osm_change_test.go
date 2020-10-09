package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/osm_change"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO update when real event log exists
var _ = XDescribe("OsmChange EventTransformer", func() {
	osmChangeConfig := event.TransformerConfig{
		TransformerName:   constants.OsmChangeTable,
		ContractAddresses: []string{test_data.OsmEthAddress()},
		ContractAbi:       constants.OsmABI(),
		Topic:             constants.OsmChangeSignature(),
	}

	It("transforms OSM change log events", func() {
		blockNumber := int64(8928180)
		osmChangeConfig.StartingBlockNumber = blockNumber
		osmChangeConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, logErr := logFetcher.FetchLogs(
			event.HexStringsToAddresses(osmChangeConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(osmChangeConfig.Topic)},
			header)
		Expect(logErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := event.ConfiguredTransformer{
			Config:      osmChangeConfig,
			Transformer: osm_change.Transformer{},
		}.NewTransformer(db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []OsmChangeModel
		err = db.Select(&dbResults, `SELECT address_id from maker.osm_change`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, "0x000000000000000000000000a950524441892a31ebddf91d3ceefa04bf454466")
		Expect(contractAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))

		msgSenderAddressID, msgSenderAddressErr := repository.GetOrCreateAddress(db, "0x000000000000000000000000a950524441892a31ebddf91d3ceefa04bf454466")
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(msgSenderAddressID, 10)))

		srcAddressID, srcAddressErr := repository.GetOrCreateAddress(db, "0x000000000000000000000000a950524441892a31ebddf91d3ceefa04bf454466")
		Expect(srcAddressErr).NotTo(HaveOccurred())
		Expect(dbResult.Src).To(Equal(strconv.FormatInt(srcAddressID, 10)))
	})
})

type OsmChangeModel struct {
	AddressID string `db:"address_id"`
	MsgSender string `db:"msg_sender"`
	Src       string
}
