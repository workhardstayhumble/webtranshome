package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const maxProductImageBytes = 2 << 20

func mediaDir() string { return filepath.Join("storage", "uploads", "products") }

func UploadProductImageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxProductImageBytes+1<<20)
		if err := r.ParseMultipartForm(3 << 20); err != nil {
			http.Error(w, "File terlalu besar atau multipart tidak valid", 413)
			return
		}
		file, head, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File wajib ada", 400)
			return
		}
		defer file.Close()
		if r.FormValue("jenis") != "gambar" {
			http.Error(w, "Hanya upload gambar produk yang didukung", 400)
			return
		}
		if head.Size > maxProductImageBytes {
			http.Error(w, "Maksimal 2 MB", 413)
			return
		}
		data, err := io.ReadAll(io.LimitReader(file, maxProductImageBytes+1))
		if err != nil || len(data) > maxProductImageBytes {
			http.Error(w, "Maksimal 2 MB", 413)
			return
		}
		detected := http.DetectContentType(data)
		ext := map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp"}[detected]
		if ext == "" {
			http.Error(w, "Hanya JPG, PNG, WebP", 415)
			return
		}
		if err := os.MkdirAll(mediaDir(), 0750); err != nil {
			http.Error(w, "Gagal menyiapkan penyimpanan", 500)
			return
		}
		id := make([]byte, 16)
		if _, err = rand.Read(id); err != nil {
			http.Error(w, "Gagal membuat nama file", 500)
			return
		}
		name := hex.EncodeToString(id) + ext
		target := filepath.Join(mediaDir(), name)
		if err := os.WriteFile(target, data, 0640); err != nil {
			http.Error(w, "Gagal menyimpan file", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"url": "/api/v1/uploads/" + name, "name": head.Filename, "size": len(data), "type": detected}})
	}
}

func ServeProductUploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/v1/uploads/")
		if name == "" || filepath.Base(name) != name || strings.Contains(name, "\\") {
			http.NotFound(w, r)
			return
		}
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" && ext != ".png" && ext != ".webp" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, filepath.Join(mediaDir(), name))
	}
}

type imageInput struct {
	ImageURL  string `json:"imageUrl"`
	IsPrimary bool   `json:"isPrimary"`
	SortOrder int64  `json:"sortOrder"`
}

func (i *imageInput) UnmarshalJSON(b []byte) error {
	var v struct {
		ImageURL   string `json:"image_url"`
		ImageURL2  string `json:"imageUrl"`
		IsPrimary  bool   `json:"is_primary"`
		IsPrimary2 bool   `json:"isPrimary"`
		SortOrder  int64  `json:"sort_order"`
		SortOrder2 int64  `json:"sortOrder"`
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.ImageURL = v.ImageURL
	if i.ImageURL == "" {
		i.ImageURL = v.ImageURL2
	}
	i.IsPrimary = v.IsPrimary || v.IsPrimary2
	i.SortOrder = v.SortOrder
	if v.SortOrder2 != 0 {
		i.SortOrder = v.SortOrder2
	}
	return nil
}

type imageOutput struct {
	ID        int64  `json:"id"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int64  `json:"sort_order"`
}

func replaceImages(r *http.Request, db sqlx.SqlConn, productID int64, images []imageInput) error {
	if err := validateImages(images); err != nil {
		return err
	}
	if images == nil {
		return nil
	}
	return db.TransactCtx(r.Context(), func(ctx context.Context, session sqlx.Session) error {
		var locked int64
		if err := session.QueryRowCtx(ctx, &locked, "SELECT id FROM crm_schema.products WHERE id=$1 FOR UPDATE", productID); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx, "DELETE FROM crm_schema.product_images WHERE product_id=$1", productID); err != nil {
			return err
		}
		for n, v := range images {
			if _, err := session.ExecCtx(ctx, "INSERT INTO crm_schema.product_images (product_id,image_url,is_primary,sort_order) VALUES ($1,$2,$3,$4)", productID, v.ImageURL, v.IsPrimary, n); err != nil {
				return err
			}
		}
		return nil
	})
}

func validateImages(images []imageInput) error {
	if images == nil {
		return nil
	}
	if len(images) > 20 {
		return errors.New("Maksimal 20 foto produk")
	}
	primary := 0
	for _, v := range images {
		if v.IsPrimary {
			primary++
		}
		if len(v.ImageURL) > 255 || !(strings.HasPrefix(v.ImageURL, "/api/v1/uploads/") || strings.HasPrefix(v.ImageURL, "http://localhost:8888/api/v1/uploads/")) {
			return errors.New("URL foto tidak diizinkan")
		}
	}
	if len(images) > 0 && primary != 1 {
		return errors.New("Wajib tepat satu foto utama")
	}
	return nil
}