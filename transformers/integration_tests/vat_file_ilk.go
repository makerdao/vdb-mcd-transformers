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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("VatFileIlk EventTransformer", func() {
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
		vatFileIlkConfig := transformer.EventTransformerConfig{
			TransformerName:   constants.VatFileIlkLabel,
			ContractAddresses: []string{test_data.VatAddress()},
			ContractAbi:       constants.VatABI(),
			Topic:             constants.VatFileIlkSignature(),
		}

		addresses = transformer.HexStringsToAddresses(vatFileIlkConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(vatFileIlkConfig.Topic)}

		initializer = shared.EventTransformer{
			Config:     vatFileIlkConfig,
			Converter:  &ilk.VatFileIlkConverter{},
			Repository: &ilk.VatFileIlkRepository{},
		}
	})

	It("fetches and transforms a Vat.file ilk 'spot' event from Kovan", func() {
		blockNumber := int64(14374954)
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

		var dbResult []vatFileIlkModel
		err = db.Select(&dbResult, `SELECT ilk_id, what, data from maker.vat_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult[0].What).To(Equal("spot"))
		Expect(dbResult[0].Data).To(Equal("210466666666666666666666666666"))
	})

	It("fetches and transforms a Vat.file ilk 'line' event from Kovan", func() {
		blockNumber := int64(14374894)
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

		var dbResults []vatFileIlkModel
		err = db.Select(&dbResults, `SELECT ilk_id, what, data from maker.vat_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x5341490000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		dbResult := dbResults[0]
		Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult.What).To(Equal("line"))
		Expect(dbResult.Data).To(Equal("10000000000000000000000000000000000000000000000000000"))
	})

	It("fetches and transforms a Vat.file ilk 'dust' event from Kovan", func() {
		blockNumber := int64(14374919)
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

		var dbResults []vatFileIlkModel
		err = db.Select(&dbResults, `SELECT ilk_id, what, data from maker.vat_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x474e542d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		dbResult := dbResults[0]
		Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult.What).To(Equal("dust"))
		Expect(dbResult.Data).To(Equal("0"))
	})
})

type vatFileIlkModel struct {
	Ilk  string `db:"ilk_id"`
	What string
	Data string
}
