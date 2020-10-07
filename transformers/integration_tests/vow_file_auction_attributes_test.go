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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_attributes"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VowFileAuctionAttributes LogNoteTransformer", func() {
	var (
		initializer event.ConfiguredTransformer
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		vowFileAuctionAttributesConfig := event.TransformerConfig{
			TransformerName:   constants.VowFileAuctionAttributesTable,
			ContractAddresses: []string{test_data.VowAddress()},
			ContractAbi:       constants.VowABI(),
			Topic:             constants.VowFileAuctionAttributesSignature(),
		}

		addresses = event.HexStringsToAddresses(vowFileAuctionAttributesConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(vowFileAuctionAttributesConfig.Topic)}

		initializer = event.ConfiguredTransformer{
			Config:      vowFileAuctionAttributesConfig,
			Transformer: auction_attributes.Transformer{},
		}
	})

	It("fetches and transforms a Vow.file auction attributes event", func() {
		blockNumber := int64(8928291)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := initializer.NewTransformer(db)
		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var msgSenderId int64
		msgSenderAddress := common.HexToAddress("0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB").Hex()
		getMsgSenderIdErr := db.Get(&msgSenderId, `SELECT id FROM public.addresses WHERE address = $1`, msgSenderAddress)
		Expect(getMsgSenderIdErr).NotTo(HaveOccurred())

		var dbResult []vowFileAuctionAttributesModel
		err = db.Select(&dbResult, `SELECT what, data, msg_sender from maker.vow_file_auction_attributes`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("sump"))
		Expect(dbResult[0].Data).To(Equal("50000000000000000000000000000000000000000000000000"))
		Expect(dbResult[0].MsgSender).To(Equal(msgSenderId))
	})
})

type vowFileAuctionAttributesModel struct {
	What      string
	Data      string
	MsgSender int64 `db:"msg_sender"`
}
