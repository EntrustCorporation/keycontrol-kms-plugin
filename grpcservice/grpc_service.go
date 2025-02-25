package grpcservice

import (
	"context"
	"kms-plugin/kmsservice"
	"kms-plugin/models"
	"kms-plugin/utils"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	kmsapi "k8s.io/kms/apis/v2"
)

// GRPCService is a grpc server that runs the kms v2 alpha1 API.
type GRPCService struct {
	addr       string
	timeout    time.Duration
	server     *grpc.Server
	kmsService kmsservice.IKeyControlKmsService
}

var _ kmsapi.KeyManagementServiceServer = (*GRPCService)(nil)

// NewGRPCService creates an instance of GRPCService.
func NewGRPCService(
	address string,
	timeout time.Duration,
	kmsService kmsservice.IKeyControlKmsService,
) *GRPCService {
	utils.Logger.Info("creating new grpc service instance")
	return &GRPCService{
		addr:       address,
		timeout:    timeout,
		kmsService: kmsService,
	}
}

// ListenAndServe accepts incoming connections on a Unix socket. It is a blocking method.
// Returns non-nil error unless Close or Shutdown is called.
func (s *GRPCService) ListenAndServe() error {
	utils.Logger.Info("attempting to serve on ", s.addr)
	_, err := os.Stat(s.addr)
	if !os.IsNotExist(err) {
		err := os.Remove(s.addr)
		if err != nil {
			utils.Logger.Error("failed to delete socket file:", err)
			return err
		}
	}
	ln, err := net.Listen("unix", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	gs := grpc.NewServer(
		grpc.ConnectionTimeout(s.timeout),
	)
	s.server = gs

	kmsapi.RegisterKeyManagementServiceServer(gs, s)
	utils.Logger.Info("litening on ", s.addr)
	return gs.Serve(ln)
}

// Shutdown performs a graceful shutdown. Doesn't accept new connections and
// blocks until all pending RPCs are finished.
func (s *GRPCService) Shutdown() {
	utils.Logger.Info("attempting to gracefully shut down grpc service")
	if s.server != nil {
		s.server.GracefulStop()
		utils.Logger.Info("stopped grpc server gracefully")
	}
}

// Close stops the server by closing all connections immediately and cancels
// all active RPCs.
func (s *GRPCService) Close() {
	if s.server != nil {
		s.server.Stop()
	}
}

// Status sends a status request to specified kms service.
func (s *GRPCService) Status(ctx context.Context, _ *kmsapi.StatusRequest) (*kmsapi.StatusResponse, error) {
	utils.Logger.Info("status request recieved")
	res, err := s.kmsService.Status(ctx)
	if err != nil {
		utils.Logger.Error("error occured on status request", err)
		return nil, err
	}
	utils.Logger.Info("status request completed", "Version", res.Version, "Healthz", res.Healthz, "KeyID", res.KeyID)
	return &kmsapi.StatusResponse{
		Version: res.Version,
		Healthz: res.Healthz,
		KeyId:   res.KeyID,
	}, nil
}

// Decrypt sends a decryption request to specified kms service.
func (s *GRPCService) Decrypt(ctx context.Context, req *kmsapi.DecryptRequest) (*kmsapi.DecryptResponse, error) {
	utils.Logger.Info("decrypt request recieved")
	plaintext, err := s.kmsService.Decrypt(ctx, req.Uid, &models.DecryptRequest{
		Ciphertext: req.Ciphertext,
		KeyID:      req.KeyId,
	})
	if err != nil {
		utils.Logger.Error("error occured on decrypt request", err)
		return nil, err
	}

	return &kmsapi.DecryptResponse{
		Plaintext: plaintext,
	}, nil
}

// Encrypt sends an encryption request to specified kms service.
func (s *GRPCService) Encrypt(ctx context.Context, req *kmsapi.EncryptRequest) (*kmsapi.EncryptResponse, error) {
	utils.Logger.Info("encrypt request recieved")
	encRes, err := s.kmsService.Encrypt(ctx, req.Uid, req.Plaintext)
	if err != nil {
		utils.Logger.Error("error occured on encrypt request", err)
		return nil, err
	}

	return &kmsapi.EncryptResponse{
		Ciphertext: encRes.Ciphertext,
		KeyId:      encRes.KeyID,
	}, nil
}
