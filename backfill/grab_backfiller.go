package backfill

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/sirupsen/logrus"
)

type GrabBackFiller interface {
	BackFillGrabStorage(startingBlock int) error
}

type grabBackFiller struct {
	blockChain        core.BlockChain
	eventsRepository  EventsRepository
	storageRepository StorageRepository
}

func NewGrabBackFiller(blockChain core.BlockChain, eventsRepository EventsRepository, storageRepository StorageRepository) GrabBackFiller {
	return grabBackFiller{
		blockChain:        blockChain,
		eventsRepository:  eventsRepository,
		storageRepository: storageRepository,
	}
}

func (backFiller grabBackFiller) BackFillGrabStorage(startingBlock int) error {
	grabs, grabsErr := backFiller.eventsRepository.GetGrabs(startingBlock)
	if grabsErr != nil {
		return fmt.Errorf("error getting grab events: %w", grabsErr)
	}
	logrus.Infof("getting storage for %d grabs", len(grabs))
	for i, grab := range grabs {
		dinkIsZero, dartIsZero, err := dinkAndDartAreZero(grab)
		if err != nil {
			return fmt.Errorf("error checking if dink/dart are zero: %w", err)
		}
		if dinkIsZero && dartIsZero {
			continue
		}

		urn, urnErr := backFiller.storageRepository.GetUrnByID(grab.UrnID)
		if urnErr != nil {
			return fmt.Errorf("error getting urn for id %d: %w", grab.UrnID, urnErr)
		}

		ilkArtExists, ilkArtErr := backFiller.storageRepository.VatIlkArtExists(urn.IlkID, grab.HeaderID)
		if ilkArtErr != nil {
			return fmt.Errorf("error checking if ilk art exists: %w", ilkArtErr)
		}
		urnArtExists, urnArtErr := backFiller.storageRepository.VatUrnArtExists(grab.UrnID, grab.HeaderID)
		if urnArtErr != nil {
			return fmt.Errorf("error checking if urn art exists: %w", urnArtErr)
		}
		urnInkExists, urnInkErr := backFiller.storageRepository.VatUrnInkExists(grab.UrnID, grab.HeaderID)
		if urnInkErr != nil {
			return fmt.Errorf("error checking if urn ink exists: %w", urnInkErr)
		}
		if !needToBackFillDiffsForGrab(dinkIsZero, dartIsZero, ilkArtExists, urnArtExists, urnInkExists) {
			continue
		}

		header, headerErr := backFiller.eventsRepository.GetHeaderByID(grab.HeaderID)
		if headerErr != nil {
			return fmt.Errorf("error getting header for id %d: %w", grab.HeaderID, headerErr)
		}

		keys, keysErr := getGrabKeys(urn, dinkIsZero, dartIsZero, ilkArtExists, urnArtExists, urnInkExists)
		if keysErr != nil {
			return fmt.Errorf("error generating storage keys: %w", keysErr)
		}

		logrus.Infof("feching diffs for grab %d, getting %d keys for urn %s", i, len(keys), urn.Urn)
		insertErr := backFiller.getAndPersistDiffs(keys, header)
		if insertErr != nil {
			return fmt.Errorf("error getting and persisting keys for grab: %w", insertErr)
		}
	}
	return nil
}

func (backFiller grabBackFiller) getAndPersistDiffs(keys []common.Hash, header core.Header) error {
	storageKeysToValues, storageErr := backFiller.blockChain.BatchGetStorageAt(VatAddress, keys,
		big.NewInt(header.BlockNumber))
	if storageErr != nil {
		return fmt.Errorf("error getting storage value: %w", storageErr)
	}
	for k, v := range storageKeysToValues {
		diff := types.RawDiff{
			HashedAddress: HashedVatAddress,
			BlockHash:     common.HexToHash(header.Hash),
			BlockHeight:   int(header.BlockNumber),
			StorageKey:    crypto.Keccak256Hash(k.Bytes()),
			StorageValue:  common.BytesToHash(v),
		}
		createDiffErr := backFiller.storageRepository.InsertDiff(diff)
		if createDiffErr != nil {
			return fmt.Errorf("error inserting diff: %w", createDiffErr)
		}
	}
	return nil
}

func getGrabKeys(urn Urn, dinkIsZero, dartIsZero, ilkArtExists, urnArtExists, urnInkExists bool) ([]common.Hash, error) {
	var keys []common.Hash
	if needToBackFillDiffsForIlkArt(dartIsZero, ilkArtExists) {
		keys = append(keys, getIlkArtKey(urn.Ilk))
	}
	if needToBackFillDiffsForUrnArt(dartIsZero, urnArtExists) {
		urnArtKey, keyErr := getUrnArtKey(urn)
		if keyErr != nil {
			return nil, fmt.Errorf("error getting urn art key: %w", keyErr)
		}
		keys = append(keys, urnArtKey)
	}
	if needToBackFillDiffsForUrnInk(dinkIsZero, urnInkExists) {
		urnInkKey, keyErr := getUrnInkKey(urn)
		if keyErr != nil {
			return nil, fmt.Errorf("error getting urn ink key: %w", keyErr)
		}
		keys = append(keys, urnInkKey)
	}
	return keys, nil
}

func needToBackFillDiffsForGrab(dinkIsZero, dartIsZero, ilkArtExists, urnArtExists, urnInkExists bool) bool {
	return needToBackFillDiffsForIlkArt(dartIsZero, ilkArtExists) ||
		needToBackFillDiffsForUrnArt(dartIsZero, urnArtExists) ||
		needToBackFillDiffsForUrnInk(dinkIsZero, urnInkExists)
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

func dinkAndDartAreZero(grab Grab) (bool, bool, error) {
	dink, dinkErr := stringToBigInt(grab.Dink)
	if dinkErr != nil {
		return false, false, fmt.Errorf("error formatting dink: %w", dinkErr)
	}
	dart, dartErr := stringToBigInt(grab.Dart)
	if dartErr != nil {
		return false, false, fmt.Errorf("error formatting dart: %w", dartErr)
	}
	return isZero(dink), isZero(dart), nil
}
