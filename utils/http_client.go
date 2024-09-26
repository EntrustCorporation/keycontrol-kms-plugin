package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"kms-plugin/models"
	"log"
	"net/http"
	"os"
	"time"
)

type KeyControlKmsHttpClient struct {
	client *http.Client
	config models.Config
}

func NewKeyControlKmsHttpClient(config *models.Config) (*KeyControlKmsHttpClient, error) {
	caCert, err := os.ReadFile(config.CaCertFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(config.CertFile, config.CertFile)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 20,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Renegotiation: tls.RenegotiateOnceAsClient,
				RootCAs:       caCertPool,
				Certificates:  []tls.Certificate{cert},
			},
		},
	}
	return &KeyControlKmsHttpClient{
		client: client,
		config: *config,
	}, nil
}

func (c *KeyControlKmsHttpClient) doRequest(url string, jsonData map[string]string) ([]byte, error) {

	jsonVal, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonVal))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, errors.New("request failed")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *KeyControlKmsHttpClient) Encrypt(plaintext string) (string, error) {

	jsonData := map[string]string{
		"plain_text": plaintext,
		"keyid_name": c.config.KeyId,
	}
	encryptUrl := "https://" + c.config.KmsServer + "/api/1.0/symm_keyid/op/encrypt/"
	res, err := c.doRequest(encryptUrl, jsonData)
	if err != nil {
		return "", err
	}
	var data models.EncResponse
	err = json.Unmarshal(res, &data)
	if err != nil {
		return "", err
	}

	return data.CipherText, nil
}

func (c *KeyControlKmsHttpClient) Decrypt(ciphertext string) (string, error) {

	jsonData := map[string]string{
		"cipher_text": ciphertext,
		"keyid_name":  c.config.KeyId,
	}
	decryptUrl := "https://" + c.config.KmsServer + "/api/1.0/symm_keyid/op/decrypt/"
	res, err := c.doRequest(decryptUrl, jsonData)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var data models.DecResponse
	err = json.Unmarshal(res, &data)
	if err != nil {
		return "", err
	}

	return data.PlainText, nil
}
