package main

import (
    "compress/gzip"
    "context"
    "encoding/csv"
    "encoding/json"
    "filepath"
    "fmt"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/satori/go.uuid"
    "log"
    "os"
    "regexp"
)

type ELBLog struct {
    Timestamp string `json:"timestamp"`
    ELBName string `json:"elb_name"`
    ClientIP string `json:"client_ip"`
    BackendIP string `json:"backend_ip"`
    RequestProcessingTime string `json:"request_processing_time"`
    BackendProcessingTime string `json:"backend_processing_time"`
    ResponseProcessingTime string `json:"response_processing_time"`
    ELBStatusCode string `json:"elb_status_code"`
    BackendStatusCode string `json:"backend_status_code"`
    ReceivedBytes string `json:"received_bytes"`
    SentBytes string `json:"sent_bytes"`
    Method string `json:"method"`
    Url string `json:"url"`
    UserAgent string `json:"user_agent"`
    SSLCipher string `json:"ssl_cipher"`
    SSLProtocol string `json:"ssl_protocol"`
}

type ALBLog struct {
    Type string `json:"type"`
    Timestamp string `json:"timestamp"`
    ELBId string `json:"elb_id"`
    ClientIP string `json:"client_ip"`
    ClientPort string `json:"client_port"`
    TargetIP string `json:"target_ip"`
    TargetPort string `json:"target_port"`
    RequestProcessingTime string `json:"request_processing_time"`
    TargetProcessingTime string `json:"target_processing_time"`
    ResponseProcessingTime string `json:"response_processing_time"`
    ELBStatusCode string `json:"elb_status_code"`
    TargetStatusCode string `json:"target_status_code"`
    ReceivedBytes string `json:"received_bytes"`
    SentBytes string `json:"sent_bytes"`
    Method string `json:"method"`
    Url string `json:"url"`
    HTTPVersion string `json:"http_version"`
    UserAgent string `json:"user_agent"`
    SSLCipher string `json:"ssl_cipher"`
    SSLProtocol string `json:"ssl_protocol"`
    TargetGroupArn string `json:"target_group_arn"`
    TraceId string `json:"trace_id"`
}

type CloudfrontLog struct {
    Timestamp string `json:"timestamp"`
    XEdgeLocation string `json:"x_edge_location"`
    SCBytes string `json:"sc_bytes"`
    CIP string `json:"c_ip"`
    CSMethod string `json:"cs_method"`
    CSHost string `json:"cs_host"`
    CSUriStem string `json:"cs_uri_stem"`
    SCStatus string `json:"sc_status"`
    CSReferer string `json:"cs_referer"`
    CSUserAgent string `json:"cs_user_agent"`
    CSUriQuery string `json:"cs_uri_query"`
    CSCookie string `json:"cs_cookie"`
    XEdgeResultType string `json:"x_edge_result_type"`
    XEdgeRequestId string `json:"x_edge_request_id"`
    XHostHeader string `json:"x_host_header"`
    CSProtocol string `json:"cs_protocol"`
    CSBytes string `json:"cs_bytes"`
    TimeTaken string `json:"time_taken"`
    XForwardedFor string `json:"x_forwarded_for"`
    SSLProtocol string `json:"ssl_protocol"`
    SSLCipher string `json:"ssl_cipher"`
    XEdgeResponseResultType string `json:"x_edge_response_result_type"`
}

var(
    REGION = os.Getenv('REGION')
    ENDPOINT = fmt.Printf("%s.data.logs.insight.rapid7.com", REGION)
    PORT = 20000
    validELBLog = regexp.MustCompile("\d+_\w+_\w{2}-\w{4,9}-[12]_.*._\d{8}T\d{4}Z_\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}_.*.log$")
    validALBLog = regexp.MustCompile("\d+_\w+_\w{2}-\w{4,9}-[12]_.*._\d{8}T\d{4}Z_\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}_.*.log.gz$")
    validCloudfrontLog = regexp.MustCompile("\w+\.\d{4}-\d{2}-\d{2}-\d{2}\.\w+\.gz$")
)

func ValidateUUID(uuid string):
    u2, err := uuid.FromString(uuid)
    if err != nil {
        fmt.Printf("Cannot validate token: %s", err)
        return false
    }
    return true
}


func ValidateELBLog(key string) {
    return validELBLog.MatchString(key)
}

func ValidateALBLog(key string) {
    return validALBLog.MatchString(key)
}

func ValidateCloudfrontLog(key string) {
    return validCloudfrontLog.MatchString(key)
}

func ParseELBUrl(input string) {
    request = strings.Split(input[11], " ")

    idx = request[1].find('/', 9)
    url = request[1][idx:]
}


func Handler(ctx context.Context, s3Event events.S3Event) (Response, error) {
    service := s3.New(session.New())
    var message string
    for _, record := range s3Event.Records {
        entity := record.S3
        object, err := service.GetObject(&s3.GetObjectInput {
            Bucket: entity.Bucket.Name,
            Key: entity.Object.Key,
        })
        if err != nil {
            log.Error(err)
        }
        reader := nil
        if filepath.Ext(entity.Object.Key) == 'gz' {
            reader, err = gzip.NewReader(object.Body)
            if err != nil {
                log.Error(err)
            }
        } else {
            reader = object.Body
        }

        if ValidateELBLog(entity.Object.Key) {
            csvReader := csv.NewReader(reader)
            for {
                line, error := reader.Read()
                if error == io.EOF {
                    break
                } else if error != nil {
                    log.Fatal(error)
                }

                message = ELBLog{
                    Timestamp: line[0],
                    ELBName: line[1],
                    ClientIP: strings.Split(line[2], ":")[0],
                    BackendIP: strings.Split(line[3], ":")[0],
                    RequestProcessingTime: line[4],
                    BackendProcessingTime: line[5],
                    ResponseProcessingTime: line[6],
                    ELBStatusCode: line[7],
                    BackendStatusCode: line[8],
                    ReceivedBytes: line[9],
                    SentBytes: line[10],
                    Method: line[11],
                    Url: ParseELBUrl(line[11]),
                    UserAgent: line[12],
                    SSLCipher: line[13],
                    SSLProtocol: line[14],
                })
            }

        } else if ValidateALBLog(entity.Object.Key) {
            csvReader := csv.NewReader(reader)
            for {
                line, error := reader.Read()
                if error == io.EOF {
                    break
                } else if error != nil {
                    log.Fatal(error)
                }

                message = ALBLog{
                    Type: line[0],
                    Timestamp: line[1],
                    ELBId: line[2],
                    ClientIP: strings.Split(line[3], ":")[0],
                    ClientPort: strings.Split(line[3], ":")[1],
                    TargetIP: strings.Split(line[4], ":")[0],
                    TargetPort: strings.Split(line[4], ":")[1],
                    RequestProcessingTime: line[5],
                    BackendProcessingTime: line[6],
                    ResponseProcessingTime: line[7],
                    ELBStatusCode: line[8],
                    TargetStatusCode: line[9],
                    ReceivedBytes: line[10],
                    SentBytes: line[11],
                    Method: line[11],
                    Url: ParseELBUrl(line[11]),
                    HTTPVersion: line[11],
                    UserAgent: line[13],
                    SSLCipher: line[14],
                    SSLProtocol: line[15],
                    TargetGroupArn: line[16],
                    TraceId: line[17],
                })
            }

        } else if ValidateCloudfrontLog(entity.Object.Key) {
            csvReader := csv.NewReader(reader)
            for {
                line, error := reader.Read()
                if error == io.EOF {
                    break
                } else if error != nil {
                    log.Fatal(error)
                }

                if len(line) < 23 {
                    continue
                }

                message = CloudfrontLog{
                    Timestamp: fmt.Printf("%sT%sZ", line[0], line[1]),
                    XEdgeLocation: line[2],
                    SCBytes: line[3],
                    CIP: line[4],
                    CSMethod: line[5],
                    CSHost: line[6],
                    CSUriStem: line[7],
                    SCStatus: line[8],
                    CSReferer: line[9],
                    CSUserAgent: line[10],
                    CSUriQuery: line[11],
                    CSCookie: line[12],
                    XEdgeResultType: line[13],
                    XEdgeRequestId: line[14],
                    XHostHeader: line[15],
                    CSProtocol: line[16],
                    CSBytes: line[17],
                    TimeTaken: line[18],
                    XForwardedFor: line[19],
                    SSLProtocol: line[20],
                    SSLCipher: line[21],
                    XEdgeResponseResultType: line[22],
                })
            }
        } else {

        }
        msg, err := json.Marshal(message)
        if err != nil {
            fmt.Println(err)
            return
        }




    }
}

func main() {
    lambda.Start(Handler)
}
