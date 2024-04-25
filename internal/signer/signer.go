package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Signer struct {
	key []byte
}

func New(key string) *Signer {
	return &Signer{[]byte(key)}
}

func (s *Signer) Hash(data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, s.key)
	h.Write(data)

	return h.Sum(nil), nil
}

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

// func SignerMiddleware(log *zap.Logger) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			hash := r.Header.Get("HashSHA256")
// 			if len(hash) == 0 {
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			data := make([]adapter.RequestMetric, 0)
// 			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 				log.Info("SIGNER", zap.Error(err))
// 				// responseWithError(w, http.StatusBadRequest, err, log)
// 				return
// 			}
// 			log.Info("SIGNER", zap.String("hash", hash))
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
