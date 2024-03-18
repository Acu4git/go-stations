package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type contextKey interface{}

// 受け取ったRequestからcontextを割り出し，User-Agentの情報を格納する
func SetUserAgent(r *http.Request) context.Context {
	ctx := r.Context()
	//User-Agentの情報を解析して構造体にまとめる
	ua := useragent.Parse(r.UserAgent())

	ctx = context.WithValue(ctx, contextKey("VersionNo"), ua.VersionNo)
	ctx = context.WithValue(ctx, contextKey("OSVersionNo"), ua.OSVersionNo)
	ctx = context.WithValue(ctx, contextKey("URL"), ua.URL)
	ctx = context.WithValue(ctx, contextKey("String"), ua.String)
	ctx = context.WithValue(ctx, contextKey("Name"), ua.Name)
	ctx = context.WithValue(ctx, contextKey("Version"), ua.Version)
	ctx = context.WithValue(ctx, contextKey("OS"), ua.OS)
	ctx = context.WithValue(ctx, contextKey("OSVersion"), ua.OSVersion)
	ctx = context.WithValue(ctx, contextKey("Device"), ua.Device)
	ctx = context.WithValue(ctx, contextKey("Mobile"), ua.Mobile)
	ctx = context.WithValue(ctx, contextKey("Tablet"), ua.Tablet)
	ctx = context.WithValue(ctx, contextKey("Desktop"), ua.Desktop)
	ctx = context.WithValue(ctx, contextKey("Bot"), ua.Bot)

	return ctx
}

// 第1引数にcontext，第2引数に取得したいUser-Agentのフィールド名を文字列で渡す
func GetUserAgent(ctx context.Context, field string) interface{} {
	//contextからvalueを抽出する際にもkeyをcontextKey型でキャストするべきか
	k := contextKey(field)
	return ctx.Value(k)
}
