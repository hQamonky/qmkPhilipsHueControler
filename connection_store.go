package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
	"fmt"
)

// Store will have two methods, to add a new hue bridge,
// and to get all existing bridges
// Each method returns an error, in case something goes wrong
// TODO : Add other methods for deleting a bridge, and
// getting a bridge by id or by username
type Store interface {
	CreateBridge(bridge *Bridge) error
	GetBridges() ([]*Bridge, error)
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateBridge(bridge *Bridge) error {
	// The first underscore means that we don't care about what's returned from
	// this insert query. We just want to know if it was inserted correctly,
	// and the error will be populated if it wasn't
	fmt.Println("Using CreateBridge :", bridge)
	_, err := store.db.Query("INSERT INTO tb_hue_bridges(id, internalipaddress, username) VALUES ($1,$2,$3)", bridge.ID, bridge.InternalIPAddress, bridge.Username)
	return err
}

func (store *dbStore) GetBridges() ([]*Bridge, error) {
	// Query the database for all bridges, and return the result to the
	// `rows` object
	rows, err := store.db.Query("SELECT id, internalipaddress, username from tb_hue_bridges")
	// We return incase of an error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of bridges
	bridges := []*Bridge{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a bridge,
		bridge := &Bridge{}
		// Populate the attributes of the bridge,
		// and return incase of an error
		if err := rows.Scan(&bridge.ID, &bridge.InternalIPAddress, &bridge.Username); err != nil {
			return nil, err
		}
		// Finally, append the result to the returned array, and repeat for
		// the next row
		bridges = append(bridges, bridge)
	}
	return bridges, nil
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

// InitStore initializes our store
/*
We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(s Store) {
	store = s
}
