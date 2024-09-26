package kmsservice

import (
	"context"
	"encoding/base64"
	"kms-plugin/models"
	"kms-plugin/utils"
)

type KeyControlKmsService struct {
	IKeyControlKmsService
	config *models.Config
	client *utils.KeyControlKmsHttpClient
}

func NewKeyControlKmsService(config *models.Config) (IKeyControlKmsService, error) {
	utils.Logger.Info("initalizing KeyControlKmsService")
	client, err := utils.NewKeyControlKmsHttpClient(config)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	utils.Logger.Info("successfully initialized KeyControlKmsService")
	return &KeyControlKmsService{
		config: config,
		client: client,
	}, nil
}

func (s *KeyControlKmsService) Status(ctx context.Context) (*models.StatusResponse, error) {
	return &models.StatusResponse{Version: "v2", Healthz: "ok", KeyID: s.config.KeyId}, nil
}

func (s *KeyControlKmsService) Encrypt(ctx context.Context, uid string, data []byte) (*models.EncryptResponse, error) {
	sEnc := base64.StdEncoding.EncodeToString(data)
	res, err := s.client.Encrypt(sEnc)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return &models.EncryptResponse{
		Ciphertext: []byte(res),
		KeyID:      s.config.KeyId,
	}, nil
}

func (s *KeyControlKmsService) Decrypt(ctx context.Context, uid string, req *models.DecryptRequest) ([]byte, error) {
	sDec := base64.StdEncoding.EncodeToString(req.Ciphertext)
	res, err := s.client.Decrypt(sDec)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return []byte(res), nil
}
