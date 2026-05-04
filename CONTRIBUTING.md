# 🤝 راهنمای مشارکت در MeliFetch

از اینکه می‌خواهید در MeliFetch مشارکت کنید متشکریم! 

---

## 📋 فهرست

- [نحوه مشارکت](#-نحوه-مشارکت)
- [گزارش باگ](#-گزارش-باگ)
- [پیشنهاد ویژگی](#-پیشنهاد-ویژگی)
- [توسعه کد](#-توسعه-کد)
- [استانداردهای کد](#-استانداردهای-کد)
- [تست](#-تست)
- [Pull Request](#-pull-request)

---

## 🚀 نحوه مشارکت

### 1. Fork کردن پروژه

```bash
# Fork در GitHub
# سپس clone کنید
git clone https://github.com/Keyaru-main/meli_fetch.git
cd melli_pet
```

### 2. ایجاد شاخه جدید

```bash
git checkout -b feature/amazing-feature
# یا
git checkout -b fix/bug-description
```

### 3. انجام تغییرات

کد خود را بنویسید و commit کنید:

```bash
git add .
git commit -m "Add amazing feature"
```

### 4. Push کردن

```bash
git push origin feature/amazing-feature
```

### 5. ایجاد Pull Request

در GitHub یک Pull Request باز کنید.

---

## 🐛 گزارش باگ

برای گزارش باگ، یک Issue جدید با برچسب `bug` باز کنید.

### اطلاعات مورد نیاز:

- **توضیح مشکل**: توضیح واضح و مختصر
- **مراحل بازتولید**: چگونه می‌توان مشکل را بازتولید کرد
- **رفتار مورد انتظار**: چه چیزی باید اتفاق بیفتد
- **رفتار واقعی**: چه چیزی اتفاق افتاد
- **محیط**: OS، نسخه Go، نسخه melifetch
- **لاگ‌ها**: خروجی خطا یا لاگ‌های مربوطه

### مثال:

```markdown
**توضیح مشکل**
دانلود فایل‌های بزرگ با خطای timeout مواجه می‌شود.

**مراحل بازتولید**
1. اجرای `melifetch fetch https://example.com/large-file.iso`
2. انتظار برای 5 دقیقه
3. خطای timeout

**رفتار مورد انتظار**
فایل باید دانلود شود.

**رفتار واقعی**
خطا: `❌ Timeout after 5m0s`

**محیط**
- OS: Ubuntu 22.04
- Go: 1.21.0
- melifetch: 2.0.0

**لاگ‌ها**
```
[لاگ‌ها اینجا]
```
```

---

## 💡 پیشنهاد ویژگی

برای پیشنهاد ویژگی جدید، یک Issue با برچسب `enhancement` باز کنید.

### اطلاعات مورد نیاز:

- **توضیح ویژگی**: چه چیزی می‌خواهید اضافه شود
- **دلیل**: چرا این ویژگی مفید است
- **مثال استفاده**: چگونه از آن استفاده می‌شود
- **جایگزین‌ها**: آیا راه حل دیگری وجود دارد

---

## 💻 توسعه کد

### راه‌اندازی محیط توسعه

```bash
# Clone
git clone https://github.com/Keyaru-main/melli_pet.git
cd melli_pet

# نصب dependencies
go mod download

# ساخت
make build

# اجرا
./melifetch --help
```

### ساختار پروژه

```
melli_pet/
├── cmd/melifetch/       # CLI entry point
├── internal/
│   ├── config/         # Configuration management
│   └── github/         # GitHub API & download logic
├── repo/
│   ├── .github/workflows/  # GitHub Actions workflows
│   └── scripts/        # Helper scripts
└── *.md                # Documentation
```

### اضافه کردن ویژگی جدید

1. **تعریف interface** در `internal/`
2. **پیاده‌سازی** منطق
3. **اضافه کردن command** در `cmd/melifetch/main.go`
4. **نوشتن تست**
5. **به‌روزرسانی مستندات**

---

## 📏 استانداردهای کد

### Go Code Style

- از `gofmt` برای فرمت کردن استفاده کنید
- از `golangci-lint` برای lint استفاده کنید
- نام‌گذاری: camelCase برای متغیرها، PascalCase برای export شده‌ها
- کامنت‌های مفید بنویسید
- خطاها را به درستی handle کنید

### مثال کد خوب:

```go
// FetchFile downloads a file from the given URL
func (f *Fetcher) FetchFile(url string, opts FetchOptions) error {
    if url == "" {
        return fmt.Errorf("URL cannot be empty")
    }
    
    // Validate options
    if err := opts.Validate(); err != nil {
        return fmt.Errorf("invalid options: %w", err)
    }
    
    // Trigger workflow
    branchName, err := f.TriggerFetchWorkflow(opts)
    if err != nil {
        return fmt.Errorf("failed to trigger workflow: %w", err)
    }
    
    return nil
}
```

### Commit Messages

از Conventional Commits استفاده کنید:

```
feat: add support for custom headers
fix: resolve timeout issue in large downloads
docs: update README with new examples
test: add tests for fetcher
refactor: simplify error handling
chore: update dependencies
```

### مثال‌های خوب:

```
feat(download): add support for custom User-Agent
fix(fetch): handle redirect loops properly
docs(readme): add CurseForge download example
test(github): add unit tests for branch management
```

---

## 🧪 تست

### اجرای تست‌ها

```bash
# تمام تست‌ها
make test

# با coverage
make test-coverage

# تست یک package خاص
go test -v ./internal/github/
```

### نوشتن تست

```go
func TestFetcherNew(t *testing.T) {
    tests := []struct {
        name    string
        token   string
        repo    string
        wantErr bool
    }{
        {
            name:    "valid input",
            token:   "token123",
            repo:    "user/repo",
            wantErr: false,
        },
        {
            name:    "invalid repo format",
            token:   "token123",
            repo:    "invalid",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewFetcher(tt.token, tt.repo)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFetcher() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## 📝 Pull Request

### قبل از ارسال PR

- [ ] کد را فرمت کنید: `make fmt`
- [ ] تست‌ها را اجرا کنید: `make test`
- [ ] مستندات را به‌روز کنید
- [ ] CHANGELOG را به‌روز کنید (اگر لازم است)
- [ ] commit message ها را بررسی کنید

### الگوی PR

```markdown
## توضیحات
توضیح مختصر از تغییرات

## نوع تغییر
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## چک‌لیست
- [ ] کد فرمت شده است
- [ ] تست‌ها pass می‌شوند
- [ ] مستندات به‌روز شده
- [ ] تست‌های جدید اضافه شده (اگر لازم است)

## تست
چگونه این تغییرات را تست کردید؟

## اسکرین‌شات (اگر لازم است)
```

### بررسی PR

PR شما توسط maintainer ها بررسی می‌شود. ممکن است:

- سوالاتی پرسیده شود
- تغییراتی درخواست شود
- پیشنهاداتی داده شود

لطفاً صبور باشید و به feedback ها پاسخ دهید.

---

## 🎨 راهنمای Style

### کد Go

```go
// خوب ✓
func (f *Fetcher) Download(url string) error {
    if url == "" {
        return ErrEmptyURL
    }
    return f.download(url)
}

// بد ✗
func (f *Fetcher) Download(url string) error {
if url=="" {
return errors.New("empty url")
}
return f.download(url)}
```

### مستندات

- از markdown استفاده کنید
- عناوین واضح
- مثال‌های کاربردی
- لینک‌های مفید

### کامنت‌ها

```go
// خوب ✓
// FetchFile downloads a file from the given URL using GitHub Actions.
// It returns an error if the download fails or times out.
func FetchFile(url string) error

// بد ✗
// fetch file
func FetchFile(url string) error
```

---

## 🏆 تشکر از مشارکت‌کنندگان

همه مشارکت‌کنندگان در فایل [CONTRIBUTORS.md](CONTRIBUTORS.md) ذکر می‌شوند.

---

## 📞 ارتباط

- **Issues**: [GitHub Issues](https://github.com/Keyaru-main/melli_pet/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Keyaru-main/melli_pet/discussions)

---

## 📄 لایسنس

با مشارکت در این پروژه، شما موافقت می‌کنید که مشارکت‌های شما تحت [MIT License](LICENSE) منتشر شوند.

---

<div align="center">

**از مشارکت شما متشکریم! 🙏**

</div>
