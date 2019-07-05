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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Bite Transformer", func() {
	biteConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.BiteLabel,
		ContractAddresses: []string{mcdConstants.CatContractAddress()},
		ContractAbi:       mcdConstants.CatABI(),
		Topic:             mcdConstants.BiteSignature(),
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
			Converter:  &bite.BiteConverter{},
			Repository: &bite.BiteRepository{},
		}
		transformer := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(biteConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(biteConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []bite.BiteModel
		err = db.Select(&dbResult, `SELECT art, ink, flip, tab, urn_id, bite_identifier from maker.bite`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Art).To(Equal("149846666666666655744"))
		urnID, err := shared.GetOrCreateUrn("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb",
			"0x4554480000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Urn).To(Equal(strconv.Itoa(urnID)))
		Expect(dbResult[0].Ink).To(Equal("1000000000000000000"))
		Expect(dbResult[0].Flip).To(Equal("2"))
		Expect(dbResult[0].Tab).To(Equal("149846666666666655744"))
		Expect(dbResult[0].Id).To(Equal(""))
	})

	It("unpacks an event log", func() {
		address := common.HexToAddress(mcdConstants.CatContractAddress())
		abi, err := geth.ParseAbi(mcdConstants.CatABI())
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &bite.BiteEntity{}

		var eventLog = test_data.EthBiteLog

		err = contract.UnpackLog(entity, "Bite", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.BiteEntity
		Expect(entity.Art).To(Equal(expectedEntity.Art))
		Expect(entity.Ilk).To(Equal(expectedEntity.Ilk))
		Expect(entity.Ink).To(Equal(expectedEntity.Ink))
		Expect(entity.Id).To(Equal(expectedEntity.Id))
	})
})
