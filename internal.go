package gocrud

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	sqlfilters "codeberg.org/mikolajgasior/gocrud/pkg/filters"
	sqlbuilder "codeberg.org/mikolajgasior/gocrud/pkg/sql/builder"
)

func (c *CRUD) builder(obj interface{}) (*sqlbuilder.Builder, error) {
	name := c.builderName(obj)

	c.buildersMu.RLock()
	builder, ok := c.builders[name]
	c.buildersMu.RUnlock()
	if ok {
		return builder, nil
	}

	builder = sqlbuilder.New(obj, sqlbuilder.Options{
		TableNamePrefix: c.tableNamePrefix,
		TagName:         c.tagName,
		Dialect:         c.dialect,
	})

	c.buildersMu.Lock()
	c.builders[name] = builder
	c.buildersMu.Unlock()

	return builder, nil
}

func (c *CRUD) builderName(obj interface{}) string {
	objValue := reflect.ValueOf(obj)
	objIndirect := reflect.Indirect(objValue)
	objType := objIndirect.Type()

	if objType.String() != "reflect.Value" {
		return objType.Name()
	}

	objType = reflect.ValueOf(obj.(reflect.Value).Interface()).Type().Elem()
	name := objType.Name()

	if !strings.Contains(name, ".") {
		return name
	}

	nameArr := strings.Split(name, ".")
	return nameArr[1]
}

func (c *CRUD) runOnDelete(ctx context.Context, obj interface{}, ids []uint64, lastDepth int8) error {
	objValue := reflect.ValueOf(obj)
	objIndirect := reflect.Indirect(objValue)
	objType := objIndirect.Type()

	var objName string
	if objType.String() != "reflect.Value" {
		objName = objType.Name()
	} else {
		objType = reflect.ValueOf(obj.(reflect.Value).Interface()).Type().Elem()
		objName = objType.Name()
		if strings.Contains(objName, ".") {
			objNameArray := strings.Split(objName, ".")
			objName = objNameArray[1]
		}
	}

	defaultParentIDField := objName + "ID"

	// tagWithValRegexp
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldKind := field.Type.Kind()

		// Only fields that are slices of struct instances
		if fieldKind != reflect.Slice || field.Type.Elem().Kind() != reflect.Struct {
			continue
		}

		// Get a field tag, search for an 'on_del' tag and determine further action.
		tag := field.Tag.Get(c.tagName)
		if tag == "" {
			continue
		}

		tags := strings.Split(tag, " ")
		tagsMap := map[string]string{}
		for _, t := range tags {
			if m := tagWithValRegexp.MatchString(t); m {
				mArr := strings.Split(t, ":")
				tagsMap[mArr[0]] = mArr[1]
			}
		}

		// Perform delete
		if tagsMap["on_del"] == "del" {
			parentIDField := defaultParentIDField
			if tagsMap["del_field"] != "" {
				parentIDField = tagsMap["del_field"]
			}

			// Delete from the child table where parent ID = id of the deleted object.
			errDelete := c.DeleteMultiple(ctx, reflect.New(field.Type.Elem()), DeleteMultipleOptions{
				Filters: &sqlfilters.Filters{
					sqlfilters.Raw: {
						Op: sqlfilters.OpAND,
						Val: []interface{}{
							fmt.Sprintf(".%s IN (?)", parentIDField),
							ids,
						},
					},
				},
				CascadeDeleteDepth: lastDepth + 1,
			})
			if errDelete != nil {
				return getCRUDFuncCRUDError("delete multiple", errDelete)
			}
		}

		// Perform update
		if tagsMap["on_del"] == "upd" {
			updateField := tagsMap["del_upd_field"]
			updateValue := tagsMap["del_upd_val"]
			if updateField == "" {
				return getUpdateFieldFromTagsCRUDError()
			}

			parentIDField := defaultParentIDField
			if tagsMap["del_field"] != "" {
				parentIDField = tagsMap["del_field"]
			}
			// Update the child table where parent ID = id of the deleted object.
			errUpdate := c.UpdateMultiple(ctx, reflect.New(field.Type.Elem()),
				map[string]interface{}{
					updateField: updateValue,
				},
				UpdateMultipleOptions{
					Filters: &sqlfilters.Filters{
						sqlfilters.Raw: {
							Op: sqlfilters.OpAND,
							Val: []interface{}{
								fmt.Sprintf(".%s IN (?)", parentIDField),
								ids,
							},
						},
					},
					ConvertValuesFromString: true,
				},
			)
			if errUpdate != nil {
				return getCRUDFuncCRUDError("update multiple", errUpdate)
			}
		}
	}

	return nil
}
