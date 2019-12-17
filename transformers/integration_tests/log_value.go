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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_value"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogValue Transformer", func() {
	logValueConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.LogValueTable,
		ContractAddresses: test_data.OsmAddresses(),
		ContractAbi:       constants.OsmABI(),
		Topic:             constants.LogValueSignature(),
	}

	It("fetches and transforms a LogValue event from Kovan chain", func() {
		blockNumber := int64(14824113)
		logValueConfig.StartingBlockNumber = blockNumber
		logValueConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.Transformer{
			Config:    logValueConfig,
			Converter: log_value.Converter{},
		}
		transformer := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(logValueConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(logValueConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []logValueModel
		err = db.Select(&dbResult, `SELECT val from maker.log_value`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Val).To(Equal("183430000000000000000"))
	})
})

type logValueModel struct {
	Val      string
	HeaderID int64
	LogID    int64 `db:"log_id"`
}
