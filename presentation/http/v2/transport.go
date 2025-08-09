package v2

// v2: ตัวอย่าง response แบบห่อ version/data
type CreateBookJSON struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type UpdateBookJSON struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookData struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BookJSON struct {
	Version string   `json:"version"` // "v2"
	Data    BookData `json:"data"`
}
