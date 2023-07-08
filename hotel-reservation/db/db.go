package db

const MongoDBNameEnvName = "MONGO_DB_NAME"

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type PaginationFilter struct {
	Limit int64
	Page  int64
}
