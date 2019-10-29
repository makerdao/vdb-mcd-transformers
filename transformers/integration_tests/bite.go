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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"strconv"
)

var _ = Describe("Bite Transformer", func() {
	biteConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.BiteLabel,
		ContractAddresses: []string{test_data.CatAddress()},
		ContractAbi:       constants.CatABI(),
		Topic:             constants.BiteSignature(),
	}

	// TODO: replace block number when there is an updated Cat bite event
	XIt("fetches and transforms a Bite event from Kovan chain", func() {
		blockNumber := int64(8956422)
		biteConfig.StartingBlockNumber = blockNumber
		biteConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.Transformer{
			Config:     biteConfig,
			Converter:  &bite.Converter{},
			Repository: &bite.Repository{},
		}
		transformer := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(biteConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(biteConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []biteModel
		err = db.Select(&dbResult, `SELECT art, ink, flip, tab, urn_id, bid_id from maker.bite`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Art).To(Equal("149846666666666655744"))
		urnID, err := shared.GetOrCreateUrn("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb",
			"0x4554480000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Urn).To(Equal(strconv.FormatInt(urnID, 10)))
		Expect(dbResult[0].Ink).To(Equal("1000000000000000000"))
		Expect(dbResult[0].Flip).To(Equal("2"))
		Expect(dbResult[0].Tab).To(Equal("149846666666666655744"))
		Expect(dbResult[0].Id).To(Equal(""))
	})
})

type biteModel struct {
	Ilk      string
	Urn      string `db:"urn_id"`
	Ink      string
	Art      string
	Tab      string
	Flip     string
	Id       string `db:"bid_id"`
	HeaderID int64
	LogID    int64 `db:"log_id"`
}
