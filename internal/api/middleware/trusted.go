package middleware

import (
	"log/slog"
	"net"
	"net/http"
)

// WithTrustedSubnet возвращает хендлер, обернутый в миддлваре проверки доверенной подсети.
func WithTrustedSubnet(trustedSubnet string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !isTructedSubnet(r.RemoteAddr, trustedSubnet) {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func isTructedSubnet(ip string, trustedSubnet string) bool {
	if ip == "" {
		return false
	}
	if trustedSubnet == "" {
		return false
	}

	_, IPNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		slog.Info("failed to parse trusted subnet", "error", err.Error())
		return false
	}

	IP := net.ParseIP(ip)

	return IPNet.Contains(IP)
}
