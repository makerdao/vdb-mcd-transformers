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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlipKick Transformer", func() {
	flipKickConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.FlipKickTable,
		ContractAddresses: []string{test_data.EthFlipAddress()},
		ContractAbi:       constants.FlipABI(),
		Topic:             constants.FlipKickSignature(),
	}

	It("fetches and transforms a FlipKick event from Kovan chain", func() {
		blockNumber := int64(8997383)
		flipKickConfig.StartingBlockNumber = blockNumber
		flipKickConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		tr := event.Transformer{
			Config:    flipKickConfig,
			Converter: flip_kick.Converter{},
		}.NewTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(flipKickConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(flipKickConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []FlipKickModel
		err = db.Select(&dbResult, `SELECT bid_id, lot, bid, tab, usr, gal, address_id FROM maker.flip_kick`)
		Expect(err).NotTo(HaveOccurred())

		flipContractAddressId, err := shared.GetOrCreateAddress(test_data.EthFlipAddress(), db)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].Lot).To(Equal("50000000000000000000"))
		Expect(dbResult[0].Tab).To(Equal("5046619216084543990261356563876808629308883826941"))
		Expect(dbResult[0].Usr).To(Equal("0x0A051CD913dFD1820dbf87a9bf62B04A129F88A5"))
		Expect(dbResult[0].Gal).To(Equal("0xA950524441892A31ebddF91d3cEEFa04Bf454466"))
		Expect(dbResult[0].AddressId).To(Equal(flipContractAddressId))
	})
})

type FlipKickModel struct {
	BidId     string `db:"bid_id"`
	Lot       string
	Bid       string
	Tab       string
	Usr       string
	Gal       string
	AddressId int64 `db:"address_id"`
	HeaderID  int64 `db:"header_id"`
	LogID     int64 `db:"log_id"`
}
