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
    "github.com/vulcanize/mcd_transformers/transformers/events/tend"
    "github.com/vulcanize/mcd_transformers/transformers/shared"
    mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = XDescribe("Tend LogNoteTransformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		tendConfig  transformer.EventTransformerConfig
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

		tendConfig = transformer.EventTransformerConfig{
			TransformerName:   mcdConstants.TendLabel,
			ContractAddresses: append(test_data.FlipperAddresses(), test_data.FlapperAddress()),
			ContractAbi:       mcdConstants.FlipABI(),
			Topic:             mcdConstants.TendSignature(),
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = transformer.HexStringsToAddresses(tendConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(tendConfig.Topic)}

		initializer = shared.LogNoteTransformer{
			Config:     tendConfig,
			Converter:  &tend.TendConverter{},
			Repository: &tend.TendRepository{},
		}
	})

	It("fetches and transforms a Flip Tend event from Kovan chain", func() {
		blockNumber := int64(8935601)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tendModel
		err = db.Select(&dbResult, `SELECT bid_id, lot, bid, contract_address FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].ContractAddress).To(Equal(""))
		Expect(dbResult[0].Bid).To(Equal("4000"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538490276 + mcdConstants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})

	It("fetches and transforms a subsequent Flip Tend event from Kovan chain for the same auction", func() {
		blockNumber := int64(8935731)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tendModel
		err = db.Select(&dbResult, `SELECT bid_id, lot, bid, contract_address FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].ContractAddress).To(Equal(""))
		Expect(dbResult[0].Bid).To(Equal("4300"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538491224 + mcdConstants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})

	It("fetches and transforms a Flap Tend event from the Kovan chain", func() {
		blockNumber := int64(9003177)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tendModel
		err = db.Select(&dbResult, `SELECT bid_id, lot, bid, contract_address FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].ContractAddress).To(Equal(""))
		Expect(dbResult[0].Bid).To(Equal("1000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538992860 + mcdConstants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})
})

type tendModel struct {
	BidId            string `db:"bid_id"`
	Lot              string
	Bid              string
	ContractAddress  string `db:"contract_address"`
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
