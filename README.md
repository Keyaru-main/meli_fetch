# 🚀 MeliFetch - GitHub Actions Download Proxy

<div align="center">

![Version](https://img.shields.io/badge/version-2.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**قدرتمندترین ابزار دانلود از طریق GitHub Actions**

دانلود فایل‌ها و صفحات وب با استفاده از GitHub Actions به عنوان پروکسی

[ویژگی‌ها](#-ویژگی‌ها) • [نصب](#-نصب) • [استفاده](#-استفاده) • [مثال‌ها](#-مثال‌ها) • [مستندات](#-مستندات)

</div>

---

## 📋 فهرست مطالب

- [معرفی](#-معرفی)
- [ویژگی‌ها](#-ویژگی‌ها)
- [نصب](#-نصب)
- [پیکربندی](#-پیکربندی)
- [استفاده](#-استفاده)
- [مثال‌های کاربردی](#-مثال‌های-کاربردی)
- [ساختار پروژه](#-ساختار-پروژه)
- [توسعه](#-توسعه)
- [عیب‌یابی](#-عیب‌یابی)
- [مشارکت](#-مشارکت)

---

## 🎯 معرفی

**MeliFetch** یک ابزار قدرتمند برای دانلود فایل‌ها و محتوای وب از طریق GitHub Actions است. این ابزار به شما امکان می‌دهد:

- 🌐 از محدودیت‌های شبکه عبور کنید
- 🔒 دانلودهای امن و ناشناس داشته باشید
- ⚡ از سرعت بالای سرورهای GitHub استفاده کنید
- 🤖 صفحات پیچیده با JavaScript را دانلود کنید
- 📦 فایل‌های بزرگ را بدون محدودیت دریافت کنید

### چرا MeliFetch؟

- **دو حالت دانلود**: curl سریع برای فایل‌های ساده، Puppeteer برای صفحات پیچیده
- **مدیریت هوشمند**: صف دانلود، retry خودکار، cleanup
- **قابلیت سفارشی‌سازی بالا**: نام فایل، selector ها، timeout و...
- **رابط کاربری ساده**: CLI قدرتمند با خروجی زیبا
- **متن‌باز و رایگان**: کاملاً رایگان و قابل توسعه

---

## ✨ ویژگی‌ها

### 🎯 حالت‌های دانلود

#### 1️⃣ Fetch Mode (curl)
- ✅ سریع و کارآمد
- ✅ مناسب فایل‌های مستقیم
- ✅ پشتیبانی از redirect
- ✅ محدودیت حجم
- ✅ retry خودکار

#### 2️⃣ Download Mode (Puppeteer)
- ✅ مرورگر headless
- ✅ اجرای JavaScript
- ✅ کلیک خودکار
- ✅ انتظار برای selector
- ✅ اسکرین‌شات خودکار
- ✅ استخراج لینک‌های دانلود

### 🛠️ قابلیت‌های پیشرفته

- **Progress Bar**: نمایش پیشرفت دانلود
- **Metadata**: ذخیره اطلاعات کامل دانلود
- **Branch Management**: مدیریت شاخه‌های موقت
- **Custom Filename**: تعیین نام دلخواه
- **User-Agent**: تنظیم User-Agent سفارشی
- **Timeout Control**: کنترل زمان انتظار
- **Auto Cleanup**: پاکسازی خودکار

---

## 📦 نصب

### پیش‌نیازها

- Go 1.25 یا بالاتر
- Git
- حساب GitHub
- یک repository در GitHub

### نصب از سورس

```bash
# کلون کردن پروژه
git clone https://github.com/Keyaru-main/meli_fetch.git
cd melli_pet

# نصب dependencies
go mod download

# ساخت برنامه
go build -o melifetch cmd/melifetch/main.go

# نصب global (اختیاری)
sudo mv melifetch /usr/local/bin/
```

### نصب سریع

```bash
# دانلود و نصب مستقیم
curl -sSL https://raw.githubusercontent.com/Keyaru-main/melli_pet/main/install.sh | bash
```

---

## ⚙️ پیکربندی

### 1. ایجاد Repository

1. یک repository جدید در GitHub بسازید (مثلاً `melli-downloads`)
2. محتویات پوشه `repo/` را در آن قرار دهید
3. workflows را فعال کنید

### 2. ایجاد GitHub Token

1. به Settings → Developer settings → Personal access tokens بروید
2. Generate new token (classic) را بزنید
3. دسترسی‌های زیر را فعال کنید:
   - `repo` (تمام دسترسی‌ها)
   - `workflow`
4. توکن را کپی کنید

### 3. پیکربندی melifetch

```bash
# تنظیم توکن
melifetch config --token YOUR_GITHUB_TOKEN

# تنظیم repository
melifetch config --repo username/repo-name

# مشاهده تنظیمات
melifetch config
```

---

## 🚀 استفاده

### دستورات اصلی

#### 1. Fetch (دانلود سریع با curl)

```bash
# دانلود ساده
melifetch fetch https://example.com/file.zip

# با نام سفارشی
melifetch fetch https://example.com/file.zip --name myfile.zip

# با محدودیت حجم
melifetch fetch https://example.com/file.zip --max-size 500

# بدون پاکسازی شاخه
melifetch fetch https://example.com/file.zip --no-cleanup
```

#### 2. Download (دانلود پیشرفته با مرورگر)

```bash
# دانلود از صفحه پیچیده
melifetch download https://www.curseforge.com/minecraft/mc-mods/mod/download/123456

# با کلیک روی دکمه
melifetch download https://example.com/download \
  --click-selector "button.download"

# با انتظار برای element
melifetch download https://example.com/download \
  --wait-selector ".download-ready" \
  --wait-time 60

# با User-Agent سفارشی
melifetch download https://example.com/download \
  --user-agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
```

#### 3. List (لیست شاخه‌ها)

```bash
# نمایش تمام شاخه‌های موقت
melifetch list
```

#### 4. Clean (پاکسازی)

```bash
# پاکسازی یک شاخه
melifetch clean fetch-1234567890

# پاکسازی تمام شاخه‌ها
melifetch clean
```

### پارامترهای مشترک

| پارامتر | کوتاه | پیش‌فرض | توضیحات |
|---------|-------|---------|---------|
| `--output` | `-o` | `.` | مسیر خروجی |
| `--name` | `-n` | - | نام فایل خروجی |
| `--no-cleanup` | - | `false` | نگه‌داشتن شاخه موقت |
| `--token` | - | - | توکن GitHub |
| `--repo` | - | - | نام repository |
| `--timeout` | - | `5` یا `10` | زمان انتظار (دقیقه) |

---

## 💡 مثال‌های کاربردی

### مثال 1: دانلود فایل JSON از API

```bash
melifetch fetch https://api.github.com/repos/golang/go/releases/latest \
  --name go-latest.json \
  --output ./downloads
```

### مثال 2: دانلود تصویر

```bash
melifetch fetch https://picsum.photos/1920/1080 \
  --name wallpaper.jpg
```

### مثال 3: دانلود از CurseForge

```bash
melifetch download https://www.curseforge.com/minecraft/mc-mods/jei/download/5043868 \
  --wait-time 45 \
  --name jei-mod.jar
```

### مثال 4: دانلود با کلیک خودکار

```bash
melifetch download https://example.com/protected-download \
  --wait-selector "#download-ready" \
  --click-selector "button#start-download" \
  --wait-time 60
```

### مثال 5: دانلود فایل بزرگ

```bash
melifetch fetch https://releases.ubuntu.com/22.04/ubuntu-22.04-desktop-amd64.iso \
  --max-size 5000 \
  --timeout 15 \
  --name ubuntu.iso
```

### مثال 6: دانلود چندین فایل

```bash
#!/bin/bash
urls=(
  "https://example.com/file1.zip"
  "https://example.com/file2.zip"
  "https://example.com/file3.zip"
)

for url in "${urls[@]}"; do
  melifetch fetch "$url" --output ./batch-downloads
  sleep 5
done
```

---

## 📁 ساختار پروژه

```
melli_pet/
├── cmd/
│   └── melifetch/
│       └── main.go              # نقطه ورود CLI
├── internal/
│   ├── config/
│   │   └── config.go            # مدیریت تنظیمات
│   └── github/
│       └── fetcher.go           # منطق اصلی دانلود
├── repo/
│   ├── .github/
│   │   └── workflows/
│   │       ├── fetch.yml        # Workflow دانلود با curl
│   │       └── download.yml     # Workflow دانلود با Puppeteer
│   └── scripts/
│       ├── advanced_download.js # اسکریپت Puppeteer
│       ├── fetch_utils.sh       # توابع کمکی bash
│       └── test_download.sh     # تست‌های خودکار
├── go.mod
├── go.sum
└── README.md
```

### توضیح فایل‌های کلیدی

#### `cmd/melifetch/main.go`
- تعریف دستورات CLI
- مدیریت flags و arguments
- رابط کاربری

#### `internal/github/fetcher.go`
- ارتباط با GitHub API
- مدیریت workflows
- دانلود فایل‌ها
- مدیریت شاخه‌ها

#### `repo/.github/workflows/fetch.yml`
- دانلود با curl
- تشخیص نوع فایل
- ذخیره metadata
- commit و push

#### `repo/.github/workflows/download.yml`
- راه‌اندازی Puppeteer
- مدیریت دانلود پیچیده
- اسکرین‌شات
- استخراج لینک‌ها

---

## 🔧 توسعه

### ساخت از سورس

```bash
# کلون
git clone https://github.com/Keyaru-main/melli_pet.git
cd melli_pet

# نصب dependencies
go mod download

# اجرای تست‌ها
go test ./...

# ساخت
go build -o melifetch cmd/melifetch/main.go

# اجرا
./melifetch --help
```

### اضافه کردن قابلیت جدید

1. کد را در `internal/` اضافه کنید
2. دستور جدید را در `cmd/melifetch/main.go` تعریف کنید
3. workflow مربوطه را در `repo/.github/workflows/` بسازید
4. مستندات را به‌روز کنید

### تست workflow ها

```bash
# کپی repo به GitHub repository خود
cd repo
git init
git remote add origin https://github.com/Keyaru-main/YOUR_REPO.git
git add .
git commit -m "Initial commit"
git push -u origin main

# تست با melifetch
melifetch config --repo Keyaru-main/YOUR_REPO
melifetch fetch https://httpbin.org/json
```

---

## 🐛 عیب‌یابی

### مشکلات رایج

#### 1. خطای Authentication

```
❌ Failed to trigger workflow: 401 Unauthorized
```

**راه‌حل:**
- توکن GitHub را بررسی کنید
- دسترسی‌های توکن را چک کنید
- توکن را دوباره تنظیم کنید: `melifetch config --token NEW_TOKEN`

#### 2. Timeout

```
❌ Timeout after 5m0s
```

**راه‌حل:**
- زمان timeout را افزایش دهید: `--timeout 15`
- اتصال اینترنت را بررسی کنید
- وضعیت GitHub Actions را چک کنید

#### 3. شاخه ایجاد نشد

```
❌ Could not determine filename. Branch may not have completed successfully
```

**راه‌حل:**
- لاگ workflow را در GitHub بررسی کنید
- URL را چک کنید
- از `--no-cleanup` برای دیباگ استفاده کنید

#### 4. فایل دانلود نشد (Puppeteer)

```
⚠️ No automatic download detected
```

**راه‌حل:**
- selector های صحیح را مشخص کنید
- زمان انتظار را افزایش دهید
- فایل `found_links.json` را بررسی کنید

### لاگ‌های دیباگ

```bash
# نگه‌داشتن شاخه برای بررسی
melifetch fetch URL --no-cleanup

# مشاهده شاخه‌ها
melifetch list

# بررسی محتوای شاخه در GitHub
# https://github.com/USERNAME/REPO/tree/BRANCH_NAME
```

---

## 🤝 مشارکت

مشارکت شما استقبال می‌شود! 

### نحوه مشارکت

1. Fork کنید
2. شاخه جدید بسازید: `git checkout -b feature/amazing-feature`
3. تغییرات را commit کنید: `git commit -m 'Add amazing feature'`
4. Push کنید: `git push origin feature/amazing-feature`
5. Pull Request باز کنید

### راهنمای کد

- از Go conventions استفاده کنید
- کد را document کنید
- تست بنویسید
- README را به‌روز کنید

---

## 📄 لایسنس

این پروژه تحت لایسنس MIT منتشر شده است. برای جزئیات بیشتر فایل [LICENSE](LICENSE) را ببینید.

---

## 🙏 تشکر

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Puppeteer](https://pptr.dev/) - Headless browser
- [GitHub Actions](https://github.com/features/actions) - CI/CD platform

---

## 📞 ارتباط

- GitHub Issues: [مشکلات و پیشنهادات](https://github.com/Keyaru-main/melli_pet/issues)
- Discussions: [بحث و گفتگو](https://github.com/Keyaru-main/melli_pet/discussions)

---

<div align="center">

**ساخته شده با ❤️ برای جامعه متن‌باز**

⭐ اگر این پروژه برایتان مفید بود، یک ستاره بدهید!

</div>
