package main

import (
	"encoding/json"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/types/ubm-api"
)

var cmdAliasMap = make(map[string]map[string][][]ubm_api.RichTextElement)

func (p *CorePlugin) setAlias(alias []ubm_api.RichTextElement, target [][]ubm_api.RichTextElement, aliasMap map[string][][]ubm_api.RichTextElement) {
	data, _ := json.Marshal(alias)
	aliasMap[string(data)] = target
}

func (p *CorePlugin) getAlias(alias []ubm_api.RichTextElement, aliasMap map[string][][]ubm_api.RichTextElement) [][]ubm_api.RichTextElement {
	data, _ := json.Marshal(alias)
	target, ok := aliasMap[string(data)]
	if ok {
		return target
	}
	return nil
}

func (p *CorePlugin) removeAlias(alias []ubm_api.RichTextElement, aliasMap map[string][][]ubm_api.RichTextElement) bool {
	data, _ := json.Marshal(alias)
	key := string(data)
	_, ok := aliasMap[key]
	if ok {
		delete(aliasMap, key)
		return true
	}
	return false
}

func (p *CorePlugin) handleAlias(rawCommand RawCommand) RawCommand {
	// find in local alias map
	key := rawCommand.Message.Chat.CID.String()
	if local, ok := cmdAliasMap[key]; ok {
		escaped := p.getAlias(rawCommand.ParsedCommand[0], local)
		if escaped != nil {
			escaped = append(escaped, rawCommand.ParsedCommand[1:]...)
			rawCommand.ParsedCommand = escaped
			return rawCommand
		}
	}
	// TODO: find in local regex map
	// find in global alias map
	key = "_"
	if local, ok := cmdAliasMap[key]; ok {
		escaped := p.getAlias(rawCommand.ParsedCommand[0], local)
		if escaped != nil {
			escaped = append(escaped, rawCommand.ParsedCommand[1:]...)
			rawCommand.ParsedCommand = escaped
			return rawCommand
		}
	}
	// TODO: find in global regex map
	return rawCommand
}

func (p *CorePlugin) handleSetAlias(command cmd.Command) {
	// regMode := false
	globalMode := false
	// var reg *regexp.Regexp
OPTIONS:
	for {
		if len(command.Args) < 2 {
			p.replyText(command.Message, "参数不足。")
			return
		}
		if command.Args[0][0].Type == "text" {
			switch command.Args[0][0].Text {
			//case "-e":
			//	regMode = true
			//	command.Args = command.Args[1:]
			case "--global":
				fallthrough
			case "-g":
				globalMode = true
				command.Args = command.Args[1:]
			case "--":
				command.Args = command.Args[1:]
				break OPTIONS
			default:
				break OPTIONS
			}
		} else {
			break
		}
	}
	if len(command.Args) < 2 {
		p.replyText(command.Message, "参数不足。")
		return
	}
	key := command.Message.Chat.CID.String()
	if globalMode {
		key = "_"
	}
	if _, ok := cmdAliasMap[key]; !ok {
		cmdAliasMap[key] = make(map[string][][]ubm_api.RichTextElement)
	}
	p.setAlias(command.Args[0], command.Args[1:], cmdAliasMap[key])
	p.replyText(command.Message, "别名已设置。")
}

func (p *CorePlugin) handleRemoveAlias(command cmd.Command) {
	// regMode := false
	globalMode := false
	// var reg *regexp.Regexp
OPTIONS:
	for {
		if len(command.Args) < 1 {
			p.replyText(command.Message, "参数不足。")
			return
		}
		if command.Args[0][0].Type == "text" {
			switch command.Args[0][0].Text {
			//case "-e":
			//	regMode = true
			//	if command.Args[1][0].Type != "text" {
			//		p.replyText(command.Message, "参数有误。")
			//		return
			//	}
			//	reg = regexp.MustCompile(command.Args[1][0].Text)
			//	command.Args = command.Args[2:]
			case "--global":
				fallthrough
			case "-g":
				globalMode = true
				command.Args = command.Args[1:]
			case "--":
				command.Args = command.Args[1:]
				break OPTIONS
			default:
				break OPTIONS
			}
		} else {
			break
		}
	}
	if len(command.Args) < 1 {
		p.replyText(command.Message, "参数不足。")
		return
	}
	key := command.Message.Chat.CID.String()
	if globalMode {
		key = "_"
	}
	if _, ok := cmdAliasMap[key]; !ok {
		p.replyText(command.Message, "没有这个别名呢。")
		return
	}
	if !p.removeAlias(command.Args[0], cmdAliasMap[key]) {
		p.replyText(command.Message, "没有这个别名呢。")
		return
	}
	p.replyText(command.Message, "别名已取消。")
}
