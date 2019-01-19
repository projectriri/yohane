package main

import (
	"encoding/json"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Filter struct {
	CmdPrefix    []string `json:"command_prefix"`
	ResponseMode uint8    `json:"response_mode"`
}

func (p *CorePlugin) IsConvertible(from types.Format, to types.Format) bool {
	if strings.ToLower(from.API) == "cmd" && strings.ToLower(to.API) == "cmd" {
		if strings.ToLower(from.Method) == "raw-cmd" && strings.ToLower(to.Method) == "cmd" &&
			utils.CheckIfVersionSatisfy(from.Version, "1.0") && utils.CheckIfVersionSatisfy("1.0", to.Version) {
			return true
		}
	}
	return false
}

func (p *CorePlugin) Convert(packet types.Packet, to types.Format) (bool, []types.Packet) {
	log.Debugf("[yohane] trying convert pkt %v", packet.Head.UUID)
	var rawCommand RawCommand
	err := json.Unmarshal(packet.Body, &rawCommand)
	if err != nil {
		return false, nil
	}
	var filter Filter
	err = json.Unmarshal([]byte(to.Protocol), &filter)
	if err != nil || filter.CmdPrefix == nil {
		return false, nil
	}
	// check whether prefix satisfy
	ok := false
	for _, v := range filter.CmdPrefix {
		if v == rawCommand.Prefix {
			ok = true
			break
		}
	}
	if ok {
		command := p.composeCommand(&rawCommand, &filter)
		command.Message = rawCommand.Message
		b, _ := json.Marshal(command)
		packet.Body = b
		packet.Head.Format = to
		return true, []types.Packet{packet}
	}
	return false, nil
}

func (p *CorePlugin) composeCommand(rawCommand *RawCommand, filter *Filter) *cmd.Command {
	// compose response according to responseMode in bit mask
	c := cmd.Command{}
	c.CmdPrefix = rawCommand.Prefix
	if filter.ResponseMode&RESPONSE_CMD != 0 {
		c.Cmd = rawCommand.ParsedCommand[0]
	}
	if filter.ResponseMode&RESPONSE_CMDSTR != 0 {
		for _, elem := range rawCommand.ParsedCommand[0] {
			if elem.Type == "text" {
				c.CmdStr += elem.Text
			}
		}
	}
	if filter.ResponseMode&RESPONSE_ARGS != 0 {
		c.Args = rawCommand.ParsedCommand[1:]
	}
	if filter.ResponseMode&RESPONSE_ARGSTXT != 0 ||
		filter.ResponseMode&RESPONSE_ARGSSTR != 0 {
		tmpArgsTxt := make([]string, 0)
		for _, aCmd := range rawCommand.ParsedCommand[1:] {
			tmp := ""
			for _, elem := range aCmd {
				if elem.Type == "text" {
					tmp += elem.Text
				}
			}
			if len(tmp) != 0 {
				tmpArgsTxt = append(tmpArgsTxt, tmp)
			}
		}
		if filter.ResponseMode&RESPONSE_ARGSTXT != 0 {
			c.ArgsTxt = tmpArgsTxt
		}
		if filter.ResponseMode&RESPONSE_ARGSSTR != 0 {
			c.ArgsStr = strings.Join(tmpArgsTxt, " ")
		}
	}

	return &c
}
