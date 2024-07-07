package utils

import (
	"strings"
	"unicode"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func RemoveAccent(str string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, str)
	return result
}

func SlugGeneratorWithId(str string, idInput primitive.ObjectID) (slug string, id primitive.ObjectID) {
	if idInput == primitive.NilObjectID {
		id = primitive.NewObjectID()
	} else {
		id = idInput
	}
	slug = RemoveAccent(str)
	slug = strings.ReplaceAll(strings.ToLower(slug), " ", "-") + "-" + id.Hex()
	return
}
