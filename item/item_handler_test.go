package item

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/sayooj/trivago/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var itemInfo = Item{
	ID:           1,
	Name:         "gggggssssss",
	Rating:       5,
	Category:     "hotel",
	Image:        "http://google.com/img11",
	Reputation:   800,
	Price:        5000,
	Availability: 10,
	Location: Location{
		City:    "tsr",
		State:   "kerala",
		Country: "india",
		ZipCode: 56766,
		Address: "sdsfsf  df fdf  df ",
	},
}

var itemsList = []Item{
	Item{
		ID:           1,
		Name:         "hotel abcd",
		Rating:       5,
		Category:     "hotel",
		Image:        "http://abc.com/img.jpg",
		Reputation:   800,
		Price:        1000,
		Availability: 10,
		Location: Location{
			City:    "abcs",
			State:   "sfddf",
			Country: "india",
			ZipCode: 78999,
			Address: "sdsfsf  df fdf  df ",
		},
	},
	Item{
		ID:           2,
		Name:         "hotel abcd",
		Rating:       5,
		Category:     "hotel",
		Image:        "http://abc.com/img.jpg",
		Reputation:   800,
		Price:        1000,
		Availability: 10,
		Location: Location{
			City:    "abcs",
			State:   "sfddf",
			Country: "india",
			ZipCode: 78999,
			Address: "sdsfsf  df fdf  df ",
		},
	},
}

var bookingInfos = BookAccommodation{
	ItemID:     1,
	PersonName: "SVR",
	NoOfRooms:  3,
}

type ctxKey struct {
	name string
}

type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) AddItem(ctx context.Context, item Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockUseCase) DeleteItem(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUseCase) GetItem(ctx context.Context, id int) (Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Item), args.Error(1)
}

func (m *MockUseCase) UpdateItem(ctx context.Context, item Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockUseCase) GetItems(ctx context.Context) ([]Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Item), args.Error(1)
}

func (m *MockUseCase) BookAccommodation(ctx context.Context, bookingInfo BookAccommodation) error {
	args := m.Called(ctx, bookingInfo)
	return args.Error(0)
}

func TestGetItemsHandler(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	uc.On("GetItems", context.Background()).Return(itemsList, nil)
	req, err := http.NewRequest("GET", "/item", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.GetItems)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.NoError(t, err)
	var item []Item
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&item)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, uint64(1), item[0].ID)
	uc.AssertExpectations(t)
}

func TestGetItemsHandlerError(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	uc.On("GetItems", context.Background()).Return([]Item{}, utils.ErrFetchError)
	req, err := http.NewRequest("GET", "/item", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.GetItems)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusInternalServerError, status)
	uc.AssertExpectations(t)
}

func ItemsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleParam := chi.URLParam(r, "id")
		itemID, err := strconv.Atoi(articleParam)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey{"id"}, itemID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestGetItemHandler(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	req, err := http.NewRequest("GET", "/item/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	uc.On("GetItem", req.Context(), 1).Return(itemInfo, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.GetItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.NoError(t, err)
	var item Item
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&item)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, uint64(1), item.ID)
	uc.AssertExpectations(t)
}

func TestGetItemHandlerBadRequest(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	req, _ := http.NewRequest("GET", "/item/bad", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "bad")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.GetItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestGetItemHandlerInternalError(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	req, _ := http.NewRequest("GET", "/item/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	uc.On("GetItem", req.Context(), 1).Return(Item{}, utils.ErrItemNotFound)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.GetItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusNotFound, status)
	uc.AssertExpectations(t)
}

func TestDeleteItemHandler(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	req, _ := http.NewRequest("DELETE", "/item/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	uc.On("DeleteItem", req.Context(), 1).Return(nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.DeleteItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusOK, status)
	uc.AssertExpectations(t)
}

func TestDeleteItemHandlerBadRequest(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	req, _ := http.NewRequest("DELETE", "/item/bad", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "bad")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.DeleteItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestDeleteItemHandlerInternalError(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	req, _ := http.NewRequest("DELETE", "/item/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	uc.On("DeleteItem", req.Context(), 1).Return(utils.ErrItemNotDeleted)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.DeleteItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusInternalServerError, status)
	uc.AssertExpectations(t)
}

func TestAddItemHandler(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	uc.On("AddItem", context.Background(), itemInfo).Return(nil)
	body, _ := os.Open("valid_mock.json")
	req, _ := http.NewRequest("POST", "/item", body)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.AddItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusCreated, status)
	uc.AssertExpectations(t)
}

func TestAddItemHandlerBadReq(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	body, _ := os.Open("invalid_mock.json")
	req, _ := http.NewRequest("POST", "/item", body)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.AddItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestAddItemHandlerInternalError(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	uc.On("AddItem", context.Background(), itemInfo).Return(utils.ErrItemNotAdded)
	body, _ := os.Open("valid_mock.json")
	req, _ := http.NewRequest("POST", "/item", body)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.AddItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestUpdateItemandler(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	body, _ := os.Open("valid_mock.json")
	req, _ := http.NewRequest("PUT", "/item/1", body)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	uc.On("UpdateItem", req.Context(), itemInfo).Return(nil)
	handler := http.HandlerFunc(ih.UpdateItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusOK, status)
	uc.AssertExpectations(t)
}

func TestUpdateItemandlerBadRequest(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	body, _ := os.Open("invalid_mock.json")
	req, _ := http.NewRequest("PUT", "/item/1", body)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ih.UpdateItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestUpdateItemandlerInernalError(t *testing.T) {
	log := logrus.New()
	uc := new(MockUseCase)
	ih := ItemsHandler{uc, log}
	body, _ := os.Open("valid_mock.json")
	req, _ := http.NewRequest("PUT", "/item/1", body)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	uc.On("UpdateItem", req.Context(), itemInfo).Return(utils.ErrItemNotUpdated)
	handler := http.HandlerFunc(ih.UpdateItem)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusInternalServerError, status)
	uc.AssertExpectations(t)
}
