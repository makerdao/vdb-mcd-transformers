//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package integration_tests

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_poke"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpotPoke Transformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	spotPokeConfig := event.TransformerConfig{
		TransformerName:   constants.SpotPokeTable,
		ContractAddresses: []string{test_data.SpotAddress()},
		ContractAbi:       constants.SpotABI(),
		Topic:             constants.SpotPokeSignature(),
	}

	It("transforms spot poke log events", func() {
		blockNumber := int64(8928374)
		spotPokeConfig.StartingBlockNumber = blockNumber
		spotPokeConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      spotPokeConfig,
			Transformer: spot_poke.Transformer{},
		}
		tr := initializer.NewTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			event.HexStringsToAddresses(spotPokeConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(spotPokeConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []spotPokeModel
		err = db.Select(&dbResult, `SELECT ilk_id, value, spot FROM maker.spot_poke`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("5341490000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult[0].Value).To(Equal("1000000000000000000.000000"))
		Expect(dbResult[0].Spot).To(Equal("100000000000000000000000000000000000000000000000000"))
	})
})

type spotPokeModel struct {
	Ilk   string `db:"ilk_id"`
	Value string
	Spot  string
}
