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
	"github.com/vulcanize/mcd_transformers/transformers/events/pip_log_value"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("Pip LogValue transformer", func() {
	var (
		db                *postgres.DB
		blockChain        core.BlockChain
		pipLogValueConfig transformer.EventTransformerConfig
		initializer       shared.LogNoteTransformer
		logFetcher        fetcher.ILogFetcher
		topics            []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		pipLogValueConfig = transformer.EventTransformerConfig{
			TransformerName: mcdConstants.PipLogValueLabel,
			ContractAbi:     mcdConstants.PipABI(),
			Topic:           mcdConstants.PipLogValueSignature(),
		}

		topics = []common.Hash{common.HexToHash(pipLogValueConfig.Topic)}

		logFetcher = fetcher.NewLogFetcher(blockChain)

		initializer = shared.LogNoteTransformer{
			Config:     pipLogValueConfig,
			Converter:  &pip_log_value.PipLogValueConverter{},
			Repository: &pip_log_value.PipLogValueRepository{},
		}
	})

	It("persists a pip eth log value event", func() {
		blockNumber := int64(10606964)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{mcdConstants.PipEthContractAddress()}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model pip_log_value.PipLogValueModel
		err = db.Get(&model, `SELECT block_number, contract_address, val, tx_idx, raw_log FROM maker.pip_log_value WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Value).To(Equal("137120000000000000000"))
		Expect(model.ContractAddress).To(Equal(addresses[0]))
	})

	It("rechecks pip eth log value event", func() {
		blockNumber := int64(10606964)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{mcdConstants.PipEthContractAddress()}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var priceFeedChecked []int
		err = db.Select(&priceFeedChecked, `SELECT pip_log_value_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(priceFeedChecked[0]).To(Equal(2))

		var model pip_log_value.PipLogValueModel
		err = db.Get(&model, `SELECT block_number, contract_address, val, tx_idx, raw_log FROM maker.pip_log_value WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Value).To(Equal("137120000000000000000"))
		Expect(model.ContractAddress).To(Equal(addresses[0]))
	})

	It("persists a pip col1 log value event", func() {
		blockNumber := int64(10606076)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{mcdConstants.PipCol1ContractAddress()}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model pip_log_value.PipLogValueModel
		err = db.Get(&model, `SELECT block_number, contract_address, val, tx_idx, raw_log FROM maker.pip_log_value WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Value).To(Equal("107360897300000000"))
		Expect(model.ContractAddress).To(Equal(addresses[0]))
	})

	It("persists a pip col2 log value event", func() {
		blockNumber := int64(10606078)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{mcdConstants.PipCol2ContractAddress()}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model pip_log_value.PipLogValueModel
		err = db.Get(&model, `SELECT block_number, contract_address, val, tx_idx, raw_log FROM maker.pip_log_value WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Value).To(Equal("165696215265000000000"))
		Expect(model.ContractAddress).To(Equal(addresses[0]))
	})

	It("persists a pip col3 log value event", func() {
		blockNumber := int64(10606080)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{mcdConstants.PipCol3ContractAddress()}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model pip_log_value.PipLogValueModel
		err = db.Get(&model, `SELECT block_number, contract_address, val, tx_idx, raw_log FROM maker.pip_log_value WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Value).To(Equal("22770561000000000"))
		Expect(model.ContractAddress).To(Equal(addresses[0]))
	})

	It("persists a pip col4 log value event", func() {
		blockNumber := int64(10606068)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{mcdConstants.PipCol4ContractAddress()}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model pip_log_value.PipLogValueModel
		err = db.Get(&model, `SELECT block_number, contract_address, val, tx_idx, raw_log FROM maker.pip_log_value WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Value).To(Equal("60872779868500000000"))
		Expect(model.ContractAddress).To(Equal(addresses[0]))
	})

	It("persists a pip col5 log value event", func() {
		blockNumber := int64(10606070)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{mcdConstants.PipCol5ContractAddress()}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model pip_log_value.PipLogValueModel
		err = db.Get(&model, `SELECT block_number, contract_address, val, tx_idx, raw_log FROM maker.pip_log_value WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.Value).To(Equal("58804714100000000"))
		Expect(model.ContractAddress).To(Equal(addresses[0]))
	})
})
