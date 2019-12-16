package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/drreflect"
)

// Creates a new Builder for a GORM based cached repository that will handle data
// of the provided data type.
//
// The data type is expected to be a struct or a pointer to a struct
func CachedRepositoryBuilder(db *gorm.DB, dataType interface{}) datarepo.Builder {
	if db == nil {
		panic("The gorm DB instance must not be nil")
	}
	builder := datarepo.CachedRepositoryBuilder(dataType).
		WithUniqueKeyDataFetcher(NewUniqueKeyDataFetcher(db, dataType)).
		WithNonUniqueKeyDataFetcher(NewNonUniqueKeyDataFetcher(db, dataType)).
		WithDataWriter(NewDataWriter(db, dataType))
	return builder
}

func NewDataWriter(db *gorm.DB, dataType interface{}) datarepo.DataWriter {
	th := drreflect.NewReflectStructTypeHandlerFromValue(dataType)
	return &dataWriter{
		db:          db,
		typeHandler: th,
	}
}

func NewUniqueKeyDataFetcher(db *gorm.DB, dataType interface{}) datarepo.DataFetcher {
	fieldToColumnMap := getFieldToColumnNames(db, dataType)
	th := drreflect.NewReflectStructTypeHandlerFromValue(dataType)
	return &uniqueDataFetcher{
		db:                db,
		typeHandler:       th,
		fieldToColumnName: fieldToColumnMap,
	}
}

func NewNonUniqueKeyDataFetcher(db *gorm.DB, dataType interface{}) datarepo.DataFetcher {
	fieldToColumnMap := getFieldToColumnNames(db, dataType)
	th := drreflect.NewReflectStructTypeHandlerFromValue(dataType)
	return &nonUniqueDataFetcher{
		db:                db,
		typeHandler:       th,
		fieldToColumnName: fieldToColumnMap,
	}
}

func getFieldToColumnNames(db *gorm.DB, model interface{}) map[string]string {
	fieldToColumn := make(map[string]string)
	modelStruct := db.NewScope(model).GetModelStruct()
	for _, v := range modelStruct.PrimaryFields {
		fieldToColumn[v.Name] = v.DBName
	}
	for _, v := range modelStruct.StructFields {
		fieldToColumn[v.Name] = v.DBName
	}
	return fieldToColumn
}
