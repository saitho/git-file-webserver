package webserver

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	webhook "github.com/GitbookIO/go-github-webhook"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
)

func processGitHubWebhook(cfg *config.Config, gitHandler *git.GitHandler) http.Handler {
	secret := cfg.Git.Update.WebHook.GitHub.Secret
	if os.Getenv("TEST_MODE") == "1" {
		secret = ""
	}
	return webhook.Handler(secret, func(event string, payload *webhook.GitHubPayload, req *http.Request) error {
		if strings.TrimSuffix(payload.Repository.CloneURL, ".git") != strings.TrimSuffix(cfg.Git.Url, ".git") {
			return fmt.Errorf("webhook clone URL does not match configured clone URL")
		}
		if err := gitHandler.DownloadRepository(); err != nil {
			return err
		}
		return nil
	})
}

func GitHubWebHookEndpoint(cfg *config.Config, gitHandler *git.GitHandler) Handler {
	return func(resp *Response, req *Request) {
		if cfg.Git.Update.Mode != config.GitUpdateModeWebhookGitHub {
			resp.Text(http.StatusUnauthorized, "Webhook is disabled.")
			return
		}

		switch cfg.Git.Update.Mode {
		case config.GitUpdateModeWebhookGitHub:
			processGitHubWebhook(cfg, gitHandler).ServeHTTP(resp, req.Request)
			break
		default:
			resp.Text(http.StatusInternalServerError, "Unknown webhook update mode.")
			break
		}
	}
}
