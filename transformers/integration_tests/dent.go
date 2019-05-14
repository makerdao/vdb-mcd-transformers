package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("Dent transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		logFetcher  fetcher.ILogFetcher
		tr          transformer.EventTransformer
		dentConfig  transformer.EventTransformerConfig
		addresses   []common.Address
		topics      []common.Hash
		initializer shared.LogNoteTransformer
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		dentConfig = transformer.EventTransformerConfig{
			TransformerName:   mcdConstants.DentLabel,
			ContractAddresses: []string{mcdConstants.FlipperContractAddress(), mcdConstants.FlopperContractAddress()},
			ContractAbi:       mcdConstants.FlipperABI(),
			Topic:             mcdConstants.DentSignature(),
		}

		addresses = transformer.HexStringsToAddresses(dentConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dentConfig.Topic)}
		logFetcher = fetcher.NewLogFetcher(blockChain)

		initializer = shared.LogNoteTransformer{
			Config:     dentConfig,
			Converter:  &dent.DentConverter{},
			Repository: &dent.DentRepository{},
		}
	})

	It("persists a flop dent log event", func() {
		blockNumber := int64(8955613)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		tr = initializer.NewLogNoteTransformer(db)
		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dent.DentModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("10000000000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("2"))
		Expect(dbResult[0].Guy).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538637780 + mcdConstants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})

	It("rechecks header for flop dent log event", func() {
		blockNumber := int64(8955613)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		tr = initializer.NewLogNoteTransformer(db)
		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var dentChecked []int
		err = db.Select(&dentChecked, `SELECT dent_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(dentChecked[0]).To(Equal(2))
	})

	It("persists a flip dent log event", func() {
		//TODO: There are currently no Flip.dent events on Kovan
	})
})
