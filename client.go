// Copyright © 2021 HyTrust, Inc. All Rights Reserved.

package main

import (
    "context"
    "log"
    "time"
    "google.golang.org/grpc"
    k8spb "src/v1beta1"
)

func main() {
    cc, err := grpc.Dial("unix:///etc/kubernetes/pki/test.sock", grpc.WithInsecure())
    if err != nil {
        log.Fatal("Socket creation failed: %v", err)
    }
    defer cc.Close()

    c := k8spb.NewKeyManagementServiceClient(cc)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    version := "1"
    r, err := c.Version(ctx, &k8spb.VersionRequest{Version: version})
    if err != nil {
        log.Fatal("Could not check version: %v", err)
    }
    log.Printf("Version: %s", r.GetVersion())
    log.Printf("Name: %s", r.GetRuntimeName())
    log.Printf("RuntimeVersion: %s", r.GetRuntimeVersion())

    b := []byte("teststring")
    encOut, err := c.Encrypt(ctx, &k8spb.EncryptRequest{Plain: b})
    if err != nil {
        log.Fatal("Could not encrypt: %v", err)
    }
    ctext := encOut.GetCipher()
    log.Printf("Cipher: %s", ctext)

    decOut, err := c.Decrypt(ctx, &k8spb.DecryptRequest{Cipher: ctext})
    if err != nil {
        log.Fatal("Could not encrypt: %v", err)
    }
    ptext := decOut.GetPlain()
    log.Printf("Plain: %s", ptext)

}
