package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Controller"
	"main/MiddleWare"
)

const f = `
		⣿ ⣿ ⣿ ⣿ ⣿ ⠟ ⠋ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⢁ ⠈ ⢻ ⢿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⠃ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠈ ⡀ ⠭ ⢿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⡟ ⠄ ⢀ ⣾ ⣿ ⣿ ⣿ ⣷ ⣶ ⣿ ⣷ ⣶ ⣶ ⡆ ⠄ ⠄ ⠄ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⡇ ⢀ ⣼ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣧ ⠄ ⠄ ⢸ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣇ ⣼ ⣿ ⣿ ⠿ ⠶ ⠙ ⣿ ⡟ ⠡ ⣴ ⣿ ⣽ ⣿ ⣧ ⠄ ⢸ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣾ ⣿ ⣿ ⣟ ⣭ ⣾ ⣿ ⣷ ⣶ ⣶ ⣴ ⣶ ⣿ ⣿ ⢄ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⡟ ⣩ ⣿ ⣿ ⣿ ⡏ ⢻ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣹ ⡋ ⠘ ⠷ ⣦ ⣀ ⣠ ⡶ ⠁ ⠈ ⠁ ⠄ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣍ ⠃ ⣴ ⣶ ⡔ ⠒ ⠄ ⣠ ⢀ ⠄ ⠄ ⠄ ⡨ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣦ ⡘ ⠿ ⣷ ⣿ ⠿ ⠟ ⠃ ⠄ ⠄ ⣠ ⡇ ⠈ ⠻ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⡿ ⠟ ⠋ ⢁ ⣷ ⣠ ⠄ ⠄ ⠄ ⠄ ⣀ ⣠ ⣾ ⡟ ⠄ ⠄ ⠄ ⠄ ⠉ ⠙ ⠻
		⡿ ⠟ ⠋ ⠁ ⠄ ⠄ ⠄ ⢸ ⣿ ⣿ ⡯ ⢓ ⣴ ⣾ ⣿ ⣿ ⡟ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄
		⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⣿ ⡟ ⣷ ⠄ ⠹ ⣿ ⣿ ⣿ ⡿ ⠁ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄
`

func main() {
	fmt.Println(f)
	router := gin.Default()
	router.Use(MiddleWare.MyCors())
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteString("看nmlgcb")
	})
	homePage := router.Group("/homePage")
	{
		homePage.GET("/index", Controller.GetIndex)
	}
	account := router.Group("/account")
	{
		account.POST("/AuthCode", Controller.GetAuth)
		account.POST("/Login", Controller.Login)
		account.POST("/Register", Controller.Register)
		account.POST("/LogOut", MiddleWare.Auth(), Controller.LogOut)
		account.GET("/Info", MiddleWare.Auth(), Controller.Info)
		account.POST("/EditInfo", MiddleWare.Auth(), Controller.EditInfo)
		account.POST("/Forget", Controller.ForgetPassword)
		account.POST("/ChangePd", MiddleWare.Auth(), Controller.ChangePassword)
	}
	stuff := router.Group("/stuff", MiddleWare.Auth())
	{
		stuff.GET("/getStuff", Controller.GetStuff)
		stuff.POST("/addStuff", Controller.AddStuff)
	}
	company := router.Group("/company", MiddleWare.Auth())
	{
		company.GET("/getJointVenture", Controller.GetJointVenture)
		company.POST("/makeFriends", Controller.MakeFriend)
	}
	order := router.Group("/order", MiddleWare.Auth())
	{
		order.GET("/getAllOrder", Controller.GetAllOrder)
		order.POST("/submitOrder", Controller.BindForm)
	}
	message := router.Group("/message", MiddleWare.Auth())
	{
		message.GET("/getMessage", Controller.GetMessage)
		message.GET("/getMessageInfo/:id", Controller.GetMessageInfo)
		message.GET("/deleteMessage/:messageId", Controller.DeleteMessage)
	}
	router.GET("/ws", MiddleWare.Auth(), Controller.BuildSocket)
	//router.RunTLS(":8081", "./Utils/cunyuqing.online_bundle.pem", "./Utils/cunyuqing.online.key")
	router.RunTLS(":8081", "./Utils/cunyuqing.online_bundle.pem", "./Utils/cunyuqing.online.key")
}
