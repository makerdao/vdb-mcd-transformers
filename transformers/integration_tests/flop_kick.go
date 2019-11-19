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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick"
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

var _ = XDescribe("FlopKick Transformer", func() {
	var (
		db             *postgres.DB
		blockChain     core.BlockChain
		flopKickConfig transformer.EventTransformerConfig
		initializer    shared.EventTransformer
		logFetcher     fetcher.ILogFetcher
		addresses      []common.Address
		topics         []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		flopKickConfig = transformer.EventTransformerConfig{
			TransformerName:   constants.FlopKickLabel,
			ContractAddresses: []string{test_data.FlopAddress()},
			ContractAbi:       constants.FlopABI(),
			Topic:             constants.FlopKickSignature(),
		}

		initializer = shared.EventTransformer{
			Config:     flopKickConfig,
			Converter:  &flop_kick.FlopKickConverter{},
			Repository: &flop_kick.FlopKickRepository{},
		}

		logFetcher = fetcher.NewLogFetcher(blockChain)
		addresses = transformer.HexStringsToAddresses(flopKickConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(flopKickConfig.Topic)}
	})

	It("fetches and transforms a FlopKick event from Kovan chain", func() {
		blockNumber := int64(8672119)
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

		var dbResult []FlopKickModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].ContractAddress).To(Equal(""))
		Expect(dbResult[0].Gal).To(Equal("0x9B870D55BaAEa9119dBFa71A92c5E26E79C4726d"))
		// this very large number appears to be derived from the data including: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		Expect(dbResult[0].Lot).To(Equal("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	})

	It("fetches and transforms another FlopKick event from Kovan chain", func() {
		blockNumber := int64(8955611)
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

		var dbResult []FlopKickModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("10000000000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("2"))
		Expect(dbResult[0].ContractAddress).To(Equal(""))
		Expect(dbResult[0].Gal).To(Equal("0x3728e9777B2a0a611ee0F89e00E01044ce4736d1"))
		Expect(dbResult[0].Lot).To(Equal("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	})
})

type FlopKickModel struct {
	BidId           string `db:"bid_id"`
	Lot             string
	Bid             string
	Gal             string
	ContractAddress string `db:"address_id"`
	HeaderID        int64  `db:"header_id"`
	LogID           int64  `db:"log_id"`
}
