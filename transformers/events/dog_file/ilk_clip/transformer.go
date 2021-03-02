package ilk_clip

import (
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.EventLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var models [] event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.)
	}
}
