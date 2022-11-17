package stream

import (
	"context"
	"os"
	"log"

	"github.com/Allan-Nava/Nomad-Deploy-Notifier/internal/bot"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
)

type Stream struct {
	nomad *api.Client
	L     hclog.Logger
	Debug bool
}

func NewStream(debug bool) *Stream {
	client, _ := api.NewClient(&api.Config{})
	return &Stream{
		nomad: client,
		L:     hclog.Default(),
		Debug: debug
	}
}

func (s *Stream) Subscribe(ctx context.Context, slack *bot.Bot) {
	if s.Debug {
		log.Println("Subscribe method()")
	}
	events := s.nomad.EventStream()

	topics := map[api.Topic][]string{
		api.Topic("Deployment"): {"*"},
	}

	eventCh, err := events.Stream(ctx, topics, 0, &api.QueryOptions{})
	if err != nil {
		s.L.Error("error creating event stream client", "error", err)
		os.Exit(1)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-eventCh:
			if event.Err != nil {
				s.L.Warn("error from event stream", "error", err)
				break
			}
			if event.IsHeartbeat() {
				continue
			}

			for _, e := range event.Events {
				deployment, err := e.Deployment()
				if err != nil {
					s.L.Error("execpted deployment", "error", err)
					continue
				}

				if err = slack.UpsertDeployMsg(*deployment); err != nil {
					s.L.Warn("error decoding payload", "error", err)
					return
				}
			}
		}
	}
}