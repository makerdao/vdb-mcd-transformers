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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_slip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("Vat slip transformer", func() {
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

	vatSlipConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatSlipLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatSlipSignature(),
	}

	It("persists vat slip event", func() {
		blockNumber := int64(14220396)
		vatSlipConfig.StartingBlockNumber = blockNumber
		vatSlipConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatSlipConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatSlipConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := shared.EventTransformer{
			Config:     vatSlipConfig,
			Converter:  &vat_slip.VatSlipConverter{},
			Repository: &vat_slip.VatSlipRepository{},
		}.NewEventTransformer(db)

		err = tr.Execute(headerSyncLogs)

		Expect(err).NotTo(HaveOccurred())
		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		var model vatSlipModel
		err = db.Get(&model, `SELECT ilk_id, usr, wad FROM maker.vat_slip WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		ilkID, err := shared.GetOrCreateIlk("0x5245502d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
		Expect(model.Usr).To(Equal("0xA36570028Ee05C41052ca3c61ceCd491E8161cC1"))
		Expect(model.Wad).To(Equal("50000000000000000000"))
	})
})

type vatSlipModel struct {
	Ilk string `db:"ilk_id"`
	Usr string
	Wad string
}
