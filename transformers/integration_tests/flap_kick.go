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
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("FlapKick Transformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
	})

	flapKickConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.FlapKickLabel,
		ContractAddresses: []string{test_data.FlapAddress()},
		ContractAbi:       constants.FlapABI(),
		Topic:             constants.FlapKickSignature(),
	}

	It("fetches and transforms a FlapKick event from Kovan chain", func() {
		blockNumber := int64(14308157)
		flapKickConfig.StartingBlockNumber = blockNumber
		flapKickConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		tr := shared.EventTransformer{
			Config:     flapKickConfig,
			Converter:  &flap_kick.FlapKickConverter{},
			Repository: &flap_kick.FlapKickRepository{},
		}.NewEventTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(flapKickConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(flapKickConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []FlapKickModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, lot FROM maker.flap_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("28"))
		Expect(dbResult[0].Lot).To(Equal("100000000000000000000000000000000000000000000"))
	})
})

type FlapKickModel struct {
	BidId    string `db:"bid_id"`
	Lot      string
	Bid      string
	HeaderID int64 `db:"header_id"`
	LogID    int64 `db:"log_id"`
}
