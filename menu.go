package lunch

import (
	"encoding/json"
)

/*
// Menu is an interface for one days Menu for one restaurant. It can hold zero or more Menu items
type Menu interface {
	// Restaurant returns the name of the restaurant
	Restaurant() string
	// Date returns the Date which the Menu is valid
	Date() Date
	// Items should return the items on the Menu. The list can be empty but should not be nil
	Items() []string
}
 */


// Menu stores a days Menu for a restaurant
type Menu struct {
	r string
	d Date
	i []string
}

// Restaurant returns the name of the restaurant
func (m Menu) Restaurant() string {
	return m.r
}

// Time returns the time.Time for which the Menu.
func (m Menu) Date() Date {
	return m.d
}

// Items returns the items in the Menu
func (m Menu) Items() []string {
	return m.i
}

// NewMenu returns a new Menu given the name of the restaurant, the Date (as year, month and day)
// and at least one Menu item
func NewMenu(restaurant string, year, month, day int, items ...string) Menu {
	d := NewDate(year, month, day)
	i := items
	return Menu{restaurant, d, i}
}

// MarshalJSON implements the Marshaller interface
func (m Menu) MarshalJSON() ([]byte, error) {
	type M struct {
		R string   `json:"restaurant"`
		D Date     `json:"date"`
		I []string `json:"items"`
	}
	return json.Marshal(M{m.r, m.d, m.i})
}

// UnmarshalJSON implements the Unmarshaller interface
func (m *Menu) UnmarshalJSON(bts []byte) error {
	type M struct {
		R string   `json:"restaurant"`
		D Date     `json:"date"`
		I []string `json:"items"`
	}
	var tmp M
	err := json.Unmarshal(bts, &tmp)
	if err != nil {
		return err
	}
	m.r, m.d, m.i = tmp.R, tmp.D, tmp.I
	return nil
}
