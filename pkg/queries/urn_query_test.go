package queries_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/pit"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
	"time"
)

var _ = Describe("Urn view", func() {
	var (
		db         *postgres.DB
		vatRepo    vat.VatStorageRepository
		pitRepo    pit.PitStorageRepository
		headerRepo repositories.HeaderRepository
		urnOne     string
		urnTwo     string
		ilkOne     string
		ilkTwo     string
		err        error
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)
		pitRepo = pit.PitStorageRepository{}
		pitRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		rand.Seed(time.Now().UnixNano())

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)
		ilkOne = test_data.RandomString(5)
		ilkTwo = test_data.RandomString(5)
	})

	It("gets an urn", func() {
		fakeBlockNo := rand.Int()
		fakeTimestamp := 12345

		setupData := getSetupData(fakeBlockNo, fakeTimestamp)
		metadata := getMetadata(ilkOne, urnOne)
		createUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)

		var result UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM get_urn_states_at_block($1)`, fakeBlockNo)
		Expect(err).NotTo(HaveOccurred())

		expectedRatio := getExpectedRatio(setupData.Ink, setupData.Spot, setupData.Art, setupData.Rate)

		expectedUrn := UrnState{
			UrnId:   urnOne,
			IlkId:   ilkOne,
			Ink:     strconv.Itoa(setupData.Ink),
			Art:     strconv.Itoa(setupData.Art),
			Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
			Safe:    expectedRatio >= 1,
			Created: sql.NullString{String: "12345", Valid: true},
			Updated: sql.NullString{String: "12345", Valid: true},
		}

		actualRatio, err := strconv.ParseFloat(result.Ratio.String, 64)
		Expect(err).NotTo(HaveOccurred())
		result.Ratio.String = ""

		Expect(result).To(Equal(expectedUrn))
		Expect(actualRatio).Should(BeNumerically("~", expectedRatio))
	})

	It("returns the correct data for multiple urns", func() {
		blockOne := rand.Int()
		timestampOne := rand.Int()
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1

		urnOneMetadata := getMetadata(ilkOne, urnOne)
		urnOneSetupData := getSetupData(blockOne, timestampOne)
		createUrn(urnOneSetupData, urnOneMetadata, vatRepo, pitRepo, headerRepo)

		urnTwoMetadata := getMetadata(ilkTwo, urnTwo)
		urnTwoSetupData := getSetupData(blockTwo, timestampTwo)
		createUrn(urnTwoSetupData, urnTwoMetadata, vatRepo, pitRepo, headerRepo)

		var result []UrnState
		err = db.Select(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM get_urn_states_at_block($1)`, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		expectedRatioOne := getExpectedRatio(urnOneSetupData.Ink, urnOneSetupData.Spot, urnOneSetupData.Art, urnOneSetupData.Rate)
		expectedRatioTwo := getExpectedRatio(urnTwoSetupData.Ink, urnTwoSetupData.Spot, urnTwoSetupData.Art, urnTwoSetupData.Rate)

		expectedUrns := []UrnState{{
			UrnId:   urnOne,
			IlkId:   ilkOne,
			Ink:     strconv.Itoa(urnOneSetupData.Ink),
			Art:     strconv.Itoa(urnOneSetupData.Art),
			Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
			Safe:    expectedRatioOne >= 1,
			Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
			Updated: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
		}, {
			UrnId:   urnTwo,
			IlkId:   ilkTwo,
			Ink:     strconv.Itoa(urnTwoSetupData.Ink),
			Art:     strconv.Itoa(urnTwoSetupData.Art),
			Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
			Safe:    expectedRatioTwo >= 1,
			Created: sql.NullString{String: strconv.Itoa(timestampTwo), Valid: true},
			Updated: sql.NullString{String: strconv.Itoa(timestampTwo), Valid: true},
		}}

		actualRatioOne, err := strconv.ParseFloat(result[0].Ratio.String, 64)
		Expect(err).NotTo(HaveOccurred())
		result[0].Ratio.String = ""

		actualRatioTwo, err := strconv.ParseFloat(result[1].Ratio.String, 64)
		Expect(err).NotTo(HaveOccurred())
		result[1].Ratio.String = ""

		Expect(result).To(Equal(expectedUrns))
		Expect(actualRatioOne).Should(BeNumerically("~", expectedRatioOne))
		Expect(actualRatioTwo).Should(BeNumerically("~", expectedRatioTwo))
	})

	It("returns urn state without timestamps if corresponding headers aren't synced", func() {
		block := rand.Int()
		timestamp := rand.Int()
		metadata := getMetadata(ilkOne, urnOne)
		setupData := getSetupData(block, timestamp)

		createUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)
		_, err = db.Exec(`DELETE FROM headers`)
		Expect(err).NotTo(HaveOccurred())

		var result UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM get_urn_states_at_block($1)`, block)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Created.String).To(BeEmpty())
		Expect(result.Updated.String).To(BeEmpty())
	})

	Describe("it includes diffs only up to given block height", func() {
		var (
			result       UrnState
			blockOne     int
			timestampOne int
			setupDataOne SetupData
			metadata     Metadata
		)

		BeforeEach(func() {
			blockOne = rand.Int()
			timestampOne = rand.Int()
			setupDataOne = getSetupData(blockOne, timestampOne)
			metadata = getMetadata(ilkOne, urnOne)
			createUrn(setupDataOne, metadata, vatRepo, pitRepo, headerRepo)
		})

		It("gets urn state as of block one", func() {
			err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
				FROM get_urn_states_at_block($1)`, blockOne)
			Expect(err).NotTo(HaveOccurred())

			expectedRatio := getExpectedRatio(setupDataOne.Ink, setupDataOne.Spot, setupDataOne.Art, setupDataOne.Rate)

			expectedUrn := UrnState{
				UrnId:   urnOne,
				IlkId:   ilkOne,
				Ink:     strconv.Itoa(setupDataOne.Ink),
				Art:     strconv.Itoa(setupDataOne.Art),
				Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
				Safe:    expectedRatio >= 1,
				Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
				Updated: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
			}

			actualRatio, err := strconv.ParseFloat(result.Ratio.String, 64)
			Expect(err).NotTo(HaveOccurred())
			result.Ratio.String = ""

			Expect(result).To(Equal(expectedUrn))
			Expect(actualRatio).Should(BeNumerically("~", expectedRatio))
		})

		It("gets urn state with updated values", func() {
			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 1
			hashTwo := test_data.RandomString(5)
			updatedInk := rand.Int()

			err = vatRepo.Create(blockTwo, hashTwo, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			fakeHeaderTwo := fakes.GetFakeHeader(int64(blockTwo))
			fakeHeaderTwo.Timestamp = strconv.Itoa(timestampTwo)
			fakeHeaderTwo.Hash = hashTwo

			_, err = headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
				FROM get_urn_states_at_block($1)`, blockTwo)
			Expect(err).NotTo(HaveOccurred())

			expectedRatio := getExpectedRatio(updatedInk, setupDataOne.Spot, setupDataOne.Art, setupDataOne.Rate)

			expectedUrn := UrnState{
				UrnId:   urnOne,
				IlkId:   ilkOne,
				Ink:     strconv.Itoa(updatedInk),
				Art:     strconv.Itoa(setupDataOne.Art),          // Not changed
				Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
				Safe:    expectedRatio >= 1,
				Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
				Updated: sql.NullString{String: strconv.Itoa(timestampTwo), Valid: true},
			}

			actualRatio, err := strconv.ParseFloat(result.Ratio.String, 64)
			Expect(err).NotTo(HaveOccurred())
			result.Ratio.String = ""

			Expect(result).To(Equal(expectedUrn))
			Expect(actualRatio).Should(BeNumerically("~", expectedRatio))
		})
	})

	It("returns null ratio and urn being safe if there is no debt", func() {
		block := rand.Int()
		setupData := getSetupData(block, 1)
		setupData.Art = 0
		metadata := getMetadata(ilkOne, urnOne)
		createUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)

		fakeHeader := fakes.GetFakeHeader(int64(block))
		_, err = headerRepo.CreateOrUpdateHeader(fakeHeader)
		Expect(err).NotTo(HaveOccurred())

		var result UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM get_urn_states_at_block($1)`, block)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Ratio.String).To(BeEmpty())
		Expect(result.Safe).To(BeTrue())
	})
})

func getExpectedRatio(ink, spot, art, rate int) float64 {
	inkXspot := float64(ink) * float64(spot)
	artXrate := float64(art) * float64(rate)
	return inkXspot / artXrate
}

// Creates urn by creating necessary state diffs and the corresponding header
func createUrn(setupData SetupData, metadata Metadata, vatRepo vat.VatStorageRepository,
	pitRepo pit.PitStorageRepository, headerRepo repositories.HeaderRepository) {

	blockNo := int(setupData.Header.BlockNumber)
	hash := setupData.Header.Hash

	// This also creates the ilk if it doesn't exist
	err := vatRepo.Create(blockNo, hash, metadata.UrnInk, strconv.Itoa(setupData.Ink))
	Expect(err).NotTo(HaveOccurred())

	err = vatRepo.Create(blockNo, hash, metadata.UrnArt, strconv.Itoa(setupData.Art))
	Expect(err).NotTo(HaveOccurred())

	err = pitRepo.Create(blockNo, hash, metadata.IlkSpot, strconv.Itoa(setupData.Spot))
	Expect(err).NotTo(HaveOccurred())

	err = vatRepo.Create(blockNo, hash, metadata.IlkRate, strconv.Itoa(setupData.Rate))
	Expect(err).NotTo(HaveOccurred())

	_, err = headerRepo.CreateOrUpdateHeader(setupData.Header)
	Expect(err).NotTo(HaveOccurred())
}

// Does not return values computed by the query (ratio, safe, updated, created)
func getSetupData(block, timestamp int) SetupData {
	fakeHeader := fakes.GetFakeHeader(int64(block))
	fakeHeader.Timestamp = strconv.Itoa(timestamp)
	fakeHeader.Hash = test_data.RandomString(5)

	return SetupData{
		Header: fakeHeader,
		Ink:    rand.Int(),
		Art:    rand.Int(),
		Spot:   rand.Int(),
		Rate:   rand.Int(),
	}
}

type SetupData struct {
	Header core.Header
	Ink    int
	Art    int
	Spot   int
	Rate   int
}

func getMetadata(ilk, urn string) Metadata {
	return Metadata{
		UrnInk: utils.GetStorageValueMetadata(vat.UrnInk,
			map[utils.Key]string{constants.Ilk: ilk, constants.Guy: urn}, utils.Uint256),
		UrnArt: utils.GetStorageValueMetadata(vat.UrnArt,
			map[utils.Key]string{constants.Ilk: ilk, constants.Guy: urn}, utils.Uint256),
		IlkSpot: utils.GetStorageValueMetadata(pit.IlkSpot,
			map[utils.Key]string{constants.Ilk: ilk}, utils.Uint256),
		IlkRate: utils.GetStorageValueMetadata(vat.IlkRate,
			map[utils.Key]string{constants.Ilk: ilk}, utils.Uint256),
	}
}

type Metadata struct {
	UrnInk  utils.StorageValueMetadata
	UrnArt  utils.StorageValueMetadata
	IlkSpot utils.StorageValueMetadata
	IlkRate utils.StorageValueMetadata
}

type UrnState struct {
	UrnId   string
	IlkId   string
	Ink     string
	Art     string
	Ratio   sql.NullString
	Safe    bool
	Created sql.NullString
	Updated sql.NullString
	// Frobs and bites collections, and ilk object, are missing
}
