package handler

import (
	"desktop-assistant/internal/app/usecase"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CoverLetterHandler struct {
	coverLetterUseCase usecase.CoverLetterGenerator
}

func NewCoverLetterHandler(coverLetterUseCase usecase.CoverLetterGenerator) *CoverLetterHandler {
	return &CoverLetterHandler{
		coverLetterUseCase: coverLetterUseCase,
	}
}

func (h *CoverLetterHandler) GeneratePersonalLetter(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	jobDescription := c.FormValue("jobDescription")
	additionalInfo := c.FormValue("additionalInfo")

	result, err := h.coverLetterUseCase.GeneratePersonalLetter(file, jobDescription, additionalInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Set the necessary headers for the PDF
	fileName := file.Filename
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileName))

	// Write the buffer to HTTP response
	c.Response().WriteHeader(http.StatusOK)
	return c.Blob(http.StatusOK, "application/pdf", result)
}
