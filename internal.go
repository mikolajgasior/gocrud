package crud

import (
	"reflect"
	"strings"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

func (c *CRUD) builder(obj interface{}) (*sqlbuilder.Builder, error) {
	name := c.builderName(obj)

	builder, ok := c.builders[name]
	if ok {
		return builder, nil
	}

	builder = sqlbuilder.New(obj, sqlbuilder.Options{
		TableNamePrefix: c.tableNamePrefix,
		TagName:         c.tagName,
	})
	if builder.Err() != nil {
		return nil, ErrCRUD{
			Op:  "builder.New",
			Err: builder.Err(),
		}
	}

	c.builders[name] = builder

	return builder, nil
}

func (c *CRUD) builderName(obj interface{}) string {
	objValue := reflect.ValueOf(obj)
	objIndirect := reflect.Indirect(objValue)
	objType := objIndirect.Type()

	if objType.String() != "reflect.Value" {
		return objType.Name()
	}

	objType = reflect.ValueOf(obj.(reflect.Value).Interface()).Type().Elem().Elem()
	name := objType.Name()

	if !strings.Contains(name, ".") {
		return name
	}

	nameArr := strings.Split(name, ".")
	return nameArr[1]
}
