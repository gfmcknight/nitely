package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

type buildInfo struct {
	gorm.Model
	Name         string `gorm:"not null;unique"`
	AbsolutePath string
	Branch       string
	Remote       string
}

type serviceInfo struct {
	gorm.Model
	Name         string `gorm:"not null;unique"`
	AbsolutePath string
	Args         string
}

type property struct {
	gorm.Model
	Name  string `gorm:"not null;unique"`
	Value string
}

func getStorageBase() string {
	return os.Getenv("NitelyPath")
}

func openAndCreateStorage() *gorm.DB {
	db, err := gorm.Open("sqlite3", filepath.Join(getStorageBase(), "build-info.db"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&buildInfo{}, &serviceInfo{}, &property{})
	return db
}

func getProperty(db *gorm.DB, name string) *string {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var prop property
	if db.Where("name = ?", name).First(&prop).RecordNotFound() {
		return nil
	}

	return &prop.Value
}

func setProperty(db *gorm.DB, name, value string) {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var prop property
	db.FirstOrInit(&prop, property{Name: name})
	prop.Value = value
	db.Save(&prop)
}

func insertBuildInfo(db *gorm.DB, info buildInfo) {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	db.Create(&info)
	if db.Error != nil {
		fmt.Println(db.Error)
	}
}

func deleteBuildInfo(db *gorm.DB, name string) {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var build buildInfo
	db.Where("name = ?", name).First(&build).Delete(&build)
}

func getBuildInfo(db *gorm.DB, name string) *buildInfo {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var build buildInfo
	if db.Where("name = ?", name).First(&build).RecordNotFound() {
		return nil
	}

	return &build
}

func getBuilds(db *gorm.DB) []*buildInfo {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var builds []*buildInfo
	db.Find(&builds)
	return builds
}

func insertServiceInfo(db *gorm.DB, info serviceInfo) {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	db.Create(&info)
	if db.Error != nil {
		fmt.Println(db.Error)
	}
}

func deleteServiceInfo(db *gorm.DB, name string) {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var service serviceInfo
	db.Where("name = ?", name).First(&service).Delete(&service)
}

func getServiceInfo(db *gorm.DB, name string) *serviceInfo {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var service serviceInfo
	if db.Where("name = ?", name).First(&service).RecordNotFound() {
		return nil
	}

	return &service
}

func getServices(db *gorm.DB) []*serviceInfo {
	if db == nil {
		db = openAndCreateStorage()
		defer db.Close()
	}

	var services []*serviceInfo
	db.Find(&services)
	return services
}
