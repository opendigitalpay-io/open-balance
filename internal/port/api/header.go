package api

type IdemHeader struct {
	IdemKey string `header:"IdemKey" binding:"required"`
}
