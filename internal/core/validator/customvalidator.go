package validator

import (
	"fmt"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/pkg/dictionary"
	"gopher/pkg/helper"
	"gopher/pkg/logparser"
	"reflect"
	"regexp"
	"strings"
)

// BindTagExtractor extract the tag of bind and field from json and bind tag
func BindTagExtractor(engine *core.Engine, model interface{}, code, entity, action string) (err error) {

	// reflect interface to get value and struct tags
	reflexType := reflect.TypeOf(model)
	reflectValue := reflect.ValueOf(model)

	for i := 0; i < reflexType.NumField(); i++ {
		// get field
		field := reflexType.Field(i)
		// get value
		value := reflectValue.Field(i).Interface()

		// to get tag bind value in field
		bindTag := field.Tag.Get("bind")

		if bindTag == "-" || bindTag == "" {
			continue
		}

		// if reflect is struct, and it will validate those fields in nested struct
		if reflect.ValueOf(value).Kind() == reflect.Struct {

			// reflect interface to get value and struct tags
			reflexTypeN := reflect.TypeOf(value)
			reflectValueN := reflect.ValueOf(value)

			if reflexTypeN.String() == "time.Time" {
				continue
			}

			for j := 0; j < reflexTypeN.NumField(); j++ {
				// get field and value
				fieldN := reflexTypeN.Field(j)
				valueN := reflectValueN.Field(j).Interface()

				// validate method to check value and compare with tags,then get proper error
				err = validationCase(err, fieldN, valueN, action)
			}
		} else {
			// validate method to check value and compare with tags,then get proper error
			err = validationCase(err, field, value, action)
		}
	}

	if err != nil {
		err = engine.ErrorLog.TickValidate(err, code, entity, action, model)
	}
	return
}

// it compare value with bind tag conditions and it get proper error message
func validationCase(err error, field reflect.StructField, value interface{}, action string) (errors error) {

	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		value = reflectPointerToValue(value)
	}

	// get bind tag to binding field
	bindTag := field.Tag.Get("bind")

	requiredIfNotNil := ""

	// get field name
	jsonTag := field.Tag.Get("json")
	regex := regexp.MustCompile(`\w+`)
	fieldName := regex.FindString(jsonTag)

	// split tags by comma and find all
	tagByAction := BindTagByAction(bindTag, action)

	allTags := strings.Split(tagByAction, ",")
	for _, v := range allTags {

		// split taq by equal to find key and value
		splitTag := strings.Split(v, "=")
		var tagKey string
		var tagValue interface{}

		if len(splitTag) > 0 {
			tagKey = splitTag[0]
		}
		if len(splitTag) > 1 {
			tagValue = splitTag[1]
		}

		// custom validation cases
		switch tagKey {

		case "required":
			strValue := fmt.Sprintf("%v", value)
			if strValue == "" || strValue == "0" {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.VisRequired, dictionary.Translate(fieldName))
			}

		case "requiredifnotnil":
			strValue := fmt.Sprintf("%v", value)
			requiredIfNotNil = strValue
			if strValue == "" || strValue == "0" {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.VisRequired, dictionary.Translate(fieldName))
			}

		case "min":
			intTagValue, _ := helper.StrToInt(tagValue.(string))
			if len(value.(string)) < intTagValue {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.MinimumAcceptedCharacterForVisV, dictionary.Translate(fieldName), tagValue)
			}

		case "max":
			intTagValue, _ := helper.StrToInt(tagValue.(string))
			if len(value.(string)) > intTagValue {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.MaximumAcceptedCharacterForVisV, dictionary.Translate(fieldName), tagValue)
			}

		case "lte":
			floatTagValue, _ := helper.StrToFloat(tagValue.(string))
			if value.(float64) < floatTagValue {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.MinimumAcceptedValueForVisV, dictionary.Translate(fieldName), tagValue)
			}

		case "gte":
			floatTagValue, _ := helper.StrToFloat(tagValue.(string))
			if value.(float64) < floatTagValue {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.MaximumAcceptedValueForVisV, dictionary.Translate(fieldName), tagValue)
			}

		case "oneof":
			types := core.MustBeInTypes[tagValue.(string)]
			if ok, _ := helper.Includes(types, value); !ok {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.AcceptedValueForVareV, dictionary.Translate(fieldName), types)
			}

		case "contain":
			if ok := strings.Contains(value.(string), tagValue.(string)); !ok {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.InvalidValueDoNotIncloudV, tagValue)
			}

		case "email":
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if !re.MatchString(value.(string)) {
				err = logparser.AddInvalidParam(err, fieldName,
					terms.VisNotValid, dictionary.Translate(fieldName))
			}
		}
	}
	if strings.Contains(tagByAction, "requiredifnotnil") && requiredIfNotNil == "" {
		err = nil
	}
	return err
}

// BindTagByAction this function separate validation per type of action
func BindTagByAction(tagvalue string, action string) (result string) {

	perAction := helper.RegexFindBetweenTwoPattern(fmt.Sprintf("%v:", action), `\|`, tagvalue)
	if perAction != "" {
		result = perAction
	}

	all := helper.RegexFindBetweenTwoPattern(`all:`, `\|`, tagvalue)
	if all != "" {
		result += "," + all
	}

	if result == "" && !strings.Contains(tagvalue, "|") && !strings.Contains(tagvalue, ":") {
		result = tagvalue
	}

	return
}

func reflectPointerToValue(v interface{}) interface{} {

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	if rv.IsValid() {
		return rv.Interface()
	}
	return ""
}
