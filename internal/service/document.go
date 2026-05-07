package service

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"unicode/utf8"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/repository"
	"github.com/google/uuid"
)

const (
	maxFileSizeBytes = 10 * 1024 * 1024 // 10 MB
)

var allowedMimeTypes = map[string]bool{
	"application/pdf":                                                          true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/plain": true,
}

type DocumentService struct {
	repo       repository.Document
	resultRepo repository.Result
	s3         repository.FileStorage
	ner        *NERService
	aiConn     *AIConnector
}

func NewDocumentService(
	repo repository.Document,
	resultRepo repository.Result,
	s3 repository.FileStorage,
	ner *NERService,
	aiConn *AIConnector,
) *DocumentService {
	return &DocumentService{
		repo:       repo,
		resultRepo: resultRepo,
		s3:         s3,
		ner:        ner,
		aiConn:     aiConn,
	}
}

func (s *DocumentService) Upload(userID int, file multipart.File, header *multipart.FileHeader) (int, error) {
	if err := s.checkBalance(userID); err != nil {
		return 0, err
	}

	if header.Size > maxFileSizeBytes {
		return 0, fmt.Errorf("file too large: max %d MB", maxFileSizeBytes/1024/1024)
	}

	mimeType := header.Header.Get("Content-Type")
	if !allowedMimeTypes[mimeType] {
		return 0, errors.New("unsupported file type: allowed PDF, DOCX, TXT")
	}

	s3Key := fmt.Sprintf("documents/%s/%s", uuid.New().String(), header.Filename)

	if err := s.s3.Upload(s3Key, file, mimeType); err != nil {
		return 0, fmt.Errorf("failed to upload to S3: %w", err)
	}

	docID, err := s.repo.Create(dokkee.Document{
		UserID:       userID,
		OriginalName: header.Filename,
		S3Key:        s3Key,
		MimeType:     mimeType,
		FileSize:     header.Size,
		Status:       "queued",
	})
	if err != nil {
		return 0, err
	}

	go s.runAnalysisPipeline(docID, s3Key)

	return docID, nil
}

func (s *DocumentService) GetByID(docID, userID int) (dokkee.Document, error) {
	doc, err := s.repo.GetByID(docID, userID)
	if err != nil {
		return dokkee.Document{}, errors.New("document not found")
	}
	return doc, nil
}

func (s *DocumentService) List(userID int, status string) ([]dokkee.Document, error) {
	return s.repo.List(userID, status)
}

func (s *DocumentService) checkBalance(userID int) error {
	balance, err := s.repo.CheckBalance(userID)
	if err != nil {
		return fmt.Errorf("failed to check balance: %w", err)
	}
	if balance < 1 {
		return errors.New("insufficient balance")
	}
	return nil
}

func (s *DocumentService) runAnalysisPipeline(docID int, s3Key string) {
	s.repo.UpdateStatus(docID, "processing", "")

	fileBytes, err := s.s3.Download(s3Key)
	if err != nil {
		s.repo.UpdateStatus(docID, "failed", "s3 download error: "+err.Error())
		return
	}

	rawText, err := extractText(fileBytes)
	if err != nil {
		s.repo.UpdateStatus(docID, "failed", "text extraction error: "+err.Error())
		return
	}

	// Обязательный этап обезличивания (ФЗ-152)
	anonymizedText := s.ner.Anonymize(rawText)

	result, err := s.aiConn.Analyze(anonymizedText)
	if err != nil {
		s.repo.UpdateStatus(docID, "failed", "ai analysis error: "+err.Error())
		return
	}

	s.resultRepo.Save(dokkee.AnalysisResult{
		DocumentID: docID,
		ResultJSON: result,
		ModelUsed:  s.aiConn.ModelName(),
	})
	s.repo.UpdateStatus(docID, "completed", "")
}

func extractText(data []byte) (string, error) {
	if utf8.Valid(data) {
		return string(data), nil
	}

	// PDF: извлечь текстовые блоки между BT и ET
	if bytes.HasPrefix(data, []byte("%PDF")) {
		return extractPDFText(data), nil
	}

	return string(data), nil
}

func extractPDFText(data []byte) string {
	var sb strings.Builder
	content := string(data)

	i := 0
	for i < len(content) {
		btIdx := strings.Index(content[i:], "BT")
		if btIdx == -1 {
			break
		}
		btIdx += i
		etIdx := strings.Index(content[btIdx:], "ET")
		if etIdx == -1 {
			break
		}
		etIdx += btIdx

		block := content[btIdx:etIdx]
		for _, line := range strings.Split(block, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasSuffix(line, "Tj") || strings.HasSuffix(line, "TJ") {
				text := extractPDFString(line)
				if text != "" {
					sb.WriteString(text)
					sb.WriteString(" ")
				}
			}
		}
		i = etIdx + 2
	}

	return sb.String()
}

func extractPDFString(line string) string {
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	return line[start+1 : end]
}
