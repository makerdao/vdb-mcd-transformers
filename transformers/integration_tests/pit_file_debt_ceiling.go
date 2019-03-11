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
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	c2 "github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	fetch "github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/pit_file/debt_ceiling"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories"
)

var _ = Describe("PitFileDebtCeiling LogNoteTransformer", func() {
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

	It("fetches and transforms a PitFileDebtCeiling event from Kovan chain", func() {
		blockNumber := int64(8535578)
		config := transformer.TransformerConfig{
			TransformerName:     constants.PitFileDebtCeilingLabel,
			ContractAddresses:   []string{test_data.KovanPitContractAddress},
			ContractAbi:         test_data.KovanPitABI,
			Topic:               test_data.KovanPitFileDebtCeilingSignature,
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

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &debt_ceiling.PitFileDebtCeilingConverter{},
			Repository: &debt_ceiling.PitFileDebtCeilingRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, c2.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []debt_ceiling.PitFileDebtCeilingModel
		err = db.Select(&dbResult, `SELECT what, data from maker.pit_file_debt_ceiling`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("Line"))
		Expect(dbResult[0].Data).To(Equal("10000000.000000000000000000"))
	})

	It("rechecks pit file debt ceiling event", func() {
		blockNumber := int64(8535578)
		config := transformer.TransformerConfig{
			TransformerName:     constants.PitFileDebtCeilingLabel,
			ContractAddresses:   []string{test_data.KovanPitContractAddress},
			ContractAbi:         test_data.KovanPitABI,
			Topic:               test_data.KovanPitFileDebtCeilingSignature,
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

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &debt_ceiling.PitFileDebtCeilingConverter{},
			Repository: &debt_ceiling.PitFileDebtCeilingRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, c2.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, c2.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var pitFileDebtCeilingChecked []int
		err = db.Select(&pitFileDebtCeilingChecked, `SELECT pit_file_debt_ceiling_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(pitFileDebtCeilingChecked[0]).To(Equal(2))

		var dbResult []debt_ceiling.PitFileDebtCeilingModel
		err = db.Select(&dbResult, `SELECT what, data from maker.pit_file_debt_ceiling`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("Line"))
		Expect(dbResult[0].Data).To(Equal("10000000.000000000000000000"))
	})
})
