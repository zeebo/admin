package admin

import (
	"fmt"
	"net/url"
	"testing"
)

func TestPagination(t *testing.T) {
	p := Pagination{10, 5, nil}

	table := []struct {
		given    int
		expected []int
	}{
		{1, []int{4, 5, 6}},
		{2, []int{3, 4, 5, 6, 7}},
		{3, []int{2, 3, 4, 5, 6, 7, 8}},
		{4, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{5, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{6, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{7, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{8, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, c := range table {
		r := p.PageList(c.given)
		if len(r) != len(c.expected) {
			t.Fatalf("Expected %v. Got %v.", c.expected, r)
		}
		for i, n := range r {
			if n != c.expected[i] {
				t.Fatalf("Expected %v. Got %v.", c.expected, r)
			}
		}
	}
}

func TestPaginationPage(t *testing.T) {
	base_values := url.Values{
		"numpage":  {"20"},
		"sort__id": {"asc"},
	}
	p := Pagination{10, 5, base_values}

	for i := 1; i < 10; i++ {
		data := p.Page(i)[1:] //strip off ?
		v, err := url.ParseQuery(data)
		if err != nil {
			t.Fatal(err)
		}
		if v["numpage"][0] != "20" || v["sort__id"][0] != "asc" ||
			v["page"][0] != fmt.Sprint(i) {
			t.Fatalf("Expected %v + page.\nGot %v.", base_values, v)
		}
	}
}
