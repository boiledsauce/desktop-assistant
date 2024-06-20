package usecase

import (
	"context"
	"desktop-assistant/infra/repository"
	"desktop-assistant/infra/repository/api"
	"desktop-assistant/infra/service"
	"log"
	"mime/multipart"
)

type CoverLetterGenerator interface {
	// GeneratePersonalLetter generates a personal letter based on the given file
	GeneratePersonalLetter(file *multipart.FileHeader, jobDescription, additionalInfo string) ([]byte, error)
}

// CoverLetterGeneratorUseCase represents the use case for generating cover letters.
type CoverLetterGeneratorUseCase struct {
	// Add fields here to represent the dependencies of the use case.
	llmClient api.AIClient
	fileRepo  repository.FileSystemRepository
	pdfReader service.PDFReader
}

// NewCoverLetterGeneratorUseCase creates a new CoverLetterGeneratorUseCase.
func NewCoverLetterGeneratorUseCase(llmClient api.AIClient, fileRepo repository.FileSystemRepository, pdfReader service.PDFReader) *CoverLetterGeneratorUseCase {
	return &CoverLetterGeneratorUseCase{
		llmClient: llmClient,
		fileRepo:  fileRepo,
		pdfReader: pdfReader,
	}
}

// GeneratePersonalLetter generates a personal letter based on the given file.
func (uc *CoverLetterGeneratorUseCase) GeneratePersonalLetter(file *multipart.FileHeader, jobDescription, additionalInfo string) ([]byte, error) {
	// Read the file contents from the form file
	_file, err := file.Open()
	if err != nil {
		log.Print("Error reading file: ", err)
		return nil, err
	}
	defer _file.Close()

	// Convert the FormFile to a File entity with the file repository
	fileEntity, err := uc.fileRepo.GetFileEntity(_file)
	if err != nil {
		log.Print("Error getting file entity: ", err)
		return nil, err
	}

	// Read the PDF file that your CV should be in
	CV, err := uc.pdfReader.ExtractText(fileEntity)
	if err != nil {
		log.Print("Error extracting text from PDF: ", err)
		return nil, err
	}

	// Combine the CV and job description to generate the personal letter
	content := "CV: " + CV + "\n JobDescription: " + jobDescription + "\n additional info: " + additionalInfo
	// AI context for generating the personal letter
	// aiContext := `You are a great cover-letter writer.
	// 	You follow all recommended guidelines in cover letter writing for software engineering jobs,
	// 	and write a compelling tailored honest letter without sounding overly fake / ai generated.
	// 	If the requirements that are listed in the job description are not met with experiences, you should not lie about having experience in them.
	// 	Instead you should mention that you are willing to learn and adapt to new technologies and that you are a fast learner.`
	// 	aiContext := `You are an expert cover letter writer for software engineering positions.
	// You adhere strictly to best practices in cover letter composition, ensuring each letter is engaging, personalized, and authentic, avoiding any semblance of being artificially generated.
	// If the job description lists requirements not met by existing experiences, do not fabricate experiences. Instead, emphasize a strong willingness to learn new technologies and adapt quickly, supported by examples demonstrating quick learning and adaptability in past roles.`
	aiContext := `You are a skilled cover letter writer specializing in software engineering applications. Your task is to compose a concise yet impactful cover letter that adheres to the highest standards of professionalism and personalization, omitting any placeholders for personal information or addresses, as the cover letter is intended for electronic submission where such details are not required.

Begin with a brief introduction that mentions how the applicant discovered the position and expresses interest in the role. The introduction should directly lead into the body of the letter.

In the body, connect the applicant's skills and experiences directly to the requirements outlined in the job description. Emphasize the applicant's ability to adapt to new technologies and environments, using general examples from past projects or roles to highlight these skills. Avoid mentioning specific metrics unless they are generic and do not imply fabrication.

Convey genuine enthusiasm for the role and the company, focusing on the applicant's alignment with the company's values and interest in its mission or projects. This should be expressed without referencing any specific, recent news unless it is well-known and directly relevant to the job.

Advise on professional formatting and readability, ensuring the letter is visually appealing and matches the tone appropriate to the company's culture. Specifically, the AI should avoid creating sections for filling in personal details such as name, address, or contact information.

Conclude with a call to action that expresses eagerness to discuss the role further in an interview, stating that the applicant looks forward to an opportunity to further discuss their qualifications.

Throughout the letter, maintain a professional yet personable tone, showcasing the applicant's communication skills alongside their technical abilities. Keep the letter concise, ideally no longer than one page, and use clear, direct language.
`

	log.Println("Message: ", content)

	// Generate the personal letter using the AI model
	generatedCvInText, err := uc.llmClient.GenerateText(context.Background(), content, aiContext)
	if err != nil {
		log.Print("Error generating text: ", err)
		return nil, err
	}

	generatedCvInPdfBytes, err := uc.pdfReader.ConvertTextToPdf(generatedCvInText)
	if err != nil {
		log.Print("Error converting text to PDF: ", err)
		return nil, err
	}

	return generatedCvInPdfBytes, nil
}
