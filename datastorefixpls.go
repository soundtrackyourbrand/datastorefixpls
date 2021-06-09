package datastorefixpls

import (
	"strings"
	"time"

	"google.golang.org/appengine/v2/datastore"
)

// SaveStruct serializes every time.Time as both with and without dot suffix
func SaveStruct(x interface{}) ([]datastore.Property, error) {
	ps, err := datastore.SaveStruct(x)
	if err != nil {
		return ps, err
	}

	return denormalize(ps), nil
}

// LoadStruct deserializes every field ending with a dot as both with and without the dot suffix
func LoadStruct(x interface{}, ps []datastore.Property) error {
	return datastore.LoadStruct(x, denormalize(ps))
}

func exist(ps []datastore.Property, needle datastore.Property) bool {
	for _, p := range ps {
		if p.Name == needle.Name {
			return true
		}
	}
	return false
}

func denormalizeName(p datastore.Property) (datastore.Property, bool) {
	if _, isTime := p.Value.(time.Time); isTime {
		if strings.HasSuffix(p.Name, ".Time") {
			p.Name = p.Name[:len(p.Name)-len("Time")]
			return p, true
		}

		if strings.HasSuffix(p.Name, ".") {
			p.Name = p.Name + "Time"
			return p, true
		}
	}

	return datastore.Property{}, false
}

func denormalize(ps []datastore.Property) []datastore.Property {
	extra := []datastore.Property{}

	for _, p := range ps {
		if prop, ok := denormalizeName(p); ok {
			extra = append(extra, prop)
		}
	}

	// make sure we dont add duplicates
	for _, p := range extra {
		if !exist(ps, p) {
			ps = append(ps, p)
		}
	}

	return ps
}
