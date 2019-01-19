package main

import (
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

var Manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:        "yohane",
		Author:      "Project Riri Staff",
		Version:     "v0.1",
		License:     "MIT",
		URL:         "https://github.com/projectriri/yohane",
		Description: "Core Extension of Little Daemon Bot Gateway.",
	},
	BuildInfo: types.BuildInfo{
		BuildTag:      BuildTag,
		BuildDate:     BuildDate,
		GitCommitSHA1: GitCommitSHA1,
		GitTag:        GitTag,
	},
}

type CorePlugin struct {
	config           Config
	allowEmptyPrefix bool
}

func (p *CorePlugin) GetManifest() types.Manifest {
	return Manifest
}

func Init(filename string, configPath string) []types.Adapter {
	// load toml config
	var config Config
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &config)
	if err != nil {
		panic(err)
	}
	p := CorePlugin{
		config: config,
	}
	p.allowEmptyPrefix = p.checkAllowEmptyPrefix()
	return []types.Adapter{&p}
}

func (p *CorePlugin) Start() {
	log.Infof("[yohane] registering consumer channel %v", p.config.ChannelUUID)
	cc := router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   ".*",
			Formats: []types.Format{
				{
					API:     "ubm-api",
					Version: "1.0",
					Method:  "receive",
				},
			},
		},
	})
	defer cc.Close()
	log.Infof("[yohane] registered consumer channel %v", cc.UUID)
	log.Infof("[yohane] registering producer channel %v", p.config.ChannelUUID)
	pc := router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer pc.Close()
	log.Infof("[yohane] registered producer channel %v", pc.UUID)
	for {
		packet := cc.Consume()
		p.produceRawCommand(packet, pc)
	}
}
