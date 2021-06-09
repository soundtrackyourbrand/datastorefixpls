package datastorefixpls

import (
	"sort"
	"testing"
	"time"

	"google.golang.org/appengine/v2/datastore"
)

var now = time.Now()

func TestAddsTimeSuffix(t *testing.T) {
	ps := []datastore.Property{
		{Name: "Field.", Value: now},
	}

	expected := []datastore.Property{
		{Name: "Field.", Value: now},
		{Name: "Field.Time", Value: now},
	}

	ps = denormalize(ps)

	if !equal(ps, expected) {
		t.Fatal("expected slices to be equal")
	}
}

func TestAddsDotSuffix(t *testing.T) {
	ps := []datastore.Property{
		{Name: "Field.Time", Value: now},
	}

	expected := []datastore.Property{
		{Name: "Field.", Value: now},
		{Name: "Field.Time", Value: now},
	}

	ps = denormalize(ps)

	if !equal(ps, expected) {
		t.Fatal("expected slices to be equal")
	}
}

func TestDontAddDuplicates(t *testing.T) {
	ps := []datastore.Property{
		{Name: "Field.", Value: now},
		{Name: "Field.Time", Value: now},
	}

	expected := []datastore.Property{
		{Name: "Field.", Value: now},
		{Name: "Field.Time", Value: now},
	}

	ps = denormalize(ps)

	if !equal(ps, expected) {
		t.Fatal("expected slices to be equal")
	}
}

func TestEqualWithoutOrderWorks(t *testing.T) {
	a := []datastore.Property{
		{Name: "A", Value: now},
		{Name: "B", Value: now},
	}
	b := []datastore.Property{
		{Name: "B", Value: now},
		{Name: "A", Value: now},
	}

	if !equal(a, b) {
		t.Fatal("expected a and b to be equal")
	}
}

func equal(a, b []datastore.Property) bool {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Name > a[j].Name
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].Name > b[j].Name
	})

	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
