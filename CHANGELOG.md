# 📝 Changelog

تمام تغییرات مهم این پروژه در این فایل مستند می‌شود.

---

## [2.0.0] - 2025-01-XX

### 🎉 نسخه جدید - بازسازی کامل

این نسخه یک بازسازی کامل پروژه است با قابلیت‌های پیشرفته.

### ✨ افزوده شده

#### CLI Commands
- **`melifetch fetch`** - دانلود سریع با curl
  - پشتیبانی از نام فایل سفارشی (`--name`)
  - محدودیت حجم قابل تنظیم (`--max-size`)
  - timeout قابل تنظیم (`--timeout`)
  
- **`melifetch download`** - دانلود پیشرفته با Puppeteer
  - انتظار برای CSS selector (`--wait-selector`)
  - کلیک خودکار (`--click-selector`)
  - User-Agent سفارشی (`--user-agent`)
  - زمان انتظار قابل تنظیم (`--wait-time`)
  
- **`melifetch list`** - لیست شاخه‌های موقت
  - نمایش شاخه‌های fetch و download
  - شمارش کل شاخه‌ها
  
- **`melifetch clean`** - پاکسازی شاخه‌ها
  - حذف یک شاخه خاص
  - حذف تمام شاخه‌های موقت
  
- **`melifetch config`** - مدیریت تنظیمات
  - نمایش تنظیمات فعلی
  - تنظیم token و repository

#### GitHub Actions Workflows

**fetch.yml (بهبود یافته)**
- تشخیص خودکار نوع فایل از Content-Type
- پشتیبانی از 25+ نوع فایل
- محاسبه سرعت دانلود
- metadata کامل با timestamp
- retry خودکار با backoff
- پشتیبانی از فایل‌های فشرده
- نمایش اطلاعات دقیق دانلود

**download.yml (کاملاً جدید)**
- مرورگر headless با Puppeteer
- اسکرین‌شات خودکار (قبل و بعد)
- استخراج لینک‌های دانلود
- ذخیره HTML صفحه
- metadata به فرمت JSON
- پشتیبانی از selector های پیچیده
- User-Agent قابل تنظیم
- timeout قابل تنظیم

#### Internal Packages

**internal/github/fetcher.go**
- `FetchOptions` و `DownloadOptions` structs
- `TriggerFetchWorkflow()` - trigger workflow fetch
- `TriggerDownloadWorkflow()` - trigger workflow download
- `WaitForBranch()` - انتظار با spinner انیمیشن
- `DownloadContent()` - دانلود با progress bar
- `copyWithProgress()` - کپی با نمایش پیشرفت
- `downloadOptionalFiles()` - دانلود فایل‌های اضافی
- `ListBranches()` - لیست شاخه‌ها با prefix
- `FetchWithOptions()` - fetch با options کامل
- `DownloadWithBrowser()` - download با مرورگر
- بررسی rate limit
- مدیریت خطا بهبود یافته

#### Scripts

**repo/scripts/advanced_download.js**
- اسکریپت Puppeteer کامل و حرفه‌ای
- لاگینگ رنگی و زیبا
- مدیریت دانلود هوشمند
- استخراج لینک‌های دانلود
- ذخیره metadata کامل
- پشتیبانی از تمام ورودی‌های workflow

**repo/scripts/fetch_utils.sh**
- توابع کمکی bash
- `fetch_with_curl()` - دانلود با curl
- `detect_extension()` - تشخیص extension
- `format_bytes()` - فرمت حجم
- `check_url()` - بررسی دسترسی
- `get_file_info()` - دریافت اطلاعات فایل
- `verify_file()` - تایید فایل
- لاگینگ رنگی

**repo/scripts/test_download.sh**
- مجموعه تست خودکار
- 5 تست مختلف
- گزارش نتایج

#### Documentation

- **README.md** - مستندات کامل فارسی
  - معرفی جامع
  - راهنمای نصب
  - مثال‌های کاربردی
  - عیب‌یابی
  
- **USAGE.md** - راهنمای استفاده
  - توضیح تمام دستورات
  - تمام flags و options
  - راهنمای CSS selectors
  - نکات و ترفندها
  
- **EXAMPLES.md** - مثال‌های واقعی
  - 30+ مثال کاربردی
  - اسکریپت‌های اتوماسیون
  - موارد خاص
  
- **CONTRIBUTING.md** - راهنمای مشارکت
  - نحوه مشارکت
  - استانداردهای کد
  - راهنمای تست
  
- **repo/README.md** - راهنمای repository
  - نحوه راه‌اندازی
  - توضیح workflows
  - سفارشی‌سازی

#### Build & Development

- **Makefile** - اتوماسیون build
  - `make build` - ساخت binary
  - `make build-all` - ساخت برای چند platform
  - `make install` - نصب
  - `make test` - اجرای تست‌ها
  - `make clean` - پاکسازی
  
- **install.sh** - اسکریپت نصب خودکار
  - بررسی پیش‌نیازها
  - نصب خودکار
  - پیکربندی
  
- **.gitignore** - فایل‌های ignore
- **LICENSE** - MIT License
- **CHANGELOG.md** - این فایل

### 🔧 تغییر یافته

- ساختار CLI کاملاً بازنویسی شد
- سیستم error handling بهبود یافت
- خروجی‌ها زیباتر و واضح‌تر شدند
- سرعت دانلود بهبود یافت
- مدیریت timeout بهتر شد

### 🐛 رفع شده

- مشکل timeout در فایل‌های بزرگ
- مشکل تشخیص نام فایل
- مشکل encoding در نام فایل‌های فارسی
- مشکل retry در خطاهای شبکه
- مشکل cleanup شاخه‌ها

### 🗑️ حذف شده

- کدهای قدیمی و استفاده نشده
- dependency های غیرضروری

---

## [1.0.0] - 2024-XX-XX

### نسخه اولیه

- دانلود ساده با curl
- workflow پایه
- CLI ساده

---

## 🔮 برنامه آینده

### [2.1.0] - آینده نزدیک

- [ ] پشتیبانی از proxy
- [ ] دانلود موازی
- [ ] resume قابلیت
- [ ] صف دانلود
- [ ] WebUI

### [3.0.0] - آینده دور

- [ ] پشتیبانی از torrent
- [ ] دانلود از YouTube
- [ ] پلاگین سیستم
- [ ] API server

---

<div align="center">

**برای تاریخچه کامل، [Releases](https://github.com/Keyaru-main/melli_pet/releases) را ببینید**

</div>
