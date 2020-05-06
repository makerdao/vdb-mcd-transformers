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

type UrnBackFiller interface {
	BackfillUrns(startingBlock int) error
}

type urnBackFiller struct {
	db               *postgres.DB
	blockChain       core.BlockChain
	diffRepository   storage.DiffRepository
	eventsRepository EventsRepository
	urnsRepository   UrnsRepository
}

func NewUrnBackFiller(blockChain core.BlockChain, diffRepository storage.DiffRepository,
	eventsRepository EventsRepository, urnsRepository UrnsRepository) UrnBackFiller {
	return urnBackFiller{
		blockChain:       blockChain,
		diffRepository:   diffRepository,
		eventsRepository: eventsRepository,
		urnsRepository:   urnsRepository,
	}
}

func (backFiller urnBackFiller) BackfillUrns(startingBlock int) error {
	urns, urnsErr := backFiller.urnsRepository.GetUrns()
	if urnsErr != nil {
		return fmt.Errorf("failed getting urns: %w", urnsErr)
	}

	logrus.Printf("getting frobs for %d urns\n", len(urns))
	for _, u := range urns {
		frobs, frobsErr := backFiller.eventsRepository.GetFrobs(u.ID, startingBlock)
		if frobsErr != nil {
			return fmt.Errorf("error getting frobs for urn %d: %w", u.ID, frobsErr)
		}

		logrus.Printf("getting diffs for %d frobs for urn %s of ilk %s", len(frobs), u.Urn, u.Ilk)
		for _, f := range frobs {
			dink, ok := big.NewInt(0).SetString(f.Dink, 10)
			if !ok {
				return errors.New("error formatting dink for urn frob")
			}
			dart, ok := big.NewInt(0).SetString(f.Dart, 10)
			if !ok {
				return errors.New("error formatting dart for urn frob")
			}
			if isZero(dink) && isZero(dart) {
				continue
			}

			header, headerErr := backFiller.eventsRepository.GetHeaderByID(f.HeaderID)
			if headerErr != nil {
				return fmt.Errorf("error getting header for id %d: %w", f.HeaderID, headerErr)
			}

			if !isZero(dart) {
				insertErr := backFiller.addVatUrnArtIfDoesNotExist(u, header)
				if insertErr != nil {
					return insertErr
				}
			}

			if !isZero(dink) {
				insertErr := backFiller.addVatUrnInkIfDoesNotExist(u, header)
				if insertErr != nil {
					return insertErr
				}
			}
		}
	}
	return nil
}

func (backFiller urnBackFiller) addVatUrnArtIfDoesNotExist(urn Urn, header core.Header) error {
	exists, existErr := backFiller.urnsRepository.VatUrnArtExists(urn.ID, int(header.Id))
	if existErr != nil {
		return fmt.Errorf("error checking if vat urn art exists: %w", existErr)
	}
	if exists {
		return nil
	}
	artKey, keyErr := getArtKey(urn)
	if keyErr != nil {
		return fmt.Errorf("error getting art key: %w", keyErr)
	}
	return backFiller.getAndPersistDiff(artKey, header)
}

func (backFiller urnBackFiller) addVatUrnInkIfDoesNotExist(urn Urn, header core.Header) error {
	exists, existErr := backFiller.urnsRepository.VatUrnInkExists(urn.ID, int(header.Id))
	if existErr != nil {
		return fmt.Errorf("error checking if vat urn ink exists: %w", existErr)
	}
	if exists {
		return nil
	}
	inkKey, keyErr := getInkKey(urn)
	if keyErr != nil {
		return fmt.Errorf("error getting ink key: %w", keyErr)
	}
	return backFiller.getAndPersistDiff(inkKey, header)
}

func (backFiller urnBackFiller) getAndPersistDiff(key common.Hash, header core.Header) error {
	storageKeysToValues, storageErr := backFiller.blockChain.BatchGetStorageAt(VatAddress,
		[]common.Hash{key}, big.NewInt(header.BlockNumber))
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
		createDiffErr := backFiller.diffRepository.CreateBackFilledStorageValue(diff)
		if createDiffErr != nil {
			return fmt.Errorf("error inserting diff: %w", createDiffErr)
		}
	}
	return nil
}

func getInkKey(u Urn) (common.Hash, error) {
	paddedGuy, padErr := utilities.PadAddress(u.Urn)
	if padErr != nil {
		return common.Hash{}, fmt.Errorf("error padding urn identifier: %w", padErr)
	}
	inkKey := storage.GetKeyForNestedMapping(storage.IndexThree, u.Ilk, paddedGuy)
	return inkKey, nil
}

func getArtKey(u Urn) (common.Hash, error) {
	inkKey, err := getInkKey(u)
	if err != nil {
		return common.Hash{}, err
	}
	artKey := storage.GetIncrementedKey(inkKey, 1)
	return artKey, nil
}

func isZero(n *big.Int) bool {
	zero := big.NewInt(0)
	return n.Cmp(zero) == 0
}
