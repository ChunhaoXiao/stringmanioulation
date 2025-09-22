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
	Fields         []string
	NewName        string
	Connector      string //string for connecting fields
	DeleteConcated bool   // if keep the fields being cancated
}

type Substring struct {
	Field      string
	StartIndex int
	Lenght     int
}

// 自定义插件配置
type MyConfig struct {
	debug             bool
	RemoveLeadingZero RemoveLeadingZero
	PaddingZeros      []PaddingZero
	ConcatFields      []ConcatFields
	Substrings        []Substring
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

		if action == "substring" {
			content := item.Get("fields")
			if content.Exists() {
				var substrings []Substring
				content.ForEach(func(key, value gjson.Result) bool {
					fmt.Println("value forsubstring###################3", value)
					fieldName := value.Get("name").String()
					startIndex := value.Get("startIndex").Int()
					length := value.Get("length").Int()
					substring := Substring{
						Field:      fieldName,
						StartIndex: int(startIndex),
						Lenght:     int(length),
					}
					substrings = append(substrings, substring)
					return true
				})
				config.Substrings = substrings
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
						NewName:        value.Get("newName").String(),
						Fields:         names,
						Connector:      value.Get("connector").String(),
						DeleteConcated: value.Get("deleteConcated").Bool(),
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
	ccc := concatField2(result, config.ConcatFields)

	fmt.Println("ccc#########################", ccc)
	return types.ActionContinue
}

/*
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
		//concatField(m, []string{}, "fullName", false)
	}
}*/

// func concatField(datas map[string]interface{}, fields []string, newName string, keepOld bool) map[string]interface{} {
// 	contactedValue := ""
// 	for _, val := range fields {
// 		value, ok := datas[val]
// 		if !ok {
// 			return datas
// 		}
// 		contactedValue = contactedValue + fmt.Sprintf("%v", value)
// 	}
// 	datas[newName] = contactedValue
// 	if keepOld == false {
// 		for _, key := range fields {
// 			delete(datas, key)
// 		}
// 	}
// 	return datas

// }

func concatField2(datas map[string]interface{}, concat []ConcatFields) map[string]interface{} {
	for _, items := range concat {
		if strings.Contains(items.NewName, ".") {
			fieldNameArr := strings.Split(items.NewName, ".")
			concatedValue := []string{}
			connector := items.Connector
			for _, fieldName := range items.Fields {
				accessor, err := dotaccess.NewAccessorDot[string](&datas, fieldName)
				if err == nil {
					concatedValue = append(concatedValue, accessor.Get())
					if items.DeleteConcated == true {
						delete(datas[fieldNameArr[0]].(map[string]interface{}), strings.Split(fieldName, ".")[1])
					}
				}
			}
			datas[fieldNameArr[0]].(map[string]interface{})[fieldNameArr[1]] = strings.Join(concatedValue, connector)
		} else {
			concatedValue := []string{}
			for _, fieldName := range items.Fields {
				value, ok := datas[fieldName]
				if ok {
					concatedValue = append(concatedValue, fmt.Sprintf("%v", value))
					if items.DeleteConcated == true {
						delete(datas, fieldName)
					}
				}
			}
			datas[items.NewName] = strings.Join(concatedValue, items.Connector)
		}
	}

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

func substrings(datas map[string]interface{}, substrings []Substring) map[string]interface{} {
	fmt.Println("substring------------------------------>", substrings)
	for _, substring := range substrings {
		field := substring.Field
		if strings.Contains(field, ".") {
			fmt.Println("###################################", substring.StartIndex)
			accessor, err := dotaccess.NewAccessorDot[string](&datas, field)
			if err == nil {
				value := accessor.Get()
				fmt.Println("value is:", value)
				fmt.Println("startINdex####", substring.StartIndex)
				substr := value[substring.StartIndex : substring.StartIndex+substring.Lenght]
				accessor.Set(substr)
			}
		} else {
			value, ok := datas[field]
			if ok {
				strValue := fmt.Sprintf("%v", value)
				datas[field] = strValue[substring.StartIndex : substring.StartIndex+substring.Lenght]
			}
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
