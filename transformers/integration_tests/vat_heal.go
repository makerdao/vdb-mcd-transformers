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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_heal"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

// TODO: Replace block number once there's a heal event on the updated Vat
var _ = XDescribe("VatHeal Transformer", func() {
	vatHealConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatHealLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatHealSignature(),
	}

	It("transforms VatHeal log events", func() {
		blockNumber := int64(10921610)
		vatHealConfig.StartingBlockNumber = blockNumber
		vatHealConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatHealConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatHealConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		tr := shared.LogNoteTransformer{
			Config:     vatHealConfig,
			Converter:  &vat_heal.VatHealConverter{},
			Repository: &vat_heal.VatHealRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatHealModel
		err = db.Select(&dbResults, `SELECT urn, v, rad from maker.vat_heal`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.Rad).To(Equal("0"))
	})
})

type vatHealModel struct {
	Rad              string
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
