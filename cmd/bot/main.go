package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"reflect"
	"strconv"

	"github.com/Allan-Nava/Nomad-Deploy-Notifier/internal/bot"
	"github.com/Allan-Nava/Nomad-Deploy-Notifier/internal/stream"
)
//
//
func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	ctx, closer := CtxWithInterrupt(context.Background())
	defer closer()
    //
	token := os.Getenv("SLACK_TOKEN")
	toChannel := os.Getenv("SLACK_CHANNEL")
    debug := os.Getenv("DEBUG")
	log.Println("Before :", reflect.TypeOf(debug))
	debugBol,_ := strconv.ParseBool(debug)
    //
	if token == ""{
		log.Fatal("SLACK_TOKEN is empty: ", token)
	}
    if toChannel == "" {
		log.Fatal("SLACK_CHANNEL is empty: ", token)
    }
    //
	slackCfg := bot.Config{
		Token:   token,
		Channel: toChannel,
	}
    // 
	stream := stream.NewStream(debugBol)
    //
	slackBot, err := bot.NewBot(slackCfg)
	if err != nil {
		panic(err)
	}

	stream.Subscribe(ctx, slackBot)

	return 0
}

func CtxWithInterrupt(ctx context.Context) (context.Context, func()) {

	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	return ctx, func() {
		signal.Stop(ch)
		cancel()
	}
}