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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("SpotFilePip LogNoteTransforer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		initializer shared.LogNoteTransformer
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
		spotFilePipConfig := transformer.EventTransformerConfig{
			TransformerName:   mcdConstants.SpotFilePipLabel,
			ContractAddresses: []string{mcdConstants.SpotContractAddress()},
			ContractAbi:       mcdConstants.SpotABI(),
			Topic:             mcdConstants.SpotFilePipSignature(),
		}

		addresses = transformer.HexStringsToAddresses(spotFilePipConfig.ContractAddresses)
		topics = []common.Hash{common.HexToHash(spotFilePipConfig.Topic)}

		initializer = shared.LogNoteTransformer{
			Config:     spotFilePipConfig,
			Converter:  pip.SpotFilePipConverter{},
			Repository: &pip.SpotFilePipRepository{},
		}
	})

	It("fetches and transforms a Spot.file pip event from Kovan", func() {
		blockNumber := int64(11257235)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		tr := initializer.NewLogNoteTransformer(db)
		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []pip.SpotFilePipModel
		err = db.Select(&dbResult, `SELECT ilk_id, pip from maker.spot_file_pip`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, ilkErr := shared.GetOrCreateIlk("434f4c332d410000000000000000000000000000000000000000000000000000", db)
		Expect(ilkErr).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].Pip).To(Equal("0xaa32EB42CBf3Bdb746b659c8DAF443f710497d80"))
	})

	It("rechecks Spot.file pip event", func() {
		blockNumber := int64(11257235)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		tr := initializer.NewLogNoteTransformer(db)
		err = tr.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var spotFilePipChecked []int
		err = db.Select(&spotFilePipChecked, `SELECT spot_file_pip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(spotFilePipChecked[0]).To(Equal(2))
	})
})
