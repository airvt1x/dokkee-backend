package service

import "regexp"

type nerPattern struct {
	re    *regexp.Regexp
	label string
}

type NERService struct {
	patterns []nerPattern
}

func NewNERService() *NERService {
	return &NERService{patterns: []nerPattern{
		// ФИО: Иванов Иван Иванович
		{regexp.MustCompile(`[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+`), "[ФИО]"},
		// Телефон: +7(999)123-45-67, 8 999 123 45 67
		{regexp.MustCompile(`(\+7|8)[\s\-]?\(?\d{3}\)?[\s\-]?\d{3}[\s\-]?\d{2}[\s\-]?\d{2}`), "[ТЕЛЕФОН]"},
		// Email
		{regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`), "[EMAIL]"},
		// ИНН (10 или 12 цифр)
		{regexp.MustCompile(`\bИНН\s*:?\s*\d{10,12}\b`), "[ИНН]"},
		// Паспорт: серия и номер (4 + 6 цифр)
		{regexp.MustCompile(`\b\d{4}\s\d{6}\b`), "[ПАСПОРТ]"},
		// Банковский счёт (20 цифр)
		{regexp.MustCompile(`\b\d{20}\b`), "[СЧЁТ]"},
	}}
}

func (n *NERService) Anonymize(text string) string {
	result := text
	for _, p := range n.patterns {
		result = p.re.ReplaceAllString(result, p.label)
	}
	return result
}
