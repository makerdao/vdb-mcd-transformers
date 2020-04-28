package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_insert"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogInsert Transformer", func() {
	config := event.TransformerConfig{
		TransformerName:   constants.LogInsertTable,
		ContractAddresses: test_data.OasisAddresses(),
		ContractAbi:       constants.OasisABI(),
		Topic:             constants.LogInsertSignature(),
	}

	It("fetches and transforms a LogInsert event for OASIS_MATCHING_MARKET_ONE contract", func() {
		blockNumber := int64(7346754)
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      config,
			Transformer: log_insert.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		oasisOneAddress := constants.GetContractAddress("OASIS_MATCHING_MARKET_ONE")
		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(oasisOneAddress)},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []logInsertModel
		err = db.Select(&dbResults, `SELECT address_id, keeper, offer_id from maker.log_insert`)
		Expect(err).NotTo(HaveOccurred())

		expectedAddressID, addressErr := shared.GetOrCreateAddress(oasisOneAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedKeeperID, keeperErr := shared.GetOrCreateAddress("0x3a32292c53bf42b6317334392bf0272da2983252", db)
		Expect(keeperErr).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.AddressID).To(Equal(expectedAddressID))
		Expect(dbResult.Keeper).To(Equal(expectedKeeperID))
		Expect(dbResult.OfferID).To(Equal("35394"))
	})

	// TODO: add test for Oasis Matching Market Two if detected
})

type logInsertModel struct {
	AddressID int64 `db:"address_id"`
	Keeper    int64
	OfferID   string `db:"offer_id"`
}
