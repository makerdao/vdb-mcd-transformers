package backfill

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

var (
	HashedVatAddress = crypto.Keccak256Hash(common.HexToAddress("0x35d1b3f3d7966a1dfe207aa4514c12a259a0492b").Bytes())
	VatAddress       = common.HexToAddress("0x35d1b3f3d7966a1dfe207aa4514c12a259a0492b")
)

type FrobBackFiller interface {
	BackFillFrobStorage(startingBlock int) error
}

type frobBackFiller struct {
	db                *postgres.DB
	blockChain        core.BlockChain
	eventsRepository  EventsRepository
	storageRepository StorageRepository
}

func NewFrobBackFiller(blockChain core.BlockChain, eventsRepository EventsRepository, storageRepository StorageRepository) FrobBackFiller {
	return frobBackFiller{
		blockChain:        blockChain,
		eventsRepository:  eventsRepository,
		storageRepository: storageRepository,
	}
}

func (backFiller frobBackFiller) BackFillFrobStorage(startingBlock int) error {
	urns, urnsErr := backFiller.storageRepository.GetUrns()
	if urnsErr != nil {
		return fmt.Errorf("failed getting urns: %w", urnsErr)
	}

	lenUrns := len(urns)
	logrus.Printf("getting frobs for %d urns\n", lenUrns)
	for i, u := range urns {
		frobs, frobsErr := backFiller.eventsRepository.GetFrobs(u.UrnID, startingBlock)
		if frobsErr != nil {
			return fmt.Errorf("error getting frobs for urn %d: %w", u.UrnID, frobsErr)
		}

		logrus.Infof("getting %d out of %d urns: %s, %s - %d frobs", i, lenUrns, u.Ilk, u.Urn, len(frobs))
		for _, f := range frobs {
			err := backFiller.backFillFrob(f, u)
			if err != nil {
				return fmt.Errorf("error backfilling frob: %w", err)
			}
		}
	}
	return nil
}

func (backFiller frobBackFiller) backFillFrob(f Frob, u Urn) error {
	dink, dinkErr := stringToBigInt(f.Dink)
	if dinkErr != nil {
		return fmt.Errorf("error formatting dink: %w", dinkErr)
	}
	dart, dartErr := stringToBigInt(f.Dart)
	if dartErr != nil {
		return fmt.Errorf("error formatting dart: %w", dartErr)
	}
	if isZero(dink) && isZero(dart) {
		return nil
	}

	header, headerErr := backFiller.eventsRepository.GetHeaderByID(f.HeaderID)
	if headerErr != nil {
		return fmt.Errorf("error getting header for id %d: %w", f.HeaderID, headerErr)
	}

	if !isZero(dart) {
		keys, keysErr := backFiller.getDartKeys(u, header)
		if keysErr != nil {
			return fmt.Errorf("error generating keys for non-zero dart: %w", keysErr)
		}

		if len(keys) > 0 {
			insertErr := backFiller.getAndPersistDiffs(keys, header)
			if insertErr != nil {
				return insertErr
			}
		}
	}

	if !isZero(dink) {
		insertErr := backFiller.addVatUrnInkIfDoesNotExist(u, header)
		if insertErr != nil {
			return insertErr
		}
	}

	return nil
}

func (backFiller frobBackFiller) getDartKeys(u Urn, header core.Header) ([]common.Hash, error) {
	var keys []common.Hash

	ilkArtExists, ilkArtErr := backFiller.storageRepository.VatIlkArtExists(u.IlkID, int(header.Id))
	if ilkArtErr != nil {
		return nil, fmt.Errorf("error checking if vat ilk art exists: %w", ilkArtExists)
	}
	if !ilkArtExists {
		keys = append(keys, getIlkArtKey(u.Ilk))
	}

	urnArtExists, urnArtErr := backFiller.storageRepository.VatUrnArtExists(u.UrnID, int(header.Id))
	if urnArtErr != nil {
		return nil, fmt.Errorf("error checking if vat urn art exists: %w", urnArtErr)
	}
	if !urnArtExists {
		urnArtKey, keyErr := getUrnArtKey(u)
		if keyErr != nil {
			return nil, fmt.Errorf("error getting urn art key: %w", keyErr)
		}
		keys = append(keys, urnArtKey)
	}

	return keys, nil
}

func (backFiller frobBackFiller) addVatUrnInkIfDoesNotExist(urn Urn, header core.Header) error {
	exists, existErr := backFiller.storageRepository.VatUrnInkExists(urn.UrnID, int(header.Id))
	if existErr != nil {
		return fmt.Errorf("error checking if vat urn ink exists: %w", existErr)
	}
	if exists {
		return nil
	}
	inkKey, keyErr := getUrnInkKey(urn)
	if keyErr != nil {
		return fmt.Errorf("error getting ink key: %w", keyErr)
	}
	return backFiller.getAndPersistDiffs([]common.Hash{inkKey}, header)
}

func (backFiller frobBackFiller) getAndPersistDiffs(keys []common.Hash, header core.Header) error {
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

func getIlkArtKey(ilk string) common.Hash {
	return storage.GetKeyForMapping(storage.IndexTwo, ilk)
}

func getUrnArtKey(u Urn) (common.Hash, error) {
	inkKey, err := getUrnInkKey(u)
	if err != nil {
		return common.Hash{}, err
	}
	artKey := storage.GetIncrementedKey(inkKey, 1)
	return artKey, nil
}

func getUrnInkKey(u Urn) (common.Hash, error) {
	paddedGuy, padErr := utilities.PadAddress(u.Urn)
	if padErr != nil {
		return common.Hash{}, fmt.Errorf("error padding urn identifier: %w", padErr)
	}
	inkKey := storage.GetKeyForNestedMapping(storage.IndexThree, u.Ilk, paddedGuy)
	return inkKey, nil
}

func isZero(n *big.Int) bool {
	zero := big.NewInt(0)
	return n.Cmp(zero) == 0
}

func stringToBigInt(s string) (*big.Int, error) {
	n, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return nil, errors.New("error formatting string as *big.Int")
	}
	return n, nil
}
