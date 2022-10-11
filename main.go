package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"github.com/joho/godotenv"
)

// Copies the item in the bucket to another bucket.

func main() {
	if len(os.Args) != 4 {
		exitErrorf("Bucket, item, and other bucket names required\nUsage: go run s3_copy_object bucket item other-bucket")
	}

	bucket := os.Args[1]
	item := os.Args[2]
	other := os.Args[3]

	// load .env file
	err := godotenv.Load(".env")
	exitErrorf("Error loading .env file",err)
	

	source := bucket + "/" + item

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY"),os.Getenv("SECRET_KEY"), ""),
	})

	// Create S3 service client
	svc := s3.New(sess, &aws.Config{Endpoint: aws.String(os.Getenv("ENDPOINT_URL"))})

	// Copy the item
	_, err = svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(other),
		CopySource: aws.String(source), Key: aws.String(item)})
	if err != nil {
		fmt.Println(err)
		exitErrorf("Unable to copy item from bucket %q to bucket %q, %v", bucket, other, err)
	}

	// Wait to see if the item got copied
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(other), Key: aws.String(item)})
	if err != nil {
		exitErrorf("Error occurred while waiting for item %q to be copied to bucket %q, %v", bucket, item, other, err)
	}

	fmt.Printf("Item %q successfully copied from bucket %q to bucket %q\n", item, bucket, other)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
  
