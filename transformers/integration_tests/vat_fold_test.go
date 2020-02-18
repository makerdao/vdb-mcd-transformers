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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_fold"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatFold Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	vatFoldConfig := event.TransformerConfig{
		TransformerName:   constants.VatFoldTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatFoldSignature(),
	}

	It("transforms VatFold log events", func() {
		blockNumber := int64(8928674)
		vatFoldConfig.StartingBlockNumber = blockNumber
		vatFoldConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vatFoldConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFoldConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := event.ConfiguredTransformer{
			Config:      vatFoldConfig,
			Transformer: vat_fold.Transformer{},
		}.NewTransformer(db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatFoldModel
		err = db.Select(&dbResults, `SELECT ilk_id, u, rate from maker.vat_fold`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		const ilk = "0x4554482d41000000000000000000000000000000000000000000000000000000"
		ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
		Expect(ilkErr).NotTo(HaveOccurred())
		Expect(dbResult.IlkID).To(Equal(ilkID))
		Expect(dbResult.U).To(Equal("0xA950524441892A31ebddF91d3cEEFa04Bf454466"))
		Expect(dbResult.Rate).To(Equal("5709754190193566796042"))
	})
})

type vatFoldModel struct {
	IlkID int64 `db:"ilk_id"`
	U     string
	Rate  string
}
