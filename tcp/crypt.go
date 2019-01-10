package tcp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
)

// шифруем сообщение
func encrypt(plaintext []byte) ([]byte, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("Ошибка в шифровании. ", err)
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.

	// nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")
	// nonce := make([]byte, 12)
	// if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	// 	panic(err.Error())
	// }

	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error("Ошибка в шифровании. ", err)
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// возвращаем сообщение
	return ciphertext, nil
}

// расшифровываем сообщение
func decrypt(ciphertext []byte) ([]byte, error) { // Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	// ciphertext, _ := hex.DecodeString("c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471")
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("Ошибка в шифровании. ", err)
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error("Ошибка в шифровании. ", err)
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Error("Ошибка в шифровании. ", err)
		panic(err.Error())
	}

	// возвращаем сообщение
	return plaintext, nil
}

// // шифруем сообщение
// func certEncrypt() {

// 	// Verifying with a custom list of root certificates.

// 	const rootPEM = `-----BEGIN CERTIFICATE-----
// 	MIIDSjCCAjKgAwIBAgIQRK+wgNajJ7qJMDmGLvhAazANBgkqhkiG9w0BAQUFADA/
// MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMT
// DkRTVCBSb290IENBIFgzMB4XDTAwMDkzMDIxMTIxOVoXDTIxMDkzMDE0MDExNVow
// PzEkMCIGA1UEChMbRGlnaXRhbCBTaWduYXR1cmUgVHJ1c3QgQ28uMRcwFQYDVQQD
// Ew5EU1QgUm9vdCBDQSBYMzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
// AN+v6ZdQCINXtMxiZfaQguzH0yxrMMpb7NnDfcdAwRgUi+DoM3ZJKuM/IUmTrE4O
// rz5Iy2Xu/NMhD2XSKtkyj4zl93ewEnu1lcCJo6m67XMuegwGMoOifooUMM0RoOEq
// OLl5CjH9UL2AZd+3UWODyOKIYepLYYHsUmu5ouJLGiifSKOeDNoJjj4XLh7dIN9b
// xiqKqy69cK3FCxolkHRyxXtqqzTWMIn/5WgTe1QLyNau7Fqckh49ZLOMxt+/yUFw
// 7BZy1SbsOFU5Q9D8/RhcQPGX69Wam40dutolucbY38EVAjqr2m7xPi71XAicPNaD
// aeQQmxkqtilX4+U9m5/wAl0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNV
// HQ8BAf8EBAMCAQYwHQYDVR0OBBYEFMSnsaR7LHH62+FLkHX/xBVghYkQMA0GCSqG
// SIb3DQEBBQUAA4IBAQCjGiybFwBcqR7uKGY3Or+Dxz9LwwmglSBd49lZRNI+DT69
// ikugdB/OEIKcdBodfpga3csTS7MgROSR6cz8faXbauX+5v3gTt23ADq1cEmv8uXr
// AvHRAosZy5Q6XkjEGB5YGV8eAlrwDPGxrancWYaLbumR9YbK+rlmM6pZW87ipxZz
// R8srzJmwN0jP41ZL9c8PDHIyh8bwRLtTcm1D9SZImlJnt1ir/md2cXjbDaJWFBM5
// JDGFoqgCWjBH4d1QB7wCCZAA62RjYJsWvIjJEubSfZGL+T0yjWW06XyxV3bqxbYo
// Ob8VZRzI9neWagqNdwvYkQsEjgfbKbYK7p2CNTUQ
// 	-----END CERTIFICATE-----`

// 	const certPEM = `-----BEGIN CERTIFICATE-----
// 	MIIGCjCCBPKgAwIBAgISAwq1gQfBeVPOm/V4D2stco5KMA0GCSqGSIb3DQEBCwUA
// 	MEoxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MSMwIQYDVQQD
// 	ExpMZXQncyBFbmNyeXB0IEF1dGhvcml0eSBYMzAeFw0xODA1MjAxNTQ3MzlaFw0x
// 	ODA4MTgxNTQ3MzlaMBsxGTAXBgNVBAMTEGFtYXRlcmFzdS5xejAuc3UwggEiMA0G
// 	CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDAvWVRXMKjLGpjuUVTkHwfpP3N5Ac3
// 	WPLDUnU8FiZFwxXfhi0gClKdIjX3p5eja+N5yPof9j4DRFszxNAHZa8eeslkKKeP
// 	tbcfkKeD12heozWQ5O1aYW2m6SHePIgTiLz9A3LmIIIl7wIcSflVqZznFhHdF7dH
// 	QTewepzmYKidSDSNTHKtB1WG75JP5pwRFVjk23yhd8pox/B1W9P0A4R3la36ZQ1s
// 	7Uxp9DeFpR7BlvriZqBJunA4xN+A7X8o3EJ6jobqEVi06RcwwyMO9vTI56Z8F9d4
// 	tP5ojAfx0TlD/K+iwY2Nct+0cgjHzhp0nDA5/K/2T8hb5vuKVwLhz0XrAgMBAAGj
// 	ggMXMIIDEzAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsG
// 	AQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFOHQiujNPKJk258BKEjiT/JP
// 	QkLlMB8GA1UdIwQYMBaAFKhKamMEfd265tE5t6ZFZe/zqOyhMG8GCCsGAQUFBwEB
// 	BGMwYTAuBggrBgEFBQcwAYYiaHR0cDovL29jc3AuaW50LXgzLmxldHNlbmNyeXB0
// 	Lm9yZzAvBggrBgEFBQcwAoYjaHR0cDovL2NlcnQuaW50LXgzLmxldHNlbmNyeXB0
// 	Lm9yZy8wGwYDVR0RBBQwEoIQYW1hdGVyYXN1LnF6MC5zdTCB/gYDVR0gBIH2MIHz
// 	MAgGBmeBDAECATCB5gYLKwYBBAGC3xMBAQEwgdYwJgYIKwYBBQUHAgEWGmh0dHA6
// 	Ly9jcHMubGV0c2VuY3J5cHQub3JnMIGrBggrBgEFBQcCAjCBngyBm1RoaXMgQ2Vy
// 	dGlmaWNhdGUgbWF5IG9ubHkgYmUgcmVsaWVkIHVwb24gYnkgUmVseWluZyBQYXJ0
// 	aWVzIGFuZCBvbmx5IGluIGFjY29yZGFuY2Ugd2l0aCB0aGUgQ2VydGlmaWNhdGUg
// 	UG9saWN5IGZvdW5kIGF0IGh0dHBzOi8vbGV0c2VuY3J5cHQub3JnL3JlcG9zaXRv
// 	cnkvMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUA23Sv7ssp7LH+yj5xbSzluaq7
// 	NveEcYPHXZ1PN7Yfv2QAAAFjfnMPLAAABAMARjBEAiB7MgfbYyHG42JyKoeWp8p+
// 	mtHUu9Yes7y4iUEYk236sAIgVYDNiEsX0R3YlDbBOHV3W4rv3PnkHXTsfxPqc23Y
// 	oUYAdgApPFGWVMg5ZbqqUPxYB9S3b79Yeily3KTDDPTlRUf0eAAAAWN+cw8/AAAE
// 	AwBHMEUCIQC7bPi5PRatD2SoTjqiwfsMfo9UDAJeFhCB1XHd4ZWSZgIgQ93H35yy
// 	PRxbb18C0pysKNtc3ctCvtYo/9kd72BwhsgwDQYJKoZIhvcNAQELBQADggEBAHIR
// 	J7SgYnPFlZa2cK4pLmrW0IY7ujYM9M3cK76PTDac9CsBaTII+wS9Xhy0m00RZCV1
// 	6HQ2jSYgmyXZ8wFc6QS1bSYBEMhlkkEmcZFOcrXXM0DFQGN99fPfL7qUwWb3FKrl
// 	0wibBlglKGQERAZtMq+49KpR2CzyaaeRszr1vQ1IkxSHobAYGjZ/18kt6btG2MXN
// 	hJFklHBjtM5dnD8tqOTi+01sJQbu0RYxMPYaZB4NlbRyo6lmQUG/7uMjN8AccDoY
// 	RCAiZeKOxjHuDZqyzcS9Suxi/x/NJX6P0dFRIEUBbOVEF0qEeNwnN+09MQbA4CFx
// 	6SH7/ZT/W9JJXKP9lrA=
// 	-----END CERTIFICATE-----`

// 	// First, create the set of root certificates. For this example we only
// 	// have one. It's also possible to omit this in order to use the
// 	// default root set of the current operating system.
// 	roots := x509.NewCertPool()
// 	for len([]byte(rootPEM)) > 0 {
// 		var block *pem.Block
// 		block, _ = pem.Decode([]byte(rootPEM))
// 		if block == nil {
// 			break
// 		}
// 		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
// 			continue
// 		}

// 		// cert, err := ParseCertificate(block.Bytes)
// 		// if err != nil {
// 		// 	continue
// 		// }
// 		println(block.Bytes)
// 	}
// 	// println(pemCerts)
// 	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
// 	if !ok {
// 		panic("failed to parse root certificate")
// 	}

// 	block, _ := pem.Decode([]byte(certPEM))
// 	if block == nil {
// 		panic("failed to parse certificate PEM")
// 	}
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		panic("failed to parse certificate: " + err.Error())
// 	}

// 	opts := x509.VerifyOptions{
// 		DNSName: "amaterasu.qz0.su",
// 		Roots:   roots,
// 	}

// 	if chain, err := cert.Verify(opts); err != nil {
// 		panic("failed to verify certificate: " + err.Error())
// 	} else {
// 		println(chain)
// 	}

// 	// Load your secret key from a safe place and reuse it across multiple
// 	// Seal/Open calls. (Obviously don't use this example key for anything
// 	// real.) If you want to convert a passphrase to a key, use a suitable
// 	// package like bcrypt or scrypt.
// 	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
// 	// key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
// 	// plaintext := []byte(message)

// 	// block, err := aes.NewCipher(key)
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }

// 	// // Never use more than 2^32 random nonces with a given key because of the risk of a repeat.

// 	// // nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")
// 	// // nonce := make([]byte, 12)
// 	// // if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 	// // 	panic(err.Error())
// 	// // }

// 	// nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

// 	// aesgcm, err := cipher.NewGCM(block)
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }

// 	// ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
// 	// fmt.Printf("%x\n", ciphertext)
// 	// возвращаем сообщение
// }
