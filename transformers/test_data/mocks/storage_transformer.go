package mocks

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockStorageTransformer struct {
	Address    common.Address
	ExecuteErr error
	PassedRow  utils.StorageDiffRow
}

func (transformer *MockStorageTransformer) Execute(row utils.StorageDiffRow) error {
	transformer.PassedRow = row
	return transformer.ExecuteErr
}

func (transformer *MockStorageTransformer) ContractAddress() common.Address {
	return transformer.Address
}

func (transformer *MockStorageTransformer) FakeTransformerInitializer(db *postgres.DB) transformer.StorageTransformer {
	return transformer
}
