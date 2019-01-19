package main

import (
	"encoding/json"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	log "github.com/sirupsen/logrus"
)

func (p *CorePlugin) handleCommand(packet types.Packet) {
	command := cmd.Command{}
	err := json.Unmarshal(packet.Body, &command)
	if err != nil {
		log.Errorf("[yohane] command %v has an incorrect body type %v", packet.Head.UUID, err)
	}
	switch command.CmdStr {
	case "yohane::alias":
		p.handleSetAlias(command)
	case "yohane::unalias":
		p.handleRemoveAlias(command)
	}
}
