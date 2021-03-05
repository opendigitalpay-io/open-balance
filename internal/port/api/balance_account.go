package api

type GetBalanceAccountResponse struct {
	Balance  int64                  `json:"balance"`
	Currency string                 `json:"currency"`
	Type     string                 `json:"type"`
	State    string                 `json:"state"`
	Metadata map[string]interface{} `json:"metadata"`
}
