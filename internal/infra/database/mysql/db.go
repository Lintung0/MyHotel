package mysql

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"backend/internal/config"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal menghubungkan ke database Mysql %v", err)
	}

	log.Println("Berhasil terhubung dengan database MySql")
	return db
}

func AutoMigrate(db *gorm.DB, models ...interface{}) {
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("gagal auto migrate tabel %v", err)
	}
	log.Println("Migrasi tabel (AutoMigrate) sukses.")
}