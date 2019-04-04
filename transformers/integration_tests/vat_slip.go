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

	c2 "github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	fetch "github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_slip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
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

	It("persists vat slip event", func() {
		blockNumber := int64(10670187)
		config := transformer.EventTransformerConfig{
			TransformerName:     constants.VatSlipLabel,
			ContractAddresses:   []string{test_data.KovanVatContractAddress},
			ContractAbi:         test_data.KovanVatABI,
			Topic:               test_data.KovanVatSlipSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		fetcher := fetch.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(
			transformer.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		tr := event.LogNoteTransformer{
			Config:     config,
			Converter:  &vat_slip.VatSlipConverter{},
			Repository: &vat_slip.VatSlipRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(logs, header, c2.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		var model vat_slip.VatSlipModel
		err = db.Get(&model, `SELECT ilk_id, usr, rad, tx_idx FROM maker.vat_slip WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		ilkID, err := shared.GetOrCreateIlk("4554482d42000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(model.Usr).To(Equal("cba1bbad5fe83cf0bc96028ae3ed8bb98b56986d00000000000000000000006d"))
		Expect(model.Rad).To(Equal("100000000000000000"))
		Expect(model.TransactionIndex).To(Equal(uint(14)))
		var headerChecked bool
		err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		Expect(headerChecked).To(BeTrue())
	})

	It("rechecks vat slip event", func() {
		blockNumber := int64(10670187)
		config := transformer.EventTransformerConfig{
			TransformerName:     constants.VatSlipLabel,
			ContractAddresses:   []string{test_data.KovanVatContractAddress},
			ContractAbi:         test_data.KovanVatABI,
			Topic:               test_data.KovanVatSlipSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		fetcher := fetch.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(
			transformer.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		tr := event.LogNoteTransformer{
			Config:     config,
			Converter:  &vat_slip.VatSlipConverter{},
			Repository: &vat_slip.VatSlipRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(logs, header, c2.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, c2.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var vatSlipChecked []int
		err = db.Select(&vatSlipChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(vatSlipChecked[0]).To(Equal(2))

		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		var model vat_slip.VatSlipModel
		err = db.Get(&model, `SELECT ilk_id, usr, rad, tx_idx FROM maker.vat_slip WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		ilkID, err := shared.GetOrCreateIlk("4554482d42000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(model.Usr).To(Equal("cba1bbad5fe83cf0bc96028ae3ed8bb98b56986d00000000000000000000006d"))
		Expect(model.Rad).To(Equal("100000000000000000"))
		Expect(model.TransactionIndex).To(Equal(uint(14)))
		var headerChecked int
		err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
		Expect(headerChecked).To(Equal(2))
	})
})
