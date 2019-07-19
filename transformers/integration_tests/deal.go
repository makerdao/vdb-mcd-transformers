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
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = XDescribe("Deal transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		dealConfig  transformer.EventTransformerConfig
		initializer shared.LogNoteTransformer
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

		dealConfig = transformer.EventTransformerConfig{
			TransformerName: mcdConstants.DealLabel,
			ContractAddresses: []string{
				test_data.FlapperAddress(),
				test_data.OldFlipperAddress(),
				test_data.FlopperAddress(),
			},
			Topic: mcdConstants.DealSignature(),
		}

		initializer = shared.LogNoteTransformer{
			Config:     dealConfig,
			Converter:  &deal.DealConverter{},
			Repository: &deal.DealRepository{},
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = transformer.HexStringsToAddresses(dealConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dealConfig.Topic)}

	})

	It("persists a flip deal log event", func() {
		// transaction: 0x05b5eabac2ace136f0f7e0efc61d7d42abe8e8938cc0f04fbf1a6ba545d59e58
		flipBlockNumber := int64(8958007)
		header, err := persistHeader(db, flipBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = flipBlockNumber
		initializer.Config.EndingBlockNumber = flipBlockNumber

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dealModel
		err = db.Select(&dbResult, `SELECT bid_id, contract_address FROM maker.deal`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].BidId).To(Equal("6"))
		Expect(dbResult[0].ContractAddress).To(Equal(test_data.OldFlipperAddress()))
	})

	It("persists a flop deal log event", func() {
		//TODO: There are currently no Flop.deal events on Kovan
	})

	It("persists a flap deal log event", func() {
		flapBlockNumber := int64(9004628)
		header, err := persistHeader(db, flapBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = flapBlockNumber
		initializer.Config.EndingBlockNumber = flapBlockNumber

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dealModel
		err = db.Select(&dbResult, `SELECT bid_id, contract_address FROM maker.deal`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].ContractAddress).To(Equal(test_data.FlapperAddress()))
	})
})

type dealModel struct {
	BidId            string `db:"bid_id"`
	ContractAddress  string `db:"contract_address"`
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
