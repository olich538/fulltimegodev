package main

import (
	"fmt"

	"github.com/olich538/fulltimegodev/hotel-reservation/types"
)

func main() {
	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "Norway",
	}

	room := types.Room{
		Type:      types.SeaSideRoomType,
		BasePrice: 88.6,
	}
	fmt.Println("seeding the DB")
}
