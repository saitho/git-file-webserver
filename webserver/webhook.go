package webserver

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	webhook "github.com/GitbookIO/go-github-webhook"
	log "github.com/sirupsen/logrus"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
)

func processGitHubWebhook(cfg *config.Config, client *git.Client) http.Handler {
	secret := cfg.Git.Update.WebHook.GitHub.Secret
	if os.Getenv("TEST_MODE") == "1" {
		log.Warningf("Test mode for hooks is enabled. Hook calls will not be verified!")
		secret = ""
	}
	return webhook.Handler(secret, func(event string, payload *webhook.GitHubPayload, req *http.Request) error {
		if strings.TrimSuffix(payload.Repository.CloneURL, ".git") != strings.TrimSuffix(cfg.Git.Url, ".git") {
			log.Errorf("webhook clone URL does not match configured clone URL. Payload: %v", payload)
			return fmt.Errorf("webhook clone URL does not match configured clone URL")
		}
		log.Debugf("Downloading repository due to Webhook trigger.")
		if err := client.DownloadRepository(); err != nil {
			return err
		}
		return nil
	})
}

func GitHubWebHookEndpoint(cfg *config.Config, client *git.Client) Handler {
	return func(resp *Response, req *Request) {
		if cfg.Git.Update.Mode != config.GitUpdateModeWebhookGitHub {
			log.Errorf("Webhook called but webhook feature is disabled.")
			resp.Text(http.StatusUnauthorized, "Webhook is disabled.")
			return
		}

		log.Debugf("Webhook request received. Header: %v", req.Request.Header)
		switch cfg.Git.Update.Mode {
		case config.GitUpdateModeWebhookGitHub:
			processGitHubWebhook(cfg, client).ServeHTTP(resp, req.Request)
		default:
			log.Errorf("Unknown webhook update mode %s.", cfg.Git.Update.Mode)
			resp.Text(http.StatusInternalServerError, "Unknown webhook update mode.")
		}
	}
}
