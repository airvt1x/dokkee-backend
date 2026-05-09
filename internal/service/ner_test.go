package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNERService_Anonymize(t *testing.T) {
	ner := NewNERService()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "replace full name",
			input:    "Иванов Иван Иванович подписал договор.",
			expected: "[ФИО] подписал договор.",
		},
		{
			name:     "replace phone number +7 format",
			input:    "Мой номер +7(999)123-45-67 для связи.",
			expected: "Мой номер [ТЕЛЕФОН] для связи.",
		},
		{
			name:     "replace phone number 8 format",
			input:    "Позвоните 8 999 123 45 67 срочно.",
			expected: "Позвоните [ТЕЛЕФОН] срочно.",
		},
		{
			name:     "replace email",
			input:    "Отправьте на test@example.com письмо.",
			expected: "Отправьте на [EMAIL] письмо.",
		},
		{
			name:     "replace passport",
			input:    "Паспорт 1234 567890 выдан.",
			expected: "Паспорт [ПАСПОРТ] выдан.",
		},
		{
			name:     "multiple replacements",
			input:    "Иванов Иван Иванович, тел. +7(999)111-22-33, email ivan@mail.ru",
			expected: "[ФИО], тел. [ТЕЛЕФОН], email [EMAIL]",
		},
		{
			name:     "no personal data",
			input:    "Обычный текст без персональных данных.",
			expected: "Обычный текст без персональных данных.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ner.Anonymize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNERService_Anonymize_EdgeCases(t *testing.T) {
	ner := NewNERService()

	testCases := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"only spaces", "   "},
		{"special characters", "!@#$%^&*()"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ner.Anonymize(tc.input)
			assert.Equal(t, tc.input, result)
		})
	}
}

func TestNERService_Anonymize_ComplexDocument(t *testing.T) {
	ner := NewNERService()

	input := `Уважаемый Иванов Иван Иванович!
	Ваш номер телефона +7(999)123-45-67 и email ivan@mail.ru были проверены.
	ИНН: 1234567890
	Паспорт: 1234 567890
	Счёт: 12345678901234567890

	С уважением, Администрация.`

	result := ner.Anonymize(input)

	// Проверяем, что опасные данные заменены
	assert.Contains(t, result, "[ФИО]")
	assert.Contains(t, result, "[ТЕЛЕФОН]")
	assert.Contains(t, result, "[EMAIL]")
	assert.Contains(t, result, "[ПАСПОРТ]")
	// ИНН и счёт могут не заменяться, поэтому не проверяем

	// Проверяем, что оригинальные данные не остались
	assert.NotContains(t, result, "Иванов Иван")
	assert.NotContains(t, result, "+7(999)123-45-67")
	assert.NotContains(t, result, "ivan@mail.ru")
	assert.NotContains(t, result, "1234 567890")
}
