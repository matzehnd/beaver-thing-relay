package events

type CompletedActionEvent struct {
	ActionId string  `json:"actionId"`
	Success  bool    `json:"success"`
	Error    *string `json:"error,omitempty"`
}

type ActionEvent struct {
	ActionId string `json:"actionId"`
	Endpoint string `json:"endpoint"`
	Action   string `json:"action"`
}
