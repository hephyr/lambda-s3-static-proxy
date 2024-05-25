package main

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	path := strings.Trim(request.Path, "/")
	//path = strings.Replace(path, os.Getenv("STAGE"), "", 1)
	if path == "" {
		path = "index.html" // Default to index.html for root
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Error loading configuration: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Printf("Error retrieving %s: %s", path, err)
		// On error, retry with index.html
		path = "index.html"
		resp, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(path),
		})
		if err != nil {
			log.Printf("Error retrieving index.html: %s", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)
	isBase64Encoded := false

	contentType := "text/html" // Default content type
	if resp.ContentType != nil {
		contentType = *resp.ContentType
		if strings.Contains(contentType, "image") {
			log.Printf("Image detected: %s", path)
			body = base64.StdEncoding.EncodeToString(bodyBytes)
			isBase64Encoded = true
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		IsBase64Encoded: isBase64Encoded,
	}, nil
}

func main() {
	lambda.Start(handler)
}
