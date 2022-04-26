package MiddleWare

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Utils"
	"net/http"
	"sync"
	"time"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		var companyId = int64(1)
		if err != nil {
			fmt.Println("有个小比崽子没登陆")
		} else {
			info, err := Utils.ParseToken(token)
			if err == nil {
				companyId = info.CompanyId
			}
		}
		c.Set("companyId", companyId)
	}
}

func Auth2() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err == nil {
			info, err := Utils.ParseToken(token)
			if err == nil {
				c.Set("companyId", info.CompanyId)
				wg := sync.WaitGroup{}
				//创建协程cookie续杯
				go func(info *Utils.CustomClaims) {
					wg.Add(1)
					defer wg.Done()
					t := time.Unix(info.ExpiresAt, 0)
					timeExceed := int(t.Sub(time.Now()).Seconds())
					if timeExceed < Utils.MAXAGE/2 {
						newToken := Utils.CreateToken(info.CompanyId)
						http.SetCookie(c.Writer, &http.Cookie{
							Name:     "token",
							Value:    newToken,
							Path:     "/",
							Domain:   c.Request.Header.Get("Origin"),
							MaxAge:   Utils.MAXAGE,
							Secure:   true,
							SameSite: 4,
						})
					}
				}(info)
				c.Next()
				wg.Wait()
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "请先登陆"})
	}
}
