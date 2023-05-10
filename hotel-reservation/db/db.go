package db

const (
	DBNAME      = "hotel-reservation"
	DBURI       = "mongodb://localhost:27017"
	TEST_DBNAME = "hotel-reservation-test"
)

// func ToObjectID(id string) (primitive.ObjectID, error) {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return oid, nil
// }
