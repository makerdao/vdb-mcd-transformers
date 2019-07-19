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

    "github.com/vulcanize/mcd_transformers/test_config"
    "github.com/vulcanize/mcd_transformers/transformers/events/flip_tick"
    "github.com/vulcanize/mcd_transformers/transformers/shared"
    mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

// Update when auction events are in kovan
var _ = XDescribe("Flip tick LogNoteTransformer", func() {
	flipTickConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.FlipTickLabel,
		ContractAddresses: test_data.FlipperAddresses(),
		ContractAbi:       mcdConstants.FlipABI(),
		Topic:             mcdConstants.FlipTickSignature(),
	}

	It("fetches and transforms a Flip tick event from Kovan chain", func() {
		blockNumber := int64(8935601)
		flipTickConfig.StartingBlockNumber = blockNumber
		flipTickConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(transformer.HexStringsToAddresses(flipTickConfig.ContractAddresses), []common.Hash{common.HexToHash(flipTickConfig.Topic)}, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := shared.LogNoteTransformer{
			Config:     flipTickConfig,
			Converter:  &flip_tick.FlipTickConverter{},
			Repository: &flip_tick.FlipTickRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flipTickModel
		err = db.Select(&dbResult, `SELECT bid_id, contract_address FROM maker.flip_tick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].ContractAddress).To(Equal(""))
		Expect(dbResult[0].BidId).To(Equal(""))
	})
})

type flipTickModel struct {
	BidId            string `db:"bid_id"`
	ContractAddress  string `db:"contract_address"`
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
