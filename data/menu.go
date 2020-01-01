package data

import (
	"encoding/json"
)

// Menu is an interface for one days menu for one restaurant. It can hold zero or more menu items
type Menu interface {
	// Restaurant returns the name of the restaurant
	Restaurant() string
	// Date returns the Date which the menu is valid
	Date() Date
	// Items should return the items on the menu. The list can be empty but should not be nil
	Items() []string
}


// Menu stores a days menu for a restaurant
type menu struct {
	r string
	d Date
	i []string
}

// Restaurant returns the name of the restaurant
func (m menu) Restaurant() string {
	return m.r
}

// Time returns the time.Time for which the menu.
func (m menu) Date() Date {
	return m.d
}

// Items returns the items in the menu
func (m menu) Items() []string {
	return m.i
}

// NewMenu returns a new menu given the name of the restaurant, the Date (as year, month and day)
// and at least one menu item
func NewMenu(restaurant string, year, month, day int, items ...string) Menu {
	d := NewDate(year, month, day)
	i := items
	return menu{restaurant, d, i}
}

// MarshalJSON implements the Marshaller interface
func (m menu) MarshalJSON() ([]byte, error) {
	type M struct {
		R string   `json:"restaurant"`
		D Date     `json:"Date"`
		I []string `json:"items"`
	}
	return json.Marshal(M{m.r, m.d, m.i})
}

// UnmarshalJSON implements the Unmarshaller interface
func (m *menu) UnmarshalJSON(bts []byte) error {
	type M struct {
		R string   `json:"restaurant"`
		D Date     `json:"Date"`
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
