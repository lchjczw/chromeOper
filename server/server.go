package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lchjczw/chromeOper/env"
	"server/config"
	"server/flow"
)

func main() {

	ctx, cancel := env.NewChrome(config.Global())
	defer cancel()

	router := gin.Default()

	router.GET("/run", func(c *gin.Context) {
		flow.RunChromeDp(c, ctx)
	})

	router.GET("/login", func(c *gin.Context) {
		userName := c.Query("user_name")
		password := c.Query("password")
		r := flow.Login(ctx, userName, password)
		c.JSON(200, gin.H{
			"code":    200,
			"message": r,
		})
	})

	// 默认启动的是 8080端口，也可以自己定义启动端口
	router.Run()
	// router.Run(":3000") for a hard coded port
}

