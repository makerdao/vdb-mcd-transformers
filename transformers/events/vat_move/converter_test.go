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

package vat_move_test

import (
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_move"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
)

var _ = Describe("Vat move converter", func() {
	var converter vat_move.VatMoveConverter

	BeforeEach(func() {
		converter = vat_move.VatMoveConverter{}
	})

	It("returns err if logs are missing topics", func() {
		badLog := core.HeaderSyncLog{}
		_, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatMoveHeaderSyncLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.VatMoveModel))
	})
})
