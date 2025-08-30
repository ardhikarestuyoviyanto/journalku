package response

import "github.com/google/uuid"

type ResToken struct {
	User ResUserToken `json:"user"`
	Exp  int64        `json:"exp"`
}

type ResUserToken struct {
	ID             uuid.UUID         `json:"id"`
	Name           string            `json:"name"`
	Email          string            `json:"email"`
	CurrentCompany *ResCurrentCompany `json:"currentCompany"`
}

type ResCurrentCompany struct {
	ID         uuid.UUID `json:"id"` // companyId
	Name       string    `json:"name"`
	Photo      *string    `json:"photo"`
	IsOwner    int64     `json:"isOwner"`
	Address    string    `json:"address"`
	Role       string    `json:"role"`
	Permission []string  `json:"permission"`
	ResMenu    []ResMenu `json:"menu"`
}

type ResMenu struct {
	ID       int64   `json:"id"`
	ParentId int64   `json:"parentId"`
	NameId   string  `json:"nameId"`
	NameEn   string  `json:"nameEn"`
	Icon     *string `json:"icon"`
	Url      string  `json:"url"`
	Order    int64   `json:"order"`
	Child    []ResMenu	`json:"child"`
}

type ResCompanyAccess struct{
	ID         uuid.UUID `json:"id"` // companyId
	Name       string    `json:"name"`
	Photo      *string    `json:"photo"`
	IsOwner    int64     `json:"isOwner"`
	Address    string    `json:"address"`
	Role       string    `json:"role"`
}

type ResInit struct{
	CompanyAccess []ResCompanyAccess `json:"companyAccess"`
}

type ResAccount struct {
	ID uuid.UUID `json:"id"`
	CategoryAccount string `json:"category_account"`
	NumberAccount string `json:"number_account"`
	Name string `json:"name"`
	Description string `json:"description"`
	IsPrimary int64 `json:"is_primary"`
	IsArchive int64 `json:"is_archive"`
}

type ResCategoryAccount struct{
	ID int64 `json:"id"`
	Value string `json:"code"` // Supaya mudah di fe
	Description string `json:"category"` // Supaya mudah di fe
}