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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_join"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PotJoin Transformer", func() {
	potJoinConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.PotJoinTable,
		ContractAddresses: []string{test_data.PotAddress()},
		ContractAbi:       constants.PotABI(),
		Topic:             constants.PotJoinSignature(),
	}

	It("transforms PotJoin log events", func() {
		blockNumber := int64(15485181)
		potJoinConfig.StartingBlockNumber = blockNumber
		potJoinConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, clientErr := getClients(ipc)
		Expect(clientErr).NotTo(HaveOccurred())
		blockChain, blockChainErr := getBlockChain(rpcClient, ethClient)
		Expect(blockChainErr).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, fetchErr := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(potJoinConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(potJoinConfig.Topic)},
			header)
		Expect(fetchErr).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := event.Transformer{
			Config:    potJoinConfig,
			Converter: pot_join.Converter{},
		}.NewTransformer(db)

		transformErr := tr.Execute(headerSyncLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult potJoinModel
		queryErr := db.Get(&dbResult, `SELECT wad from maker.pot_join`)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(dbResult.Wad).To(Equal("4719670301595647258"))
	})
})

type potJoinModel struct {
	Wad string
}
