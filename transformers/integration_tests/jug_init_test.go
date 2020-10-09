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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_init"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JugInit EventTransformer", func() {
	jugInitConfig := event.TransformerConfig{
		TransformerName:   constants.JugInitTable,
		ContractAddresses: []string{test_data.JugAddress()},
		ContractAbi:       constants.JugABI(),
		Topic:             constants.JugInitSignature(),
	}

	It("transforms jug init log events", func() {
		blockNumber := int64(8928180)
		jugInitConfig.StartingBlockNumber = blockNumber
		jugInitConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(jugInitConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugInitConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := event.ConfiguredTransformer{
			Config:      jugInitConfig,
			Transformer: jug_init.Transformer{},
		}.NewTransformer(db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult jugInitModel
		err = db.Get(&dbResult, `SELECT ilk_id, msg_sender from maker.jug_init`)
		Expect(err).NotTo(HaveOccurred())

		expectedIlkID, ilkErr := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(ilkErr).NotTo(HaveOccurred())

		expectedMsgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, eventLogs[0].Log.Topics[1].Hex())
		Expect(msgSenderErr).NotTo(HaveOccurred())

		Expect(dbResult.MsgSenderID).To(Equal(expectedMsgSenderID))
		Expect(dbResult.IlkID).To(Equal(expectedIlkID))
	})
})

type jugInitModel struct {
	IlkID       int64 `db:"ilk_id"`
	MsgSenderID int64 `db:"msg_sender"`
}
