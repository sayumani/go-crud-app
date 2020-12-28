package item

import (
	"testing"
)

func TestValidateRequiredItem(t *testing.T) {
	item := Item{}
	invalidFields := item.ValidateRequiredItem()
	if len(invalidFields) != 7 {
		t.Errorf("Expected 7 got %d", len(invalidFields))
	}
	if invalidFields[0].Name != "name" {
		t.Errorf("Expected name got %s", invalidFields[0].Name)
	}
	if invalidFields[0].Reason != "name required" {
		t.Errorf("Expected name required got %s", invalidFields[0].Reason)
	}
}

func TestValidateFields(t *testing.T) {
	item := Item{
		Name:         "tets",
		Rating:       7,
		Category:     "abc",
		Image:        "google.com",
		Reputation:   1100,
		Price:        1000,
		Availability: 10,
	}
	validateErr := item.ValidateFields()
	if validateErr[0].Name != "name" {
		t.Errorf("Expected name got %s", validateErr[0].Name)
	}
	if validateErr[1].Name != "rating" {
		t.Errorf("Expected ratingname got %s", validateErr[1].Name)
	}
	if validateErr[2].Name != "category" {
		t.Errorf("Expected category got %s", validateErr[2].Name)
	}
	if validateErr[3].Name != "Image" {
		t.Errorf("Expected Image got %s", validateErr[3].Name)
	}
	if validateErr[4].Name != "reputation" {
		t.Errorf("Expected reputation got %s", validateErr[4].Name)
	}

}
