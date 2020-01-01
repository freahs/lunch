package data

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

// Store holds all menus for all restaurants. All operations are thread safe.
type Store struct {
	menus []Menu
	*sync.RWMutex
}

// MarshalJSON implements the Marshaller interface
func (s Store) MarshalJSON() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	type S struct {
		M []Menu `json:"menus"`
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
	s.menus = make([]Menu, 0, len(tmp.M))
	for _, m := range tmp.M {
		s.AddMenu(m)
	}
	return nil
}

// NewStore returns a new, empty, store.go.
func NewStore() *Store {
	return &Store{[]Menu{}, &sync.RWMutex{}}
}

// AddMenu adds a menu to the store.go
func (s *Store) AddMenu(menu Menu) {
	s.Lock()
	defer s.Unlock()

	i := sort.Search(len(s.menus), func(i int) bool {
		d1, d2 := s.menus[i].Date(), menu.Date()
		if d2.Equal(d1) {
			return menu.Restaurant() < s.menus[i].Restaurant()
		}
		return d2.Before(d1)
	})
	s.menus = append(s.menus, nil)
	copy(s.menus[i+1:], s.menus[i:])
	s.menus[i] = menu
}

// Menus returns the menus in the store.go. Note that values (i.e. copies) are returned, not pointers.
func (s *Store) Menus() []Menu {
	s.RLock()
	defer s.RUnlock()

	res := make([]Menu, len(s.menus))
	for i, m := range s.menus {
		res[i] = m
	}

	return res
}

// FilterName filters
func (s *Store) FilterName(name string) *Store {
	s.RLock()
	defer s.RUnlock()

	menus := make([]Menu, 0, len(s.menus))
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

	d2 := NewDate(year, month, day)
	start, stop := 0, len(s.menus)

	cmpLT := func(i int) bool {
		d1 := s.menus[i].Date()
		return d2.Before(d1)
	}
	cmpLE := func(i int) bool {
		d1 := s.menus[i].Date()
		return d2.Equal(d1) || d2.Before(d1)
	}

	// Want function f for which f(i)=1 => f(i+1)=1. Search will return the lowest index i for
	// which f holds. Using Search to find a start and stop index...
	switch f {
	case FilterLt:
		//     the lowest index i for which d2 <= menus[i] is true
		stop = sort.Search(len(s.menus), cmpLE)
	case FilterLe:
		stop = sort.Search(len(s.menus), cmpLT)
	case FilterEq:
		start = sort.Search(len(s.menus), cmpLE)
		stop = sort.Search(len(s.menus), cmpLT)
	case FilterGe:
		start = sort.Search(len(s.menus), cmpLE)
	case FilterGt:
		start = sort.Search(len(s.menus), cmpLT)
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

// LoadStore loads a store.go saved to disk
func LoadStore(filename string) (*Store, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	store := NewStore()
	err = decoder.Decode(store)
	if err != nil {
		return nil, err
	}
	return store, nil
}
