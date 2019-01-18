package main

import (
	"encoding/json"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func (p *CorePlugin) IsConvertible(from types.Format, to types.Format) bool {
	if strings.ToLower(from.API) == "cmd" && strings.ToLower(to.API) == "cmd" {
		if strings.ToLower(from.Method) == "cmd" && strings.ToLower(to.Method) == "cmd" &&
			utils.CheckIfVersionSatisfy(from.Version, "1.0") && utils.CheckIfVersionSatisfy("1.0", to.Version) {
			return true
		}
	}
	return false
}

func (p *CorePlugin) Convert(packet types.Packet, to types.Format) (bool, []types.Packet) {
	log.Debugf("[yohane] trying convert pkt %v", packet.Head.UUID)
	var command cmd.Command
	err := json.Unmarshal(packet.Body, &command)
	if err != nil {
		return false, nil
	}
	if ok, _ := regexp.MatchString(to.Protocol, command.CmdPrefix); ok {
		return true, []types.Packet{packet}
	}
	return false, nil
}
