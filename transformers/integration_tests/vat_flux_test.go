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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_flux"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatFlux EventTransformer", func() {
	vatFluxConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatFluxTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatFluxSignature(),
	}

	It("transforms VatFlux log events", func() {
		blockNumber := int64(9046454)
		vatFluxConfig.StartingBlockNumber = blockNumber
		vatFluxConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatFluxConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFluxConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		initializer := event.Transformer{
			Config:    vatFluxConfig,
			Converter: vat_flux.Converter{},
		}
		transformer := initializer.NewTransformer(db)

		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatFluxModel
		err = db.Select(&dbResult, `SELECT ilk_id, src, dst, wad from maker.vat_flux`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult[0].Src).To(Equal("0xF6Aa7E9C9805b683e3c70FFd68e9d0F2E632307d"))
		Expect(dbResult[0].Dst).To(Equal("0xd4E133dD36442CeE6da6438d480558a631325edd"))
		Expect(dbResult[0].Wad).To(Equal("1485170438407164677"))
	})
})

type vatFluxModel struct {
	Ilk string `db:"ilk_id"`
	Src string
	Dst string
	Wad string
}
