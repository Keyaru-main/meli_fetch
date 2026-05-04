# 📖 راهنمای استفاده MeliFetch

راهنمای سریع و کامل برای استفاده از melifetch

---

## 🎯 شروع سریع

### نصب و پیکربندی (5 دقیقه)

```bash
# 1. ساخت برنامه
go build -o melifetch cmd/melifetch/main.go

# 2. تنظیم GitHub token
melifetch config --token ghp_xxxxxxxxxxxxxxxxxxxx

# 3. تنظیم repository
melifetch config --repo username/repo-name

# 4. اولین دانلود
melifetch fetch https://httpbin.org/json
```

---

## 📚 دستورات

### `melifetch fetch` - دانلود سریع

دانلود فایل‌ها با curl (سریع و ساده)

```bash
melifetch fetch [URL] [FLAGS]
```

**Flags:**

| Flag | کوتاه | نوع | پیش‌فرض | توضیحات |
|------|-------|-----|---------|---------|
| `--output` | `-o` | string | `.` | مسیر ذخیره فایل |
| `--name` | `-n` | string | - | نام فایل خروجی |
| `--type` | `-t` | string | `web` | نوع محتوا (web/file) |
| `--max-size` | - | int | `100` | حداکثر حجم (MB) |
| `--timeout` | - | int | `5` | زمان انتظار (دقیقه) |
| `--no-cleanup` | - | bool | `false` | نگه‌داشتن شاخه |
| `--token` | - | string | - | توکن GitHub |
| `--repo` | - | string | - | نام repository |

**مثال‌ها:**

```bash
# دانلود ساده
melifetch fetch https://example.com/file.zip

# با نام سفارشی
melifetch fetch https://example.com/file.zip -n myfile.zip

# با مسیر خروجی
melifetch fetch https://example.com/file.zip -o ./downloads

# محدودیت حجم 500MB
melifetch fetch https://example.com/large-file.iso --max-size 500

# timeout 15 دقیقه
melifetch fetch https://example.com/slow-server.zip --timeout 15
```

---

### `melifetch download` - دانلود پیشرفته

دانلود با مرورگر headless (برای صفحات پیچیده)

```bash
melifetch download [URL] [FLAGS]
```

**Flags:**

| Flag | کوتاه | نوع | پیش‌فرض | توضیحات |
|------|-------|-----|---------|---------|
| `--output` | `-o` | string | `.` | مسیر ذخیره فایل |
| `--name` | `-n` | string | - | نام فایل خروجی |
| `--wait-selector` | - | string | - | CSS selector برای انتظار |
| `--click-selector` | - | string | - | CSS selector برای کلیک |
| `--wait-time` | - | int | `30` | زمان انتظار (ثانیه) |
| `--user-agent` | - | string | - | User-Agent سفارشی |
| `--timeout` | - | int | `10` | زمان انتظار (دقیقه) |
| `--no-cleanup` | - | bool | `false` | نگه‌داشتن شاخه |
| `--token` | - | string | - | توکن GitHub |
| `--repo` | - | string | - | نام repository |

**مثال‌ها:**

```bash
# دانلود ساده
melifetch download https://example.com/download-page

# با کلیک روی دکمه
melifetch download https://example.com/page \
  --click-selector "button.download"

# انتظار برای element
melifetch download https://example.com/page \
  --wait-selector ".ready" \
  --wait-time 60

# User-Agent سفارشی
melifetch download https://example.com/page \
  --user-agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

# ترکیب همه
melifetch download https://example.com/complex \
  --wait-selector "#download-ready" \
  --click-selector "a.download-link" \
  --wait-time 90 \
  --name result.jar \
  --output ./downloads
```

---

### `melifetch list` - لیست شاخه‌ها

نمایش شاخه‌های موقت

```bash
melifetch list [FLAGS]
```

**خروجی:**

```
📥 Fetch branches (2):
   • fetch-1704123456
   • fetch-1704123789

🌐 Download branches (1):
   • download-1704124000

Total: 3 branches
```

---

### `melifetch clean` - پاکسازی

حذف شاخه‌های موقت

```bash
# حذف یک شاخه
melifetch clean [BRANCH_NAME]

# حذف همه
melifetch clean
```

**مثال‌ها:**

```bash
# حذف یک شاخه خاص
melifetch clean fetch-1704123456

# حذف تمام شاخه‌ها
melifetch clean
```

---

### `melifetch config` - پیکربندی

مدیریت تنظیمات

```bash
melifetch config [FLAGS]
```

**Flags:**

| Flag | توضیحات |
|------|---------|
| `--token` | تنظیم GitHub token |
| `--repo` | تنظیم repository |

**مثال‌ها:**

```bash
# مشاهده تنظیمات
melifetch config

# تنظیم token
melifetch config --token ghp_xxxxxxxxxxxxxxxxxxxx

# تنظیم repo
melifetch config --repo username/repo-name

# تنظیم همزمان
melifetch config --token TOKEN --repo username/repo
```

---

## 🎨 CSS Selectors

برای استفاده از `--wait-selector` و `--click-selector` باید CSS selector صحیح را بدانید.

### یافتن Selector

1. صفحه را در مرورگر باز کنید
2. F12 را بزنید (Developer Tools)
3. روی element مورد نظر راست‌کلیک کنید
4. Copy → Copy selector

### مثال‌های Selector

```css
/* ID */
#download-button

/* Class */
.btn-download

/* Tag */
button

/* Attribute */
[data-download="true"]

/* ترکیبی */
div.container > button.download
a[href*="download"]
```

### مثال‌های واقعی

```bash
# CurseForge
melifetch download "https://www.curseforge.com/..." \
  --wait-selector ".download-button"

# GitHub Releases
melifetch download "https://github.com/.../releases/..." \
  --click-selector "a[href*='.zip']"

# SourceForge
melifetch download "https://sourceforge.net/..." \
  --wait-selector "#download" \
  --click-selector "a.direct-download"
```

---

## 🔍 دیباگ و عیب‌یابی

### نگه‌داشتن شاخه برای بررسی

```bash
melifetch fetch URL --no-cleanup
```

سپس در GitHub به آدرس زیر بروید:
```
https://github.com/USERNAME/REPO/tree/BRANCH_NAME
```

### مشاهده لاگ‌های GitHub Actions

1. به repository خود بروید
2. تب Actions را باز کنید
3. آخرین workflow run را انتخاب کنید
4. لاگ‌ها را بررسی کنید

### بررسی فایل‌های دانلود شده

در حالت download، فایل‌های زیر ذخیره می‌شوند:

- `metadata.json` - اطلاعات دانلود
- `page_initial.png` - اسکرین‌شات اولیه
- `page_final.png` - اسکرین‌شات نهایی
- `found_links.json` - لینک‌های یافت شده
- `page.html` - HTML صفحه

---

## 💡 نکات و ترفندها

### 1. دانلود سریع‌تر

```bash
# استفاده از fetch به جای download
melifetch fetch URL  # سریع‌تر
melifetch download URL  # کندتر ولی قدرتمندتر
```

### 2. دانلود چندین فایل

```bash
#!/bin/bash
for url in $(cat urls.txt); do
  melifetch fetch "$url" -o ./batch
  sleep 5  # فاصله بین درخواست‌ها
done
```

### 3. دانلود با retry

```bash
#!/bin/bash
max_retries=3
for i in $(seq 1 $max_retries); do
  if melifetch fetch "$URL"; then
    break
  fi
  echo "Retry $i/$max_retries"
  sleep 10
done
```

### 4. پاکسازی خودکار

```bash
# اضافه کردن به crontab
0 0 * * * melifetch clean  # هر شب نیمه‌شب
```

### 5. استفاده در اسکریپت

```bash
#!/bin/bash
set -e

# دانلود
melifetch fetch "$URL" -n output.zip -o /tmp

# استخراج
unzip /tmp/output.zip -d ./extracted

# پردازش
process_files ./extracted
```

---

## 📊 محدودیت‌ها

### GitHub Actions

- **زمان اجرا**: حداکثر 6 ساعت
- **حجم فایل**: محدودیت سخت‌افزاری (معمولاً تا 2GB)
- **تعداد workflow**: 1000 درخواست در ساعت
- **فضای ذخیره‌سازی**: بستگی به پلن GitHub

### توصیه‌ها

- برای فایل‌های بزرگ از `--max-size` استفاده کنید
- بین درخواست‌ها فاصله بگذارید
- شاخه‌های قدیمی را پاک کنید
- از حالت fetch برای فایل‌های ساده استفاده کنید

---

## 🆘 کمک بیشتر

- [README.md](README.md) - مستندات کامل
- [EXAMPLES.md](EXAMPLES.md) - مثال‌های بیشتر
- [GitHub Issues](https://github.com/Keyaru-main/melli_pet/issues) - گزارش مشکل
- [Discussions](https://github.com/Keyaru-main/melli_pet/discussions) - پرسش و پاسخ

---

<div align="center">

**سوال دارید? Issue باز کنید!**

[📝 گزارش مشکل](https://github.com/Keyaru-main/melli_pet/issues/new) • [💬 پرسش](https://github.com/Keyaru-main/melli_pet/discussions/new)

</div>
