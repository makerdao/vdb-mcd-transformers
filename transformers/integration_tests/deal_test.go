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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deal transformer", func() {
	var (
		dealConfig  event.TransformerConfig
		initializer event.ConfiguredTransformer
		logFetcher  fetcher.ILogFetcher
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		dealConfig = event.TransformerConfig{
			TransformerName: constants.DealTable,
			ContractAddresses: []string{
				test_data.FlapV100Address(),
				test_data.FlipEthV100Address(),
				test_data.FlopV101Address(),
			},
			Topic: constants.DealSignature(),
		}

		initializer = event.ConfiguredTransformer{
			Config:      dealConfig,
			Transformer: deal.Transformer{},
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = event.HexStringsToAddresses(dealConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dealConfig.Topic)}
	})

	It("persists a flip deal log event", func() {
		flipBlockNumber := int64(8997455)
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

		var dbResult dealModel
		err := db.Get(&dbResult, `SELECT bid_id, address_id, msg_sender FROM maker.deal`)
		Expect(err).NotTo(HaveOccurred())

		flipAddressID, flipAddressErr := shared.GetOrCreateAddress(test_data.FlipEthV100Address(), db)
		Expect(flipAddressErr).NotTo(HaveOccurred())
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress("0x00aBe7471ec9b6953A3BD0ed3C06c46F29Aa4280", db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		Expect(dbResult.BidID).To(Equal("115"))
		Expect(dbResult.AddressID).To(Equal(flipAddressID))
		Expect(dbResult.MsgSenderID).To(Equal(msgSenderID))
	})

	It("persists a flop deal log event", func() {
		blockNumber := int64(9763527)
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

		var dbResult dealModel
		err := db.Get(&dbResult, `SELECT bid_id, address_id, msg_sender FROM maker.deal`)
		Expect(err).NotTo(HaveOccurred())

		flopAddressID, addressErr := shared.GetOrCreateAddress(test_data.FlopV101Address(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress("0x06C36BEA54A74dB813Af0fc136c2E8d0B08e2FB1", db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		Expect(dbResult.BidID).To(Equal("102"))
		Expect(dbResult.AddressID).To(Equal(flopAddressID))
		Expect(dbResult.MsgSenderID).To(Equal(msgSenderID))
	})

	It("persists a flap deal log event", func() {
		blockNumber := int64(9635404)
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

		var dbResult dealModel
		err := db.Get(&dbResult, `SELECT bid_id, address_id, msg_sender FROM maker.deal ORDER BY log_id`)
		Expect(err).NotTo(HaveOccurred())

		flapAddressID, addressErr := shared.GetOrCreateAddress(test_data.FlapV100Address(), db)
		Expect(addressErr).NotTo(HaveOccurred())
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress("0xFDc7768e92B479F27dD11635c9207d542177ae72", db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		Expect(dbResult.BidID).To(Equal("48"))
		Expect(dbResult.AddressID).To(Equal(flapAddressID))
		Expect(dbResult.MsgSenderID).To(Equal(msgSenderID))
	})
})

type dealModel struct {
	BidID       string `db:"bid_id"`
	AddressID   int64  `db:"address_id"`
	MsgSenderID int64  `db:"msg_sender"`
}
