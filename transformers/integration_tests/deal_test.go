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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deal transformer", func() {
	var (
		dealConfig  transformer.EventTransformerConfig
		initializer event.ConfiguredTransformer
		logFetcher  fetcher.ILogFetcher
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		dealConfig = transformer.EventTransformerConfig{
			TransformerName: constants.DealTable,
			ContractAddresses: []string{
				test_data.FlapAddress(),
				test_data.EthFlipAddress(),
				test_data.FlopAddress(),
			},
			Topic: constants.DealSignature(),
		}

		initializer = event.ConfiguredTransformer{
			Config:      dealConfig,
			Transformer: deal.Transformer{},
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = transformer.HexStringsToAddresses(dealConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dealConfig.Topic)}
	})

	It("persists a flip deal log event", func() {
		flipBlockNumber := int64(14887716)
		header, headerErr := persistHeader(db, flipBlockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = flipBlockNumber
		initializer.Config.EndingBlockNumber = flipBlockNumber

		logs, fetchErr := logFetcher.FetchLogs(addresses, topics, header)
		Expect(fetchErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		transformErr := transformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult []dealModel
		err := db.Select(&dbResult, `SELECT bid_id, address_id FROM maker.deal`)
		Expect(err).NotTo(HaveOccurred())

		flipAddressID, flipAddressErr := shared.GetOrCreateAddress(test_data.EthFlipAddress(), db)
		Expect(flipAddressErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].BidID).To(Equal("15"))
		Expect(dbResult[0].AddressID).To(Equal(flipAddressID))
	})

	It("persists a flop deal log event", func() {
		blockNumber := int64(15788320)
		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, fetchErr := logFetcher.FetchLogs(addresses, topics, header)
		Expect(fetchErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		transformErr := transformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult []dealModel
		err := db.Select(&dbResult, `SELECT bid_id, address_id FROM maker.deal`)
		Expect(err).NotTo(HaveOccurred())

		flopAddressID, addressErr := shared.GetOrCreateAddress(test_data.FlopAddress(), db)
		Expect(addressErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].BidID).To(Equal("989"))
		Expect(dbResult[0].AddressID).To(Equal(flopAddressID))
	})

	It("persists a flap deal log event", func() {
		blockNumber := int64(15240030)
		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, fetchErr := logFetcher.FetchLogs(addresses, topics, header)
		Expect(fetchErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		transformErr := transformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult []dealModel
		err := db.Select(&dbResult, `SELECT bid_id, address_id FROM maker.deal ORDER BY log_id`)
		Expect(err).NotTo(HaveOccurred())

		flapAddressID, addressErr := shared.GetOrCreateAddress(test_data.FlapAddress(), db)
		Expect(addressErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		Expect(dbResult[0].BidID).To(Equal("506"))
		Expect(dbResult[0].AddressID).To(Equal(flapAddressID))
	})
})

type dealModel struct {
	BidID     string `db:"bid_id"`
	AddressID int64  `db:"address_id"`
}
