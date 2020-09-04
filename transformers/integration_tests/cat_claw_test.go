package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_claw"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat Claw transformer", func() {
	var logFetcher fetcher.ILogFetcher

	var catClawConfig = event.TransformerConfig{
		ContractAddresses: []string{test_data.Cat110Address()},
		ContractAbi:       constants.Cat110ABI(),
	}

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		logFetcher = fetcher.NewLogFetcher(blockChain)
	})


	It("persists a cat_claw event", func() {
		catClawBlockNumber := int64(10773034)
		header, err := persistHeader(db, catClawBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catClawConfig.TransformerName = constants.CatClawTable
		catClawConfig.Topic = constants.CatClawSignature()
		catClawConfig.StartingBlockNumber = catClawBlockNumber
		catClawConfig.EndingBlockNumber = catClawBlockNumber

		initializer := event.ConfiguredTransformer{
			Config:      catClawConfig,
			Transformer: cat_claw.Transformer{},
		}

		t := initializer.NewTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catClawConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catClawConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = t.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult catClawModel
		err = db.Get(&dbResult, `SELECT address_id, msg_sender, rad FROM maker.cat_claw`)
		Expect(err).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress(test_data.Cat110Address(), db)
		Expect(addressErr).NotTo(HaveOccurred())

		msgSender := "0xF32836B9E1f47a0515c6Ec431592D5EbC276407f"
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(msgSender, db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		Expect(dbResult.AddressID).To(Equal(addressID))
		Expect(dbResult.MsgSender).To(Equal(msgSenderID))
		Expect(dbResult.Rad).To(Equal("164878299999999999999707170900220955677367419868"))
	})
})

type catClawModel struct {
	MsgSender int64 `db:"msg_sender"`
	AddressID int64 `db:"address_id"`
	Rad string
}
