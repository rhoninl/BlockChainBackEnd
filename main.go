package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Controller"
	"main/MiddleWare"
)

const f = `
                  _oo0oo_
                 088888880
                 88" . "88
                 (| -_- |)
                  0\ = /0
               ___/'---'\___
             .' \\\\| | \\\\'.
            / \\\\|||:||| \\\ \
           /_ ||||| -:- ||||| _\
          |   | \\\\\-//////|   |  
          | \_|  ''\---/''  |_/ |
          \  .-\__  '-'  __/-.  /
        ___'. .'  /--.--\  '. .'___
     ."" '<  '.___\_<|>_/___.' >'  "".
    | | : '-  \'.;'\ _ /';.'/ - ' : | |
    \  \ '_.   \_ __\ /__ _/   ._' /  /
====='-.____'.___ \_____/___.-'____.-'=====
`

func main() {
	fmt.Println(f)
	router := gin.Default()
	homePage := router.Group("/homePage")
	{
		homePage.GET("/index", Controller.GetIndex)
	}
	account := router.Group("/account")
	{
		account.POST("/Login", Controller.Login)
		account.POST("/Register", Controller.Register)
		account.POST("/LogOut", MiddleWare.Auth(), Controller.LogOut)
		account.GET("/Info", MiddleWare.Auth(), Controller.Info)
	}

	stuff := router.Group("/stuff")
	{
		stuff.GET("/getStuff", MiddleWare.Auth(), Controller.GetStuff)
		stuff.POST("/addStuff", MiddleWare.Auth(), Controller.AddStuff)
	}
	company := router.Group("/company")
	{
		company.GET("/getJointVenture", MiddleWare.Auth(), Controller.GetJointVenture)
	}
	order := router.Group("/order")
	{
		order.GET("/getAllOrder", MiddleWare.Auth(), Controller.GetAllOrder)
		order.POST("/submitOrder", MiddleWare.Auth(), Controller.BindForm)
	}
	router.GET("/ws", MiddleWare.Auth(), Controller.BuildSocket)
	router.Run(":8080")
}
