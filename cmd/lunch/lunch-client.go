package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/freahs/lunch"
	"github.com/freahs/lunch/loaders"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	datePrinter       = color.New(color.FgBlue).Add(color.Bold).SprintFunc()
	restaurantPrinter = color.New(color.FgMagenta).SprintFunc()
)

type ClientError struct {
	s string
}

func (e ClientError) Error() string {
	return e.s
}

func NewClientError(s string, i ...interface{}) ClientError {
	return ClientError{fmt.Sprintf(s, i...)}
}

type Client struct {
	Store        *lunch.Store
	Loaders      []loaders.Loader
	start        lunch.Date
	days, offset int
}

func (c *Client) Update(name string) error {
	for _, l := range c.Loaders {
		if l.Name() == name {
			err := l.Load(c.Store)
			if err != nil {
				return NewClientError("error while loading %v: err", l.Name(), err)
			}
			return nil
		}
	}
	return NewClientError("invalid loader %v", name)
}

func (c *Client) UpdateAll() error {
	failed := map[string]error{}
	for _, l := range c.Loaders {
		err := l.Load(c.Store)
		if err != nil {
			failed[l.Name()] = err
		}
	}
	if len(failed) > 0 {
		return NewClientError("the following loaders failed: %v", failed)
	}
	return nil
}

func (c *Client) parse(args []string) error {
	if len(args) == 0 {
		return nil
	}
	if args[0] == "--offset" {
		if len(args) == 1 {
			return NewClientError("missing value for --offset")
		}
		offset, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		c.offset = offset
		return c.parse(args[2:])
	}
	return nil
}

func (c *Client) Parse() {
	err := c.parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) print() {
	start := c.start.Add(0, 0, c.days*c.offset)
	stop := c.start.Add(0, 0, c.days*(c.offset+1))

	var prev *lunch.Date = nil

	for _, m := range c.Store.FilterDate(lunch.FilterGe, start).FilterDate(lunch.FilterLt, stop).Menus() {
		d := m.Date()
		if prev == nil || !prev.Equal(d) {
			prev = &d
			fmt.Printf("%s:\n", datePrinter(d))
		}
		fmt.Printf("    %s\n", restaurantPrinter(m.Restaurant()))
		for _, i := range m.Items() {
			fmt.Printf("        • %s\n", i)
		}
	}
}

func NewClient(store *lunch.Store, loaders []loaders.Loader) (*Client, error) {
	y, w := lunch.Now().Week()
	client := Client{store, loaders, lunch.Week(y, w), 7, 0}
	return &client, nil
}

func getStore(dir, filename string) (*lunch.Store, error) {
	info, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("error while creating directory '%v': %v", dir, err)
		}
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("could not stat path '%v': %v", dir, err)
	}

	if !info.Mode().IsDir() {
		return nil, fmt.Errorf("'%v' is not a directory", dir)
	}
	path := filepath.Join(dir, "client_data.json")
	store, err := lunch.LoadStore(path)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func main() {
	store := lunch.NewStore()
	c, err := NewClient(store, []loaders.Loader{
		loaders.NewPrime(),
	})
	if err != nil {
		log.Fatal(err)
	}
	c.UpdateAll()
	c.Parse()
	c.print()
}
