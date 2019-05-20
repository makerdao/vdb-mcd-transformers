package queries

import (
	"database/sql"
	"math/rand"
	"strconv"
	"strings"

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
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/vow"
	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/base"
	ilk2 "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk"
	vow2 "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/vow"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/debt_ceiling"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Files query", func() {
	Describe("ilk files", func() {
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

		It("returns all ilk files for ilk", func() {
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

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.ilk_files($1)`, test_helpers.FakeIlk.Name)
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthCatFileChopLog.Address.Hex()),
					IlkName: relevantIlkName,
					What:    catFileChopLump.What,
					Data:    catFileChopLump.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthCatFileFlipLog.Address.Hex()),
					IlkName: relevantIlkName,
					What:    catFileFlip.What,
					Data:    catFileFlip.Flip,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthJugFileIlkLog.Address.Hex()),
					IlkName: relevantIlkName,
					What:    jugFile.What,
					Data:    jugFile.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: relevantIlkName,
					What:    vatFile.What,
					Data:    vatFile.Data,
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

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.ilk_files($1)`, test_helpers.FakeIlk.Name)
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: relevantIlkName,
					What:    fileBlockOne.What,
					Data:    fileBlockOne.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: relevantIlkName,
					What:    fileBlockTwo.What,
					Data:    fileBlockTwo.Data,
				},
			))
		})

		It("does not include files for a different ilk", func() {
			relevantFile := test_data.VatFileIlkDustModel
			relevantFile.Ilk = test_helpers.FakeIlk.Hex
			relevantFile.Data = strconv.Itoa(rand.Int())

			irrelevantFile := test_data.VatFileIlkDustModel
			irrelevantFile.Ilk = test_helpers.AnotherFakeIlk.Hex
			irrelevantFile.Data = strconv.Itoa(rand.Int())
			irrelevantFile.TransactionIndex = test_data.VatFileIlkDustModel.TransactionIndex + 1

			vatBlockOneErr := vatFileRepo.Create(headerOneId, []interface{}{relevantFile, irrelevantFile})
			Expect(vatBlockOneErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.ilk_files($1)`, test_helpers.FakeIlk.Name)
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: relevantIlkName,
					What:    relevantFile.What,
					Data:    relevantFile.Data,
				},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.ilk_files()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.ilk_files() does not exist"))
		})
	})

	Describe("address files", func() {
		var (
			db                  *postgres.DB
			catFileChopLumpRepo chop_lump.CatFileChopLumpRepository
			catFileFlipRepo     flip.CatFileFlipRepository
			catFileVowRepo      vow.CatFileVowRepository
			err                 error
			headerOneId         int64
			jugFileBaseRepo     base.JugFileBaseRepository
			jugFileIlkRepo      ilk2.JugFileIlkRepository
			jugFileVowRepo      vow2.JugFileVowRepository
			populatedIlkName    = sql.NullString{
				String: test_helpers.FakeIlk.Name,
				Valid:  true,
			}
			emptyIlkName           = sql.NullString{}
			vatFileDebtCeilingRepo debt_ceiling.VatFileDebtCeilingRepository
			vatFileIlkRepo         ilk.VatFileIlkRepository
			headerRepo             datastore.HeaderRepository
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			catFileChopLumpRepo = chop_lump.CatFileChopLumpRepository{}
			catFileChopLumpRepo.SetDB(db)
			catFileFlipRepo = flip.CatFileFlipRepository{}
			catFileFlipRepo.SetDB(db)
			catFileVowRepo = vow.CatFileVowRepository{}
			catFileVowRepo.SetDB(db)
			headerRepo = repositories.NewHeaderRepository(db)
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err = headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())
			jugFileBaseRepo = base.JugFileBaseRepository{}
			jugFileBaseRepo.SetDB(db)
			jugFileIlkRepo = ilk2.JugFileIlkRepository{}
			jugFileIlkRepo.SetDB(db)
			jugFileVowRepo = vow2.JugFileVowRepository{}
			jugFileVowRepo.SetDB(db)
			vatFileDebtCeilingRepo = debt_ceiling.VatFileDebtCeilingRepository{}
			vatFileDebtCeilingRepo.SetDB(db)
			vatFileIlkRepo = ilk.VatFileIlkRepository{}
			vatFileIlkRepo.SetDB(db)
		})

		It("returns all files for cat contract address", func() {
			catFileChopLump := test_data.CatFileChopModel
			catFileChopLump.Ilk = test_helpers.FakeIlk.Hex
			chopLumpErr := catFileChopLumpRepo.Create(headerOneId, []interface{}{catFileChopLump})
			Expect(chopLumpErr).NotTo(HaveOccurred())

			catFileFlip := test_data.CatFileFlipModel
			catFileFlip.Ilk = test_helpers.FakeIlk.Hex
			catFileFlipErr := catFileFlipRepo.Create(headerOneId, []interface{}{catFileFlip})
			Expect(catFileFlipErr).NotTo(HaveOccurred())

			catFileVowErr := catFileVowRepo.Create(headerOneId, []interface{}{test_data.CatFileVowModel})
			Expect(catFileVowErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.address_files($1)`, constants.CatContractAddress())
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthCatFileChopLog.Address.Hex()),
					IlkName: populatedIlkName,
					What:    catFileChopLump.What,
					Data:    catFileChopLump.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthCatFileFlipLog.Address.Hex()),
					IlkName: populatedIlkName,
					What:    catFileFlip.What,
					Data:    catFileFlip.Flip,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthCatFileVowLog.Address.Hex()),
					IlkName: emptyIlkName,
					What:    test_data.CatFileVowModel.What,
					Data:    test_data.CatFileVowModel.Data,
				},
			))
		})

		It("gets all files for jug contract address", func() {
			jugFileBaseErr := jugFileBaseRepo.Create(headerOneId, []interface{}{test_data.JugFileBaseModel})
			Expect(jugFileBaseErr).NotTo(HaveOccurred())

			jugFileIlk := test_data.JugFileIlkModel
			jugFileIlk.Ilk = test_helpers.FakeIlk.Hex
			jugFileIlkErr := jugFileIlkRepo.Create(headerOneId, []interface{}{jugFileIlk})
			Expect(jugFileIlkErr).NotTo(HaveOccurred())

			jugFileVowErr := jugFileVowRepo.Create(headerOneId, []interface{}{test_data.JugFileVowModel})
			Expect(jugFileVowErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.address_files($1)`, constants.JugContractAddress())
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthJugFileBaseLog.Address.Hex()),
					IlkName: emptyIlkName,
					What:    test_data.JugFileBaseModel.What,
					Data:    test_data.JugFileBaseModel.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthJugFileIlkLog.Address.Hex()),
					IlkName: populatedIlkName,
					What:    jugFileIlk.What,
					Data:    jugFileIlk.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthJugFileVowLog.Address.Hex()),
					IlkName: emptyIlkName,
					What:    test_data.JugFileVowModel.What,
					Data:    test_data.JugFileVowModel.Data,
				},
			))
		})

		It("gets all files for vat contract address", func() {
			vatFileDebtCeilingErr := vatFileDebtCeilingRepo.Create(headerOneId, []interface{}{test_data.VatFileDebtCeilingModel})
			Expect(vatFileDebtCeilingErr).NotTo(HaveOccurred())

			vatFileIlk := test_data.VatFileIlkDustModel
			vatFileIlk.Ilk = test_helpers.FakeIlk.Hex
			vatFileIlkErr := vatFileIlkRepo.Create(headerOneId, []interface{}{vatFileIlk})
			Expect(vatFileIlkErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.address_files($1)`, constants.VatContractAddress())
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileDebtCeilingLog.Address.Hex()),
					IlkName: emptyIlkName,
					What:    test_data.VatFileDebtCeilingModel.What,
					Data:    test_data.VatFileDebtCeilingModel.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: populatedIlkName,
					What:    vatFileIlk.What,
					Data:    vatFileIlk.Data,
				},
			))
		})

		It("includes results across blocks", func() {
			fileBlockOne := test_data.VatFileIlkDustModel
			fileBlockOne.Ilk = test_helpers.FakeIlk.Hex
			fileBlockOne.Data = strconv.Itoa(rand.Int())
			fileBlockOneErr := vatFileIlkRepo.Create(headerOneId, []interface{}{fileBlockOne})
			Expect(fileBlockOneErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeader(2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			fileBlockTwo := test_data.VatFileIlkDustModel
			fileBlockTwo.Ilk = test_helpers.FakeIlk.Hex
			fileBlockTwo.Data = strconv.Itoa(rand.Int())
			fileBlockTwoErr := vatFileIlkRepo.Create(headerTwoId, []interface{}{fileBlockTwo})
			Expect(fileBlockTwoErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.address_files($1)`, constants.VatContractAddress())
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: populatedIlkName,
					What:    fileBlockOne.What,
					Data:    fileBlockOne.Data,
				},
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: populatedIlkName,
					What:    fileBlockTwo.What,
					Data:    fileBlockTwo.Data,
				},
			))
		})

		It("includes results across ilks", func() {
			ilkFile := test_data.VatFileIlkDustModel
			ilkFile.Ilk = test_helpers.FakeIlk.Hex

			anotherIlkFile := test_data.VatFileIlkDustModel
			anotherIlkFile.Ilk = test_helpers.AnotherFakeIlk.Hex
			anotherIlkFile.TransactionIndex = test_data.VatFileIlkDustModel.TransactionIndex + 1

			vatBlockOneErr := vatFileIlkRepo.Create(headerOneId, []interface{}{ilkFile, anotherIlkFile})
			Expect(vatBlockOneErr).NotTo(HaveOccurred())

			var actualFiles []test_helpers.FileEvent
			filesErr := db.Select(&actualFiles, `SELECT id, ilk_name, what, data FROM api.address_files($1)`, constants.VatContractAddress())
			Expect(filesErr).NotTo(HaveOccurred())

			Expect(actualFiles).To(ConsistOf(
				test_helpers.FileEvent{
					Id:      strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: populatedIlkName,
					What:    ilkFile.What,
					Data:    ilkFile.Data,
				},
				test_helpers.FileEvent{
					Id: strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: sql.NullString{
						String: test_helpers.AnotherFakeIlk.Name,
						Valid:  true,
					},
					What: anotherIlkFile.What,
					Data: anotherIlkFile.Data,
				},
			))
		})

		It("includes results with different address case", func() {
			fileModel := test_data.VatFileDebtCeilingModel
			createErr := vatFileDebtCeilingRepo.Create(headerOneId, []interface{}{fileModel})
			Expect(createErr).NotTo(HaveOccurred())

			address := test_data.EthVatFileDebtCeilingLog.Address.Hex()
			lowerCaseAddress := strings.ToLower(address)
			upperCaseAddress := strings.ToUpper(address)

			var lowerCaseAddressFiles []test_helpers.FileEvent
			lowerAddressErr := db.Select(&lowerCaseAddressFiles, `SELECT id, ilk_name, what, data FROM api.address_files($1)`, lowerCaseAddress)
			Expect(lowerAddressErr).NotTo(HaveOccurred())

			var upperCaseAddressFiles []test_helpers.FileEvent
			upperAddressErr := db.Select(&upperCaseAddressFiles, `SELECT id, ilk_name, what, data FROM api.address_files($1)`, upperCaseAddress)
			Expect(upperAddressErr).NotTo(HaveOccurred())

			Expect(lowerCaseAddress).NotTo(BeEmpty())
			Expect(lowerCaseAddressFiles).To(ConsistOf(upperCaseAddressFiles))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.address_files()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.address_files() does not exist"))
		})
	})
})
