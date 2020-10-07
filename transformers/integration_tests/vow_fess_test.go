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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_fess"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VowFess EventTransformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	vowFessConfig := event.TransformerConfig{
		TransformerName:   constants.VowFessTable,
		ContractAddresses: []string{test_data.VowAddress()},
		ContractAbi:       constants.VowABI(),
		Topic:             constants.VowFessSignature(),
	}

	It("transforms VowFess log events", func() {
		blockNumber := int64(8997324)
		vowFessConfig.StartingBlockNumber = blockNumber
		vowFessConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vowFessConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vowFessConfig.Topic)},
			header)
		Expect(len(logs)).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      vowFessConfig,
			Transformer: vow_fess.Transformer{},
		}.NewTransformer(db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		msgSenderID, msgSenderErr := shared.GetOrCreateAddress("0x78F2c2AF65126834c51822F56Be0d7469D7A523E", db)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		var dbResult vowFessModel
		err = db.Get(&dbResult, `SELECT tab, msg_sender from maker.vow_fess`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult.Tab).To(Equal("4466031366353941646208178591268931635087392443453"))
		Expect(dbResult.MsgSenderID).To(Equal(msgSenderID))
	})
})

type vowFessModel struct {
	Tab         string `db:"tab"`
	MsgSenderID int64  `db:"msg_sender"`
}
