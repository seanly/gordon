package command

import (
	"context"
	"fmt"

	"github.com/kyoh86/ask"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/kyoh86/gordon/internal/hub"
	"github.com/morikuni/aec"
)

func q(text string) string {
	return aec.GreenF.With(aec.Bold).Apply("?") + " " + aec.Bold.Apply(text)
}

func Setup(_ context.Context, ev Env, cfg *env.Config, force bool) error {
	user := ev.GithubUser()
	if user == "" || force {
		if err := ask.Default(ev.GithubUser()).Message(q("Enter your GitHub user ID")).StringVar(&user).Do(); err != nil {
			return fmt.Errorf("asking GitHub user ID: %w", err)
		}

		opt, err := cfg.Property("github.user")
		if err != nil {
			return err
		}

		return opt.Set(user)
	}
	token, _ := hub.GetGithubToken(ev.GithubHost(), user)
	if token == "" || force {
		if err := ask.Default(token).Hidden(true).Message(q("Enter your GitHub Private Access Token")).StringVar(&token).Do(); err != nil {
			return fmt.Errorf("asking GitHub Private Access Token: %w", err)
		}
	}

	return hub.SetGithubToken(ev.GithubHost(), user, token)
}
