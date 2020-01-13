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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_suck"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatSuck Transformer", func() {
	vatSuckConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatSuckTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatSuckSignature(),
	}

	It("transforms VatSuck log events", func() {
		blockNumber := int64(14911513)
		vatSuckConfig.StartingBlockNumber = blockNumber
		vatSuckConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatSuckConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatSuckConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      vatSuckConfig,
			Transformer: vat_suck.Transformer{},
		}.NewTransformer(db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatSuckModel
		err = db.Select(&dbResults, `SELECT u, v, rad from maker.vat_suck`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.U).To(Equal("0x0F4Cbe6CBA918b7488C26E29d9ECd7368F38EA3b"))
		Expect(dbResult.V).To(Equal("0xEA190DBDC7adF265260ec4dA6e9675Fd4f5A78bb"))
		Expect(dbResult.Rad).To(Equal("9186142869752396102379486822757141551229585"))
	})
})

type vatSuckModel struct {
	U   string
	V   string
	Rad string
}
