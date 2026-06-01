package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// GetTLSConfig возвращает конфиг для mTLS
// isServer = true, если настраиваем сервер (требуем сертификат от клиента)
// isServer = false, если настраиваем клиент (предоставляем свой сертификат)
func GetTLSConfig(isServer bool) (*tls.Config, error) {
	// Загружаем наш корневой сертификат CA
	caCert, err := os.ReadFile("../../certs/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ca.crt: %w", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Загружаем сертификат и ключ (используем пару server для сервера, client для клиента)
	certFile := "../../certs/client.crt"
	keyFile := "../../certs/client.key"
	if isServer {
		certFile = "../../certs/server.crt"
		keyFile = "../../certs/server.key"
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки пары сертификатов: %w", err)
	}

	if isServer {
		return &tls.Config{
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert, // Требуем сертификат клиента!
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}, nil
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}
