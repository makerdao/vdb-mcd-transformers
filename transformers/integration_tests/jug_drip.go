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
    "github.com/vulcanize/mcd_transformers/transformers/test_data"
    "strconv"

    "github.com/ethereum/go-ethereum/common"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
    "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
    "github.com/vulcanize/vulcanizedb/pkg/core"
    "github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

    "github.com/vulcanize/mcd_transformers/test_config"
    "github.com/vulcanize/mcd_transformers/transformers/events/jug_drip"
    "github.com/vulcanize/mcd_transformers/transformers/shared"
    mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = XDescribe("JugDrip Transformer", func() {
	var (
		db            *postgres.DB
		blockChain    core.BlockChain
		jugDripConfig transformer.EventTransformerConfig
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		jugDripConfig = transformer.EventTransformerConfig{
			ContractAddresses: []string{test_data.JugAddress()},
			ContractAbi:       mcdConstants.JugABI(),
			Topic:             mcdConstants.JugDripSignature(),
		}
	})

	It("transforms JugDrip log events", func() {
		blockNumber := int64(11144455)
		jugDripConfig.StartingBlockNumber = blockNumber
		jugDripConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     jugDripConfig,
			Converter:  &jug_drip.JugDripConverter{},
			Repository: &jug_drip.JugDripRepository{},
		}
		tr := initializer.NewLogNoteTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(jugDripConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugDripConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []jugDripModel
		err = db.Select(&dbResults, `SELECT ilk_id from maker.jug_drip`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		ilkID, err := shared.GetOrCreateIlk("434f4c312d410000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Ilk).To(Equal(strconv.Itoa(ilkID)))
	})

	It("rechecks jug drip event", func() {
		blockNumber := int64(11144455)
		jugDripConfig.StartingBlockNumber = blockNumber
		jugDripConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     jugDripConfig,
			Converter:  &jug_drip.JugDripConverter{},
			Repository: &jug_drip.JugDripRepository{},
		}
		tr := initializer.NewLogNoteTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, err := f.FetchLogs(
			transformer.HexStringsToAddresses(jugDripConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(jugDripConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		err = tr.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var jugDripChecked []int
		err = db.Select(&jugDripChecked, `SELECT jug_drip FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())
	})
})

type jugDripModel struct {
	Ilk              string `db:"ilk_id"`
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
