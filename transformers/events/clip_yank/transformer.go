package clip_yank

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]ClipYankEntity, error) {
	panic("implement me")
}
func (t Transformer) ToModels(contractAbi string, ethLog []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	panic("implement me")

}
