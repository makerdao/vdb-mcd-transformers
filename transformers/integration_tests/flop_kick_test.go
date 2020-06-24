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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlopKick Transformer", func() {
	var (
		flopKickConfig event.TransformerConfig
		initializer    event.ConfiguredTransformer
		logFetcher     fetcher.ILogFetcher
		addresses      []common.Address
		topics         []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		flopKickConfig = event.TransformerConfig{
			TransformerName:   constants.FlopKickTable,
			ContractAddresses: []string{test_data.FlopAddress()},
			ContractAbi:       constants.FlopABI(),
			Topic:             constants.FlopKickSignature(),
		}

		initializer = event.ConfiguredTransformer{
			Config:      flopKickConfig,
			Transformer: flop_kick.Transformer{},
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = event.HexStringsToAddresses(flopKickConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(flopKickConfig.Topic)}
	})

	It("fetches and transforms a FlopKick event", func() {
		blockNumber := int64(9724477)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		logs, fetchErr := logFetcher.FetchLogs(addresses, topics, header)
		Expect(fetchErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		transformErr := transformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult FlopKickModel
		err := db.Get(&dbResult, `SELECT bid, bid_id, address_id, gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		flopAddressID, flopAddressErr := shared.GetOrCreateAddress(test_data.FlopAddress(), db)
		Expect(flopAddressErr).NotTo(HaveOccurred())

		expectedModel := FlopKickModel{
			BidId:     "86",
			Lot:       "250000000000000000000",
			Bid:       "50000000000000000000000000000000000000000000000000",
			Gal:       "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
			AddressID: strconv.FormatInt(flopAddressID, 10),
		}
		Expect(dbResult).To(Equal(expectedModel))
	})
})

type FlopKickModel struct {
	BidId     string `db:"bid_id"`
	Lot       string
	Bid       string
	Gal       string
	AddressID string `db:"address_id"`
}
