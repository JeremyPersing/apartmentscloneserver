package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var BucketName = "apartments-clone"
var bucketURL = "https://" + BucketName + ".s3.us-east-2.amazonaws.com/"

func InitializeS3() {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	customProvider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(customProvider),
		config.WithRegion("us-east-2"),
	)
	if err != nil {
		log.Fatal(err)
	}

	S3Client = s3.NewFromConfig(cfg)
}

func UploadBase64Image(base64ImageSrc string, name string) map[string]string {
	i := strings.Index(base64ImageSrc, ",")

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader((base64ImageSrc[i+1:])))

	url := bucketURL + name

	uploader := manager.NewUploader(S3Client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &BucketName,
		Key:    &name,
		Body:   decoder,
	})

	if err != nil {
		fmt.Print("some error happened", err)
	}

	urlMap := make(map[string]string)
	urlMap["url"] = url

	return urlMap
}
