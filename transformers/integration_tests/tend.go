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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("Tend EventTransformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		tendConfig  transformer.EventTransformerConfig
		initializer shared.EventTransformer
		logFetcher  fetcher.ILogFetcher
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		tendConfig = transformer.EventTransformerConfig{
			TransformerName:   constants.TendLabel,
			ContractAddresses: []string{test_data.FlapAddress()},
			ContractAbi:       constants.FlapABI(),
			Topic:             constants.TendSignature(),
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = transformer.HexStringsToAddresses(tendConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(tendConfig.Topic)}

		initializer = shared.EventTransformer{
			Config:     tendConfig,
			Converter:  &tend.TendConverter{},
			Repository: &tend.TendRepository{},
		}
	})

	It("fetches and transforms a Tend event from the Kovan chain", func() {
		blockNumber := int64(14308157)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewEventTransformer(db)
		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tendModel
		err = db.Select(&dbResult, `SELECT bid_id, lot, bid FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("193362866030587"))
		Expect(dbResult[0].BidId).To(Equal("27"))
		Expect(dbResult[0].Lot).To(Equal("100000000000000000000000000000000000000000000"))
	})
})

type tendModel struct {
	BidId string `db:"bid_id"`
	Lot   string
	Bid   string
}
