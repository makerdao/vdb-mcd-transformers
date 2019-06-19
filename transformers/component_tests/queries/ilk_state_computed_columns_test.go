package queries

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Ilk state computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          = "fakeAddress"
		fakeHeader       core.Header
		headerId         int64
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

		ilkValues := test_helpers.GetIlkValues(0)
		test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
			test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("ilk_state_frobs", func() {
		It("returns relevant frobs for an ilk_state", func() {
			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.VatFrobModelWithPositiveDart
			frobEvent.Urn = fakeGuy
			frobEvent.Ilk = test_helpers.FakeIlk.Hex
			insertFrobErr := frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			getFrobsErr := db.Select(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.ilk_state_frobs(
                        (SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
                         FROM api.get_ilk($1, $2))
                    )`, test_helpers.FakeIlk.Identifier, fakeBlock)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			expectedFrobs := []test_helpers.FrobEvent{{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: frobEvent.Urn,
				Dink:          frobEvent.Dink,
				Dart:          frobEvent.Dart,
			}}

			Expect(actualFrobs).To(Equal(expectedFrobs))
		})
	})

	Describe("ilks_state_ilk_file_events", func() {
		It("returns ilk file events for an ilk state", func() {
			fileRepo := ilk.VatFileIlkRepository{}
			fileRepo.SetDB(db)
			fileEvent := test_data.VatFileIlkDustModel
			fileEvent.Ilk = test_helpers.FakeIlk.Hex
			insertFileErr := fileRepo.Create(headerId, []interface{}{fileEvent})
			Expect(insertFileErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.IlkFileEvent
			getFilesErr := db.Select(&actualFiles,
				`SELECT ilk_identifier, what, data FROM api.ilk_state_ilk_file_events(
                        (SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
                         FROM api.get_ilk($1, $2))
                    )`, test_helpers.FakeIlk.Identifier, fakeBlock)
			Expect(getFilesErr).NotTo(HaveOccurred())

			expectedFiles := []test_helpers.IlkFileEvent{{
				IlkIdentifier: test_helpers.GetValidNullString(test_helpers.FakeIlk.Identifier),
				What:          fileEvent.What,
				Data:          fileEvent.Data,
			}}

			Expect(actualFiles).To(Equal(expectedFiles))
		})
	})

	Describe("ilk_state_bites", func() {
		It("returns bite event for an ilk state", func() {
			biteRepo := bite.BiteRepository{}
			biteRepo.SetDB(db)
			biteEvent := test_data.BiteModel
			biteEvent.Ilk = test_helpers.FakeIlk.Hex
			insertBiteErr := biteRepo.Create(headerId, []interface{}{biteEvent})
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			getBitesErr := db.Select(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.ilk_state_bites(
					(SELECT (ilk_identifier, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated)::api.ilk_state
					FROM api.get_ilk($1, $2))
				)`,
				test_helpers.FakeIlk.Identifier,
				fakeBlock)
			Expect(getBitesErr).NotTo(HaveOccurred())

			expectedBites := []test_helpers.BiteEvent{{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: biteEvent.Urn,
				Ink:           biteEvent.Ink,
				Art:           biteEvent.Art,
				Tab:           biteEvent.Tab,
			}}

			Expect(actualBites).To(Equal(expectedBites))
		})
	})
})
