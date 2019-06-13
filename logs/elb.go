package logs

type ELBLog struct {
	Timestamp             string `json:"timestamp" regexp:"[^ ]*"`
	_                     string `regexp:" "`
	ELBName               string `json:"elb_name" regexp:"[^ ]*"`
	_                     string `regexp:" "`
	ClientIP              string `json:"client_ip" regexp:"[^ ]*"`
	_                     string `regexp:":"`
	ClientPort            string `json:"client_port" regexp:"[0-9]*"`
	_                     string `regexp:" "`
	BackendIP             string `json:"backend_ip" regexp:"[^ ]*"`
	_                     string `regexp:":"`
	BackendPort           string `json:"backend_port" regexp:"[0-9]*"`
	_                     string `regexp:" "`
	RequestProcessingTime string `json:"request_processing_time" regexp:"[.0-9]*"`
	_                     string `regexp:" "`
	BackendProcessingTime string `json:"backend_processing_time" regexp:"[.0-9]*"`
	_                     string `regexp:" "`
	ClientResponseTime    string `json:"client_response_time" regexp:"[.0-9]*"`
	_                     string `regexp:" "`
	ELBResponseCode       string `json:"elb_response_code" regexp:"-|[0-9]*"`
	_                     string `regexp:" "`
	BackendResponseCode   string `json:"backend_response_code" regexp:"-|[0-9]*"`
	_                     string `regexp:" "`
	ReceivedBytes         string `json:"received_bytes" regexp:"[-0-9]*"`
	_                     string `regexp:" "`
	SentBytes             string `json:"sent_bytes" regexp:"[-0-9]*"`
	_                     string `regexp:" \"`
	Method                string `json:"method" regexp:"[^ ]*"`
	_                     string `regexp:" "`
	Url                   string `json:"url" regexp:"[^ ]*"`
	_                     string `regexp:" "`
	Protocol              string `json:"protocol" regexp:"- |[^ ]*"`
	_                     string `regexp:"\".*"`
}
