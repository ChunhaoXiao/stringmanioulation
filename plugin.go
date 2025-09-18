package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

// 自定义插件配置
type MyConfig struct {
	mockEnable bool
	actions    []map[string]gjson.Result
}

// 在控制台插件配置中填写的yaml配置会自动转换为json，此处直接从json这个参数里解析配置即可
func parseConfig(json gjson.Result, config *MyConfig, log logs.Log) error {
	// 解析出配置，更新到config中
	config.mockEnable = json.Get("mockEnable").Bool()
	ac := json.Get("actions")
	var actions []map[string]gjson.Result
	for _, item := range ac.Array() {
		actions = append(actions, item.Map())
	}
	config.actions = actions
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
	if config.mockEnable {
		proxywasm.SendHttpResponse(200, nil, []byte("hello world"), -1)
	}
	return types.HeaderContinue
}

func onHttpRequestBody(ctx wrapper.HttpContext, config MyConfig, body []byte) types.Action {
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	// for k, v := range result {
	// 	//fmt.Println(k, "===>", v)
	// }
	updateNestedMap(result, "name", false)

	fmt.Println("updated result######", result)
	for _, actionType := range config.actions {
		fmt.Println(actionType["type"])
		fmt.Println("fields", actionType["field"])
	}

	concatField(result, []string{"firstName", "lastName"}, "fullName", false)
	fmt.Println("after  concat result", result)

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

func paddingZero(datas map[string]interface{}, fields []string, zeroCount int) map[string]interface{} {
	for _, val := range fields {
		value, ok := datas[val]
		if ok {
			strval := fmt.Sprintf("%v", value)
			length := len(strval) + zeroCount
			floatValue, err := strconv.ParseFloat(strval, 64)
			if err == nil {
				s, _ := fmt.Printf("%*d", length, floatValue)
				datas[val] = s
			} else {
				paddedStr := strings.Repeat("0", zeroCount) + strval
				datas[val] = paddedStr
			}
		}
	}
	return datas
}
