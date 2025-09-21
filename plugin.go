package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/claytonsingh/golib/dotaccess"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"

	logs "github.com/higress-group/wasm-go/pkg/log"
	"github.com/higress-group/wasm-go/pkg/wrapper"
	"github.com/tidwall/gjson"
)

func main() {

	// wrapper.SetCtx(
	// 	// 插件名称
	// 	"my-pluginsss",
	// 	// 为解析插件配置，设置自定义函数
	// 	wrapper.ParseConfigBy(parseConfig),
	// 	// 为处理请求头，设置自定义函数
	// 	wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
	// 	wrapper.ProcessRequestBody(onHttpRequestBody),
	// )
}

func init() {
	wrapper.SetCtx(
		// 插件名称
		"my-plugin",
		// 为解析插件配置，设置自定义函数
		wrapper.ParseConfigBy(parseConfig),
		// 为处理请求头，设置自定义函数
		wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
		wrapper.ProcessRequestBody(onHttpRequestBody),
	)
}

type RemoveLeadingZero struct {
	Fields    []string
	Recursive bool
}

type PaddingZero struct {
	Field     string
	ZeroCount int
	Recursive bool
}

type ConcatFields struct {
	// Fields       []string
	// NewFieldName string
	// Recursive    bool
	// KeepFields   bool
	Fields  []string
	NewName string
}

// 自定义插件配置
type MyConfig struct {
	debug             bool
	RemoveLeadingZero RemoveLeadingZero
	PaddingZeros      []PaddingZero
	ConcatFields      []ConcatFields
}

// 在控制台插件配置中填写的yaml配置会自动转换为json，此处直接从json这个参数里解析配置即可
func parseConfig(json gjson.Result, config *MyConfig, log logs.Log) error {
	// 解析出配置，更新到config中
	//config.mockEnable = json.Get("mockEnable").Bool()
	ac := json.Get("actions")
	//var actions []map[string]gjson.Result
	//myconf := MyConfig{}
	for _, item := range ac.Array() {
		action := item.Get("type").String()
		if action == "removeLeadingZero" {
			rmLeadingZero := RemoveLeadingZero{
				Recursive: item.Get("recursive").Bool(),
			}
			if fields := item.Get("fields"); fields.Exists() && fields.IsArray() {
				fields.ForEach(func(_, value gjson.Result) bool {
					rmLeadingZero.Fields = append(rmLeadingZero.Fields, value.String())
					return true
				})
			}
			config.RemoveLeadingZero = rmLeadingZero
		}

		if action == "paddingZero" {
			fieldMap := item.Get("fields")
			if fieldMap.Exists() {
				var paddingZeros []PaddingZero
				fieldMap.ForEach(func(key, value gjson.Result) bool {
					paddingZeros = append(paddingZeros, PaddingZero{Field: key.String(), ZeroCount: int(value.Int())})
					return true
				})
				config.PaddingZeros = paddingZeros
			}
		}

		if action == "concat" {
			content := item.Get("concatContent")
			if content.Exists() {
				var concats []ConcatFields
				content.ForEach(func(key, value gjson.Result) bool {
					var names []string
					fields := value.Get("concatFields")
					if fields.Exists() {
						fields.ForEach(func(key, value gjson.Result) bool {
							names = append(names, value.String())
							return true
						})
					}
					concat := ConcatFields{
						NewName: value.Get("newName").String(),
						Fields:  names,
					}
					concats = append(concats, concat)
					return true
				})
				config.ConcatFields = concats
			}

		}
		//actions = append(actions, item.Map())
	}

	fmt.Println("cfg=========================>", config)

	//fmt.Println("actions============================>", actions)
	//config.actions = actions
	/*
		action1 := map[string]any{
			"field": []string{"name", "email"},
			"type":  "removeLeadingZero",
		}
		action2 := map[string]any{
			"field": map[string]int{
				"username": 20,
				"company":  30,
			},
			"type": "paddingZero",
		}
		action3 := map[string]any{
			"field":   []string{"firstName", "lastName"},
			"type":    "contact",
			"newName": "fullName",
		}
		config.actions = []map[string]any{
			action1, action2, action3,
		}*/
	return nil
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config MyConfig, log logs.Log) types.Action {
	proxywasm.AddHttpRequestHeader("hello", "world")
	// if config.mockEnable {
	// 	proxywasm.SendHttpResponse(200, nil, []byte("hello world"), -1)
	// }
	return types.HeaderContinue
}

func onHttpRequestBody(ctx wrapper.HttpContext, config MyConfig, body []byte) types.Action {
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	// for k, v := range result {
	// 	//fmt.Println(k, "===>", v)
	// }
	//updateNestedMap(result, "name", false)

	/*
		fmt.Println("updated result######", result)
		for _, actionType := range config.actions {
			fmt.Println(actionType["type"])
			fmt.Println("fields", actionType["field"])
		}

		//concatField(result, []string{"firstName", "lastName"}, "fullName", false)
		//fmt.Println("after  concat result", result)
		paddingAction := config.actions[0]
		for field, count := range paddingAction {

		}
		//paddingZero(result, )*/
	// pzr := paddingZero(result, config.PaddingZeros)
	// fmt.Println("padding Zero result########################", pzr)

	//concatField2(result)
	ccc := trimLeadingZeros(result, config.RemoveLeadingZero) //paddingZero(result, config.PaddingZeros) //concatField2(result, config)

	fmt.Println("ccc#########################", ccc)
	return types.ActionContinue
}

func updateNestedMap(m map[string]interface{}, targetKey string, recursive bool) {
	for key, val := range m {
		// if key == targetKey {
		// 	m[key] = newValue
		// 	return // Value updated, no need to go deeper for this key
		// }

		// If the value is another map, recurse
		if recursive {
			if nestedMap, ok := val.(map[string]interface{}); ok {
				updateNestedMap(nestedMap, targetKey, recursive)
			} else {
				if targetKey == key {
					m[key] = "asss"
				}

			}
		} else {
			if targetKey == key {
				m[key] = "asss"
			}

		}
		concatField(m, []string{}, "fullName", false)
	}
}

func concatField(datas map[string]interface{}, fields []string, newName string, keepOld bool) map[string]interface{} {
	contactedValue := ""
	for _, val := range fields {
		value, ok := datas[val]
		if !ok {
			return datas
		}
		contactedValue = contactedValue + fmt.Sprintf("%v", value)
	}
	datas[newName] = contactedValue
	if keepOld == false {
		for _, key := range fields {
			delete(datas, key)
		}
	}
	return datas

}

func concatField2(datas map[string]interface{}, config MyConfig) map[string]interface{} {
	fmt.Println("input data###############", datas)
	concat := config.ConcatFields
	p := &datas
	for _, items := range concat {
		//.Println("item=====================>", items.Fields)
		concatedValue := ""
		//newNames := strings.Split(items.NewName, ".")
		if strings.Contains(items.NewName, ".") {
			fieldNameArr := strings.Split(items.NewName, ".")
			datas[fieldNameArr[0]] = map[string]interface{}{
				fieldNameArr[1]: nil,
			}
		}

		for _, fieldName := range items.Fields {
			fmt.Println("fname^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^", fieldName)

			if strings.Contains(fieldName, ".") {
				accessor, errs := dotaccess.NewAccessorDot[any](p, fieldName)
				if errs != nil {
					return datas
				}
				concatedValue = concatedValue + fmt.Sprintf("%v", accessor.Get())

				//newAccessor.Set(concatedValue)
				//datas[items.NewName] = concatedValue
				//fmt.Println("accessor$$$$$$$$$$$$$$$$$$$$$$$$$", nameAccessor)
				//fmt.Println("get field value==================>:", nameAccessor.Get())

			} else {
				value, ok := datas[fieldName]
				if !ok {
					return datas
				}
				concatedValue = concatedValue + fmt.Sprintf("%v", value)
				datas[items.NewName] = concatedValue
			}

		}

	}
	/*
		for key, val := range datas {
			if contactField.Recursive {
				if nestedMap, ok := val.(map[string]interface{}); ok {
					contactedValue := ""
					for _, val := range contactField.Fields {
						value, _ := datas[val]
						// if !ok {
						// 	return datas
						// }
						contactedValue = contactedValue + fmt.Sprintf("%v", value)
					}
					datas[key] = contactedValue
					concatField2(nestedMap, contactField)
				} else {
					// if targetKey == key {
					// 	m[key] = "asss"
					// }
					datas[key] = val

				}
			} else {

				contactedValue := ""
				for _, val := range contactField.Fields {
					value, _ := datas[val]
					contactedValue = contactedValue + fmt.Sprintf("%v", value)
				}
				datas[key] = contactedValue

			}
			//concatField(m, []string{}, "fullName", false)
		} */
	/*
		contactedValue := ""
		for _, val := range fields {
			value, ok := datas[val]
			if !ok {
				return datas
			}
			contactedValue = contactedValue + fmt.Sprintf("%v", value)
		}
		datas[newName] = contactedValue
		if keepOld == false {
			for _, key := range fields {
				delete(datas, key)
			}
		}*/
	return datas

}

func trimLeadingZeros(datas map[string]interface{}, removeLeadingZero RemoveLeadingZero) map[string]interface{} {

	for _, field := range removeLeadingZero.Fields {
		fmt.Println("field================>", field)

		if strings.Contains(field, ".") {
			accessor, _ := dotaccess.NewAccessorDot[string](&datas, field)
			fieldValue := accessor.Get()
			val := removeLeadingZeros(fieldValue)
			accessor.Set(val)
		} else {
			val := removeLeadingZeros(fmt.Sprintf("%v", datas[field]))
			datas[field] = val
		}
	}
	return datas

}

func paddingZero(datas map[string]interface{}, fields []PaddingZero) map[string]interface{} {

	for _, obj := range fields {
		value, _ := datas[obj.Field]
		if strings.Contains(obj.Field, ".") {
			accessor, _ := dotaccess.NewAccessorDot[string](&datas, obj.Field)
			value = accessor.Get()
		}
		strval := fmt.Sprintf("%v", value)
		paddedStr := strings.Repeat("0", obj.ZeroCount) + strval
		if strings.Contains(obj.Field, ".") {
			accessor, _ := dotaccess.NewAccessorDot[any](&datas, obj.Field)
			accessor.Set(paddedStr)

		} else {
			datas[obj.Field] = paddedStr
		}
	}
	return datas
}

func removeLeadingZeros(s string) string {
	if s == "" {
		return s
	}
	// Keep at least one digit
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}
	if i == len(s) { // all zeros
		return "0"
	}
	return s[i:]
}
