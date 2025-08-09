package dto

// ชุด DTO ที่ “use case” รับ/คืน (ไม่ผูกกับ HTTP/JSON)

type CreateBookCommand struct {
	Title  string
	Author string
}

type UpdateBookCommand struct {
	ID     uint
	Title  string
	Author string
}

type BookReadModel struct {
	ID        uint
	Title     string
	Author    string
	CreatedAt string // ฟอร์แมตแล้ว (เช่น RFC3339Nano) เพื่อส่งต่อ presentation ได้เลย
	UpdatedAt string
}
