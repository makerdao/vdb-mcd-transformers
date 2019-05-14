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

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_flux"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("VatFlux LogNoteTransformer", func() {
	vatFluxConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatFluxLabel,
		ContractAddresses: []string{mcdConstants.VatContractAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatFluxSignature(),
	}

	// TODO: Replace block number once there's a flux event on the updated Vat
	XIt("transforms VatFlux log events", func() {
		blockNumber := int64(9004474)
		vatFluxConfig.StartingBlockNumber = blockNumber
		vatFluxConfig.EndingBlockNumber = blockNumber

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
			transformer.HexStringsToAddresses(vatFluxConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFluxConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     vatFluxConfig,
			Converter:  &vat_flux.VatFluxConverter{},
			Repository: &vat_flux.VatFluxRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vat_flux.VatFluxModel
		err = db.Select(&dbResult, `SELECT ilk_id, src, dst, rad from maker.vat_flux`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("5245500000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].Src).To(Equal("000000000000000000000000c0851f73cc8dd5c0765e71980ec7e7fd1ef74434"))
		Expect(dbResult[0].Dst).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
		Expect(dbResult[0].Wad).To(Equal("1800000000000000000000000000000000000000000000"))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})

	// TODO: Replace block number once there's a flux event on the updated Vat
	XIt("rechecks vat flux event", func() {
		blockNumber := int64(9004474)
		vatFluxConfig.StartingBlockNumber = blockNumber
		vatFluxConfig.EndingBlockNumber = blockNumber

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
			transformer.HexStringsToAddresses(vatFluxConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFluxConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     vatFluxConfig,
			Converter:  &vat_flux.VatFluxConverter{},
			Repository: &vat_flux.VatFluxRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var vatFluxChecked []int
		err = db.Select(&vatFluxChecked, `SELECT vat_flux_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(vatFluxChecked[0]).To(Equal(2))

		var dbResult []vat_flux.VatFluxModel
		err = db.Select(&dbResult, `SELECT ilk_id, src, dst, rad from maker.vat_flux`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("5245500000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].Src).To(Equal("000000000000000000000000c0851f73cc8dd5c0765e71980ec7e7fd1ef74434"))
		Expect(dbResult[0].Dst).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
		Expect(dbResult[0].Wad).To(Equal("1800000000000000000000000000000000000000000000"))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})
})
