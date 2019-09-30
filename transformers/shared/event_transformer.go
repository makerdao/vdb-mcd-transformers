// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package shared

import (
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type EventTransformer struct {
	Config     transformer.EventTransformerConfig
	Converter  Converter
	Repository SharedRepository
}

func (tr EventTransformer) NewEventTransformer(db *postgres.DB) transformer.EventTransformer {
	tr.Repository.SetDB(db)
	return tr
}

func (tr EventTransformer) Execute(logs []core.HeaderSyncLog) error {
	transformerName := tr.Config.TransformerName

	// No matching logs, mark the header as checked for this type of logs
	if len(logs) < 1 {
		return nil
	}

	models, err := tr.Converter.ToModels(tr.Config.ContractAbi, logs)
	if err != nil {
		log.Printf("Error converting logs in %v: %v", transformerName, err)
		return err
	}

	err = tr.Repository.Create(models)
	if err != nil {
		log.Printf("Error persisting %v record: %v", transformerName, err)
		return err
	}
	return nil
}

func (tr EventTransformer) GetName() string {
	return tr.Config.TransformerName
}

func (tr EventTransformer) GetConfig() transformer.EventTransformerConfig {
	return tr.Config
}
