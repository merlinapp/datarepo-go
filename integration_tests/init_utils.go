package integration_tests

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"log"
	"os"
	"time"
)

type ConnectionFactory func() *gorm.DB

var TestConnectionFactory ConnectionFactory

func init() {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	connectionUrl := dbParams(dbUsername, dbPassword, dbHost, dbPort, dbName)
	log.Println("Connecting to DB: ", dbHost+":"+dbPort)

	noTxDbParam := os.Getenv("NO_TXDB")
	if noTxDbParam == "true" {
		TestConnectionFactory = func() *gorm.DB {
			return ConnectLocalDb(connectionUrl)
		}
	} else {
		log.Println("Using TXDB - performing initial DB Automigrate using regular DB connection")
		db := ConnectLocalDb(connectionUrl)
		autoMigrate(db)
		db.Close()
		log.Println("DB Automigrate completed and connection closed")

		txdb.Register("txdb", "mysql", connectionUrl)
		TestConnectionFactory = ConnectTxDb
	}
}

func ConnectTxDb() *gorm.DB {
	sqlDb, err := sql.Open("txdb", fmt.Sprintf("connection_%d", time.Now().UnixNano()))
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}

	db, err := gorm.Open("mysql", sqlDb)

	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	return db
}

func ConnectLocalDb(connectionUrl string) *gorm.DB {
	db, err := gorm.Open("mysql", connectionUrl)

	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	autoMigrate(db)
	return db
}

func dbParams(params ...string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", params[0], params[1], params[2], params[3], params[4])
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.Author{})
	db.AutoMigrate(&model.Book{})
	db.AutoMigrate(&model.BookType{})
	db.AutoMigrate(&model.BookCategory{})
}
