package idformatter

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
)

const defaultGlobalPrefix = "1"

var regexForFormat = regexp.MustCompile("[a-z]+")
var regexForParse = regexp.MustCompile("^[a-z]*")

type TypeIDFormatter struct {
	globalPrefix   string
	shortTypeCache map[string]string
	plainTypes     map[string]string
}

func NewFormatter() api.IDFormatter {
	formatter := TypeIDFormatter{
		globalPrefix:   defaultGlobalPrefix,
		shortTypeCache: make(map[string]string),
		plainTypes:     make(map[string]string),
	}
	return &formatter
}

func (formatter *TypeIDFormatter) FormatID(id, idType string, schemas *client.Schemas) string {
	if id == "" || len(id) == 0 {
		return ""
	}

	_, ok := formatter.plainTypes[idType]
	if ok {
		return id
	}

	shortType := formatter.getShortType(idType, schemas)

	rune, _ := utf8.DecodeRuneInString(id)
	if unicode.IsDigit(rune) {
		return shortType + "!" + id
	}
	return shortType + id
}

func (formatter *TypeIDFormatter) ParseID(id string) string {
	if id == "" || len(id) == 0 {
		return ""
	}

	rune, _ := utf8.DecodeRuneInString(id)
	if unicode.IsLetter(rune) && !strings.HasPrefix(id, formatter.globalPrefix) {
		return id
	}

	if !strings.HasPrefix(id, formatter.globalPrefix) {
		return ""
	}

	id = id[len(formatter.globalPrefix):]

	rune, _ = utf8.DecodeRuneInString(id)
	if len(id) == 0 || !unicode.IsLetter(rune) {
		return ""
	}

	parsedID := regexForParse.ReplaceAllString(id, "")
	if strings.HasPrefix(parsedID, "!") {
		return parsedID[1:]
	}
	return parsedID
}

func (formatter *TypeIDFormatter) getShortType(idType string, schemas *client.Schemas) string {
	orginalType := idType
	schemaType, ok := formatter.shortTypeCache[idType]
	if ok {
		return schemaType
	}

	schema, schemaExists := schemas.CheckSchema(idType)

	if schemaExists {
		idType = schema.Id
	}

	shortType := formatter.globalPrefix + string(idType[0]) + regexForFormat.ReplaceAllString(idType, "")
	shortType = strings.ToLower(shortType)
	formatter.shortTypeCache[orginalType] = shortType
	return shortType
}
