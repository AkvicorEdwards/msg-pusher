package mod

import (
	"net/http"
)

type Model interface {
	// URL is url path prefix
	URL() string
	// Name is mod name
	Name() string
	// Table is the table name used in the database
	Table() string
	// Key is the key for extra
	//     The message sender can use the extra in data to carry the information required by the specified mod.
	Key() string
	// Handle is HTTP handler, handle http requests
	Handle(http.ResponseWriter, *http.Request)
	// Prepare is used to process the general information passed in and return the processed data
	Prepare(id int64, data *MessageModel) Package
	// GetTarget returns the information of the key corresponding to the id
	GetTarget(id int64) Target
}
