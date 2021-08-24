package main

import (
	"fmt"
	"net/http"
	"neurone-am-simulator-v2/controller"
	"neurone-am-simulator-v2/hash"
	"neurone-am-simulator-v2/memory"
	"os"
	"time"

	"neurone-am-simulator-v2/middleware"
	"neurone-am-simulator-v2/util"

	"github.com/gin-gonic/gin"
)

func backgroundTask(name string) {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-hash.Channels[name]:
			fmt.Println("stoppingggg")
			time.Sleep(1 * time.Second)
			ticker.Stop()
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t, name)
		}
	}

}

func modify(a int) {

	hash.Channels["s1"] <- true
}
func main() {
	time.Local = time.UTC
	// Cargar Variables de entorno:
	util.LoadEnv()
	memory.Setup()

	app := gin.Default()
	app.Use(middleware.CorsMiddleware())

	base := app.Group("/api/")
	controller.Routes(base)

	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Servicio no encontrado."})
	})

	http.ListenAndServe(os.Getenv("ADDR"), app)
	// app.GET("/init/:id", func(c *gin.Context) {

	// 	id := c.Param("id")
	// 	hash.Channels[id] = make(chan bool)
	// 	go backgroundTask(id)

	// })

	// app.GET("/stop/:id", func(c *gin.Context) {
	// 	id := c.Param("id")
	// 	hash.Channels[id] <- true

	// 	c.JSON(200, "OK")
	// })

	// go modify(1)
	// fmt.Println(<-hash.Channels["s1"])

	// modify(2)
	// fmt.Println(hash.Channels)
	// fmt.Println("Init background job")
	// // done := make(chan bool)
	// // go backgroundTask("s1")
	// go backgroundTask("s2")

	// time.Sleep(5 * time.Second)
	// hash.Channels["s1"] <- true
	// fmt.Println("Background job stopped")
	// time.Sleep(5 * time.Second)
	// hash.Channels["s2"] <- true
	// time.Sleep(5 * time.Second)

}
