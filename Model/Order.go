package Model

import (
	"fmt"
	"log"
	"main/Utils"
	"sync"
)

func GetAllOrder(companyId int64) ([]Utils.Order, error) {
	var orderInfos []Utils.Order
	var orderInfo Utils.Order
	var landTransportCompanyId, seaTransportCompanyId int64
	template := `Select OrderId, StartDate, LandTransportCompanyId, SeaTransportCompanyId, OrderStatus From Orders Where ClientCompanyId = ?`
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

func RecordOrder(info Utils.OrderInfo) (int64, bool, error) {
	affair, err := Utils.DB().Begin()
	if err != nil {
		return 0, false, err
	}
	template := `Insert Into Orders Set ClientCompanyId = ?,StartDate = now() ,OrderStatus = ?`
	rows, err := affair.Exec(template, info.ClientCompanyId, "新增")
	if err != nil {
		log.Println("[RecordOrder]Orders出错了", err)
		affair.Rollback()
		return 0, false, err
	}
	info.OrderId, _ = rows.LastInsertId()
	wg := sync.WaitGroup{}
	template = `Insert Into Cargo Set CargoName = ? , CargoModel = ? , CargoSize = ? , CargoNum = ? , Category = ? , Weight = ? `
	for _, item := range info.Cargos {
		wg.Add(1)
		rows, err := affair.Exec(template, item.CargoName, item.CargoModel, item.CargoSize, item.CargoNum, item.Category, item.CargoWeight)
		if err != nil {
			log.Println("[RecordOrder]Cargo出错了", err)
			affair.Rollback()
			return 0, false, nil
		}
		cargoId, _ := rows.LastInsertId()
		go func(cargoId int64) {
			template := `Insert Into Order_Cargo Set OrderId = ?,CargoId = ?`
			affair.Exec(template, info.OrderId, cargoId)
			wg.Done()
		}(cargoId)
	}
	template = `Insert Into Address Set Country = ?,City = ?,Address = ?`
	rows, err = affair.Exec(template, info.SendAddress.Country, info.SendAddress.City, info.SendAddress.Address)
	sendAddressId, _ := rows.LastInsertId()
	rows, err = affair.Exec(template, info.ReceiveAddress.Country, info.ReceiveAddress.City, info.ReceiveAddress.Address)
	receiveAddressId, _ := rows.LastInsertId()

	template = `Insert Into OrderInfo Set OrderId = ?,StartAddressId = ? ,EndAddressId = ? ,Phone= ?,Email = ?,Fax = ? , HopeReachDate = ? , INCOTERMS = ? , UnStackable = ? , Perishable =?,Dangerous = ? , Clearance = ? , Other = ?, deliveryDate = ?`
	_, err = affair.Exec(template, info.OrderId, sendAddressId, receiveAddressId, info.Phone, info.Email, info.Fax, info.HopeReachDate, info.Incoterms, info.UnStackable, info.Perishable, info.Dangerous, info.Clearance, info.Other, info.DeliveryDate)
	wg.Wait()
	affair.Commit()
	return info.OrderId, true, nil
}

func CheckOrderCompany(orderId, companyId int64) bool {
	template := `Select ClientCompanyId From Orders Where OrderId = ?`
	rows, err := Utils.DB().Query(template, orderId)
	if err != nil || !rows.Next() {
		log.Println("[CheckOrderCompany] make a mistake", err)
		return false
	}
	var id int64
	rows.Scan(&id)
	return id == companyId
}

func GetCompanyBargain(orderId, companyId int64) ([]Utils.Bargain, error) {
	template := `Select B.CompanyId,COALESCE(Price,0),COALESCE(isPass,-1) From 
	(Select * From Bargain Where OrderId = ?) A Right Join (
		Select TargetCompanyId as companyId From Relation Where CompanyId = ? And isDelete = 0
		UNION
		Select CompanyId as companyId From Relation Where TargetCompanyId = ? And isDelete = 0
	) B On A.CompanyId = B.companyId`
	rows, err := Utils.DB().Query(template, orderId, companyId, companyId)
	if err != nil {
		return nil, err
	}
	var bargains []Utils.Bargain
	var bargain Utils.Bargain
	for rows.Next() {
		rows.Scan(&bargain.CompanyId, &bargain.Price, &bargain.Status)
		bargains = append(bargains, bargain)
	}
	return bargains, nil
}

func GetOrderInfo(orderId int64) (Utils.OrderInfo, error) {
	var info Utils.OrderInfo
	template := `Select OrderId, StartAddressId, EndAddressId, Phone, Email, Fax, HopeReachDate,deliveryDate,INCOTERMS, UnStackable, Perishable, Dangerous, Clearance, Other From OrderInfo Where OrderId = ? limit 1`
	rows, err := Utils.DB().Query(template, orderId)
	if err != nil {
		return info, err
	}
	defer rows.Close()
	if !rows.Next() {
		return info, fmt.Errorf("订单未查询到")
	}
	var startAddressId, endAddressId int64
	err = rows.Scan(&info.OrderId, &startAddressId, &endAddressId, &info.Phone, &info.Email, &info.Fax, &info.HopeReachDate, &info.DeliveryDate, &info.Incoterms, &info.UnStackable, &info.Perishable, &info.Dangerous, &info.Clearance, &info.Other)
	if err != nil {
		log.Println("[GetOrderInfo] make a mistake ", err)
		return info, err
	}
	defer rows.Close()
	template = `Select Country, City, Address From Address Where AddressId = ?`
	rows, err = Utils.DB().Query(template, startAddressId)
	if err != nil || !rows.Next() {
		log.Println("[GetOrderInfo] make a mistake ", err)
		return info, err
	}
	defer rows.Close()
	rows.Scan(&info.SendAddress.Country, &info.SendAddress.City, &info.SendAddress.Address)
	defer rows.Close()
	rows, err = Utils.DB().Query(template, endAddressId)
	if err != nil || !rows.Next() {
		log.Println("[GetOrderInfo] make a mistake ", err)
		return info, err
	}
	defer rows.Close()
	rows.Scan(&info.ReceiveAddress.Country, &info.ReceiveAddress.City, &info.ReceiveAddress.Address)
	defer rows.Close()
	template = `Select CargoId From Order_Cargo Where OrderId = ?`
	rows, err = Utils.DB().Query(template, orderId)
	if err != nil {
		log.Println("[GetOrderInfo] make a mistake ", err)
		return info, err
	}
	defer rows.Close()
	template = `Select CargoName, CargoModel, CargoNum, Category, Weight, CargoSize From Cargo Where CargoId = ? limit 1`
	var cargoId int64
	var cargo Utils.Cargo
	for rows.Next() {
		rows.Scan(&cargoId)
		rows1, err := Utils.DB().Query(template, cargoId)
		if err != nil || !rows1.Next() {
			log.Println("[GetOrderInfo] make a mistake ", err)
			continue
		}
		rows1.Scan(&cargo.CargoName, &cargo.CargoModel, &cargo.CargoNum, &cargo.Category, &cargo.CargoWeight, &cargo.CargoSize)
		cargo.CargoId = cargoId
		info.Cargos = append(info.Cargos, cargo)
		rows1.Close()
	}
	return info, nil
}

func ReplyBargain(bargain Utils.ReplyBargain, companyId int64) bool {
	template := `Update Bargain Set isPass = 1 , Price = ? ,ReplyTime = now() Where OrderId = ? And CompanyId = ?`
	result, err := Utils.DB().Exec(template, bargain.Bargain, bargain.OrderId, companyId)
	if err != nil {
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}

func AskFroBargain(companyId, orderId int64) bool {
	template := `Insert Into Bargain Set ReplyTime = now() Where OrderId = ? , CompanyId = ? limit 1`
	result, err := Utils.DB().Exec(template, companyId, orderId)
	if err != nil {
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}
