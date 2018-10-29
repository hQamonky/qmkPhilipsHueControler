package main

import (
	"database/sql"
	"strings"
	"testing" // The "testify/suite" package is used to make the test suite

	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	/*
		The suite is defined as a struct, with the store and db as its
		attributes. Any variables that are to be shared between tests in a
		suite should be stored as attributes of the suite instance
	*/
	store *dbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	/*
		The database connection is opened in the setup, and
		stored as an instance variable,
		as is the higher level `store`, that wraps the `db`
	*/
	connString := "host=localhost port=5432 user=postgres password=qmk dbname=qmkPhilipsHueController sslmode=disable"
	//connString := "dbname=qmkPhilipsHueController sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.store = &dbStore{db: db}
}

func (s *StoreSuite) SetupTest() {
	/*
		We delete all entries from the table before each test runs, to ensure a
		consistent state before our tests run. In more complex applications, this
		is sometimes achieved in the form of migrations
	*/
	_, err := s.db.Query("DELETE FROM tb_hue_bridges")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	// Close the connection after all tests in the suite finish
	s.db.Close()
}

// This is the actual "test" as seen by Go, which runs the tests defined below
func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateBridge() {
	// Add a bridge through the store `CreateBridge` method
	s.store.CreateBridge(&Bridge{
		ID:                "test ID",
		InternalIPAddress: "test IP Address",
		Username:          "test Username",
	})

	// Query the database for the entry we just created
	res, err := s.db.Query(`SELECT COUNT(*) FROM tb_hue_bridges WHERE id='test ID' AND internalipaddress='test IP Address' AND username=''`)
	if err != nil {
		s.T().Fatal(err)
	}

	// Get the count result
	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	// Assert that there must be one entry with the properties of the bridge that we just inserted (since the database was empty before this)
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *StoreSuite) TestGetBridge() {
	// Insert a sample bridge into the table
	_, err := s.db.Query(`INSERT INTO tb_hue_bridges (id, internalipaddress, username) VALUES('ID','10.0.0.3','Username')`)
	if err != nil {
		s.T().Fatal(err)
	}

	// Get the list of bridges through the stores `GetBridges` method
	bridges, err := s.store.GetBridges()
	if err != nil {
		s.T().Fatal(err)
	}
	// Postgresql adds additionnal space characters, so we have to get rid of those
	bridges[0].ID = strings.Replace(bridges[0].ID, " ", "", -1)
	bridges[0].InternalIPAddress = strings.Replace(bridges[0].InternalIPAddress, " ", "", -1)
	bridges[0].Username = strings.Replace(bridges[0].Username, " ", "", -1)

	// Assert that the count of bridges received must be 1
	nBridges := len(bridges)
	if nBridges != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", nBridges)
	}

	// Assert that the details of the bridge is the same as the one we inserted
	expectedBridge := Bridge{"ID", "10.0.0.3", "Username"}
	if *bridges[0] != expectedBridge {
		s.T().Errorf("incorrect details, expected %v, got %v", expectedBridge, *bridges[0])
	}
}
