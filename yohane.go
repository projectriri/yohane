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
	pc               *router.ProducerChannel
	cc               *router.ConsumerChannel
}

func (p *CorePlugin) GetManifest() types.Manifest {
	return Manifest
}

func Init(filename string, configPath string) []types.Converter {
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
	return []types.Converter{&p}
}

func (p *CorePlugin) Start() {
	log.Infof("[yohane] registering consumer channel %v", p.config.ChannelUUID)
	p.cc = router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   ".*",
			Formats: []types.Format{
				{
					API:     "ubm-api",
					Version: "1.0",
					Method:  "receive",
				},
				{
					API:      "cmd",
					Version:  "1.0",
					Method:   "cmd",
					Protocol: `{"command_prefix":["/"],"response_mode":6}`,
				},
			},
		},
	})
	defer p.cc.Close()
	log.Infof("[yohane] registered consumer channel %v", p.cc.UUID)
	log.Infof("[yohane] registering producer channel %v", p.config.ChannelUUID)
	p.pc = router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer p.pc.Close()
	log.Infof("[yohane] registered producer channel %v", p.pc.UUID)
	for {
		packet := p.cc.Consume()
		switch packet.Head.Format.API {
		case "ubm-api":
			p.produceRawCommand(packet)
		case "cmd":
			p.handleCommand(packet)
		}
	}
}

func (p *CorePlugin) replyText(message *ubm_api.Message, text string) {
	msg := ubm_api.UBM{
		Type: "message",
		Message: &ubm_api.Message{
			Type: "rich_text",
			RichText: &ubm_api.RichText{
				{
					Type: "text",
					Text: text,
				},
			},
			CID:     &message.Chat.CID,
			ReplyID: message.ID,
		},
	}
	b, _ := json.Marshal(msg)
	p.pc.Produce(types.Packet{
		Head: types.Head{
			From: "yohane",
			To:   message.Chat.CID.Messenger,
			UUID: utils.GenerateUUID(),
			Format: types.Format{
				API:      "ubm-api",
				Method:   "send",
				Version:  "1.0",
				Protocol: "",
			},
		},
		Body: b,
	})
}
