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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dent transformer", func() {
	var (
		logFetcher  fetcher.ILogFetcher
		tr          event.ITransformer
		dentConfig  event.TransformerConfig
		addresses   []common.Address
		topics      []common.Hash
		initializer event.ConfiguredTransformer
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		dentConfig = event.TransformerConfig{
			TransformerName:   constants.DentTable,
			ContractAddresses: append(test_data.FlipV100Addresses(), test_data.FlopV101Address()),
			ContractAbi:       constants.FlipV100ABI(),
			Topic:             constants.DentSignature(),
		}

		addresses = event.HexStringsToAddresses(dentConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dentConfig.Topic)}
		logFetcher = fetcher.NewLogFetcher(blockChain)

		initializer = event.ConfiguredTransformer{
			Config:      dentConfig,
			Transformer: dent.Transformer{},
		}
	})

	It("persists a flop dent log event", func() {
		blockNumber := int64(9758937)
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

		var dbResult dentModel
		err = db.Get(&dbResult, `SELECT bid, bid_id, lot, msg_sender, address_id FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		msgSender := common.HexToAddress("0xe06ac4777f04ac7638f736a0b95f7bfeadcee556").Hex()
		msgSenderId, msgSenderErr := shared.GetOrCreateAddress(msgSender, db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		flopContractAddressId, addressErr := shared.GetOrCreateAddress(test_data.FlopV101Address(), db)
		Expect(addressErr).NotTo(HaveOccurred())

		expectedModel := dentModel{
			BidId:     "90",
			Lot:       "176522506619593998233",
			Bid:       "50000000000000000000000000000000000000000000000000",
			MsgSender: msgSenderId,
			AddressId: flopContractAddressId,
		}
		Expect(dbResult).To(Equal(expectedModel))
	})

	It("persists a flip dent log event", func() {
		blockNumber := int64(9003162)
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

		var dbResult dentModel
		err = db.Get(&dbResult, `SELECT bid, bid_id, lot, msg_sender, address_id FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		flipContractAddressId, err := shared.GetOrCreateAddress(test_data.FlipEthV100Address(), db)
		Expect(err).NotTo(HaveOccurred())

		msgSender := common.HexToAddress("0xabe7471ec9b6953a3bd0ed3c06c46f29aa4280").Hex()
		msgSenderId, msgSenderErr := shared.GetOrCreateAddress(msgSender, db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		Expect(dbResult.Bid).To(Equal("111871106928171434728687324748784117143125320430"))
		Expect(dbResult.BidId).To(Equal("119"))
		Expect(dbResult.Lot).To(Equal("903984178994823415"))
		Expect(dbResult.MsgSender).To(Equal(msgSenderId))
		Expect(dbResult.AddressId).To(Equal(flipContractAddressId))
	})
})

type dentModel struct {
	BidId     string `db:"bid_id"`
	Lot       string
	Bid       string
	AddressId int64 `db:"address_id"`
	MsgSender int64 `db:"msg_sender"`
}
