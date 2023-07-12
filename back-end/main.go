package main

import (
        "github.com/gin-gonic/gin"
        "net/http"
        "fmt"
        "os"
)

func main() {
        // 获取主机名
        hostname, err := os.Hostname()
        if err != nil {
            fmt.Println("无法获取主机名：", err)
            return
        }

        r := gin.Default() // 初始化gin

        // /api/v1/version
        r.GET("/api/v1/version", func(ctx *gin.Context) {
                ctx.JSON(http.StatusOK, gin.H{
                        "message": "Current API version: v1, Current hostname is: " + hostname,
                })
        })
        
        // /api/v2/version
        r.GET("/api/v2/version", func(ctx *gin.Context) {
                ctx.JSON(http.StatusOK, gin.H{
                        "message": "Current API version: v2, Current hostname is: " + hostname,
                })
        })
        
        // /api/v3/version
        r.GET("/api/v3/version", func(ctx *gin.Context) {
                ctx.JSON(http.StatusOK, gin.H{
                        "message": "Current API version: v3, Current hostname is: " + hostname,
                })
        })
        
        r.Run() // 启动
}
