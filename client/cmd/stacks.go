package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"net"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/apps"
	"github.com/deis/deis/client/controller/api"
)

var createNewApp = func(client *client.Client, appId string) (api.App, error) {


	fmt.Print("Creating Application... ")
	quit := progress()
	app, err := apps.New(client, appId)

	quit <- true
	<-quit

	if err != nil {
		return api.App{}, err
	}

	fmt.Printf("done, created %s\n", app.ID)
	return app, nil
}

func StackCreate(appId, stackName string) error {
	c, err := client.New()
	app, err := createNewApp(c, appId)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
	}

	if err != nil {
		return  fmt.Errorf("Fail to get client config: %s", err)
	}

	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.ControllerURL.Host , 2222), sshConfig)
	if err != nil {
		return  fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	err = session.Run(fmt.Sprintf("stack-init %s %s", app.ID, stackName))
	return nil
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
