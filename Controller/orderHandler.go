package Controller

import (
	"github.com/gin-gonic/gin"
	"log"
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
	if !Model.CheckOrderStatus(info.OrderId, "议价") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "当前订单未处于议价状态"})
		return
	}
	if Model.CheckBargainSent(info.OrderId, info.TargetCompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请勿重复发送询价"})
		return
	}
	if !Model.CheckCompanyFriend(companyId.(int64), info.TargetCompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "尚未与目标公司建交"})
		return
	}
	if !Model.SendMessageTo(3, info.OrderId, info.TargetCompanyId, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "消息发送失败"})
		return
	}
	if !Model.AskFroBargain(info.TargetCompanyId, info.OrderId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "发送失败"})
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

func ReplyBargain(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var info Utils.ReplyBargain
	c.Bind(&info)
	if !Model.CheckOrderCanBargain(info.OrderId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "订单状态不对"})
		return
	}
	if !Model.ReplyBargain(info, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "回复失败"})
		return
	}
	clientCompanyId, err := Model.GetOrderClientId(info.OrderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "订单信息获取失败"})
		return
	}
	fromCompanyName, _ := Model.GetCompanyBasicInfo(companyId.(int64))
	var str string
	if info.IsPass {
		str = `您的订单(id:` + strconv.FormatInt(info.OrderId, 10) + `)得到了新的报价
	` + fromCompanyName + `(id : ` + strconv.FormatInt(companyId.(int64), 10) + `)` + `向您报价` + strconv.FormatInt(info.Bargain, 10) +
			`元`
	} else {
		str = `您的订单(id:` + strconv.FormatInt(info.OrderId, 10) + `)的报价请求被` + fromCompanyName + `(id : ` + strconv.FormatInt(companyId.(int64), 10) + `)拒绝`
	}
	go Model.SendMessageTo(4, str, clientCompanyId, 0)
	c.JSON(http.StatusCreated, nil)
}

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

func SubmitCompanyChoose(c *gin.Context) {
	iCompanyId, _ := c.Get("companyId")
	companyId := iCompanyId.(int64)
	var form Utils.OrderCompany
	c.Bind(&form)
	if !Model.CheckOrderCompany(form.OrderId, companyId) {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "当前帐号与该订单所属帐号不符"})
		return
	}
	if !Model.CheckCompanyFriend(companyId, form.SeaCompanyId) || !Model.CheckCompanyFriend(companyId, form.LandCompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "当前帐号与陆代公司或者船代公司不是好友"})
		return
	}
	_, landType := Model.GetCompanyBasicInfo(form.LandCompanyId)
	_, seaType := Model.GetCompanyBasicInfo(form.SeaCompanyId)
	if landType != "陆运公司" || seaType != "船代" {
		log.Println(landType, seaType)
		c.JSON(http.StatusBadRequest, gin.H{"message": "未选择船代或者陆代"})
		return
	}
	if !Model.UpdateOrderAgent(form) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusCreated, nil)
}
