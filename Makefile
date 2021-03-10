#
# Copyright (c) 2021 HyTrust, Inc. All Rights Reserved.
#
THISDIR = $(shell pwd)

GOCMD = /usr/local/go/bin/go
GOBUILD = $(GOCMD) build

BUILDDIR = kmsplugin-build
WORKSPACE = $(THISDIR)/$(BUILDDIR)

KMSPLUGIN_SERVERGO = server.go
KMSPLUGIN_LINUX = kms-plugin-server
KMSPLUGIN_WINDOWS = kms-plugin-server.exe
KMSPLUGIN_MAC = kms-plugin-server.app

all:
	@echo "Compiling kmsplugin for Linux, Windows & Mac..."
	@mkdir $(WORKSPACE)
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(WORKSPACE)/$(KMSPLUGIN_LINUX) $(KMSPLUGIN_SERVERGO)
	@env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(WORKSPACE)/$(KMSPLUGIN_WINDOWS) $(KMSPLUGIN_SERVERGO)
	@env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(WORKSPACE)/$(KMSPLUGIN_MAC) $(KMSPLUGIN_SERVERGO)
	@echo "Please find respective Linux, Windows & Mac kms-plugin binaries, $(KMSPLUGIN_LINUX), $(KMSPLUGIN_WINDOWS) & $(KMSPLUGIN_MAC) at $(WORKSPACE)"

install:

clean:
	@echo "Clearing old Workspace if any at $(WORKSPACE).."
	@/usr/bin/rm -rf $(WORKSPACE)
	@echo "Clean up complete..."
