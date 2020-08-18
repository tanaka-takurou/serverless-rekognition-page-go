package main

import (
	"os"
	"fmt"
	"log"
	"strings"
	"context"
	"encoding/json"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

const layout       string = "2006-01-02 15:04"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "detectmoderation" :
			if i, ok := d["image"]; ok {
				r, e := detectModeration(i)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detecttext" :
			if i, ok := d["image"]; ok {
				r, e := detectText(i)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectfaces" :
			if i, ok := d["image"]; ok {
				r, e := detectFaces(i)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectlabels" :
			if i, ok := d["image"]; ok {
				r, e := detectLabels(i)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: 500,
			Body: string(jsonBytes),
		}, nil
	}
	return Response {
		StatusCode: 200,
		Body: string(jsonBytes),
	}, nil
}

func detectModeration(img string)(string, error) {
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}
	svc := rekognition.New(session.New(), &aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})

	input := &rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
	}
	res, err2 := svc.DetectModerationLabels(input)
	if err2 != nil {
		return "", err2
	}
	if len(res.ModerationLabels) < 1 {
		return "No ModerationLabel", nil
	}
	results, err3 := json.Marshal(res.ModerationLabels)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func detectText(img string)(string, error) {
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}
	svc := rekognition.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
	}
	res, err2 := svc.DetectText(input)
	if err2 != nil {
		return "", err2
	}
	if len(res.TextDetections) < 1 {
		return "No TextDetection", nil
	}
	results, err3 := json.Marshal(res.TextDetections)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func detectFaces(img string)(string, error) {
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}
	svc := rekognition.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	input := &rekognition.DetectFacesInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
	}
	res, err2 := svc.DetectFaces(input)
	if err2 != nil {
		return "", err2
	}
	if len(res.FaceDetails) < 1 {
		return "No FaceDetails", nil
	}
	results, err3 := json.Marshal(res.FaceDetails)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func detectLabels(img string)(string, error) {
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}
	svc := rekognition.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
		MaxLabels: aws.Int64(10),
		MinConfidence: aws.Float64(60.0),
	}
	res, err2 := svc.DetectLabels(input)
	if err2 != nil {
		return "", err2
	}
	if len(res.Labels) < 1 {
		return "No Labels", nil
	}
	results, err3 := json.Marshal(res.Labels)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func main() {
	lambda.Start(HandleRequest)
}
