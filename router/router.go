package router

import (
	"github.com/go-chi/chi"
	"github.com/sayooj/trivago/item"
)

//ItemsRoutes set the routes for the Item
func ItemsRoutes(h *item.ItemsHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/", h.GetItems)                    //GET /item
		r.Get("/{id}", h.GetItem)                 //GET /item/56
		r.Post("/", h.AddItem)                    //POST /item
		r.Put("/{id}", h.UpdateItem)              //PUT /item/56
		r.Delete("/{id}", h.DeleteItem)           //DELETE /item/56
		r.Post("/{id}/book", h.BookAccommodation) //POST /item/56/booking
	})
	return r
}
