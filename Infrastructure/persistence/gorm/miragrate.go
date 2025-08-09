package gormp

import "gorm.io/gorm"

// AutoMigrateTables สร้าง/อัปเดตตารางตามโครงสร้าง record ในชั้น infrastructure
func AutoMigrateTables(database *gorm.DB) error {
	return database.AutoMigrate(&bookRecord{})
}

// EnsureIndexes สร้าง unique index ป้องกันชื่อซ้ำ (เฉพาะที่ยังไม่ถูก soft delete)
func EnsureIndexes(database *gorm.DB) error {
	return database.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ux_books_title_active
        ON public.books (lower(title)) WHERE deleted_at IS NULL;`).Error
}
