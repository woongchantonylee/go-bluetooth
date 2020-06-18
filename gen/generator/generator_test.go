package generator

import (
	"fmt"
	"testing"

	"github.com/woongchantonylee/go-bluetooth/gen"
	"github.com/woongchantonylee/go-bluetooth/gen/util"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {

	TplPath = "../../gen/generator/tpl/%s.go.tpl"
	outdir := "../../test/out"

	bluezApi, err := gen.Parse("../../src/bluez/doc", []string{}, false)
	if err != nil {
		t.Fatal(err)
	}

	err = util.Mkdir("../../test")
	if err != nil {
		t.Fatal(err)
	}
	err = util.Mkdir(outdir)
	if err != nil {
		t.Fatal(err)
	}

	err = Generate(bluezApi, outdir, true, true)
	if err != nil {
		t.Fatal(err)
	}

	assert.DirExists(t, outdir)
	assert.DirExists(t, fmt.Sprintf("%s/profile/adapter", outdir))

}
