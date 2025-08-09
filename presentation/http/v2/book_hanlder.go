package v2

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nuba55yo/go-101-CleanCRUD/application/usecase"
	"github.com/nuba55yo/go-101-CleanCRUD/domain"
)

// @Summary Create book (v2)
// @Tags books
// @Accept json
// @Produce json
// @Param body body CreateBookJSON true "payload"
// @Success 201 {object} BookJSON
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books [post]
func CreateBook(bookUseCase usecase.BookUseCase) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		var requestBody CreateBookJSON
		if bindError := requestContext.ShouldBindJSON(&requestBody); bindError != nil {
			requestContext.JSON(http.StatusBadRequest, gin.H{"error": bindError.Error()})
			return
		}
		readModel, createError := bookUseCase.Create(requestContext, MapCreateJSONToCommand(requestBody))
		if createError != nil {
			switch {
			case errors.Is(createError, domain.ErrTitleExists):
				requestContext.JSON(http.StatusConflict, gin.H{"error": "title already exists"})
			case errors.Is(createError, domain.ErrBadInput):
				requestContext.JSON(http.StatusBadRequest, gin.H{"error": "title and author are required"})
			default:
				requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
			}
			return
		}
		requestContext.JSON(http.StatusCreated, MapReadModelToJSON(readModel))
	}
}

// @Summary List books (v2)
// @Tags books
// @Produce json
// @Success 200 {array} BookJSON
// @Router /books [get]
func ListBooks(bookUseCase usecase.BookUseCase) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		readModels, listError := bookUseCase.List(requestContext)
		if listError != nil {
			requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "cannot list books"})
			return
		}
		requestContext.JSON(http.StatusOK, MapReadModelsToJSON(readModels))
	}
}

// @Summary Get book by id (v2)
// @Tags books
// @Produce json
// @Param id path int true "book id"
// @Success 200 {object} BookJSON
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func GetBookByID(bookUseCase usecase.BookUseCase) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		idText := requestContext.Param("id")
		idNumber, convertError := strconv.Atoi(idText)
		if convertError != nil {
			requestContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		readModel, getError := bookUseCase.Get(requestContext, uint(idNumber))
		if errors.Is(getError, domain.ErrNotFound) {
			requestContext.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if getError != nil {
			requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "get failed"})
			return
		}
		requestContext.JSON(http.StatusOK, MapReadModelToJSON(readModel))
	}
}

// @Summary Update book (v2)
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "book id"
// @Param body body UpdateBookJSON true "payload"
// @Success 200 {object} BookJSON
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books/{id} [put]
func UpdateBook(bookUseCase usecase.BookUseCase) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		idText := requestContext.Param("id")
		idNumber, convertError := strconv.Atoi(idText)
		if convertError != nil {
			requestContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var requestBody UpdateBookJSON
		if bindError := requestContext.ShouldBindJSON(&requestBody); bindError != nil {
			requestContext.JSON(http.StatusBadRequest, gin.H{"error": bindError.Error()})
			return
		}
		readModel, updateError := bookUseCase.Update(requestContext, MapUpdateJSONToCommand(uint(idNumber), requestBody))
		if updateError != nil {
			switch {
			case errors.Is(updateError, domain.ErrTitleExists):
				requestContext.JSON(http.StatusConflict, gin.H{"error": "title already exists"})
			case errors.Is(updateError, domain.ErrBadInput):
				requestContext.JSON(http.StatusBadRequest, gin.H{"error": "title and author are required"})
			case errors.Is(updateError, domain.ErrNotFound):
				requestContext.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			default:
				requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			}
			return
		}
		requestContext.JSON(http.StatusOK, MapReadModelToJSON(readModel))
	}
}

// @Summary Delete book (v2) (soft delete)
// @Tags books
// @Param id path int true "book id"
// @Success 204
// @Router /books/{id} [delete]
func DeleteBook(bookUseCase usecase.BookUseCase) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		idText := requestContext.Param("id")
		idNumber, convertError := strconv.Atoi(idText)
		if convertError != nil {
			requestContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if deleteError := bookUseCase.Delete(requestContext, uint(idNumber)); deleteError != nil {
			requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}
		requestContext.Status(http.StatusNoContent)
	}
}
