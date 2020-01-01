package data

import (
	"fmt"
	"testing"
)

func TestDate(t *testing.T) {
	type TCdate struct {
		d1, d2   date
		expected []bool
	}

	tests := []TCdate{
		{date{2019, 9, 16}, date{2019, 10, 16}, []bool{true, true, false, false, false}},
		{date{2019, 10, 16}, date{2019, 10, 16}, []bool{false, true, true, true, false}},
		{date{2019, 11, 16}, date{2019, 10, 16}, []bool{false, false, false, true, true}},
		{date{2019, 10, 15}, date{2019, 10, 16}, []bool{true, true, false, false, false}},
		{date{2019, 10, 17}, date{2019, 10, 16}, []bool{false, false, false, true, true}},
		{date{2020, 1, 1}, date{2019, 10, 16}, []bool{false, false, false, true, true}},
		{date{2018, 12, 32}, date{2019, 10, 16}, []bool{true, true, false, false, false}},
	}

	for _, tc := range tests {
		fnames := [5]string{"lt", "le", "eq", "ge", "gt"}
		fn := [5]func(date) bool{tc.d1.lt, tc.d1.le, tc.d1.eq, tc.d1.ge, tc.d1.gt}
		for i := 0; i < 5; i++ {

			name := fmt.Sprintf("date.%s", fnames[i])
			t.Run(name, func(t *testing.T) {
				if res := fn[i](tc.d2); res != tc.expected[i] {
					ops := [5]string{"<", "<=", "==", ">=", ">"}
					t.Errorf("%v %s %v: got %v, expected %v", tc.d1, ops[i], tc.d2, res, tc.expected[i])
				}
			})

		}
	}
}

func TestMenu_less(t *testing.T) {
	type fields struct {
		r string
		d date
		i []string
	}

	other := &Menu{"B", date{2019, 10, 16}, []string{}}

	tests := []struct {
		fields fields
		other  *Menu
		want   bool
	}{
		{fields{"A", date{2019, 9, 16}, []string{}}, other, true},
		{fields{"A", date{2019, 10, 16}, []string{}}, other, true},
		{fields{"C", date{2019, 10, 16}, []string{}}, other, false},
		{fields{"A", date{2019, 11, 16}, []string{}}, other, false},
		{fields{"A", date{2019, 10, 15}, []string{}}, other, true},
		{fields{"A", date{2019, 10, 17}, []string{}}, other, false},
		{fields{"A", date{2020, 1, 1}, []string{}}, other, false},
		{fields{"A", date{2018, 12, 32}, []string{}}, other, true},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("Menu.less() %d	", i)
		t.Run(name, func(t *testing.T) {
			m := Menu{
				r: tt.fields.r,
				d: tt.fields.d,
				i: tt.fields.i,
			}
			if got := m.less(tt.other); got != tt.want {
				t.Errorf("Menu.less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore(t *testing.T) {

	checkNames := func(s *Store, expected ...string) {
		if len(s.menus) != len(expected) {
			t.Errorf("len(m.menus) != len(expected): %d != %d", len(s.menus), len(expected))
		}

		for i, m := range s.menus {
			if i >= len(expected) {
				t.Errorf("unexpected: s.menus[%d].Restaurant() = %v", i, s.menus[i].Restaurant())
			} else if m.Restaurant() != expected[i] {
				t.Errorf("s.menus[%d].Restaurant() != expected[%d]: (%v != %v)", i, i, s.menus[i].Restaurant(), expected[i])
			}
		}
	}

	s := NewStore()
	t.Run("Ensure menus are inserted in the right order", func(t *testing.T) {
		s.AddMenu(NewMenu("c", 2019, 10, 16))
		checkNames(s, "c")
		s.AddMenu(NewMenu("i", 2020, 10, 17))
		checkNames(s, "c", "i")
		s.AddMenu(NewMenu("d", 2019, 10, 17))
		checkNames(s, "c", "d", "i")
		s.AddMenu(NewMenu("f", 2019, 10, 17))
		checkNames(s, "c", "d", "f", "i")
		s.AddMenu(NewMenu("e", 2019, 10, 17))
		checkNames(s, "c", "d", "e", "f", "i")
		s.AddMenu(NewMenu("a", 2018, 10, 17))
		checkNames(s, "a", "c", "d", "e", "f", "i")
		s.AddMenu(NewMenu("h", 2019, 11, 17))
		checkNames(s, "a", "c", "d", "e", "f", "h", "i")
		s.AddMenu(NewMenu("b", 2019, 9, 17))
		checkNames(s, "a", "b", "c", "d", "e", "f", "h", "i")
		s.AddMenu(NewMenu("g", 2019, 10, 18))
		checkNames(s, "a", "b", "c", "d", "e", "f", "g", "h", "i")
	})

	t.Run("Ensure filter works as expected", func(t *testing.T) {
		checkNames(s.FilterDate(FilterLt, 2019, 10, 17), "a", "b", "c")
		checkNames(s.FilterDate(FilterLe, 2019, 10, 17), "a", "b", "c", "d", "e", "f")
		checkNames(s.FilterDate(FilterEq, 2019, 10, 17), "d", "e", "f")
		checkNames(s.FilterDate(FilterGe, 2019, 10, 17), "d", "e", "f", "g", "h", "i")
		checkNames(s.FilterDate(FilterGt, 2019, 10, 17), "g", "h", "i")
	})
}
