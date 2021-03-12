package connect

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	host       string
	sftpClient *sftp.Client
	session    *ssh.Session
}

func (conn *Connection) Init(host, passwd string) error {
	conn.host = host
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(passwd))

	clientConfig = &ssh.ClientConfig{
		User:    "root",
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:22", conn.host)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	// create ssh session for run cmd
	if conn.session, err = sshClient.NewSession(); err != nil {
		return err
	}

	// create sftp client
	if conn.sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return err
	}
	return nil
}

func (conn *Connection) IsExist(remoteFileDir string) bool {
	_, err := conn.sftpClient.Open(remoteFileDir)
	if err != nil {
		log.Print(remoteFileDir, " 不存在")
		return false
	}
	log.Print(remoteFileDir, " 存在")
	return true
}

// download file from sftp and copy to localFileDir
// localFileDir: such as data/0/ip/errorInfo.log
func (conn *Connection) GetFile(remoteFileDir, localDataDir, fileName string) {
	os.MkdirAll(localDataDir, os.ModePerm)
	dstFile, _ := os.Create(localDataDir + "/" + fileName)
	defer dstFile.Close()

	if srcFile, err := conn.sftpClient.Open(remoteFileDir); err == nil {
		if _, err = srcFile.WriteTo(dstFile); err != nil {
			log.Print(remoteFileDir, err)
		}
		defer srcFile.Close()
	} else {
		log.Print(remoteFileDir, err)
	}
	// log.Printf("GetFile %v success. save in %v/%v", remoteFileDir, localDataDir, fileName)
}

// run cmd
func (conn *Connection) RunCmd(cmd string) {
	var stdOut, stdErr bytes.Buffer

	conn.session.Stdout = &stdOut
	conn.session.Stderr = &stdErr

	conn.session.Run(cmd)

	log.Printf("%v: %v", conn.host, stdOut)
	log.Printf("%v: %v", conn.host, stdErr)
}
