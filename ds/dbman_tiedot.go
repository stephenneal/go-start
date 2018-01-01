/* DB manager */
package ds

import (
	"encoding/json"
	"fmt"
	"os"

    "github.com/op/go-logging"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/mitchellh/mapstructure"
)

type (
	DbManager struct {
		location string
		tdConn   *db.DB // Operate on this database
	}
	Col struct {
		*db.Col
	}
)

var (
	mgr *DbManager
	log *logging.Logger
)

func init() {
    log = logging.MustGetLogger("dbman_tiedot")
}

func ConnectDb(dbName string) *DbManager {
	dbDir := dbName
	var err error
	db, err := db.OpenDB(dbDir)
	if err != nil {
		panic(err)
	}
	mgr = &DbManager{
		location: dbDir,
		tdConn:   db,
	}
	return mgr
}

// Create a collection
func (m *DbManager) CreateCol(name string) *Col {
	if err := m.tdConn.Create(name); err != nil {
		panic(err)
	}
	return &Col{m.tdConn.Use(name)}
}

// Drop (delete) a collection
func (m *DbManager) DropCol(name string) {
	if err := m.tdConn.Drop(name); err != nil {
		panic(err)
	}
}

// Scrub (repair and compact) a collection
func (m *DbManager) ScrubCol(name string) {
	if err := m.tdConn.Scrub(name); err != nil {
		panic(err)
	}
}

// Get a collection
func UseCol(name string) *Col {
	return &Col{mgr.tdConn.Use(name)}
}

func ExistsDoc(col *Col, criteria string) (bool, error) {
	var query interface{}
	json.Unmarshal([]byte(criteria), &query)

	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

	if err := db.EvalQuery(query, col.Col, &queryResult); err != nil {
		return false, err
	}
	return len(queryResult) > 0, nil
}

func FindDoc(col *Col, criteria string, result interface{}) error {
	var query interface{}
	json.Unmarshal([]byte(criteria), &query)

	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

	if err := db.EvalQuery(query, col.Col, &queryResult); err != nil {
		panic(err)
	}

	count := len(queryResult)
	if count > 1 {
		log.Panicf("found more than one document:%d", count)
	}
	for id := range queryResult {
		// To get query result document, simply read it
		readBack, err := col.Read(id)
		if err != nil {
			return err
		}
		if err := mapstructure.Decode(readBack, result); err != nil {
			return err
		}
		return nil
	}
	// TODO return error (no records found)...?
	return nil
}

func InsertDoc(col *Col, criteria string, doc map[string]interface{}) (int, error) {
	// Only insert if document does not already exist
	if exists, err := ExistsDoc(col, criteria); err != nil {
		return 0, err
	} else if exists {
		return 0, fmt.Errorf("document already exists:%s", criteria)
	}
	return col.Insert(doc)
}

// Drop all collections
func (m *DbManager) CleanDb() {
	if _, err := os.Stat(m.location); err == nil {
		m.tdConn.Drop(MODEL0)
	}
}

// Gracefully close database
func (m *DbManager) CloseDb() {
	if err := m.tdConn.Close(); err != nil {
		panic(err)
	}
}

func (m *DbManager) PrintIndexes(col *db.Col) {
	// What collections do I have?
	for _, name := range col.AllIndexes() {
		log.Infof("Index: %s\n", name)
	}
}

func (m *DbManager) PrintCols() {
	// What collections do I have?
    log.Info("List all collections...")
	for _, name := range m.tdConn.AllCols() {
		log.Infof("Collection: %s\n", name)
	}
    log.Info("Finished listing all collections")
}
