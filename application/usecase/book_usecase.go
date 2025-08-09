package usecase

import (
	"context"
	"strings"
	"time"

	// เปลี่ยนโมดูลให้ตรงกับของคุณ ถ้าไม่ใช่ path นี้
	"github.com/nuba55yo/go-101-CleanCRUD/application/dto"
	"github.com/nuba55yo/go-101-CleanCRUD/application/interfaces"
	"github.com/nuba55yo/go-101-CleanCRUD/domain"
)

// BookUseCase = พอร์ตเข้า (ให้ presentation เรียก)
type BookUseCase interface {
	Create(requestContext context.Context, command dto.CreateBookCommand) (dto.BookReadModel, error)
	Update(requestContext context.Context, command dto.UpdateBookCommand) (dto.BookReadModel, error)
	Get(requestContext context.Context, id uint) (dto.BookReadModel, error)
	List(requestContext context.Context) ([]dto.BookReadModel, error)
	Delete(requestContext context.Context, id uint) error
}

// bookUseCase = implementation ของพอร์ตข้างบน
type bookUseCase struct {
	bookRepository interfaces.BookRepository
	clock          interfaces.Clock
	logger         interfaces.Logger
}

// NewBookUseCase ประกอบ dependencies ให้พร้อมใช้
func NewBookUseCase(
	bookRepository interfaces.BookRepository,
	clock interfaces.Clock,
	logger interfaces.Logger,
) BookUseCase {
	return &bookUseCase{
		bookRepository: bookRepository,
		clock:          clock,
		logger:         logger,
	}
}

// Create: ตรวจ input, เช็คชื่อซ้ำ (ไม่นับเล่มที่ลบแบบ soft delete), เซฟ, แล้วคืน ReadModel
func (useCase *bookUseCase) Create(
	requestContext context.Context,
	command dto.CreateBookCommand,
) (dto.BookReadModel, error) {

	title := strings.TrimSpace(command.Title)
	author := strings.TrimSpace(command.Author)
	if title == "" || author == "" {
		return dto.BookReadModel{}, domain.ErrBadInput
	}

	isDuplicate, existsError := useCase.bookRepository.
		ExistsActiveByTitle(strings.ToLower(title), nil)
	if existsError != nil {
		return dto.BookReadModel{}, existsError
	}
	if isDuplicate {
		return dto.BookReadModel{}, domain.ErrTitleExists
	}

	now := useCase.clock.Now()
	entity := domain.Book{
		Title:     title,
		Author:    author,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if createError := useCase.bookRepository.Create(&entity); createError != nil {
		return dto.BookReadModel{}, createError
	}

	useCase.logger.Info(requestContext, "book created",
		"id", entity.ID, "title", entity.Title, "author", entity.Author)

	return dto.BookReadModel{
		ID:        entity.ID,
		Title:     entity.Title,
		Author:    entity.Author,
		CreatedAt: now.Format(time.RFC3339Nano),
		UpdatedAt: now.Format(time.RFC3339Nano),
	}, nil
}

// Update: ตรวจ input, เช็คชื่อซ้ำ (ยกเว้นเล่มตัวเอง), โหลดของเดิม, แก้ไข, เซฟ
func (useCase *bookUseCase) Update(
	requestContext context.Context,
	command dto.UpdateBookCommand,
) (dto.BookReadModel, error) {

	title := strings.TrimSpace(command.Title)
	author := strings.TrimSpace(command.Author)
	if title == "" || author == "" {
		return dto.BookReadModel{}, domain.ErrBadInput
	}

	isDuplicate, existsError := useCase.bookRepository.
		ExistsActiveByTitle(strings.ToLower(title), &command.ID)
	if existsError != nil {
		return dto.BookReadModel{}, existsError
	}
	if isDuplicate {
		return dto.BookReadModel{}, domain.ErrTitleExists
	}

	currentEntity, getError := useCase.bookRepository.GetByID(command.ID)
	if getError != nil {
		return dto.BookReadModel{}, getError // รวมทั้งกรณี ErrNotFound
	}

	currentEntity.Title = title
	currentEntity.Author = author
	currentEntity.UpdatedAt = useCase.clock.Now()

	if updateError := useCase.bookRepository.Update(&currentEntity); updateError != nil {
		return dto.BookReadModel{}, updateError
	}

	return dto.BookReadModel{
		ID:        currentEntity.ID,
		Title:     currentEntity.Title,
		Author:    currentEntity.Author,
		CreatedAt: currentEntity.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: currentEntity.UpdatedAt.Format(time.RFC3339Nano),
	}, nil
}

// Get: ดึงเล่มเดียวแล้วแปลงเป็น ReadModel
func (useCase *bookUseCase) Get(
	requestContext context.Context,
	id uint,
) (dto.BookReadModel, error) {

	entity, getError := useCase.bookRepository.GetByID(id)
	if getError != nil {
		return dto.BookReadModel{}, getError // รวมทั้งกรณี ErrNotFound
	}

	return dto.BookReadModel{
		ID:        entity.ID,
		Title:     entity.Title,
		Author:    entity.Author,
		CreatedAt: entity.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: entity.UpdatedAt.Format(time.RFC3339Nano),
	}, nil
}

// List: ดึงทั้งหมดแล้ว map เป็น ReadModel slice
func (useCase *bookUseCase) List(
	requestContext context.Context,
) ([]dto.BookReadModel, error) {

	entities, listError := useCase.bookRepository.List()
	if listError != nil {
		return nil, listError
	}

	readModels := make([]dto.BookReadModel, 0, len(entities))
	for _, entity := range entities {
		readModels = append(readModels, dto.BookReadModel{
			ID:        entity.ID,
			Title:     entity.Title,
			Author:    entity.Author,
			CreatedAt: entity.CreatedAt.Format(time.RFC3339Nano),
			UpdatedAt: entity.UpdatedAt.Format(time.RFC3339Nano),
		})
	}
	return readModels, nil
}

// Delete: ลบแบบ soft delete
func (useCase *bookUseCase) Delete(
	requestContext context.Context,
	id uint,
) error {
	return useCase.bookRepository.SoftDelete(id)
}
