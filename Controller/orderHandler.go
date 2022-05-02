package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"main/Utils"
	"net/http"
	"strconv"
)

func GetAllOrder(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	orderInfo, err := Model.GetAllOrder(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, orderInfo)
}

func BindOrder(c *gin.Context) {
	var orders Utils.OrderInfo
	c.Bind(&orders)
	companyId, _ := c.Get("companyId")
	orders.ClientCompanyId = companyId.(int64)
	id, ok, err := Model.RecordOrder(orders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求异常"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"orderId": id})
}

func AskForPrice(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var info Utils.AskPrice
	c.Bind(&info)
	if !Model.CheckOrderCompany(info.OrderId, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该订单不属于当前公司"})
		return
	}
	if !Model.CheckCompanyFriend(companyId.(int64), info.TargetCompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "尚未与目标公司建交"})
		return
	}
	text := `123`
	if !Model.SendMessageTo(3, text, info.TargetCompanyId, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "消息发送失败"})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func GetAllBargain(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	sOrderId := c.Query("orderId")
	orderId, err := strconv.ParseInt(sOrderId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求错误"})
		return
	}
	if !Model.CheckOrderCompany(orderId, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该订单不属于当前公司"})
		return
	}
	info, err := Model.GetCompanyBargain(orderId, companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	for n, item := range info {
		info[n].CompanyName, info[n].CompanyType = Model.GetCompanyBasicInfo(item.CompanyId)
	}
	c.JSON(http.StatusOK, info)
}

//func BargainReply(c *gin.Context) {
//	companyId, _ := c.Get("companyId")
//	c.Bind()
//	Model.CheckMessageAuth()
//	c.JSON(http.StatusCreated, nil)
//}

func GetOrderInfo(c *gin.Context) {
	sOrderId := c.Query("orderId")
	orderId, err := strconv.ParseInt(sOrderId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求异常"})
		return
	}
	info, err := Model.GetOrderInfo(orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, info)
}
