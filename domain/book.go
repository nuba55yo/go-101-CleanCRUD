package domain

import "time"

// Book = เอนทิตีหลักของธุรกิจ (ไม่ผูกกับ HTTP/DB/Framework ใด ๆ)
// เก็บแค่ “สภาพจริง” ของหนังสือในระบบ
type Book struct {
	ID        uint
	Title     string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time // ใช้ soft delete; ถ้ายังไม่ลบจะเป็น nil
}
