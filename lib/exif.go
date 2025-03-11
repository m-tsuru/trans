package lib

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/m-tsuru/trans/structs"
)

func GetExifData(path string) (*exif2.Exif, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	meta, err := imagemeta.Decode(f)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func GetImageMetadata(path string) (*structs.ImageMetadata, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	mime := http.DetectContentType(buffer)

	fileExif, err := GetExifData(path)
	if errors.Is(imagemeta.ErrNoExif, err) || errors.Is(imagemeta.ErrImageTypeNotFound, err) {
		return &structs.ImageMetadata{
			Path:         path,
			SaveDateTime: fileInfo.ModTime(),
			ExifDateTime: time.Unix(0, 0),
			Ext:          filepath.Ext(path),
			MIME:         mime,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	ExifDateTime := fileExif.DateTimeOriginal()
	if ExifDateTime.IsZero() {
		ExifDateTime = time.Unix(0, 0)
	}

	return &structs.ImageMetadata{
		Path:         path,
		SaveDateTime: fileInfo.ModTime(),
		ExifDateTime: ExifDateTime,
		Ext:          filepath.Ext(path),
		MIME:         mime,
	}, nil
}

func GetImageMetadataList(dir string) (structs.ImagesMetadataList, error) {
	var l structs.ImagesMetadataList
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		ignoreFileMatch, e := filepath.Match("*/._*", path)
		if e != nil {
			return e
		}
		if d.IsDir() == false && ignoreFileMatch == false {
			meta, e := GetImageMetadata(path)
			if err != nil {
				return e
			}
			if meta != nil {
				l = append(l, *meta)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return l, nil
}
