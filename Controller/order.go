package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"main/Utils"
	"net/http"
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

func BindForm(c *gin.Context) {
	var orders Utils.OrderInfo
	c.Bind(&orders)
	id, ok, err := Model.RecordOrder(orders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求异常"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orderId": id})
}
