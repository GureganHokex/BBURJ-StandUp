package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IsLocalAsset(path string) bool {
	return IsLocalUpload(path) || strings.HasPrefix(path, "/static/")
}

// ImportFromURL downloads a remote image into uploads and returns its public path.
func (u *Uploader) ImportFromURL(ctx context.Context, rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" || IsLocalAsset(rawURL) {
		return rawURL, nil
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "BBURJ-StandUp/1.0 (+image-import)")
	req.Header.Set("Accept", "image/*")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("image returned status %d", resp.StatusCode)
	}

	limited := io.LimitReader(resp.Body, u.maxBytes+1)
	sniff := make([]byte, 512)
	n, err := limited.Read(sniff)
	if err != nil && err != io.EOF {
		return "", err
	}
	if n == 0 {
		return "", fmt.Errorf("empty image")
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

	written := int64(n)
	if _, err := dst.Write(sniff); err != nil {
		_ = os.Remove(destPath)
		return "", err
	}
	copied, err := io.Copy(dst, limited)
	if err != nil {
		_ = os.Remove(destPath)
		return "", err
	}
	written += copied
	if written > u.maxBytes {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("file too large (max %d MB)", u.maxBytes/(1024*1024))
	}

	return u.urlPrefix + "/" + name, nil
}
