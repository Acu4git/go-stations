package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	// TODO: サーバーをlistenする

	// Version 1
	// err = http.ListenAndServe(port, mux)
	// if err != nil {
	// 	return err
	// }

	// Version 2 (Graceful Shutdown用)
	// Station6(Go基礎編): 以下にGraceful Shutdown処理を追加
	//引数に与えた信号が送られたときにctx.Doneチャンネルが閉じられる
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer func() {
		//os.KILLの場合この文が表示されないことがわかる
		fmt.Println("stop func worked.")
		stop()
	}()

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	//サーバーを非同期で実行（Graceful Shutdownを行うため）
	go func() {
		err = srv.ListenAndServe() //http.ErrServerClosed
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	//信号を受け付けるまでブロック
	<-ctx.Done()
	fmt.Println("Done channel is closed.")

	//cancel関数が呼ばれるか，指定した時間が経過するとctx.Doneチャンネルが閉じられる
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		fmt.Println("cancel func worked.")
		cancel()
	}()

	//この関数を呼び出すと新しいリクエストを受け付けなくなる．
	//現在実行しているタスクが完了するか，ctxに指定した時間が経過するとシャットダウンする．
	err = srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
