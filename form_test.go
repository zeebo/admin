package admin

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func BenchmarkLoadFlat(b *testing.B) {
	var x struct {
		X int
		Y bool
		Z string
	}
	var values = url.Values{"X": {"20"}, "Y": {"true"}, "Z": {"hello"}}
	for i := 0; i < b.N; i++ {
		Load(values, &x)
	}
}

func BenchmarkLoadNested(b *testing.B) {
	var x struct {
		X struct {
			Y struct {
				Z struct {
					A string
				}
			}
		}
	}
	var values = url.Values{"X.Y.Z.A": {"hello"}}
	for i := 0; i < b.N; i++ {
		Load(values, &x)
	}
}

func BenchmarkLoadLong(b *testing.B) {
	var x struct {
		X1  int
		X2  int
		X3  int
		X4  int
		X5  int
		X6  int
		X7  int
		X8  int
		X9  int
		X10 int
		X11 int
		X12 int
		X13 int
		X14 int
		X15 int
		X16 int
		X17 int
		X18 int
		X19 int
		X20 int
	}
	var values = url.Values{}
	for i := 1; i <= 20; i++ {
		values.Add(fmt.Sprintf("X%d", i), "20")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Load(values, &x)
	}
}

func BenchmarkLoadAlloc(b *testing.B) {
	var x struct {
		X *******int
	}
	var values = url.Values{"X": {"20"}}
	for i := 0; i < b.N; i++ {
		Load(values, &x)
	}
}

func BenchmarkLoadErrors(b *testing.B) {
	var x struct {
		X int
		Y int
		Z bool
	}
	var values = url.Values{"X": {"t"}, "Y": {"t"}, "Z": {"wtf"}}
	for i := 0; i < b.N; i++ {
		Load(values, &x)
	}
}

func BenchmarkLoadCataErrors(b *testing.B) {
	var x struct {
		X int
	}
	var values = url.Values{"X.Y": {"t"}}
	for i := 0; i < b.N; i++ {
		Load(values, &x)
	}
}

//compare two d types deeply.
func compare(one, two d) bool {
	if len(one) != len(two) {
		return false
	}

	for k, v1 := range one {
		v2, ex := two[k]
		if !ex {
			return false
		}

		//check strings
		if v1s, ok := v1.(string); ok {
			if v2s, ok := v2.(string); ok {
				if v1s == v2s {
					continue
				}
			}
		}

		//check dicts
		if v1d, ok := v1.(d); ok {
			if v2d, ok := v2.(d); ok {
				if compare(v1d, v2d) {
					continue
				}
			}
		}

		return false
	}
	return true
}

func compareErrs(one LoadingErrors, two []string) bool {
	if len(one) != len(two) {
		return false
	}

	for k := range one {
		//search for k in two
		var found bool
		for _, v := range two {
			if v == k {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func TestUnflatten(t *testing.T) {
	table := []struct {
		data     url.Values
		expected d
	}{
		{url.Values{"A": {"a"}, "B": {"b"}}, d{"A": "a", "B": "b"}},
		{url.Values{"A.B": {"b"}, "A.C": {"c"}}, d{"A": d{"B": "b", "C": "c"}}},
		{url.Values{"A.B.C": {"c"}, "A.B.D": {"d"}, "A.E": {"e"}}, d{"A": d{"E": "e", "B": d{"C": "c", "D": "d"}}}},
		{url.Values{"A.A.A": {"a"}}, d{"A": d{"A": d{"A": "a"}}}},
		{url.Values{"A": {"a", "b"}}, d{"A": "a"}},
		{url.Values{"A": {"a"}, "B.C": {"c"}, "B.D": {"d"}}, d{"A": "a", "B": d{"C": "c", "D": "d"}}},
	}

	for _, c := range table {
		ret := unflatten(c.data, "")
		if !compare(ret, c.expected) {
			t.Fatalf("Test case failed: %s\nExpected: %s\nGot: %s\n", c.data, c.expected, ret)
		}
	}
}

func TestLoadBasic(t *testing.T) {
	var x struct {
		X string
	}

	_, err := Load(url.Values{"X": {"hello"}}, &x)
	if err != nil {
		t.Fatal(err)
	}

	if x.X != "hello" {
		t.Fatalf("Expected %q. Got %q", "hello", x.X)
	}
}

func TestLoadAliasedType(t *testing.T) {
	type T string
	var x struct {
		X T
	}

	_, err := Load(url.Values{"X": {"hello"}}, &x)
	if err != nil {
		t.Fatal(err)
	}

	if string(x.X) != "hello" {
		t.Fatalf("Expected %q. Got %q", "hello", x.X)
	}
}

func TestLoadNested(t *testing.T) {
	var x struct {
		X string
		Z struct {
			B bool
		}
	}

	_, err := Load(url.Values{"X": {"hello"}, "Z.B": {"true"}}, &x)
	if err != nil {
		t.Fatal(err)
	}

	if x.X != "hello" {
		t.Fatalf("Expected %q. Got %q", "hello", x.X)
	}

	if x.Z.B != true {
		t.Fatalf("Expected %v. Got %v", true, x.Z.B)
	}
}

func TestLoadPointer(t *testing.T) {
	var x struct {
		X *string
	}

	_, err := Load(url.Values{"X": {"hello"}}, &x)
	if err != nil {
		t.Fatal(err)
	}

	if *x.X != "hello" {
		t.Fatalf("Expected %q. Got %q", "hello", *x.X)
	}
}

func TestErrors(t *testing.T) {
	//save some typing (punny!)
	type F []string

	table := []struct {
		data url.Values
		errs F
	}{
		{url.Values{"X": {"hello"}, "Z.B": {"true"}}, F{}},
		{url.Values{"X": {"hello"}, "Z.B": {"wtf"}}, F{"Z.B"}},
		{url.Values{"X": {"hello"}, "Y": {"twenty"}, "Z.B": {"wtf"}}, F{"Y", "Z.B"}},
	}
	var x struct {
		X string
		Y int
		Z struct {
			B bool
		}
	}

	for _, c := range table {
		val, err := Load(c.data, &x)
		if err != nil {
			t.Fatalf("Error while loading %v:\n%s", c.data, err)
		}
		if !compareErrs(val, c.errs) {
			t.Fatalf("Errors did not agree.\nExpected: %v\nGot %v", c.errs, val)
		}
	}
}

func TestLoadInvalidSchema(t *testing.T) {
	var x struct {
		X string
	}

	lerrs, err := Load(url.Values{"X.X": {"hello"}}, &x)
	if err == nil {
		t.Fatal("Expected an error loading dictionary into string")
	}

	if !compareErrs(lerrs, []string{}) {
		t.Fatalf("Expected no LoadingErrors. Got %v", lerrs)
	}
}

func TestLoadManyNested(t *testing.T) {
	var x struct {
		X struct {
			X struct {
				X struct {
					X int
				}
			}
		}
	}

	_, err := Load(url.Values{"X.X.X.X": {"20"}}, &x)
	if err != nil {
		t.Fatal(err)
	}

	if x.X.X.X.X != 20 {
		t.Fatalf("Expected %d. Got %d", 20, x.X.X.X.X)
	}

	val, err := Load(url.Values{"X.X.X.X": {"twenty"}}, &x)
	if err != nil {
		t.Fatal(err)
	}

	if _, ex := val["X.X.X.X"]; !ex {
		t.Fatalf("Expected key at 'X.X.X.X'. Got %v", val)
	}
}

func TestLoadIntoString(t *testing.T) {
	type F string
	var (
		x   F
		val = "hello"
	)

	if err := loadInto(reflect.ValueOf(&x), val); err != nil {
		t.Fatal(err)
	}

	if string(x) != val {
		t.Fatalf("Error loading in. Got %q expected %q", x, val)
	}
}

func TestLoadIntoInt(t *testing.T) {
	type F int
	var x F

	if err := loadInto(reflect.ValueOf(&x), "20"); err != nil {
		t.Fatal(err)
	}

	if int(x) != 20 {
		t.Fatalf("Error loading in. Got %v expected %v", x, "20")
	}
}

func TestLoadIntoUint(t *testing.T) {
	type F uint
	var x F

	if err := loadInto(reflect.ValueOf(&x), "20"); err != nil {
		t.Fatal(err)
	}

	if uint(x) != 20 {
		t.Fatalf("Error loading in. Got %v expected %v", x, "20")
	}
}

func TestLoadIntoFloat(t *testing.T) {
	type F float32
	var x F

	if err := loadInto(reflect.ValueOf(&x), "20"); err != nil {
		t.Fatal(err)
	}

	if float32(x) != 20 {
		t.Fatalf("Error loading in. Got %v expected %v", x, "20")
	}
}

func TestLoadIntoBool(t *testing.T) {
	type F bool
	var x F

	if err := loadInto(reflect.ValueOf(&x), "true"); err != nil {
		t.Fatal(err)
	}

	if bool(x) != true {
		t.Fatalf("Error loading in. Got %v expected %v", x, true)
	}
}

func TestLoadIntoPointer(t *testing.T) {
	var x *string

	if err := loadInto(reflect.ValueOf(&x), "hello"); err != nil {
		t.Fatal(err)
	}

	if *x != "hello" {
		t.Fatalf("Error loading in. Got %q expected %q", *x, "hello")
	}
}

func TestLoadIntoFailures(t *testing.T) {
	var x string
	if err := loadInto(reflect.ValueOf(x), "hello"); err == nil {
		t.Fatal("expected loadInto to fail on type that cannot be set")
	}
}
