package main

import (
	"os"
	"fmt"
	"log"
	"strings"
	"context"
	"encoding/json"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
)

type APIResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

var cfg aws.Config
var rekognitionClient *rekognition.Client

const layout string = "2006-01-02 15:04"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "detectmoderation" :
			if i, ok := d["image"]; ok {
				r, e := detectModeration(ctx, i)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detecttext" :
			if i, ok := d["image"]; ok {
				r, e := detectText(ctx, i)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectfaces" :
			if i, ok := d["image"]; ok {
				r, e := detectFaces(ctx, i)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectlabels" :
			if i, ok := d["image"]; ok {
				r, e := detectLabels(ctx, i)
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

func detectModeration(ctx context.Context, img string)(string, error) {
	if rekognitionClient == nil {
		rekognitionClient = rekognition.New(cfg)
	}
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}

	input := &rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
	}
	req := rekognitionClient.DetectModerationLabelsRequest(input)
	res, err2 := req.Send(ctx)
	if err2 != nil {
		return "", err2
	}
	if len(res.DetectModerationLabelsOutput.ModerationLabels) < 1 {
		return "No ModerationLabel", nil
	}
	results, err3 := json.Marshal(res.DetectModerationLabelsOutput.ModerationLabels)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func detectText(ctx context.Context, img string)(string, error) {
	if rekognitionClient == nil {
		rekognitionClient = rekognition.New(cfg)
	}
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}

	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
	}
	req := rekognitionClient.DetectTextRequest(input)
	res, err2 := req.Send(ctx)
	if err2 != nil {
		return "", err2
	}
	if len(res.DetectTextOutput.TextDetections) < 1 {
		return "No TextDetection", nil
	}
	results, err3 := json.Marshal(res.DetectTextOutput.TextDetections)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func detectFaces(ctx context.Context, img string)(string, error) {
	if rekognitionClient == nil {
		rekognitionClient = rekognition.New(cfg)
	}
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}

	input := &rekognition.DetectFacesInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
	}
	req := rekognitionClient.DetectFacesRequest(input)
	res, err2 := req.Send(ctx)
	if err2 != nil {
		return "", err2
	}
	if len(res.DetectFacesOutput.FaceDetails) < 1 {
		return "No FaceDetails", nil
	}
	results, err3 := json.Marshal(res.DetectFacesOutput.FaceDetails)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func detectLabels(ctx context.Context, img string)(string, error) {
	if rekognitionClient == nil {
		rekognitionClient = rekognition.New(cfg)
	}
	b64data := img[strings.IndexByte(img, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}

	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: data,
		},
		MaxLabels: aws.Int64(10),
		MinConfidence: aws.Float64(60.0),
	}
	req := rekognitionClient.DetectLabelsRequest(input)
	res, err2 := req.Send(ctx)
	if err2 != nil {
		return "", err2
	}
	if len(res.DetectLabelsOutput.Labels) < 1 {
		return "No Labels", nil
	}
	results, err3 := json.Marshal(res.DetectLabelsOutput.Labels)
	if err3 != nil {
		return "", err3
	}
	return string(results), nil
}

func init() {
	var err error
	cfg, err = external.LoadDefaultAWSConfig()
	cfg.Region = os.Getenv("REGION")
	if err != nil {
		log.Print(err)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
