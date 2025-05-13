package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name         string         
	Email        string         
	Password     string         
	PasswordView string         
	PhoneNumber  string         `gorm:"default:null"`
	GoogleId     string         `gorm:"default:null"`
	Token        string         `gorm:"default:null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Users) TableName() string{
	return "users"
}

type Province struct{
	ID  		int64  `gorm:"type:int;primaryKey"`
	Name 		string
	Latitude 	float64
	Longitude 	float64
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Province) TableName() string{
	return "province"
}

type Regency struct{
	ID 		   		int64  `gorm:"type:int;primaryKey"`
	ProvinceId 		int64
	Name 			string
	Latitude 		float64
	Longitude 		float64
	CreatedAt    	time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    	time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    	gorm.DeletedAt `gorm:"index"`
}

func (Regency) TableName() string{
	return "regency"
}

type SubDistrict struct{
	ID 			 int64  `gorm:"type:int;primaryKey"`
	RegencyId    int64
	Name         string
	Latitude     float64
	Longitude    float64
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (SubDistrict) TableName() string{
	return "sub_district"
}

type Company struct {
	ID 				uuid.UUID  	  `gorm:"type:uuid;primaryKey"`
	ProvinceId 		int64
	RegencyId 		int64
	SubDistrictId 	int64
	Name         	string    
	Address		 	string 	   	
	Photo 			*string 	   `gorm:"default:null"`     
	CreatedAt    	time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    	time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    	gorm.DeletedAt `gorm:"index"`
}

func (Company) TableName() string{
	return "company"
}

type UserCompanyAccess struct{
	ID				int64 			`gorm:"type:int;primaryKey;autoIncrement"`	
	UserId			uuid.UUID		`gorm:"type:uuid"`
	CompanyId		uuid.UUID 		`gorm:"type:uuid"`
	RoleId			uuid.UUID		`gorm:"type:uuid"`
	IsOwner			int // Enum 1 or 0
	CurrentCompany 	int // Enum 1 or 0
	CreatedAt    	time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    	time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    	gorm.DeletedAt `gorm:"index"`
}

func (UserCompanyAccess) TableName() string{
	return "user_company_access"
}

type Menu struct{
	ID 		 		int64		   `gorm:"type:int;primaryKey;autoIncrement"`
	ParentId 		int64	       `gorm:"type:int"`
	NameId	 		string
	NameEn	 		string
	Icon	 		*string		   `gorm:"default:null"`
	Url		 		string		   `gorm:"default:null"`
	Order			int64
	CreatedAt    	time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    	time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    	gorm.DeletedAt `gorm:"index"`
}

func (Menu) TableName()string{
	return "menu"
}

type Role struct{
	ID				uuid.UUID	`gorm:"type:uuid;primaryKey"`
	CompanyId 		uuid.UUID	`gorm:"type:uuid"`
	Name			string		
	CreatedAt    	time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    	time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    	gorm.DeletedAt `gorm:"index"`
}

func (Role) TableName()string{
	return "role"
}

type Permission struct{
	ID 				int64	`gorm:"type:int;primaryKey;autoIncrement"`
	MenuId 			int64
	Code			string
	NameId 			string
	NameEn  		string
	PermissionView	int64		   `gorm:"default:null"`
	CreatedAt    	time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    	time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    	gorm.DeletedAt `gorm:"index"`	
}

func (Permission) TableName()string{
	return "permission"
}

type RoleHasPermission struct{
	ID 				int64	`gorm:"type:int;primaryKey;autoIncrement"`
	RoleId 			uuid.UUID
	PermissionId 	int64
	CreatedAt    	time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    	time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    	gorm.DeletedAt `gorm:"index"`	
}

func (RoleHasPermission) TableName()string{
	return "role_has_permission"
}