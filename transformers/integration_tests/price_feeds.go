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
	"github.com/vulcanize/mcd_transformers/transformers/events/price_feeds"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
)

var _ = Describe("Price feeds transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		config      transformer.EventTransformerConfig
		fetcher     *fetch.Fetcher
		initializer event.LogNoteTransformer
		topics      []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		config = transformer.EventTransformerConfig{
			TransformerName: constants.PriceFeedLabel,
			ContractAddresses: []string{
				test_data.KovanPipEthContractAddress,
				test_data.KovanPipRepContractAddress,
			},
			ContractAbi:         test_data.KovanMedianizerABI,
			Topic:               test_data.KovanLogMedianPriceSignature,
			StartingBlockNumber: 0,
			EndingBlockNumber:   -1,
		}

		topics = []common.Hash{common.HexToHash(config.Topic)}

		fetcher = fetch.NewFetcher(blockChain)

		initializer = event.LogNoteTransformer{
			Config:     config,
			Converter:  &price_feeds.PriceFeedConverter{},
			Repository: &price_feeds.PriceFeedRepository{},
		}
	})

	It("persists a ETH/USD price feed event", func() {
		blockNumber := int64(10648994)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{test_data.KovanPipEthContractAddress}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, c2.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, age, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("140705000000000000000"))
		Expect(model.Age).To(Equal("1553885404"))
		Expect(model.MedianizerAddress).To(Equal(addresses[0]))
	})

	It("rechecks price feed event", func() {
		blockNumber := int64(10648994)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{test_data.KovanPipEthContractAddress}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, c2.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, c2.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var priceFeedChecked []int
		err = db.Select(&priceFeedChecked, `SELECT price_feeds_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(priceFeedChecked[0]).To(Equal(2))

		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, age, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("140705000000000000000"))
		Expect(model.Age).To(Equal("1553885404"))
		Expect(model.MedianizerAddress).To(Equal(addresses[0]))
	})

	It("persists a REP/USD price feed event", func() {
		blockNumber := int64(10653439)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{test_data.KovanPipRepContractAddress}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, c2.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, age, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("16629421281200000000"))
		Expect(model.Age).To(Equal("1553925852"))
		Expect(model.MedianizerAddress).To(Equal(addresses[0]))
	})
})
