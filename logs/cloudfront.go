package logs

type CloudfrontLog struct {
	Date 							string `json:"date" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Time		               string `json:"time" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Location           		string `json:"location" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Bytes                 	string `json:"bytes" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	RequestIP               string `json:"request_ip" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Method           			string `json:"method" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Host                  	string `json:"host" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Uri               		string `json:"uri" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	RequestStatus           string `json:"request_status" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Referer               	string `json:"referer" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	UserAgent             	string `json:"user_agent" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	QueryString             string `json:"query_string" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Cookie       	         string `json:"cookie" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	ResultType         		string `json:"result_type" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	RequestId         		string `json:"request_id" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	HostHeader            	string `json:"host_header" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	RequestProtocol         string `json:"request_protocol" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	RequestBytes            string `json:"request_bytes" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	TimeTaken               string `json:"time_taken" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	XForwardedFor           string `json:"x_forwarded_for" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	SSLProtocol             string `json:"ssl_protocol" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	SSLCipher               string `json:"ssl_cipher" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	ResponseResultType 		string `json:"response_result_type" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	ProtocolVersion 			string `json:"protocol_version" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	Status 						string `json:"status" regexp:"[^ \t]+"`
	_								string `regexp:"\s+"`
	EncryptedFields			string `json:"encrypted_fields" regexp:"[^ \t]+"`
}
