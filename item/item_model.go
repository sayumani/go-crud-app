package item

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sayooj/trivago/utils"
)

// Item struct
type Item struct {
	ID              uint64   `json:"id"`
	Name            string   `json:"name"`
	Rating          uint     `json:"rating"`
	Category        string   `json:"category"`
	Location        Location `json:"location"`
	Image           string   `json:"image"`
	Reputation      uint64   `json:"reputation"`
	ReputationBadge string   `json:"reputationBadge"`
	Price           uint64   `json:"price"`
	Availability    uint     `json:"availability"`
}

// Location struct
type Location struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	ZipCode uint64 `json:"zip_code"`
	Address string `json:"address"`
}

// BookAccommodation struct
type BookAccommodation struct {
	ItemID     uint64 `json:"item_id"`
	PersonName string `json:"person_name"`
	NoOfRooms  uint   `json:"no_of_rooms"`
}

// ValidateRequiredItem validates the item
func (i Item) ValidateRequiredItem() []utils.InvalidParams {
	validationErr := []utils.InvalidParams{}
	invalidItem := utils.InvalidParams{}
	// required validation
	if i.Name == "" {
		invalidItem.Name = "name"
		invalidItem.Reason = "name required"
		validationErr = append(validationErr, invalidItem)
	}
	if i.Rating == 0 {
		invalidItem.Name = "rating"
		invalidItem.Reason = "rating required"
		validationErr = append(validationErr, invalidItem)
	}
	if i.Category == "" {
		invalidItem.Name = "category"
		invalidItem.Reason = "category required"
		validationErr = append(validationErr, invalidItem)
	}
	if i.Image == "" {
		invalidItem.Name = "image"
		invalidItem.Reason = "image required"
		validationErr = append(validationErr, invalidItem)
	}
	if i.Reputation == 0 {
		invalidItem.Name = "reputation"
		invalidItem.Reason = "reputation required"
		validationErr = append(validationErr, invalidItem)
	}
	if i.Price == 0 {
		invalidItem.Name = "price"
		invalidItem.Reason = "price required"
		validationErr = append(validationErr, invalidItem)
	}
	if i.Availability == 0 {
		invalidItem.Name = "availability"
		invalidItem.Reason = "availability required"
		validationErr = append(validationErr, invalidItem)
	}
	return validationErr
}

// ValidateFields validate fields
func (i Item) ValidateFields() []utils.InvalidParams {
	validationErr := []utils.InvalidParams{}
	invalidItem := utils.InvalidParams{}

	if i.Name != "" {
		var isInvalidName = false
		if len(i.Name) < 10 {
			invalidItem.Name = "name"
			invalidItem.Reason = "Name should be 10 char long"
			validationErr = append(validationErr, invalidItem)
		}
		for _, key := range invalidHotelNames {
			if strings.Contains(i.Name, key) {
				isInvalidName = true
				break
			}
		}
		if isInvalidName {
			invalidItem.Reason = "Name should not contain [Free, Offer, Book, Website]"
			validationErr = append(validationErr, invalidItem)
		}
	}

	if i.Rating < 0 || i.Rating > 5 {
		invalidItem.Name = "rating"
		invalidItem.Reason = "Rating should be >=0 and <=5"
		validationErr = append(validationErr, invalidItem)
	}

	if i.Category != "" {
		validCategory := false
		for _, key := range invalidCategory {
			if i.Category == key {
				validCategory = true
				break
			}
		}
		if !validCategory {
			invalidItem.Name = "category"
			invalidItem.Reason = `category should any of [hotel, alternative, hostel, lodge, resort, guest-house]`
			validationErr = append(validationErr, invalidItem)
		}
	}

	if i.Image != "" {
		u, err := url.Parse(i.Image)
		valid := (err == nil) && u.Scheme != "" && u.Host != ""
		if !valid {
			invalidItem.Name = "Image"
			invalidItem.Reason = `image should be a valid url`
			validationErr = append(validationErr, invalidItem)
		}
	}

	if i.Reputation < 0 || i.Reputation > 1000 {
		invalidItem.Name = "reputation"
		invalidItem.Reason = `reputation should be >= 0 and <= 1000`
		validationErr = append(validationErr, invalidItem)
	}

	if i.Location.ZipCode != 0 {
		// re := regexp.MustCompile(`\[0-9][0-9][0-9][0-9][0-9]`)
		strZip := fmt.Sprint(i.Location.ZipCode)
		// valid := re.MatchString(strZip)
		valid := len(strZip) == 5
		if !valid {
			invalidItem.Name = "zip code"
			invalidItem.Reason = `zip code  should be 5 digits`
			validationErr = append(validationErr, invalidItem)
		}
	}
	return validationErr
}
