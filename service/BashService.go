package service

import "golang.org/x/crypto/ssh"

func BashCommand(ip, username, password, command string) []byte {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, _ := ssh.Dial("tcp", ip+":22", config)
	defer client.Close()
	session, _ := client.NewSession()
	defer session.Close()
	output, _ := session.CombinedOutput(command)
	return output
}
