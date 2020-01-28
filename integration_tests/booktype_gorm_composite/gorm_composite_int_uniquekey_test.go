package booktype_gorm_composite

import (
	"github.com/merlinapp/datarepo-go/integration_tests/booktype_gorm_composite/testdomain"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"github.com/satori/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GormCompositeIntegrationUniqueKeyTestSuite struct {
	suite.Suite
	system *testdomain.SystemInstance
}

func TestGormCompositeIntegrationUniqueKeyTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(GormCompositeIntegrationUniqueKeyTestSuite))
}

func (s *GormCompositeIntegrationUniqueKeyTestSuite) TestGetNonExistentBookType() {
	ctx := s.system.Ctx

	Convey("Scenario: Get non-existent book type", s.T(), func() {
		Convey("Given no book types in the system", func() {
			// no-op: database starts with no data in the test transaction, so no need
			// to do anything

			Convey("When a book type that doesn't exist is fetched", func() {
				bookTypeId := uuid.NewV4().String()
				s.system.BookTypeCacheStore.ClearStats()
				result, err := s.system.BookTypeRepo.FindByKey(ctx, "ID", bookTypeId)

				Convey("Then the result should be empty, "+
					"And there should be a cache miss", func() {
					So(err, ShouldBeNil)
					So(result.IsEmpty(), ShouldBeTrue)
					So(s.system.BookTypeCacheStore.Miss(), ShouldEqual, 1)
				})
			})
		})
	})
}

func (s *GormCompositeIntegrationUniqueKeyTestSuite) TestGetExistentBookType() {
	ctx := s.system.Ctx
	typeName := "Digital Book"

	Convey("Scenario: Get existent book type", s.T(), func() {
		Convey("Given a book type exists in the system", func() {
			bookType, _ := testdomain.CreateBookType(s.system, typeName)
			bookType.ClearCacheData(ctx)

			Convey("When the book type is fetched twice", func() {

				s.system.BookTypeCacheStore.ClearStats()
				result, err := s.system.BookTypeRepo.FindByKey(ctx, "ID", bookType.BookTypeId)
				result2, err2 := s.system.BookTypeRepo.FindByKey(ctx, "ID", bookType.BookTypeId)

				Convey("Then the book type should be fetched in both occasions, "+
					"And there should be a cache miss, "+
					"And there should be a cache hit", func() {
					So(err, ShouldBeNil)
					So(err2, ShouldBeNil)
					So(result.IsEmpty(), ShouldBeFalse)
					So(result2.IsEmpty(), ShouldBeFalse)
					So(s.system.BookTypeCacheStore.Miss(), ShouldEqual, 1)
					So(s.system.BookTypeCacheStore.Hits(), ShouldEqual, 1)
					So(bookType.VerifyBookTypeIsCached(ctx), ShouldBeTrue)
				})
			})
		})
	})
}

func (s *GormCompositeIntegrationUniqueKeyTestSuite) TestGetExistentAndNonExistentBookTypes() {
	ctx := s.system.Ctx
	typeName := "Digital Book"

	Convey("Scenario: Get existent and non-existent book types", s.T(), func() {
		Convey("Given a book type exists in the system", func() {
			bookType, _ := testdomain.CreateBookType(s.system, typeName)

			Convey("When the book type is queried with one that doesn't exist", func() {

				s.system.BookTypeCacheStore.ClearStats()
				bookTypeId2 := uuid.NewV4().String()
				ids := []string{bookType.BookTypeId, bookTypeId2}
				results, err := s.system.BookTypeRepo.FindByKeys(ctx, "ID", ids)

				Convey("Then the book types should contain the type that exists, "+
					"And there should be a cache miss", func() {
					So(err, ShouldBeNil)
					var b model.BookType
					results[0].InjectResult(&b)
					So(b.ID, ShouldEqual, bookType.BookTypeId)
					So(b.Name, ShouldEqual, bookType.DBBookType.Name)
					So(results[1].IsEmpty(), ShouldBeTrue)
					So(s.system.BookTypeCacheStore.Miss(), ShouldEqual, 1)
					So(bookType.VerifyBookTypeIsCached(ctx), ShouldBeTrue)
				})
			})
		})
	})
}

func (s *GormCompositeIntegrationUniqueKeyTestSuite) SetupTest() {
	s.system = startSystemForIntegrationTests()
	prepareTestDB()
}

func (s *GormCompositeIntegrationUniqueKeyTestSuite) TearDownTest() {
	rollbackTestDb()
}
