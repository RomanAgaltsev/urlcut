// Пакет cert предназначен для создания нового сертификата и приватного ключа TLS.
package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	// CertPEM содержит имя файла сертификата TLS.
	CertPEM string = "cert.pem"

	// PrivateKeyPEM содержит имя файла приватного ключа TLS.
	PrivateKeyPEM string = "privatekey.pem"
)

// NewCertificate создает новый сертификат и приватный ключ TLS.
func NewCertificate() error {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// сохраняем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	err = writeCertFile(CertPEM, "CERTIFICATE", certBytes)
	if err != nil {
		return err
	}

	err = writeCertFile(PrivateKeyPEM, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
	if err != nil {
		return err
	}

	return nil
}

func writeCertFile(fname string, ftype string, fbytes []byte) error {
	var b bytes.Buffer

	err := pem.Encode(&b, &pem.Block{
		Type:  ftype,
		Bytes: fbytes,
	})
	if err != nil {
		return err
	}

	if err = os.WriteFile(fname, b.Bytes(), 0o666); err != nil {
		return err
	}

	return nil
}
