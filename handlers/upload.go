package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Presigner encapsulates the Amazon Simple Storage Service (Amazon S3) presign actions
// used in the examples.
// It contains PresignClient, a client that is used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.
type Presigner struct {
	PresignClient *s3.PresignClient
}

func GetSignedURL(c *gin.Context) {
	// Initialize AWS session
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			},
		}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load AWS config"})
		return
	}

	// Create S3 service client
	svc := s3.NewFromConfig(cfg)

	// Generate a unique file name
	menuID := c.Param("menuId")
	fileName := generateUniqueFileName(menuID)

	// Set the expiration time for the signed URL
	expiration := time.Now().Add(15 * time.Minute)

	// Set the S3 bucket name and key
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	objectKey := "uploads/" + menuID + "/" + fileName

	// Create a Presigner instance
	presigner := Presigner{
		PresignClient: s3.NewPresignClient(svc),
	}

	// Generate the presigned URL using PutObject method
	presignedReq, err := presigner.PutObject(bucketName, objectKey, int64(time.Until(expiration).Seconds()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate signed URL"})
		return
	}

	// Return the signed URL to the client
	c.JSON(http.StatusOK, gin.H{"url": presignedReq.URL})
}

// Helper function to generate a unique file name
func generateUniqueFileName(menuID string) string {
	// Implement your logic to generate a unique file name here
	// You can use the menuID variable to generate the file name
	return menuID + "-" + uuid.New().String()
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner Presigner) PutObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}
