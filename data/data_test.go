package data

import (
	"fmt"
	"testing"
)


func TestDate(t *testing.T) {
	type TCdate struct {
		d1, d2   Date
		expected []bool
	}

	tests := []TCdate{
		{Date{2019, 9, 16}, Date{2019, 10, 16}, []bool{true, false, false}},
		{Date{2019, 10, 16}, Date{2019, 10, 16}, []bool{false, true, false}},
		{Date{2019, 11, 16}, Date{2019, 10, 16}, []bool{false, false, true}},
		{Date{2019, 10, 15}, Date{2019, 10, 16}, []bool{true, false, false}},
		{Date{2019, 10, 17}, Date{2019, 10, 16}, []bool{false, false, true}},
		{Date{2020, 1, 1}, Date{2019, 10, 16}, []bool{false, false, true}},
		{Date{2018, 12, 32}, Date{2019, 10, 16}, []bool{true, false, false}},
	}

	for _, tc := range tests {
		fnames := [3]string{"Before", "Equal", "After"}
		fn := [3]func(Date) bool{tc.d1.Before, tc.d1.Equal, tc.d1.After}
		for i := 0; i < 3; i++ {

			name := fmt.Sprintf("Date.%s", fnames[i])
			t.Run(name, func(t *testing.T) {
				if res := fn[i](tc.d2); res != tc.expected[i] {
					ops := [3]string{"<", "==", ">"}
					t.Errorf("%v %s %v: got %v, expected %v", tc.d1, ops[i], tc.d2, res, tc.expected[i])
				}
			})

		}
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
