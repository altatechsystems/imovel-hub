package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCSClient handles Google Cloud Storage operations
type GCSClient struct {
	client     *storage.Client
	bucketName string
}

// NewGCSClient creates a new GCS client
func NewGCSClient(ctx context.Context, projectID, bucketName, credentialsFile string) (*GCSClient, error) {
	var opts []option.ClientOption
	if credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCSClient{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Upload uploads data to GCS and returns the public URL
func (c *GCSClient) Upload(ctx context.Context, objectPath string, data []byte, contentType string) (string, error) {
	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(objectPath)

	// Create writer
	w := obj.NewWriter(ctx)
	w.ContentType = contentType
	w.CacheControl = "public, max-age=31536000" // Cache for 1 year

	// Write data
	if _, err := w.Write(data); err != nil {
		w.Close()
		return "", fmt.Errorf("failed to write data: %w", err)
	}

	// Close writer (this commits the upload)
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Make object publicly readable
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set ACL: %w", err)
	}

	// Return public URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectPath)
	return url, nil
}

// UploadFromReader uploads data from an io.Reader to GCS
func (c *GCSClient) UploadFromReader(ctx context.Context, objectPath string, reader io.Reader, contentType string) (string, error) {
	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(objectPath)

	// Create writer
	w := obj.NewWriter(ctx)
	w.ContentType = contentType
	w.CacheControl = "public, max-age=31536000"

	// Copy data
	if _, err := io.Copy(w, reader); err != nil {
		w.Close()
		return "", fmt.Errorf("failed to copy data: %w", err)
	}

	// Close writer
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Make object publicly readable
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set ACL: %w", err)
	}

	// Return public URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectPath)
	return url, nil
}

// Delete deletes an object from GCS
func (c *GCSClient) Delete(ctx context.Context, objectPath string) error {
	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(objectPath)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// GetSignedURL generates a signed URL for temporary access
// Note: Requires proper GCS credentials with signing capability
func (c *GCSClient) GetSignedURL(ctx context.Context, objectPath string, expiration time.Duration) (string, error) {
	bucket := c.client.Bucket(c.bucketName)

	// For now, return the public URL since we're making objects publicly readable
	// In production, you would use storage.SignedURL() function with proper credentials
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectPath)

	// TODO: Implement actual signed URLs when needed for private objects
	// This requires using storage.SignedURL() with a service account that has signing permissions

	_ = bucket // suppress unused warning
	_ = expiration

	return url, nil
}

// Close closes the GCS client
func (c *GCSClient) Close() error {
	return c.client.Close()
}
