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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_join"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotJoin Transformer", func() {
	potJoinConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.PotJoinTable,
		ContractAddresses: []string{test_data.PotAddress()},
		ContractAbi:       constants.PotABI(),
		Topic:             constants.PotJoinSignature(),
	}

	It("transforms PotJoin log events", func() {
		blockNumber := int64(9131808)
		potJoinConfig.StartingBlockNumber = blockNumber
		potJoinConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, fetchErr := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(potJoinConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(potJoinConfig.Topic)},
			header)
		Expect(fetchErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      potJoinConfig,
			Transformer: pot_join.Transformer{},
		}.NewTransformer(db)

		transformErr := tr.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult potJoinModel
		queryErr := db.Get(&dbResult, `SELECT msg_sender, wad from maker.pot_join`)
		Expect(queryErr).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress("0x06AF07097C9Eeb7fD685c692751D5C66dB49c215", db)
		Expect(addressErr).NotTo(HaveOccurred())
		Expect(dbResult.MsgSender).To(Equal(strconv.FormatInt(addressID, 10)))
		Expect(dbResult.Wad).To(Equal("1065452848025509462649"))
	})
})

type potJoinModel struct {
	MsgSender string `db:"msg_sender"`
	Wad       string
}
