package model

import (
	"devopscenter/helper"
)

type Menu struct {
	MenuCode    string `json:"MenuCode"`
	SystemCode  string `json:"SystemCode"`
	RelativeUrl string `json:"RelativeUrl"`
	MenuNameCN  string `json:"MenuNameCN"`
	MenuNameEN  string `json:"MenuNameEN"`
	ParentCode  string `json:"ParentCode"`
	Category    string `json:"Category"`
	SortNo      int    `json:"SortNo"`
	IsPublic    int    `json:"IsPublic"`
	IsDelete    int    `json:"IsDelete"`
	IsValid     int    `json:"IsValid"`
	CreateDate  string `json:"CreateDate"`
	ModifyDate  string `json:"ModifyDate"`
}

func (m Menu) ListMenus(name string) ([]*Menu, error) {
	query := `SELECT a.MenuCode, a.SystemCode, a.RelativeUrl, a.MenuNameCN, a.MenuNameEN, a.ParentCode, a.Category, a.SortNo, a.IsPublic,
        a.IsDelete, a.IsValid, a.CreateDate, a.ModifyDate
        FROM auth_menu a 
        WHERE a.IsDelete=0 AND a.IsValid=1 
            AND EXISTS(
                SELECT 1 FROM auth_menu_role b 
                WHERE b.IsDelete=0 AND b.IsValid=1 AND a.MenuCode=b.MenuCode 
                    AND b.SystemCode='devops' 
                    AND b.RoleCode IN(
                        SELECT r.rolecode FROM auth_role r 
                        WHERE r.IsDelete=0 AND r.IsValid=1 
                            AND r.rolecode IN (
                                SELECT ru.RoleCode FROM auth_role_user ru 
                                WHERE ru.IsDelete=0 AND ru.IsValid=1 
                                    AND EXISTS(
                                        SELECT 1 FROM auth_user_info ui 
                                        WHERE ui.STATUS=1 AND ui.EmployeeNo=ru.EmployeeNo AND ui.ADAccount=?
                                    )
                            )
                    )
            ) 
        ORDER BY a.SortNo`

	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]*Menu, 0)
	for rows.Next() {
		obj := &Menu{}
		err = rows.Scan(&obj.MenuCode, &obj.SystemCode, &obj.RelativeUrl, &obj.MenuNameCN, &obj.MenuNameEN, &obj.ParentCode, &obj.Category, &obj.SortNo, &obj.IsPublic,
			&obj.IsDelete, &obj.IsValid, &obj.CreateDate, &obj.ModifyDate)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}

	return data, nil
}
