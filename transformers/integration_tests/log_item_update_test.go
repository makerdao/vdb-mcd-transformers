package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_item_update"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogItemUpdate Transformer", func() {
	logItemUpdateConfig := event.TransformerConfig{
		TransformerName:   constants.LogItemUpdateTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogItemUpdateSignature(),
	}

	It("fetches and transforms a LogItemUpdate event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9613377)
		logItemUpdateConfig.StartingBlockNumber = blockNumber
		logItemUpdateConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logItemUpdateConfig,
			Transformer: log_item_update.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisOneAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisOneAddress)},
			[]common.Hash{common.HexToHash(logItemUpdateConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []logItemUpdateModel
		err = db.Select(&dbResult, `SELECT offer_id from maker.log_item_update`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(15))
		var offerIds []string
		for _, d := range dbResult {
			offerIds = append(offerIds, d.OfferID)
		}
		Expect(offerIds).To(ConsistOf([]string{
			"604350", "662672", "581874", "680605", "666869",
			"581386", "581829", "662659", "606999", "807535",
			"809995", "637558", "559973", "548250", "772805",
		}))
	})
	It("fetches and transforms a LogItemUpdate event for OASIS_MATCHING_MARKET_TWO contract", func() {
		blockNumber := int64(9827320)
		logItemUpdateConfig.StartingBlockNumber = blockNumber
		logItemUpdateConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logItemUpdateConfig,
			Transformer: log_item_update.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisTwoAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_TWO")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisTwoAddress)},
			[]common.Hash{common.HexToHash(logItemUpdateConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []logItemUpdateModel
		err = db.Select(&dbResult, `SELECT offer_id from maker.log_item_update`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(3))
		var offerIds []string
		for _, d := range dbResult {
			offerIds = append(offerIds, d.OfferID)
		}
		Expect(offerIds).To(ConsistOf([]string{"228713", "228696", "228695"}))
	})
})

type logItemUpdateModel struct {
	OfferID  string `db:"offer_id"`
	HeaderID int64
	LogID    int64 `db:"log_id"`
}
