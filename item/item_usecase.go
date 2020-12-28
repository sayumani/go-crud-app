package item

import (
	"context"
	"fmt"

	"github.com/sayooj/trivago/utils"
)

//ItemsUseCaseInterface interface
type ItemsUseCaseInterface interface {
	AddItem(ctx context.Context, Item Item) error
	DeleteItem(ctx context.Context, id int) error
	GetItem(ctx context.Context, id int) (Item, error)
	UpdateItem(ctx context.Context, item Item) error
	GetItems(ctx context.Context) ([]Item, error)
	BookAccommodation(ctx context.Context, bookingInfo BookAccommodation) error
}

//ItemsUseCase struct
type ItemsUseCase struct {
	itemRepo ItemsRepositoryInterface
}

//AddItem method
func (u *ItemsUseCase) AddItem(ctx context.Context, item Item) error {
	err := u.itemRepo.AddItem(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

//DeleteItem delete Item
func (u *ItemsUseCase) DeleteItem(ctx context.Context, id int) error {
	_, err := u.itemRepo.GetItem(ctx, id)
	if err != nil {
		return fmt.Errorf("Item not found %w", utils.ErrItemNotFound)
	}
	err = u.itemRepo.DeleteItem(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

//GetItem gets a Item with id
func (u *ItemsUseCase) GetItem(ctx context.Context, id int) (Item, error) {
	item, err := u.itemRepo.GetItem(ctx, id)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

//UpdateItem updates a Item with id
func (u *ItemsUseCase) UpdateItem(ctx context.Context, item Item) error {
	itemInfo, err := u.itemRepo.GetItem(ctx, int(item.ID))
	if err != nil {
		return fmt.Errorf("Item not found %w", utils.ErrItemNotFound)
	}
	if item.Name != "" {
		itemInfo.Name = item.Name
	}
	if item.Rating != 0 {
		itemInfo.Rating = item.Rating
	}
	if item.Category != "" {
		itemInfo.Category = item.Category
	}
	if item.Image != "" {
		itemInfo.Image = item.Image
	}
	if item.Reputation != 0 {
		itemInfo.Reputation = item.Reputation
	}
	if item.Price != 0 {
		itemInfo.Price = item.Price
	}
	if item.Availability != 0 {
		itemInfo.Availability = item.Availability
	}
	if item.Location.City != "" {
		itemInfo.Location.City = item.Location.City
	}
	if item.Location.State != "" {
		itemInfo.Location.State = item.Location.State
	}
	if item.Location.Country != "" {
		itemInfo.Location.Country = item.Location.Country
	}
	if item.Location.ZipCode != 0 {
		itemInfo.Location.ZipCode = item.Location.ZipCode
	}
	if item.Location.Address != "" {
		itemInfo.Location.Address = item.Location.Address
	}
	err = u.itemRepo.UpdateItem(ctx, itemInfo)
	if err != nil {
		return err
	}
	return nil
}

//GetItems returns items
func (u *ItemsUseCase) GetItems(ctx context.Context) ([]Item, error) {
	items, err := u.itemRepo.GetItems(ctx)
	if err != nil {
		return []Item{}, err
	}
	return items, nil
}

// BookAccommodation book accommodation
func (u *ItemsUseCase) BookAccommodation(ctx context.Context, bookingInfo BookAccommodation) error {
	itemInfo, err := u.itemRepo.GetItem(ctx, int(bookingInfo.ItemID))
	if err != nil {
		return fmt.Errorf("Item not found %w", utils.ErrItemNotFound)
	}
	if itemInfo.Availability == 0 {
		return fmt.Errorf("Item not found %w", utils.ErrRoomsNotEnough)
	}
	if bookingInfo.NoOfRooms > itemInfo.Availability {
		return fmt.Errorf("Item not found %w", utils.ErrRoomsNotEnough)
	}
	err = u.itemRepo.BookAccommodation(ctx, bookingInfo)
	if err != nil {
		return err
	}
	return nil
}

//NewItemsUseCase method
func NewItemsUseCase(repo *ItemsRepository) *ItemsUseCase {
	return &ItemsUseCase{repo}
}
