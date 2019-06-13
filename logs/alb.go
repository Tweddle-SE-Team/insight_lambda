package logs

type ALBLog struct {
	Type                   string `json:"type" regexp:"[^ ]*"`
	_                      string `regexp:" "`
	Timestamp              string `json:"timestamp" regexp:"[^ ]*"`
	_                      string `regexp:" "`
	ELBId                  string `json:"elb_id" regexp:"[^ ]*"`
	_                      string `regexp:" "`
	ClientIP               string `json:"client_ip" regexp:"[^ ]*"`
	_                      string `regexp:":"`
	ClientPort             string `json:"client_port" regexp:"[0-9]*"`
	_                      string `regexp:" "`
	TargetIP               string `json:"target_ip" regexp:"[^ ]*"`
	_                      string `regexp:"[:-]"`
	TargetPort             string `json:"target_port" regexp:"[0-9]*"`
	_                      string `regexp:" "`
	RequestProcessingTime  string `json:"request_processing_time" regexp:"[-.0-9]*"`
	_                      string `regexp:" "`
	TargetProcessingTime   string `json:"target_processing_time" regexp:"[-.0-9]*"`
	_                      string `regexp:" "`
	ResponseProcessingTime string `json:"response_processing_time" regexp:"[-.0-9]*"`
	_                      string `regexp:" "`
	ELBResponseCode        string `json:"elb_response_code" regexp:"|[-0-9]*"`
	_                      string `regexp:" "`
	TargetStatusCode       string `json:"target_status_code" regexp:"-|[-0-9]*"`
	_                      string `regexp:" "`
	ReceivedBytes          string `json:"received_bytes" regexp:"[-0-9]*"`
	_                      string `regexp:" "`
	SentBytes              string `json:"sent_bytes" regexp:"[-0-9]*"`
	_                      string `regexp:" \""`
	Method                 string `json:"method" regexp:"[^ ]*"`
	_                      string `regexp:" "`
	Url                    string `json:"url" regexp:"[^ ]*"`
	_                      string `regexp:" "`
	HTTPVersion            string `json:"http_version" regexp:"- |[^ ]*"`
	_                      string `regexp:"\" \""`
	UserAgent              string `json:"user_agent" regexp:"[^\"]*"`
	_                      string `regexp:"\" "`
	SSLCipher              string `json:"ssl_cipher" regexp:"[A-Z0-9-]+"`
	_                      string `regexp:" "`
	SSLProtocol            string `json:"ssl_protocol" regexp:"[A-Za-z0-9.-]*"`
	_                      string `regexp:" "`
	TargetGroupArn         string `json:"target_group_arn" regexp:"[^ ]*"`
	_                      string `regexp:" \""`
	TraceId                string `json:"trace_id" regexp:"[^\"]*"`
	_                      string `regexp:"\" \""`
	DomainName             string `json:"domain_name" regexp:"[^\"]*"`
	_                      string `regexp:"\" \""`
	ChosenCertArn          string `json:"chosen_cert_arn" regexp:"[^\"]*"`
	_                      string `regexp:"\" "`
	MatchedRulePriority    string `json:"matched_rule_priority" regexp:"[-.0-9]*"`
	_                      string `regexp:" "`
	RequestCreationTime    string `json:"request_creation_time" regexp:"[^ ]*"`
	_                      string `regexp:" \""`
	ActionsExecuted        string `json:"actions_executed" regexp:"[^\"]*"`
	_                      string `regexp:"\" \""`
	RedirectUrl            string `json:"redirect_url" regexp:"[^\"]*"`
	_                      string `regexp:"\""`
	LambdaErrorReason      string `json:"lambda_error_reason" regexp:"$| \"[^ ]*\""`
	NewField               string `json:"new_field" regexp:".*"`
}
