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

	fmt.Printf("Starting application perform checks now in %v:%v\n", config.Host, config.Port)

	go backgroundTaskSync()

	router := initRouter()
	router.Run(config.Host + ":" + config.Port)

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
		app.CheckRepo(element.URL, true, element.Script)
	}
}

func backgroundTaskSync() {
	fmt.Printf("Running process in backgroud looking for repository changes, configured to interval %v seconds\n", config.RefreshInterval)
	ticker := time.NewTicker(time.Duration(config.RefreshInterval) * time.Second)
	for _ = range ticker.C {
		for _, element := range config.Repositories {
			app.CheckRepo(element.URL, element.AutoSync, element.Script)
		}
	}
}
