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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deal transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		dealConfig  transformer.EventTransformerConfig
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

		dealConfig = transformer.EventTransformerConfig{
			TransformerName: constants.DealLabel,
			ContractAddresses: []string{
				test_data.FlapAddress(),
				test_data.EthFlipAddress(),
				test_data.FlopAddress(),
			},
			Topic: constants.DealSignature(),
		}

		initializer = shared.EventTransformer{
			Config:     dealConfig,
			Converter:  &deal.DealConverter{},
			Repository: &deal.DealRepository{},
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = transformer.HexStringsToAddresses(dealConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dealConfig.Topic)}

	})

	It("persists a flip deal log event", func() {
		flipBlockNumber := int64(14784093)
		header, err := persistHeader(db, flipBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = flipBlockNumber
		initializer.Config.EndingBlockNumber = flipBlockNumber

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewEventTransformer(db)
		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dealModel
		err = db.Select(&dbResult, `SELECT bid_id, address_id FROM maker.deal`)
		Expect(err).NotTo(HaveOccurred())

		flipContractAddressId, err := shared.GetOrCreateAddress(test_data.EthFlipAddress(), db)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].AddressId).To(Equal(flipContractAddressId))
	})

	It("persists a flop deal log event", func() {
		//TODO: There are currently no Flop.deal events on Kovan
	})

	It("persists a flap deal log event", func() {
		//TODO: There are currently no Flap.deal events on Kovan
	})
})

type dealModel struct {
	BidId     string `db:"bid_id"`
	AddressId int64  `db:"address_id"`
}
