package lib

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Copy(from string, to string) error {
	toFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFile.Close()

	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	io.Copy(toFile, fromFile)
	return nil
}

func CompareFileSum(from string, to string) (bool, *[]byte, *[]byte, error) {
	fromData, err := os.Open(from)
	if err != nil {
		return false, nil, nil, fmt.Errorf("failed to calculate check sum (%s): %s", from, err)
	}
	fH := sha256.New()
	if _, err := io.Copy(fH, fromData); err != nil {
		return false, nil, nil, fmt.Errorf("failed to calculate check sum (%s): %s", from, err)
	}

	toData, err := os.Open(to)
	if err != nil {
		return false, nil, nil, fmt.Errorf("failed to calculate check sum (%s): %s", to, err)
	}
	tH := sha256.New()
	if _, err := io.Copy(tH, toData); err != nil {
		return false, nil, nil, fmt.Errorf("failed to calculate check sum (%s): %s", to, err)
	}

	fromHash := fH.Sum(nil)
	toHash := tH.Sum(nil)

	if bytes.Equal(fromHash, toHash) {
		return true, &fromHash, &toHash, nil
	} else {
		return false, &fromHash, &toHash, nil
	}
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(ex)
}

func formatSize(size int64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	var result string
	sizeFloat := float64(size)

	switch {
	case sizeFloat >= PB:
		result = fmt.Sprintf("%.2f PB", sizeFloat/PB)
	case sizeFloat >= TB:
		result = fmt.Sprintf("%.2f TB", sizeFloat/TB)
	case sizeFloat >= GB:
		result = fmt.Sprintf("%.2f GB", sizeFloat/GB)
	case sizeFloat >= MB:
		result = fmt.Sprintf("%.2f MB", sizeFloat/MB)
	case sizeFloat >= KB:
		result = fmt.Sprintf("%.2f KB", sizeFloat/KB)
	default:
		result = fmt.Sprintf("%d B", size)
	}

	return result
}

func copyFileWithProgress(from string, to string) error {
	// コピー元ファイルを開く
	fromFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("failed to open original file: %w", err)
	}
	defer fromFile.Close()

	// コピー元ファイルの情報を取得
	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get original file stats: %w", err)
	}

	// ディレクトリが存在しない場合は作成
	toDir := filepath.Dir(to)
	err = os.MkdirAll(toDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	// プログレスバーを作成
	bar := progressbar.NewOptions64(
		fromFileInfo.Size(),
		progressbar.OptionSetDescription("Copy"),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: "-", BarStart: "[", BarEnd: "]"}),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(false),
	)

	// コピー先ファイルを作成
	toFile, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("failed to create new file: %w", err)
	}
	defer toFile.Close()

	// コピーを実行しながらプログレスバーを更新
	_, err = io.Copy(io.MultiWriter(toFile, bar), fromFile)
	if err != nil {
		return fmt.Errorf("failed to copy data: %s", err)
	}

	return nil
}
