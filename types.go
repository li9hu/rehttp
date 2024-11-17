package rehttp

type Result struct {
	StatusCode   int     `json:"statusCode,omitempty"`
	URL          string  `json:"URL,omitempty"`
	Duration     float64 `json:"duration,omitempty"`
	ResponseBody string  `json:"responseBody,omitempty"`
	Err          error   `json:"err,omitempty"`
}

//type Request struct {
//	Data   []byte
//	URL    string
//	Header map[string]string
//	Method string
//}
