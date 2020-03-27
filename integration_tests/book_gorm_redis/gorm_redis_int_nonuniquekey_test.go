package book_gorm_redis

import (
	"github.com/merlinapp/datarepo-go/integration_tests/book_gorm_redis/testdomain"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GormRedisIntegrationNonUniqueKeyTestSuite struct {
	suite.Suite
	system *testdomain.SystemInstance
}

const EmptyStatus = ""
const InProgressStatus = "In Progress"
const CompletedStatus = "Completed"

func TestGormRedisIntegrationNonUniqueKeyTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(GormRedisIntegrationNonUniqueKeyTestSuite))
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TestGetEmptyBooksForAuthor() {
	ctx := s.system.Ctx

	Convey("Scenario: Get books for author with no books", s.T(), func() {
		Convey("Given no books in the system", func() {
			// no-op: database starts with no data in the test transaction, so no need
			// to do anything

			Convey("When the books for an author that has no books are fetched", func() {
				author := testdomain.CreateAuthor(s.system)
				s.system.BookCacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				books, err := author.GetBooks(ctx)

				Convey("Then the result should be empty, "+
					"And the database should've been queried, "+
					"And there should be a cache miss", func() {
					So(err, ShouldBeNil)
					So(books, ShouldBeEmpty)
					So(s.system.BookCacheStore.Miss(), ShouldEqual, 1)
					So(s.system.NonUniqueKeyDataFetcher.Reads(), ShouldEqual, 1)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TestGetEmptyBooksForAuthors() {
	Convey("Scenario: Get books for authors with no books", s.T(), func() {
		Convey("Given no books in the system", func() {
			// no-op: database starts with no data in the test transaction, so no need
			// to do anything

			Convey("When the books for 2 authors that have no books are fetched", func() {
				author1 := testdomain.CreateAuthor(s.system)
				author2 := testdomain.CreateAuthor(s.system)
				ids := []string{author1.AuthorId, author2.AuthorId}
				s.system.BookCacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				books, err := testdomain.GetBooksForAuthors(s.system, ids)

				Convey("Then the results should be empty, "+
					"And the database should've been queried for two missing ids, "+
					"And there should be 2 cache misses", func() {
					So(err, ShouldBeNil)
					So(len(books), ShouldEqual, 2)
					So(books[0], ShouldBeEmpty)
					So(books[1], ShouldBeEmpty)
					So(s.system.BookCacheStore.Miss(), ShouldEqual, 2)
					So(s.system.NonUniqueKeyDataFetcher.Reads(), ShouldEqual, 2)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TestGetAuthorBooksFromCache() {
	ctx := s.system.Ctx

	Convey("Scenario: Retrieve author books from cache", s.T(), func() {
		Convey("Given an author with 2 books in the system, "+
			"And I query for the books of the author once", func() {
			author := testdomain.CreateAuthor(s.system)
			book1, _ := author.CreateBook(ctx, EmptyStatus)
			book2, _ := author.CreateBook(ctx, CompletedStatus)

			books, _ := author.GetBooks(ctx)

			Convey("When the books for the author are retrieved a 2nd time", func() {
				s.system.BookCacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				books2, err := author.GetBooks(ctx)

				Convey("The books should be retrieved successfully, "+
					"And the data should've been retrieved from the cache", func() {
					So(err, ShouldBeNil)
					So(s.system.BookCacheStore.Miss(), ShouldEqual, 0)
					So(s.system.NonUniqueKeyDataFetcher.Reads(), ShouldEqual, 0)
					ContainsBooks(books, book1, book2)
					ContainsBooks(books2, book1, book2)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TestGetCreateAndGetAuthorBooks() {
	ctx := s.system.Ctx

	Convey("Scenario: Retrieve author books from cache", s.T(), func() {
		Convey("Given an author with 2 books in the system, "+
			"And I query for the books of the author once", func() {
			author := testdomain.CreateAuthor(s.system)
			book1, _ := author.CreateBook(ctx, EmptyStatus)
			book2, _ := author.CreateBook(ctx, CompletedStatus)

			books, _ := author.GetBooks(ctx)

			Convey("When I create a 3rd book and query for the books again", func() {
				book3, _ := author.CreateBook(ctx, InProgressStatus)
				s.system.BookCacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				books2, err := author.GetBooks(ctx)

				Convey("The books should be retrieved successfully, "+
					"And the data should've been retrieved from the cache", func() {
					So(err, ShouldBeNil)
					So(s.system.BookCacheStore.Hits(), ShouldEqual, 1)
					So(s.system.BookCacheStore.Miss(), ShouldEqual, 0)
					So(s.system.NonUniqueKeyDataFetcher.Reads(), ShouldEqual, 0)
					ContainsBooks(books, book1, book2)
					ContainsBooks(books2, book1, book2, book3)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TestGetCreateAndGetTwoAuthorsBooks() {
	ctx := s.system.Ctx

	Convey("Scenario: Retrieve books for 2 authors from cache with creation", s.T(), func() {
		Convey("Given an author with 2 books in the system, "+
			"And an author without books, "+
			"And I query for the books of the authors once", func() {
			author1 := testdomain.CreateAuthor(s.system)
			author2 := testdomain.CreateAuthor(s.system)
			book1, _ := author1.CreateBook(ctx, EmptyStatus)
			book2, _ := author1.CreateBook(ctx, CompletedStatus)

			authorIds := []string{author1.AuthorId, author2.AuthorId}

			books, _ := testdomain.GetBooksForAuthors(s.system, authorIds)

			Convey("When I create a 3rd book for the 1st author and query for the books again", func() {
				book3, _ := author1.CreateBook(ctx, InProgressStatus)
				s.system.BookCacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				books2, err := testdomain.GetBooksForAuthors(s.system, authorIds)

				Convey("The books should be retrieved successfully, "+
					"And the data should've been retrieved from the cache", func() {
					So(err, ShouldBeNil)
					So(s.system.BookCacheStore.Hits(), ShouldEqual, 2)
					So(s.system.BookCacheStore.Miss(), ShouldEqual, 0)
					So(s.system.NonUniqueKeyDataFetcher.Reads(), ShouldEqual, 0)
					So(len(books2), ShouldEqual, 2)
					ContainsBooks(books[0], book1, book2)
					ContainsBooks(books2[0], book1, book2, book3)
					So(books2[1], ShouldBeEmpty)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TestGetUpdateAndGetTwoAuthorsBooks() {
	ctx := s.system.Ctx

	Convey("Scenario: Retrieve books for 2 authors from cache with update", s.T(), func() {
		Convey("Given an author with 2 books in the system, "+
			"And an author without books, "+
			"And I query for the books of the authors once", func() {
			author1 := testdomain.CreateAuthor(s.system)
			author2 := testdomain.CreateAuthor(s.system)
			book1, _ := author1.CreateBook(ctx, EmptyStatus)
			book2, _ := author1.CreateBook(ctx, CompletedStatus)

			authorIds := []string{author1.AuthorId, author2.AuthorId}

			testdomain.GetBooksForAuthors(s.system, authorIds)

			Convey("When I update the 1st book for the 1st author and query for the books again", func() {
				book1.UpdateStatus(ctx, InProgressStatus)
				s.system.BookCacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				books2, err := testdomain.GetBooksForAuthors(s.system, authorIds)

				Convey("The books should be retrieved successfully, "+
					"And the data should've been retrieved from the cache", func() {
					So(err, ShouldBeNil)
					So(s.system.BookCacheStore.Hits(), ShouldEqual, 2)
					So(s.system.BookCacheStore.Miss(), ShouldEqual, 0)
					So(s.system.NonUniqueKeyDataFetcher.Reads(), ShouldEqual, 0)
					So(len(books2), ShouldEqual, 2)
					ContainsBooks(books2[0], book1, book2)
					So(books2[1], ShouldBeEmpty)
				})
			})
		})
	})
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TestGetPartialUpdateAndGetTwoAuthorsBooks() {
	ctx := s.system.Ctx

	Convey("Scenario: Retrieve books for 2 authors from cache with partial update", s.T(), func() {
		Convey("Given an author with 2 books in the system, "+
			"And an author without books, "+
			"And I query for the books of the authors once", func() {
			author1 := testdomain.CreateAuthor(s.system)
			author2 := testdomain.CreateAuthor(s.system)
			book1, _ := author1.CreateBook(ctx, EmptyStatus)
			book2, _ := author1.CreateBook(ctx, CompletedStatus)

			authorIds := []string{author1.AuthorId, author2.AuthorId}

			testdomain.GetBooksForAuthors(s.system, authorIds)

			Convey("When I do a partial update of 1st book for the 1st author and query for the books again", func() {
				book1.PartialUpdateStatus(ctx, InProgressStatus)
				s.system.BookCacheStore.ClearStats()
				s.system.NonUniqueKeyDataFetcher.ClearStats()
				books2, err := testdomain.GetBooksForAuthors(s.system, authorIds)

				Convey("The books should be retrieved successfully, "+
					"And the data should've been retrieved from the cache", func() {
					So(err, ShouldBeNil)
					So(s.system.BookCacheStore.Hits(), ShouldEqual, 2)
					So(s.system.BookCacheStore.Miss(), ShouldEqual, 0)
					So(s.system.NonUniqueKeyDataFetcher.Reads(), ShouldEqual, 0)
					So(len(books2), ShouldEqual, 2)
					ContainsBooks(books2[0], book1, book2)
					So(books2[1], ShouldBeEmpty)
				})
			})
		})
	})
}

func ContainsBooks(books []*model.Book, booksToFind ...*testdomain.Book) {
	for _, b := range booksToFind {
		So(books, ShouldContain, b.DBBook)
	}
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) SetupTest() {
	s.system = startSystemForIntegrationTests()
	prepareTestDB()
}

func (s *GormRedisIntegrationNonUniqueKeyTestSuite) TearDownTest() {
	rollbackTestDb()
}
