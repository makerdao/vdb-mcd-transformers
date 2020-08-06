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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_flog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VowFlog EventTransformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	vowFlogConfig := event.TransformerConfig{
		TransformerName:   constants.VowFlogTable,
		ContractAddresses: []string{test_data.VowAddress()},
		ContractAbi:       constants.VowABI(),
		Topic:             constants.VowFlogSignature(),
	}

	It("transforms VowFlog log events", func() {
		blockNumber := int64(9242813)
		vowFlogConfig.StartingBlockNumber = blockNumber
		vowFlogConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vowFlogConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vowFlogConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      vowFlogConfig,
			Transformer: vow_flog.Transformer{},
		}.NewTransformer(db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		msgSender := shared.GetChecksumAddressString("0x00000000000000000000000022e86ab483084053562ce713e94431c29d1adb8b")
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(msgSender, db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		var dbResult vowFlogModel
		err = db.Get(&dbResult, `SELECT msg_sender, era from maker.vow_flog`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult.MsgSender).To(Equal(msgSenderID))
		Expect(dbResult.Era).To(Equal("1577965150"))
	})
})

type vowFlogModel struct {
	MsgSender int64 `db:"msg_sender"`
	Era       string
}
