package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/m-tsuru/trans/structs"
	"github.com/rs/zerolog/log"
)

func Import(dir string, importProfile structs.ImportProfile) error {
	// ディレクトリの存在確認
	stat, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return err
	} else if !stat.IsDir() {
		return fmt.Errorf("%s is not directory", dir)
	}

	// すべてのファイルを再帰的に探索
	list, err := GetImageMetadataList(dir)
	if err != nil {
		return err
	}

	// Statistics
	copied := 0
	copiedBytes := 0
	start := time.Now()

	for _, v := range list {
		chk, path, err := JudgeCopy(importProfile.Target, importProfile.Patterns, v)
		if err != nil {
			return err
		}
		if chk {
			stat, _ = os.Stat(v.Path)
			copied += 1
			copiedBytes += int(stat.Size())
			log.Info().Msgf("Copy: %s -> %s (Size: %s)", v.Path, *path, formatSize(stat.Size()))
			if Exists(*path) {
				switch importProfile.Samefile {
				case "filename":
					// コピー先に同じファイル名があったとき，スキップする
					log.Info().Msgf("Skip (Same Filename) - %s", v.Path)
					continue
				case "sha256":
					// コピー先に同じファイル名があったとき，ハッシュ値で判定する
					// 異なるハッシュ値であった場合，上書きする
					sameHash, _, _, err := CompareFileSum(v.Path, *path)
					if err != nil {
						return err
					}
					if sameHash {
						log.Info().Msgf("Skip (Same Hash) - %s", v.Path)
						continue
					}
				case "overwrite":
					// 上書きする
					fallthrough
				default:
					// 上書きする
				}
			}
			// コピーする
			err := copyFileWithProgress(v.Path, *path)
			if err != nil {
				return err
			}
		} else {
			log.Info().Msgf("Skip - %s", v.Path)
		}
	}

	processTime := time.Since(start)

	log.Info().Msgf(
		"[Stats] count: %d file(s); duration: %s; size: %s",
		copied,
		processTime.String(),
		formatSize(int64(copiedBytes)),
	)

	return nil
}

func JudgeCopy(targetBasePath string, patterns structs.Patterns, metadata structs.ImageMetadata) (bool, *string, error) {
	baseName := filepath.Base(metadata.Path)
	for _, v := range patterns {

		if len(v.Extensions) < 1 && len(v.Mime) < 1 {
			return false, nil, errors.New("valid import pattern does not exist")
		}

		var sortFormatted string
		if metadata.ExifDateTime.IsZero() == true || v.Datetime == "file" {
			sortFormatted = metadata.SaveDateTime.Format(v.Sort)
		} else {
			sortFormatted = metadata.ExifDateTime.Format(v.Sort)
		}
		targetPath := filepath.Join(targetBasePath, sortFormatted, baseName)
		if len(v.Extensions) > 0 && len(metadata.Ext) > 0 {
			if slices.Contains(v.Extensions, metadata.Ext[1:]) {
				return true, &targetPath, nil
			}
		} else if len(v.Mime) > 0 {
			if slices.Contains(v.Mime, metadata.MIME) {
				return true, &targetPath, nil
			}
		} else {
			return false, nil, nil
		}
	}
	return false, nil, nil
}
