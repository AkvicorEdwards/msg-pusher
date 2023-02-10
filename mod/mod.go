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
	Handle(http.ResponseWriter, *http.Request)
	Prepare(id int64, data *MessageModel) Package
	GetTarget(id int64) Target
}
