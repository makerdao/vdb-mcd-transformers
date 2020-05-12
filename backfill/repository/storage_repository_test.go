package repository_test

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Urns repository", func() {
	var (
		db   = test_config.NewTestDB(test_config.NewTestNode())
		repo repository.StorageRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = repository.NewStorageRepository(db)
	})

	Describe("GetUrnByID", func() {
		It("returns urn for matching id", func() {
			fakeIlk := test_data.RandomString(64)
			fakeIlkIdentifier := "ETH-A"
			var ilkID int
			ilkErr := db.Get(&ilkID, `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`, fakeIlk, fakeIlkIdentifier)
			Expect(ilkErr).NotTo(HaveOccurred())

			fakeUrn := test_data.RandomString(40)
			var urnID int
			urnErr := db.Get(&urnID, `INSERT INTO maker.urns (ilk_id, identifier) VALUES ($1, $2) RETURNING id`, ilkID, fakeUrn)
			Expect(urnErr).NotTo(HaveOccurred())

			urns, err := repo.GetUrnByID(urnID)

			Expect(err).NotTo(HaveOccurred())
			Expect(urns).To(Equal(repository.Urn{
				UrnID: urnID,
				Ilk:   fakeIlk,
				IlkID: ilkID,
				Urn:   fakeUrn,
			}))
		})
	})

	Describe("InsertDiff", func() {
		It("inserts a back filled diff", func() {
			diff := types.RawDiff{
				HashedAddress: common.HexToHash(test_data.RandomString(64)),
				BlockHash:     common.HexToHash(test_data.RandomString(64)),
				BlockHeight:   rand.Int(),
				StorageKey:    common.HexToHash(test_data.RandomString(64)),
				StorageValue:  common.HexToHash(test_data.RandomString(64)),
			}

			err := repo.InsertDiff(diff)

			Expect(err).NotTo(HaveOccurred())
			var persistedDiff types.PersistedDiff
			readErr := db.Get(&persistedDiff, `SELECT * FROM public.storage_diff LIMIT 1`)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(persistedDiff.RawDiff).To(Equal(diff))
			Expect(persistedDiff.FromBackfill).To(BeTrue())
			Expect(persistedDiff.Checked).To(BeFalse())
		})

		It("doesn't duplicate diffs", func() {
			diff := types.RawDiff{
				HashedAddress: common.HexToHash(test_data.RandomString(64)),
				BlockHash:     common.HexToHash(test_data.RandomString(64)),
				BlockHeight:   rand.Int(),
				StorageKey:    common.HexToHash(test_data.RandomString(64)),
				StorageValue:  common.HexToHash(test_data.RandomString(64)),
			}

			err := repo.InsertDiff(diff)
			Expect(err).NotTo(HaveOccurred())

			errTwo := repo.InsertDiff(diff)
			Expect(errTwo).NotTo(HaveOccurred())

			var count int
			readErr := db.Get(&count, `SELECT count(*) FROM public.storage_diff`)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("VatIlkArtExists", func() {
		var (
			headerID, diffID int64
			ilkID            int
		)

		BeforeEach(func() {
			headerID = test_data.CreateTestHeader(db)
			diffID = test_helpers.CreateFakeDiffRecord(db)

			fakeIlk := test_data.RandomString(64)
			fakeIlkIdentifier := "ETH-A"
			ilkErr := db.Get(&ilkID, `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`, fakeIlk, fakeIlkIdentifier)
			Expect(ilkErr).NotTo(HaveOccurred())
		})

		It("returns true if vat_ilk_art for same ilk and header exists", func() {
			_, insertErr := db.Exec(`INSERT INTO maker.vat_ilk_art (diff_id, header_id, ilk_id, art) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, ilkID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatIlkArtExists(ilkID, int(headerID))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("returns false if vat_ilk_art only exists for same ilk at a different block", func() {
			_, insertErr := db.Exec(`INSERT INTO maker.vat_ilk_art (diff_id, header_id, ilk_id, art) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, ilkID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatIlkArtExists(ilkID, int(rand.Int31()))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("returns false if vat_urn_art only has rows for same header that are different urns", func() {
			var anotherIlkID int
			anotherIlkErr := db.Get(&anotherIlkID, `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2)
				RETURNING id`, test_data.RandomString(40), "BAT-A")
			Expect(anotherIlkErr).NotTo(HaveOccurred())

			_, insertErr := db.Exec(`INSERT INTO maker.vat_ilk_art (diff_id, header_id, ilk_id, art) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, anotherIlkID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatIlkArtExists(ilkID, int(headerID))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})

	Describe("VatUrnArtExists", func() {
		var (
			headerID, diffID int64
			urnID, ilkID     int
		)

		BeforeEach(func() {
			headerID = test_data.CreateTestHeader(db)
			diffID = test_helpers.CreateFakeDiffRecord(db)

			fakeIlk := test_data.RandomString(64)
			fakeIlkIdentifier := "ETH-A"
			ilkErr := db.Get(&ilkID, `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`, fakeIlk, fakeIlkIdentifier)
			Expect(ilkErr).NotTo(HaveOccurred())

			fakeUrn := test_data.RandomString(40)
			urnErr := db.Get(&urnID, `INSERT INTO maker.urns (ilk_id, identifier) VALUES ($1, $2) RETURNING id`, ilkID, fakeUrn)
			Expect(urnErr).NotTo(HaveOccurred())
		})

		It("returns true if vat_urn_art for same urn and header exists", func() {
			_, insertErr := db.Exec(`INSERT INTO maker.vat_urn_art (diff_id, header_id, urn_id, art) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, urnID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatUrnArtExists(urnID, int(headerID))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("returns false if vat_urn_art only exists for same urn at a different block", func() {
			_, insertErr := db.Exec(`INSERT INTO maker.vat_urn_art (diff_id, header_id, urn_id, art) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, urnID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatUrnArtExists(urnID, int(rand.Int31()))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("returns false if vat_urn_art only has rows for same header that are different urns", func() {
			var anotherUrnID int
			anotherUrnErr := db.Get(&anotherUrnID, `INSERT INTO maker.urns (ilk_id, identifier) VALUES ($1, $2)
				RETURNING id`, ilkID, test_data.RandomString(40))
			Expect(anotherUrnErr).NotTo(HaveOccurred())

			_, insertErr := db.Exec(`INSERT INTO maker.vat_urn_art (diff_id, header_id, urn_id, art) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, anotherUrnID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatUrnArtExists(urnID, int(headerID))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})

	Describe("VatUrnInkExists", func() {
		var (
			headerID, diffID int64
			urnID, ilkID     int
		)

		BeforeEach(func() {
			headerID = test_data.CreateTestHeader(db)
			diffID = test_helpers.CreateFakeDiffRecord(db)

			fakeIlk := test_data.RandomString(64)
			fakeIlkIdentifier := "ETH-A"
			ilkErr := db.Get(&ilkID, `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`, fakeIlk, fakeIlkIdentifier)
			Expect(ilkErr).NotTo(HaveOccurred())

			fakeUrn := test_data.RandomString(40)
			urnErr := db.Get(&urnID, `INSERT INTO maker.urns (ilk_id, identifier) VALUES ($1, $2) RETURNING id`, ilkID, fakeUrn)
			Expect(urnErr).NotTo(HaveOccurred())
		})

		It("returns true if vat_urn_ink for same urn and header exists", func() {
			_, insertErr := db.Exec(`INSERT INTO maker.vat_urn_ink (diff_id, header_id, urn_id, ink) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, urnID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatUrnInkExists(urnID, int(headerID))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("returns false if vat_urn_ink only exists for same urn at a different block", func() {
			_, insertErr := db.Exec(`INSERT INTO maker.vat_urn_ink (diff_id, header_id, urn_id, ink) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, urnID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatUrnInkExists(urnID, int(rand.Int31()))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("returns false if vat_urn_ink only has rows for same header that are different urns", func() {
			anotherFakeUrn := test_data.RandomString(40)
			var anotherUrnID int
			anotherUrnErr := db.Get(&anotherUrnID, `INSERT INTO maker.urns (ilk_id, identifier) VALUES ($1, $2) RETURNING id`, ilkID, anotherFakeUrn)
			Expect(anotherUrnErr).NotTo(HaveOccurred())

			_, insertErr := db.Exec(`INSERT INTO maker.vat_urn_ink (diff_id, header_id, urn_id, ink) VALUES
                            ($1, $2, $3, $4)`, diffID, headerID, anotherUrnID, 0)
			Expect(insertErr).NotTo(HaveOccurred())

			exists, err := repo.VatUrnInkExists(urnID, int(headerID))

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})
})
