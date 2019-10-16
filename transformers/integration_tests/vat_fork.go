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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fork"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = XDescribe("Vat fork transformer", func() {
	// TODO: Update when event exists in kovan
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

	vatForkConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatForkLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatForkSignature(),
	}

	It("fetches and transforms vat fork event", func() {
		blockNumber := int64(13234996)
		vatForkConfig.StartingBlockNumber = blockNumber
		vatForkConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.EventTransformer{
			Config:     vatForkConfig,
			Converter:  &vat_fork.VatForkConverter{},
			Repository: &vat_fork.VatForkRepository{},
		}
		tr := initializer.NewEventTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(vatForkConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatForkConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())
		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult vatForkModel
		err = db.Get(&dbResult, `SELECT ilk_id, src, dst, dink, dart FROM maker.vat_fork`)
		Expect(err).NotTo(HaveOccurred())

		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Ilk).To(Equal(ilkID))
		Expect(dbResult.Src).To(Equal("0xAE21412A422279B72aA8641a3D5F1da4BF6cfD30"))
		Expect(dbResult.Dst).To(Equal("0xdB33dFD3D61308C33C63209845DaD3e6bfb2c674"))
		Expect(dbResult.Dink).To(Equal(0))
		Expect(dbResult.Dart).To(Equal(0))
	})
})

type vatForkModel struct {
	Ilk  int64 `db:"ilk_id"`
	Src  string
	Dst  string
	Dink int
	Dart int
}
