package services

import (
	"net/url"
	"regexp"
	"strings"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/repository"
)

// Utility Responsible for parsing filters from a url.

type FilterColumns map[string]string

type ValidFilters []string

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func MakeFilterColumns(filters ValidFilters) FilterColumns {
	returnMap := make(map[string]string)
	for _, key := range filters {
		returnMap[key] = toSnakeCase(key)
	}
	return returnMap
}

// Validates the url.Values against the FilterColumns
// Returns the converted Filters used by repositories for searching.
// Error if any url value is provided that does not match any of the FilterColumns
func ParseURLFilters(urlParam url.Values, fc FilterColumns) (repository.Filters, error) {

	returnKeys := make(repository.Filters)

	errors := []error{}

	for urlKey, urlParamValue := range urlParam {
		columnName, present := fc[urlKey]

		if !present {
			errors = append(errors, apperror.BadRequest("Invalid Request Parameter: %s", urlKey))
		}

		if len(urlParamValue) > 0 {
			returnKeys[columnName] = urlParamValue[0]
		}
	}

	if len(errors) != 0 {
		return nil, apperror.Of(errors)
	}

	return returnKeys, nil
}
