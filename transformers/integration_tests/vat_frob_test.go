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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat frob Transformer", func() {
	var (
		logFetcher    fetcher.ILogFetcher
		vatFrobConfig event.TransformerConfig
		initializer   event.ConfiguredTransformer
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		logFetcher = fetcher.NewLogFetcher(blockChain)
		vatFrobConfig = event.TransformerConfig{
			TransformerName:   constants.VatFrobTable,
			ContractAddresses: []string{test_data.VatAddress()},
			ContractAbi:       constants.VatABI(),
			Topic:             constants.VatFrobSignature(),
		}

		initializer = event.ConfiguredTransformer{
			Config:      vatFrobConfig,
			Transformer: vat_frob.Transformer{},
		}
	})

	It("fetches and transforms a vat frob event", func() {
		blockNumber := int64(8928903)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vatFrobConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFrobConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatFrobModel
		err = db.Select(&dbResult, `SELECT urn_id, v, w, dink, dart from maker.vat_frob`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		urnID, err := shared.GetOrCreateUrn("0xfC24B0E4A2ff825cd25E85D658b70b51Bf1D0902",
			"0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult[0].Urn).To(Equal(strconv.FormatInt(urnID, 10)))
		Expect(dbResult[0].V).To(Equal("0xfC24B0E4A2ff825cd25E85D658b70b51Bf1D0902"))
		Expect(dbResult[0].W).To(Equal("0xfC24B0E4A2ff825cd25E85D658b70b51Bf1D0902"))
		Expect(dbResult[0].Dink).To(Equal("50000000000000000"))
		Expect(dbResult[0].Dart).To(Equal("0"))
	})
})

type vatFrobModel struct {
	Ilk  string
	Urn  string `db:"urn_id"`
	V    string
	W    string
	Dink string
	Dart string
}
