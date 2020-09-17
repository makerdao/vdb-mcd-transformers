package cmd_test

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/cmd"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assigns Ilk Lump to zero", func() {
	var (
		generator            cmd.ZeroValueDiffGenerator
		mockDiffRepo         mocks.MockStorageDiffRepository
		mockHeaderRepository fakes.MockHeaderRepository
		blockNumber          int

		ethA = "0x4554482d41000000000000000000000000000000000000000000000000000000"
		batA = "0x4241542d41000000000000000000000000000000000000000000000000000000"

		ethAKey = crypto.Keccak256Hash(common.FromHex(ethA + cat.IlksMappingIndex))
		batAKey = crypto.Keccak256Hash(common.FromHex(batA + cat.IlksMappingIndex))

		ethALumpKey = storage.GetIncrementedKey(ethAKey, 2).Hex()
		batALumpKey = storage.GetIncrementedKey(batAKey, 2).Hex()

		ilks = []string{ethA, batA}
	)

	Describe("CreateZeroValueIlkLumpDiff", func() {
		BeforeEach(func() {
			mockHeaderRepository = fakes.MockHeaderRepository{}
			mockDiffRepo = mocks.MockStorageDiffRepository{}
			blockNumber = rand.Int()
			generator = cmd.ZeroValueDiffGenerator{
				HeaderRepository: &mockHeaderRepository,
				DiffRepo:         &mockDiffRepo,
			}
		})

		It("looks up the header for the given block number", func() {
			err := generator.CreateZeroValueIlkLumpDiff(blockNumber, ilks)
			Expect(err).NotTo(HaveOccurred())
			Expect(mockHeaderRepository.GetHeaderPassedBlockNumber).To(Equal(int64(blockNumber)))
		})

		It("returns an error if getting the header fails", func() {
			mockHeaderRepository.GetHeaderByBlockNumberError = fakes.FakeError
			err := generator.CreateZeroValueIlkLumpDiff(blockNumber, ilks)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("inserts an artificial diff into the database for each ilk at the given block", func() {
			mockHeaderRepository.GetHeaderByBlockNumberReturnHash = fakes.FakeHash.Hex()

			err := generator.CreateZeroValueIlkLumpDiff(blockNumber, ilks)
			Expect(err).NotTo(HaveOccurred())

			ethAIlkLumpRawDiff := types.RawDiff{
				HashedAddress: common.HexToHash(test_data.Cat110Address()),
				BlockHash:     fakes.FakeHash,
				BlockHeight:   blockNumber,
				StorageKey:    common.HexToHash(ethALumpKey),
				StorageValue:  common.Hash{},
			}

			batAIlkLumpRawDiff := types.RawDiff{
				HashedAddress: common.HexToHash(test_data.Cat110Address()),
				BlockHash:     fakes.FakeHash,
				BlockHeight:   blockNumber,
				StorageKey:    common.HexToHash(batALumpKey),
				StorageValue:  common.Hash{},
			}
			expectedDiffs := []types.RawDiff{ethAIlkLumpRawDiff, batAIlkLumpRawDiff}
			Expect(mockDiffRepo.CreatePassedRawDiffs).To(Equal(expectedDiffs))
		})
	})
})
