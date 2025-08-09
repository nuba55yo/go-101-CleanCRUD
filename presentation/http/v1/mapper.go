package v1

import "github.com/nuba55yo/go-101-CleanCRUD/application/dto"

func MapCreateJSONToCommand(requestBody CreateBookJSON) dto.CreateBookCommand {
	return dto.CreateBookCommand{Title: requestBody.Title, Author: requestBody.Author}
}

func MapUpdateJSONToCommand(id uint, requestBody UpdateBookJSON) dto.UpdateBookCommand {
	return dto.UpdateBookCommand{ID: id, Title: requestBody.Title, Author: requestBody.Author}
}

func MapReadModelToJSON(readModel dto.BookReadModel) BookJSON {
	return BookJSON{
		ID:        readModel.ID,
		Title:     readModel.Title,
		Author:    readModel.Author,
		CreatedAt: readModel.CreatedAt,
		UpdatedAt: readModel.UpdatedAt,
	}
}

func MapReadModelsToJSON(readModels []dto.BookReadModel) []BookJSON {
	result := make([]BookJSON, 0, len(readModels))
	for _, m := range readModels {
		result = append(result, MapReadModelToJSON(m))
	}
	return result
}
