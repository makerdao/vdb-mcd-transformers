package cmd

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_0_0"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	storage2 "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var blockNumber int
var assignIlkLumpToZeroCmd = &cobra.Command{
	Use:   "assignIlkLumpToZero",
	Short: "Create a diff with a zero ilk lump value.",
	Long: `Inserts an artificial storage_diff with the storage_key for cat_ilk_lump, and a
storage_value of zero for all ilks for the given block height. This storage_diff will then be picked up in the execute
process to be transformed into maker.cat_ilk_lump record.

Note: this command should only be run for the block height when the Cat contract upgrade change takes place, i.e. this
should only be run for the block after the first cat_ilk_dunk was set.
`,
	Example: "",
	//"./vdb-mcd-transformers assignIlkLumpToZero --blockHeight=1000",
	ValidArgs: nil,
	Args:      nil,
	PreRun:    setViperConfigs,
	Run: func(cmd *cobra.Command, args []string) {
		err := zeroOutIlkLump()
		if err != nil {
			logrus.Fatalf("error inserting a zero ilk lump diff: %s", err.Error())
		}
		logrus.Infof("Successfully created a zero-value ilk lump diffs for block %d", blockNumber)
		return
	},
}

func init() {
	rootCmd.AddCommand(assignIlkLumpToZeroCmd)
	assignIlkLumpToZeroCmd.Flags().IntVarP(&blockNumber, "block-number", "b", -1, "the block associated with artificial ilk lump diffs")
	assignIlkLumpToZeroCmd.MarkFlagRequired("block-number")
}

func zeroOutIlkLump() error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	generator := NewZeroValueDiffGenerator(db)
	return generator.CreateZeroValueIlkLumpDiff(blockNumber)
}

type ZeroValueDiffGenerator struct {
	MakerStorageRepo storage.IMakerStorageRepository
	HeaderRepository datastore.HeaderRepository
	DiffRepo         storage2.DiffRepository
	ilks             []string
}

func NewZeroValueDiffGenerator(db postgres.DB) ZeroValueDiffGenerator {
	makerStorageRepo := storage.MakerStorageRepository{}
	makerStorageRepo.SetDB(&db)
	return ZeroValueDiffGenerator{
		MakerStorageRepo: &makerStorageRepo,
		HeaderRepository: repositories.NewHeaderRepository(&db),
		DiffRepo:         storage2.NewDiffRepository(&db),
	}
}

func (generator ZeroValueDiffGenerator) CreateZeroValueIlkLumpDiff(blockNumber int) error {
	keys, getKeysErr := generator.getIlkLumpKeys()
	if getKeysErr != nil {
		return fmt.Errorf("error getting ilk lump keys %w", getKeysErr)
	}

	header, getHeaderErr := generator.HeaderRepository.GetHeaderByBlockNumber(int64(blockNumber))
	if getHeaderErr != nil {
		return fmt.Errorf("error gettting header %w", getHeaderErr)
	}

	return generator.createDiff(keys, header)
}

func (generator *ZeroValueDiffGenerator) getIlkLumpKeys() ([]string, error) {
	ilks, getIlksErr := generator.MakerStorageRepo.GetIlks()
	if getIlksErr != nil {
		return nil, fmt.Errorf("error retriving ilks from db %w", getIlksErr)
	}
	generator.ilks = ilks

	var keys []string
	for _, ilk := range ilks {
		keys = append(keys, v1_0_0.GetIlkLumpKey(ilk).Hex())
	}

	return keys, nil
}

func (generator ZeroValueDiffGenerator) createDiff(keys []string, header core.Header) error {
	fmt.Println(generator.ilks)
	logrus.Infof("Creating zero value lump diffs for the following ilks: %s", strings.Join(generator.ilks, ", "))
	for _, key := range keys {
		rawDiff := types.RawDiff{
			HashedAddress: common.HexToHash(test_data.Cat110Address()),
			BlockHash:     common.HexToHash(header.Hash),
			BlockHeight:   int(header.BlockNumber),
			StorageKey:    common.HexToHash(key),
			StorageValue:  common.Hash{},
		}
		_, createDiffErr := generator.DiffRepo.CreateStorageDiff(rawDiff)
		if createDiffErr != nil {
			return createDiffErr
		}
	}

	return nil
}
