package queries

import (
	"database/sql"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip"
	ilk2 "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Ilk File Events Query", func() {
	var (
		catFileChopLumpRepo chop_lump.CatFileChopLumpRepository
		catFileFlipRepo     flip.CatFileFlipRepository
		db                  *postgres.DB
		err                 error
		headerOneId         int64
		headerRepo          datastore.HeaderRepository
		jugFileRepo         ilk2.JugFileIlkRepository
		relevantIlkName     = sql.NullString{
			String: test_helpers.FakeIlk.Name,
			Valid:  true,
		}
		vatFileRepo ilk.VatFileIlkRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		catFileChopLumpRepo = chop_lump.CatFileChopLumpRepository{}
		catFileChopLumpRepo.SetDB(db)
		catFileFlipRepo = flip.CatFileFlipRepository{}
		catFileFlipRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		headerOne := fakes.GetFakeHeader(1)
		headerOneId, err = headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(err).NotTo(HaveOccurred())
		jugFileRepo = ilk2.JugFileIlkRepository{}
		jugFileRepo.SetDB(db)
		vatFileRepo = ilk.VatFileIlkRepository{}
		vatFileRepo.SetDB(db)
	})

	It("returns all ilk file events for ilk", func() {
		catFileChopLump := test_data.CatFileChopModel
		catFileChopLump.Ilk = test_helpers.FakeIlk.Hex
		chopLumpErr := catFileChopLumpRepo.Create(headerOneId, []interface{}{catFileChopLump})
		Expect(chopLumpErr).NotTo(HaveOccurred())

		catFileFlip := test_data.CatFileFlipModel
		catFileFlip.Ilk = test_helpers.FakeIlk.Hex
		flipErr := catFileFlipRepo.Create(headerOneId, []interface{}{catFileFlip})
		Expect(flipErr).NotTo(HaveOccurred())

		jugFile := test_data.JugFileIlkModel
		jugFile.Ilk = test_helpers.FakeIlk.Hex
		jugErr := jugFileRepo.Create(headerOneId, []interface{}{jugFile})
		Expect(jugErr).NotTo(HaveOccurred())

		vatFile := test_data.VatFileIlkDustModel
		vatFile.Ilk = test_helpers.FakeIlk.Hex
		vatErr := vatFileRepo.Create(headerOneId, []interface{}{vatFile})
		Expect(vatErr).NotTo(HaveOccurred())

		var actualFiles []test_helpers.IlkFileEvent
		filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1)`, test_helpers.FakeIlk.Name)
		Expect(filesErr).NotTo(HaveOccurred())

		Expect(actualFiles).To(ConsistOf(
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkName,
				What:          catFileChopLump.What,
				Data:          catFileChopLump.Data,
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkName,
				What:          catFileFlip.What,
				Data:          catFileFlip.Flip,
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkName,
				What:          jugFile.What,
				Data:          jugFile.Data,
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkName,
				What:          vatFile.What,
				Data:          vatFile.Data,
			},
		))
	})

	It("includes results across blocks", func() {
		fileBlockOne := test_data.VatFileIlkDustModel
		fileBlockOne.Ilk = test_helpers.FakeIlk.Hex
		fileBlockOne.Data = strconv.Itoa(rand.Int())
		fileBlockOneErr := vatFileRepo.Create(headerOneId, []interface{}{fileBlockOne})
		Expect(fileBlockOneErr).NotTo(HaveOccurred())

		headerTwo := fakes.GetFakeHeader(2)
		headerTwo.Hash = "anotherHash"
		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		fileBlockTwo := test_data.VatFileIlkDustModel
		fileBlockTwo.Ilk = test_helpers.FakeIlk.Hex
		fileBlockTwo.Data = strconv.Itoa(rand.Int())
		fileBlockTwoErr := vatFileRepo.Create(headerTwoId, []interface{}{fileBlockTwo})
		Expect(fileBlockTwoErr).NotTo(HaveOccurred())

		var actualFiles []test_helpers.IlkFileEvent
		filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1)`, test_helpers.FakeIlk.Name)
		Expect(filesErr).NotTo(HaveOccurred())

		Expect(actualFiles).To(ConsistOf(
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkName,
				What:          fileBlockOne.What,
				Data:          fileBlockOne.Data,
			},
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkName,
				What:          fileBlockTwo.What,
				Data:          fileBlockTwo.Data,
			},
		))
	})

	It("does not include ilk file events for a different ilk", func() {
		relevantFile := test_data.VatFileIlkDustModel
		relevantFile.Ilk = test_helpers.FakeIlk.Hex
		relevantFile.Data = strconv.Itoa(rand.Int())

		irrelevantFile := test_data.VatFileIlkDustModel
		irrelevantFile.Ilk = test_helpers.AnotherFakeIlk.Hex
		irrelevantFile.Data = strconv.Itoa(rand.Int())
		irrelevantFile.TransactionIndex = test_data.VatFileIlkDustModel.TransactionIndex + 1

		vatBlockOneErr := vatFileRepo.Create(headerOneId, []interface{}{relevantFile, irrelevantFile})
		Expect(vatBlockOneErr).NotTo(HaveOccurred())

		var actualFiles []test_helpers.IlkFileEvent
		filesErr := db.Select(&actualFiles, `SELECT ilk_identifier, what, data FROM api.all_ilk_file_events($1)`, test_helpers.FakeIlk.Name)
		Expect(filesErr).NotTo(HaveOccurred())

		Expect(actualFiles).To(ConsistOf(
			test_helpers.IlkFileEvent{
				IlkIdentifier: relevantIlkName,
				What:          relevantFile.What,
				Data:          relevantFile.Data,
			},
		))
	})

	It("fails if no argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.all_ilk_file_events()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.all_ilk_file_events() does not exist"))
	})
})
