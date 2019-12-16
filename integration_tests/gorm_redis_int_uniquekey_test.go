package integration_tests

import (
	"github.com/merlinapp/datarepo-go/integration_tests/bootstrap"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"github.com/merlinapp/datarepo-go/integration_tests/testdomain"
	"github.com/satori/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GormRedisIntegrationUniqueKeyTestSuite struct {
	suite.Suite
	system *bootstrap.SystemInstance
}

func TestGormRedisIntegrationUniqueKeyTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(GormRedisIntegrationUniqueKeyTestSuite))
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) TestGetNonExistentBook() {
	ctx := s.system.Ctx

	Convey("Scenario: Get non-existent book", s.T(), func() {
		Convey("Given no books in the system", func() {
			// no-op: database starts with no data in the test transaction, so no need
			// to do anything

			Convey("When a book is that doesn't exist is fetched", func() {
				bookId := uuid.NewV4().String()
				s.system.CacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				result, err := s.system.CachedRepo.FindByKey(ctx, "ID", bookId)

				Convey("Then the result should be empty, "+
					"And the database should've been queried, "+
					"And there should be a cache miss", func() {
					So(err, ShouldBeNil)
					So(result.IsEmpty(), ShouldBeTrue)
					So(s.system.CacheStore.Miss(), ShouldEqual, 1)
					So(s.system.UniqueKeyDataFetcher.Reads(), ShouldEqual, 1)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) TestGetNonExistentBooks() {
	ctx := s.system.Ctx

	Convey("Scenario: Get non-existent books", s.T(), func() {
		Convey("Given no books in the system", func() {
			// no-op: database starts with no data in the test transaction, so no need
			// to do anything

			Convey("When multiple books that don't exist are fetched", func() {
				bookId1 := uuid.NewV4().String()
				bookId2 := uuid.NewV4().String()
				ids := []string{bookId1, bookId2}
				s.system.CacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				results, err := s.system.CachedRepo.FindByKeys(ctx, "ID", ids)

				Convey("Then the results should be empty, "+
					"And the database should've been queried for 2 missing ids, "+
					"And there should be 2 cache misses", func() {
					So(err, ShouldBeNil)
					So(results[0].IsEmpty(), ShouldBeTrue)
					So(results[1].IsEmpty(), ShouldBeTrue)
					So(s.system.CacheStore.Miss(), ShouldEqual, 2)
					So(s.system.UniqueKeyDataFetcher.Reads(), ShouldEqual, 2)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) TestCreateBook() {
	ctx := s.system.Ctx

	Convey("Scenario: Create a book", s.T(), func() {
		Convey("Given no books in the system", func() {
			// no-op: database starts with no data in the test transaction, so no need
			// to do anything

			Convey("When a book is created", func() {
				author := testdomain.CreateAuthor(s.system)
				book, err := author.CreateBook(ctx, EmptyStatus)

				Convey("The book should be created successfully, "+
					"And the book should appear in the database, "+
					"And the book should appear in the cache", func() {
					So(err, ShouldBeNil)
					So(book.DBBook.AuthorID, ShouldEqual, author.AuthorId)
					So(book.VerifyBookExists(ctx), ShouldBeTrue)
					So(book.VerifyBookIsCached(ctx), ShouldBeTrue)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) TestUpdateBook() {
	ctx := s.system.Ctx

	Convey("Scenario: Update a book", s.T(), func() {
		Convey("Given a books exists the system", func() {
			author := testdomain.CreateAuthor(s.system)
			book, _ := author.CreateBook(ctx, EmptyStatus)

			Convey("When the book is updated", func() {
				err := book.UpdateStatus(ctx, InProgressStatus)

				Convey("The book should be updated successfully, "+
					"And the updated book should appear in the database, "+
					"And the updated book should appear in the cache", func() {
					So(err, ShouldBeNil)
					So(book.DBBook.AuthorID, ShouldEqual, author.AuthorId)
					So(book.VerifyBookExists(ctx), ShouldBeTrue)
					So(book.VerifyBookIsCached(ctx), ShouldBeTrue)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) TestGetExistentAndNonExistentBooks() {
	ctx := s.system.Ctx

	Convey("Scenario: Get an existent and a non-existent book", s.T(), func() {
		Convey("Given 1 book in the system", func() {
			author := testdomain.CreateAuthor(s.system)
			book, _ := author.CreateBook(ctx, EmptyStatus)

			Convey("When the book is queried with a book that doesn't exist", func() {
				bookId2 := uuid.NewV4().String()
				ids := []string{book.BookId, bookId2}
				s.system.CacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				results, err := s.system.CachedRepo.FindByKeys(ctx, "ID", ids)

				Convey("Then the results should contain the existent book, "+
					"And the database should've been queried for 1 missing id, "+
					"And there should be 1 cache misses", func() {
					So(err, ShouldBeNil)
					var b model.Book
					results[0].InjectResult(&b)
					So(b.ID, ShouldEqual, book.DBBook.ID)
					So(b.AuthorID, ShouldEqual, book.DBBook.AuthorID)
					So(b.Status, ShouldEqual, book.DBBook.Status)
					So(results[1].IsEmpty(), ShouldBeTrue)
					So(s.system.CacheStore.Miss(), ShouldEqual, 1)
					So(s.system.UniqueKeyDataFetcher.Reads(), ShouldEqual, 1)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) TestGetExistentAndNonExistentBooksNoCache() {
	ctx := s.system.Ctx

	Convey("Scenario: Get an existent and a non-existent book", s.T(), func() {
		Convey("Given 1 book in the system", func() {
			author := testdomain.CreateAuthor(s.system)
			book, _ := author.CreateBook(ctx, EmptyStatus)

			Convey("When the first book is evicted from the cache, "+
				"And the book is queried with a book that doesn't exist", func() {
				book.ClearCacheData(ctx)
				bookId2 := uuid.NewV4().String()
				ids := []string{book.BookId, bookId2}
				s.system.CacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				results, err := s.system.CachedRepo.FindByKeys(ctx, "ID", ids)

				Convey("Then the results should contain the existent book, "+
					"And the database should've been queried for 2 missing ids, "+
					"And there should be 2 cache misses", func() {
					So(err, ShouldBeNil)
					var b model.Book
					results[0].InjectResult(&b)
					So(b.ID, ShouldEqual, book.DBBook.ID)
					So(b.AuthorID, ShouldEqual, book.DBBook.AuthorID)
					So(b.Status, ShouldEqual, book.DBBook.Status)
					So(results[1].IsEmpty(), ShouldBeTrue)
					So(s.system.CacheStore.Miss(), ShouldEqual, 2)
					So(s.system.UniqueKeyDataFetcher.Reads(), ShouldEqual, 2)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) SetupTest() {
	s.system = bootstrap.StartSystemForIntegrationTests()
	bootstrap.PrepareTestDB()
}

func (s *GormRedisIntegrationUniqueKeyTestSuite) TearDownTest() {
	bootstrap.RollbackTestDb()
}
