package item

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddItemSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	item := Item{
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
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO item`).WithArgs(item.Name, item.Rating, item.Category, item.Image, item.Reputation, item.Price, item.Availability).WillReturnRows(sqlmock.NewRows([]string{"item_id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO item_location`).WithArgs(item.ID, item.Location.City, item.Location.State, item.Location.Country, item.Location.ZipCode, item.Location.Address).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repo := NewItemsRepository(db)
	resp := repo.AddItem(context.Background(), item)
	assert.NoError(t, resp)
}

func TestAddItemFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	item := Item{
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
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO item`).WithArgs(item.Name, item.Rating, item.Category, item.Image, item.Reputation, item.Price, item.Availability).WillReturnError(errors.New("error"))
	mock.ExpectExec(`INSERT INTO item_location`).WithArgs(item.ID, item.Location.City, item.Location.State, item.Location.Country, item.Location.ZipCode, item.Location.Address).WillReturnError(errors.New("error"))
	mock.ExpectCommit()
	repo := NewItemsRepository(db)
	resp := repo.AddItem(context.Background(), item)
	assert.Error(t, resp)
}

func TestDeleteItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec(`DELETE FROM item`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewItemsRepository(db)
	resp := repo.DeleteItem(context.Background(), 1)
	assert.NoError(t, resp)
}

func TestDeleteItemError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec(`DELETE FROM item`).WithArgs(1).WillReturnError(errors.New("error"))
	repo := NewItemsRepository(db)
	resp := repo.DeleteItem(context.Background(), 1)
	assert.Error(t, resp)
}

func TestGetItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(`SELECT`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"item_id", "name", "rating", "category", "reputation", "reputation_badge", "price", "availability", "image", "city", "state", "country", "zip_code", "address"}).AddRow(1, "test", 5, "hotel", 600, "yellow", 1000, 10, "http://sc.com", "fdfd", "dffd", "fdfdf", 67888, "dfdfdf dfd d "))
	repo := NewItemsRepository(db)
	resp, err := repo.GetItem(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, resp.ID, uint64(1))
}

func TestGetItemError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(`SELECT`).WithArgs(1).WillReturnError(errors.New("error"))
	repo := NewItemsRepository(db)
	_, err = repo.GetItem(context.Background(), 1)
	assert.Error(t, err)
}

func TestGetItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{"item_id", "name", "rating", "category", "reputation", "reputation_badge", "price", "availability", "image", "city", "state", "country", "zip_code", "address"}).
		AddRow(1, "test", 5, "hotel", 600, "yellow", 1000, 10, "http://sc.com", "fdfd", "dffd", "fdfdf", 67888, "dfdfdf dfd d ").AddRow(2, "test", 5, "hotel", 600, "yellow", 1000, 10, "http://sc.com", "fdfd", "dffd", "fdfdf", 67888, "dfdfdf dfd d "))
	repo := NewItemsRepository(db)
	resp, err := repo.GetItems(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, resp[0].ID, uint64(1))
	assert.Equal(t, resp[1].ID, uint64(2))
}

func TestGetItemsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(`SELECT`).WillReturnError(errors.New("Error"))
	repo := NewItemsRepository(db)
	_, err = repo.GetItems(context.Background())
	assert.Error(t, err)
}

func TestUpdateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	item := Item{
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
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE`).WithArgs(item.ID, item.Name, item.Rating, item.Category, item.Image, item.Reputation, item.Price, item.Availability).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE`).WithArgs(item.ID, item.Location.City, item.Location.State, item.Location.Country, item.Location.ZipCode, item.Location.Address).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repo := NewItemsRepository(db)
	resp := repo.UpdateItem(context.Background(), item)
	assert.NoError(t, resp)
}

func TestUpdateItemError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	item := Item{
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
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE`).WithArgs(item.ID, item.Name, item.Rating, item.Category, item.Image, item.Reputation, item.Price, item.Availability).WillReturnError(errors.New("error"))
	mock.ExpectExec(`UPDATE`).WithArgs(item.ID, item.Location.City, item.Location.State, item.Location.Country, item.Location.ZipCode, item.Location.Address).WillReturnError(errors.New("error"))
	mock.ExpectCommit()
	repo := NewItemsRepository(db)
	resp := repo.UpdateItem(context.Background(), item)
	assert.Error(t, resp)
}

func TestBookAccommodation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	item := BookAccommodation{
		ItemID:     1,
		PersonName: "Svr",
		NoOfRooms:  3,
	}
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE`).WithArgs(item.ItemID, item.NoOfRooms).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT`).WithArgs(item.ItemID, item.PersonName, item.NoOfRooms).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repo := NewItemsRepository(db)
	resp := repo.BookAccommodation(context.Background(), item)
	assert.NoError(t, resp)
}

func TestBookAccommodationError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	item := BookAccommodation{
		ItemID:     1,
		PersonName: "Svr",
		NoOfRooms:  3,
	}
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE`).WithArgs(item.ItemID, item.NoOfRooms).WillReturnError(errors.New("error"))
	mock.ExpectExec(`INSERT`).WithArgs(item.ItemID, item.PersonName, item.NoOfRooms).WillReturnError(errors.New("error"))
	mock.ExpectCommit()
	repo := NewItemsRepository(db)
	resp := repo.BookAccommodation(context.Background(), item)
	assert.Error(t, resp)
}
