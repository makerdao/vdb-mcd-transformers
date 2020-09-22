package cmd

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_0_0"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	blockNumber int
	ilks        []string
)

var assignIlkLumpToZeroCmd = &cobra.Command{
	Use:   "assignIlkLumpToZero",
	Short: "Create a diff with a zero ilk lump value.",
	Long: `Inserts artificial cat_ilk_lump storage_diffs with storage_values of zero for the given ilks at the given block
height. The storage_diffs will then be picked up in the execute process to be transformed into maker.cat_ilk_lump records.
The block-number (b) and ilks (i) arguments are required.

Note: the ilk argument(s) passed in need to be the hex representation of the ilk.

Note: this command should only be run for the block height when the Cat contract upgrade change takes place, i.e. this
should only be run for the block where the first cat_ilk_dunk was set for the given ilks.
`,
	Example: `./vdb-mcd-transformers assignIlkLumpToZero --blockHeight=10769102
		--ilks 0x4241542d41000000000000000000000000000000000000000000000000000000
		--ilks 0x4554482d41000000000000000000000000000000000000000000000000000000
	`,
	PreRun: setViperConfigs,
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
	assignIlkLumpToZeroCmd.Flags().StringSliceVarP(&ilks, "ilks", "i", []string{}, "the ilk to set lump to zero")
	assignIlkLumpToZeroCmd.MarkFlagRequired("ilk")
}

func zeroOutIlkLump() error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	generator := NewZeroValueDiffGenerator(db)
	return generator.CreateZeroValueIlkLumpDiff(blockNumber, ilks)
}

type ZeroValueDiffGenerator struct {
	HeaderRepository datastore.HeaderRepository
	DiffRepo         storage.DiffRepository
}

func NewZeroValueDiffGenerator(db postgres.DB) ZeroValueDiffGenerator {
	return ZeroValueDiffGenerator{
		HeaderRepository: repositories.NewHeaderRepository(&db),
		DiffRepo:         storage.NewDiffRepository(&db),
	}
}

func (generator ZeroValueDiffGenerator) CreateZeroValueIlkLumpDiff(blockNumber int, ilks []string) error {
	keys, getKeysErr := generator.getIlkLumpKey(ilks)
	if getKeysErr != nil {
		return fmt.Errorf("error getting ilk lump keys %w", getKeysErr)
	}

	header, getHeaderErr := generator.HeaderRepository.GetHeaderByBlockNumber(int64(blockNumber))
	if getHeaderErr != nil {
		return fmt.Errorf("error gettting header %w", getHeaderErr)
	}
	logrus.Infof("Creating zero value lump diffs for the following ilks: %s", strings.Join(ilks, ", "))
	return generator.createDiffs(keys, header)
}

func (generator *ZeroValueDiffGenerator) getIlkLumpKey(ilks []string) ([]string, error) {
	var keys []string
	for _, ilk := range ilks {
		keys = append(keys, v1_0_0.GetIlkLumpKey(ilk).Hex())
	}

	return keys, nil
}

func (generator ZeroValueDiffGenerator) createDiffs(keys []string, header core.Header) error {
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
