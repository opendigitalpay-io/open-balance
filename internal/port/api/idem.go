package api

type IdemStartRequest struct {
	IdemKey    string `json:"idemKey"`
}

type IdemEndRequest struct {
	IdemKey    string `json:"idemKey"`
	Response   string `json:"response"`
}