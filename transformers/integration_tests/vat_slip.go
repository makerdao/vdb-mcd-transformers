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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_slip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
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
		TransformerName:   mcdConstants.VatSlipLabel,
		ContractAddresses: []string{mcdConstants.VatContractAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatSlipSignature(),
	}

	It("persists vat slip event", func() {
		blockNumber := int64(11400744)
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

		tr := shared.LogNoteTransformer{
			Config:     vatSlipConfig,
			Converter:  &vat_slip.VatSlipConverter{},
			Repository: &vat_slip.VatSlipRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(logs, header)

		Expect(err).NotTo(HaveOccurred())
		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		var model vatSlipModel
		err = db.Get(&model, `SELECT ilk_id, usr, wad, tx_idx FROM maker.vat_slip WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(model.Usr).To(Equal("0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d"))
		Expect(model.Wad).To(Equal("-8000000000000000"))
		Expect(model.TransactionIndex).To(Equal(uint(1)))
		var headerChecked bool
		err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		Expect(headerChecked).To(BeTrue())
	})
})

type vatSlipModel struct {
	Ilk              string `db:"ilk_id"`
	Usr              string
	Wad              string
	TransactionIndex uint   `db:"tx_idx"`
	LogIndex         uint   `db:"log_idx"`
	Raw              []byte `db:"raw_log"`
}
