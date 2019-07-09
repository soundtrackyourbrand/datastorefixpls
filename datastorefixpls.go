package datastorefixpls

import (
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
)

// SaveStruct serializes every time.Time as both with and without dot suffix
func SaveStruct(x interface{}) ([]datastore.Property, error) {
	ps, err := datastore.SaveStruct(x)
	if err != nil {
		return ps, err
	}

	extra := []datastore.Property{}

	for _, p := range ps {
		if _, isTime := p.Value.(time.Time); isTime {
			extra = append(extra, datastore.Property{
				Name:     correspondingName(p.Name),
				Value:    p.Value,
				NoIndex:  p.NoIndex,
				Multiple: p.Multiple,
			})
		}
	}

	// make sure we dont add duplicates
	for _, p := range extra {
		if !exist(ps, p) {
			ps = append(ps, p)
		}
	}

	return ps, nil
}

// LoadStruct deserializes every field ending with a dot as both with and without the dot suffix
func LoadStruct(x interface{}, ps []datastore.Property) error {
	extra := []datastore.Property{}

	for _, p := range ps {
		if strings.HasSuffix(p.Name, ".") {
			extra = append(extra, datastore.Property{
				Name:     p.Name[:len(p.Name)-1],
				Value:    p.Value,
				NoIndex:  p.NoIndex,
				Multiple: p.Multiple,
			})
		}
	}

	return datastore.LoadStruct(x, append(ps, extra...))
}

func exist(ps []datastore.Property, needle datastore.Property) bool {
	for _, p := range ps {
		if p.Name == needle.Name {
			return true
		}
	}
	return false
}

func correspondingName(name string) string {
	if strings.HasSuffix(name, ".") {
		return name[:len(name)-1]
	}

	return name + "."
}
