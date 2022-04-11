package Model

import (
	"log"
	"main/Utils"
	"time"
)

func GetAllOrder(companyId string) ([]Utils.Order, error) {
	var orderInfos []Utils.Order
	var orderInfo Utils.Order
	var landTransportCompanyId, seaTransportCompanyId string
	template := `Select OrderId, StartDate, LandTransportCompanyId, SeaTransportCompanyId, OrderStatus From ShippingTraceability.Orders Where ClientCompanyId = ?`
	rows, err := Utils.DB().Query(template, companyId)
	if err != nil {
		log.Println("[GetAllOrder]数据库异常", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&orderInfo.OrderId, &orderInfo.StartDate, &landTransportCompanyId, &seaTransportCompanyId, &orderInfo.Status)
		orderInfo.LandTransCompanyName, _ = GetCompanyBasicInfo(landTransportCompanyId)
		orderInfo.SeaTransCompanyName, _ = GetCompanyBasicInfo(seaTransportCompanyId)
		orderInfo.ClientCompanyName, _ = GetCompanyBasicInfo(companyId)
		orderInfos = append(orderInfos, orderInfo)
	}
	return orderInfos, nil
}

//GetCompanyBasicInfo 通过Id获取企业的类型以及名称
func GetCompanyBasicInfo(companyId string) (string, string) {
	if companyId == "" {
		return "", ""
	}
	companyName, err := Utils.RDB().Get(companyId + "#name").Result()
	companyType, err := Utils.RDB().Get(companyId + "#type").Result()
	if err != nil { //Redis中没有找到则进行查找
		basicInfo, _ := CompanyBasicInfo(companyId)
		Utils.RDB().Set(companyId+"#name", basicInfo.CompanyType, time.Minute*5)
		Utils.RDB().Set(companyId+"#name", basicInfo.CompanyName, time.Minute*5)
		return basicInfo.CompanyName, basicInfo.CompanyType
	}
	return companyName, companyType
}
