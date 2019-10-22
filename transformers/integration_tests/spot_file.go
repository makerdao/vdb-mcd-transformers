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
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/mat"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("SpotFile EventTransformers", func() {
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
			initializer shared.EventTransformer
			logs        []types.Log
			topics      []common.Hash
			tr          transformer.EventTransformer
		)

		BeforeEach(func() {
			blockNumber = int64(13773237)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			spotFileMatConfig := transformer.EventTransformerConfig{
				TransformerName:     constants.SpotFileMatLabel,
				ContractAddresses:   []string{test_data.SpotAddress()},
				ContractAbi:         constants.SpotABI(),
				Topic:               constants.SpotFileMatSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(spotFileMatConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(spotFileMatConfig.Topic)}

			initializer = shared.EventTransformer{
				Config:     spotFileMatConfig,
				Converter:  mat.SpotFileMatConverter{},
				Repository: &mat.SpotFileMatRepository{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

			tr = initializer.NewEventTransformer(db)
			executeErr := tr.Execute(headerSyncLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Spot.file mat event from Kovan", func() {
			var dbResult spotFileMatModel
			getSpotErr := db.Get(&dbResult, `SELECT ilk_id, what, data FROM maker.spot_file_mat`)
			Expect(getSpotErr).NotTo(HaveOccurred())

			ilkID, ilkErr := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
			Expect(dbResult.What).To(Equal("mat"))
			Expect(dbResult.Data).To(Equal("1500000000000000000000000000"))
		})
	})

	Describe("Spot file pip", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			initializer shared.EventTransformer
			logs        []types.Log
			topics      []common.Hash
			tr          transformer.EventTransformer
		)

		BeforeEach(func() {
			blockNumber = int64(13772999)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			spotFilePipConfig := transformer.EventTransformerConfig{
				TransformerName:     constants.SpotFilePipLabel,
				ContractAddresses:   []string{test_data.SpotAddress()},
				ContractAbi:         constants.SpotABI(),
				Topic:               constants.SpotFilePipSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(spotFilePipConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(spotFilePipConfig.Topic)}

			initializer = shared.EventTransformer{
				Config:     spotFilePipConfig,
				Converter:  pip.SpotFilePipConverter{},
				Repository: &pip.SpotFilePipRepository{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

			tr = initializer.NewEventTransformer(db)
			executeErr := tr.Execute(headerSyncLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Spot.file pip event from Kovan", func() {
			var dbResult spotFilePipModel
			getSpotErr := db.Get(&dbResult, `SELECT ilk_id, pip from maker.spot_file_pip`)
			Expect(getSpotErr).NotTo(HaveOccurred())

			ilkID, ilkErr := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
			Expect(dbResult.Pip).To(Equal("0x75dD74e8afE8110C8320eD397CcCff3B8134d981"))
		})
	})
})

type spotFileMatModel struct {
	Ilk  string `db:"ilk_id"`
	What string
	Data string
}

type spotFilePipModel struct {
	Ilk string `db:"ilk_id"`
	Pip string
}
