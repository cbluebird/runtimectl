package dao

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"runtimectl/config"
	"runtimectl/model"
	"time"
)

var DB *gorm.DB

func Init() {
	user := config.Config.GetString("database.user")
	pass := config.Config.GetString("database.pass")
	port := config.Config.GetString("database.port")
	host := config.Config.GetString("database.host")
	name := config.Config.GetString("database.name")
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=require TimeZone=Asia/Shanghai", host, user, pass, name, port)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
}

func GetOrganization(name string) (model.Organization, error) {
	var organization model.Organization
	result := DB.First(&organization, "name = ?", name)
	if result.Error != nil {
		log.Println("Error getting organization:", result.Error)
		return model.Organization{}, result.Error
	}
	return organization, nil
}

func CreateOrUpdateTemplateRepository(name, kind string) error {
	organization, err := GetOrganization("labring")
	if err != nil {
		log.Println("Error getting organization:", err)
		return err
	}

	tmp := model.TemplateRepository{}
	result := DB.Model(&model.TemplateRepository{}).Where(&model.TemplateRepository{
		Name:            name,
		Kind:            kind,
		OrganizationUid: organization.UID,
	}).First(&tmp).Error

	if result == nil {
		return nil
	}

	repo := model.TemplateRepository{
		Name:            name,
		Kind:            kind,
		IconId:          name,
		OrganizationUid: organization.UID,
		IsPublic:        true,
	}
	return DB.Create(&repo).Error
}

func GetTemplateRepository(name string) *model.TemplateRepository {
	t := &model.TemplateRepository{}
	DB.Model(&model.TemplateRepository{}).Where(&model.TemplateRepository{Name: name}).First(t)
	return t
}

func CreateOrUpdateTemplate(version, repoUid, image, config, state string, deleteTime time.Time) error {
	flag := state == "active"
	template := model.Template{
		Name:                  version,
		TemplateRepositoryUid: repoUid,
		Image:                 image,
		Config:                config,
	}
	if flag {
		template.DeletedTime = nil
	} else {
		template.DeletedTime = &deleteTime
	}

	log.Println("updating template:", template.DeletedTime)

	tmp := model.Template{}
	result := DB.Model(&model.Template{}).Where(&model.Template{
		Name:                  version,
		TemplateRepositoryUid: repoUid,
		Image:                 image,
		Config:                config,
	}).First(&tmp).Error

	if errors.Is(result, gorm.ErrRecordNotFound) {
		if flag {
			sql := `INSERT INTO "Template" ("name", "templateRepositoryUid", "image", "config", "isDeleted", "deletedAt", "createdAt", "updatedAt", "parentUid") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
			if err := DB.Exec(sql, template.Name, template.TemplateRepositoryUid, template.Image, template.Config, false, template.DeletedTime, time.Now(), time.Now(), template.ParentUid).Error; err != nil {
				return err
			}
		} else {
			sql := `INSERT INTO "Template" ("name", "templateRepositoryUid", "image", "config", "isDeleted", "deletedAt", "createdAt", "updatedAt", "parentUid") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
			if err := DB.Exec(sql, template.Name, template.TemplateRepositoryUid, template.Image, template.Config, nil, template.DeletedTime, time.Now(), time.Now(), template.ParentUid).Error; err != nil {
				return err
			}
		}
		return nil
	}

	if flag {
		sql := `UPDATE "Template" SET "image" = ?, "config" = ?, "isDeleted" = ?, "deletedAt" = ?, "updatedAt" = ? WHERE "uid" = ?`
		if err := DB.Exec(sql, template.Image, template.Config, false, template.DeletedTime, time.Now(), tmp.UID).Error; err != nil {
			return err
		}
	} else {
		sql := `UPDATE "Template" SET "image" = ?, "config" = ?, "isDeleted" = ?, "deletedAt" = ?, "updatedAt" = ? WHERE "uid" = ?`
		if err := DB.Exec(sql, template.Image, template.Config, nil, template.DeletedTime, time.Now(), tmp.UID).Error; err != nil {
			return err
		}
	}

	return nil
}

func GetTemplates() (map[string]string, error) {
	var templates []model.Template
	result := DB.Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}

	data := make(map[string]string)
	for _, template := range templates {
		var repo model.TemplateRepository
		if err := DB.First(&repo, "uid = ?", template.TemplateRepositoryUid).Error; err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%s-%s", repo.Name, template.Name)
		data[key] = template.UID
	}
	return data, nil
}

func GetTemplateID(class, version, image string) (string, error) {
	var template model.Template
	var templateRepository model.TemplateRepository
	result := DB.Model(&model.TemplateRepository{}).Where(&model.TemplateRepository{Name: class}).First(&templateRepository).Error
	if result != nil {
		log.Println("Error getting template repository:", result)
		return "", result
	}
	result = DB.Model(&model.Template{}).Where(&model.Template{
		TemplateRepositoryUid: templateRepository.UID,
		Name:                  version,
		Image:                 image,
	}).First(&template).Error
	return template.UID, result
}
