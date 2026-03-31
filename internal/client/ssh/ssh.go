package ssh

import (
	"judo/internal/config"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

func NewSSHClient(cfg config.SSHConf) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", cfg.Host+":"+cfg.Port, config)
	if err != nil {
		return nil, err
	}

	return &SSHClient{
		client: client,
	}, nil
}

func (s *SSHClient) MigrateOnServer() {

}

func (s *SSHClient) Close() {
	s.client.Close()
}
