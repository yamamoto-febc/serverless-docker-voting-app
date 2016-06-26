package main

import (
	"fmt"
	"os"
	"net/http"

	"github.com/bfirsh/go-dcgi"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types/container"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	certPath := os.Getenv("DOCKER_CERT_PATH")
	hostConfig := &container.HostConfig{
		NetworkMode: "serverlessdockervotingapp_default",
		Binds: []string{fmt.Sprintf("%s:%s", certPath, certPath)},
	}
	inheritEnv := []string{"DOCKER_HOST", "DOCKER_MACHINE_NAME", "DOCKER_TLS_VERIFY", "DOCKER_CERT_PATH"}

	http.Handle("/vote/", &dcgi.Handler{
		Image:      "bfirsh/serverless-vote",
		Client:     cli,
		HostConfig: hostConfig,
		InheritEnv: inheritEnv,
		Root:       "/vote", // strip /vote from all URLs
	})
	http.Handle("/result/", &dcgi.Handler{
		Image:      "bfirsh/serverless-result",
		Client:     cli,
		HostConfig: hostConfig,
		InheritEnv: inheritEnv,
		Root:       "/result",
	})
	http.ListenAndServe(":80", nil)
}
