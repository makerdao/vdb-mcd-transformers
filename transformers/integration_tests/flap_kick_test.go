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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
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

	It("fetches and transforms a FlapKick event", func() {
		blockNumber := int64(9600250)
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

		var dbResult FlapKickModel
		err = db.Get(&dbResult, `SELECT address_id, bid, bid_id, lot FROM maker.flap_kick`)
		Expect(err).NotTo(HaveOccurred())

		flapAddressID, addressErr := shared.GetOrCreateAddress(test_data.FlapAddress(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedModel := FlapKickModel{
			AddressID: flapAddressID,
			Bid:       "0",
			BidId:     "39",
			Lot:       "10000000000000000000000000000000000000000000000000",
		}
		Expect(dbResult).To(Equal(expectedModel))
	})
})

type FlapKickModel struct {
	AddressID int64 `db:"address_id"`
	Bid       string
	BidId     string `db:"bid_id"`
	Lot       string
}
