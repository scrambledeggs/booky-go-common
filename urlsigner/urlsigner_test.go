package urlsigner

import (
	"context"
	"net/url"
	"testing"
)

func setTestCreds(t *testing.T) {
	t.Setenv("AMZ_ACCESS_KEY_ID", "test-access-key")
	t.Setenv("AMZ_SECRET_ACCESS_KEY", "test-secret-key")
	t.Setenv("AMZ_REGION", "us-east-1")
}

func TestGeneratePresignedURLAgainstAWS(t *testing.T) {
	setTestCreds(t)

	signedURL, err := GeneratePresignedURL(context.Background(), "my-bucket", "path/to/object.csv", 604800)
	if err != nil {
		t.Fatalf("GeneratePresignedURL returned error: %v", err)
	}

	parsed, err := url.Parse(signedURL)
	if err != nil {
		t.Fatalf("url.Parse returned error: %v", err)
	}
	if got, want := parsed.Host, "my-bucket.s3.us-east-1.amazonaws.com"; got != want {
		t.Errorf("Host = %q, want %q", got, want)
	}
	if got, want := parsed.Path, "/path/to/object.csv"; got != want {
		t.Errorf("Path = %q, want %q", got, want)
	}
	if got, want := parsed.Query().Get("X-Amz-Expires"), "604800"; got != want {
		t.Errorf("X-Amz-Expires = %q, want %q", got, want)
	}
}

func TestGeneratePresignedURLAgainstCustomEndpoint(t *testing.T) {
	setTestCreds(t)
	t.Setenv("AWS_ENDPOINT_URL_S3", "https://abc123.r2.cloudflarestorage.com")

	signedURL, err := GeneratePresignedURL(context.Background(), "my-bucket", "path/to/object.csv", 1500)
	if err != nil {
		t.Fatalf("GeneratePresignedURL returned error: %v", err)
	}

	parsed, err := url.Parse(signedURL)
	if err != nil {
		t.Fatalf("url.Parse returned error: %v", err)
	}
	// Path-style addressing: bucket stays in the path, host is the raw endpoint.
	if got, want := parsed.Host, "abc123.r2.cloudflarestorage.com"; got != want {
		t.Errorf("Host = %q, want %q", got, want)
	}
	if got, want := parsed.Path, "/my-bucket/path/to/object.csv"; got != want {
		t.Errorf("Path = %q, want %q", got, want)
	}
	if got, want := parsed.Query().Get("X-Amz-Expires"), "1500"; got != want {
		t.Errorf("X-Amz-Expires = %q, want %q", got, want)
	}
}
