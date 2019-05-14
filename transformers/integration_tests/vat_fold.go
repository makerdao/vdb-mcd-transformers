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
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fold"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("VatFold Transformer", func() {
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

	vatFoldConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatFoldLabel,
		ContractAddresses: []string{mcdConstants.VatContractAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             test_data.KovanVatFoldSignature,
	}

	// TODO: Replace block number once there is a fold event on the updated Vat
	XIt("transforms VatFold log events", func() {
		blockNumber := int64(9367233)
		vatFoldConfig.StartingBlockNumber = blockNumber
		vatFoldConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatFoldConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFoldConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := shared.LogNoteTransformer{
			Config:     vatFoldConfig,
			Converter:  &vat_fold.VatFoldConverter{},
			Repository: &vat_fold.VatFoldRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vat_fold.VatFoldModel
		err = db.Select(&dbResults, `SELECT urn_id, rate from maker.vat_fold`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		ilkID, err := shared.GetOrCreateIlk("5245500000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		urnID, err := shared.GetOrCreateUrn("0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1", ilkID, db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Urn).To(Equal(strconv.Itoa(urnID)))
		Expect(dbResult.Rate).To(Equal("0.000000000000000000000000000"))
	})

	// TODO: Replace block number once there is a fold event on the updated Vat
	XIt("rechecks vat fold event", func() {
		blockNumber := int64(9367233)
		vatFoldConfig.StartingBlockNumber = blockNumber
		vatFoldConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatFoldConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFoldConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := shared.LogNoteTransformer{
			Config:     vatFoldConfig,
			Converter:  &vat_fold.VatFoldConverter{},
			Repository: &vat_fold.VatFoldRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var vatFoldChecked []int
		err = db.Select(&vatFoldChecked, `SELECT vat_fold_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(vatFoldChecked[0]).To(Equal(2))

		var dbResults []vat_fold.VatFoldModel
		err = db.Select(&dbResults, `SELECT urn_id, rate from maker.vat_fold`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		ilkID, err := shared.GetOrCreateIlk("5245500000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		urnID, err := shared.GetOrCreateUrn("0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1", ilkID, db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Urn).To(Equal(strconv.Itoa(urnID)))
		Expect(dbResult.Rate).To(Equal("0.000000000000000000000000000"))
	})
})
