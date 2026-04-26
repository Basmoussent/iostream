//go:build windows

package driver

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// zadigURL points at a pinned release so the binary we run is reproducible.
// Update both URL and SHA together when bumping the version.
const (
	zadigURL    = "https://github.com/pbatard/libwdi/releases/download/v1.5.1/zadig-2.9.exe"
	zadigSHA256 = "" // empty = skip checksum verification (TODO before release)
	zadigName   = "zadig-2.9.exe"
)

func setup() error {
	cached, err := ensureZadig()
	if err != nil {
		return fmt.Errorf("driver: %w", err)
	}

	fmt.Printf("driver: launching Zadig (%s)\n", cached)
	fmt.Println("driver: a UAC prompt will appear — click Yes, then in Zadig:")
	fmt.Println("        1. Options -> List All Devices")
	fmt.Println("        2. select the Apple iPhone interface (VID 05AC)")
	fmt.Println("        3. pick WinUSB in the right dropdown -> Replace Driver")
	fmt.Println()

	// PowerShell's Start-Process -Verb RunAs triggers UAC elevation; -Wait
	// keeps us blocked until Zadig closes so the caller can react to errors.
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf(`Start-Process -FilePath %q -Verb RunAs -Wait`, cached))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("driver: launch zadig: %w", err)
	}
	return nil
}

// ensureZadig downloads zadig-x.y.exe to %LOCALAPPDATA%\iostream\ on first
// run and returns its absolute path. Subsequent runs reuse the cached copy.
func ensureZadig() (string, error) {
	cacheDir, err := iostreamCacheDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}
	target := filepath.Join(cacheDir, zadigName)

	if ok, err := fileMatchesChecksum(target, zadigSHA256); err == nil && ok {
		return target, nil
	}

	fmt.Printf("driver: downloading %s\n", zadigURL)
	if err := downloadFile(zadigURL, target); err != nil {
		return "", fmt.Errorf("download zadig: %w", err)
	}
	if zadigSHA256 != "" {
		if ok, err := fileMatchesChecksum(target, zadigSHA256); !ok || err != nil {
			_ = os.Remove(target)
			return "", fmt.Errorf("zadig checksum mismatch — refusing to launch (err=%v)", err)
		}
	}
	return target, nil
}

func iostreamCacheDir() (string, error) {
	root := os.Getenv("LOCALAPPDATA")
	if root == "" {
		// Fallback for unusual setups (running as SYSTEM, etc.).
		dir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		root = dir
	}
	return filepath.Join(root, "iostream"), nil
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http %d", resp.StatusCode)
	}
	tmp := dest + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dest)
}

func fileMatchesChecksum(path, want string) (bool, error) {
	if want == "" {
		// No expected checksum → cache is "valid" only if file exists.
		_, err := os.Stat(path)
		return err == nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}
	return hex.EncodeToString(h.Sum(nil)) == want, nil
}
