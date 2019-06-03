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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_suck"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

// TODO: Replace block number once there's a suck event on the updated Vat
var _ = XDescribe("VatSuck Transformer", func() {
	vatSuckConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatSuckLabel,
		ContractAddresses: []string{mcdConstants.VatContractAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatSuckSignature(),
	}

	It("transforms VatSuck log events", func() {
		blockNumber := int64(0)
		vatSuckConfig.StartingBlockNumber = blockNumber
		vatSuckConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatSuckConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatSuckConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		tr := shared.LogNoteTransformer{
			Config:     vatSuckConfig,
			Converter:  &vat_suck.VatSuckConverter{},
			Repository: &vat_suck.VatSuckRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vat_suck.VatSuckModel
		err = db.Select(&dbResults, `SELECT u, v, rad from maker.vat_suck`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.U).To(Equal(""))
		Expect(dbResult.V).To(Equal(""))
		Expect(dbResult.Rad).To(Equal(""))
	})

	It("rechecks vat suck event", func() {
		blockNumber := int64(0)
		vatSuckConfig.StartingBlockNumber = blockNumber
		vatSuckConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatSuckConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatSuckConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		tr := shared.LogNoteTransformer{
			Config:     vatSuckConfig,
			Converter:  &vat_suck.VatSuckConverter{},
			Repository: &vat_suck.VatSuckRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())
		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var vatSuckChecked []int
		err = db.Select(&vatSuckChecked, `SELECT vat_suck_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(vatSuckChecked[0]).To(Equal(2))

		var dbResults []vat_suck.VatSuckModel
		err = db.Select(&dbResults, `SELECT u, v, rad from maker.vat_suck`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.U).To(Equal(""))
		Expect(dbResult.V).To(Equal(""))
		Expect(dbResult.Rad).To(Equal(""))
	})
})
