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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlipKick Transformer", func() {
	flipKickConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.FlipKickLabel,
		ContractAddresses: []string{test_data.EthFlipAddress()},
		ContractAbi:       constants.FlipABI(),
		Topic:             constants.FlipKickSignature(),
	}

	// TODO: Update when updated kick event exists in kovan
	XIt("fetches and transforms a FlipKick event from Kovan chain", func() {
		blockNumber := int64(8956476)
		flipKickConfig.StartingBlockNumber = blockNumber
		flipKickConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		tr := shared.EventTransformer{
			Config:     flipKickConfig,
			Converter:  &flip_kick.FlipKickConverter{},
			Repository: &flip_kick.FlipKickRepository{},
		}.NewEventTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(flipKickConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(flipKickConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []FlipKickModel
		err = db.Select(&dbResult, `SELECT bid_id, lot, bid, tab, usr, gal, contract_address FROM maker.flip_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))
		Expect(dbResult[0].Tab).To(Equal(""))
		Expect(dbResult[0].Usr).To(Equal(""))
		Expect(dbResult[0].Gal).To(Equal("0x3728e9777B2a0a611ee0F89e00E01044ce4736d1"))
		Expect(dbResult[0].ContractAddress).To(Equal(""))
	})
})

type FlipKickModel struct {
	BidId           string `db:"bid_id"`
	Lot             string
	Bid             string
	Tab             string
	Usr             string
	Gal             string
	ContractAddress string `db:"address_id"`
	HeaderID        int64  `db:"header_id"`
	LogID           int64  `db:"log_id"`
}
