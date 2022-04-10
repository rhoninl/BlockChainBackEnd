package Model

import (
	"log"
	"main/Utils"
)

func GetStuff(companyId string) ([]Utils.Stuff, error) {
	template := `Select StaffId, StaffName, StaffJob From ShippingTraceability.Staff Where CompanyId = ?`
	rows, err := Utils.DB().Query(template, companyId)
	if err != nil {
		log.Println("[GetStuff]", err)
		return nil, err
	}
	defer rows.Close()
	var stuffs []Utils.Stuff
	var stuff Utils.Stuff
	for rows.Next() {
		rows.Scan(&stuff.StuffId, &stuff.StuffName, &stuff.StuffJob)
		stuffs = append(stuffs, stuff)
	}
	return stuffs, nil
}
