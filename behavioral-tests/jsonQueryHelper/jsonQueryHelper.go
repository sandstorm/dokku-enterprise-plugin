package jsonQueryHelper

import (
	"fmt"
	"strings"
	"reflect"
	"time"
	"regexp"
	"github.com/elgs/gojq"
	"github.com/DATA-DOG/godog/gherkin"
	"math"
	"strconv"
)

func AssertJsonStructure(jsonString string, comparators *gherkin.DataTable) error {
	jsonQuery, err := gojq.NewStringQuery(jsonString)

	if err != nil {
		return fmt.Errorf("Cannot deserialize JSON response: %s", string(jsonString))
	}

	for i, value := range comparators.Rows {
		if len(value.Cells) != 3 {
			return fmt.Errorf("every comparison needs to be written in exactly 3 columns; but in row %d I found %d.", i + 1, len(value.Cells))
		}

		fieldPath := value.Cells[0].Value
		comparator := value.Cells[1].Value
		operand := value.Cells[2].Value

		value, err := jsonQuery.Query(fieldPath)

		if comparator != "is empty" && err != nil {
			return fmt.Errorf("Error doing query %v: %v. JSON: %v", fieldPath, err, jsonString)
		}

		switch comparator {
		case "equals":
			switch value.(type) {
			case string:
				if strings.Compare(value.(string), operand) != 0 {
					return fmt.Errorf("String '%s' is not equal to expected string '%s' at path %s. (line %d)", value, operand, fieldPath, i)
				}
			case float64:
				floatOperand, _ := strconv.ParseFloat(operand, 64)
				if math.Abs(value.(float64) - floatOperand) > 0.0000001 {
					return fmt.Errorf("Number '%d' is not equal to expected number '%d' at path %s. (line %d)", value, operand, fieldPath, i)
				}
			default:
				return fmt.Errorf("equals comparison for unknown type value in line %d - type was: %s", i, reflect.TypeOf(value))
			}
		case "contains":
			switch value.(type) {
			case string:
				if !strings.Contains(value.(string), operand) {
					return fmt.Errorf("String '%s' does not contain '%s' at path %s. (line %d)", value, operand, fieldPath, i)
				}
			default:
				return fmt.Errorf("contains comparison for unknown type value in line %d - type was: %s", i, reflect.TypeOf(value))
			}
		case "is a date":
			switch value.(type) {
			case string:
				_, e := time.Parse(time.RFC3339, value.(string))
				if e != nil {
					return fmt.Errorf("Timestamp '%s' could not be parsed at path %s. (line %d)", value, fieldPath, i)
				}
			default:
				return fmt.Errorf("'is a date' comparison for unknown type value in line %d - type was: %s", i, reflect.TypeOf(value))
			}

		case "matches regex":
			switch value.(type) {
			case string:
				hasMatched, err := regexp.MatchString(operand, value.(string))
				if err != nil {
					return err
				}
				if !hasMatched {
					return fmt.Errorf("String '%s' does not match to expected regex '%s' at path %s. (line %d)", value, operand, fieldPath, i)
				}
			default:
				return fmt.Errorf("matches regex' comparison in line %d only works for type string - type was: %s", i, reflect.TypeOf(value))
			}
		case "is empty":
			if value != nil && len(value.([]interface{})) > 0 {
				return fmt.Errorf("'%v' is not empty at path %s. (line %d)", value, fieldPath, i)
			}
		case "count":
			operandInt, _ := strconv.Atoi(operand)

			switch value.(type) {
			case map[string]interface{}:
				if len(value.(map[string]interface{})) != operandInt {
					return fmt.Errorf("Value '%s' expected to have a count of %d, but is %d at path %s. (line %d)", value, operandInt, len(value.(map[string]interface{})), fieldPath, i)
				}
			default:
				return fmt.Errorf("count comparison in line %d only works for type map - type was: %s", i, reflect.TypeOf(value))
			}




		default:
			return fmt.Errorf("Comparator %s not supported (line %d)", comparator, i)
		}

	}

	return nil
}