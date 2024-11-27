package dao

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"runtimectl/config"
	"runtimectl/model"
)

var DB *gorm.DB

func Init() {
	user := config.Config.GetString("database.user")
	pass := config.Config.GetString("database.pass")
	port := config.Config.GetString("database.port")
	host := config.Config.GetString("database.host")
	name := config.Config.GetString("database.name")
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v TimeZone=Asia/Shanghai", host, user, pass, name, port)
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
		OrganizationUid: organization.UID,
	}
	return DB.Create(&repo).Error
}

func GetTemplateRepository(name string) *model.TemplateRepository {
	t := &model.TemplateRepository{}
	DB.Model(&model.TemplateRepository{}).Where(&model.TemplateRepository{Name: name}).First(t)
	return t
}

func CreateOrUpdateTemplate(version, repoUid, image, config string) error {
	template := model.Template{
		Name:                  version,
		TemplateRepositoryUid: repoUid,
		Image:                 image,
		Config:                config,
	}

	tmp := model.Template{}
	result := DB.Model(&model.Template{}).Where(&model.Template{
		Name:                  version,
		TemplateRepositoryUid: repoUid,
		Image:                 image,
		Config:                config,
	}).First(&tmp).Error

	if result == nil {
		return nil
	}

	log.Println("result:", tmp)

	result = DB.Model(&model.Template{}).Where(&model.Template{
		Name:                  version,
		TemplateRepositoryUid: repoUid,
	}).First(&tmp).Error

	if errors.Is(result, gorm.ErrRecordNotFound) {
		if result = DB.Create(&template).Error; result != nil {
			return result
		}
		return nil
	}

	return DB.Model(&model.Template{}).Where(&model.Template{UID: tmp.UID}).Updates(&template).Error
}
