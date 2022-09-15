// Copyright © 2021 HyTrust, Inc. All Rights Reserved.
package main

import (
    "os"
    b64 "encoding/base64"
    "flag"
    "errors"
    "time"
    "net"
    "bytes"
    "os/signal"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "net/http"
    "log"
    "syscall"
    "encoding/json"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
    "golang.org/x/net/trace"
    k8spb "k8s.io/apiserver/pkg/storage/value/encrypt/envelope/v1beta1"
//    k8spb "src/v1beta1"
)

const (
    socketPath     = "/var/tmp/run.sock"
    netProtocol    = "unix"
    apiVersion     = "v1beta1"
    runtimeName    = "HyTrustKMS"
    runtimeVersion = "1.0"
)

// Make the first letter of struct elements upper case to export them.
// Export each json with a lowercase letter
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
    KeyIdName  string `json:"keyid_name"`
    KeyId      string `json:"keyid"`
    Result     string `json:"result"`
    PlainText  string `json:"plain_text"`
}

type KeyManagementServiceServer struct {
    *grpc.Server
    net.Listener
    pathToUnixSocket string
    config Config
    client *http.Client
}

func New(pathToUnixSocket string, configFilePath string) (*KeyManagementServiceServer, error) {
    keyManagementServiceServer := new(KeyManagementServiceServer)
    keyManagementServiceServer.pathToUnixSocket = pathToUnixSocket

    var config Config
    err := config.getConfig(configFilePath)
    if err != nil {
        log.Fatal("Error: %v", err)
        return nil, err
    }
    keyManagementServiceServer.config = config

    caCert, err := ioutil.ReadFile(config.CaCertFile)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    cert, err := tls.LoadX509KeyPair(config.CertFile, config.CertFile)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    client := &http.Client{
            Timeout: time.Second * 20,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    Renegotiation: tls.RenegotiateOnceAsClient,
                    RootCAs: caCertPool,
                    Certificates: []tls.Certificate{cert},
                },
            },
    }
    keyManagementServiceServer.client = client
    return keyManagementServiceServer, nil
}

func doRequest(url string, jsonData map[string]string, config *Config, client *http.Client) ([]byte, error) {

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

    r, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    defer r.Body.Close()

    if 200 != r.StatusCode {
        log.Fatal("Response status: %d", r.StatusCode)
        return nil, errors.New("Request failed")
    }
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    return body, nil
}

func encrypt(config *Config, client *http.Client, plaintext string) (string, error) {

    jsonData := map[string]string{
        "plain_text": plaintext,
        "keyid_name": config.KeyId,
    }
    encryptUrl := "https://" + config.KmsServer + "/api/1.0/symm_keyid/op/encrypt/"
    res, err := doRequest(encryptUrl, jsonData, config, client)
    if err != nil {
        log.Fatal(err)
        return "", err
    }
    var data EncResponse
    err = json.Unmarshal(res, &data)
    if err != nil {
        log.Fatal(err)
        return "", err
    }

    return data.CipherText, nil
}

func decrypt(config *Config, client *http.Client, ciphertext string) (string, error) {
   
    jsonData := map[string]string {
        "cipher_text": ciphertext,
        "keyid_name": config.KeyId,
    }
    decryptUrl := "https://" + config.KmsServer + "/api/1.0/symm_keyid/op/decrypt/"
    res, err := doRequest(decryptUrl, jsonData, config, client)
    if err != nil {
        log.Fatal(err)
        return "", err
    }

    var data DecResponse
    err = json.Unmarshal(res, &data)
    if err != nil {
        log.Fatal(err)
        return "", err
    }
    
    return data.PlainText, nil
}

func (config *Config) getConfig(configFile string) (error) {
    configJson, err := os.Open(configFile)
    if err != nil {
        log.Fatal(err)
        return err
    }
    defer configJson.Close()

    jsonVal, err := ioutil.ReadAll(configJson)
    if err != nil {
        log.Fatal(err)
        return err
    }

    err = json.Unmarshal(jsonVal, &config)
    if err != nil {
        log.Fatal(err)
        return err
    }
    return nil
}

func main() {

    var (
        debugListenAddr = flag.String("debug-listen-addr", "127.0.0.1:7901", "HTTP listen address.")
    )

    sockFile := flag.String("sockFile", "", "unix Domain socket that gRpc server listens to")
    confFile := flag.String("confFile", "", "config File location for HyTrust KMS plugin")
    flag.Parse()

    if len(*sockFile) == 0 {
        log.Fatal("No sockFile specified")
    }

    if len(*confFile) == 0 {
        log.Fatal("No config file specified")
    }

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGTERM)

    hytrustKeyServer, err := New(*sockFile, *confFile)
    if err != nil {
        log.Fatal("Failed to start")
    }

    err = hytrustKeyServer.deleteSocketFile()
    if err != nil {
        log.Fatal("Failed to clean socketFile")
    }
    listener, err := net.Listen(netProtocol, hytrustKeyServer.pathToUnixSocket)
    if err != nil {
        log.Fatal("Failed to start listner: %c", err)
    }
    hytrustKeyServer.Listener = listener

    server := grpc.NewServer()
    k8spb.RegisterKeyManagementServiceServer(server, hytrustKeyServer)
    hytrustKeyServer.Server = server

    go server.Serve(listener)
    trace.AuthRequest = func(req *http.Request) (any, sensitive bool) { return true, true }
    log.Println("KeyManagementServiceServer service started successfully.")

    go func() {
        for {
            sig := <-sigChan
            if sig == syscall.SIGTERM {
                log.Println("force stop")
                log.Println("Shutting down grpc service")
                server.GracefulStop()
                os.Exit(0)
            }
        }
    }()
    log.Fatal(http.ListenAndServe(*debugListenAddr, nil))
}

func (kmsServer *KeyManagementServiceServer) Version(ctx context.Context, request *k8spb.VersionRequest) (*k8spb.VersionResponse, error) {
    return &k8spb.VersionResponse{Version: apiVersion, RuntimeName: runtimeName, RuntimeVersion: runtimeVersion}, nil
}

func (kmsServer *KeyManagementServiceServer) Decrypt(ctx context.Context, request *k8spb.DecryptRequest) (*k8spb.DecryptResponse, error) {

    plainText, err := decrypt(&kmsServer.config, kmsServer.client, string(request.Cipher))
    if err != nil {
        log.Fatal("Failed to decrypt data. Error: %v", err)
        return &k8spb.DecryptResponse{}, err
    }
    sDec, _ := b64.StdEncoding.DecodeString(plainText)
    return &k8spb.DecryptResponse{Plain: []byte(sDec)}, nil
}
func (kmsServer *KeyManagementServiceServer) Encrypt(ctx context.Context, request *k8spb.EncryptRequest) (*k8spb.EncryptResponse, error) {

    sEnc := b64.StdEncoding.EncodeToString([]byte(request.Plain))
    cipherText, err := encrypt(&kmsServer.config, kmsServer.client, sEnc)
    if err != nil {
        log.Fatal("Failed to encrypt data. Error: %v", err)
        return &k8spb.EncryptResponse{}, err
    }
   return &k8spb.EncryptResponse{Cipher: []byte(cipherText)}, nil
}

/*cleanSockFile function cleans the unix socker created for the gRPC server. */
func (kmsServer *KeyManagementServiceServer) deleteSocketFile() error {

    _, err := os.Stat(kmsServer.pathToUnixSocket)
    if !os.IsNotExist(err) {
        err := os.Remove(kmsServer.pathToUnixSocket)
        if err != nil {
            log.Fatal("Failed to delete socket file: %v", err)
        }
        return err
    }
    return nil
}
