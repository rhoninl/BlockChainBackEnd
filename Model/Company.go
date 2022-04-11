package Model

import (
	"log"
	"main/Utils"
)

func GetJointVenture(companyId string) ([]Utils.CompanyList, error) {
	template := `Select CompanyId, TargetCompanyId From ShippingTraceability.Relation Where CompanyId = ? Or TargetCompanyId = ?`
	rows, err := Utils.DB().Query(template, companyId, companyId)
	if err != nil {
		log.Println("[GetJointVenture]数据库异常", err)
		return nil, err
	}
	defer rows.Close()
	var companyList []Utils.CompanyList
	var company Utils.CompanyList
	var aCompanyId, bCompanyId string
	for rows.Next() {
		rows.Scan(&aCompanyId, &bCompanyId)
		if companyId == aCompanyId {
			aCompanyId = bCompanyId
		}
		company.CompanyId = aCompanyId
		company.CompanyName, company.CompanyType = GetCompanyBasicInfo(aCompanyId)
		companyList = append(companyList, company)
	}
	return companyList, nil
}
