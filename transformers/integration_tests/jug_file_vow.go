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
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/vow"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("Jug File Vow LogNoteTransformer", func() {
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

	jugFileVowConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.JugFileVowLabel,
		ContractAddresses: []string{mcdConstants.JugContractAddress()},
		ContractAbi:       mcdConstants.JugABI(),
		Topic:             mcdConstants.JugFileVowSignature(),
	}

	It("transforms JugFileVow log events", func() {
		blockNumber := int64(10771086)
		jugFileVowConfig.StartingBlockNumber = blockNumber
		jugFileVowConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     jugFileVowConfig,
			Converter:  &vow.JugFileVowConverter{},
			Repository: &vow.JugFileVowRepository{},
		}
		tr := initializer.NewLogNoteTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(jugFileVowConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugFileVowConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vow.JugFileVowModel
		err = db.Select(&dbResult, `SELECT what, data FROM maker.jug_file_vow`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("vow"))
		Expect(dbResult[0].Data).To(Equal("0xa2c0d575cb4e1f145830326420e0ccfab8bebc1d"))
	})

	It("rechecks jug file vow event", func() {
		blockNumber := int64(10771086)
		jugFileVowConfig.StartingBlockNumber = blockNumber
		jugFileVowConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     jugFileVowConfig,
			Converter:  &vow.JugFileVowConverter{},
			Repository: &vow.JugFileVowRepository{},
		}
		tr := initializer.NewLogNoteTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(jugFileVowConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugFileVowConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var jugFileVowChecked []int
		err = db.Select(&jugFileVowChecked, `SELECT jug_file_vow_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(jugFileVowChecked[0]).To(Equal(2))
	})
})
