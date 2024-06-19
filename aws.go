package main

import (
	"context"
	// "fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var client *s3.Client

func initS3Client() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client = s3.NewFromConfig(cfg)
}

func batchDeleteObjects(bucketName string, objectKeys []string, batchSize int) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}

	for i := 0; i < len(objectIds); i += batchSize {
		end := i + batchSize
		if end > len(objectIds) {
			end = len(objectIds)
		}
		// fmt.Println(i, end)
		output, err := client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{Objects: objectIds[i:end]},
		})
		if err != nil {
			log.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", bucketName, err)
			return err
		} else {
			log.Printf("Deleted %v objects.\n", len(output.Deleted))
		}
	}
	return nil
}
