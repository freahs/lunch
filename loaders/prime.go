package loaders

import (
	"fmt"
	"github.com/freahs/lunch"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Prime implements the Loader interface and loads menus from an URL
type Prime struct {
	URL string
}

// NewPrime returns a new Prime
func NewPrime() Prime {
	return Prime{"https://www.primeburger.se/lunchmeny/"}
}

// Name implements Loader
func (p Prime) Name() string {
	return "prime"
}

func (p Prime) parseWeekString(s string) (int, error) {
	parts := strings.Split(s, " ")
	if len(parts) < 2 {
		return -1, fmt.Errorf("can't get week number from %s", s)
	}
	week, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, fmt.Errorf("illegal week number: '%s'", s)
	}
	return week, nil
}

// Load implements Loader
func (p Prime) Load(store *lunch.Store) error {
	res, err := http.Get(p.URL)
	if err != nil {
		return err
	} else if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	now := time.Now()

	defer func() {
		switch r := recover().(type) {
		case nil: // do nothing
		case error:
			err = r
		default:
			err = fmt.Errorf("%v", r)
		}
	}()

	doc.Find("h2").Each(
		func(i int, weekSelector *goquery.Selection) {
			week, err := p.parseWeekString(weekSelector.Text())
			if err != nil {
				panic(err)
			}
			year, currentWeek := now.ISOWeek()
			if week < currentWeek-20 {
				year++
			} else if week > currentWeek+20 {
				year--
			}
			menuSelector := weekSelector.Parent().Find("tr")
			if menuSelector.Length() == 5 {
				menuSelector.Each(func(i int, itemSelector *goquery.Selection) {
					day, err := dayFromString(itemSelector.Find(".column-1").Text())
					if err != nil {
						panic(err)
					}
					item := itemSelector.Find(".column-2").Text()
					y, m, d := parseISOWeek(year, week, day).Date()
					store.AddMenu(lunch.NewMenu("prime", y, int(m), d, item))
				})
			}
		})
	return err
}
