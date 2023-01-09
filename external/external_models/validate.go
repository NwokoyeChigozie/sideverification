package external_models

type ValidateOnDBReq struct {
	Table string      `json:"table"`
	Type  string      `json:"type"`
	Query string      `json:"query"`
	Value interface{} `json:"value"`
}

type ValidateOnDBReqModel struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    bool   `json:"data"`
}
