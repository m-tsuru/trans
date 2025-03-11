# trans

Simple, Easy, Fine Copy.

## tl;dr

1. Download the latest executable for your machine architecture.

> Now, Windows is testing. Please build from source.

- [Releases](https://github.com/m-tsuru/trans/releases)

2. Make it Executable.

```
$ chmod +x trans-{YOUR_ARCHITECTURE}
```

3. Create `$HOME/.trans/config.yaml`

Here is a sample of following conditions.

- Copy from: /Volumes/MEMORY_CARDS/DCIM/*/*.(jpg|jpeg|JPG|JPEG...)
- Copy to: ~/Pictures/yyyy-mm-dd/*.(jpg|jpeg|JPG|JPEG...)
- Only and All JPEG files will be copy.

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

4. Run.

```sh
trans-{YOUR-ARCHITECTURE} import --profile default

# 4:05AM INF Read Configuration File Path: /Users/User/.trans/config.yaml
# 4:05AM INF Read Configuration Successfully: /Users/User/.trans/config.yaml
# 4:05AM INF [Profile: default] - Base Directory: /volumes/MEMORY_CARD/DCIM -> ~/Pictures
# ...
```

## Download & Run


## 免責事項 / Disclaimer

デベロッパ，およびコントリビュータは、データの損失、破損、またはその他の不具合について、一切の責任を負いません。データの管理およびバックアップはお客様の責任において行ってください。当社のサービスの使用または使用不能によって生じるいかなる直接的、間接的、特別、付随的、または結果的損害についても、当社は一切の責任を負わないものとします。

We shall not be held liable for any loss, corruption, or other issues related to data. It is the sole responsibility of the user to manage and back up their data. We disclaim any liability for direct, indirect, special, incidental, or consequential damages arising from the use or inability to use our services.
