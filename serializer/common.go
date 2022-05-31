package serializer

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg,omitempty"`
	Error string      `json:"error,omitempty"`
}
