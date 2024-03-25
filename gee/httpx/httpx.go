package httpx

import (
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
	"reflect"
)

func Parse(ctx *fasthttp.RequestCtx, v any) error {
	//TODO
	params := make(map[string]any)
	ctx.VisitUserValuesAll(func(k any, v any) {
		params[k.(string)] = v
	})
	if err := mapToPerson(params, v); err != nil {
		return err
	}
	return nil
}

func mapToPerson(m map[string]interface{}, v any) error {
	t := reflect.TypeOf(v).Elem()
	req := reflect.ValueOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		structKey := t.Field(i).Tag.Get("json")
		if val, ok := m[structKey]; ok {
			ToAny(val, req.Field(i))
		}
	}
	return nil
}

func ToAny(i any, value reflect.Value) {
	switch value.Kind().String() {
	case "string":
		value.SetString(cast.ToString(i))
	case "int32", "int", "int8":
		value.SetInt(int64(cast.ToInt(i)))
	case "int64":
		value.SetInt(cast.ToInt64(i))
	case "bool":
		value.SetBool(cast.ToBool(i))
	default:
		panic("type not found")
	}
}
