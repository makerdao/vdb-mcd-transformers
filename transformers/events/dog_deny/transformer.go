package dog_deny

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

//func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]DogDenyEntity, error) {
//	return []DogDenyEntity{}, nil
//}
func (t Transformer) ToModels(contractAbi string, ethLog []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	panic("implement me")
}
