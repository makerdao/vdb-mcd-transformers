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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_address"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VowFileAuctionAddress LogNoteTransformer", func() {
	var (
		initializer event.ConfiguredTransformer
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		config := event.TransformerConfig{
			TransformerName:   constants.VowFileAuctionAddressTable,
			ContractAddresses: []string{test_data.VowAddress()},
			ContractAbi:       constants.VowABI(),
			Topic:             constants.VowFileAuctionAddressSignature(),
		}

		addresses = event.HexStringsToAddresses(config.ContractAddresses)
		topics = []common.Hash{common.HexToHash(config.Topic)}

		initializer = event.ConfiguredTransformer{
			Config:      config,
			Transformer: auction_address.Transformer{},
		}
	})

	It("fetches and transforms a Vow.file auction address event", func() {
		blockNumber := int64(9017707)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := initializer.NewTransformer(db)
		executeErr := tr.Execute(eventLogs)
		Expect(executeErr).NotTo(HaveOccurred())

		var dbResult []vowFileAuctionAddressModel
		getVowFileErr := db.Select(&dbResult, `SELECT what, data from maker.vow_file_auction_address`)
		Expect(getVowFileErr).NotTo(HaveOccurred())

		var addressId int64
		address := common.HexToAddress("0x4d95a049d5b0b7d32058cd3f2163015747522e99").Hex()
		getAddressIdErr := db.Get(&addressId, `SELECT id FROM public.addresses WHERE address = $1`, address)
		Expect(getAddressIdErr).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("flopper"))
		Expect(dbResult[0].Data).To(Equal(addressId))
	})
})

type vowFileAuctionAddressModel struct {
	What string
	Data int64
}
