// Package web Webレイヤー テンプレートエンジン
package web

import (
	"bytes"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TemplateRenderer テンプレートレンダラーインターフェース
type TemplateRenderer interface {
	Render(w http.ResponseWriter, name string, data any) error
	RenderWithLayout(w http.ResponseWriter, name, layout string, data any) error
}

// TemplateEngine テンプレートエンジン
type TemplateEngine struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
	funcMap   template.FuncMap
	baseDir   string
}

// NewTemplateEngine テンプレートエンジン生成
func NewTemplateEngine(baseDir string) *TemplateEngine {
	return &TemplateEngine{
		templates: make(map[string]*template.Template),
		funcMap:   defaultFuncMap(),
		baseDir:   baseDir,
	}
}

// Load テンプレート読み込み
func (e *TemplateEngine) Load() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// レイアウトテンプレート
	layoutPattern := filepath.Join(e.baseDir, "layouts", "*.html")
	layouts, err := filepath.Glob(layoutPattern)
	if err != nil {
		return err
	}

	// コンポーネントテンプレート
	componentPattern := filepath.Join(e.baseDir, "components", "*.html")
	components, err := filepath.Glob(componentPattern)
	if err != nil {
		return err
	}

	// ページテンプレート
	pagesDir := filepath.Join(e.baseDir, "pages")
	err = filepath.WalkDir(pagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		// テンプレート名 pagesディレクトリからの相対パス
		relPath, err := filepath.Rel(e.baseDir, path)
		if err != nil {
			return err
		}

		// レイアウト + コンポーネント + ページ
		files := append(layouts, components...)
		files = append(files, path)

		tmpl, err := template.New(filepath.Base(path)).Funcs(e.funcMap).ParseFiles(files...)
		if err != nil {
			return err
		}

		e.templates[relPath] = tmpl
		return nil
	})

	return err
}

// Render テンプレートレンダリング TemplateRendererインターフェース実装
func (e *TemplateEngine) Render(w http.ResponseWriter, name string, data any) error {
	e.mu.RLock()
	tmpl, ok := e.templates[name]
	e.mu.RUnlock()

	if !ok {
		return &TemplateError{Name: name, Message: "テンプレートが見つかりません"}
	}

	// バッファリング エラー時のレスポンス汚染防止
	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, "base", data); err != nil {
		return &TemplateError{Name: name, Message: err.Error()}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := buf.WriteTo(w)
	return err
}

// RenderWithLayout 指定レイアウトでテンプレートレンダリング
func (e *TemplateEngine) RenderWithLayout(w http.ResponseWriter, name, layout string, data any) error {
	e.mu.RLock()
	tmpl, ok := e.templates[name]
	e.mu.RUnlock()

	if !ok {
		return &TemplateError{Name: name, Message: "テンプレートが見つかりません"}
	}

	// バッファリング エラー時のレスポンス汚染防止
	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, layout, data); err != nil {
		return &TemplateError{Name: name, Message: err.Error()}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := buf.WriteTo(w)
	return err
}

// RenderToWriter io.Writer向けレンダリング
func (e *TemplateEngine) RenderToWriter(w io.Writer, name string, data any) error {
	e.mu.RLock()
	tmpl, ok := e.templates[name]
	e.mu.RUnlock()

	if !ok {
		return &TemplateError{Name: name, Message: "テンプレートが見つかりません"}
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, "base", data); err != nil {
		return &TemplateError{Name: name, Message: err.Error()}
	}

	_, err := buf.WriteTo(w)
	return err
}

// RenderPartial 部分テンプレートレンダリング HTMX用
func (e *TemplateEngine) RenderPartial(w io.Writer, name, block string, data any) error {
	e.mu.RLock()
	tmpl, ok := e.templates[name]
	e.mu.RUnlock()

	if !ok {
		return &TemplateError{Name: name, Message: "テンプレートが見つかりません"}
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, block, data); err != nil {
		return &TemplateError{Name: name, Message: err.Error()}
	}

	_, err := buf.WriteTo(w)
	return err
}

// TemplateError テンプレートエラー
type TemplateError struct {
	Name    string
	Message string
}

// Error エラーメッセージ
func (e *TemplateError) Error() string {
	return "テンプレートエラー [" + e.Name + "]: " + e.Message
}

// jst 日本標準時タイムゾーン
var jst = time.FixedZone("JST", 9*60*60)

// toJST UTC時刻を日本時間に変換
func toJST(t time.Time) time.Time {
	return t.In(jst)
}

// defaultFuncMap デフォルトテンプレート関数
func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		// 安全なHTML出力
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		// 日付フォーマット YYYY/MM/DD（日本時間）
		"formatDate": func(t any) string {
			switch v := t.(type) {
			case string:
				// RFC3339形式からの変換を試行
				if parsed, err := time.Parse(time.RFC3339, v); err == nil {
					return toJST(parsed).Format("2006/01/02")
				}
				// YYYY-MM-DD形式からの変換を試行
				if parsed, err := time.Parse("2006-01-02", v); err == nil {
					return parsed.Format("2006/01/02")
				}
				return v
			case time.Time:
				return toJST(v).Format("2006/01/02")
			default:
				return ""
			}
		},
		// 日付フォーマット input用 YYYY-MM-DD（日本時間）
		"formatDateInput": func(t any) string {
			switch v := t.(type) {
			case string:
				// RFC3339形式からの変換を試行
				if parsed, err := time.Parse(time.RFC3339, v); err == nil {
					return toJST(parsed).Format("2006-01-02")
				}
				return v
			case time.Time:
				return toJST(v).Format("2006-01-02")
			default:
				return ""
			}
		},
		// 日時フォーマット YYYY/MM/DD HH:mm（日本時間）
		"formatDateTime": func(t any) string {
			switch v := t.(type) {
			case string:
				// RFC3339形式からの変換を試行
				if parsed, err := time.Parse(time.RFC3339, v); err == nil {
					return toJST(parsed).Format("2006/01/02 15:04")
				}
				return v
			case time.Time:
				return toJST(v).Format("2006/01/02 15:04")
			default:
				return ""
			}
		},
		// 時刻フォーマット HH:mm（日本時間）
		"formatTime": func(t any) string {
			switch v := t.(type) {
			case string:
				// HH:MM:SS形式
				if parsed, err := time.Parse("15:04:05", v); err == nil {
					return parsed.Format("15:04")
				}
				// HH:MM形式
				if parsed, err := time.Parse("15:04", v); err == nil {
					return parsed.Format("15:04")
				}
				return v
			case time.Time:
				return toJST(v).Format("15:04")
			default:
				return ""
			}
		},
		// 数値加算
		"add": func(a, b int) int {
			return a + b
		},
		// 数値減算
		"sub": func(a, b int) int {
			return a - b
		},
		// 数値乗算
		"mul": func(a, b int) int {
			return a * b
		},
		// シーケンス生成 start以上end以下
		"seq": func(start, end int) []int {
			result := make([]int, 0, end-start+1)
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		// イテレート生成 start以上end未満
		"iterate": func(start, end int) []int {
			result := make([]int, 0, end-start)
			for i := start; i < end; i++ {
				result = append(result, i)
			}
			return result
		},
		// 大文字変換
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		// 小文字変換
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		// タイトルケース変換
		"title": func(s string) string {
			return cases.Title(language.Japanese).String(s)
		},
	}
}
