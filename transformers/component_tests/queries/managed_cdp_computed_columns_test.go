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
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Managed CDP computed columns", func() {
	var (
		db               *postgres.DB
		fakeHeader       core.Header
		headerRepository repositories.HeaderRepository
		storageValues    map[string]interface{}
		fakeCdpi         int
		blockNumber      int
	)

	BeforeEach(func() {
		blockNumber = rand.Int()
		fakeCdpi = rand.Int()

		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeHeader = fakes.GetFakeHeader(int64(blockNumber))
		_, headerOneErr := headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		storageValues = test_helpers.GetCdpManagerStorageValues(1, test_helpers.FakeIlk.Hex, test_data.FakeUrn, fakeCdpi)
		cdpErr := test_helpers.CreateManagedCdp(db, fakeHeader, storageValues, test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("managed_cdp_ilk", func() {
		It("returns ilk_state for a managed_cdp", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp,
				fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result, `
				SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
				FROM api.managed_cdp_ilk(
					(SELECT (id, cdpi, usr, urn_identifier, ilk_identifier, created)::api.managed_cdp
					 FROM api.managed_cdp
					 WHERE cdpi = $1))
			`, fakeCdpi)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("managed_cdp_urn", func() {
		It("returns urn_state for a managed_cdp", func() {
			urnSetupData := test_helpers.GetUrnSetupData(blockNumber, 1)
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, test_data.FakeUrn)
			vatRepository := vat.VatStorageRepository{}
			vatRepository.SetDB(db)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)
			expectedUrn := test_helpers.UrnState{
				UrnIdentifier: test_data.FakeUrn,
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
			}

			var actualUrn test_helpers.UrnState
			getUrnErr := db.Get(&actualUrn, `
				SELECT urn_identifier, ilk_identifier
				FROM api.managed_cdp_urn(
					(SELECT (id, cdpi, usr, urn_identifier, ilk_identifier, created)::api.managed_cdp
					 FROM api.managed_cdp
					 WHERE cdpi = $1))
			`, fakeCdpi)

			Expect(getUrnErr).NotTo(HaveOccurred())
			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})
	})
})
