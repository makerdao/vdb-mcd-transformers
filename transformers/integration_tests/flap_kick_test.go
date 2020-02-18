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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flap_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlapKick Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	flapKickConfig := event.TransformerConfig{
		TransformerName:   constants.FlapKickTable,
		ContractAddresses: []string{test_data.FlapAddress()},
		ContractAbi:       constants.FlapABI(),
		Topic:             constants.FlapKickSignature(),
	}
	//TODO: There are no flap kick events on mainnet
	XIt("fetches and transforms a FlapKick event", func() {
		blockNumber := int64(14883695)
		flapKickConfig.StartingBlockNumber = blockNumber
		flapKickConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		tr := event.ConfiguredTransformer{
			Config:      flapKickConfig,
			Transformer: flap_kick.Transformer{},
		}.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(flapKickConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(flapKickConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []FlapKickModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, lot FROM maker.flap_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("46"))
		Expect(dbResult[0].Lot).To(Equal("100000000000000000000000000000000000000000000"))
	})
})

type FlapKickModel struct {
	BidId    string `db:"bid_id"`
	Lot      string
	Bid      string
	HeaderID int64 `db:"header_id"`
	LogID    int64 `db:"log_id"`
}
