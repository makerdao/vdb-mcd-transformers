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
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/mat"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("SpotFile LogNoteTransformers", func() {
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

	Describe("Spot file mat", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			initializer shared.LogNoteTransformer
			logs        []types.Log
			topics      []common.Hash
			tr          transformer.EventTransformer
		)

		BeforeEach(func() {
			blockNumber = int64(11257385)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			spotFileMatConfig := transformer.EventTransformerConfig{
				TransformerName:     mcdConstants.SpotFileMatLabel,
				ContractAddresses:   []string{mcdConstants.SpotContractAddress()},
				ContractAbi:         mcdConstants.SpotABI(),
				Topic:               mcdConstants.SpotFileMatSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(spotFileMatConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(spotFileMatConfig.Topic)}

			initializer = shared.LogNoteTransformer{
				Config:     spotFileMatConfig,
				Converter:  mat.SpotFileMatConverter{},
				Repository: &mat.SpotFileMatRepository{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			tr = initializer.NewLogNoteTransformer(db)
			executeErr := tr.Execute(logs, header)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Spot.file mat event from Kovan", func() {
			var dbResult mat.SpotFileMatModel
			getSpotErr := db.Get(&dbResult, `SELECT ilk_id, what, data FROM maker.spot_file_mat`)
			Expect(getSpotErr).NotTo(HaveOccurred())

			ilkID, ilkErr := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbResult.What).To(Equal("mat"))
			Expect(dbResult.Data).To(Equal("1500000000000000000000000000"))
		})

		It("rechecks Spot.file mat event", func() {
			recheckErr := tr.Execute(logs, header)
			Expect(recheckErr).NotTo(HaveOccurred())

			var headerID int64
			getHeaderErr := db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
			Expect(getHeaderErr).NotTo(HaveOccurred())

			var spotFileMatChecked int
			getSpotCheckedErr := db.Get(&spotFileMatChecked, `SELECT spot_file_mat_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(getSpotCheckedErr).NotTo(HaveOccurred())

			Expect(spotFileMatChecked).To(Equal(2))
		})
	})

	Describe("Spot file pip", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			initializer shared.LogNoteTransformer
			logs        []types.Log
			topics      []common.Hash
			tr          transformer.EventTransformer
		)

		BeforeEach(func() {
			blockNumber = int64(11257235)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			spotFilePipConfig := transformer.EventTransformerConfig{
				TransformerName:     mcdConstants.SpotFilePipLabel,
				ContractAddresses:   []string{mcdConstants.SpotContractAddress()},
				ContractAbi:         mcdConstants.SpotABI(),
				Topic:               mcdConstants.SpotFilePipSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(spotFilePipConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(spotFilePipConfig.Topic)}

			initializer = shared.LogNoteTransformer{
				Config:     spotFilePipConfig,
				Converter:  pip.SpotFilePipConverter{},
				Repository: &pip.SpotFilePipRepository{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			tr = initializer.NewLogNoteTransformer(db)
			executeErr := tr.Execute(logs, header)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Spot.file pip event from Kovan", func() {
			var dbResult pip.SpotFilePipModel
			getSpotErr := db.Get(&dbResult, `SELECT ilk_id, pip from maker.spot_file_pip`)
			Expect(getSpotErr).NotTo(HaveOccurred())

			ilkID, ilkErr := shared.GetOrCreateIlk("0x434f4c332d410000000000000000000000000000000000000000000000000000", db)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbResult.Pip).To(Equal("0xaa32EB42CBf3Bdb746b659c8DAF443f710497d80"))
		})

		It("rechecks Spot.file pip event", func() {
			recheckErr := tr.Execute(logs, header)
			Expect(recheckErr).NotTo(HaveOccurred())

			var headerID int64
			getHeaderErr := db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
			Expect(getHeaderErr).NotTo(HaveOccurred())

			var spotFilePipChecked int
			getSpotCheckedErr := db.Get(&spotFilePipChecked, `SELECT spot_file_pip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(getSpotCheckedErr).NotTo(HaveOccurred())

			Expect(spotFilePipChecked).To(Equal(2))
		})
	})
})
