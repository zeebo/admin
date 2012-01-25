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

func BenchmarkCreateValues(b *testing.B) {
	type flat struct {
		X string
		Y int
		Z bool
	}
	type nested struct {
		X string
		Y struct {
			Z int
		}
		Q flat
	}
	var x = nested{"food", struct{ Z int }{-8000}, flat{"doof", 20, false}}
	for i := 0; i < b.N; i++ {
		CreateValues(x)
	}
}

func BenchmarkCreateEmptyValues(b *testing.B) {
	type flat struct {
		X **string
		Y *int
		Z ***bool
	}
	type nested struct {
		X *string
		Y **struct {
			Z *int
		}
		Q ***flat
	}
	var x nested
	for i := 0; i < b.N; i++ {
		CreateEmptyValues(x)
	}
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
			t.Errorf("Test case failed: %s\nExpected: %s\nGot: %s\n", c.data, c.expected, ret)
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

func TestLoadPointerEmpty(t *testing.T) {
	var x struct {
		X *string
	}

	_, err := Load(url.Values{"X": {""}}, &x)
	if err != nil {
		t.Fatal(err)
	}

	if *x.X != "" {
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
			t.Errorf("Error while loading %v:\n%s", c.data, err)
		}
		if !compareErrs(val, c.errs) {
			t.Errorf("Errors did not agree.\nExpected: %v\nGot %v", c.errs, val)
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

func TestLoadInvalidTypes(t *testing.T) {
	var (
		x1 = struct{ X func() }{}
		x2 = struct{ X []string }{}
		x3 = struct{ X interface{} }{}
		x4 = struct{ X map[string]string }{}
		x5 = struct{ X [5]string }{}
		x6 = struct{ X chan string }{}
		x7 = struct{ X uintptr }{}
		x8 = struct{ X complex64 }{}
		x9 = struct{ X complex128 }{}
	)

	if _, err := Load(url.Values{"X": {"data"}}, &x1); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x1)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x2); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x2)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x3); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x3)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x4); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x4)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x5); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x5)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x6); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x6)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x7); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x7)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x8); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x8)
	}

	if _, err := Load(url.Values{"X": {"data"}}, &x9); err == nil {
		t.Errorf("Error. Allowed loading into a %T", x9)
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

func TestCreateValuesValid(t *testing.T) {
	type flat struct {
		X string
		Y int
		Z bool
	}
	type nested struct {
		X string
		Y struct {
			Z int
		}
		Q flat
	}
	type F map[string]string
	table := []struct {
		data     interface{}
		expected F
	}{
		{flat{"foo", 2, true}, F{"X": "foo", "Y": "2", "Z": "true"}},
		{flat{"foob", -2, false}, F{"X": "foob", "Y": "-2", "Z": "false"}},
		{nested{"foo", struct{ Z int }{8}, flat{"foob", 2, true}}, F{"X": "foo", "Y.Z": "8", "Q.X": "foob", "Q.Y": "2", "Q.Z": "true"}},
		{nested{"food", struct{ Z int }{-8000}, flat{"doof", 20, false}}, F{"X": "food", "Y.Z": "-8000", "Q.X": "doof", "Q.Y": "20", "Q.Z": "false"}},
	}

	for _, c := range table {
		ret, err := CreateValues(c.data)
		if err != nil {
			t.Errorf("Error processing: %s\n%s", err, c.data)
		}

		if !compareString(ret, c.expected) {
			t.Errorf("Test case failed: %s\nExpected: %s\nGot: %s\n", c.data, c.expected, ret)
		}
	}
}

func TestCreateValuesInValid(t *testing.T) {
	var x struct {
		X *int
	}

	_, err := CreateValues(x)
	if err == nil {
		t.Fatal("Expected error for nil pointer")
	}
}

func TestCreateEmptyValuesValid(t *testing.T) {
	type F map[string]string
	table := []struct {
		data     interface{}
		expected F
	}{
		{struct {
			X *int
			Y **int
			Z **bool
		}{}, F{"X": "", "Y": "", "Z": ""}},
		{struct {
			X *int
			Y int
			Z *******int
		}{}, F{"X": "", "Y": "", "Z": ""}},
		{struct {
			X **struct {
				Y *int
				Z string
			}
		}{}, F{"X.Y": "", "X.Z": ""}},
	}

	for _, c := range table {
		ret, err := CreateEmptyValues(c.data)
		if err != nil {
			t.Errorf("Error processing: %s\n%s", err, c.data)
		}

		if !compareString(ret, c.expected) {
			t.Errorf("Test case failed: %s\nExpected: %s\nGot: %s\n", c.data, c.expected, ret)
		}
	}
}

func TestCreateEmptyValuesInvalidTypes(t *testing.T) {
	var (
		x1 = struct{ X func() }{}
		x2 = struct{ X []string }{}
		x3 = struct{ X interface{} }{}
		x4 = struct{ X map[string]string }{}
		x5 = struct{ X [5]string }{}
		x6 = struct{ X chan string }{}
		x7 = struct{ X uintptr }{}
		x8 = struct{ X complex64 }{}
		x9 = struct{ X complex128 }{}
	)

	if _, err := CreateEmptyValues(x1); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x1)
	}

	if _, err := CreateEmptyValues(x2); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x2)
	}

	if _, err := CreateEmptyValues(x3); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x3)
	}

	if _, err := CreateEmptyValues(x4); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x4)
	}

	if _, err := CreateEmptyValues(x5); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x5)
	}

	if _, err := CreateEmptyValues(x6); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x6)
	}

	if _, err := CreateEmptyValues(x7); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x7)
	}

	if _, err := CreateEmptyValues(x8); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x8)
	}

	if _, err := CreateEmptyValues(x9); err == nil {
		t.Errorf("Error. Creating an empty value with a %T", x9)
	}
}

func TestCreateEmptyValuesInvalid(t *testing.T) {
	r, err := CreateEmptyValues(struct {
		X []string
	}{})
	if err == nil {
		t.Fatalf("Expected an error because slices are unsupported. Got a %s", r)
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

func compareString(one, two map[string]string) bool {
	if len(one) != len(two) {
		return false
	}
	for k, v1 := range one {
		v2, ex := two[k]
		if !ex {
			return false
		}
		if v1 != v2 {
			return false
		}
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
