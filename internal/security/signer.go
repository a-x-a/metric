// Package security реализует функции, связанные с безопасностью при передаче данных.
package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/adapter"
)

// Signer подписывает и проверяет переданные данные.
type Signer struct {
	key []byte
}

// New создаёт новый экземпляр Signer.
//
// Параметры:
//   - key: ключ для подписи.
//
// Возвращаемое значение:
//   - экземпляр Signer или nil в случае ошибки.
func NewSigner(key string) *Signer {
	if len(key) == 0 {
		return nil
	}

	return &Signer{[]byte(key)}
}

// Hash подписывает данные и возвращает хэш.
//
// Параметры:
//   - data: данные для подписи.
//
// Возвращаемое значение:
//   - хэш или ошибку.
func (s *Signer) Hash(data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, s.key)
	h.Write(data)

	return h.Sum(nil), nil
}

// Verify проверяет подпись данных.
//
// Параметры:
//   - data: данные для подписи.
//   - hash: хэш.
//
// Возвращаемое значение:
//   - true, если подпись верна, false в противном случае или ошибку.
func (s *Signer) Verify(data []byte, hash string) (bool, error) {
	mac1, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}

	mac2, err := s.Hash(data)
	if err != nil {
		return false, err
	}

	return hmac.Equal(mac1, mac2), nil
}

// SignerMiddleware HTTP middleware для проверки подписи данных.
//
// Параметры:
//   - log: логгер.
//   - key: ключ для подписи.
func SignerMiddleware(log *zap.Logger, key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sgnr := NewSigner(key)
			if sgnr == nil {
				next.ServeHTTP(w, r)
				return
			}

			hash := r.Header.Get("HashSHA256")
			if len(hash) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Info("SIGNER", zap.String("hash received", hash))

			buf, _ := io.ReadAll(r.Body)
			rdr1 := io.NopCloser(bytes.NewBuffer(buf))
			rdr2 := io.NopCloser(bytes.NewBuffer(buf))

			data := make([]adapter.RequestMetric, 0)
			if err := json.NewDecoder(rdr1).Decode(&data); err != nil {
				log.Info("SIGNER", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			b, err := json.Marshal(data)
			if err != nil {
				log.Info("SIGNER", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if ok, err := sgnr.Verify(b, hash); !ok || err != nil {
				log.Info("SIGNER", zap.String("hash is not valid", hash))
				log.Info("SIGNER", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Info("SIGNER", zap.String("hash is valid", hash))

			r.Body = rdr2
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
