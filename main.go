package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	gin "github.com/gin-gonic/gin"

	app "repoear.io/app"

	model "repoear.io/model"
)

var (
	err    error
	config model.Config
)

func main() {
	//init all envs
	config, err = app.LoadConfiguration()
	if err != nil {
		fmt.Printf("Failed to load configuration %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("Starting application perform checks\n")

	go backgroundTaskSync()

	httpPort := os.Getenv("REPOEAR_HOST_PORT")
	if httpPort == "" {
		httpPort = "8000"
	}
	router := initRouter()
	router.Run(":" + httpPort)

}

func initRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	})

	r.POST("/sync", func(c *gin.Context) {
		go backgroundTaskForceSync()
		c.JSON(http.StatusOK, gin.H{"status": "sync has been trigged"})
		return
	})
	return r
}

func backgroundTaskForceSync() {
	fmt.Println("Forcing sync, the process will run for all elements of config.yml")
	for _, element := range config.Repositories {
		app.CheckRepo(element.URL, true, element.Script, element.Force)
	}
}

func backgroundTaskSync() {
	fmt.Printf("Running process in backgroud looking for repository changes, configured to interval %v seconds\n", config.Refresh)
	ticker := time.NewTicker(time.Duration(config.Refresh) * time.Second)
	for _ = range ticker.C {
		for _, element := range config.Repositories {
			app.CheckRepo(element.URL, element.Sync, element.Script, element.Force)
		}
	}
}
