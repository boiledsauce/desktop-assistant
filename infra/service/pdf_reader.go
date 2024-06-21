package service

import (
	"bytes"
	"desktop-assistant/internal/domain"
	"log"
	"strings"

	"github.com/go-pdf/fpdf"
	"github.com/ledongthuc/pdf"
)

// PDFReaderImpl implements the interfaces.PDFReader interface.
type PDFReader interface {
	ExtractText(*domain.File) (string, error)
	ConvertTextToPdf(string) ([]byte, error)
}

type PDFReaderImpl struct {
}

func NewPDFReader() *PDFReaderImpl {
	return &PDFReaderImpl{}
}

func (p *PDFReaderImpl) ExtractText(file *domain.File) (string, error) {
	// Create a byte reader from the PDF data
	reader := bytes.NewReader(file.Content)

	// Open the PDF from the byte reader
	pdfReader, err := pdf.NewReader(reader, int64(len(file.Content)))
	if err != nil {
		log.Println("Error creating PDF reader:", err)
		return "", err
	}

	// Extract text from the first page
	page := pdfReader.Page(1)
	if page.V.IsNull() {
		log.Println("Page not found")
		return "", err
	}

	text, err := page.GetPlainText(nil)
	if err != nil {
		log.Println("Error extracting text:", err)
		return "", err
	}

	log.Println("Extracted text:", text)
	return text, nil
}

func (p *PDFReaderImpl) ConvertTextToPdf(text string) ([]byte, error) {
	pdfWriter := fpdf.New("P", "cm", "A4", "")
	pdfWriter.AddPage()
	// Enable swedish letter encoding
	pdfWriter.AddUTF8Font("Roboto", "", "./fonts/Roboto-Regular.ttf")
	pdfWriter.SetFont("Roboto", "", 11) // Set font size to 11

	pdfWriter.SetLeftMargin(2.54)
	pdfWriter.SetRightMargin(2.54)
	pdfWriter.SetTopMargin(2.54)
	pdfWriter.SetAutoPageBreak(true, 2.54)

	pdfWriter.Cell(0, 0, "Name")
	pdfWriter.Ln(0.5)
	pdfWriter.Cell(0, 0, "Adress")
	pdfWriter.Ln(0.5)
	pdfWriter.Cell(0, 0, "User Postal")
	pdfWriter.Ln(1)

	paragraphs := strings.Split(text, "\n\n")
	for _, paragraph := range paragraphs {
		trimmedText := strings.TrimSpace(paragraph)
		lineHeight := 0.5 // Approximation for 1.15 times the font height in cm (11 pt * 0.0352778 cm/pt * 1.15)
		pdfWriter.MultiCell(0, lineHeight, trimmedText, "", "L", false)
		pdfWriter.Ln(lineHeight)
	}

	buf := new(bytes.Buffer)
	err := pdfWriter.Output(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
