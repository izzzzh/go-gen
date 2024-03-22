package httpx

import (
	"github.com/valyala/fasthttp"
	"reflect"
	"strconv"
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
	// 使用反射来动态地设置结构体的字段
	val := reflect.ValueOf(v).Elem()
	t := reflect.TypeOf(v).Elem()
	for key, value := range m {
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Tag.Get("json") == key {
				field := val.Field(i)
				if field.Kind() == reflect.String {
					if strVal, ok := value.(string); ok {
						field.SetString(strVal)
					}
				} else if field.Kind() == reflect.Int {
					parseInt, err := strconv.ParseInt(value.(string), 10, 64)
					if err != nil {
						return err
					}
					field.SetInt(parseInt)
				}
			}
		}
	}

	return nil
}
