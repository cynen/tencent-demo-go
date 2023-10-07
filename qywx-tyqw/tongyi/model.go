package tongyi

// =============== 通义千问查询结构体 ==========================
// resp body
// =========================================================

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type HistoryMessage struct {
	Role    string `json:"role"` // user/bot
	Content string `json:"content"`
}
type History struct {
}
type Parameters struct {
	ResultFormat string  `json:"result_format"`
	TopP         float64 `json:"top_p"`
	TopK         float64 `json:"top_k"`
	Seed         int32   `json:"seed"`
	Temperature  float64 `json:"temperature"`
	EnableSearch bool    `json:"enable_search"`
}
type Input struct {
	Messages []Message `json:"messages"`
	History  []History `json:"history"`
	Prompt   string    `json:"prompt"`
}

type ReqBody struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

// =============== 通义千问响应结构体 ==========================
// resp body
// =========================================================

type Resp struct {
	Output    OutPut `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestId string `json:"request_id"`
}

type OutPut struct {
	Text         string   `json:"text"`
	FinishReason string   `json:"finish_reason"`
	Choices      []Choice `json:"choices"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type Usage struct {
	OutputTokens int `json:"output_tokens"`
	InputTokens  int `json:"input_tokens"`
}
