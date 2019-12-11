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

	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Managed CDP trigger-populated table", func() {
	var (
		headerRepo             repositories.HeaderRepository
		repo                   cdp_manager.CdpManagerStorageRepository
		fakeCdpi               = rand.Int()
		headerOne              core.Header
		blockOne, timestampOne int
		diffID                 int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		repo = cdp_manager.CdpManagerStorageRepository{}
		repo.SetDB(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	It("stores the state of each managed CDP, unique by cdpi", func() {
		fakeIlk := test_helpers.FakeIlk.Hex
		fakeUrn := test_data.FakeUrn

		cdpManagerStorageValues1 := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr1 := test_helpers.CreateManagedCdp(db, diffID, headerOne.Id, cdpManagerStorageValues1, test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr1).NotTo(HaveOccurred())

		fakeCdpi2 := fakeCdpi + 1
		cdpManagerStorageValues2 := test_helpers.GetCdpManagerStorageValues(2, fakeIlk, fakeUrn, fakeCdpi2)
		cdpErr2 := test_helpers.CreateManagedCdp(db, diffID, headerOne.Id, cdpManagerStorageValues2, test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi2)))
		Expect(cdpErr2).NotTo(HaveOccurred())

		expectedCdp1 := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, headerOne.Timestamp, cdpManagerStorageValues1)
		expectedCdp2 := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, headerOne.Timestamp, cdpManagerStorageValues2)

		var actualCdps []test_helpers.ManagedCdp
		queryErr := db.Select(&actualCdps,
			`SELECT usr, cdpi, urn_identifier, ilk_identifier, created FROM api.managed_cdp`)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(actualCdps).To(ConsistOf([]test_helpers.ManagedCdp{expectedCdp1, expectedCdp2}))
	})

	It("stores the latest owner of the CDP according to block number", func() {
		fakeIlk := test_helpers.FakeIlk.Hex
		fakeUrn := test_data.FakeUrn

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		newOwner := "0x16Fb96a5fa0427Af0C8F7cF1eB4870231c8154B6"
		_, ownsErr := db.Exec(cdp_manager.InsertOwnsQuery, diffID, headerTwo.Id, fakeCdpi, newOwner)
		Expect(ownsErr).NotTo(HaveOccurred())

		cdpManagerStorageValues := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr := test_helpers.CreateManagedCdp(db, diffID, headerOne.Id, cdpManagerStorageValues, test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr).NotTo(HaveOccurred())

		cdpManagerStorageValues[cdp_manager.Owns] = newOwner
		expectedCdp := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, headerOne.Timestamp, cdpManagerStorageValues)

		var actualCdps []test_helpers.ManagedCdp
		queryErr := db.Select(&actualCdps,
			`SELECT usr, cdpi, urn_identifier, ilk_identifier, created FROM api.managed_cdp WHERE cdpi = $1`,
			fakeCdpi)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(actualCdps).To(Equal([]test_helpers.ManagedCdp{expectedCdp}))
	})

	It("gets time created based on when cdpi changed", func() {
		fakeIlk := test_helpers.FakeIlk.Hex
		fakeUrn := test_data.FakeUrn

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		_, cdpiErr := db.Exec(cdp_manager.InsertCdpiQuery, diffID, headerOne.Id, fakeCdpi)
		Expect(cdpiErr).NotTo(HaveOccurred())

		cdpManagerStorageValues := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr := test_helpers.CreateManagedCdp(db, diffID, headerTwo.Id, cdpManagerStorageValues, test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr).NotTo(HaveOccurred())

		expectedCdp := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, headerOne.Timestamp, cdpManagerStorageValues)

		var actualCdp test_helpers.ManagedCdp
		queryErr := db.Get(&actualCdp,
			`SELECT usr, cdpi, urn_identifier, ilk_identifier, created FROM api.managed_cdp WHERE cdpi = $1`,
			fakeCdpi)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedCdp).To(Equal(actualCdp))
	})
})
