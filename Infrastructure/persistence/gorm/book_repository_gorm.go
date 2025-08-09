package gormp

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/nuba55yo/go-101-CleanCRUD/application/interfaces"
	"github.com/nuba55yo/go-101-CleanCRUD/domain"
)

// bookRecord = โครงสร้างตารางจริงในฐานข้อมูล (เลเยอร์ infrastructure)
// แยกออกจาก domain.Book เพื่อให้ mapping/constraint เป็นเรื่องฝั่ง infra
type bookRecord struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"not null"`
	Author    string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (bookRecord) TableName() string { return "books" }

// BookRepositoryGorm = อแดปเตอร์ที่ implement พอร์ต interfaces.BookRepository
type BookRepositoryGorm struct {
	database *gorm.DB
}

func NewBookRepositoryGorm(database *gorm.DB) interfaces.BookRepository {
	return &BookRepositoryGorm{database: database}
}

func toDomain(record bookRecord) domain.Book {
	var deletedAt *time.Time
	if record.DeletedAt.Valid {
		t := record.DeletedAt.Time
		deletedAt = &t
	}
	return domain.Book{
		ID:        record.ID,
		Title:     record.Title,
		Author:    record.Author,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (repository *BookRepositoryGorm) List() ([]domain.Book, error) {
	var records []bookRecord
	if err := repository.database.Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]domain.Book, 0, len(records))
	for _, r := range records {
		result = append(result, toDomain(r))
	}
	return result, nil
}

func (repository *BookRepositoryGorm) GetByID(id uint) (domain.Book, error) {
	var record bookRecord
	if err := repository.database.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Book{}, domain.ErrNotFound
		}
		return domain.Book{}, err
	}
	return toDomain(record), nil
}

func (repository *BookRepositoryGorm) ExistsActiveByTitle(title string, excludeID *uint) (bool, error) {
	query := repository.database.
		Model(&bookRecord{}).
		Where("lower(title) = ? AND deleted_at IS NULL", strings.ToLower(title))
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *BookRepositoryGorm) Create(book *domain.Book) error {
	record := bookRecord{
		Title:     book.Title,
		Author:    book.Author,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
	if err := repository.database.Create(&record).Error; err != nil {
		return err
	}
	book.ID = record.ID
	return nil
}

func (repository *BookRepositoryGorm) Update(book *domain.Book) error {
	return repository.database.
		Model(&bookRecord{}).
		Where("id = ?", book.ID).
		Updates(map[string]any{
			"title":      book.Title,
			"author":     book.Author,
			"updated_at": book.UpdatedAt,
		}).Error
}

func (repository *BookRepositoryGorm) SoftDelete(id uint) error {
	return repository.database.Delete(&bookRecord{}, id).Error
}
