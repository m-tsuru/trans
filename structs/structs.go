package structs

import "time"

type Conf struct {
	//
	// Conf は config.yaml をパースするときに読み出される構造体です．
	//
	Import []ImportProfile `yaml:"import"`
}

type Pattern struct { // インポートパターン
	Name       string   `yaml:"name"`                 // パターン名
	Extensions []string `yaml:"extensions,omitempty"` // 一致させる拡張子
	Mime       []string `yaml:"mime,omitempty"`       // 一致させる MIME Type (拡張子が指定されている場合は無視)
	Sort       string   `yaml:"sort"`                 // コピー先への振り分け方（Target/Sort），Go 日付フォーマットで指定
	Datetime   string   `yaml:"datetime,omitempty"`   // 日付取得元 (Exif (exif), ファイル日時 (file))
}

type Patterns []Pattern

type ImportProfile struct { // インポートプロファイル
	Name     string `yaml:"name"`
	Default  bool   `yaml:"default"`  // プロファイル名
	Original string `yaml:"original"` // コピー元ディレクトリ
	Target   string `yaml:"target"`   // コピー先ディレクトリ

	Samefile string `yaml:"samefile,omitempty"`
	// 同一ファイルを見つけた場合に，どのように対応するか
	// コピー先のファイルパスの先にファイルが存在した場合の動作を指定する
	//
	// コピー先ファイル名一致スキップ (filename)
	// SHA-256 ハッシュ値スキップ (sha256)
	// 常に上書き (overwrite), 指定無しの場合はこれが選択される

	Patterns Patterns `yaml:"patterns"`
}

type ImagesMetadataList []ImageMetadata

type ImageMetadata struct {
	Path         string
	SaveDateTime time.Time
	ExifDateTime time.Time
	Ext          string
	MIME         string
}
