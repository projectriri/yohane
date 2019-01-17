package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

var manifest = types.Manifest{
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
	return manifest
}

func (p *CorePlugin) Init(filename string, configPath string) {
	// load toml config
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &p.config)
	if err != nil {
		panic(err)
	}
	p.allowEmptyPrefix = p.checkAllowEmptyPrefix()
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
		req := ubm_api.UBM{}
		err := json.Unmarshal(packet.Body, &req)
		if err != nil {
			log.Errorf("[yohane] message %v has an incorrect body type %v", packet.Head.UUID, err)
		}
		c := p.produceCommand(&req)
		if c != nil {
			b, _ := json.Marshal(c)
			pc.Produce(types.Packet{
				Head: types.Head{
					From: packet.Head.From,
					To:   packet.Head.To,
					UUID: utils.GenerateUUID(),
					Format: types.Format{
						API:      "cmd",
						Method:   "cmd",
						Version:  "1.0",
						Protocol: "",
					},
				},
				Body: b,
			})
		}
	}
}

var PluginInstance types.Adapter = &CorePlugin{}
