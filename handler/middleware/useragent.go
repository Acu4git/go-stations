package middleware

import (
	"context"
	"net/http"
	"reflect"

	"github.com/mileusna/useragent"
)

type contextKey interface{}

// 受け取ったRequestからcontextを割り出し，User-Agentの情報を格納する
func SetUserAgent(r *http.Request) context.Context {
	ctx := r.Context()
	//User-Agentの情報を解析して構造体にまとめる
	ua := useragent.Parse(r.UserAgent())
	//uaの型情報
	uaTypeInfo := reflect.TypeOf(ua)
	//uaの各フィールドの値の情報をもつ
	uaValueInfo := reflect.ValueOf(ua)
	//debug
	// fmt.Println("uaTypeInfo:", uaTypeInfo.Name())
	// fmt.Println("uaValueInfo:", uaValueInfo)
	// fmt.Println("-----------------------------------------------")
	for i := 0; i < uaTypeInfo.NumField(); i++ {
		k := contextKey(uaTypeInfo.Field(i).Name)
		v := uaValueInfo.Field(i)
		ctx = context.WithValue(ctx, k, v)
		//debug
		// fmt.Println("Field", i, "name :", k)
		// fmt.Println("Field", i, "value:", v)
		// fmt.Println("-----------------------------------------------")
	}
	return ctx
}

// 第1引数にcontext，第2引数に取得したいUser-Agentのフィールド名を文字列で渡す
func GetUserAgent(ctx context.Context, field string) interface{} {
	//contextからvalueを抽出する際にもkeyをcontextKey型でキャストするべきか
	k := contextKey(field)
	return ctx.Value(k)
}
