package input

import (
	"fmt"
	"html"
	"reflect"
	"strconv"
	"strings"
)

type Options struct {
	RestrictFields  map[string]struct{}
	ExcludeFields   map[string]struct{}
	TagName         string
	IDPrefix        string
	NamePrefix      string
	OverwriteValues map[string]string
	Values          bool
}

func BuildInputs(obj interface{}, options *Options) ([]Input, []string) {
	v := reflect.ValueOf(obj)
	i := reflect.Indirect(v)
	s := i.Type()
	elem := v.Elem()

	tagName := DefaultTagName
	if options != nil && options.TagName != "" {
		tagName = options.TagName
	}

	inputs := []Input{}
	fieldsOrder := []string{}

	for j := 0; j < s.NumField(); j++ {
		field := s.Field(j)
		fieldKind := field.Type.Kind()

		// check if only specified field should be generated
		if options != nil {
			if len(options.RestrictFields) > 0 {
				_, ok := options.RestrictFields[field.Name]
				if !ok {
					continue
				}
			}
			if len(options.ExcludeFields) > 0 {
				_, ok := options.ExcludeFields[field.Name]
				if ok {
					continue
				}
			}
		}

		if !isInt(fieldKind) && fieldKind != reflect.String && fieldKind != reflect.Bool {
			continue
		}

		input := Input{
			FieldName: field.Name,
			Name:      field.Name,
		}

		if options != nil {
			if options.Values {
				if fieldKind == reflect.Bool && elem.Field(j).Bool() {
					input.Value = "true"
					input.Checked = true
				}
				if fieldKind == reflect.String {
					input.Value = elem.Field(j).String()
				}
				if isInt(fieldKind) {
					input.Value = fmt.Sprintf("%d", elem.Field(j).Int())
				}
			}

			if len(options.OverwriteValues) > 0 {
				overwriteValue, ok := options.OverwriteValues[field.Name]
				if ok {
					input.Value = overwriteValue
				}
			}

			if options.IDPrefix != "" {
				input.ID = options.IDPrefix + field.Name
			}

			if options.NamePrefix != "" {
				input.Name = options.NamePrefix + field.Name
			}
		}

		tagVal := field.Tag.Get(tagName)
		tagRegexpVal := field.Tag.Get(tagName + "_regexp")

		if tagRegexpVal != "" {
			input.Pattern = html.EscapeString(tagRegexpVal)
		}

		validationAttrs, inputType := attributes(tagVal)
		input.InputType = inputType

		parseValidationAttrs(validationAttrs, &input)

		if fieldKind == reflect.Bool {
			input.InputType = TypeCheckbox
		} else if isInt(fieldKind) {
			input.InputType = TypeNumber
		} else if input.InputType == "" {
			input.InputType = TypeText
		}

		if input.InputType == TypePassword {
			input.Value = ""
		}

		inputs = append(inputs, input)
		fieldsOrder = append(fieldsOrder, field.Name)
	}

	return inputs, fieldsOrder
}

func parseValidationAttrs(attrs string, input *Input) {
	if strings.Contains(attrs, "required") {
		input.Required = true
	}

	if strings.Contains(attrs, "pattern=") {
		start := strings.Index(attrs, `pattern="`) + 9
		end := strings.Index(attrs[start:], `"`)
		if end > 0 {
			input.Pattern = attrs[start : start+end]
		}
	}

	if strings.Contains(attrs, "minlength=") {
		start := strings.Index(attrs, `minlength="`) + 11
		end := strings.Index(attrs[start:], `"`)
		if end > 0 {
			val, _ := strconv.Atoi(attrs[start : start+end])
			input.MinLength = val
		}
	}

	if strings.Contains(attrs, "maxlength=") {
		start := strings.Index(attrs, `maxlength="`) + 11
		end := strings.Index(attrs[start:], `"`)
		if end > 0 {
			val, _ := strconv.Atoi(attrs[start : start+end])
			input.MaxLength = val
		}
	}

	if strings.Contains(attrs, "min=") {
		start := strings.Index(attrs, `min="`) + 5
		end := strings.Index(attrs[start:], `"`)
		if end > 0 {
			val, _ := strconv.Atoi(attrs[start : start+end])
			input.Min = val
		}
	}

	if strings.Contains(attrs, "max=") {
		start := strings.Index(attrs, `max="`) + 5
		end := strings.Index(attrs[start:], `"`)
		if end > 0 {
			val, _ := strconv.Atoi(attrs[start : start+end])
			input.Max = val
		}
	}
}

func attributes(tag string) (string, string) {
	attrs := ""
	inputType := TypeText

	opts := strings.SplitN(tag, " ", -1)
	for _, opt := range opts {
		if opt == "req" {
			attrs = attrs + " required"
		}
		if opt == "uiemail" {
			inputType = TypeEmail
			continue
		}
		if opt == "uitextarea" {
			inputType = TypeTextarea
		}
		if opt == "uipassword" {
			inputType = TypePassword
		}
		for _, valOpt := range []string{"len", "val", "regexp"} {
			if !strings.HasPrefix(opt, valOpt+":") {
				continue
			}

			val := strings.Replace(opt, valOpt+":", "", 1)
			if valOpt == "regexp" {
				attrs = attrs + fmt.Sprintf(` pattern="%s"`, html.EscapeString(val))
				continue
			}

			minMax := strings.Split(val, ",")
			if minMax[0] != "" {
				min, err := strconv.Atoi(minMax[0])
				if err == nil {
					if valOpt == "len" {
						attrs = attrs + fmt.Sprintf(` minlength="%d"`, min)
					}
					if valOpt == "val" {
						attrs = attrs + fmt.Sprintf(` min="%d"`, min)
					}
				}
			}
			if len(minMax) > 1 && minMax[1] != "" {
				max, err := strconv.Atoi(minMax[1])
				if err == nil {
					if valOpt == "len" {
						attrs = attrs + fmt.Sprintf(` maxlength="%d"`, max)
					}
					if valOpt == "val" {
						attrs = attrs + fmt.Sprintf(` max="%d"`, max)
					}
				}
			}
		}
	}

	return attrs, inputType
}

func isInt(k reflect.Kind) bool {
	return k == reflect.Int64 || k == reflect.Int32 || k == reflect.Int16 || k == reflect.Int8 || k == reflect.Int || k == reflect.Uint64 || k == reflect.Uint32 || k == reflect.Uint16 || k == reflect.Uint8 || k == reflect.Uint
}
