package middleware

import (
	"fmt"
	"net/http"
)

// recover機能がないハンドラを引数に渡して，recover機能付きのハンドラに改造して返してもらう
func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovered from panic:", err)
				http.Error(w, "サーバー内でエラーが起きました，ごめんちょ", http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
	//fnをhttp.HandlerFunc型でキャストしている
	//HandlerFunc自体にServerHTTPがあり，その中で元の関数fnを呼び出している
	//構造体を宣言することなくhttp.Handlerを実装できる
}
