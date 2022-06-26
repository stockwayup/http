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
		Host     string `json:"host"      required:"true"`
		Port     string `json:"port"      default:"6379"`
		User     string `json:"user"      default:"guest"`
		Password string `json:"password"  default:"guest"`

		Queues struct {
			Requests  *pubsub.Cfg `json:"requests"`
			Responses *pubsub.Cfg `json:"responses"`
		} `json:"queues"`
	} `json:"rmq"`
	ListenPort string `json:"listen_port" default:"8000"`
	DebugMode  bool   `json:"debug_mode"  default:"false"`
	EnableCors bool   `json:"enable_cors" default:"false"`
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
