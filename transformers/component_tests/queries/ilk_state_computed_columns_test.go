package queries

import (
	"database/sql"
	"math/rand"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Ilk state computed columns", func() {
	var (
		db         *postgres.DB
		fakeBlock  int
		fakeGuy    = "fakeAddress"
		fakeHeader core.Header
		headerId   int64
		ilkID      int64
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
		ilkValues := test_helpers.GetIlkValues(0)
		createIlkAtBlock(fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
			test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

		getIlkErr := db.Get(&ilkID, `SELECT id FROM maker.ilks WHERE ilk = $1`, test_helpers.FakeIlk.Hex)
		Expect(getIlkErr).NotTo(HaveOccurred())
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
				`SELECT ilk_name, urn_id, dink, dart FROM api.ilk_state_frobs(
                        (SELECT (ilk_id, ilk_name, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated)::api.ilk_state
                         FROM api.get_ilk($1, $2))
                    )`, fakeBlock, ilkID)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			expectedFrobs := []test_helpers.FrobEvent{{
				IlkName: test_helpers.FakeIlk.Name,
				UrnId:   frobEvent.Urn,
				Dink:    frobEvent.Dink,
				Dart:    frobEvent.Dart,
			}}

			Expect(actualFrobs).To(Equal(expectedFrobs))
		})
	})

	Describe("ilks_state_files", func() {
		It("returns file event for an ilk state", func() {
			fileRepo := ilk.VatFileIlkRepository{}
			fileRepo.SetDB(db)
			fileEvent := test_data.VatFileIlkDustModel
			fileEvent.Ilk = test_helpers.FakeIlk.Hex
			insertFileErr := fileRepo.Create(headerId, []interface{}{fileEvent})
			Expect(insertFileErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.FileEvent
			getFilesErr := db.Select(&actualFiles,
				`SELECT id, ilk_name, what, data FROM api.ilk_state_files(
                        (SELECT (ilk_id, ilk_name, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated)::api.ilk_state
                         FROM api.get_ilk($1, $2))
                    )`, fakeBlock, ilkID)
			Expect(getFilesErr).NotTo(HaveOccurred())

			expectedFiles := []test_helpers.FileEvent{{
				Id: strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
				IlkName: sql.NullString{
					String: test_helpers.FakeIlk.Name,
					Valid:  true,
				},
				What: fileEvent.What,
				Data: fileEvent.Data,
			}}

			Expect(actualFiles).To(Equal(expectedFiles))
		})
	})
})
