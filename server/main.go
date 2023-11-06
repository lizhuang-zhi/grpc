package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"server/common"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Use CORS middleware with explicit origin
	config := cors.Config{
		AllowOrigins:     []string{"http://124.221.157.89:8080"}, // Add your frontend URL here
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}
	router.Use(cors.New(config))

	// 声明一个通配符路由，将所有请求都指向前端的入口页面
	router.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	router.GET("/query", func(c *gin.Context) {
		// 休眠
		time.Sleep(time.Second * 3)

		user, ok := c.GetQuery("user")
		if !ok {
			user = "Guest"
		}
		c.JSON(200, gin.H{
			"message": "Hello world!",
			"data":    user,
		})
	})

	router.GET("/list", func(c *gin.Context) {
		id, ok := c.GetQuery("id")
		if !ok {
			id = "no id"
		}
		c.JSON(200, gin.H{
			"params": id,
			"data":   common.DataList(),
		})
	})

	router.POST("/form", func(c *gin.Context) {
		var formReqParam struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Hobby string `json:"hobby"`
		}

		err := c.BindJSON(&formReqParam)
		if err != nil {
			c.JSON(200, gin.H{
				"message": err.Error(),
				"data":    formReqParam,
			})
		}

		// 休眠
		time.Sleep(time.Second * 1)

		c.JSON(200, gin.H{
			"message": "success",
			"data":    formReqParam,
		})
	})

	// 接收数据上报
	router.POST("/report", func(c *gin.Context) {
		var reportParam struct {
			Title string      `json:"name"`
			Msg   string      `json:"msg"`
			Data  interface{} `json:"data"`
		}

		err := c.BindJSON(&reportParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "无效的 JSON 数据"})
			return
		}

		// 存储到文件
		err = saveToFile(reportParam)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": "上报成功",
			"data":    reportParam,
		})
	})

	router.GET("/stop", func(c *gin.Context) {
		// 终止服务
		os.Exit(0)
	})

	// 注意，这里不再使用 router.Run("0.0.0.0:12222")，而是使用 http.ListenAndServe
	// 这是因为我们将 NoRoute 处理交给了前端，不再由 Gin 处理
	if err := http.ListenAndServe(":12222", router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func saveToFile(data interface{}) error {
	// 文件路径
	filePath := "./record.txt"

	// 如果文件不存在，创建文件
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("无法创建文件: %v", err)
		}
	}

	// 将数据转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("无法转换为 JSON 字符串: %v", err)
	}

	// 将 JSON 字符串写入文件
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("无法写入文件: %v", err)
	}

	return nil
}
