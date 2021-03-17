package main

import (
	"code.gitea.io/sdk/gitea"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	normalclient "github.com/gzlj/gitea-demo/pkg/client"
	"os"
)

var (
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
)

func init() {
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	logger = log.With(logger, "caller", log.DefaultCaller)
}

func main() {

	giteaServer := os.Getenv("GITEA_SERVER")
	if len(giteaServer) == 0 {
		level.Error(logger).Log("msg", "exiting", "err", "no gitea server is set. Please set env \"GITEA_SERVER\" like https://gitea.192.168.208.140.nip.io.")
		return
	}

	username := os.Getenv("GITEA_USERNAME")
	password := os.Getenv("GITEA_PASSWORD")

	// 创建gitea客户端
	c := normalclient.NewGiteaClient("https://gitea.192.168.208.140.nip.io", true)
	c.SetBasicAuth(username, password)

	// 获取用户信息
	user, _, err := c.GetMyUserInfo()
	if err != nil {
		level.Error(logger).Log("msg", "failed to get user info", "err", err)
		return
	}
	level.Info(logger).Log("msg", "get use info", "user", user.ID)

	// 获取单个仓库
	dockerRepo, _, err := c.GetRepo("ansible", "docker")
	if err != nil {
		level.Error(logger).Log("msg", "failed to get repo", "err", err)
	} else {
		level.Error(logger).Log("msg", "GetRepo() get repo", "id", dockerRepo.ID, "fullname", dockerRepo.FullName, "name", dockerRepo.Name, "owner", dockerRepo.Owner)
	}

	// 列出仓库
	repos, _, err := c.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		level.Error(logger).Log("msg", "failed to list repo", "err", err)
	}
	for _, r := range repos {
		level.Info(logger).Log("msg", "get repo", "id", r.ID, "fullname", r.FullName, "name", r.Name)
	}

	// 列出org级别的hook
	hooks, _, err := c.ListOrgHooks("ansible", gitea.ListHooksOptions{})
	if err != nil {
		level.Error(logger).Log("msg", "failed to list org hooks", "err", err)
	}
	for _, h := range hooks {
		level.Info(logger).Log("msg", "get org hook", "ID", h.ID, "hook_url", h.URL, "config", h.Config, "type", h.Type, "active", h.Active,
			"event", h.Events)
	}

	// 创建org级别的hook
	_, _, err = c.CreateOrgHook("ansible", gitea.CreateHookOption{
		Type:   "gitea",
		Events: []string{"push"},
		Config: map[string]string{
			"content_type": "json",
			"url":          "http://newomain.com/webhook.php",
		},
		Active: true,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to create org hook", "err", err)
	}

}
