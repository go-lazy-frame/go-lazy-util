//             ,%%%%%%%%,
//           ,%%/\%%%%/\%%
//          ,%%%\c "" J/%%%
// %.       %%%%/ o  o \%%%
// `%%.     %%%%    _  |%%%
//  `%%     `%%%%(__Y__)%%'
//  //       ;%%%%`\-/%%%'
// ((       /  `%%%%%%%'
//  \\    .'          |
//   \\  /       \  | |
//    \\/攻城狮保佑) | |
//     \         /_ | |__
//     (___________)))))))                   `\/'
/*
 * 修订记录:
 * long.qian 2021-10-12 10:15 创建
 */

/**
 * @author long.qian
 */

package go_lazy_util

import (
	"encoding/json"
	"fmt"
)

var (
	JsonUtil = new(jsonUtil)
)

type jsonUtil struct {
}

// ParseObjToJsonString 解析对象为 JSON 字符串
func (receiver *jsonUtil) ParseObjToJsonString(obj interface{}) string {
	if obj == nil {
		return ""
	}
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
