package MiddleWare

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		//if err == nil {
		//获取CompanyId
		info, err := Utils.ParseToken(token)
		var companyId = int64(1)
		if token != "" && err == nil {
			companyId = info.CompanyId
		} else {
			fmt.Println("有个小比崽子没登陆")
		}
		//获取成功则向后传输CompanyId
		c.Set("companyId", companyId)
		//wg := sync.WaitGroup{} //同步
		//协程执行闭包中的 续杯程序
		//go func(info *Utils.CustomClaims) {
		//	wg.Add(1)
		//	defer wg.Done()
		//	t := time.Unix(info.ExpiresAt, 0)
		//	timeExceed := int(t.Sub(time.Now()).Seconds())
		//	// 当token有效期小于一半时，生成新的token
		//	if timeExceed < Utils.MAXAGE/2 {
		//		//c.SetCookie("token", Utils.CreateToken(info.CompanyId), Utils.MAXAGE, "/", "", false, false)
		//		//http.SetCookie(c.Writer, &http.Cookie{})
		//	}
		//}(info)
		c.Next() //继续执行后续程序
		//wg.Wait()
		return
		//}
	}
}
