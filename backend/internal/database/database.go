package database

import (
	"fmt"
	"log"
	"nexushub-personal/internal/config"
	"nexushub-personal/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DB.User,
		config.AppConfig.DB.Password,
		config.AppConfig.DB.Host,
		config.AppConfig.DB.Port,
		config.AppConfig.DB.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate models
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// Initialize default user
	if err := initDefaultUser(); err != nil {
		return fmt.Errorf("failed to initialize default user: %v", err)
	}

	// Initialize default theme
	if err := initDefaultTheme(); err != nil {
		return fmt.Errorf("failed to initialize default theme: %v", err)
	}

	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Note{},
		&model.File{},
		&model.Task{},
		&model.Bookmark{},
		&model.Theme{},
		&model.ChatMessage{},
		&model.Collection{},
		&model.Event{},
	)
}

func initDefaultUser() error {
	var count int64
	if err := DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		defaultUser := &model.User{
			ID:       1,
			Username: "admin",
			Nickname: "NexusHub User",
			Bio:      "Welcome to NexusHub Personal Workstation",
		}
		if err := DB.Create(defaultUser).Error; err != nil {
			return err
		}
		log.Println("Default user created successfully")
	}

	return nil
}

func initDefaultTheme() error {
	var count int64
	if err := DB.Model(&model.Theme{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		defaultTheme := &model.Theme{
			UserID:         1,
			ThemeName:      "dark",
			PrimaryColor:   "#000000",
			SecondaryColor: "#ffffff",
			ThemeTemplate:  "default",
			MusicVolume:    0.5,
		}
		if err := DB.Create(defaultTheme).Error; err != nil {
			return err
		}
		log.Println("Default theme created successfully")
	}

	return nil
}
