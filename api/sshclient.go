package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

func sshRun(user string, host string, port string, privateKeyPath string, privateKeyPassword string, cmd string) {
	log.Debugf("ssh %v@%v:%v", user, host, port)

	// read private key file
	pemBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Errorf("Read private key failed %v", err)
		return
	}

	// read pem block
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		log.Errorf("Pem decode failed, no key found in %s", privateKeyPath)
		return
	}

	if x509.IsEncryptedPEMBlock(pemBlock) {
		pemBlock.Bytes, err = x509.DecryptPEMBlock(pemBlock, []byte(privateKeyPassword))
		if err != nil {
			log.Errorf("Decrypting private PEM data failed %v", err)
			return
		}

		key, err := parsePemBlock(pemBlock)
		if err != nil {
			log.Errorf("Parse PEM block failed %v", err)
			return
		}

		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			log.Errorf("New Signer from key failed %v", err)
			return
		}

		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}

		server := host + ":" + port

		conn, err := ssh.Dial("tcp", server, config)
		if err != nil {
			log.Errorf("Failed to dial %v", err)
			return
		}
		defer conn.Close()

		session, err := conn.NewSession()
		if err != nil {
			log.Errorf("Failed to create session %v", err)
			return
		}
		defer session.Close()

		var stdoutBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Run(cmd)

		fmt.Println(stdoutBuf.String())
	}
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "DSA PRIVATE KEY":
		return ssh.ParseDSAPrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("Unsupported key type %q", block.Type)
	}
}
