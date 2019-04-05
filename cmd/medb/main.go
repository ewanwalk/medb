package medb

import "encoder-backend/pkg/database"

var (
	Version string
	Build   string
)

func main() {

	database.Connect()
	// TODO encoder again..

}
