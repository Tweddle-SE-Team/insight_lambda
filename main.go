package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tweddle-SE-Team/insight_go"
	"github.com/Tweddle-SE-Team/insight_goclient"
	"github.com/Tweddle-SE-Team/insight_lambda/logs"
	"github.com/alexflint/go-restructure"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	CLOUDFRONT_LOG_REGEXP  = "\\w+\\.\\d{4}-\\d{2}-\\d{2}-\\d{2}\\.\\w+\\.gz$"
	ALB_LOG_REGEXP         = "\\d+_\\w+_\\w{2}-\\w{4,9}-[12]_.*._\\d{8}T\\d{4}Z_\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}_.*.log.gz$"
	ELB_LOG_REGEXP         = "\\d+_\\w+_\\w{2}-\\w{4,9}-[12]_.*._\\d{8}T\\d{4}Z_\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}_.*.log$"
	INSIGHT_LOG_NAME_ENV   = "INSIGHT_LOG_NAME"
	INSIGHT_API_KEY_ENV    = "INSIGHT_API_KEY"
	INSIGHT_API_REGION_ENV = "INSIGHT_API_REGION"
)

func GetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatal(fmt.Errorf("%s environment variable is not set", name))
	}
	return value
}

func Handler(ctx context.Context, s3Event events.S3Event) {
	insightLogName := GetEnv(INSIGHT_LOG_NAME_ENV)
	insightApiKey := GetEnv(INSIGHT_API_KEY_ENV)
	insightApiRegion := GetEnv(INSIGHT_API_REGION_ENV)

	cloudfrontLogRegexp := regexp.MustCompile(CLOUDFRONT_LOG_REGEXP)
	albLogRegexp := regexp.MustCompile(ALB_LOG_REGEXP)
	elbLogRegexp := regexp.MustCompile(ELB_LOG_REGEXP)

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(cfg)
	adminClient, err := insight_goclient.NewInsightClient(insightApiKey, insightApiRegion)
	if err != nil {
		log.Fatal(err)
	}
	var ioReader io.Reader
	for _, record := range s3Event.Records {
		entity := record.S3
		request := client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(entity.Bucket.Name),
			Key:    aws.String(entity.Object.Key),
		})
		s3response, err := request.Send(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if filepath.Ext(entity.Object.Key) == ".gz" {
			ioReader, err = gzip.NewReader(s3response.Body)
			if err != nil {
				log.Println(err)
				continue
			}
		} else {
			ioReader = s3response.Body
		}

		logsetName := strings.Replace(filepath.Dir(entity.Object.Key), "/", "-", -1)
		if index := strings.Index(logsetName, "-AWSLogs-"); index != -1 {
			logsetName = logsetName[:index]
		}

		insightToken, err := adminClient.GetLogToken(logsetName, insightLogName)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Found token for %s logset", logsetName)

		insight, err := insight_go.Connect(insightApiRegion, insightToken)
		if err != nil {
			log.Println(err)
			continue
		}

		defer insight.Close()

		scanner := bufio.NewScanner(ioReader)

		switch {
		case albLogRegexp.MatchString(entity.Object.Key):
			logEntry := logs.ALBLog{}
			for scanner.Scan() {
				restructure.Find(&logEntry, scanner.Text())
			}
			msg, err := json.Marshal(logEntry)
			if err != nil {
				log.Println(err)
			}
			insight.Println(string(msg))
		case elbLogRegexp.MatchString(entity.Object.Key):
			logEntry := logs.ELBLog{}
			for scanner.Scan() {
				restructure.Find(&logEntry, scanner.Text())
			}
			msg, err := json.Marshal(logEntry)
			if err != nil {
				log.Println(err)
			}
			insight.Println(string(msg))
		case cloudfrontLogRegexp.MatchString(entity.Object.Key):
			logEntry := logs.CloudfrontLog{}
			for scanner.Scan() {
				restructure.Find(&logEntry, scanner.Text())
			}
			msg, err := json.Marshal(logEntry)
			if err != nil {
				log.Println(err)
			}
			insight.Println(string(msg))
		default:
			log.Println("Couldn't parse log file")
			continue
		}

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func readJsonFromFile(inputFile string) []byte {
	inputJson, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Println(err)
	}
	return inputJson
}

func main() {
	//inputJson := readJsonFromFile("./s3-event.json")
	//var inputEvent events.S3Event
	//if err := json.Unmarshal(inputJson, &inputEvent); err != nil {
	//	log.Println(err)
	//}
	//Handler(context.TODO(), inputEvent)
	lambda.Start(Handler)
}
