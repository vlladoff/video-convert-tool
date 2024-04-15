package helper

import (
	"reflect"
)

func GetStructTagByField(structure any, fieldName, tagName string) (tag string) {
	st := reflect.TypeOf(structure)
	if stf, ok := st.FieldByName(fieldName); ok {
		tag = stf.Tag.Get(tagName)
	}

	return tag
}
