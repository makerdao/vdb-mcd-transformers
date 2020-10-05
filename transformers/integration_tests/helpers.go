// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package integration_tests

import (
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
)

// Persist the header for a given block to postgres. Returns the header if successful.
func persistHeader(db *postgres.DB, blockNumber int64, blockChain core.BlockChain) (core.Header, error) {
	header, err := blockChain.GetHeaderByNumber(blockNumber)
	if err != nil {
		return core.Header{}, err
	}
	headerRepository := repositories.NewHeaderRepository(db)
	id, err := headerRepository.CreateOrUpdateHeader(header)
	header.Id = id
	return header, err
}
