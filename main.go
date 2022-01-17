package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "zulip plugin"
	app.Usage = "zulip plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "url",
			Usage:  "zulip webhook url",
			EnvVar: "PLUGIN_URL",
		},
		cli.StringFlag{
			Name:   "type",
			Usage:  "type of zulip message to be sent, either \"private\" or \"stream\"",
			EnvVar: "PLUGIN_TYPE",
		},
		cli.StringFlag{
			Name:   "to",
			Usage:  "zulip recipient, either stream name or user id",
			EnvVar: "PLUGIN_RECIPIENT",
		},
		cli.StringFlag{
			Name:   "topic",
			Usage:  "zulip stream topic",
			EnvVar: "PLUGIN_TOPIC",
		},
		cli.StringFlag{
			Name:   "bot_email",
			Usage:  "zulip bot email",
			EnvVar: "PLUGIN_BOT_EMAIL",
		},
		cli.StringFlag{
			Name:   "bot_apikey",
			Usage:  "zulip bot apikey",
			EnvVar: "PLUGIN_BOT_APIKEY",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author username",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.author.avatar",
			Usage:  "git author avatar",
			EnvVar: "DRONE_COMMIT_AUTHOR_AVATAR",
		},
		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR_NAME",
		},
		cli.StringFlag{
			Name:   "commit.pull",
			Usage:  "git pull request",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.StringFlag{
			Name:   "build.deployTo",
			Usage:  "environment deployed to",
			EnvVar: "DRONE_DEPLOY_TO",
		},
		cli.StringFlag{
			Name:   "stage.status",
			Usage:  "stage status",
			EnvVar: "DRONE_STAGE_STATUS",
		},
		cli.StringFlag{
			Name:   "stage.name",
			Usage:  "stage name",
			EnvVar: "DRONE_STAGE_NAME",
		},
		cli.StringFlag{
			Name:   "stage.type",
			Usage:  "stage type",
			EnvVar: "DRONE_STAGE_TYPE",
		},
		cli.StringFlag{
			Name:   "stage.kind",
			Usage:  "stage kind",
			EnvVar: "DRONE_STAGE_KIND",
		},
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env")
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Commit:   c.String("commit.sha"),
			Ref:      c.String("commit.ref"),
			Branch:   c.String("commit.branch"),
			Author:   c.String("commit.author.email"),
			Message:  c.String("commit.message"),
			Link:     c.String("build.link"),
			DeployTo: c.String("build.deployTo"),
			Started:  c.Int64("build.started"),
			Created:  c.Int64("build.created"),
		},
		Config: Config{
			URL:       c.String("url"),
			Type:      c.String("type"),
			To:        c.String("to"),
			Topic:     c.String("topic"),
			BotEmail:  c.String("bot_email"),
			BotApikey: c.String("bot_apikey"),
		},
		Stage: Stage{
			Type:   c.String("stage.type"),
			Name:   c.String("stage.name"),
			Status: c.String("stage.status"),
			Kind:   c.String("stage.kind"),
		},
	}
	return plugin.Exec()
}
