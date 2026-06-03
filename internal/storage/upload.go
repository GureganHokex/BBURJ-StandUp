package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var allowedMimes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type Uploader struct {
	dir        string
	maxBytes   int64
	urlPrefix  string
}

func NewUploader(dir string, maxMB int) (*Uploader, error) {
	if maxMB <= 0 {
		maxMB = 10
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &Uploader{
		dir:       abs,
		maxBytes:  int64(maxMB) * 1024 * 1024,
		urlPrefix: "/uploads",
	}, nil
}

func (u *Uploader) Save(file *multipart.FileHeader) (string, error) {
	if file.Size > u.maxBytes {
		return "", fmt.Errorf("file too large (max %d MB)", u.maxBytes/(1024*1024))
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	sniff := make([]byte, 512)
	n, err := src.Read(sniff)
	if err != nil && err != io.EOF {
		return "", err
	}
	if n == 0 {
		return "", fmt.Errorf("empty file")
	}
	sniff = sniff[:n]

	mime := http.DetectContentType(sniff)
	ext, ok := allowedMimes[mime]
	if !ok {
		return "", fmt.Errorf("unsupported image type: %s", mime)
	}

	name := randomName() + ext
	destPath := filepath.Join(u.dir, name)

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := dst.Write(sniff); err != nil {
		return "", err
	}
	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(destPath)
		return "", err
	}

	return u.urlPrefix + "/" + name, nil
}

func randomName() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func IsLocalUpload(path string) bool {
	return strings.HasPrefix(path, "/uploads/")
}

func LocalPath(uploadDir, publicURL string) string {
	base := strings.TrimPrefix(publicURL, "/uploads/")
	return filepath.Join(uploadDir, filepath.Base(base))
}
