// Package web Webレイヤー ミドルウェア
package web

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	authDomain "shiftmaster/internal/modules/auth/domain"
)

// Middleware ミドルウェア型
type Middleware func(http.Handler) http.Handler

// Chain ミドルウェアチェーン
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// Logger リクエストログミドルウェア
func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// レスポンスラッパー
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			logger.Info("HTTPリクエスト",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration", time.Since(start).String(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

// Recover パニックリカバリーミドルウェア
func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("パニック発生", "error", err, "path", r.URL.Path)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// CORS CORSミドルウェア
func CORS(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// オリジンチェック
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, HX-Request, HX-Target, HX-Trigger")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// プリフライトリクエスト処理
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders セキュリティヘッダーミドルウェア
func SecurityHeaders() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter レスポンスラッパー ステータスコード取得用
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader ステータスコード記録
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// contextKey コンテキストキー型
type contextKey string

const (
	// ContextKeyClaims クレームキー
	ContextKeyClaims contextKey = "claims"
)

// TokenValidator トークン検証インターフェース
type TokenValidator interface {
	ValidateAccessToken(token string) (*authDomain.Claims, error)
}

// Auth 認証ミドルウェア
func Auth(validator TokenValidator, logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authorizationヘッダー取得
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Cookieからトークン取得を試みる
				cookie, err := r.Cookie("access_token")
				if err != nil || cookie.Value == "" {
					http.Error(w, "認証が必要です", http.StatusUnauthorized)
					return
				}
				authHeader = "Bearer " + cookie.Value
			}

			// Bearer トークン抽出
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "不正な認証ヘッダーです", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// トークン検証
			claims, err := validator.ValidateAccessToken(token)
			if err != nil {
				logger.Warn("トークン検証失敗", "error", err)
				http.Error(w, "無効なトークンです", http.StatusUnauthorized)
				return
			}

			// コンテキストにクレーム設定
			ctx := context.WithValue(r.Context(), ContextKeyClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthOptional オプション認証ミドルウェア 認証なしでもアクセス可能
func AuthOptional(validator TokenValidator, logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				cookie, err := r.Cookie("access_token")
				if err == nil && cookie.Value != "" {
					authHeader = "Bearer " + cookie.Value
				}
			}

			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					token := parts[1]
					claims, err := validator.ValidateAccessToken(token)
					if err == nil {
						ctx := context.WithValue(r.Context(), ContextKeyClaims, claims)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole ロール要求ミドルウェア（認可）
func RequireRole(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
			if claims == nil {
				http.Error(w, "認証が必要です", http.StatusUnauthorized)
				return
			}

			// ロールチェック
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "この操作を行う権限がありません", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireSuperAdmin スーパー管理者要求ミドルウェア（全テナント管理）
func RequireSuperAdmin() Middleware {
	return RequireRole("super_admin")
}

// RequireAdmin テナント管理者以上要求ミドルウェア（認可）
func RequireAdmin() Middleware {
	return RequireRole("super_admin", "admin")
}

// RequireManager マネージャー以上要求ミドルウェア（認可）
func RequireManager() Middleware {
	return RequireRole("super_admin", "admin", "manager")
}

// GetClaimsFromContext コンテキストからクレーム取得
func GetClaimsFromContext(ctx context.Context) *authDomain.Claims {
	claims, ok := ctx.Value(ContextKeyClaims).(*authDomain.Claims)
	if !ok {
		return nil
	}
	return claims
}

// RedirectIfNotAuthenticated 未認証時リダイレクトミドルウェア ページ用
func RedirectIfNotAuthenticated(validator TokenValidator, loginPath string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				cookie, err := r.Cookie("access_token")
				if err != nil || cookie.Value == "" {
					http.Redirect(w, r, loginPath, http.StatusFound)
					return
				}
				authHeader = "Bearer " + cookie.Value
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Redirect(w, r, loginPath, http.StatusFound)
				return
			}

			token := parts[1]
			claims, err := validator.ValidateAccessToken(token)
			if err != nil {
				http.Redirect(w, r, loginPath, http.StatusFound)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RedirectIfAuthenticated 認証済み時リダイレクトミドルウェア ログインページ用
func RedirectIfAuthenticated(validator TokenValidator, redirectPath string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err == nil && cookie.Value != "" {
				_, err := validator.ValidateAccessToken(cookie.Value)
				if err == nil {
					http.Redirect(w, r, redirectPath, http.StatusFound)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
