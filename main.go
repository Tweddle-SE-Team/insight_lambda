package main

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/bsphere/le_go"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ELBLog struct {
	Timestamp              string `json:"timestamp"`
	ELBName                string `json:"elb_name"`
	ClientIP               string `json:"client_ip"`
	ClientPort             string `json:"client_port"`
	BackendIP              string `json:"backend_ip"`
	BackendPort            string `json:"backend_port"`
	RequestProcessingTime  string `json:"request_processing_time"`
	BackendProcessingTime  string `json:"backend_processing_time"`
	ResponseProcessingTime string `json:"response_processing_time"`
	ELBResponseCode        string `json:"elb_response_code"`
	BackendStatusCode      string `json:"backend_status_code"`
	ReceivedBytes          string `json:"received_bytes"`
	SentBytes              string `json:"sent_bytes"`
	Method                 string `json:"method"`
	Url                    string `json:"url"`
	HTTPVersion            string `json:"http_version"`
	UserAgent              string `json:"user_agent"`
	SSLCipher              string `json:"ssl_cipher"`
	SSLProtocol            string `json:"ssl_protocol"`
}

type ALBLog struct {
	Type                   string `json:"type"`
	Timestamp              string `json:"timestamp"`
	ELBId                  string `json:"elb_id"`
	ClientIP               string `json:"client_ip"`
	ClientPort             string `json:"client_port"`
	TargetIP               string `json:"target_ip"`
	TargetPort             string `json:"target_port"`
	RequestProcessingTime  string `json:"request_processing_time"`
	TargetProcessingTime   string `json:"target_processing_time"`
	ResponseProcessingTime string `json:"response_processing_time"`
	ELBResponseCode        string `json:"elb_response_code"`
	TargetStatusCode       string `json:"target_status_code"`
	ReceivedBytes          string `json:"received_bytes"`
	SentBytes              string `json:"sent_bytes"`
	Method                 string `json:"method"`
	Url                    string `json:"url"`
	HTTPVersion            string `json:"http_version"`
	UserAgent              string `json:"user_agent"`
	SSLCipher              string `json:"ssl_cipher"`
	SSLProtocol            string `json:"ssl_protocol"`
	TargetGroupArn         string `json:"target_group_arn"`
	TraceId                string `json:"trace_id"`
}

type CloudfrontLog struct {
	Timestamp               string `json:"timestamp"`
	XEdgeLocation           string `json:"x_edge_location"`
	SCBytes                 string `json:"sc_bytes"`
	CIP                     string `json:"c_ip"`
	CSMethod                string `json:"cs_method"`
	CSHost                  string `json:"cs_host"`
	CSUriStem               string `json:"cs_uri_stem"`
	SCStatus                string `json:"sc_status"`
	CSReferer               string `json:"cs_referer"`
	CSUserAgent             string `json:"cs_user_agent"`
	CSUriQuery              string `json:"cs_uri_query"`
	CSCookie                string `json:"cs_cookie"`
	XEdgeResultType         string `json:"x_edge_result_type"`
	XEdgeRequestId          string `json:"x_edge_request_id"`
	XHostHeader             string `json:"x_host_header"`
	CSProtocol              string `json:"cs_protocol"`
	CSBytes                 string `json:"cs_bytes"`
	TimeTaken               string `json:"time_taken"`
	XForwardedFor           string `json:"x_forwarded_for"`
	SSLProtocol             string `json:"ssl_protocol"`
	SSLCipher               string `json:"ssl_cipher"`
	XEdgeResponseResultType string `json:"x_edge_response_result_type"`
}

var (
	validELBLog        *regexp.Regexp = regexp.MustCompile("\\d+_\\w+_\\w{2}-\\w{4,9}-[12]_.*._d{8}T\\d{4}Z_\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}_.*.log$")
	validALBLog        *regexp.Regexp = regexp.MustCompile("\\d+_\\w+_\\w{2}-\\w{4,9}-[12]_.*._\\d{8}T\\d{4}Z_\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}_.*.log.gz$")
	validCloudfrontLog *regexp.Regexp = regexp.MustCompile("\\w+\\.\\d{4}-\\d{2}-\\d{2}-\\d{2}\\.\\w+\\.gz$")
)

func ValidateELBLog(key string) bool {
	return validELBLog.MatchString(key)
}

func ValidateALBLog(key string) bool {
	return validALBLog.MatchString(key)
}

func ValidateCloudfrontLog(key string) bool {
	return validCloudfrontLog.MatchString(key)
}

func handler(ctx context.Context, s3Event events.S3Event) {
	session := session.New()
	s3Service := s3.New(session)
	ssmService := ssm.New(session)
	var logentriesToken string
	var line []string
	var msg []byte
	var ioReader io.Reader
	for _, record := range s3Event.Records {
		entity := record.S3
		object, err := s3Service.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(entity.Bucket.Name),
			Key:    aws.String(entity.Object.Key),
		})
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(entity.Object.Key) == ".gz" {
			ioReader, err = gzip.NewReader(object.Body)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			ioReader = object.Body
		}

		ssmPrefix := filepath.Dir(entity.Object.Key))

		if index := strings.Index(ssmPrefix, "/AWSLogs/"); index != -1 {
			ssmPrefix = ssmPrefix[:index]
		}

		resp, err := ssmService.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(fmt.Sprintf("/%s/error_logs_token", filepath.Dir(entity.Object.Key))),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatal(err)
		}
		logentriesToken = *resp.Parameter.Value

		le, err := le_go.Connect(logentriesToken)
		if err != nil {
			log.Fatal(err)
		}

		defer le.Close()

		if ValidateELBLog(entity.Object.Key) {
			csvReader := csv.NewReader(ioReader)
			csvReader.Comma = '\t'
			csvReader.Comment = '#'
			csvReader.LazyQuotes = true
			for {
				line, err = csvReader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}

				clientAddr := strings.Split(line[2], ":")
				clientIP, clientPort := clientAddr[0], clientAddr[1]

				backendAddr := strings.Split(line[3], ":")
				backendIP, backendPort := backendAddr[0], backendAddr[1]

				lbRequest, err := strconv.Unquote(line[12])
				if err != nil {
					log.Fatal(err)
				}
				lbRequestParts := strings.Split(lbRequest, " ")
				method, url, http_version := lbRequestParts[0], lbRequestParts[1], lbRequestParts[2]

				message := ELBLog{
					Timestamp:              line[0],
					ELBName:                line[1],
					ClientIP:               clientIP,
					ClientPort:             clientPort,
					BackendIP:              backendIP,
					BackendPort:            backendPort,
					RequestProcessingTime:  line[4],
					BackendProcessingTime:  line[5],
					ResponseProcessingTime: line[6],
					ELBResponseCode:        line[7],
					BackendStatusCode:      line[8],
					ReceivedBytes:          line[9],
					SentBytes:              line[10],
					Method:                 method,
					Url:                    url,
					HTTPVersion:            http_version,
					UserAgent:              line[12],
					SSLCipher:              line[13],
					SSLProtocol:            line[14],
				}
				msg, err = json.Marshal(message)
				if err != nil {
					fmt.Println(err)
				}
				le.Println(string(msg))
			}

		} else if ValidateALBLog(entity.Object.Key) {
			csvReader := csv.NewReader(ioReader)
			csvReader.Comma = '\t'
			csvReader.Comment = '#'
			csvReader.LazyQuotes = true
			for {
				line, err := csvReader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}

				clientAddr := strings.Split(line[3], ":")
				clientIP, clientPort := clientAddr[0], clientAddr[1]

				targetAddr := strings.Split(line[4], ":")
				targetIP, targetPort := targetAddr[0], targetAddr[1]

				lbRequest, err := strconv.Unquote(line[12])
				if err != nil {
					log.Fatal(err)
				}
				lbRequestParts := strings.Split(lbRequest, " ")
				method, url, http_version := lbRequestParts[0], lbRequestParts[1], lbRequestParts[2]

				message := ALBLog{
					Type:                   line[0],
					Timestamp:              line[1],
					ELBId:                  line[2],
					ClientIP:               clientIP,
					ClientPort:             clientPort,
					TargetIP:               targetIP,
					TargetPort:             targetPort,
					RequestProcessingTime:  line[5],
					TargetProcessingTime:   line[6],
					ResponseProcessingTime: line[7],
					ELBResponseCode:        line[8],
					TargetStatusCode:       line[9],
					ReceivedBytes:          line[10],
					SentBytes:              line[11],
					Method:                 method,
					Url:                    url,
					HTTPVersion:            http_version,
					UserAgent:              line[13],
					SSLCipher:              line[14],
					SSLProtocol:            line[15],
					TargetGroupArn:         line[16],
					TraceId:                line[17],
				}

				msg, err = json.Marshal(message)
				if err != nil {
					fmt.Println(err)
				}
				le.Println(string(msg))
			}

		} else if ValidateCloudfrontLog(entity.Object.Key) {
			csvReader := csv.NewReader(ioReader)
			csvReader.Comma = '\t'
			csvReader.Comment = '#'
			csvReader.LazyQuotes = true
			for {
				line, err := csvReader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}

				if len(line) < 23 {
					continue
				}

				message := CloudfrontLog{
					Timestamp:               fmt.Sprintf("%sT%sZ", line[0], line[1]),
					XEdgeLocation:           line[2],
					SCBytes:                 line[3],
					CIP:                     line[4],
					CSMethod:                line[5],
					CSHost:                  line[6],
					CSUriStem:               line[7],
					SCStatus:                line[8],
					CSReferer:               line[9],
					CSUserAgent:             line[10],
					CSUriQuery:              line[11],
					CSCookie:                line[12],
					XEdgeResultType:         line[13],
					XEdgeRequestId:          line[14],
					XHostHeader:             line[15],
					CSProtocol:              line[16],
					CSBytes:                 line[17],
					TimeTaken:               line[18],
					XForwardedFor:           line[19],
					SSLProtocol:             line[20],
					SSLCipher:               line[21],
					XEdgeResponseResultType: line[22],
				}
				msg, err = json.Marshal(message)
				if err != nil {
					fmt.Println(err)
				}
				le.Println(string(msg))
			}
		} else {

		}
	}
}

func readJsonFromFile(inputFile string) []byte {
	inputJson, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
	}
	return inputJson
}

func main() {
	//inputJson := readJsonFromFile("./s3-event.json")
	//var inputEvent events.S3Event
	//if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
	//	fmt.Println(err)
	//}
	//handler(nil, inputEvent)
	lambda.Start(handler)
}
