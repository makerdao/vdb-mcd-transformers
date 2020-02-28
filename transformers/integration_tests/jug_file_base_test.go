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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/base"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jug File Base EventTransformer", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	jugFileBaseConfig := event.TransformerConfig{
		TransformerName:   constants.JugFileBaseTable,
		ContractAddresses: []string{test_data.JugAddress()},
		ContractAbi:       constants.JugABI(),
		Topic:             constants.JugFileBaseSignature(),
	}

	It("transforms jug file base log events", func() {
		blockNumber := int64(8928298)
		jugFileBaseConfig.StartingBlockNumber = blockNumber
		jugFileBaseConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      jugFileBaseConfig,
			Transformer: base.Transformer{},
		}
		tr := initializer.NewTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			event.HexStringsToAddresses(jugFileBaseConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugFileBaseConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []jugFileBaseModel
		err = db.Select(&dbResult, `SELECT what, data FROM maker.jug_file_base`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("base"))
		Expect(dbResult[0].Data).To(Equal("0"))
	})
})

type jugFileBaseModel struct {
	What string
	Data string
}
