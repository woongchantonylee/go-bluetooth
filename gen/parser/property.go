package parser

import (
	"errors"
	"regexp"
	"strings"

	"github.com/woongchantonylee/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

const propBaseRegexp = `(bool|boolean|byte|string|int16|uint16|uint16_t|uint32|dict|object|array\{.*?) ([A-Z].+?)`

type PropertyParser struct {
	model *types.Property
	debug bool
}

// NewPropertyParser
func NewPropertyParser(debug bool) PropertyParser {
	return PropertyParser{
		model: new(types.Property),
		debug: debug,
	}
}

func (g *PropertyParser) Parse(raw []byte) (*types.Property, error) {

	var err error
	property := g.model
	// log.Debugf("prop raw -> %s", raw)

	re1 := regexp.MustCompile(`[ \t]*?` + propBaseRegexp + `( \[[^\]]*\].*)?\n((?s).+)`)
	matches2 := re1.FindAllSubmatch(raw, -1)

	// log.Warnf("m1 %s", matches2)

	if len(matches2) == 0 || len(matches2[0]) == 1 {
		re1 = regexp.MustCompile(`[ \t]*?` + propBaseRegexp + `\n((?s).+)`)
		matches2 = re1.FindAllSubmatch(raw, -1)
		// log.Warnf("m2 %s", matches2)
	}

	if len(matches2) == 0 {
		log.Debugf("prop raw -> %s", raw)
		return property, errors.New("No property found")
	}

	flags := []types.Flag{}
	flagListRaw := string(matches2[0][3])
	flagList := strings.Split(strings.Trim(flagListRaw, "[] "), ",")

	for _, f := range flagList {

		// int16 Handle [read-write, optional] (Server Only)
		if strings.Contains(f, "]") {
			f = strings.Split(f, "]")[0]
		}

		var flag types.Flag
		switch f {
		case "readonly":
			{
				flag = types.FlagReadOnly
			}
		case "readwrite":
			{
				flag = types.FlagReadWrite
			}
		case "experimental":
			{
				flag = types.FlagExperimental
			}
		}

		if flag > 0 {
			flags = append(flags, flag)
		}

	}

	docs := string(matches2[0][4])
	docs = strings.Replace(docs, " \t\n", "", -1)
	docs = strings.Trim(docs, " \t\n")

	name := string(matches2[0][2])

	if strings.Contains(name, "optional") {
		name = strings.Replace(name, " (optional)", "", -1)
		docs = "(optional) " + docs
	}

	name = strings.Replace(name, " \t\n", "", -1)

	property.Type = string(matches2[0][1])
	property.Name = name
	property.Flags = flags
	property.Docs = docs

	if g.debug {
		log.Debugf("\t - %s", property)
	}
	return property, err
}
