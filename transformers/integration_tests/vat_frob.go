// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("Vat frob Transformer", func() {
	var (
		db            *postgres.DB
		blockChain    core.BlockChain
		logFetcher    fetcher.ILogFetcher
		vatFrobConfig transformer.EventTransformerConfig
		initializer   shared.LogNoteTransformer
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
			TransformerName:   mcdConstants.VatFrobLabel,
			ContractAddresses: []string{mcdConstants.VatContractAddress()},
			ContractAbi:       mcdConstants.VatABI(),
			Topic:             mcdConstants.VatFrobSignature(),
		}

		initializer = shared.LogNoteTransformer{
			Config:     vatFrobConfig,
			Converter:  &vat_frob.VatFrobConverter{},
			Repository: &vat_frob.VatFrobRepository{},
		}
	})

	It("fetches and transforms a vat frob event from Kovan chain", func() {
		blockNumber := int64(11400742)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatFrobConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFrobConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vat_frob.VatFrobModel
		err = db.Select(&dbResult, `SELECT urn_id, v, w, dink, dart from maker.vat_frob`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		urnID, err := shared.GetOrCreateUrn("0xa09408f055a8F7aFD33e9F17f82f33739bE2693c", ilkID, db)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult[0].Urn).To(Equal(strconv.Itoa(urnID)))
		Expect(dbResult[0].V).To(Equal("0xA09408f055A8F7AFD33e9F17f82f33739bE2693c"))
		Expect(dbResult[0].W).To(Equal("0xA09408f055A8F7AFD33e9F17f82f33739bE2693c"))
		Expect(dbResult[0].Dink).To(Equal("10000000000000000"))
		Expect(dbResult[0].Dart).To(Equal("500000000000000000"))
	})
})
