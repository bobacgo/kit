package tag

import "reflect"

var basicKind = []reflect.Kind{
	reflect.Bool,
	reflect.String,
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
}

var multiElement = []reflect.Kind{
	reflect.Slice, reflect.Array,
	reflect.Map,
	reflect.Struct,
}
