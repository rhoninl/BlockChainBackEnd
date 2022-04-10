package Model

import (
	"log"
	"main/Utils"
)

func Login(account Utils.Account) (Utils.Account, bool, error) {
	var info Utils.Account
	template := `Select PassWord,CompanyId From ShippingTraceability.Account Where Account = ? limit 1`
	rows, err := Utils.DB().Query(template, account.Account)
	if err != nil {
		log.Println("[Login]make Mistake", err)
		return info, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return info, false, nil
	}
	rows.Scan(&info.Password, &info.CompanyId)
	return info, true, nil
}
