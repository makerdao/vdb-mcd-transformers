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
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Urn state computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          = "fakeAddress"
		fakeHeader       core.Header
		headerId         int64
		vatRepository    vat.VatStorageRepository
		catRepository    cat.CatStorageRepository
		jugRepository    jug.JugStorageRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("urn_state_ilk", func() {
		It("returns the ilk for an urn", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			fakeGuy := "fakeAddress"
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			ilkRate, convertRateErr := strconv.Atoi(ilkValues[vat.IlkRate])
			Expect(convertRateErr).NotTo(HaveOccurred())
			urnSetupData.Rate = ilkRate
			ilkSpot, convertSpotErr := strconv.Atoi(ilkValues[vat.IlkSpot])
			Expect(convertSpotErr).NotTo(HaveOccurred())
			urnSetupData.Spot = ilkSpot
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
					FROM api.urn_state_ilk(
					(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
					FROM api.get_urn($1, $2, $3)))`, test_helpers.FakeIlk.Identifier, fakeGuy, fakeHeader.BlockNumber)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("urn_state_frobs", func() {
		It("returns frobs for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.VatFrobModelWithPositiveDart
			frobEvent.Urn = fakeGuy
			frobEvent.Ilk = test_helpers.FakeIlk.Hex
			insertFrobErr := frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs test_helpers.FrobEvent
			getFrobsErr := db.Get(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
                        (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
                         FROM api.all_urns($1))
                    )`, fakeBlock)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			expectedFrobs := test_helpers.FrobEvent{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Dink:          frobEvent.Dink,
				Dart:          frobEvent.Dart,
			}

			Expect(actualFrobs).To(Equal(expectedFrobs))
		})
	})

	Describe("urn_state_bites", func() {
		It("returns bites for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			biteRepo := bite.BiteRepository{}
			biteRepo.SetDB(db)
			biteEvent := test_data.BiteModel
			biteEvent.Urn = fakeGuy
			biteEvent.Ilk = test_helpers.FakeIlk.Hex
			insertBiteErr := biteRepo.Create(headerId, []interface{}{biteEvent})
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites test_helpers.BiteEvent
			getBitesErr := db.Get(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
				    (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
				    FROM api.all_urns($1)))`,
				fakeBlock)
			Expect(getBitesErr).NotTo(HaveOccurred())

			expectedBites := test_helpers.BiteEvent{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Ink:           biteEvent.Ink,
				Art:           biteEvent.Art,
				Tab:           biteEvent.Tab,
			}

			Expect(actualBites).To(Equal(expectedBites))
		})
	})
})
