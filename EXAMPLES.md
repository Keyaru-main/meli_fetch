# 💡 مثال‌های کاربردی MeliFetch

مجموعه‌ای از مثال‌های واقعی و کاربردی

---

## 📑 فهرست

- [دانلود فایل‌های ساده](#-دانلود-فایل‌های-ساده)
- [دانلود از سایت‌های پیچیده](#-دانلود-از-سایت‌های-پیچیده)
- [دانلود دسته‌جمعی](#-دانلود-دسته‌جمعی)
- [اتوماسیون](#-اتوماسیون)
- [موارد خاص](#-موارد-خاص)

---

## 📥 دانلود فایل‌های ساده

### 1. دانلود JSON از API

```bash
# GitHub API
melifetch fetch https://api.github.com/repos/golang/go/releases/latest \
  --name go-latest.json

# JSONPlaceholder
melifetch fetch https://jsonplaceholder.typicode.com/posts \
  --name posts.json
```

### 2. دانلود تصاویر

```bash
# تصویر تصادفی
melifetch fetch https://picsum.photos/1920/1080 \
  --name wallpaper.jpg

# لوگو
melifetch fetch https://github.com/golang/go/raw/master/doc/gopher/frontpage.png \
  --name gopher.png
```

### 3. دانلود فایل‌های متنی

```bash
# README
melifetch fetch https://raw.githubusercontent.com/golang/go/master/README.md \
  --name go-readme.md

# License
melifetch fetch https://raw.githubusercontent.com/torvalds/linux/master/COPYING \
  --name linux-license.txt
```

### 4. دانلود آرشیو

```bash
# ZIP از GitHub
melifetch fetch https://github.com/golang/go/archive/refs/heads/master.zip \
  --name go-source.zip \
  --max-size 200

# TAR.GZ
melifetch fetch https://golang.org/dl/go1.21.0.linux-amd64.tar.gz \
  --name go.tar.gz \
  --max-size 150
```

---

## 🌐 دانلود از سایت‌های پیچیده

### 1. CurseForge (Minecraft Mods)

```bash
# دانلود مود
melifetch download \
  "https://www.curseforge.com/minecraft/mc-mods/jei/download/5043868" \
  --wait-time 45 \
  --name jei-mod.jar \
  --output ./minecraft/mods

# با کلیک خودکار
melifetch download \
  "https://www.curseforge.com/minecraft/mc-mods/optifine/download/4907048" \
  --click-selector "a.download-button" \
  --wait-time 60 \
  --name optifine.jar
```

### 2. SourceForge

```bash
# دانلود پروژه
melifetch download \
  "https://sourceforge.net/projects/sevenzip/files/latest/download" \
  --wait-selector "#download" \
  --wait-time 30 \
  --name 7zip.exe
```

### 3. MediaFire

```bash
# دانلود فایل
melifetch download \
  "https://www.mediafire.com/file/xxxxxxxxx/file.zip/file" \
  --wait-selector "#downloadButton" \
  --click-selector "#downloadButton" \
  --wait-time 60
```

### 4. Google Drive (لینک عمومی)

```bash
# دانلود فایل کوچک
melifetch download \
  "https://drive.google.com/uc?export=download&id=FILE_ID" \
  --wait-time 30 \
  --name gdrive-file.pdf
```

### 5. Dropbox

```bash
# دانلود فایل اشتراکی
melifetch download \
  "https://www.dropbox.com/s/xxxxxxxxx/file.zip?dl=1" \
  --wait-time 30
```

---

## 📦 دانلود دسته‌جمعی

### 1. دانلود از لیست URL

```bash
#!/bin/bash
# download_list.sh

# فایل urls.txt:
# https://example.com/file1.zip
# https://example.com/file2.zip
# https://example.com/file3.zip

while IFS= read -r url; do
  echo "Downloading: $url"
  melifetch fetch "$url" --output ./downloads
  sleep 5  # فاصله بین درخواست‌ها
done < urls.txt

echo "All downloads completed!"
```

### 2. دانلود با نام‌گذاری خودکار

```bash
#!/bin/bash
# download_numbered.sh

urls=(
  "https://example.com/file1.zip"
  "https://example.com/file2.zip"
  "https://example.com/file3.zip"
)

counter=1
for url in "${urls[@]}"; do
  filename=$(printf "file_%03d.zip" $counter)
  melifetch fetch "$url" --name "$filename" --output ./batch
  ((counter++))
  sleep 3
done
```

### 3. دانلود موازی (محدود)

```bash
#!/bin/bash
# parallel_download.sh

max_parallel=3
urls_file="urls.txt"

download_file() {
  local url=$1
  local name=$(basename "$url")
  melifetch fetch "$url" --name "$name" --output ./parallel
}

export -f download_file

cat "$urls_file" | xargs -P $max_parallel -I {} bash -c 'download_file "{}"'
```

### 4. دانلود با retry

```bash
#!/bin/bash
# download_with_retry.sh

download_with_retry() {
  local url=$1
  local max_retries=3
  local retry_delay=10
  
  for i in $(seq 1 $max_retries); do
    echo "Attempt $i/$max_retries: $url"
    
    if melifetch fetch "$url" --output ./downloads; then
      echo "✓ Success"
      return 0
    fi
    
    if [ $i -lt $max_retries ]; then
      echo "⚠️  Failed, retrying in ${retry_delay}s..."
      sleep $retry_delay
    fi
  done
  
  echo "❌ Failed after $max_retries attempts"
  return 1
}

# استفاده
download_with_retry "https://example.com/file.zip"
```

---

## 🤖 اتوماسیون

### 1. دانلود روزانه

```bash
#!/bin/bash
# daily_download.sh

# اضافه به crontab:
# 0 2 * * * /path/to/daily_download.sh

DATE=$(date +%Y%m%d)
OUTPUT_DIR="./downloads/$DATE"

mkdir -p "$OUTPUT_DIR"

# دانلود فایل‌های روزانه
melifetch fetch "https://example.com/daily/report.pdf" \
  --name "report_$DATE.pdf" \
  --output "$OUTPUT_DIR"

# پاکسازی فایل‌های قدیمی (بیش از 7 روز)
find ./downloads -type d -mtime +7 -exec rm -rf {} +

# پاکسازی شاخه‌های موقت
melifetch clean
```

### 2. مانیتورینگ و دانلود خودکار

```bash
#!/bin/bash
# monitor_and_download.sh

WATCH_URL="https://api.github.com/repos/owner/repo/releases/latest"
LAST_VERSION_FILE=".last_version"

get_latest_version() {
  melifetch fetch "$WATCH_URL" --name temp.json -o /tmp
  jq -r '.tag_name' /tmp/temp.json
}

download_release() {
  local version=$1
  local download_url=$(jq -r '.assets[0].browser_download_url' /tmp/temp.json)
  
  echo "New version detected: $version"
  melifetch fetch "$download_url" \
    --name "release_$version.zip" \
    --output ./releases
}

# بررسی نسخه جدید
current_version=$(get_latest_version)
last_version=$(cat "$LAST_VERSION_FILE" 2>/dev/null || echo "")

if [ "$current_version" != "$last_version" ]; then
  download_release "$current_version"
  echo "$current_version" > "$LAST_VERSION_FILE"
else
  echo "No new version"
fi
```

### 3. دانلود با اعلان

```bash
#!/bin/bash
# download_with_notification.sh

send_notification() {
  local title=$1
  local message=$2
  
  # Linux (notify-send)
  notify-send "$title" "$message"
  
  # macOS (osascript)
  # osascript -e "display notification \"$message\" with title \"$title\""
}

URL="https://example.com/large-file.iso"
FILENAME="ubuntu.iso"

send_notification "Download Started" "Downloading $FILENAME..."

if melifetch fetch "$URL" --name "$FILENAME" --timeout 30; then
  send_notification "Download Complete" "$FILENAME downloaded successfully!"
else
  send_notification "Download Failed" "Failed to download $FILENAME"
fi
```

### 4. آرشیو خودکار

```bash
#!/bin/bash
# auto_archive.sh

DOWNLOAD_DIR="./downloads"
ARCHIVE_DIR="./archive"
DATE=$(date +%Y%m)

# دانلود فایل‌ها
melifetch fetch "https://example.com/file1.zip" -o "$DOWNLOAD_DIR"
melifetch fetch "https://example.com/file2.zip" -o "$DOWNLOAD_DIR"

# آرشیو
mkdir -p "$ARCHIVE_DIR"
tar -czf "$ARCHIVE_DIR/downloads_$DATE.tar.gz" -C "$DOWNLOAD_DIR" .

# پاکسازی
rm -rf "$DOWNLOAD_DIR"/*

echo "Archived to: $ARCHIVE_DIR/downloads_$DATE.tar.gz"
```

---

## 🎯 موارد خاص

### 1. دانلود با User-Agent سفارشی

```bash
# Windows Chrome
melifetch download "https://example.com/file" \
  --user-agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

# Mobile Safari
melifetch download "https://example.com/file" \
  --user-agent "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1"

# Bot
melifetch download "https://example.com/file" \
  --user-agent "MeliFetch-Bot/2.0"
```

### 2. دانلود فایل‌های بزرگ

```bash
# ISO لینوکس (4GB)
melifetch fetch \
  "https://releases.ubuntu.com/22.04/ubuntu-22.04-desktop-amd64.iso" \
  --max-size 5000 \
  --timeout 30 \
  --name ubuntu.iso

# بازی (10GB+)
melifetch fetch \
  "https://example.com/game-installer.exe" \
  --max-size 15000 \
  --timeout 60 \
  --name game.exe
```

### 3. دانلود با selector پیچیده

```bash
# انتظار برای چند element
melifetch download "https://example.com/page" \
  --wait-selector "div.content.loaded[data-ready='true']" \
  --click-selector "button.download:not(.disabled)" \
  --wait-time 90

# استفاده از attribute selector
melifetch download "https://example.com/page" \
  --click-selector "a[href*='download'][data-file-id='12345']"
```

### 4. دانلود و پردازش

```bash
#!/bin/bash
# download_and_process.sh

# دانلود
melifetch fetch "https://example.com/data.json" \
  --name data.json \
  --output /tmp

# پردازش با jq
cat /tmp/data.json | jq '.items[] | select(.active == true)' > filtered.json

# آپلود به جای دیگر
curl -X POST -H "Content-Type: application/json" \
  -d @filtered.json \
  https://api.example.com/upload

# پاکسازی
rm /tmp/data.json filtered.json
```

### 5. دانلود با لاگ کامل

```bash
#!/bin/bash
# download_with_logging.sh

LOG_FILE="download_$(date +%Y%m%d_%H%M%S).log"

{
  echo "=== Download Started at $(date) ==="
  echo "URL: $1"
  echo ""
  
  melifetch fetch "$1" --output ./downloads 2>&1
  
  echo ""
  echo "=== Download Finished at $(date) ==="
} | tee "$LOG_FILE"

# بررسی نتیجه
if [ ${PIPESTATUS[0]} -eq 0 ]; then
  echo "✓ Success" | tee -a "$LOG_FILE"
else
  echo "✗ Failed" | tee -a "$LOG_FILE"
fi
```

### 6. دانلود با محدودیت زمانی

```bash
#!/bin/bash
# timed_download.sh

MAX_TIME=300  # 5 minutes

timeout $MAX_TIME melifetch fetch "$URL" --output ./downloads

case $? in
  0)
    echo "✓ Download completed"
    ;;
  124)
    echo "✗ Download timed out after ${MAX_TIME}s"
    ;;
  *)
    echo "✗ Download failed"
    ;;
esac
```

---

## 🔧 اسکریپت‌های کاربردی

### مدیر دانلود ساده

```bash
#!/bin/bash
# download_manager.sh

QUEUE_FILE="download_queue.txt"
COMPLETED_FILE="completed.txt"
FAILED_FILE="failed.txt"

add_to_queue() {
  echo "$1" >> "$QUEUE_FILE"
  echo "Added to queue: $1"
}

process_queue() {
  if [ ! -f "$QUEUE_FILE" ]; then
    echo "Queue is empty"
    return
  fi
  
  while IFS= read -r url; do
    echo "Processing: $url"
    
    if melifetch fetch "$url" --output ./downloads; then
      echo "$url" >> "$COMPLETED_FILE"
      echo "✓ Completed"
    else
      echo "$url" >> "$FAILED_FILE"
      echo "✗ Failed"
    fi
    
    sleep 5
  done < "$QUEUE_FILE"
  
  rm "$QUEUE_FILE"
  echo "Queue processed!"
}

case "$1" in
  add)
    add_to_queue "$2"
    ;;
  process)
    process_queue
    ;;
  *)
    echo "Usage: $0 {add|process} [url]"
    ;;
esac
```

### استفاده:

```bash
# اضافه کردن به صف
./download_manager.sh add "https://example.com/file1.zip"
./download_manager.sh add "https://example.com/file2.zip"

# پردازش صف
./download_manager.sh process
```

---

## 📚 منابع بیشتر

- [README.md](README.md) - مستندات اصلی
- [USAGE.md](USAGE.md) - راهنمای استفاده
- [GitHub Issues](https://github.com/Keyaru-main/melli_pet/issues) - مشکلات و پیشنهادات

---

<div align="center">

**مثال جدیدی دارید؟ Pull Request بفرستید!**

[🔧 مشارکت](CONTRIBUTING.md) • [💬 بحث](https://github.com/Keyaru-main/melli_pet/discussions)

</div>
