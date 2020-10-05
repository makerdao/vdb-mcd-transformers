package zero_value_diff

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
	"github.com/sirupsen/logrus"
)

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
			Address:      common.HexToAddress(test_data.Cat110Address()),
			BlockHash:    common.HexToHash(header.Hash),
			BlockHeight:  int(header.BlockNumber),
			StorageKey:   common.HexToHash(key),
			StorageValue: common.Hash{},
		}
		_, createDiffErr := generator.DiffRepo.CreateStorageDiff(rawDiff)
		if createDiffErr != nil {
			return createDiffErr
		}
	}

	return nil
}
