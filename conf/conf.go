package conf

import (
	"os"

	pubsub "github.com/soulgarden/rmq-pubsub"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env string `json:"env" required:"true"`
	RMQ struct {
		Host     string `json:"host"     required:"true"`
		Port     string `default:"6379"  json:"port"`
		User     string `default:"guest" json:"user"`
		Password string `default:"guest" json:"password"`

		Queues struct {
			Requests  *pubsub.Cfg `json:"requests"`
			Responses *pubsub.Cfg `json:"responses"`
		} `json:"queues"`
	} `json:"rmq"`
	ListenPort string `default:"8000"  json:"listen_port"`
	DebugMode  bool   `default:"false" json:"debug_mode"`
	EnableCors bool   `default:"false" json:"enable_cors"`
}

func New() *Config {
	c := &Config{}
	path := os.Getenv("CFG_PATH")

	if path == "" {
		path = "./conf/config.json"
	}

	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(c, path); err != nil {
		log.Fatal().Err(err).Msg("config validation errors")
	}

	return c
}
