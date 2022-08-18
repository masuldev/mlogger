package mlogger

import (
	"fmt"
	"github.com/goccy/go-json"
	"reflect"
)

func typeCheck(m interface{}) reflect.Kind {
	return reflect.ValueOf(m).Kind()
}

func AppendKey() {

}

func Write(m interface{}) string {
	switch typeCheck(m) {
	case reflect.Struct:
		//var data map[string]string

		mm, _ := json.Marshal(m)
		return string(mm)
		//json.Unmarshal(mm, &data)
		//
		//var b bytes.Buffer
		//for key, value := range data {
		//	b.WriteString("{ ")
		//	b.WriteString(key)
		//	b.WriteString(": ")
		//	b.WriteString(value)
		//	b.WriteString(" }")
		//}
		//
		//return b.String()
	default:
		return fmt.Sprintf("%s", m)
	}
}
