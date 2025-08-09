package interfaces

import (
	"github.com/nuba55yo/go-101-CleanCRUD/domain"
)

// BookRepository คือพอร์ตออกจาก use case ไปยังเลเยอร์ persistence
// เลเยอร์ infrastructure ต้อง implement อินเทอร์เฟซนี้ (เช่น GORM repository)
type BookRepository interface {
	List() ([]domain.Book, error)
	GetByID(id uint) (domain.Book, error)
	ExistsActiveByTitle(title string, excludeID *uint) (bool, error)
	Create(book *domain.Book) error
	Update(book *domain.Book) error
	SoftDelete(id uint) error
}
