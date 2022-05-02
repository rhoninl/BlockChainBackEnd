package Model

import (
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
	template = `Insert Into Cargo Set CargoName = ? , CargoModel = ? , Size = ? , CargoNum = ? , CategoryId = ? , Weight = ? `
	for _, item := range info.Cargos {
		wg.Add(1)
		rows, err := affair.Exec(template, item.CargoName, item.CargoModel, item.CargoSize, item.CargoNum, item.CategoryId, item.CargoWeight)
		if err != nil {
			log.Println("[RecordOrder]Cargo出错了", err)
			affair.Rollback()
			return 0, false, nil
		}
		cargoId, _ := rows.LastInsertId()
		go func(cargoId int64) {
			template := `Insert Into Order_Cargo Set OrderId = ?,Cargo = ?`
			affair.Exec(template, info.OrderId, cargoId)
			wg.Done()
		}(cargoId)
	}
	template = `Insert Into Address Set Country = ?,City = ?,Address = ?`
	rows, err = affair.Exec(template, info.SendAddress.Country, info.SendAddress.City, info.SendAddress.Address)
	sendAddressId, _ := rows.LastInsertId()
	rows, err = affair.Exec(template, info.ReceiveAddress.Country, info.ReceiveAddress.City, info.ReceiveAddress.Address)
	receiveAddressId, _ := rows.LastInsertId()

	template = `Insert Into OrderInfo Set OrderId = ?,StartAddressId = ? ,EndAddressId = ? ,Phone= ?,Email = ?,Fax = ? , HopeReachDate = ? , INCOTERMS = ? , UnStackable = ? , Perishable =?,Dangerous = ? , Clearance = ? , Other = ?`
	_, err = affair.Exec(template, info.OrderId, sendAddressId, receiveAddressId, info.Phone, info.Email, info.Fax, info.HopeReachDate, info.Incoterms, info.UnStackable, info.Perishable, info.Dangerous, info.Clearance, info.Other)
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
