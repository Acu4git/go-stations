package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func checkAuth(r *http.Request) bool {
	uid, pwd, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return uid == os.Getenv("BASIC_AUTH_USER_ID") && pwd == os.Getenv("BASIC_AUTH_PASSWORD")
}

func BasicAuthMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(r) {
			// WWW-Authenticate	  => サーバが提供する認証方式等の情報を示すヘッダ名
			// Basic 			  => 認証方式
			// realm			  => リソースが属するURI空間の名前

			//疑問：ルーティングごとにrealmの名前を指定する必要がある？(if文などで)
			w.Header().Add("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		_, err := fmt.Println("Successful Basic Authentication")
		if err != nil {
			log.Fatal(err)
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
