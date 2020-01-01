package main

import (
	"fmt"
	"sync"

	"github.com/freahs/lunch-server/data"
)

type storable interface {
	SetID(int)
}

type indexedTable struct {
	nextID int
	items  map[int]interface{}
	*sync.Mutex
}

func newTable() *indexedTable {
	return &indexedTable{1, make(map[int]interface{}), &sync.Mutex{}}
}

func (t *indexedTable) Add(item storable) interface{} {
	t.Lock()
	defer t.Unlock()
	var id int
	id, t.nextID = t.nextID, t.nextID+1
	item.SetID(id)
	t.items[id] = item
	return item
}

func (t *indexedTable) Get(id int) (item interface{}, err error) {
	t.Lock()
	defer t.Unlock()
	if i, ok := t.items[id]; ok {
		return i, nil
	}
	return nil, fmt.Errorf("no such id %d", id)
}

func (t *indexedTable) Del(id int) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.items[id]; ok {
		delete(t.items, id)
		return nil
	}
	return fmt.Errorf("no such id %d", id)
}

type index struct {
	m map[interface{}]interface{}
}

func newIndex() *index {
	return &index{make(map[interface{}]interface{})}
}

func (c *index) Add(k interface{}, v interface{}) {
	c.m[k] = v
}

func (c *index) Get(k interface{}) (interface{}, error) {
	v, ok := c.m[k]
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("no such key in index %v", k)
}

func (c *index) Del(k interface{}) {
	if _, ok := c.m[k]; ok {
		delete(c.m, k)
	}
}

type database struct {
	restaurants       *indexedTable
	menues            *indexedTable
	menuesRestaurants map[int][]int
	*sync.Mutex
}

/* func newDB() *database {
	return &database{newTable(), newTable(), newIndex(), &sync.Mutex{}}
}

func (db *database) AddRestaurant(name string) restaurant {
	res := db.restaurants.Add(&restaurant{Name: name}).(*restaurant)
	return *res
}
*/

/*
func (db *database) AddMenu(restaurantID int, m menu) menu {
	key := fmt.Sprintf("%d%04d%02d%02d", restaurantID, m.Year, m.Month, m.Day)
	if db.menuesIndex.Get(key)

	res := db.menues.Add(&menu{Name: name}).(*menu)
	return *res
}

func (db *database) GetMenu(id int) (menu, error) {
	res, err := db.menues.Get(id)
	if err != nil {
		return menu{}, err
	}
	return *res.(*menu), nil
}
*/

func (db *database) DelMenu(id int) error {
	return db.menues.Del(id)
}

/*
GET   /restaurant (query)
GET   /restaurant/{id} (query)
POST  /restaurant/{id}/menu
GET   /restaurant/{id}/menu/{id}
PUT   /restaurant/{id}/menu/{id}
POST  /menu/{id}
*/

/* func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main() {
	initEvents()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
*/
/*
func main() {
	r := newRestaurant("hej")
	for i := 1; i < 6; i++ {
		d := date{Year: 100, Month: 10, Day: i}
		r.AddMenu(menu{ID: i + 100, date: d})
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	r2 := newRestaurant("test")
	err = json.Unmarshal(b, &r2)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range r2.Menus {
		fmt.Printf("%+v: %+v\n", k, v)
	}

} */

func main() {
	s := data.NewStore()
	s.AddMenu(data.NewMenu("a", 1, 2, 3, "item5", "item6"))
	s.AddMenu(data.NewMenu("a", 1, 2, 4, "item1", "item2", "item3"))
	s.AddMenu(data.NewMenu("a", 1, 2, 5, "item4"))
	s.Test()

	//fmt.Println(s.Filter(data.FilterLt, 1, 2, 5).Menus())

}
