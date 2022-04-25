package MiddleWare

import (
	"github.com/gin-gonic/gin"
	"main/Utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.Set("companyId", int64(1))
			return
		}
		info, err := Utils.ParseToken(token)
		if err != nil {
			c.Set("companyId", int64(1))
			return
		}
		//获取成功则向后传输CompanyId
		c.Set("companyId", info.CompanyId)
		//wg := sync.WaitGroup{} //同步
		//协程执行闭包中的 续杯程序
		//go func(info *Utils.CustomClaims) {
		//	wg.Add(1)
		//	defer wg.Done()
		//	t := time.Unix(info.ExpiresAt, 0)
		//	timeExceed := int(t.Sub(time.Now()).Seconds())
		//	// 当token有效期小于一半时，生成新的token
		//	if timeExceed < Utils.MAXAGE/2 {
		//		c.SetCookie("token", Utils.CreateToken(info.CompanyId), Utils.MAXAGE, "/", "", false, false)
		//	}
		//}(info)

		c.Next() //继续执行后续程序
		//wg.Wait()
		return
		//}
		//c.Set("companyId", )
		//}
		//c.JSON(http.StatusBadRequest, gin.H{"message": "请先登陆"})
		//c.Abort()
		//}
	}
}
