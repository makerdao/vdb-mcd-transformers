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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/tick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	mcdConstants "github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tick EventTransformer", func() {
	var (
		tickConfig  transformer.EventTransformerConfig
		initializer event.Transformer
		logFetcher  fetcher.ILogFetcher
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		tickConfig = transformer.EventTransformerConfig{
			TransformerName:   mcdConstants.TickTable,
			ContractAddresses: append(test_data.FlipAddresses(), test_data.FlopAddress()),
			ContractAbi:       mcdConstants.FlipABI(),
			Topic:             mcdConstants.TickSignature(),
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = transformer.HexStringsToAddresses(tickConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(tickConfig.Topic)}

		initializer = event.Transformer{
			Config:    tickConfig,
			Converter: tick.Converter{},
		}
	})

	It("fetches and transforms a flip tick event", func() {
		blockNumber := int64(8974350)
		tickConfig.StartingBlockNumber = blockNumber
		tickConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())
		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tickModel
		err = db.Select(&dbResult, `SELECT bid_id, address_id FROM maker.tick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		flipContractAddressId, err := shared.GetOrCreateAddress(test_data.EthFlipAddress(), db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].AddressID).To(Equal(flipContractAddressId))
		Expect(dbResult[0].BidId).To(Equal("15"))
	})

	// Todo: fill this in with flap tick event data from kovan
	XIt("fetches and transforms a flap tick event from the Kovan chain", func() {
		blockNumber := int64(8935601)
		tickConfig.StartingBlockNumber = blockNumber
		tickConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())
		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tickModel
		err = db.Select(&dbResult, `SELECT bid_id, address_id FROM maker.tick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].AddressID).To(Equal(""))
		Expect(dbResult[0].BidId).To(Equal(""))
	})
})

type tickModel struct {
	BidId            string `db:"bid_id"`
	AddressID        int64  `db:"address_id"`
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
