package ssh

import (
	"fmt"
	"os"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Dir    string
	client *scp.Client
	user   User
}

func NewClient(dir string, u User) *Client {
	clientConfig, _ := auth.PrivateKey(u.Name, u.RSAPath, ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(u.Client, &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		panic("Couldn't establish a connection to the remote serve")
	}
	return &Client{
		Dir:    dir,
		user:   u,
		client: &client,
	}
}

func (c Client) Upload(path string) bool {
	// Open a file
	f, _ := os.Open(path)

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)
	p := c.user.Dir + strings.Replace(path, c.Dir, "", 1)
	fmt.Println(fmt.Sprintf("%s => %s", path, p))
	err := c.client.CopyFile(f, p, c.user.Chmod)

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
	return true
}

func (c Client) Created(dir string) bool {
	return true
}

func (c Client) Close() {
	// Close client connection after the file has been copied
	c.client.Close()
}
