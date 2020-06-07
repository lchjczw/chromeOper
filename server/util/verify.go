package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"
)

func RandNum(num int) string {
	randNum := ""
	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		n := rand.Intn(9)
		randNum = randNum + fmt.Sprintf("%d", n)
	}

	return randNum
}

func ShowVerify() {
	addr := "127.0.0.1:" + RandNum(4)
	go Listen(addr, "img")
	Verify("http://"+addr+"/img/verify.png", 250, 75)
}

func InputVerify() string {
	fmt.Printf("请输入验证码：")
	//fmt.Println("请输入验证码：")
	verify := ""
	fmt.Scanln(&verify)
	return verify
}

func Listen(addr, staticDir string) {
	gin.SetMode(gin.ReleaseMode)

	// 创建记录日志的文件
	f, _ := os.Create("log.bat")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.Static("/"+staticDir, "./"+staticDir)
	r.GET("/verify", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "ok",
		})
	})
	r.Run(addr)
}

func Verify(verifyUrl string, width, height int) bool {
	tpl := `
	<html>
		<head><title>verify</title></head>
		<img src="%s" alt="verify" width="%d" height="%d">
	</html>
	`
	tpl = fmt.Sprintf(tpl, verifyUrl, width, height)

	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("data:text/html,"+url.PathEscape(tpl), "", width+40, height+40)
	if err != nil {
		log.Fatal(err)
	}
	//time.Sleep(time.Second*8)
	// Wait until UI window is closed
	go func() {
		defer ui.Close()
		<-ui.Done()
	}()

	return true
}
