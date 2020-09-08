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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/box"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/chop_lump_dunk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat File transformer", func() {
	var logFetcher fetcher.ILogFetcher

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		logFetcher = fetcher.NewLogFetcher(blockChain)
	})

	Describe("1.0.0 Cat Contract", func() {
		var catFileConfig = event.TransformerConfig{
			ContractAddresses: []string{test_data.Cat100Address()},
			ContractAbi:       constants.Cat100ABI(),
		}

		It("persists a chop lump event (lump)", func() {
			chopLumpBlockNumber := int64(8928392)
			header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
			Expect(err).NotTo(HaveOccurred())
			catFileConfig.TransformerName = constants.CatFileChopLumpDunkTable
			catFileConfig.Topic = constants.CatFileChopLumpSignature()
			catFileConfig.StartingBlockNumber = chopLumpBlockNumber
			catFileConfig.EndingBlockNumber = chopLumpBlockNumber

			initializer := event.ConfiguredTransformer{
				Config:      catFileConfig,
				Transformer: chop_lump_dunk.Transformer{},
			}
			transformer := initializer.NewTransformer(db)

			logs, err := logFetcher.FetchLogs(
				[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
				[]common.Hash{common.HexToHash(catFileConfig.Topic)},
				header)
			Expect(err).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			err = transformer.Execute(eventLogs)
			Expect(err).NotTo(HaveOccurred())

			var dbResult catFileChopLumpModel
			err = db.Get(&dbResult, `SELECT what, ilk_id, data FROM maker.cat_file_chop_lump_dunk`)
			Expect(err).NotTo(HaveOccurred())

			ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
			Expect(dbResult.What).To(Equal("lump"))
			Expect(dbResult.Data).To(Equal("50000000000000000000"))
		})

		It("persists a chop lump event (chop)", func() {
			chopLumpBlockNumber := int64(8928383)
			header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
			Expect(err).NotTo(HaveOccurred())
			catFileConfig.TransformerName = constants.CatFileChopLumpDunkTable
			catFileConfig.Topic = constants.CatFileChopLumpSignature()
			catFileConfig.StartingBlockNumber = chopLumpBlockNumber
			catFileConfig.EndingBlockNumber = chopLumpBlockNumber

			initializer := event.ConfiguredTransformer{
				Config:      catFileConfig,
				Transformer: chop_lump_dunk.Transformer{},
			}
			transformer := initializer.NewTransformer(db)

			logs, err := logFetcher.FetchLogs(
				[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
				[]common.Hash{common.HexToHash(catFileConfig.Topic)},
				header)
			Expect(err).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			err = transformer.Execute(eventLogs)
			Expect(err).NotTo(HaveOccurred())

			var dbResult catFileChopLumpModel
			err = db.Get(&dbResult, `SELECT what, ilk_id, data FROM maker.cat_file_chop_lump_dunk`)
			Expect(err).NotTo(HaveOccurred())

			ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
			Expect(dbResult.What).To(Equal("chop"))
			Expect(dbResult.Data).To(Equal("1130000000000000000000000000"))
		})

		It("persists a flip event", func() {
			flipBlockNumber := int64(8928180)
			header, err := persistHeader(db, flipBlockNumber, blockChain)
			Expect(err).NotTo(HaveOccurred())
			catFileConfig.TransformerName = constants.CatFileFlipTable
			catFileConfig.Topic = constants.CatFileFlipSignature()
			catFileConfig.StartingBlockNumber = flipBlockNumber
			catFileConfig.EndingBlockNumber = flipBlockNumber

			initializer := event.ConfiguredTransformer{
				Config:      catFileConfig,
				Transformer: flip.Transformer{},
			}

			t := initializer.NewTransformer(db)

			logs, err := logFetcher.FetchLogs(
				[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
				[]common.Hash{common.HexToHash(catFileConfig.Topic)},
				header)
			Expect(err).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			err = t.Execute(eventLogs)
			Expect(err).NotTo(HaveOccurred())

			var dbResult catFileFlipModel
			err = db.Get(&dbResult, `SELECT ilk_id, what, flip FROM maker.cat_file_flip`)
			Expect(err).NotTo(HaveOccurred())

			ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
			Expect(dbResult.What).To(Equal("flip"))
			Expect(dbResult.Flip).To(Equal(test_data.FlipEthV100Address()))
		})

		It("persists a vow event", func() {
			vowBlockNumber := int64(8928165)
			header, err := persistHeader(db, vowBlockNumber, blockChain)
			Expect(err).NotTo(HaveOccurred())
			catFileConfig.TransformerName = constants.CatFileVowTable
			catFileConfig.Topic = constants.CatFileVowSignature()
			catFileConfig.StartingBlockNumber = vowBlockNumber
			catFileConfig.EndingBlockNumber = vowBlockNumber

			initializer := event.ConfiguredTransformer{
				Config:      catFileConfig,
				Transformer: vow.Transformer{},
			}
			t := initializer.NewTransformer(db)

			logs, err := logFetcher.FetchLogs(
				[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
				[]common.Hash{common.HexToHash(catFileConfig.Topic)},
				header)
			Expect(err).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			err = t.Execute(eventLogs)
			Expect(err).NotTo(HaveOccurred())

			var headerID int64
			err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, vowBlockNumber)
			Expect(err).NotTo(HaveOccurred())

			var dbResult catFileVowModel
			err = db.Get(&dbResult, `SELECT what, data FROM maker.cat_file_vow`)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.What).To(Equal("vow"))
			Expect(dbResult.Data).To(Equal(test_data.VowAddress()))
		})
	})

	Describe("1.1.0 Cat Contract", func() {
		var catFileConfig = event.TransformerConfig{
			ContractAddresses: []string{test_data.Cat110Address()},
			ContractAbi:       constants.Cat110ABI(),
		}

		It("persists a box event", func() {
			boxBlockNumber := int64(10769102)
			header, err := persistHeader(db, boxBlockNumber, blockChain)
			Expect(err).NotTo(HaveOccurred())
			catFileConfig.TransformerName = constants.CatFileBoxTable
			catFileConfig.Topic = constants.CatFileBoxSignature()
			catFileConfig.StartingBlockNumber = boxBlockNumber
			catFileConfig.EndingBlockNumber = boxBlockNumber

			initializer := event.ConfiguredTransformer{
				Config:      catFileConfig,
				Transformer: box.Transformer{},
			}
			transformer := initializer.NewTransformer(db)

			logs, err := logFetcher.FetchLogs(
				[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
				[]common.Hash{common.HexToHash(catFileConfig.Topic)},
				header)
			Expect(err).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			err = transformer.Execute(eventLogs)
			Expect(err).NotTo(HaveOccurred())

			var dbResult catFileBoxModel
			err = db.Get(&dbResult, `SELECT address_id, msg_sender, what, data FROM maker.cat_file_box`)
			Expect(err).NotTo(HaveOccurred())

			addressID, err := shared.GetOrCreateAddress(test_data.Cat110Address(), db)
			Expect(err).NotTo(HaveOccurred())
			msgSenderID, err := shared.GetOrCreateAddress("0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB", db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.AddressID).To(Equal(addressID))
			Expect(dbResult.MsgSenderID).To(Equal(msgSenderID))
			Expect(dbResult.What).To(Equal("box"))
			Expect(dbResult.Data).To(Equal("30000000000000000000000000000000000000000000000000000"))
		})

		It("persists a chop lump dunk event (dunk)", func() {
			chopDunkBlockNumber := int64(10769102)
			header, err := persistHeader(db, chopDunkBlockNumber, blockChain)
			Expect(err).NotTo(HaveOccurred())
			catFileConfig.TransformerName = constants.CatFileChopLumpDunkTable
			catFileConfig.Topic = constants.CatFileChopDunkSignature()
			catFileConfig.StartingBlockNumber = chopDunkBlockNumber
			catFileConfig.EndingBlockNumber = chopDunkBlockNumber

			initializer := event.ConfiguredTransformer{
				Config:      catFileConfig,
				Transformer: chop_lump_dunk.Transformer{},
			}
			transformer := initializer.NewTransformer(db)

			logs, err := logFetcher.FetchLogs(
				[]common.Address{common.HexToAddress(catFileConfig.ContractAddresses[0])},
				[]common.Hash{common.HexToHash(catFileConfig.Topic)},
				header)
			Expect(err).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			err = transformer.Execute(eventLogs)
			Expect(err).NotTo(HaveOccurred())

			var dbResults []catFileChopLumpDunkModel
			msgSenderID, err := shared.GetOrCreateAddress("0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB", db)
			Expect(err).NotTo(HaveOccurred())
			err = db.Select(&dbResults, `SELECT address_id, msg_sender, what, data FROM maker.cat_file_chop_lump_dunk`)
			Expect(err).NotTo(HaveOccurred())
			addressID, err := shared.GetOrCreateAddress(test_data.Cat110Address(), db)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbResults[0].AddressID).To(Equal(addressID))
			Expect(dbResults[0].MsgSenderID).To(Equal(msgSenderID))
			Expect(dbResults[0].What).Should(Or(Equal("dunk"), Equal("chop")))
			Expect(dbResults[0].Data).Should(Or(Equal("1130000000000000000"),
				Equal("50000000000000000000000000000000000000000000000000")))
			Expect(len(dbResults)).To(Equal(18))
		})
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

type catFileBoxModel struct {
	AddressID   int64 `db:"address_id"`
	MsgSenderID int64 `db:"msg_sender"`
	What        string
	Data        string
}
type catFileChopLumpDunkModel struct {
	AddressID   int64 `db:"address_id"`
	MsgSenderID int64 `db:"msg_sender"`
	What        string
	Data        string
}
