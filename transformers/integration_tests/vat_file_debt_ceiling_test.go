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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/debt_ceiling"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatFileDebtCeiling EventTransformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	vatFileDebtCeilingConfig := event.TransformerConfig{
		TransformerName:   constants.VatFileDebtCeilingTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatFileDebtCeilingSignature(),
	}

	It("fetches and transforms a VatFileDebtCeiling event", func() {
		blockNumber := int64(8957370)
		vatFileDebtCeilingConfig.StartingBlockNumber = blockNumber
		vatFileDebtCeilingConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vatFileDebtCeilingConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFileDebtCeilingConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		initializer := event.ConfiguredTransformer{
			Config:      vatFileDebtCeilingConfig,
			Transformer: debt_ceiling.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatFileDebtCeilingModel
		err = db.Select(&dbResult, `SELECT what, data from maker.vat_file_debt_ceiling`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("Line"))
		Expect(dbResult[0].Data).To(Equal("153000000000000000000000000000000000000000000000000000"))
	})
})

type vatFileDebtCeilingModel struct {
	What string
	Data string
}
