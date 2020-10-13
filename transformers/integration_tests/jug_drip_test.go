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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_drip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JugDrip Transformer", func() {
	var jugDripConfig event.TransformerConfig

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		jugDripConfig = event.TransformerConfig{
			ContractAddresses: []string{test_data.JugAddress()},
			ContractAbi:       constants.JugABI(),
			Topic:             constants.JugDripSignature(),
		}
	})

	It("transforms JugDrip log events", func() {
		blockNumber := int64(8928358)
		jugDripConfig.StartingBlockNumber = blockNumber
		jugDripConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      jugDripConfig,
			Transformer: jug_drip.Transformer{},
		}
		tr := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(jugDripConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugDripConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var result jugDripModel
		err = db.Get(&result, `SELECT msg_sender, ilk_id from maker.jug_drip`)
		Expect(err).NotTo(HaveOccurred())

		msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB")
		Expect(msgSenderErr).NotTo(HaveOccurred())
		ilkID, ilkErr := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(ilkErr).NotTo(HaveOccurred())
		expectedResult := jugDripModel{
			IlkID:     ilkID,
			MsgSender: msgSenderID,
		}

		Expect(result).To(Equal(expectedResult))
	})
})

type jugDripModel struct {
	IlkID     int64 `db:"ilk_id"`
	MsgSender int64 `db:"msg_sender"`
}
