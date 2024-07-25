package taskstruct

type Task struct {
	Method  string            `json:"method" binding:"required"`
	URL     string            `json:"url" binding:"required"`
	Headers map[string]string `json:"headers" `
	Status  *TaskStatus       `json:"status"`
}

type TaskStatus struct {
	ID             string            `json:"id" binding:"required"`
	Status         string            `json:"status"`
	HTTPStatusCode int               `json:"httpStatusCode"`
	Headers        map[string]string `json:"headers"`
	Length         int               `json:"length"`
}
