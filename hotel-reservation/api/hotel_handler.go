package api

import (
	"github.com/fulltimegodev/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrNotFound("rooms")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return ErrNotFound("hotel")
	}
	return c.JSON(hotel)
}

type ResourceResp struct {
	Results int   `json:"results"`
	Data    any   `json:"data"`
	Page    int64 `json:"page"`
}

type HotelQueryParams struct {
	db.PaginationFilter

	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{
		"rating": params.Rating,
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(),
		primitive.M(filter), &params.PaginationFilter)
	if err != nil {
		return ErrNotFound("hotels")
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    params.Page,
	}
	return c.JSON(resp)
}
