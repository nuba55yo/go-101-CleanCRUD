package v1

// โครง JSON สำหรับ HTTP v1 (ใช้ bind/Swagger)
type CreateBookJSON struct {
	Title  string `json:"title"  example:"Domain-Driven Design"`
	Author string `json:"author" example:"Eric Evans"`
}

type UpdateBookJSON struct {
	Title  string `json:"title"  example:"DDD 2nd"`
	Author string `json:"author" example:"Eric Evans"`
}

type BookJSON struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
