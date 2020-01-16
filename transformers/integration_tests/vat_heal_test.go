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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_heal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatHeal Transformer", func() {
	vatHealConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatHealTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatHealSignature(),
	}

	It("transforms VatHeal log events", func() {
		blockNumber := int64(8957867)
		vatHealConfig.StartingBlockNumber = blockNumber
		vatHealConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatHealConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatHealConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      vatHealConfig,
			Transformer: vat_heal.Transformer{},
		}.NewTransformer(db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult vatHealModel
		err = db.Get(&dbResult, `SELECT rad from maker.vat_heal`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult.Rad).To(Equal("331024328758114038198514075052115218961684178"))
	})
})

type vatHealModel struct {
	Rad string
}
