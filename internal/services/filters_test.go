package services

import (
	"fmt"
	"net/url"
	"testing"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/repository"
)

var parseUrlFiltersTests = []struct {
	urlParams       url.Values
	fc              FilterColumns
	expectedFilters repository.Filters
	epxectedError   error
}{
	// Success case
	{
		urlParams: url.Values{
			"FirstName":        []string{"Rob"},
			"FirstAndLastName": []string{"Rob Test"},
		},
		fc: FilterColumns{"FirstName": "first_name", "FirstAndLastName": "first_and_last_name"},
		expectedFilters: repository.Filters{
			"first_name":          "Rob",
			"first_and_last_name": "Rob Test",
		},
		epxectedError: nil,
	},
	// Single Error Case
	{
		urlParams: url.Values{
			"FirstName": []string{"Rob"},
			"FirstNa":   []string{"Rob Test"},
		},
		fc:              FilterColumns{"FirstName": "first_name", "FirstAndLastName": "first_and_last_name"},
		expectedFilters: nil,
		epxectedError:   apperror.BadRequest("Invalid Request Parameter: FirstNa"),
	},
	// Multi Error Case
	{
		urlParams: url.Values{
			"FirstNa": []string{"Rob"},
			"FooBar":  []string{"Rob Test"},
		},
		fc:              FilterColumns{"FirstName": "first_name", "FirstAndLastName": "first_and_last_name"},
		expectedFilters: nil,
		epxectedError:   apperror.Of([]error{apperror.BadRequest("Invalid Request Parameter: FirstNa"), apperror.BadRequest("Invalid Request Parameter: FooBar")}),
	},
}

func Test_ParseURLFilters(t *testing.T) {

	for idx, tc := range parseUrlFiltersTests {
		t.Run(fmt.Sprintf("ParseURLFilters Test Case: %d", idx), func(t *testing.T) {
			outFilters, outError := ParseURLFilters(tc.urlParams, tc.fc)

			// Wonder if reflect.DeepEquals is better to use to compare errors...
			if fmt.Sprint(outError) != fmt.Sprint(tc.epxectedError) {
				t.Errorf("Out Error was not as expected Want: %s Got: %s", tc.epxectedError, outError)
			}

			if fmt.Sprint(tc.expectedFilters) != fmt.Sprint(outFilters) {
				t.Errorf("Out Filters was not as expected Want: %s Got: %s", tc.expectedFilters, outFilters)
			}

		})
	}
}

func Test_MakeFilterColumns(t *testing.T) {
	validFilters := MakeFilterColumns(ValidFilters{"FirstName", "FirstAndLastName", "Email", ""})

	expected := FilterColumns{
		"FirstName": "first_name",
		"Email":     "email",
		"":          "",
	}

	for key, expectedVal := range expected {
		if value := validFilters[key]; value != expectedVal {
			t.Errorf("Filter Was not as expected. Want: %s got %s", expectedVal, value)
		}
	}
}
