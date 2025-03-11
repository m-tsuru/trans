package lib

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/m-tsuru/trans/structs"
)

func Check(dir string, importProfile structs.ImportProfile) error {
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

	count := 0
	countNotCopied := 0
	countHashIncorrect := 0
	countSkiped := 0

	for _, v := range list {
		chk, path, err := JudgeCopy(importProfile.Target, importProfile.Patterns, v)
		if err != nil {
			return err
		}
		if chk {
			count++
			log.Info().Msgf("Check: %s -> %s", v.Path, *path)
			if !Exists(*path) {
				log.Error().Msgf("[Not be Copied]: %s", v.Path)
				countNotCopied++
			} else {
				ok, from, to, err := CompareFileSum(v.Path, *path)
				if err != nil {
					return err
				}
				log.Info().Msgf("# Original: %x", from)
				log.Info().Msgf("# Copy    : %x", to)
				if !ok {
					log.Error().Msgf("[Invalid Hash]: %s", v.Path)
					countHashIncorrect++
				}
			}
		} else {
			countSkiped++
			log.Info().Msgf("Skip : %s", v.Path)
		}
	}

	countError := countHashIncorrect + countNotCopied

	log.Info().Msgf("")
	log.Info().Msgf("====================")
	log.Info().Msgf("")
	if countError != 0 {
		log.Error().Msgf("Check Completed, Result: [ NG ]")
	} else {
		log.Info().Msgf("Check Completed, Result: [ OK ]")
	}
	log.Info().Msgf("")
	log.Info().Msgf("Statistics:")
	log.Info().Msgf("All Files: %d file(s), Target Files: %d file(s), Skip Files: %d file(s)", len(list), count, countSkiped)
	log.Info().Msgf("Invalid Hash: %d file(s), Not Copied: %d file(s), total: %d file(s)", countHashIncorrect, countNotCopied, countError)
	log.Info().Msgf("")
	log.Info().Msgf("====================")

	return nil
}
