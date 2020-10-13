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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/bite"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bite Transformer", func() {
	Context("fetches and transforms a Bite event from Cat v1.0.0", func() {
		blockNumber := int64(8997324)
		urn := "0x0A051CD913dFD1820dbf87a9bf62B04A129F88A5"
		ilk := "0x4554482d41000000000000000000000000000000000000000000000000000000"
		expectedResult := biteModel{
			Ink:  "50000000000000000000",
			Art:  "4460522851157616216837",
			Tab:  "4466031366353941646208178591268931635087392443453",
			Flip: test_data.FlipEthV100Address(),
			Id:   "112",
		}
		biteIntegrationTest(blockNumber, test_data.Cat100Address(), constants.Cat100ABI(), ilk, urn, expectedResult)
	})

	Context("fetches and transforms a Bite event from Cat v1.1.0", func() {
		blockNumber := int64(10782907)
		urn := "0x8b7b68b93cb709976f4fefdc05408039e9927246"
		ilk := "0x4241542d41000000000000000000000000000000000000000000000000000000"
		expectedResult := biteModel{
			Ink:  "2479962706275246500000",
			Art:  "504881251104361771684",
			Tab:  "515000000000000000000921515165752588303020959580",
			Flip: test_data.FlipBatV110Address(),
			Id:   "1",
		}
		biteIntegrationTest(blockNumber, test_data.Cat110Address(), constants.Cat110ABI(), ilk, urn, expectedResult)
	})
})

func biteIntegrationTest(blockNumber int64, contractAddressHex, contractABI, ilk, urn string, expectedResult biteModel) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		biteConfig := event.TransformerConfig{
			ContractAbi:         contractABI,
			ContractAddresses:   []string{contractAddressHex},
			EndingBlockNumber:   blockNumber,
			StartingBlockNumber: blockNumber,
			Topic:               constants.BiteSignature(),
			TransformerName:     constants.BiteTable,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      biteConfig,
			Transformer: bite.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(biteConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(biteConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult biteModel
		err = db.Get(&dbResult, `SELECT art, ink, flip, tab, urn_id, bid_id, address_id from maker.bite`)
		Expect(err).NotTo(HaveOccurred())

		urnID, urnErr := shared.GetOrCreateUrn(urn, ilk, db)
		Expect(urnErr).NotTo(HaveOccurred())
		expectedResult.Urn = urnID

		addressID, addressErr := repository.GetOrCreateAddress(db, contractAddressHex)
		Expect(addressErr).NotTo(HaveOccurred())
		expectedResult.AddressID = addressID

		Expect(dbResult).To(Equal(expectedResult))
	})
}

type biteModel struct {
	Urn       int64 `db:"urn_id"`
	Ink       string
	Art       string
	Tab       string
	Flip      string
	Id        string `db:"bid_id"`
	AddressID int64  `db:"address_id"`
}
