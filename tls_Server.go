package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {

	CreateServerCertAndKey("127.0.0.1", "server.crt", "server.key")

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cer}}
	config.Rand = rand.Reader

	listener, err := tls.Listen("tcp", ":8080", &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}

		log.Printf("server: conn: echo %q\n", string(buf[:n]))
		n, err = conn.Write(buf[:n])
		log.Printf("server: conn: wrote %d bytes", n)

		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Println("server: conn: closed")
}

func CreateServerCertAndKey(host string, certFileName string, keyFileName string) error {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Server Cert"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(time.Hour * 24 * 365 * 10), // 10 Years
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		IPAddresses: []net.IP{net.ParseIP(host)},
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}
	certOut, err := os.Create(certFileName)
	if err != nil {
		return err
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}
	if err := certOut.Close(); err != nil {
		return err
	}

	keyOut, err := os.OpenFile(keyFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}
	if err := keyOut.Close(); err != nil {
		return err
	}
	return nil
}

func ca() {
	// Load the private key for the default authority
	caKeyData, err := ioutil.ReadFile("ca.key")
	if err != nil {
		log.Fatal(err)
	}
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyData)
	if err != nil {
		log.Fatal(err)
	}

	// Load the server certificate
	serverCertData, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatal(err)
	}
	serverCert, err := x509.ParseCertificate(serverCertData)
	if err != nil {
		log.Fatal(err)
	}

	// Create the template for the signed certificate
	certTemplate := x509.Certificate{
		SerialNumber:          serverCert.SerialNumber,
		Subject:               serverCert.Subject,
		SignatureAlgorithm:    serverCert.SignatureAlgorithm,
		PublicKeyAlgorithm:    serverCert.PublicKeyAlgorithm,
		PublicKey:             serverCert.PublicKey,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour * 365),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		SubjectKeyId:          serverCert.SubjectKeyId,
	}

	// Sign the certificate with the default authority
	signedCert, err := x509.CreateCertificate(nil, &certTemplate, serverCert, &caKey.PublicKey, caKey)
	if err != nil {
		log.Fatal(err)
	}

	// Write the signed certificate to disk
	err = ioutil.WriteFile("signed_server.crt", signedCert, 0644)
	if err != nil {
		log.Fatal(err)
	}

}
