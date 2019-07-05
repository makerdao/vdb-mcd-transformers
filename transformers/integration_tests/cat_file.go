// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package integration_tests

import (
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/vow"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("Cat File transformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
		rpcClient  client.RpcClient
		err        error
		ethClient  *ethclient.Client
		logFetcher fetcher.ILogFetcher
	)

	var catFileConfig = transformer.EventTransformerConfig{
		ContractAddresses: []string{mcdConstants.CatContractAddress()},
		ContractAbi:       mcdConstants.CatABI(),
	}

	BeforeEach(func() {
		rpcClient, ethClient, err = getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		logFetcher = fetcher.NewLogFetcher(blockChain)
	})

	It("persists a chop lump event (lump)", func() {
		chopLumpBlockNumber := int64(11257529)
		header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = mcdConstants.CatFileChopLumpLabel
		catFileConfig.Topic = mcdConstants.CatFileChopLumpSignature()
		catFileConfig.StartingBlockNumber = chopLumpBlockNumber
		catFileConfig.EndingBlockNumber = chopLumpBlockNumber

		initializer := shared.LogNoteTransformer{
			Config:     catFileConfig,
			Converter:  &chop_lump.CatFileChopLumpConverter{},
			Repository: &chop_lump.CatFileChopLumpRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []catFileChopLumpModel
		err = db.Select(&dbResult, `SELECT what, ilk_id, data, log_idx FROM maker.cat_file_chop_lump`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		sort.Sort(byLogIndexChopLump(dbResult))

		ilkID, err := shared.GetOrCreateIlk("0x434f4c352d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].What).To(Equal("lump"))
		Expect(dbResult[0].Data).To(Equal("10000000000000000000"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(3)))
	})

	It("persists a chop lump event (chop)", func() {
		chopLumpBlockNumber := int64(11257512)
		header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = mcdConstants.CatFileChopLumpLabel
		catFileConfig.Topic = mcdConstants.CatFileChopLumpSignature()
		catFileConfig.StartingBlockNumber = chopLumpBlockNumber
		catFileConfig.EndingBlockNumber = chopLumpBlockNumber

		initializer := shared.LogNoteTransformer{
			Config:     catFileConfig,
			Converter:  &chop_lump.CatFileChopLumpConverter{},
			Repository: &chop_lump.CatFileChopLumpRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []catFileChopLumpModel
		err = db.Select(&dbResult, `SELECT what, ilk_id, data, log_idx FROM maker.cat_file_chop_lump`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		sort.Sort(byLogIndexChopLump(dbResult))

		ilkID, err := shared.GetOrCreateIlk("0x434f4c352d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].What).To(Equal("chop"))
		Expect(dbResult[0].Data).To(Equal("1050000000000000000000000000"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(3)))
	})

	It("persists a flip event", func() {
		flipBlockNumber := int64(11257255)
		header, err := persistHeader(db, flipBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = mcdConstants.CatFileFlipLabel
		catFileConfig.Topic = mcdConstants.CatFileFlipSignature()
		catFileConfig.StartingBlockNumber = flipBlockNumber
		catFileConfig.EndingBlockNumber = flipBlockNumber

		initializer := shared.LogNoteTransformer{
			Config:     catFileConfig,
			Converter:  &flip.CatFileFlipConverter{},
			Repository: &flip.CatFileFlipRepository{},
		}

		t := initializer.NewLogNoteTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = t.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []catFileFlipModel
		err = db.Select(&dbResult, `SELECT ilk_id, what, flip FROM maker.cat_file_flip`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x434f4c352d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].What).To(Equal("flip"))
		Expect(dbResult[0].Flip).To(Equal("0x3fd71145A9597335DeECB3823286d8d5f0d29B82"))
	})

	It("rechecks a flip event", func() {
		flipBlockNumber := int64(11257255)
		header, err := persistHeader(db, flipBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = mcdConstants.CatFileFlipLabel
		catFileConfig.Topic = mcdConstants.CatFileFlipSignature()
		catFileConfig.StartingBlockNumber = flipBlockNumber
		catFileConfig.EndingBlockNumber = flipBlockNumber

		initializer := shared.LogNoteTransformer{
			Config:     catFileConfig,
			Converter:  &flip.CatFileFlipConverter{},
			Repository: &flip.CatFileFlipRepository{},
		}

		t := initializer.NewLogNoteTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = t.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		err = t.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, flipBlockNumber)
		Expect(err).NotTo(HaveOccurred())

		var catFlipChecked []int
		err = db.Select(&catFlipChecked, `SELECT cat_file_flip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(catFlipChecked[0]).To(Equal(2))
	})

	It("persists a vow event", func() {
		vowBlockNumber := int64(11257175)
		header, err := persistHeader(db, vowBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = mcdConstants.CatFileVowLabel
		catFileConfig.Topic = mcdConstants.CatFileVowSignature()
		catFileConfig.StartingBlockNumber = vowBlockNumber
		catFileConfig.EndingBlockNumber = vowBlockNumber

		initializer := shared.LogNoteTransformer{
			Config:     catFileConfig,
			Converter:  &vow.CatFileVowConverter{},
			Repository: &vow.CatFileVowRepository{},
		}
		t := initializer.NewLogNoteTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = t.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		err = t.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, vowBlockNumber)
		Expect(err).NotTo(HaveOccurred())

		var catVowChecked []int
		err = db.Select(&catVowChecked, `SELECT cat_file_vow_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(catVowChecked[0]).To(Equal(2))
	})

	It("rechecks a vow event", func() {
		vowBlockNumber := int64(11257175)
		header, err := persistHeader(db, vowBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = mcdConstants.CatFileVowLabel
		catFileConfig.Topic = mcdConstants.CatFileVowSignature()
		catFileConfig.StartingBlockNumber = vowBlockNumber
		catFileConfig.EndingBlockNumber = vowBlockNumber

		initializer := shared.LogNoteTransformer{
			Config:     catFileConfig,
			Converter:  &vow.CatFileVowConverter{},
			Repository: &vow.CatFileVowRepository{},
		}
		t := initializer.NewLogNoteTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = t.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []catFileVowModel
		err = db.Select(&dbResult, `SELECT what, data, log_idx FROM maker.cat_file_vow`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("vow"))
		Expect(dbResult[0].Data).To(Equal("0x528C10D35D9612d1667D2D9b8915f80A573a800B"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(2)))
	})
})

type catFileChopLumpModel struct {
	Ilk              string `db:"ilk_id"`
	What             string
	Data             string
	TransactionIndex uint   `db:"tx_idx"`
	LogIndex         uint   `db:"log_idx"`
	Raw              []byte `db:"raw_log"`
}

type catFileFlipModel struct {
	Ilk              string `db:"ilk_id"`
	What             string
	Flip             string
	TransactionIndex uint   `db:"tx_idx"`
	LogIndex         uint   `db:"log_idx"`
	Raw              []byte `db:"raw_log"`
}

type catFileVowModel struct {
	What             string
	Data             string
	TransactionIndex uint   `db:"tx_idx"`
	LogIndex         uint   `db:"log_idx"`
	Raw              []byte `db:"raw_log"`
}

type byLogIndexChopLump []catFileChopLumpModel

func (c byLogIndexChopLump) Len() int           { return len(c) }
func (c byLogIndexChopLump) Less(i, j int) bool { return c[i].LogIndex < c[j].LogIndex }
func (c byLogIndexChopLump) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
