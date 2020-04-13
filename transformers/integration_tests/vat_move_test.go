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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("VatMove EventTransformer", func() {
	vatMoveConfig := event.TransformerConfig{
		TransformerName:   constants.VatMoveTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatMoveSignature(),
	}

	It("transforms VatMove log events", func() {
		blockNumber := int64(9284211)
		vatMoveConfig.StartingBlockNumber = blockNumber
		vatMoveConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vatMoveConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatMoveConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.ConfiguredTransformer{
			Config:      vatMoveConfig,
			Transformer: vat_move.Transformer{},
		}.NewTransformer(db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatMoveModel
		err = db.Select(&dbResults, `SELECT src, dst, rad from maker.vat_move`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(4))
		sort.Sort(byRad(dbResults))
		dbResult := dbResults[2]
		Expect(dbResult.Src).To(Equal("0x9759A6Ac90977b93B58547b4A71c78317f391A28"))
		Expect(dbResult.Dst).To(Equal("0x5d3a536E4D6DbD6114cc1Ead35777bAB948E3643"))
		Expect(dbResult.Rad).To(Equal("8170620340000000000000000000000000000000000000000"))
	})
})

type byRad []vatMoveModel

func (b byRad) Len() int {
	return len(b)
}

func (b byRad) Less(i, j int) bool {
	return b[i].Rad < b[j].Rad
}

func (b byRad) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type vatMoveModel struct {
	Src string
	Dst string
	Rad string
}
