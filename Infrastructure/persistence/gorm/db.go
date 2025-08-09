package gormp

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Open เปิดการเชื่อมต่อฐานข้อมูลด้วย DSN จากตัวแปรแวดล้อม DB_DSN
// ตัวอย่าง .env:
// DB_DSN=host=localhost user=postgres password=postgres dbname=books port=5432 sslmode=disable TimeZone=Asia/Bangkok
func Open() (*gorm.DB, error) {
	dataSourceName := os.Getenv("DB_DSN")
	return gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
}
