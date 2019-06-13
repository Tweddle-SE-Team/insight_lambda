package logs

type CloudfrontLog struct {
	Date               string `json:"date" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Time               string `json:"time" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Location           string `json:"location" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Bytes              string `json:"bytes" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	RequestIP          string `json:"request_ip" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Method             string `json:"method" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Host               string `json:"host" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Uri                string `json:"uri" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	RequestStatus      string `json:"request_status" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Referer            string `json:"referer" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	UserAgent          string `json:"user_agent" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	QueryString        string `json:"query_string" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Cookie             string `json:"cookie" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	ResultType         string `json:"result_type" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	RequestId          string `json:"request_id" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	HostHeader         string `json:"host_header" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	RequestProtocol    string `json:"request_protocol" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	RequestBytes       string `json:"request_bytes" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	TimeTaken          string `json:"time_taken" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	XForwardedFor      string `json:"x_forwarded_for" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	SSLProtocol        string `json:"ssl_protocol" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	SSLCipher          string `json:"ssl_cipher" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	ResponseResultType string `json:"response_result_type" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	ProtocolVersion    string `json:"protocol_version" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	Status             string `json:"status" regexp:"[^ ]*"`
	_                  string `regexp:"\t+"`
	EncryptedFields    string `json:"encrypted_fields" regexp:"[^ ]*"`
}
