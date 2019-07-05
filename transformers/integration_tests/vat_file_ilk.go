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
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("VatFileIlk LogNoteTransformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		initializer shared.LogNoteTransformer
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
			TransformerName:   mcdConstants.VatFileIlkLabel,
			ContractAddresses: []string{mcdConstants.VatContractAddress()},
			ContractAbi:       mcdConstants.VatABI(),
			Topic:             mcdConstants.VatFileIlkSignature(),
		}

		addresses = transformer.HexStringsToAddresses(vatFileIlkConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(vatFileIlkConfig.Topic)}

		initializer = shared.LogNoteTransformer{
			Config:     vatFileIlkConfig,
			Converter:  &ilk.VatFileIlkConverter{},
			Repository: &ilk.VatFileIlkRepository{},
		}
	})

	It("fetches and transforms a Vat.file ilk 'spot' event from Kovan", func() {
		blockNumber := int64(11257491)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		tr := initializer.NewLogNoteTransformer(db)
		err = tr.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatFileIlkModel
		err = db.Select(&dbResult, `SELECT ilk_id, what, data from maker.vat_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x434f4c352d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].What).To(Equal("spot"))
		Expect(dbResult[0].Data).To(Equal("46877063947368421052631578"))
	})

	It("fetches and transforms a Vat.file ilk 'line' event from Kovan", func() {
		blockNumber := int64(11257434)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		tr := initializer.NewLogNoteTransformer(db)
		err = tr.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatFileIlkModel
		err = db.Select(&dbResults, `SELECT ilk_id, what, data from maker.vat_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x434f4c352d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		sort.Sort(byLogIndexVatFileIlk(dbResults))
		dbResult := dbResults[0]
		Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult.What).To(Equal("line"))
		Expect(dbResult.Data).To(Equal("200000000000000000000000000000000000000000000000000"))
	})
})

type vatFileIlkModel struct {
	Ilk              string `db:"ilk_id"`
	What             string
	Data             string
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

type byLogIndexVatFileIlk []vatFileIlkModel

func (c byLogIndexVatFileIlk) Len() int           { return len(c) }
func (c byLogIndexVatFileIlk) Less(i, j int) bool { return c[i].LogIndex < c[j].LogIndex }
func (c byLogIndexVatFileIlk) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
