package ssh

import (
	"context"
	"fmt"
	"judo/internal/config"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	cfg    config.Config
	client *ssh.Client
}

func NewSSHClient(cfg config.Config) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: cfg.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.SSH.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", cfg.SSH.Host+":"+cfg.SSH.Port, config)
	if err != nil {
		return nil, err
	}

	fmt.Println("SSH соединение установлено")

	return &SSHClient{
		client: client,
		cfg:    cfg,
	}, nil
}

func (c *SSHClient) ConnectRemoteDB(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := c.client.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *SSHClient) Close() {
	c.client.Close()
}
