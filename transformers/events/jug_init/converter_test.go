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

package jug_init_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_init"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jug init converter", func() {
	var converter = jug_init.JugInitConverter{}

	It("returns err if log is missing topics", func() {
		incompleteLog := core.HeaderSyncLog{}
		_, err := converter.ToModels(constants.JugABI(), []core.HeaderSyncLog{incompleteLog})
		Expect(err).To(HaveOccurred())
	})

	It("convert a log to an insertion model", func() {
		models, err := converter.ToModels(constants.JugABI(), []core.HeaderSyncLog{test_data.JugInitHeaderSyncLog})
		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]shared.InsertionModel{test_data.JugInitModel}))
	})
})
