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

package queries

import (
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Flip state computed columns", func() {
	var (
		db               *postgres.DB
		fakeHeader       core.Header
		headerRepository repositories.HeaderRepository
		contractAddress  = "contract address"
		fakeIlk          = test_helpers.FakeIlk.Hex
		fakeUrn          = test_data.FlipKickModel.Usr
		fakeBidId        int
		blockNumber      int
	)

	BeforeEach(func() {
		fakeBidId = rand.Int()
		blockNumber = rand.Int()
		timestamp := int(rand.Int31())

		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		flipKickRepo := flip_kick.FlipKickRepository{}
		flipKickRepo.SetDB(db)
		dealRepo := deal.DealRepository{}
		dealRepo.SetDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeHeader = fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))
		headerId, headerOneErr := headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, fakeHeader, flipStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		_, _, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			IlkHex:           fakeIlk,
			UrnGuy:           fakeUrn,
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerId,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flip_state_ilk", func() {
		It("returns ilk_state for a flip_state", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(fakeIlk, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result, `
				SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
				FROM api.flip_state_ilk(
					(SELECT (block_height, bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated)::api.flip_state
					 FROM api.get_flip($1, $2, $3))
			)`, fakeBidId, fakeIlk, blockNumber)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("flip_state_urn", func() {
		It("returns urn_state for a flip_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData(blockNumber, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(fakeIlk, fakeUrn)
			vatRepository := vat.VatStorageRepository{}
			vatRepository.SetDB(db)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			var actualUrn test_helpers.UrnState
			getUrnErr := db.Get(&actualUrn, `
				SELECT urn_identifier, ilk_identifier
				FROM api.flip_state_urn(
					(SELECT (block_height, bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated)::api.flip_state
					 FROM api.get_flip($1, $2, $3))
			)`, fakeBidId, fakeIlk, blockNumber)

			Expect(getUrnErr).NotTo(HaveOccurred())

			expectedUrn := test_helpers.UrnState{
				UrnIdentifier: fakeUrn,
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
			}

			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})
	})
})
