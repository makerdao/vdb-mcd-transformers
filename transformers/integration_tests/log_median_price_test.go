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
	Context("LogMedianPrice event on Median BAT", func() {
		val := "179042302500000000"
		age := "1588036910"
		logMedianPriceIntegrationTest(9957995, test_data.MedianBatAddress(), val, age)
	})

	Context("LogMedianPrice event on Median COMP", func() {
		val := "136240785441850000000"
		age := "1601305976"
		logMedianPriceIntegrationTest(10951676, test_data.MedianCompAddress(), val, age)
	})

	Context("LogMedianPrice event on Median ETH", func() {
		val := "192578360000000000000"
		age := "1588003362"
		logMedianPriceIntegrationTest(9955467, test_data.MedianEthAddress(), val, age)
	})

	Context("LogMedianPrice event on Median KNC", func() {
		val := "1692228609000000000"
		age := "1594112655"
		logMedianPriceIntegrationTest(10411362, test_data.MedianKncAddress(), val, age)
	})

	Context("LogMedianPrice event on Median LINK", func() {
		val := "10658915841150000000"
		age := "1601302284"
		logMedianPriceIntegrationTest(10951425, test_data.MedianLinkAddress(), val, age)
	})

	Context("LogMedianPrice event on Median MANA", func() {
		val := "42028498300000000"
		age := "1595917269"
		logMedianPriceIntegrationTest(10546364, test_data.MedianManaAddress(), val, age)
	})

	Context("LogMedianPrice event on Median USDT-A", func() {
		val := "1000732490700000000"
		age := "1599686517"
		logMedianPriceIntegrationTest(10829955, test_data.MedianUsdtAddress(), val, age)
	})

	Context("LogMedianPrice event on Median WBTC", func() {
		val := "8830300000000000000000"
		age := "1588280393"
		logMedianPriceIntegrationTest(9976164, test_data.MedianWbtcAddress(), val, age)
	})

	Context("LogMedianPrice event on Median ZRX", func() {
		val := "324434000000000000"
		age := "1593308874"
		logMedianPriceIntegrationTest(10351386, test_data.MedianZrxAddress(), val, age)
	})
})

func logMedianPriceIntegrationTest(blockNumber int64, contractAddressHex, val, age string) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		config := event.TransformerConfig{
			ContractAbi:         constants.MedianABI(),
			ContractAddresses:   []string{contractAddressHex},
			EndingBlockNumber:   blockNumber,
			StartingBlockNumber: blockNumber,
			TransformerName:     constants.LogMedianPriceTable,
		}
		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_median_price.Transformer{},
		}

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, fetchErr := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(contractAddressHex)},
			[]common.Hash{common.HexToHash(constants.LogMedianPriceSignature())},
			header)
		Expect(fetchErr).NotTo(HaveOccurred())
		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		executeErr := transformer.Execute(eventLogs)
		Expect(executeErr).NotTo(HaveOccurred())

		var dbResults []logMedianPriceModel
		queryErr := db.Select(&dbResults, `SELECT address_id, val, age from maker.log_median_price`)
		Expect(queryErr).NotTo(HaveOccurred())
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(contractAddressHex, db)
		Expect(contractAddressErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		Expect(dbResults[0].AddressID).To(Equal(contractAddressID))
		Expect(dbResults[0].Val).To(Equal(val))
		Expect(dbResults[0].Age).To(Equal(age))
	})
}

type logMedianPriceModel struct {
	Val       string
	Age       string
	AddressID int64 `db:"address_id"`
}
