package bookcategory_gorm_memory

import (
	"github.com/merlinapp/datarepo-go/integration_tests/bookcategory_gorm_numericid/testdomain"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GormNumericIdIntegrationUniqueKeyTestSuite struct {
	suite.Suite
	system *testdomain.SystemInstance
}

func TestGormNumericIdIntegrationUniqueKeyTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(GormNumericIdIntegrationUniqueKeyTestSuite))
}

func (s *GormNumericIdIntegrationUniqueKeyTestSuite) TestGetNonExistentBookCategory() {
	ctx := s.system.Ctx

	Convey("Scenario: Get non-existent book category", s.T(), func() {
		Convey("Given no book categories exist in the system", func() {
			// no-op: database starts with no data in the test transaction, so no need
			// to do anything

			Convey("When a book category that doesn't exist is fetched", func() {
				bookCategoryId := 1
				s.system.BookCategoryCacheStore.ClearStats()
				result, err := s.system.BookCategoryRepo.FindByKey(ctx, "ID", bookCategoryId)

				Convey("Then the result should be empty, "+
					"And there should be a cache miss", func() {
					So(err, ShouldBeNil)
					So(result.IsEmpty(), ShouldBeTrue)
					So(s.system.BookCategoryCacheStore.Miss(), ShouldEqual, 1)
				})
			})
		})
	})
}

func (s *GormNumericIdIntegrationUniqueKeyTestSuite) TestGetExistentBookCategory() {
	ctx := s.system.Ctx
	categoryName := "Fiction"

	Convey("Scenario: Get existent book category", s.T(), func() {
		Convey("Given a book category exists in the system", func() {
			bookCategory, _ := testdomain.CreateBookCategory(s.system, categoryName)
			bookCategory.ClearCacheData(ctx)

			Convey("When the book category is fetched twice", func() {

				s.system.BookCategoryCacheStore.ClearStats()
				result, err := s.system.BookCategoryRepo.FindByKey(ctx, "ID", bookCategory.BookCategoryId)
				result2, err2 := s.system.BookCategoryRepo.FindByKey(ctx, "ID", bookCategory.BookCategoryId)

				Convey("Then the book category should be fetched in both occasions, "+
					"And there should be a cache miss, "+
					"And there should be a cache hit", func() {
					So(err, ShouldBeNil)
					So(err2, ShouldBeNil)
					So(result.IsEmpty(), ShouldBeFalse)
					So(result2.IsEmpty(), ShouldBeFalse)
					So(s.system.BookCategoryCacheStore.Miss(), ShouldEqual, 1)
					So(s.system.BookCategoryCacheStore.Hits(), ShouldEqual, 1)
					So(bookCategory.VerifyBookCategoryIsCached(ctx), ShouldBeTrue)
				})
			})
		})
	})
}

func (s *GormNumericIdIntegrationUniqueKeyTestSuite) TestGetExistentAndNonExistentBookCategories() {
	ctx := s.system.Ctx
	categoryName := "Fiction"

	Convey("Scenario: Get existent and non-existent book categories", s.T(), func() {
		Convey("Given a book category exists in the system", func() {
			bookCategory, _ := testdomain.CreateBookCategory(s.system, categoryName)

			Convey("When the book category is queried with one that doesn't exist", func() {

				s.system.BookCategoryCacheStore.ClearStats()
				bookTypeId2 := 1000
				ids := []int{bookCategory.BookCategoryId, bookTypeId2}
				results, err := s.system.BookCategoryRepo.FindByKeys(ctx, "ID", ids)

				Convey("Then the book categories should contain the category that exists, "+
					"And there should be a cache miss", func() {
					So(err, ShouldBeNil)
					var b model.BookCategory
					results[0].InjectResult(&b)
					So(b.ID, ShouldEqual, bookCategory.BookCategoryId)
					So(b.Name, ShouldEqual, bookCategory.DBBookCategory.Name)
					So(results[1].IsEmpty(), ShouldBeTrue)
					So(s.system.BookCategoryCacheStore.Miss(), ShouldEqual, 1)
					So(bookCategory.VerifyBookCategoryIsCached(ctx), ShouldBeTrue)
				})
			})
		})
	})
}

func (s *GormNumericIdIntegrationUniqueKeyTestSuite) SetupTest() {
	s.system = startSystemForIntegrationTests()
	prepareTestDB()
}

func (s *GormNumericIdIntegrationUniqueKeyTestSuite) TearDownTest() {
	rollbackTestDb()
}
