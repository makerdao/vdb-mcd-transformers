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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_move"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

var _ = Describe("VatMove EventTransformer", func() {
	vatMoveConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatMoveLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatMoveSignature(),
	}

	It("transforms VatMove log events", func() {
		blockNumber := int64(14440769)
		vatMoveConfig.StartingBlockNumber = blockNumber
		vatMoveConfig.EndingBlockNumber = blockNumber

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
			transformer.HexStringsToAddresses(vatMoveConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatMoveConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := shared.EventTransformer{
			Config:     vatMoveConfig,
			Converter:  &vat_move.VatMoveConverter{},
			Repository: &vat_move.VatMoveRepository{},
		}.NewEventTransformer(db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatMoveModel
		err = db.Select(&dbResults, `SELECT src, dst, rad from maker.vat_move`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.Src).To(Equal("0x9E0d5a6a836a6C323Cf45Eb07Cb40CFc81664eec"))
		Expect(dbResult.Dst).To(Equal("0xa16b0a95a0894C685A995F0F56beDad982CD55A2"))
		Expect(dbResult.Rad).To(Equal("40000000000000000000000000000000000000000000000"))
	})
})

type vatMoveModel struct {
	Src string
	Dst string
	Rad string
}
