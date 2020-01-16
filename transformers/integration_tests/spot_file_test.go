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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/mat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/pip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpotFile EventTransformers", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	Describe("Spot file mat", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			initializer event.ConfiguredTransformer
			logs        []types.Log
			topics      []common.Hash
			tr          transformer.EventTransformer
		)

		BeforeEach(func() {
			blockNumber = int64(8928336)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			spotFileMatConfig := transformer.EventTransformerConfig{
				TransformerName:     constants.SpotFileMatTable,
				ContractAddresses:   []string{test_data.SpotAddress()},
				ContractAbi:         constants.SpotABI(),
				Topic:               constants.SpotFileMatSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(spotFileMatConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(spotFileMatConfig.Topic)}

			initializer = event.ConfiguredTransformer{
				Config:      spotFileMatConfig,
				Transformer: mat.Transformer{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			tr = initializer.NewTransformer(db)
			executeErr := tr.Execute(eventLogs)
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

	XDescribe("Spot file par", func() {
		// TODO: Add when there is a Spot file par event available
	})

	Describe("Spot file pip", func() {
		var (
			addresses   []common.Address
			blockNumber int64
			header      core.Header
			initializer event.ConfiguredTransformer
			logs        []types.Log
			topics      []common.Hash
			tr          transformer.EventTransformer
		)

		BeforeEach(func() {
			blockNumber = int64(8928180)
			var insertHeaderErr error
			header, insertHeaderErr = persistHeader(db, blockNumber, blockChain)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			spotFilePipConfig := transformer.EventTransformerConfig{
				TransformerName:     constants.SpotFilePipTable,
				ContractAddresses:   []string{test_data.SpotAddress()},
				ContractAbi:         constants.SpotABI(),
				Topic:               constants.SpotFilePipSignature(),
				StartingBlockNumber: blockNumber,
				EndingBlockNumber:   blockNumber,
			}

			addresses = transformer.HexStringsToAddresses(spotFilePipConfig.ContractAddresses)
			topics = []common.Hash{common.HexToHash(spotFilePipConfig.Topic)}

			initializer = event.ConfiguredTransformer{
				Config:      spotFilePipConfig,
				Transformer: pip.Transformer{},
			}

			logFetcher := fetcher.NewLogFetcher(blockChain)
			var fetcherErr error
			logs, fetcherErr = logFetcher.FetchLogs(addresses, topics, header)
			Expect(fetcherErr).NotTo(HaveOccurred())

			eventLogs := test_data.CreateLogs(header.Id, logs, db)

			tr = initializer.NewTransformer(db)
			executeErr := tr.Execute(eventLogs)
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("fetches and transforms a Spot.file pip event from Kovan", func() {
			var dbResult spotFilePipModel
			getSpotErr := db.Get(&dbResult, `SELECT ilk_id, pip from maker.spot_file_pip`)
			Expect(getSpotErr).NotTo(HaveOccurred())

			ilkID, ilkErr := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
			Expect(dbResult.Pip).To(Equal("0x81FE72B5A8d1A857d176C3E7d5Bd2679A9B85763"))
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
