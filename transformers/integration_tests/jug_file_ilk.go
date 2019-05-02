package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"strconv"

	c2 "github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	fetch "github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Jug File Ilk LogNoteTransformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
	})

	It("transforms jug file ilk log events", func() {
		blockNumber := int64(10691361)
		config := transformer.EventTransformerConfig{
			TransformerName:     constants.JugFileIlkLabel,
			ContractAddresses:   []string{test_data.KovanJugContractAddress},
			ContractAbi:         test_data.KovanJugABI,
			Topic:               test_data.KovanJugFileIlkSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     config,
			Converter:  &ilk.JugFileIlkConverter{},
			Repository: &ilk.JugFileIlkRepository{},
		}
		tr := initializer.NewLogNoteTransformer(db)

		f := fetch.NewFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, c2.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []ilk.JugFileIlkModel
		err = db.Select(&dbResult, `SELECT ilk_id, what, data FROM maker.jug_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("4554482d43000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].What).To(Equal("duty"))
		Expect(dbResult[0].Data).To(Equal("1000000000158153903837946257"))
	})

	It("rechecks jug file ilk event", func() {
		blockNumber := int64(10691361)
		config := transformer.EventTransformerConfig{
			TransformerName:     constants.JugFileIlkLabel,
			ContractAddresses:   []string{test_data.KovanJugContractAddress},
			ContractAbi:         test_data.KovanJugABI,
			Topic:               test_data.KovanJugFileIlkSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     config,
			Converter:  &ilk.JugFileIlkConverter{},
			Repository: &ilk.JugFileIlkRepository{},
		}
		tr := initializer.NewLogNoteTransformer(db)

		f := fetch.NewFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, c2.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, c2.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var jugFileIlkChecked []int
		err = db.Select(&jugFileIlkChecked, `SELECT jug_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(jugFileIlkChecked[0]).To(Equal(2))
	})
})
