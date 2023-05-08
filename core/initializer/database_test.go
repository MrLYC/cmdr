package initializer_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/initializer"
	"github.com/mrlyc/cmdr/core/mock"
)

type TestModel struct {
	ID   int    `storm:"increment"`
	Name string `storm:"index" json:"name"`
}

var _ = Describe("Database", func() {
	var (
		ctrl   *gomock.Controller
		db     *mock.MockDatabase
		models = map[core.ModelType]interface{}{
			core.ModelTypeUnknown: TestModel{},
		}
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		db = mock.NewMockDatabase(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("DatabaseMigrator", func() {
		var (
			migrator *initializer.DatabaseMigrator
		)

		BeforeEach(func() {
			migrator = initializer.NewDatabaseMigrator(func() (core.Database, error) {
				return db, nil
			}, models)
		})

		It("should migrate models", func() {
			histories := make(map[string][]string)

			db.EXPECT().Init(gomock.Any()).DoAndReturn(func(data interface{}) error {
				key := fmt.Sprintf("%T", data)
				histories[key] = append(histories[key], "Init")
				return nil
			}).AnyTimes()
			db.EXPECT().ReIndex(gomock.Any()).DoAndReturn(func(data interface{}) error {
				key := fmt.Sprintf("%T", data)
				histories[key] = append(histories[key], "ReIndex")
				return nil
			}).AnyTimes()
			db.EXPECT().Close().Return(nil)

			Expect(migrator.Init(false)).To(Succeed())
			Expect(histories["initializer_test.TestModel"]).To(Equal([]string{"Init", "ReIndex"}))
		})
	})

})
