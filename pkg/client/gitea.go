package client

import (
	"code.gitea.io/sdk/gitea"
	"crypto/tls"
	"net/http"
)

func NewGiteaClient(url string, skipTlsVerify bool) *gitea.Client{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTlsVerify},
	}
	client := &http.Client{
		Transport: tr,
	}
	c := gitea.NewClientWithHTTP(url, client)
	return c
}
