package api

import (
	"fmt"
	"testing"

	"github.com/woongchantonylee/go-bluetooth/bluez"
	"github.com/woongchantonylee/go-bluetooth/props"
	"github.com/woongchantonylee/go-bluetooth/util"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	IgnoreFlag        bool                   `dbus:"ignore"`
	ToOmit            map[string]interface{} `dbus:"omitEmpty,writable"`
	Ignored           string                 `dbus:"ignore"`
	IgnoredByProperty []string               `dbus:"ignore=IgnoreFlag"`
	Avail             string
}

func (s testStruct) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	util.StructToMap(s, m)
	return m, nil
}

func (s testStruct) Lock()   {}
func (s testStruct) Unlock() {}

func TestParseTag(t *testing.T) {

	s := testStruct{
		IgnoreFlag:        true,
		Ignored:           "foo",
		IgnoredByProperty: []string{"bar"},
		Avail:             "foo",
	}

	prop := &DBusProperties{
		props:       make(map[string]bluez.Properties),
		propsConfig: make(map[string]map[string]*props.PropInfo),
	}

	err := prop.AddProperties("test", s)
	if err != nil {
		t.Fatal(err)
	}

	err = prop.parseProperties()
	if err != nil {
		t.Fatal(err)
	}

	cfg := prop.propsConfig["test"]

	for field, cfg := range cfg {
		fmt.Printf("%s: %++v\n", field, cfg)
	}

	assert.True(t, cfg["ToOmit"].Skip)
	assert.True(t, cfg["ToOmit"].Writable)
	assert.True(t, cfg["Ignored"].Skip)
	assert.True(t, cfg["IgnoredByProperty"].Skip)
	assert.Equal(t, "foo", cfg["Avail"].Value)

}
