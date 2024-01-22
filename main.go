package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	const port = ":8080"
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	fmt.Println("server listen: ", "http://localhost"+port)
	// 在运行例程中初始化服务器
	// 它不会阻止下面的优雅关机处理
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 监听中断信号
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
