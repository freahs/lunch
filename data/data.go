package data

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

type Filter int

const (
	// FilterLt represents strictly less than
	FilterLt Filter = iota
	// FilterLe represents less than or equal
	FilterLe
	// FilterEq represents equal
	FilterEq
	// FilterGe represent greater than or equal
	FilterGe
	// FilterGt represents strictly greater than
	FilterGt
)

// date in the form of Year, Month and Day
type date struct {
	Year, Month, Day int
}

func (d date) String() string {
	return fmt.Sprintf("%4d%2d%2d", d.Year, d.Month, d.Day)
}

// lt returns true if d < d2
func (d date) lt(d2 date) bool {
	if d.Year < d2.Year {
		return true
	}
	if d.Year > d2.Year {
		return false
	}
	if d.Month < d2.Month {
		return true
	}
	if d.Month > d2.Month {
		return false
	}
	return d.Day < d2.Day
}

// eq returns true if d == d2
func (d date) eq(d2 date) bool {
	return d.Year == d2.Year && d.Month == d2.Month && d.Day == d2.Day
}

// le returns true if d <= d2
func (d date) le(d2 date) bool {
	if d.Year < d2.Year {
		return true
	}
	if d.Year > d2.Year {
		return false
	}
	if d.Month < d2.Month {
		return true
	}
	if d.Month > d2.Month {
		return false
	}
	return d.Day <= d2.Day
}

// Ge returns true if d >= d2
func (d date) ge(d2 date) bool {
	return !d.lt(d2)
}

// Gt returns true if d > d2
func (d date) gt(d2 date) bool {
	return !d.le(d2)
}

// Menu stores a days menu for a restaurant
type Menu struct {
	r string
	d date
	i []string
}

// Restaurant returns the name of the restaurant
func (m Menu) Restaurant() string {
	return m.r
}

// Date returns the date for which the menu is valid
func (m Menu) Date() (year int, month int, day int) {
	return m.d.Year, m.d.Month, m.d.Day
}

// Items returns the items in the menu
func (m Menu) Items() []string {
	return m.i
}

func (m Menu) less(other *Menu) bool {
	if m.d.Year != other.d.Year {
		return m.d.Year < other.d.Year
	}
	if m.d.Month != other.d.Month {
		return m.d.Month < other.d.Month
	}
	if m.d.Day != other.d.Day {
		return m.d.Day < other.d.Day
	}
	return m.r < other.r
}

// NewMenu returns a new menu given the name of the restaurant, the date (as year, month and day)
// and at least one menu item
func NewMenu(restaurant string, year, month, day int, items ...string) Menu {
	d := date{year, month, day}
	i := items
	return Menu{restaurant, d, i}
}

// MarshalJSON implements the Marshaller interface
func (m Menu) MarshalJSON() ([]byte, error) {
	type M struct {
		R string   `json:"restaurant"`
		D date     `json:"date"`
		I []string `json:"items"`
	}
	return json.Marshal(M{m.r, m.d, m.i})
}

// UnmarshalJSON implements the Unmarshaller interface
func (m *Menu) UnmarshalJSON(bts []byte) error {
	type M struct {
		R string   `json:"restaurant"`
		D date     `json:"date"`
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

// Store holds all menus for all restaurants. All operations are thread safe.
type Store struct {
	menus []*Menu
	*sync.RWMutex
}

// MarshalJSON implements the Marshaller interface
func (s Store) MarshalJSON() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	type S struct {
		M []*Menu `json:"menus"`
	}
	return json.Marshal(S{s.menus})
}

// UnmarshalJSON implements the Unmarshaller interface
func (s *Store) UnmarshalJSON(bts []byte) error {
	s.Lock()
	defer s.Unlock()

	type S struct {
		M []Menu `json:"menus"`
	}
	var tmp S
	err := json.Unmarshal(bts, &tmp)
	if err != nil {
		return err
	}
	s.menus = make([]*Menu, 0, len(tmp.M))
	for _, m := range tmp.M {
		s.AddMenu(m)
	}
	return nil
}

// NewStore returns a new, empty, store.
func NewStore() *Store {
	return &Store{[]*Menu{}, &sync.RWMutex{}}
}

// AddMenu adds a menu to the store
func (s *Store) AddMenu(menu Menu) {
	s.Lock()
	defer s.Unlock()

	i := sort.Search(len(s.menus), func(i int) bool {
		return !s.menus[i].less(&menu)
	})
	s.menus = append(s.menus, nil)
	copy(s.menus[i+1:], s.menus[i:])
	s.menus[i] = &menu
}

// Menus returns the menus in the store. Note that values (i.e. copies) are returned, not pointers.
func (s *Store) Menus() []Menu {
	s.RLock()
	defer s.RUnlock()

	res := make([]Menu, len(s.menus))
	for i, m := range s.menus {
		res[i] = *m
	}

	return res
}

// FilterName filters
func (s *Store) FilterName(name string) *Store {
	s.RLock()
	defer s.RUnlock()

	menus := make([]*Menu, 0, len(s.menus))
	for _, m := range s.menus {
		if m.Restaurant() == name {
			menus = append(menus, m)
		}
	}
	return &Store{menus, s.RWMutex}
}

// FilterDate filters
func (s *Store) FilterDate(f Filter, year, month, day int) *Store {
	s.RLock()
	defer s.RUnlock()

	d2 := &date{year, month, day}
	start, stop := 0, len(s.menus)

	search := func(compareFunc func(date) bool) int {
		return sort.Search(len(s.menus), func(i int) bool {
			return compareFunc(s.menus[i].d)
		})
	}

	// Want function f for which f(i)=1 => f(i+1)=1. Search will return the lowest index i for
	// which f holds. Using Search to find a start and stop index...
	switch f {
	case FilterLt:
		//     the lowest index i for which d2 <= menus[i] is true
		stop = search(d2.le)
	case FilterLe:
		stop = search(d2.lt)
	case FilterEq:
		start = search(d2.le)
		stop = search(d2.lt)
	case FilterGe:
		start = search(d2.le)
	case FilterGt:
		start = search(d2.lt)
	default:
		panic(fmt.Errorf("illegal filter '%v'", f))
	}
	return &Store{s.menus[start:stop], s.RWMutex}
}

// Save saves the storage to disk
func (s *Store) Save(filename string) error {
	s.RLock()
	defer s.RUnlock()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	err = encoder.Encode(s)
	if err != nil {
		return err
	}
	return nil
}

// LoadStore loads a store saved to disk
func LoadStore(filename string) (*Store, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	s := NewStore()
	err = decoder.Decode(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
