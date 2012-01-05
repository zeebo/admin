package admin

import "testing"

type T struct{}

func TestRegister(t *testing.T) {
	Register(T{}, "T", nil)

	ans := GetType("T")
	if _, ok := ans.(*T); !ok {
		t.Fatalf("Type incorrect. Expected *admin.T, got %T", ans)
	}
}
