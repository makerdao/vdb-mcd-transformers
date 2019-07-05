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
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_init"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("JugInit LogNoteTransformer", func() {
	jugInitConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.JugInitLabel,
		ContractAddresses: []string{mcdConstants.JugContractAddress()},
		ContractAbi:       mcdConstants.JugABI(),
		Topic:             mcdConstants.JugInitSignature(),
	}

	It("transforms jug init log events", func() {
		blockNumber := int64(11257255)
		jugInitConfig.StartingBlockNumber = blockNumber
		jugInitConfig.EndingBlockNumber = blockNumber

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
			transformer.HexStringsToAddresses(jugInitConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugInitConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := shared.LogNoteTransformer{
			Config:     jugInitConfig,
			Converter:  &jug_init.JugInitConverter{},
			Repository: &jug_init.JugInitRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []JugInitModel
		err = db.Select(&dbResults, `SELECT ilk_id from maker.jug_init`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		ilkID, err := shared.GetOrCreateIlk("0x434f4c352d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
	})

	It("rechecks jug init event", func() {
		blockNumber := int64(11257255)
		jugInitConfig.StartingBlockNumber = blockNumber
		jugInitConfig.EndingBlockNumber = blockNumber

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
			transformer.HexStringsToAddresses(jugInitConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugInitConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := shared.LogNoteTransformer{
			Config:     jugInitConfig,
			Converter:  &jug_init.JugInitConverter{},
			Repository: &jug_init.JugInitRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var jugInitChecked []int
		err = db.Select(&jugInitChecked, `SELECT jug_init_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(jugInitChecked[0]).To(Equal(2))

		var dbResults []JugInitModel
		err = db.Select(&dbResults, `SELECT ilk_id from maker.jug_init`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		ilkID, err := shared.GetOrCreateIlk("0x434f4c352d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
	})
})

type JugInitModel struct {
	Ilk              string `db:"ilk_id"`
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
