# trans

Simple, Easy, Fine Copy.

[Zenn の記事](https://zenn.dev/sasakulari/articles/763d889ef7c201)

## tl;dr

1. 最新のリリースをダウンロードする

> Windows 版は開発中です．試したい場合は，ビルドしてください

- [Releases](https://github.com/m-tsuru/trans/releases)

2. 実行ファイルに，実行権限を与える

```
$ chmod +x trans-{YOUR_ARCHITECTURE}
```

3. `$HOME/.trans/config.yaml` を作成する

以下の例は，次のような条件で転送を行う場合の `config.yaml` です

- 転送元ディレクトリ: /Volumes/MEMORY_CARDS/DCIM/*/*.(jpg|jpeg|JPG|JPEG...)
- 転送先ディレクトリ: ~/Pictures/yyyy-mm-dd/*.(jpg|jpeg|JPG|JPEG...)
- 転送するファイル: MIME タイプが `image/jpeg` のもの

```yaml
import:
  - name: "default"
    default: true
    original: /Volumes/MEMORY_CARD/DCIM
    # original: E:/DCIM (for Windows)
    target: ~/Pictures
    # target: F:/Images (for Windows)
    patterns:
      - name: "jpeg"
        mime:
          - "image/jpeg"
        sort: "2006-01-02/" # the time format used by Go is available.
        datetime: exif
```

4. 実行する

```sh
trans-{YOUR-ARCHITECTURE} import --profile default

# 4:05AM INF Read Configuration File Path: /Users/User/.trans/config.yaml
# 4:05AM INF Read Configuration Successfully: /Users/User/.trans/config.yaml
# 4:05AM INF [Profile: default] - Base Directory: /volumes/MEMORY_CARD/DCIM -> ~/Pictures
# ...
```

## 免責事項

デベロッパ，およびコントリビュータは、データの損失、破損、またはその他の不具合について、一切の責任を負いません。データの管理およびバックアップはユーザの責任において行ってください。サービスの使用または使用不能によって生じるいかなる直接的、間接的、特別、付随的、または結果的損害についても、当社は一切の責任を負わないものとします。
