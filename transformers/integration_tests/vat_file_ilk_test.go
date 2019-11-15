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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/ilk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatFileIlk EventTransformer", func() {
	var (
		initializer event.Transformer
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		vatFileIlkConfig := transformer.EventTransformerConfig{
			TransformerName:   constants.VatFileIlkTable,
			ContractAddresses: []string{test_data.VatAddress()},
			ContractAbi:       constants.VatABI(),
			Topic:             constants.VatFileIlkSignature(),
		}

		addresses = transformer.HexStringsToAddresses(vatFileIlkConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(vatFileIlkConfig.Topic)}

		initializer = event.Transformer{
			Config:    vatFileIlkConfig,
			Converter: ilk.Converter{},
		}
	})

	It("fetches and transforms a Vat.file ilk 'spot' event", func() {
		blockNumber := int64(8928374)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := initializer.NewTransformer(db)
		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatFileIlkModel
		err = db.Select(&dbResult, `SELECT ilk_id, what, data from maker.vat_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x5341490000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult[0].What).To(Equal("spot"))
		// TODO: verify this should actually be zero
		Expect(dbResult[0].Data).To(Equal("100000000000000000000000000000000000000000000000000"))
	})

	// TODO: add when ilk line set
	XIt("fetches and transforms a Vat.file ilk 'line' event from Kovan", func() {
		blockNumber := int64(14374894)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := initializer.NewTransformer(db)
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
		Expect(dbResult.Data).To(Equal("100000000000000000000000000000000000000000000000000000"))
	})

	It("fetches and transforms a Vat.file ilk 'dust' event", func() {
		blockNumber := int64(8928348)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := initializer.NewTransformer(db)
		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatFileIlkModel
		err = db.Select(&dbResults, `SELECT ilk_id, what, data from maker.vat_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		dbResult := dbResults[0]
		Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(dbResult.What).To(Equal("dust"))
		Expect(dbResult.Data).To(Equal("20000000000000000000000000000000000000000000000"))
	})
})

type vatFileIlkModel struct {
	Ilk  string `db:"ilk_id"`
	What string
	Data string
}
