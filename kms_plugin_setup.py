#!/usr/bin/python3

# Copyright (c) 2021 HyTrust, Inc. All Rights Reserved.

import sys
import base64
import json
import argparse

class KMSPluginSetup():
    """
    KMSPluginSetup class
    """
    def __init__(self):
        parser = argparse.ArgumentParser(
                    description='KMS Plugin Setup Script',
                    usage='''kms_plugin_setup.py <command> <certificate_bundle_path>
            Commands:
                show_client_cert  Display SSL certificate for KMS Plugin
                show_ca_cert      Display CA certificate to verify KeyControl
            ''')

        parser.add_argument('command', help='Subcommand to run')
        args = parser.parse_args(sys.argv[1:2])
        if not hasattr(self, args.command):
            print('Unrecognized command: ', args.command)
            parser.print_help()
            exit(1)
        getattr(self, args.command)()

    def show_client_cert(self):
        """
        Displyay client certificate from bundle
        """
        self.parse_certificate(client_cert=True)

    def show_ca_cert(self):
        """
        Displyay CA certificate from bundle
        """
        self.parse_certificate(ca_cert=True)

    def parse_certificate(self, client_cert=False, ca_cert=False):
        """
        Parse given certificate bundle to create certificate files
        """
        parser = argparse.ArgumentParser(
                        description='Certificate bundle downloaded from KeyControl')
        parser.add_argument('cert_bundle')
        args = parser.parse_args(sys.argv[2:])
        cert_bundle = args.cert_bundle
        try:
            client_cert_file = None
            ca_cert_file = None
            with open(cert_bundle, 'r') as f:
                cert_json = f.read()
            cert_json = json.loads(cert_json)
            for pemfile in cert_json['certificates'].keys():
                if pemfile != 'cacert.pem':
                    client_cert_file = pemfile
                else:
                    ca_cert_file = pemfile
            if client_cert and client_cert_file:
                client_cert_file = base64.b64decode(cert_json['certificates'][client_cert_file]).strip()
                print(client_cert_file.decode('utf-8'))
            if ca_cert and ca_cert_file:
                ca_cert_file = base64.b64decode(cert_json['certificates'][ca_cert_file]).strip()
                print(ca_cert_file.decode('utf-8'))
        except Exception as exc:
            print("Invalid certificate bundle. Error: ", exc)
            exit(1)

if __name__ == '__main__':
    KMSPluginSetup()
