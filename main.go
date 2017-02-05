package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	godocomo "github.com/kyokomi/go-docomo/docomo"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/akari"
	"github.com/kyokomi/slackbot/plugins/cron/v2"
	"github.com/kyokomi/slackbot/plugins/docomo"
	"github.com/kyokomi/slackbot/plugins/esa"
	"github.com/kyokomi/slackbot/plugins/kohaimage"
	"github.com/kyokomi/slackbot/plugins/lgtm"
	"github.com/kyokomi/slackbot/plugins/naruhodo"
	"github.com/kyokomi/slackbot/plugins/router"
	"github.com/kyokomi/slackbot/plugins/suddendeath"
	"github.com/kyokomi/slackbot/plugins/sysstd"
)

type Config struct {
	redisToGoURL string
	slackToken   string
	docomoAPIKey string
	esaTeam      string
	esaToken     string
}

func init() {
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
}

func setupSlackBot(cfg Config) error {
	// setup repository
	addr, password, err := parseRedisURL(cfg.redisToGoURL)
	if err != nil {
		return err
	}

	cronRepository := cron.NewRedisRepository(addr, password, 0, "")
	docomoRepository := slackbot.NewRedisRepository(addr, password, 0)

	cronCtx, err := cron.NewContext(cronRepository)
	if err != nil {
		return err
	}

	botCtx, err := slackbot.NewBotContextNotSysstd(cfg.slackToken)
	if err != nil {
		return err
	}
	cronCtx.AllRefresh(botCtx)

	sysPlugin := sysstd.NewPlugin(botCtx.PluginManager())
	sysPlugin.SetTimezone("JST")
	botCtx.AddPlugin("sysstd", sysPlugin)
	botCtx.AddPlugin("cron", cron.NewPlugin(cronCtx))
	botCtx.AddPlugin("router", router.NewPlugin(botCtx.Client, docomoRepository))
	if cfg.docomoAPIKey != "" {
		botCtx.AddPlugin("docomo", docomo.NewPlugin(godocomo.NewClient(cfg.docomoAPIKey), docomoRepository))
	}
	botCtx.AddPlugin("akari", akari.NewPlugin())
	botCtx.AddPlugin("naruhodo", naruhodo.NewPlugin())
	botCtx.AddPlugin("lgtm", lgtm.NewPlugin())
	botCtx.AddPlugin("koha", kohaimage.NewPlugin(kohaimage.NewKohaAPI()))
	botCtx.AddPlugin("suddendeath", suddendeath.NewPlugin())
	if cfg.esaTeam != "" && cfg.esaToken != "" {
		botCtx.AddPlugin("esa", esa.NewPlugin(cfg.esaTeam, cfg.esaToken))
	}

	botCtx.WebSocketRTM()
	return nil
}

func main() {
	var serverAddr, slackToken, dApiKey, esaTeam, esaToken string
	flag.StringVar(&serverAddr, "addr", ":8080", "serverのaddr")
	flag.StringVar(&slackToken, "slackToken", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.StringVar(&dApiKey, "docomo", os.Getenv("DOCOMO_APIKEY"), "DocomoのAPIKEY")
	flag.StringVar(&esaTeam, "esa-team", os.Getenv("ESA_TEAM"), "esaのチーム名")
	flag.StringVar(&esaToken, "esa-token", os.Getenv("ESA_TOKEN"), "esaのToken")
	flag.Parse()

	cfg := Config{
		redisToGoURL: os.Getenv("REDISTOGO_URL"),
		slackToken:   slackToken,
		docomoAPIKey: dApiKey,
		esaTeam:      esaTeam,
		esaToken:     esaToken,
	}
	if err := setupSlackBot(cfg); err != nil {
		log.Fatalln(err)
	}

	// listen and serve for healthy
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("PONG"))
	})
	log.Println("Starting on", serverAddr)
	http.ListenAndServe(serverAddr, nil)
}
