package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client wraps an S3 client configured for SeaweedFS.
type Client struct {
	s3Client *s3.Client
	bucket   string
}

// NewClient creates an S3 client pointing at a SeaweedFS endpoint.
func NewClient(endpoint, accessKey, secretKey, bucket string) (*Client, error) {
	s3Client := s3.New(s3.Options{
		BaseEndpoint: aws.String("http://" + endpoint),
		Region:       "us-east-1", // required by SDK but ignored by SeaweedFS
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		UsePathStyle: true, // SeaweedFS requires path-style
	})
	return &Client{s3Client: s3Client, bucket: bucket}, nil
}

// PutJSON marshals data as JSON and uploads it to the given key.
func (c *Client) PutJSON(ctx context.Context, key string, data any) error {
	body, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	contentType := "application/json"
	_, err = c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &c.bucket,
		Key:         &key,
		Body:        bytes.NewReader(body),
		ContentType: &contentType,
	})
	if err != nil {
		return fmt.Errorf("put object %s: %w", key, err)
	}
	return nil
}

// GetJSON downloads an object and unmarshals its JSON content into out.
func (c *Client) GetJSON(ctx context.Context, key string, out any) error {
	resp, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("get object %s: %w", key, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read object %s: %w", key, err)
	}
	return json.Unmarshal(body, out)
}
