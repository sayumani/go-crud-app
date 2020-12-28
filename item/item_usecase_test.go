package item

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sayooj/trivago/utils"

	"github.com/stretchr/testify/mock"
)

var item = Item{
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
}

var items = []Item{
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

var bookingInfo = BookAccommodation{
	ItemID:     1,
	PersonName: "SVR",
	NoOfRooms:  3,
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) AddItem(ctx context.Context, item Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRepo) DeleteItem(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepo) GetItem(ctx context.Context, id int) (Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Item), args.Error(1)
}

func (m *MockRepo) UpdateItem(ctx context.Context, item Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRepo) GetItems(ctx context.Context) ([]Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Item), args.Error(1)
}

func (m *MockRepo) BookAccommodation(ctx context.Context, bookingInfo BookAccommodation) error {
	args := m.Called(ctx, bookingInfo)
	return args.Error(0)
}

func TestAddItem(t *testing.T) {
	repo := new(MockRepo)
	repo.On("AddItem", context.Background(), item).Return(nil)
	uc := ItemsUseCase{repo}
	uc.AddItem(context.Background(), item)
	repo.AssertExpectations(t)
}

func TestAddFail(t *testing.T) {
	repo := new(MockRepo)
	repo.On("AddItem", context.Background(), item).Return(errors.New("Error"))
	uc := ItemsUseCase{repo}
	uc.AddItem(context.Background(), item)
	repo.AssertExpectations(t)
}

func TestDeleteItemSuccess(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	repo.On("DeleteItem", context.Background(), 1).Return(nil)
	uc := ItemsUseCase{repo}
	uc.DeleteItem(context.Background(), 1)
	repo.AssertExpectations(t)
}

func TestDeleteItemItemNotFound(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(Item{}, utils.ErrItemNotFound)
	// repo.On("DeleteItem", context.Background(), 1).Return(nil)
	uc := ItemsUseCase{repo}
	uc.DeleteItem(context.Background(), 1)
	repo.AssertExpectations(t)
}

func TestDeleteItemFail(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	repo.On("DeleteItem", context.Background(), 1).Return(utils.ErrItemNotDeleted)
	uc := ItemsUseCase{repo}
	uc.DeleteItem(context.Background(), 1)
	repo.AssertExpectations(t)
}

func TestGetItemSuccess(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	uc := ItemsUseCase{repo}
	res, err := uc.GetItem(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), res.ID)
	repo.AssertExpectations(t)
}

func TestGetItemFail(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(Item{}, utils.ErrItemNotFound)
	uc := ItemsUseCase{repo}
	_, err := uc.GetItem(context.Background(), 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateItemSuccess(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	repo.On("UpdateItem", context.Background(), item).Return(nil)
	uc := ItemsUseCase{repo}
	err := uc.UpdateItem(context.Background(), item)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateItemFail(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	repo.On("UpdateItem", context.Background(), item).Return(utils.ErrItemNotUpdated)
	uc := ItemsUseCase{repo}
	err := uc.UpdateItem(context.Background(), item)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetItemsSuccess(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItems", context.Background()).Return(items, nil)
	uc := ItemsUseCase{repo}
	res, err := uc.GetItems(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, res[0].ID, uint64(1))
	assert.Equal(t, res[1].ID, uint64(2))
	repo.AssertExpectations(t)
}

func TestGetItemsFail(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItems", context.Background()).Return([]Item{}, utils.ErrFetchError)
	uc := ItemsUseCase{repo}
	_, err := uc.GetItems(context.Background())
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestBookAccommodationSuccess(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	repo.On("BookAccommodation", context.Background(), bookingInfo).Return(nil)
	uc := ItemsUseCase{repo}
	err := uc.BookAccommodation(context.Background(), bookingInfo)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookAccommodationZeroRooms(t *testing.T) {
	repo := new(MockRepo)
	newitem := item
	newitem.Availability = 0
	repo.On("GetItem", context.Background(), 1).Return(newitem, nil)
	uc := ItemsUseCase{repo}
	err := uc.BookAccommodation(context.Background(), bookingInfo)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestBookAccommodationNotEnoughRooms(t *testing.T) {
	repo := new(MockRepo)
	newBooking := bookingInfo
	newBooking.NoOfRooms = 11
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	uc := ItemsUseCase{repo}
	err := uc.BookAccommodation(context.Background(), newBooking)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestBookAccommodationFail(t *testing.T) {
	repo := new(MockRepo)
	repo.On("GetItem", context.Background(), 1).Return(item, nil)
	repo.On("BookAccommodation", context.Background(), bookingInfo).Return(utils.ErrBookingFailed)
	uc := ItemsUseCase{repo}
	err := uc.BookAccommodation(context.Background(), bookingInfo)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
