package models

type Config struct {
	KmsServer  string `json:"kmsServer"`
	CertFile   string `json:"certFile"`
	CaCertFile string `json:"caCertFile"`
	KeyId      string `json:"keyId"`
}

type EncResponse struct {
	KeyIdName  string `json:"keyid_name"`
	KeyId      string `json:"keyid"`
	Result     string `json:"result"`
	CipherText string `json:"cipher_text"`
}

type DecResponse struct {
	KeyIdName string `json:"keyid_name"`
	KeyId     string `json:"keyid"`
	Result    string `json:"result"`
	PlainText string `json:"plain_text"`
}

// EncryptResponse is the response from the Envelope service when encrypting data.
type EncryptResponse struct {
	Ciphertext  []byte
	KeyID       string
	Annotations map[string][]byte
}

// DecryptRequest is the request to the Envelope service when decrypting data.
type DecryptRequest struct {
	Ciphertext  []byte
	KeyID       string
	Annotations map[string][]byte
}

// StatusResponse is the response from the Envelope service when getting the status of the service.
type StatusResponse struct {
	Version string
	Healthz string
	KeyID   string
}
