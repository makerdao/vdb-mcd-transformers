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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file"
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

var _ = Describe("VowFile LogNoteTransforer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		initializer shared.EventTransformer
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
		vowFileConfig := transformer.EventTransformerConfig{
			TransformerName:   constants.VowFileLabel,
			ContractAddresses: []string{test_data.VowAddress()},
			ContractAbi:       constants.VowABI(),
			Topic:             constants.VowFileSignature(),
		}

		addresses = transformer.HexStringsToAddresses(vowFileConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(vowFileConfig.Topic)}

		initializer = shared.EventTransformer{
			Config:     vowFileConfig,
			Converter:  vow_file.VowFileConverter{},
			Repository: &vow_file.VowFileRepository{},
		}
	})

	It("fetches and transforms a Vow.file event from Kovan", func() {
		blockNumber := int64(14374777)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := initializer.NewEventTransformer(db)
		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vowFileModel
		err = db.Select(&dbResult, `SELECT what, data from maker.vow_file`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("sump"))
		Expect(dbResult[0].Data).To(Equal("100000000000000000000000000000000000000000000"))
	})
})

type vowFileModel struct {
	What string
	Data string
}
