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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogValue Transformer", func() {
	logValueConfig := event.TransformerConfig{
		TransformerName:   constants.LogValueTable,
		ContractAddresses: test_data.OsmAddresses(),
		ContractAbi:       constants.OsmABI(),
		Topic:             constants.LogValueSignature(),
	}

	It("fetches and transforms a LogValue event", func() {
		blockNumber := int64(9290757)
		logValueConfig.StartingBlockNumber = blockNumber
		logValueConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      logValueConfig,
			Transformer: log_value.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		var addresses []common.Address
		for _, addr := range logValueConfig.ContractAddresses {
			addresses = append(addresses, common.HexToAddress(addr))
		}
		logs, err := logFetcher.FetchLogs(
			addresses,
			[]common.Hash{common.HexToHash(logValueConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []string
		err = db.Select(&dbResults, `SELECT val from maker.log_value`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(2))
		Expect(dbResults).To(ConsistOf("213417546000000000", "160720000000000000000"))
	})
})
