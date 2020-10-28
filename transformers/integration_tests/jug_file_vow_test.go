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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jug File Vow EventTransformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	jugFileVowConfig := event.TransformerConfig{
		TransformerName:   constants.JugFileVowTable,
		ContractAddresses: []string{test_data.JugAddress()},
		ContractAbi:       constants.JugABI(),
		Topic:             constants.JugFileVowSignature(),
	}

	It("transforms JugFileVow log events", func() {
		blockNumber := int64(8928163)
		jugFileVowConfig.StartingBlockNumber = blockNumber
		jugFileVowConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      jugFileVowConfig,
			Transformer: vow.Transformer{},
		}
		tr := initializer.NewTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			event.HexStringsToAddresses(jugFileVowConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugFileVowConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult jugFileVowModel
		err = db.Get(&dbResult, `SELECT msg_sender, what, data FROM maker.jug_file_vow`)
		Expect(err).NotTo(HaveOccurred())

		msgSender := shared.GetChecksumAddressString("0x000000000000000000000000baa65281c2fa2baacb2cb550ba051525a480d3f4")
		msgSenderId, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
		Expect(msgSenderErr).NotTo(HaveOccurred())

		Expect(dbResult.MsgSender).To(Equal(msgSenderId))
		Expect(dbResult.What).To(Equal("vow"))
		Expect(dbResult.Data).To(Equal(test_data.VowAddress()))
	})
})

type jugFileVowModel struct {
	MsgSender int64 `db:"msg_sender"`
	What      string
	Data      string
}
