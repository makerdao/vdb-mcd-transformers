package shared

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/sirupsen/logrus"
)

var VatAddress = common.HexToAddress("0x35d1b3f3d7966a1dfe207aa4514c12a259a0492b")

type DartDink struct {
	Dart     string
	Dink     string
	HeaderID int64
	UrnID    int64
}

type DartDinkRetriever interface {
	RetrieveDartDinkDiffs(dartDink DartDink) error
}

type dartDinkRetriever struct {
	blockChain        core.BlockChain
	eventsRepository  repository.EventsRepository
	headerRepository  datastore.HeaderRepository
	storageRepository repository.StorageRepository
}

func NewDartDinkRetriever(blockChain core.BlockChain, eventsRepository repository.EventsRepository, headerRepository datastore.HeaderRepository, storageRepository repository.StorageRepository) DartDinkRetriever {
	return dartDinkRetriever{
		blockChain:        blockChain,
		eventsRepository:  eventsRepository,
		headerRepository:  headerRepository,
		storageRepository: storageRepository,
	}
}

func (retriever dartDinkRetriever) RetrieveDartDinkDiffs(dartDink DartDink) error {
	dartIsZero, dinkIsZero, err := dartAndDinkAreZero(dartDink.Dart, dartDink.Dink)
	if err != nil {
		return fmt.Errorf("error checking if dart/dink are zero: %w", err)
	}
	// if dart and dink are zero, there is no delta - so we wouldn't expect a diff
	if dartIsZero && dinkIsZero {
		return nil
	}

	urn, urnErr := retriever.storageRepository.GetUrnByID(dartDink.UrnID)
	if urnErr != nil {
		return fmt.Errorf("failed getting urn: %w", urnErr)
	}

	ilkArtExists, ilkArtErr := retriever.storageRepository.VatIlkArtExists(urn.IlkID, dartDink.HeaderID)
	if ilkArtErr != nil {
		return fmt.Errorf("error checking if ilk art exists: %w", ilkArtErr)
	}
	urnArtExists, urnArtErr := retriever.storageRepository.VatUrnArtExists(dartDink.UrnID, dartDink.HeaderID)
	if urnArtErr != nil {
		return fmt.Errorf("error checking if urn art exists: %w", urnArtErr)
	}
	urnInkExists, urnInkErr := retriever.storageRepository.VatUrnInkExists(dartDink.UrnID, dartDink.HeaderID)
	if urnInkErr != nil {
		return fmt.Errorf("error checking if urn ink exists: %w", urnInkErr)
	}
	if !needToBackFillDiffsForDartDink(dartIsZero, dinkIsZero, ilkArtExists, urnArtExists, urnInkExists) {
		return nil
	}

	header, headerErr := retriever.headerRepository.GetHeaderByID(dartDink.HeaderID)
	if headerErr != nil {
		return fmt.Errorf("error getting header for id %d: %w", dartDink.HeaderID, headerErr)
	}

	keys, keysErr := getDartDinkKeys(urn, dartIsZero, dinkIsZero, ilkArtExists, urnArtExists, urnInkExists)
	if keysErr != nil {
		return fmt.Errorf("error generating storage keys: %w", keysErr)
	}

	logrus.Infof("fetching %d keys for urn %s", len(keys), urn.Urn)
	insertErr := retriever.getAndPersistVatDiffs(keys, header)
	if insertErr != nil {
		return fmt.Errorf("error getting and persisting keys: %w", insertErr)
	}

	return nil
}

func (retriever dartDinkRetriever) getAndPersistVatDiffs(keys []common.Hash, header core.Header) error {
	storageKeysToValues, storageErr := retriever.blockChain.BatchGetStorageAt(VatAddress, keys,
		big.NewInt(header.BlockNumber))
	if storageErr != nil {
		return fmt.Errorf("error getting storage value: %w", storageErr)
	}
	for k, v := range storageKeysToValues {
		diff := types.RawDiff{
			Address:      VatAddress,
			BlockHash:    common.HexToHash(header.Hash),
			BlockHeight:  int(header.BlockNumber),
			StorageKey:   k,
			StorageValue: common.BytesToHash(v),
		}
		createDiffErr := retriever.storageRepository.InsertDiff(diff)
		if createDiffErr != nil {
			return fmt.Errorf("error inserting diff: %w", createDiffErr)
		}
	}
	return nil
}

func dartAndDinkAreZero(dart, dink string) (bool, bool, error) {
	dinkInt, dinkErr := StringToBigInt(dink)
	if dinkErr != nil {
		return false, false, fmt.Errorf("error formatting dink: %w", dinkErr)
	}
	dartInt, dartErr := StringToBigInt(dart)
	if dartErr != nil {
		return false, false, fmt.Errorf("error formatting dart: %w", dartErr)
	}
	return IsZero(dartInt), IsZero(dinkInt), nil
}

func needToBackFillDiffsForDartDink(dartIsZero, dinkIsZero, ilkArtExists, urnArtExists, urnInkExists bool) bool {
	return needToBackFillDiffsForIlkArt(dartIsZero, ilkArtExists) ||
		needToBackFillDiffsForUrnArt(dartIsZero, urnArtExists) ||
		needToBackFillDiffsForUrnInk(dinkIsZero, urnInkExists)
}

func getDartDinkKeys(urn repository.Urn, dartIsZero, dinkIsZero, ilkArtExists, urnArtExists, urnInkExists bool) ([]common.Hash, error) {
	var keys []common.Hash
	if needToBackFillDiffsForIlkArt(dartIsZero, ilkArtExists) {
		keys = append(keys, vat.GetIlkArtKey(urn.Ilk))
	}
	if needToBackFillDiffsForUrnArt(dartIsZero, urnArtExists) {
		urnArtKey, keyErr := GetUrnArtKey(urn)
		if keyErr != nil {
			return nil, fmt.Errorf("error getting urn art key: %w", keyErr)
		}
		keys = append(keys, urnArtKey)
	}
	if needToBackFillDiffsForUrnInk(dinkIsZero, urnInkExists) {
		urnInkKey, keyErr := GetUrnInkKey(urn)
		if keyErr != nil {
			return nil, fmt.Errorf("error getting urn ink key: %w", keyErr)
		}
		keys = append(keys, urnInkKey)
	}
	return keys, nil
}

func needToBackFillDiffsForUrnInk(dinkIsZero, urnInkExists bool) bool {
	return !dinkIsZero && !urnInkExists
}

func needToBackFillDiffsForUrnArt(dartIsZero, urnArtExists bool) bool {
	return !dartIsZero && !urnArtExists
}

func needToBackFillDiffsForIlkArt(dartIsZero, ilkArtExists bool) bool {
	return !dartIsZero && !ilkArtExists
}
