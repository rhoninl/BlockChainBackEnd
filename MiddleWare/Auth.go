package MiddleWare

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/Utils"
	"net/http"
	"sync"
	"time"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err == nil {
			// 获取CompanyId
			info, err := Utils.ParseToken(token)
			if err == nil {
				//获取成功则向后传输CompanyId
				c.Set("companyId", info.CompanyId)
				wg := sync.WaitGroup{} //同步
				//协程执行闭包中的 续杯程序
				go func(info *Utils.CustomClaims) {
					wg.Add(1)
					defer wg.Done()
					t := time.Unix(info.ExpiresAt, 0)
					timeExceed := int(t.Sub(time.Now()).Seconds())
					// 当token有效期小于一半时，生成新的token
					if timeExceed < Utils.MAXAGE/2 {
						c.SetCookie("token", Utils.CreateToken(info.CompanyId), Utils.MAXAGE, "/", "", false, false)
					}
				}(info)
				c.Next() //继续执行后续程序
				wg.Wait()
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "请先登陆"})
		c.Abort()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Method", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Origin", "*")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}

}
