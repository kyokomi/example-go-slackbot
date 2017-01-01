package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	godocomo "github.com/kyokomi/go-docomo/docomo"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/akari"
	"github.com/kyokomi/slackbot/plugins/cron"
	"github.com/kyokomi/slackbot/plugins/docomo"
	"github.com/kyokomi/slackbot/plugins/kohaimage"
	"github.com/kyokomi/slackbot/plugins/lgtm"
	"github.com/kyokomi/slackbot/plugins/naruhodo"
	"github.com/kyokomi/slackbot/plugins/suddendeath"
	"github.com/kyokomi/slackbot/plugins/sysstd"
)

func init() {
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
}

func setupSlackBot(redisToGoURL, slackToken, docomoAPIKey string) error {
	// setup repository
	addr, password, err := parseRedisURL(redisToGoURL)
	if err != nil {
		return err
	}
	cronRepository := cron.NewRedisRepository(addr, password, 1)
	docomoRepository := slackbot.NewRedisRepository(addr, password, 1)

	cronCtx := cron.NewCronContext(cronRepository)
	defer cronCtx.Close()

	botCtx, err := slackbot.NewBotContextNotSysstd(slackToken)
	if err != nil {
		return err
	}
	cronCtx.AllRefreshCron(botCtx)

	sysPlugin := sysstd.NewPlugin(botCtx.PluginManager())
	sysPlugin.SetTimezone("JST")
	botCtx.AddPlugin("sysstd", sysPlugin)
	botCtx.AddPlugin("cron", cron.NewPlugin(cronCtx))
	if docomoAPIKey != "" {
		botCtx.AddPlugin("docomo", docomo.NewPlugin(godocomo.NewClient(docomoAPIKey), docomoRepository))
	}
	botCtx.AddPlugin("akari", akari.NewPlugin())
	botCtx.AddPlugin("naruhodo", naruhodo.NewPlugin())
	botCtx.AddPlugin("lgtm", lgtm.NewPlugin())
	botCtx.AddPlugin("koha", kohaimage.NewPlugin(kohaimage.NewKohaAPI()))
	botCtx.AddPlugin("suddendeath", suddendeath.NewPlugin())

	botCtx.WebSocketRTM()
	return nil
}

func main() {
	var serverAddr, slackToken, dApiKey string
	flag.StringVar(&serverAddr, "addr", ":8080", "serverのaddr")
	flag.StringVar(&slackToken, "slackToken", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.StringVar(&dApiKey, "docomo", os.Getenv("DOCOMO_APIKEY"), "DocomoのAPIKEY")
	flag.Parse()

	if err := setupSlackBot(os.Getenv("REDISTOGO_URL"), slackToken, dApiKey); err != nil {
		log.Fatalln(err)
	}

	// listen and serve for healthy
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("PONG"))
	})
	log.Println("Starting on", serverAddr)
	http.ListenAndServe(serverAddr, nil)
}
