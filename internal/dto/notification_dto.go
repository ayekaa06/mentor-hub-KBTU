package dto

// SendReminderRequest — отправка напоминания ментору или студенту.
type SendReminderRequest struct {
	TargetID string `json:"target_id" validate:"required,uuid"`
	Message  string `json:"message"   validate:"required,min=5,max=500"`
}

// CreateQuestionRequest — вопрос freshman-а ментору.
type CreateQuestionRequest struct {
	Body string `json:"body" validate:"required,min=5,max=2000"`
}

// AnswerQuestionRequest — ответ ментора на вопрос.
type AnswerQuestionRequest struct {
	Answer string `json:"answer" validate:"required,min=1,max=2000"`
}

// MarkNotificationReadRequest — отметить уведомление прочитанным.
// Используется для bulk-операций (если нужно отметить несколько).
type MarkNotificationReadRequest struct {
	IDs []string `json:"ids" validate:"required,min=1,dive,uuid"`
}
