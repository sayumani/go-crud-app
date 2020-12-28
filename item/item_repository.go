package item

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sayooj/trivago/utils"
)

//ItemsRepositoryInterface interface
type ItemsRepositoryInterface interface {
	AddItem(ctx context.Context, p Item) error
	DeleteItem(ctx context.Context, id int) error
	GetItem(ctx context.Context, id int) (Item, error)
	UpdateItem(ctx context.Context, item Item) error
	GetItems(ctx context.Context) ([]Item, error)
	BookAccommodation(ctx context.Context, bookingInfo BookAccommodation) error
}

//ItemsRepository struct
type ItemsRepository struct {
	db *sql.DB
}

//AddItem adds a Item to db
func (r *ItemsRepository) AddItem(ctx context.Context, item Item) error {
	tx, err := r.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction%w", utils.ErrTransactionBeginFailed)
	}
	itemQuery := `INSERT INTO item(name, rating, category, image, reputation , price , availability) VALUES($1 , $2 , $3 , $4 , $5 , $6 ,$7) RETURNING item_id`
	err = tx.QueryRowContext(ctx, itemQuery, item.Name, item.Rating, item.Category, item.Image, item.Reputation, item.Price, item.Availability).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("Error occured during insertion %w", utils.ErrItemNotAdded)
	}
	locationQry := `INSERT INTO item_location(item_id , city, state, country, zip_code, address ) VALUES($1 , $2 , $3 , $4 , $5 , $6 )`
	result, err := tx.ExecContext(ctx, locationQry, item.ID, item.Location.City, item.Location.State, item.Location.Country, item.Location.ZipCode, item.Location.Address)
	if err != nil {
		return fmt.Errorf("Error occured during insertion %w", utils.ErrItemNotAdded)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error occured during insertion %w", utils.ErrItemNotAdded)
	}
	if rows == 0 {
		return fmt.Errorf("Error occured during insertion %w", sql.ErrNoRows)
	}
	tx.Commit()
	return nil
}

//DeleteItem delete Item from db
func (r *ItemsRepository) DeleteItem(ctx context.Context, id int) error {
	query := "DELETE FROM item WHERE item_id=$1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete product %w", utils.ErrItemNotDeleted)
	}
	return nil
}

//GetItem gets a Item based on id
func (r *ItemsRepository) GetItem(ctx context.Context, id int) (Item, error) {
	var item Item
	query := `
	SELECT
		item.item_id,
		item.name,
		item.rating,
		item.category,
		item.reputation,
		CASE
			WHEN item.reputation <= 500 THEN 'red'
			WHEN item.reputation <= 799 THEN 'yellow'
        	ELSE 'green'
		END AS reputation_badge,
		item.price,
		item.availability,
		item.image,
		item_location.city,
		item_location.state,
		item_location.country,
		item_location.zip_code,
		item_location.address
	FROM
		item
	INNER JOIN
		item_location
	ON 
		item.item_id = item_location.item_id
	WHERE
		item.item_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name, &item.Rating, &item.Category, &item.Reputation, &item.ReputationBadge, &item.Price, &item.Availability, &item.Image, &item.Location.City, &item.Location.State, &item.Location.Country, &item.Location.ZipCode, &item.Location.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return Item{}, fmt.Errorf("Item not found %w", utils.ErrItemNotFound)
		}
		return Item{}, fmt.Errorf("Failed to fetch Item%w", utils.ErrFetchError)
	}
	return item, nil
}

//UpdateItem  updates a Item with sku code
func (r *ItemsRepository) UpdateItem(ctx context.Context, item Item) error {
	tx, err := r.db.Begin()
	defer func() {
		if err != nil {
			// rolling back if error occured
			tx.Rollback()
		}
	}()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction%w", utils.ErrTransactionBeginFailed)
	}

	// update item details
	itemQry := `UPDATE item SET name = $2, rating = $3, category=$4 , image =$5 , reputation =$6 , price=$7 , availability = $8 WHERE item_id = $1;`
	_, err = tx.ExecContext(ctx, itemQry, item.ID, item.Name, item.Rating, item.Category, item.Image, item.Reputation, item.Price, item.Availability)
	if err != nil {
		return fmt.Errorf("Error occured while updating the Item %w", utils.ErrItemNotUpdated)
	}

	// update location
	locationQry := `UPDATE item_location SET city = $2, state = $3, country=$4 , zip_code =$5 , address =$6  WHERE item_id = $1;`
	_, err = tx.ExecContext(ctx, locationQry, item.ID, item.Location.City, item.Location.State, item.Location.Country, item.Location.ZipCode, item.Location.Address)
	if err != nil {
		return fmt.Errorf("Error occured while updating the Item %w", utils.ErrItemNotUpdated)
	}

	tx.Commit()
	return nil
}

//GetItems with limits
func (r *ItemsRepository) GetItems(ctx context.Context) ([]Item, error) {
	query := `
	SELECT
		item.item_id,
		item.name,
		item.rating,
		item.category,
		item.reputation,
		CASE
			WHEN item.reputation <= 500 THEN 'red'
			WHEN item.reputation <= 799 THEN 'yellow'
        	ELSE 'green'
		END AS reputation_badge,
		item.price,
		item.availability,
		item.image,
		item_location.city,
		item_location.state,
		item_location.country,
		item_location.zip_code,
		item_location.address
	FROM
		item
	LEFT JOIN
		item_location
	ON 
		item.item_id = item_location.item_id;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return []Item{}, fmt.Errorf("Error occured while fetching record%w", utils.ErrFetchError)
	}
	defer rows.Close()
	items := []Item{}
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Rating, &i.Category, &i.Reputation, &i.ReputationBadge, &i.Price, &i.Availability, &i.Image, &i.Location.City, &i.Location.State, &i.Location.Country, &i.Location.ZipCode, &i.Location.Address); err != nil {
			return nil, fmt.Errorf("Error occured while fetching record%w", utils.ErrFetchError)
		}
		items = append(items, i)
	}
	return items, nil
}

// BookAccommodation book accomodation
func (r *ItemsRepository) BookAccommodation(ctx context.Context, bookingInfo BookAccommodation) error {
	tx, err := r.db.Begin()
	defer func() {
		if err != nil {
			// rolling back if error occured
			tx.Rollback()
		}
	}()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction%w", utils.ErrTransactionBeginFailed)
	}

	// update item availability
	itemQry := `UPDATE item SET availability = availability - $2 WHERE item_id = $1;`
	_, err = tx.ExecContext(ctx, itemQry, bookingInfo.ItemID, bookingInfo.NoOfRooms)
	if err != nil {
		return fmt.Errorf("Error occured while updating the Item %w", utils.ErrBookingFailed)
	}
	// creating booking record
	bookingQry := `INSERT INTO item_booking(item_id , person_name , no_of_rooms) VALUES ($1, $2, $3);`
	_, err = tx.ExecContext(ctx, bookingQry, bookingInfo.ItemID, bookingInfo.PersonName, bookingInfo.NoOfRooms)
	if err != nil {
		return fmt.Errorf("Error occured while updating the Item %w", utils.ErrBookingFailed)
	}
	tx.Commit()
	return nil
}

//NewItemsRepository method
func NewItemsRepository(db *sql.DB) *ItemsRepository {
	return &ItemsRepository{db}
}
