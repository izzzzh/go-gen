package types

type HelloReq struct {
	ID int `json:"id"`
}

type HelloResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
