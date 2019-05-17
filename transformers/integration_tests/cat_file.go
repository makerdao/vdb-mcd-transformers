// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
<<<<<<< HEAD
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
=======
	"sort"
	"strconv"

	c2 "github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	fetch "github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
>>>>>>> Fixes integration tests
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

	It("persists a chop lump event", func() {
		chopLumpBlockNumber := int64(10771104)
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

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []chop_lump.CatFileChopLumpModel
		err = db.Select(&dbResult, `SELECT what, ilk_id, data, log_idx FROM maker.cat_file_chop_lump`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		sort.Sort(byLogIndexChopLump(dbResult))

		ilkID, err := shared.GetOrCreateIlk("4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].What).To(Equal("lump"))
		Expect(dbResult[0].Data).To(Equal("10000000000000000000000000000000000000000000000000"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(1)))

		Expect(dbResult[1].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[1].What).To(Equal("chop"))
		Expect(dbResult[1].Data).To(Equal("1000000000000000000000000000"))
		Expect(dbResult[1].LogIndex).To(Equal(uint(2)))
	})

	It("rechecks header for chop lump event", func() {
		chopLumpBlockNumber := int64(10691268)
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

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, chopLumpBlockNumber)
		Expect(err).NotTo(HaveOccurred())

		var catChopLumpChecked []int
		err = db.Select(&catChopLumpChecked, `SELECT cat_file_chop_lump_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(catChopLumpChecked[0]).To(Equal(2))
	})

	It("persists a flip event", func() {
		flipBlockNumber := int64(10771104)
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

		err = t.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flip.CatFileFlipModel
		err = db.Select(&dbResult, `SELECT ilk_id, what, flip FROM maker.cat_file_flip`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].What).To(Equal("flip"))
		Expect(dbResult[0].Flip).To(Equal("0x259D562D7d14c11efCF4fc1678F29D3f618B68aD"))
	})

	It("rechecks a flip event", func() {
		flipBlockNumber := int64(10771104)
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

		err = t.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = t.Execute(logs, header, constants.HeaderRecheck)
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
		vowBlockNumber := int64(10771088)
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

		err = t.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vow.CatFileVowModel
		err = db.Select(&dbResult, `SELECT what, data, log_idx FROM maker.cat_file_vow`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("vow"))
		Expect(dbResult[0].Data).To(Equal("0xa2c0D575CB4e1F145830326420e0CcFab8BeBc1d"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(2)))
	})
})

type byLogIndexChopLump []chop_lump.CatFileChopLumpModel

func (c byLogIndexChopLump) Len() int           { return len(c) }
func (c byLogIndexChopLump) Less(i, j int) bool { return c[i].LogIndex < c[j].LogIndex }
func (c byLogIndexChopLump) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
