package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_sorted_offer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogSortedOffer Transformer", func() {
	logSortedOfferConfig := event.TransformerConfig{
		TransformerName:   constants.LogSortedOfferTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogSortedOfferSignature(),
	}

	It("fetches and transforms a LogSortedOffer event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(9440502)
		logSortedOfferConfig.StartingBlockNumber = blockNumber
		logSortedOfferConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logSortedOfferConfig,
			Transformer: log_sorted_offer.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisOneAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisOneAddress)},
			[]common.Hash{common.HexToHash(logSortedOfferConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []string
		err = db.Select(&dbResult, `SELECT offer_id FROM maker.log_sorted_offer`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult).To(ConsistOf("811647", "811648"))
	})

	It("fetches and transforms a LogSortedOffer event for OASIS_MATCHING_MARKET_TWO contract", func() {
		blockNumber := int64(9881525)
		logSortedOfferConfig.StartingBlockNumber = blockNumber
		logSortedOfferConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logSortedOfferConfig,
			Transformer: log_sorted_offer.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisTwoAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_TWO")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisTwoAddress)},
			[]common.Hash{common.HexToHash(logSortedOfferConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []string
		err = db.Select(&dbResult, `SELECT offer_id from maker.log_sorted_offer`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult).To(ConsistOf("253097"))
	})
})
