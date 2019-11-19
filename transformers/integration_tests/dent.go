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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dent"
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

var _ = XDescribe("Dent transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		logFetcher  fetcher.ILogFetcher
		tr          transformer.EventTransformer
		dentConfig  transformer.EventTransformerConfig
		addresses   []common.Address
		topics      []common.Hash
		initializer shared.EventTransformer
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		dentConfig = transformer.EventTransformerConfig{
			TransformerName:   constants.DentLabel,
			ContractAddresses: append(test_data.FlipAddresses(), test_data.FlopAddress()),
			ContractAbi:       constants.FlipABI(),
			Topic:             constants.DentSignature(),
		}

		addresses = transformer.HexStringsToAddresses(dentConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(dentConfig.Topic)}
		logFetcher = fetcher.NewLogFetcher(blockChain)

		initializer = shared.EventTransformer{
			Config:     dentConfig,
			Converter:  &dent.DentConverter{},
			Repository: &dent.DentRepository{},
		}
	})

	It("persists a flop dent log event", func() {
		blockNumber := int64(8955613)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr = initializer.NewEventTransformer(db)
		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []dentModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("10000000000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("2"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000000000000"))
		Expect(dbResult[0].ContractAddress).To(Equal(""))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.dent`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538637780 + test_data.TTL
		Expect(dbTic).To(Equal(actualTic))
	})

	It("persists a flip dent log event", func() {
		//TODO: There are currently no Flip.dent events on Kovan
	})
})

type dentModel struct {
	BidId           string `db:"bid_id"`
	Lot             string
	Bid             string
	ContractAddress string `db:"contract_address"`
}
