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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/vow"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	"strconv"
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
		ContractAddresses: []string{test_data.CatAddress()},
		ContractAbi:       constants.CatABI(),
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
		chopLumpBlockNumber := int64(13773400)
		header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = constants.CatFileChopLumpLabel
		catFileConfig.Topic = constants.CatFileChopLumpSignature()
		catFileConfig.StartingBlockNumber = chopLumpBlockNumber
		catFileConfig.EndingBlockNumber = chopLumpBlockNumber

		initializer := shared.EventTransformer{
			Config:     catFileConfig,
			Converter:  &chop_lump.CatFileChopLumpConverter{},
			Repository: &chop_lump.CatFileChopLumpRepository{},
		}
		transformer := initializer.NewEventTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []catFileChopLumpModel
		err = db.Select(&dbResult, `SELECT what, ilk_id, data FROM maker.cat_file_chop_lump`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))

		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult[0].What).To(Equal("lump"))
		Expect(dbResult[0].Data).To(Equal("1500000000000000000"))
	})

	It("persists a chop lump event (chop)", func() {
		chopLumpBlockNumber := int64(13773376)
		header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = constants.CatFileChopLumpLabel
		catFileConfig.Topic = constants.CatFileChopLumpSignature()
		catFileConfig.StartingBlockNumber = chopLumpBlockNumber
		catFileConfig.EndingBlockNumber = chopLumpBlockNumber

		initializer := shared.EventTransformer{
			Config:     catFileConfig,
			Converter:  &chop_lump.CatFileChopLumpConverter{},
			Repository: &chop_lump.CatFileChopLumpRepository{},
		}
		transformer := initializer.NewEventTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []catFileChopLumpModel
		err = db.Select(&dbResult, `SELECT what, ilk_id, data FROM maker.cat_file_chop_lump`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))

		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult[0].What).To(Equal("chop"))
		Expect(dbResult[0].Data).To(Equal("1050000000000000000000000000"))
	})

	It("persists a flip event", func() {
		flipBlockNumber := int64(13772999)
		header, err := persistHeader(db, flipBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = constants.CatFileFlipLabel
		catFileConfig.Topic = constants.CatFileFlipSignature()
		catFileConfig.StartingBlockNumber = flipBlockNumber
		catFileConfig.EndingBlockNumber = flipBlockNumber

		initializer := shared.EventTransformer{
			Config:     catFileConfig,
			Converter:  &flip.CatFileFlipConverter{},
			Repository: &flip.CatFileFlipRepository{},
		}

		t := initializer.NewEventTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = t.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []catFileFlipModel
		err = db.Select(&dbResult, `SELECT ilk_id, what, flip FROM maker.cat_file_flip`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult[0].What).To(Equal("flip"))
		Expect(dbResult[0].Flip).To(Equal("0x494d6664A6b305F1f6dbdED879f01E5DC1EA8B55"))
	})

	It("persists a vow event", func() {
		vowBlockNumber := int64(13772981)
		header, err := persistHeader(db, vowBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		catFileConfig.TransformerName = constants.CatFileVowLabel
		catFileConfig.Topic = constants.CatFileVowSignature()
		catFileConfig.StartingBlockNumber = vowBlockNumber
		catFileConfig.EndingBlockNumber = vowBlockNumber

		initializer := shared.EventTransformer{
			Config:     catFileConfig,
			Converter:  &vow.CatFileVowConverter{},
			Repository: &vow.CatFileVowRepository{},
		}
		t := initializer.NewEventTransformer(db)

		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(catFileConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = t.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, vowBlockNumber)
		Expect(err).NotTo(HaveOccurred())

		var dbResult catFileVowModel
		err = db.Get(&dbResult, `SELECT what, data FROM maker.cat_file_vow`)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.What).To(Equal("vow"))
		Expect(dbResult.Data).To(Equal("0xdC02a6b2eCd2ed6C54a5EC1F1585FE82137D31dD"))
	})
})

type catFileChopLumpModel struct {
	Ilk  string `db:"ilk_id"`
	What string
	Data string
}

type catFileFlipModel struct {
	Ilk  string `db:"ilk_id"`
	What string
	Flip string
}

type catFileVowModel struct {
	What string
	Data string
}
