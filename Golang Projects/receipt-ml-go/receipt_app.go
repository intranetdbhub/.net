package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db        *sql.DB
	templates *template.Template
	uploadDir = "uploads"
)

/* =========================
   Models
   ========================= */

type Receipt struct {
	ID           int64
	Merchant     *string
	PurchaseDate *time.Time
	Subtotal     *float64
	Tax          *float64
	Fees         *float64
	Total        *float64
	Currency     *string
	NotVisible   map[string]bool
	OCRText      string
	OriginalFile string
	StoragePath  string
	Status       string
	CreatedAt    time.Time
	Items        []ReceiptItem
	Pages        []ReceiptPage
}

type ReceiptItem struct {
	ID        int64
	ReceiptID int64
	Label     string
	Quantity  float64
	UnitPrice *float64
	LineTotal *float64
	IsFee     bool
}

type ReceiptPage struct {
	ID        int64
	ReceiptID int64
	PageIndex int
	ImagePath string
}

/* =========================
   Boot
   ========================= */

func main() {
	var err error
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN env var is required, e.g. user:pass@tcp(127.0.0.1:3306)/receiptsdb?parseTime=true&charset=utf8mb4")
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	if err = os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", listHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/upload-multi", uploadMultiHandler)
	http.HandleFunc("/receipt/", viewHandler) // /receipt/{id}

	// Exports
	http.HandleFunc("/export/receipts.csv", exportReceiptsCSV)
	http.HandleFunc("/export/receipt_items.csv", exportItemsCSV)
	http.HandleFunc("/export/receipts.json", exportReceiptsJSON)
	http.HandleFunc("/receiptjson/", exportReceiptJSONHandler) // /receiptjson/{id}
	http.HandleFunc("/receiptcsv/", exportReceiptCSVHandler)   // /receiptcsv/{id}

	// Static files (CSS) and uploaded assets (images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	addr := ":8080"
	log.Printf("Receipt OCR app running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

/* =========================
   Handlers
   ========================= */

type listRow struct {
	ID           int64
	Merchant     sql.NullString
	PurchaseDate sql.NullTime
	Total        sql.NullFloat64
	Status       string
	OriginalFile sql.NullString
	CreatedAt    time.Time
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, merchant, purchase_date, total, status, original_filename, created_at
		FROM receipts ORDER BY created_at DESC LIMIT 500`)
	if err != nil {
		httpError(w, err)
		return
	}
	defer rows.Close()

	var list []listRow
	for rows.Next() {
		var row listRow
		if err := rows.Scan(&row.ID, &row.Merchant, &row.PurchaseDate, &row.Total, &row.Status, &row.OriginalFile, &row.CreatedAt); err != nil {
			httpError(w, err)
			return
		}
		list = append(list, row)
	}
	render(w, "list.html", map[string]any{"Receipts": list})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "upload.html", nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		httpError(w, err)
		return
	}
	// Extract file
	f, fh, err := r.FormFile("receipt")
	if err != nil {
		httpError(w, errors.New("please choose a file"))
		return
	}
	defer f.Close()

	receiptID, err := processOneUpload(f, fh.Filename)
	if err != nil {
		httpError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/receipt/%d", receiptID), http.StatusSeeOther)
}

// Bulk upload: input name="receipts" multiple
func uploadMultiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "upload.html", map[string]any{"Multi": true})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(256 << 20); err != nil {
		httpError(w, err)
		return
	}
	mf := r.MultipartForm
	files := mf.File["receipts"]
	if len(files) == 0 {
		httpError(w, errors.New("no files provided"))
		return
	}

	var firstID int64
	for i, fh := range files {
		f, err := fh.Open()
		if err != nil {
			httpError(w, err)
			return
		}
		recID, err := processOneUpload(f, fh.Filename)
		f.Close()
		if err != nil {
			httpError(w, fmt.Errorf("upload failed for %s: %w", fh.Filename, err))
			return
		}
		if i == 0 {
			firstID = recID
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/receipt/%d", firstID), http.StatusSeeOther)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// URL: /receipt/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/receipt/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	idStr := parts[0]
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if r.Method == http.MethodPost {
		if err := saveEdits(r, id); err != nil {
			httpError(w, err)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	rec, items, pages, err := loadReceipt(id)
	if err != nil {
		httpError(w, err)
		return
	}
	render(w, "view.html", map[string]any{"R": rec, "Items": items, "Pages": pages})
}

/* =========================
   Upload core
   ========================= */

func processOneUpload(file multipart.File, filename string) (int64, error) {
	// Save original
	ext := strings.ToLower(filepath.Ext(filename))
	tempName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), sanitizeName(filename))
	destPath := filepath.Join(uploadDir, tempName)
	out, err := os.Create(destPath)
	if err != nil {
		return 0, err
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		return 0, err
	}
	out.Close()

	// Normalize to images
	imgPaths, err := normalizeToImages(destPath, ext)
	if err != nil {
		return 0, fmt.Errorf("convert error: %w", err)
	}

	// OCR all pages
	var ocrBuf strings.Builder
	for _, img := range imgPaths {
		text, err := ocrImage(img)
		if err != nil {
			return 0, fmt.Errorf("ocr error on %s: %w", img, err)
		}
		ocrBuf.WriteString(text)
		ocrBuf.WriteString("\n")
	}
	ocrText := normalizeWhitespace(ocrBuf.String())

	// Parse
	parsed := parseReceipt(ocrText)

	// Insert receipt + items + pages
	receiptID, err := insertReceipt(filename, toSlash(destPath), ocrText, parsed)
	if err != nil {
		return 0, err
	}
	if err := insertItems(receiptID, parsed.Items); err != nil {
		return 0, err
	}
	if err := insertPages(receiptID, imgPaths); err != nil {
		return 0, err
	}
	return receiptID, nil
}

/* =========================
   DB Ops
   ========================= */

func insertReceipt(origName, storagePath, ocr string, p ParsedReceipt) (int64, error) {
	nv, _ := json.Marshal(p.NotVisible)
	res, err := db.Exec(`INSERT INTO receipts
	(merchant, purchase_date, subtotal, tax, fees, total, currency, not_visible, ocr_text, original_filename, storage_path, status)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		nilOrStrPtr(p.Merchant),
		nilOrDatePtr(p.PurchaseDate),
		nilOrFloatPtr(p.Subtotal),
		nilOrFloatPtr(p.Tax),
		nilOrFloatPtr(p.Fees),
		nilOrFloatPtr(p.Total),
		nilOrStrPtr(p.Currency),
		string(nv),
		ocr,
		origName,
		storagePath,
		p.Status,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertItems(receiptID int64, items []ParsedItem) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`INSERT INTO receipt_items
	(receipt_id, label, quantity, unit_price, line_total, is_fee)
	VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, it := range items {
		_, err = stmt.Exec(receiptID, it.Label, it.Quantity, nilOrFloatPtr(it.UnitPrice), nilOrFloatPtr(it.LineTotal), it.IsFee)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func insertPages(receiptID int64, imgPaths []string) error {
	if len(imgPaths) == 0 {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`INSERT INTO receipt_pages (receipt_id, page_index, image_path) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i, p := range imgPaths {
		_, err = stmt.Exec(receiptID, i, toSlash(p))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func loadReceipt(id int64) (Receipt, []ReceiptItem, []ReceiptPage, error) {
	var rec Receipt
	var notVisRaw sql.NullString
	err := db.QueryRow(`SELECT id, merchant, purchase_date, subtotal, tax, fees, total, currency, not_visible, ocr_text, original_filename, storage_path, status, created_at
		FROM receipts WHERE id=?`, id).
		Scan(&rec.ID, &nullableString{&rec.Merchant}, &nullableTime{&rec.PurchaseDate}, &nullableFloat{&rec.Subtotal},
			&nullableFloat{&rec.Tax}, &nullableFloat{&rec.Fees}, &nullableFloat{&rec.Total}, &nullableString{&rec.Currency},
			&notVisRaw, &rec.OCRText, &rec.OriginalFile, &rec.StoragePath, &rec.Status, &rec.CreatedAt)
	if err != nil {
		return rec, nil, nil, err
	}
	rec.NotVisible = map[string]bool{}
	if notVisRaw.Valid && notVisRaw.String != "" {
		_ = json.Unmarshal([]byte(notVisRaw.String), &rec.NotVisible)
	}

	rows, err := db.Query(`SELECT id, receipt_id, label, quantity, unit_price, line_total, is_fee
		FROM receipt_items WHERE receipt_id=? ORDER BY id ASC`, id)
	if err != nil {
		return rec, nil, nil, err
	}
	defer rows.Close()
	var items []ReceiptItem
	for rows.Next() {
		var it ReceiptItem
		if err := rows.Scan(&it.ID, &it.ReceiptID, &it.Label, &it.Quantity, &nullableFloat{&it.UnitPrice}, &nullableFloat{&it.LineTotal}, &it.IsFee); err != nil {
			return rec, nil, nil, err
		}
		items = append(items, it)
	}

	pgs, err := db.Query(`SELECT id, receipt_id, page_index, image_path FROM receipt_pages WHERE receipt_id=? ORDER BY page_index ASC`, id)
	if err != nil {
		return rec, nil, nil, err
	}
	defer pgs.Close()
	var pages []ReceiptPage
	for pgs.Next() {
		var pg ReceiptPage
		if err := pgs.Scan(&pg.ID, &pg.ReceiptID, &pg.PageIndex, &pg.ImagePath); err != nil {
			return rec, nil, nil, err
		}
		pages = append(pages, pg)
	}
	return rec, items, pages, nil
}

func saveEdits(r *http.Request, id int64) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	merchant := strings.TrimSpace(r.FormValue("merchant"))
	dateStr := strings.TrimSpace(r.FormValue("purchase_date"))
	subtotal := parseFormMoney(r.FormValue("subtotal"))
	tax := parseFormMoney(r.FormValue("tax"))
	fees := parseFormMoney(r.FormValue("fees"))
	total := parseFormMoney(r.FormValue("total"))
	currency := strings.TrimSpace(r.FormValue("currency"))

	var dt *time.Time
	if dateStr != "" {
		if t, err := parseDateAny(dateStr); err == nil {
			dt = &t
		}
	}
	_, err := db.Exec(`UPDATE receipts SET merchant=?, purchase_date=?, subtotal=?, tax=?, fees=?, total=?, currency=?, status='parsed' WHERE id=?`,
		nullIfEmpty(merchant), dt, subtotal, tax, fees, total, nullIfEmpty(currency), id)
	return err
}

/* =========================
   OCR + Conversions
   ========================= */

func normalizeToImages(path, ext string) ([]string, error) {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	switch ext {
	case "jpg", "jpeg", "png":
		return []string{path}, nil
	case "heic", "heif":
		out := strings.TrimSuffix(path, "."+ext) + ".jpg"
		if err := convertWithMagick(path, out); err != nil {
			if err2 := convertWithHeifConvert(path, out); err2 != nil {
				return nil, fmt.Errorf("HEIC->JPG failed: magick: %v, heif-convert: %v", err, err2)
			}
		}
		return []string{out}, nil
	case "pdf":
		base := strings.TrimSuffix(path, "."+ext)
		if err := convertPDFWithPdftoppm(path, base); err != nil {
			outGlob := base + "-%d.png"
			if err2 := convertPDFWithMagick(path, outGlob); err2 != nil {
				return nil, fmt.Errorf("PDF->PNG failed: pdftoppm: %v, magick: %v", err, err2)
			}
		}
		var imgs []string
		// collect outBase-n.png (poppler) or outBase-0.png,1.png (magick)
		for i := 0; i < 1000; i++ {
			fp := fmt.Sprintf("%s-%d.png", base, i)
			if _, err := os.Stat(fp); err == nil {
				imgs = append(imgs, fp)
				continue
			}
			// also try i+1 in case numbering started at 1
			fp1 := fmt.Sprintf("%s-%d.png", base, i+1)
			if _, err := os.Stat(fp1); err == nil {
				imgs = append(imgs, fp1)
				continue
			}
			if i > 0 {
				break
			}
		}
		if len(imgs) == 0 {
			return nil, errors.New("no PDF pages converted")
		}
		return imgs, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func ocrImage(imgPath string) (string, error) {
	if _, err := exec.LookPath("tesseract"); err != nil {
		return "", errors.New("tesseract not found in PATH")
	}
	cmd := exec.Command("tesseract", imgPath, "stdout", "--psm", "6")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("tesseract error: %v, output: %s", err, out.String())
	}
	return out.String(), nil
}

func convertPDFWithPdftoppm(pdfPath, outBase string) error {
	if _, err := exec.LookPath("pdftoppm"); err != nil {
		return fmt.Errorf("pdftoppm not found")
	}
	cmd := exec.Command("pdftoppm", "-png", pdfPath, outBase)
	return cmd.Run()
}

func convertPDFWithMagick(pdfPath, outGlob string) error {
	if _, err := exec.LookPath("magick"); err != nil {
		return fmt.Errorf("magick not found")
	}
	cmd := exec.Command("magick", "-density", "300", pdfPath, outGlob)
	return cmd.Run()
}

func convertWithMagick(in, out string) error {
	if _, err := exec.LookPath("magick"); err != nil {
		return fmt.Errorf("magick not found")
	}
	cmd := exec.Command("magick", in, out)
	return cmd.Run()
}

func convertWithHeifConvert(in, out string) error {
	if _, err := exec.LookPath("heif-convert"); err != nil {
		return fmt.Errorf("heif-convert not found")
	}
	cmd := exec.Command("heif-convert", in, out)
	return cmd.Run()
}

/* =========================
   Parsing (heuristic MVP)
   ========================= */

type ParsedReceipt struct {
	Merchant     *string
	PurchaseDate *time.Time
	Subtotal     *float64
	Tax          *float64
	Fees         *float64
	Total        *float64
	Currency     *string
	Status       string
	Items        []ParsedItem
	NotVisible   map[string]bool
}

type ParsedItem struct {
	Label     string
	Quantity  float64
	UnitPrice *float64
	LineTotal *float64
	IsFee     bool
}

func parseReceipt(ocr string) ParsedReceipt {
	lines := splitNonEmptyLines(ocr)
	p := ParsedReceipt{Status: "needs_review", NotVisible: map[string]bool{}}

	p.Merchant = guessMerchant(lines)
	if p.Merchant == nil {
		p.NotVisible["merchant"] = true
	}
	if dt := findDate(lines); dt != nil {
		p.PurchaseDate = dt
	} else {
		p.NotVisible["purchase_date"] = true
	}
	if cur := findCurrency(lines); cur != nil {
		p.Currency = cur
	}

	subtotal, tax, fees, total := findTotals(lines)
	if subtotal == nil {
		p.NotVisible["subtotal"] = true
	}
	if tax == nil {
		p.NotVisible["tax"] = true
	}
	if fees == nil {
		p.NotVisible["fees"] = true
	}
	if total == nil {
		p.NotVisible["total"] = true
	}
	p.Subtotal, p.Tax, p.Fees, p.Total = subtotal, tax, fees, total
	p.Items = findItems(lines)
	return p
}

func guessMerchant(lines []string) *string {
	skip := regexp.MustCompile(`(?i)(receipt|invoice|thank|order|visa|mastercard|debit|credit|cash|change|total|subtotal|tax|date)`)
	for i := 0; i < len(lines) && i < 8; i++ {
		L := strings.TrimSpace(lines[i])
		if L == "" {
			continue
		}
		if len([]rune(L)) < 3 {
			continue
		}
		if skip.MatchString(L) {
			continue
		}
		upperRatio := ratioUpper(L)
		if upperRatio > 0.4 || hasBizToken(L) {
			s := cleanRightPrice(L)
			return &s
		}
	}
	return nil
}

func findDate(lines []string) *time.Time {
	datePatterns := []string{
		`(\d{4})[-/](\d{1,2})[-/](\d{1,2})`,
		`(\d{1,2})[-/](\d{1,2})[-/](\d{2,4})`,
		`(?i)(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*[ .-]+(\d{1,2})[, ]+(\d{2,4})`,
		`(\d{1,2})[ .-]+(?i)(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*[, ]+(\d{2,4})`,
	}
	for _, L := range lines {
		for _, pat := range datePatterns {
			re := regexp.MustCompile(pat)
			m := re.FindString(L)
			if m != "" {
				if t, err := parseDateAny(m); err == nil {
					return &t
				}
			}
		}
	}
	return nil
}

func parseDateAny(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	formats := []string{
		"2006-1-2", "2006/1/2",
		"1/2/2006", "01/02/2006",
		"2/1/2006", "02/01/2006",
		"1-2-2006", "01-02-2006", "2-1-2006",
		"Jan 2 2006", "Jan 2, 2006", "2 Jan 2006", "02 Jan 2006",
		"2006.01.02", "1.2.2006",
		"01-02-06", "1-2-06", "01/02/06", "1/2/06",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	s = strings.ReplaceAll(s, ",", "")
	for _, f := range []string{"Jan 2 2006", "2 Jan 2006"} {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown date format: %s", s)
}

func findCurrency(lines []string) *string {
	cands := []string{"USD", "$", "EUR", "€", "GBP", "£", "MXN"}
	for _, L := range lines {
		for _, c := range cands {
			if strings.Contains(L, c) {
				cc := c
				return &cc
			}
		}
	}
	return nil
}

var moneyRE = regexp.MustCompile(`[-+]?\d{1,3}(?:,\d{3})*(?:\.\d{2})|\d+\.\d{2}`)

func findTotals(lines []string) (subtotal, tax, fees, total *float64) {
	for _, L := range lines {
		LL := strings.ToLower(L)
		if strings.Contains(LL, "subtotal") {
			if v := lastMoney(L); v != nil {
				subtotal = v
			}
		}
		if strings.Contains(LL, "tax") {
			if v := lastMoney(L); v != nil {
				tax = v
			}
		}
		if strings.Contains(LL, "fee") || strings.Contains(LL, "fees") || strings.Contains(LL, "surcharge") || strings.Contains(LL, "service") {
			if v := lastMoney(L); v != nil {
				fees = v
			}
		}
		if strings.Contains(LL, "total") || strings.Contains(LL, "amount due") || strings.Contains(LL, "grand total") {
			if v := lastMoney(L); v != nil {
				total = v
			}
		}
	}
	return
}

func findItems(lines []string) []ParsedItem {
	var items []ParsedItem
	itemLine := regexp.MustCompile(`^(?P<label>.+?)\s+(?P<price>[-+]?\d{1,3}(?:,\d{3})*(?:\.\d{2})|\d+\.\d{2})$`)
	skip := regexp.MustCompile(`(?i)(subtotal|total|tax|change|cash|visa|mastercard|debit|credit|balance|amount due|thank)`)
	for _, L := range lines {
		if skip.MatchString(L) {
			continue
		}
		L = strings.TrimSpace(L)
		if L == "" {
			continue
		}
		if m := itemLine.FindStringSubmatch(L); m != nil {
			label := strings.TrimSpace(m[1])
			price := parseMoney(m[2])
			items = append(items, ParsedItem{
				Label:     label,
				Quantity:  1.0,
				UnitPrice: price,
				LineTotal: price,
				IsFee:     looksLikeFee(label),
			})
		}
	}
	return items
}

func looksLikeFee(label string) bool {
	l := strings.ToLower(label)
	return strings.Contains(l, "fee") || strings.Contains(l, "surcharge") || strings.Contains(l, "service")
}

/* =========================
   Exports
   ========================= */

func exportReceiptsCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=receipts.csv")

	rows, err := db.Query(`SELECT id, merchant, purchase_date, subtotal, tax, fees, total, currency, original_filename, created_at FROM receipts ORDER BY id`)
	if err != nil {
		httpError(w, err)
		return
	}
	defer rows.Close()

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"id", "merchant", "purchase_date", "subtotal", "tax", "fees", "total", "currency", "original_filename", "created_at"})
	for rows.Next() {
		var id int64
		var merchant sql.NullString
		var dt sql.NullTime
		var subtotal, tax, fees, total sql.NullFloat64
		var currency, orig sql.NullString
		var created time.Time
		if err := rows.Scan(&id, &merchant, &dt, &subtotal, &tax, &fees, &total, &currency, &orig, &created); err != nil {
			httpError(w, err)
			return
		}
		_ = cw.Write([]string{
			strconv.FormatInt(id, 10), ns(merchant), nt(dt), nf(subtotal), nf(tax), nf(fees), nf(total), ns(currency), ns(orig), created.Format(time.RFC3339),
		})
	}
	cw.Flush()
}

func exportItemsCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=receipt_items.csv")

	rows, err := db.Query(`SELECT id, receipt_id, label, quantity, unit_price, line_total, is_fee FROM receipt_items ORDER BY receipt_id, id`)
	if err != nil {
		httpError(w, err)
		return
	}
	defer rows.Close()

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"id", "receipt_id", "label", "quantity", "unit_price", "line_total", "is_fee"})
	for rows.Next() {
		var id, rid int64
		var label string
		var qty float64
		var unit, line sql.NullFloat64
		var isFee bool
		if err := rows.Scan(&id, &rid, &label, &qty, &unit, &line, &isFee); err != nil {
			httpError(w, err)
			return
		}
		_ = cw.Write([]string{
			strconv.FormatInt(id, 10), strconv.FormatInt(rid, 10), label, fmt.Sprintf("%.2f", qty), nf(unit), nf(line), strconv.FormatBool(isFee),
		})
	}
	cw.Flush()
}

func exportReceiptsJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT id FROM receipts ORDER BY id`)
	if err != nil {
		httpError(w, err)
		return
	}
	defer rows.Close()
	var out []any
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			httpError(w, err)
			return
		}
		rec, items, pages, err := loadReceipt(id)
		if err != nil {
			httpError(w, err)
			return
		}
		out = append(out, map[string]any{"receipt": rec, "items": items, "pages": pages})
	}
	_ = json.NewEncoder(w).Encode(out)
}

func exportReceiptJSONHandler(w http.ResponseWriter, r *http.Request) {
	// /receiptjson/{id}
	id, ok := tailID(r.URL.Path, "/receiptjson/")
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	rec, items, pages, err := loadReceipt(id)
	if err != nil {
		httpError(w, err)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"receipt": rec, "items": items, "pages": pages})
}

func exportReceiptCSVHandler(w http.ResponseWriter, r *http.Request) {
	// /receiptcsv/{id}
	id, ok := tailID(r.URL.Path, "/receiptcsv/")
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=receipt_%d.csv", id))

	rec, items, pages, err := loadReceipt(id)
	if err != nil {
		httpError(w, err)
		return
	}
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"field", "value"})
	_ = cw.Write([]string{"id", strconv.FormatInt(rec.ID, 10)})
	_ = cw.Write([]string{"merchant", sv(rec.Merchant)})
	_ = cw.Write([]string{"purchase_date", st(rec.PurchaseDate)})
	_ = cw.Write([]string{"subtotal", sf(rec.Subtotal)})
	_ = cw.Write([]string{"tax", sf(rec.Tax)})
	_ = cw.Write([]string{"fees", sf(rec.Fees)})
	_ = cw.Write([]string{"total", sf(rec.Total)})
	_ = cw.Write([]string{"currency", sv(rec.Currency)})
	_ = cw.Write([]string{"original_filename", rec.OriginalFile})
	_ = cw.Write([]string{"created_at", rec.CreatedAt.Format(time.RFC3339)})
	_ = cw.Write([]string{"pages", fmt.Sprintf("%d", len(pages))})
	_ = cw.Write([]string{"items", fmt.Sprintf("%d", len(items))})
	_ = cw.Write([]string{"status", rec.Status})

	_ = cw.Write([]string{"", ""})
	_ = cw.Write([]string{"items_header", "label|qty|unit|line|is_fee"})
	for _, it := range items {
		_ = cw.Write([]string{"item", fmt.Sprintf("%s|%.2f|%s|%s|%t", it.Label, it.Quantity, pf(it.UnitPrice), pf(it.LineTotal), it.IsFee)})
	}
	cw.Flush()
}

/* =========================
   Helpers
   ========================= */

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Always execute the base layout; the page's {{define "content"}} overrides the block.
	if err := templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		httpError(w, err)
	}
}

func httpError(w http.ResponseWriter, err error) {
	log.Println("[error]", err)
	http.Error(w, err.Error(), 500)
}

func sanitizeName(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	return regexp.MustCompile(`[^a-zA-Z0-9._-]`).ReplaceAllString(s, "")
}

func splitNonEmptyLines(s string) []string {
	raw := strings.Split(s, "\n")
	var out []string
	for _, L := range raw {
		L = strings.TrimSpace(L)
		if L != "" {
			out = append(out, L)
		}
	}
	return out
}

func ratioUpper(s string) float64 {
	count := 0
	total := 0
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			count++
		}
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			total++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(count) / float64(total)
}

func hasBizToken(s string) bool {
	s = strings.ToLower(s)
	toks := []string{"inc", "llc", "corp", "ltd", "co.", "restaurant", "store"}
	for _, t := range toks {
		if strings.Contains(s, t) {
			return true
		}
	}
	return false
}

func cleanRightPrice(s string) string {
	m := moneyRE.FindAllStringIndex(s, -1)
	if len(m) == 0 {
		return s
	}
	last := m[len(m)-1]
	if last[1] == len(s) {
		return strings.TrimSpace(s[:last[0]])
	}
	return s
}

func lastMoney(s string) *float64 {
	ms := moneyRE.FindAllString(s, -1)
	if len(ms) == 0 {
		return nil
	}
	return parseMoney(ms[len(ms)-1])
}

func parseMoney(s string) *float64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}

func normalizeWhitespace(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	re := regexp.MustCompile(`[ \t]+`)
	return re.ReplaceAllString(s, " ")
}

func parseFormMoney(s string) *float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return parseMoney(s)
}

func nullIfEmpty(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func nilOrStrPtr(s *string) *string        { return s }
func nilOrDatePtr(t *time.Time) *time.Time { return t }
func nilOrFloatPtr(f *float64) *float64    { return f }
func toSlash(p string) string              { return filepath.ToSlash(p) }

// scanners

type nullableString struct{ p **string }

func (n *nullableString) Scan(src any) error {
	if src == nil {
		*n.p = nil
		return nil
	}
	switch v := src.(type) {
	case []byte:
		s := string(v)
		*n.p = &s
	case string:
		s := v
		*n.p = &s
	default:
		return fmt.Errorf("bad type for string")
	}
	return nil
}

type nullableFloat struct{ p **float64 }

func (n *nullableFloat) Scan(src any) error {
	if src == nil {
		*n.p = nil
		return nil
	}
	switch v := src.(type) {
	case []byte:
		f, _ := strconv.ParseFloat(string(v), 64)
		*n.p = &f
	case float64:
		f := v
		*n.p = &f
	case int64:
		f := float64(v)
		*n.p = &f
	default:
		return fmt.Errorf("bad type for float")
	}
	return nil
}

type nullableTime struct{ p **time.Time }

func (n *nullableTime) Scan(src any) error {
	if src == nil {
		*n.p = nil
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		t := v
		*n.p = &t
	default:
		return fmt.Errorf("bad type for time")
	}
	return nil
}

// export helpers
func ns(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
func nt(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("2006-01-02")
	}
	return ""
}
func nf(f sql.NullFloat64) string {
	if f.Valid {
		return fmt.Sprintf("%.2f", f.Float64)
	}
	return ""
}
func sv(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}
func st(p *time.Time) string {
	if p != nil {
		return p.Format("2006-01-02")
	}
	return ""
}
func sf(p *float64) string {
	if p != nil {
		return fmt.Sprintf("%.2f", *p)
	}
	return ""
}
func pf(p *float64) string {
	if p != nil {
		return fmt.Sprintf("%.2f", *p)
	}
	return ""
}
func tailID(path, prefix string) (int64, bool) {
	s := strings.TrimPrefix(path, prefix)
	if s == path || s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
