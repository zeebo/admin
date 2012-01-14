package admin

import "testing"

func TestPagination(t *testing.T) {
	p := Pagination{10, 5}

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
