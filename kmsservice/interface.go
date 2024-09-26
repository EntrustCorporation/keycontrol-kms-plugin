package kmsservice

import (
	"context"
	"kms-plugin/models"
)

// Service allows encrypting and decrypting data using an external Key Management Service.
type IKeyControlKmsService interface {
	// Decrypt a given bytearray to obtain the original data as bytes.
	Decrypt(ctx context.Context, uid string, req *models.DecryptRequest) ([]byte, error)
	// Encrypt bytes to a ciphertext.
	Encrypt(ctx context.Context, uid string, data []byte) (*models.EncryptResponse, error)
	// Status returns the status of the KMS.
	Status(ctx context.Context) (*models.StatusResponse, error)
}
