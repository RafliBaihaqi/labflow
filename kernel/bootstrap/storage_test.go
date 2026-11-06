package bootstrap

import "testing"

func TestResolveStorageConfig(t *testing.T) {
	t.Run("S3_ENDPOINT wins", func(t *testing.T) {
		t.Setenv("S3_ENDPOINT", "https://s3.example.com")
		t.Setenv("MINIO_ENDPOINT", "minio:9000")
		got := resolveStorageConfig()
		if got.endpoint != "https://s3.example.com" || !got.pathStyle {
			t.Fatalf("unexpected config: %+v", got)
		}
	})
	t.Run("MINIO_ENDPOINT gets a scheme and path-style", func(t *testing.T) {
		t.Setenv("S3_ENDPOINT", "")
		t.Setenv("MINIO_ENDPOINT", "minio:9000")
		t.Setenv("MINIO_ACCESS", "u")
		t.Setenv("MINIO_SECRET", "p")
		got := resolveStorageConfig()
		if got.endpoint != "http://minio:9000" || !got.pathStyle || got.accessKey != "u" || got.secretKey != "p" {
			t.Fatalf("unexpected config: %+v", got)
		}
	})
	t.Run("falls back to R2", func(t *testing.T) {
		t.Setenv("S3_ENDPOINT", "")
		t.Setenv("MINIO_ENDPOINT", "")
		t.Setenv("ACCOUNT_ID", "acc")
		got := resolveStorageConfig()
		if got.endpoint != "https://acc.r2.cloudflarestorage.com" || got.pathStyle {
			t.Fatalf("unexpected config: %+v", got)
		}
	})
}
