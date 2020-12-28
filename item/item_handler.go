package item

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sayooj/trivago/utils"
	"github.com/sirupsen/logrus"
)

//ItemsHandler handler for items
type ItemsHandler struct {
	useCase ItemsUseCaseInterface
	logger  *logrus.Logger
}

//GetItems get all items
func (h *ItemsHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	products, err := h.useCase.GetItems(r.Context())
	if errors.Is(err, utils.ErrFetchError) {
		h.logger.Info("An error occured while fetching products")
		utils.RespondWithError(w, http.StatusInternalServerError, "An error occured while fetching products")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, products)

}

//GetItem get product based on id
func (h *ItemsHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id number")
		return
	}
	product, err := h.useCase.GetItem(r.Context(), itemID)
	if errors.Is(err, utils.ErrFetchError) {
		h.logger.Info("Error occured while fetching the item")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error occured while fetching the item")
		return
	}
	if errors.Is(err, utils.ErrItemNotFound) {
		h.logger.Info("Item not found")
		utils.RespondWithError(w, http.StatusNotFound, "Item not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, product)
}

//AddItem add a item
func (h *ItemsHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		h.logger.Info("Invalid request payload")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	invalidParams := item.ValidateRequiredItem()
	if len(invalidParams) > 0 {
		h.logger.Info("Invalid request payload")
		utils.RespondWithValidationError(w, http.StatusBadRequest, invalidParams)
		return
	}
	invalidParams = item.ValidateFields()
	if len(invalidParams) > 0 {
		h.logger.Info("Invalid request payload")
		utils.RespondWithValidationError(w, http.StatusBadRequest, invalidParams)
		return
	}
	err := h.useCase.AddItem(r.Context(), item)
	if errors.Is(err, utils.ErrItemNotAdded) {
		h.logger.Info("An error occured while adding item to db")
		utils.RespondWithError(w, http.StatusInternalServerError, "An error occured while adding item to db")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, nil)
}

//UpdateItem update a item based on id
func (h *ItemsHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	id := chi.URLParam(r, "id")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		h.logger.Info("Invalid request payload")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	invalidParams := item.ValidateFields()
	if len(invalidParams) > 0 {
		h.logger.Info("Invalid request payload")
		utils.RespondWithValidationError(w, http.StatusBadRequest, invalidParams)
		return
	}
	itemID, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Info("Invalid id")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id number")
		return
	}
	item.ID = uint64(itemID)
	err = h.useCase.UpdateItem(r.Context(), item)
	if errors.Is(err, utils.ErrItemNotUpdated) {
		h.logger.Info("An error occured while updating the product")
		utils.RespondWithError(w, http.StatusInternalServerError, "An error occured while updating the product")
		return
	}
	if errors.Is(err, utils.ErrItemNotFound) {
		h.logger.Info("Item not found")
		utils.RespondWithError(w, http.StatusNotFound, "Item not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, nil)
}

//DeleteItem delete a item based on id
func (h *ItemsHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Info("Invalid id number")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id number")
		return
	}
	err = h.useCase.DeleteItem(r.Context(), itemID)
	if errors.Is(err, utils.ErrItemNotDeleted) {
		h.logger.Info("Failed to delete product")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	if errors.Is(err, utils.ErrItemNotFound) {
		h.logger.Info("Item not found")
		utils.RespondWithError(w, http.StatusNotFound, "Item not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, nil)
}

// BookAccommodation func
func (h *ItemsHandler) BookAccommodation(w http.ResponseWriter, r *http.Request) {
	var bookingInfo BookAccommodation
	id := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id number")
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&bookingInfo); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	bookingInfo.ItemID = uint64(itemID)
	err = h.useCase.BookAccommodation(r.Context(), bookingInfo)
	if err != nil {
		if errors.Is(err, utils.ErrRoomsNotEnough) {
			utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"status": http.StatusOK, "message": "rooms not available"})
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete product")
	}
	utils.RespondWithJSON(w, http.StatusOK, nil)
}

//NewItemsHandler method
func NewItemsHandler(useCase *ItemsUseCase, log *logrus.Logger) *ItemsHandler {
	return &ItemsHandler{useCase, log}
}
