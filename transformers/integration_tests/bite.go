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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/bite"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bite Transformer", func() {
	biteConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.BiteLabel,
		ContractAddresses: []string{test_data.CatAddress()},
		ContractAbi:       constants.CatABI(),
		Topic:             constants.BiteSignature(),
	}

	It("fetches and transforms a Bite event from Kovan chain", func() {
		blockNumber := int64(14783840)
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
		Expect(dbResult[0].Art).To(Equal("20990706494289402873"))
		urnID, err := shared.GetOrCreateUrn("0xA17B62E20d2d700950cf95ceE5d8cC910850Dd98",
			"0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Urn).To(Equal(strconv.FormatInt(urnID, 10)))
		Expect(dbResult[0].Ink).To(Equal("170000000000000000"))
		Expect(dbResult[0].Flip).To(Equal("0x6eb58e32E3f4407F074DE5FEE7454f1749Be57B2"))
		Expect(dbResult[0].Tab).To(Equal("21001400664939036955849236599648053637998064448"))
		Expect(dbResult[0].Id).To(Equal("1"))
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
