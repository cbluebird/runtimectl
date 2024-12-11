package model

import (
	"time"
)

type Organization struct {
	UID       string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ID        string     `gorm:"unique;not null"`
	CreatedAt time.Time  `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
	Name      string     `gorm:"not null"`
}

func (Organization) TableName() string {
	return "Organization"
}

type TemplateRepository struct {
	UID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DeletedAt       *time.Time `gorm:"column:deletedAt;index"`
	CreatedAt       time.Time  `gorm:"column:createdAt;not null;default:current_timestamp"`
	UpdatedAt       time.Time  `gorm:"column:updatedAt;not null"`
	Name            string     `gorm:"not null"`
	Description     string
	OrganizationUid string `gorm:"column:organizationUid;not null"`
	IsPublic        bool   `gorm:"column:isPublic;default:false;not null"`
	IconId          string `gorm:"column:iconId"`
	Kind            string `gorm:"not null"`
}

func (TemplateRepository) TableName() string {
	return "TemplateRepository"
}

type Template struct {
	UID                   string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name                  string     `gorm:"not null;uniqueIndex:idx_version_template_repo_uid"`
	TemplateRepositoryUid string     `gorm:"column:templateRepositoryUid;not null;uniqueIndex:idx_version_template_repo_uid"`
	DevboxReleaseImage    *string    `gorm:"column:devboxReleaseImage"`
	Image                 string     `gorm:"not null"`
	Config                string     `gorm:"not null"`
	DeletedTime           *time.Time `gorm:"column:deletedAt;index"`
	CreatedAt             time.Time  `gorm:"column:createdAt;not null;default:current_timestamp"`
	UpdatedAt             time.Time  `gorm:"column:updatedAt;not null"`
	ParentUid             *string    `gorm:"column:parentUid;type:uuid"`
	IsDeleted             bool       `gorm:"column:isDeleted;default:null"`
}

func (Template) TableName() string {
	return "Template"
}
