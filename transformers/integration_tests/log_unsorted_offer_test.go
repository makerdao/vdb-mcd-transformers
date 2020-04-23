package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_unsorted_offer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogUnsortedOffer Transformer", func() {
	logUnsortedOfferConfig := event.TransformerConfig{
		TransformerName:   constants.LogUnsortedOfferTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogUnsortedOfferSignature(),
	}

	It("fetches and transforms a LogUnsortedOffer event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9243052)
		logUnsortedOfferConfig.StartingBlockNumber = blockNumber
		logUnsortedOfferConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logUnsortedOfferConfig,
			Transformer: log_unsorted_offer.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisOneAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisOneAddress)},
			[]common.Hash{common.HexToHash(logUnsortedOfferConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []string
		err = db.Select(&dbResult, `SELECT offer_id FROM maker.log_unsorted_offer`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult).To(ConsistOf("717050"))
	})
})
