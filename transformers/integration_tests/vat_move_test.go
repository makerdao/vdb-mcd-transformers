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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_move"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatMove EventTransformer", func() {
	vatMoveConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatMoveTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatMoveSignature(),
	}

	It("transforms VatMove log events", func() {
		blockNumber := int64(9046898)
		vatMoveConfig.StartingBlockNumber = blockNumber
		vatMoveConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatMoveConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatMoveConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.Transformer{
			Config:    vatMoveConfig,
			Converter: vat_move.Converter{},
		}.NewTransformer(db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatMoveModel
		err = db.Select(&dbResults, `SELECT src, dst, rad from maker.vat_move ORDER BY rad ASC LIMIT 1`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.Src).To(Equal("0xc73e0383F3Aff3215E6f04B0331D58CeCf0Ab849"))
		Expect(dbResult.Dst).To(Equal("0x9759A6Ac90977b93B58547b4A71c78317f391A28"))
		Expect(dbResult.Rad).To(Equal("14967448188835982956000000000000000000000000000"))
	})
})

type vatMoveModel struct {
	Src string
	Dst string
	Rad string
}
