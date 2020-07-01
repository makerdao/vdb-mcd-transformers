package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_median_price"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogMedianPrice Transformer", func() {
	config := event.TransformerConfig{
		TransformerName:   constants.LogMedianPriceTable,
		ContractAddresses: test_data.MedianAddresses(),
		ContractAbi:       constants.MedianABI(),
		Topic:             constants.LogMedianPriceSignature(),
	}

	It("fetches and transforms a LogMedianPrice event on Median ETH", func() {
		blockNumber := int64(9955467)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_median_price.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		medianEthAddress := constants.GetContractAddress("MEDIAN_ETH")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(medianEthAddress)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logMedianPriceModel
		err = db.Select(&dbResults, `SELECT address_id, val, age from maker.log_median_price`)
		Expect(err).NotTo(HaveOccurred())
		expectedAddressID, addressErr := shared.GetOrCreateAddress(medianEthAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		Expect(dbResults[0].AddressID).To(Equal(expectedAddressID))
		Expect(dbResults[0].Val).To(Equal("192578360000000000000"))
		Expect(dbResults[0].Age).To(Equal("1588003362"))
	})

	It("fetches and transforms a LogMedianPrice event on Median BAT", func() {
		blockNumber := int64(9957995)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_median_price.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		medianBatAddress := constants.GetContractAddress("MEDIAN_BAT")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(medianBatAddress)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logMedianPriceModel
		err = db.Select(&dbResults, `SELECT address_id, val, age from maker.log_median_price`)
		Expect(err).NotTo(HaveOccurred())
		expectedAddressID, addressErr := shared.GetOrCreateAddress(medianBatAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		Expect(dbResults[0].AddressID).To(Equal(expectedAddressID))
		Expect(dbResults[0].Val).To(Equal("179042302500000000"))
		Expect(dbResults[0].Age).To(Equal("1588036910"))
	})

	It("fetches and transforms a LogMedianPrice event on Median WBTC", func() {
		blockNumber := int64(9976164)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_median_price.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		medianWbtcAddress := constants.GetContractAddress("MEDIAN_WBTC")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(medianWbtcAddress)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logMedianPriceModel
		err = db.Select(&dbResults, `SELECT address_id, val, age from maker.log_median_price`)
		Expect(err).NotTo(HaveOccurred())
		expectedAddressID, addressErr := shared.GetOrCreateAddress(medianWbtcAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		Expect(dbResults[0].AddressID).To(Equal(expectedAddressID))
		Expect(dbResults[0].Val).To(Equal("8830300000000000000000"))
		Expect(dbResults[0].Age).To(Equal("1588280393"))
	})

	It("fetches and transforms a LogMedianPrice event on Median ZRX", func() {
		blockNumber := int64(10351386)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_median_price.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		medianZRXAddress := constants.GetContractAddress("MEDIAN_ZRX")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(medianZRXAddress)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logMedianPriceModel
		err = db.Select(&dbResults, `SELECT address_id, val, age from maker.log_median_price`)
		Expect(err).NotTo(HaveOccurred())
		expectedAddressID, addressErr := shared.GetOrCreateAddress(medianZRXAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		Expect(dbResults[0].AddressID).To(Equal(expectedAddressID))
		Expect(dbResults[0].Val).To(Equal("324434000000000000"))
		Expect(dbResults[0].Age).To(Equal("1593308874"))
	})
})

type logMedianPriceModel struct {
	Val       string
	Age       string
	AddressID int64 `db:"address_id"`
}
