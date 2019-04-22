package database

import (
	"os"
	"path"

	"github.com/protosio/protos/config"
	"github.com/protosio/protos/core"
	"github.com/protosio/protos/util"

	gobEncoding "encoding/gob"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

var gconfig = config.Get()
var log = util.GetLogger("db")

// db - package wide db reference
var db *storm.DB

// Exists checks if the database file exists on disk
func Exists() bool {
	dbpath := path.Join(gconfig.WorkDir, "protos.db")
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Open opens a a boltdb database
func Open() {

	var err error
	dbpath := path.Join(gconfig.WorkDir, "protos.db")
	log.Info("Opening database [", dbpath, "]")
	db, err = storm.Open(dbpath, storm.Codec(gob.Codec))
	if err != nil {
		log.Fatalf("Failed to open database at path %s, %s", dbpath, err.Error())
	}

}

// Close closes the boltdb database
func Close() {
	log.Info("Closing database")
	db.Close()
}

// Save writes a new value for a specific key in a bucket
func Save(data interface{}) error {
	return db.Save(data)
}

// One retrieves one record from the database based on the field name
func One(fieldName string, value interface{}, to interface{}) error {
	return db.One(fieldName, value, to)
}

// All retrieves all records for a specific type
func All(to interface{}) error {
	return db.All(to)
}

// Remove removes a record of specific type
func Remove(data interface{}) error {
	return db.DeleteStruct(data)
}

//
// DB implementation the implements the core DB interface
//

// CreateDatabase returns a database instance that implements the core DB interface
func CreateDatabase() core.DB {
	return &database{}
}

type database struct {
	s *storm.DB
}

// Open opens a a boltdb database
func (db *database) Open() {

	var err error
	dbpath := path.Join(gconfig.WorkDir, "protos.db")
	log.Info("Opening database [", dbpath, "]")
	db.s, err = storm.Open(dbpath, storm.Codec(gob.Codec))
	if err != nil {
		log.Fatalf("Failed to open database at path %s, %s", dbpath, err.Error())
	}

}

// Close closes the boltdb database
func (db *database) Close() {
	log.Info("Closing database")
	db.s.Close()
}

// Save writes a new value for a specific key in a bucket
func (db *database) Save(data interface{}) error {
	return db.s.Save(data)
}

// One retrieves one record from the database based on the field name
func (db *database) One(fieldName string, value interface{}, to interface{}) error {
	return db.s.One(fieldName, value, to)
}

// All retrieves all records for a specific type
func (db *database) All(to interface{}) error {
	return db.s.All(to)
}

// Remove removes a record of specific type
func (db *database) Remove(data interface{}) error {
	return db.s.DeleteStruct(data)
}

func (db *database) Register(structure interface{}) {
	gobEncoding.Register(structure)
}
