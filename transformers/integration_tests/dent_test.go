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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dent"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dent transformer", func() {
	var (
		logFetcher  fetcher.ILogFetcher
		tr          transformer.EventTransformer
		dentConfig  transformer.EventTransformerConfig
		addresses   []common.Address
		topics      []common.Hash
		initializer event.ConfiguredTransformer
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		dentConfig = transformer.EventTransformerConfig{
			TransformerName:   constants.DentTable,
			ContractAddresses: append(test_data.FlipAddresses(), test_data.FlopAddress()),
			ContractAbi:       constants.FlipABI(),
			Topic:             constants.DentSignature(),
		}

		addresses = transformer.HexStringsToAddresses(dentConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dentConfig.Topic)}
		logFetcher = fetcher.NewLogFetcher(blockChain)

		initializer = event.ConfiguredTransformer{
			Config:      dentConfig,
			Transformer: dent.Transformer{},
		}
	})

	It("persists a flop dent log event", func() {
		blockNumber := int64(15788325)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, fetchErr := logFetcher.FetchLogs(addresses, topics, header)
		Expect(fetchErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr = initializer.NewTransformer(db)
		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dentModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, lot, address_id FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		flipContractAddressId, addressErr := shared.GetOrCreateAddress(test_data.FlopAddress(), db)
		Expect(addressErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("100000000000000000000000000000000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("1000"))
		Expect(dbResult[0].Lot).To(Equal("229531987479152"))
		Expect(dbResult[0].AddressId).To(Equal(flipContractAddressId))
	})

	It("persists a flip dent log event", func() {
		blockNumber := int64(14887563)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr = initializer.NewTransformer(db)
		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dentModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, lot, address_id FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		flipContractAddressId, err := shared.GetOrCreateAddress(test_data.EthFlipAddress(), db)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("22600003035575947920904950765775383978062700848"))
		Expect(dbResult[0].BidId).To(Equal("15"))
		Expect(dbResult[0].Lot).To(Equal("125730197694441991"))
		Expect(dbResult[0].AddressId).To(Equal(flipContractAddressId))
	})
})

type dentModel struct {
	BidId     string `db:"bid_id"`
	Lot       string
	Bid       string
	AddressId int64 `db:"address_id"`
}
