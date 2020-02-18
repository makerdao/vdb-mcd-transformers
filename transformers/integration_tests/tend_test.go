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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/tend"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tend EventTransformer", func() {
	var (
		tendConfig  event.TransformerConfig
		initializer event.ConfiguredTransformer
		logFetcher  fetcher.ILogFetcher
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		tendConfig = event.TransformerConfig{
			TransformerName:   constants.TendTable,
			ContractAddresses: []string{test_data.EthFlipAddress()},
			ContractAbi:       constants.FlipABI(),
			Topic:             constants.TendSignature(),
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = event.HexStringsToAddresses(tendConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(tendConfig.Topic)}

		initializer = event.ConfiguredTransformer{
			Config:      tendConfig,
			Transformer: tend.Transformer{},
		}
	})

	It("fetches and transforms a Flip Tend event", func() {
		blockNumber := int64(9004844)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tendModel
		err = db.Select(&dbResult, `SELECT bid_id, lot, bid FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("76840636079422693500873675445736719538580144543"))
		Expect(dbResult[0].BidId).To(Equal("121"))
		Expect(dbResult[0].Lot).To(Equal("700000000000000000"))
	})

	//TODO: There are currently no Flap Tend events
	It("fetches and transforms a Flap Tend event", func() {})
})

type tendModel struct {
	BidId string `db:"bid_id"`
	Lot   string
	Bid   string
}
