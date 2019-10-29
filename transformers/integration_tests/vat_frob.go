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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("Vat frob Transformer", func() {
	var (
		db            *postgres.DB
		blockChain    core.BlockChain
		logFetcher    fetcher.ILogFetcher
		vatFrobConfig transformer.EventTransformerConfig
		initializer   shared.EventTransformer
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		logFetcher = fetcher.NewLogFetcher(blockChain)
		vatFrobConfig = transformer.EventTransformerConfig{
			TransformerName:   constants.VatFrobLabel,
			ContractAddresses: []string{test_data.VatAddress()},
			ContractAbi:       constants.VatABI(),
			Topic:             constants.VatFrobSignature(),
		}

		initializer = shared.EventTransformer{
			Config:     vatFrobConfig,
			Converter:  &vat_frob.VatFrobConverter{},
			Repository: &vat_frob.VatFrobRepository{},
		}
	})

	It("fetches and transforms a vat frob event from Kovan chain", func() {
		blockNumber := int64(14438235)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatFrobConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFrobConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := initializer.NewEventTransformer(db)
		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatFrobModel
		err = db.Select(&dbResult, `SELECT urn_id, v, w, dink, dart from maker.vat_frob`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		urnID, err := shared.GetOrCreateUrn("0xE89a5a093FE1281430Fb766d809128f9c1f4bB59",
			"0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult[0].Urn).To(Equal(strconv.FormatInt(urnID, 10)))
		Expect(dbResult[0].V).To(Equal("0xE89a5a093FE1281430Fb766d809128f9c1f4bB59"))
		Expect(dbResult[0].W).To(Equal("0xE89a5a093FE1281430Fb766d809128f9c1f4bB59"))
		Expect(dbResult[0].Dink).To(Equal("0"))
		Expect(dbResult[0].Dart).To(Equal("-9996066618124872"))
	})
})

type vatFrobModel struct {
	Ilk  string
	Urn  string `db:"urn_id"`
	V    string
	W    string
	Dink string
	Dart string
}
