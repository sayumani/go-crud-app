package utils

import "github.com/pkg/errors"

var (
	//ErrItemNotFound when item not found in db
	ErrItemNotFound = errors.New("Item not found")
	//ErrItemNotAdded when an error occred during item insertion
	ErrItemNotAdded = errors.New("Error occured while adding item to db")
	//ErrDBConnectionError when db connection creation failed
	ErrDBConnectionError = errors.New("Failed to create a connection")
	//ErrFetchError when db access failed
	ErrFetchError = errors.New("Error occured while accessing the information from db")
	//ErrItemNotDeleted when a item not deleted
	ErrItemNotDeleted = errors.New("Error occured while deleting the item")
	//ErrItemNotUpdated when a item is not updated
	ErrItemNotUpdated = errors.New("Error occured while updating the item")
	// ErrRoomsNotAvailable when Rooms not available
	ErrRoomsNotAvailable = errors.New("Rooms not available")
	// ErrRoomsNotEnough when enough Rooms not available
	ErrRoomsNotEnough = errors.New("Enough Rooms not available")
	// ErrBookingFailed when Booking failed
	ErrBookingFailed = errors.New("Booking failed")
	//ErrTransactionBeginFailed when transaction begin failed
	ErrTransactionBeginFailed = errors.New("Failed to begin transaction")
	//ErrStatementCreationFailed when statement creation failed
	ErrStatementCreationFailed = errors.New("Failed to create the statement")
)

// ErrorModel struct
type ErrorModel struct {
	Type          string          `json:"type"`
	Title         string          `json:"title"`
	InvalidParams []InvalidParams `json:"invalid-params"`
}

// InvalidParams struct
type InvalidParams struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
