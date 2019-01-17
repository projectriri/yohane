package main

import (
	"encoding/json"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"strings"
)

func (p *CorePlugin) produceCommand(req *ubm_api.UBM) *cmd.Command {
	if req.Type != "message" || req.Message == nil {
		return nil
	}
	if req.Message.Type != "rich_text" || req.Message.RichText == nil {
		return nil
	}

	// deep copy richTexts
	var richTexts ubm_api.RichText
	t, _ := json.Marshal(*req.Message.RichText)
	json.Unmarshal(t, &richTexts)
	// Trim all leading white characters
	for i := 0; i < len(richTexts) && richTexts[i].Type == "text"; i++ {
		richTexts[i].Text = strings.TrimLeftFunc(richTexts[i].Text, p.isWhiteChar)
		if len(richTexts[i].Text) == 0 {
			richTexts = richTexts[1:]
			i--
		} else {
			break
		}
	}
	if len(richTexts) == 0 {
		return nil
	}

	pfx := ""
	if richTexts[0].Type == "text" {
		// If the first rich text element is text, trim the command prefix
		var ok bool
		pfx, ok = p.checkPrefix(richTexts[0].Text)
		if !ok {
			return nil
		} else {
			if richTexts[0].Text == pfx {
				richTexts = richTexts[1:]
			} else {
				richTexts[0].Text = richTexts[0].Text[len(pfx):]
			}
		}
	} else {
		// else check allowEmptyPrefix
		if !p.allowEmptyPrefix {
			return nil
		}
	}

	// Process this command
	parsedCommand := make([][]ubm_api.RichTextElement, 1)
	parsedCommand[0] = make([]ubm_api.RichTextElement, 0)
	// Process rich text array
	lastEscape := false
	lastWhiteChar := false
	inQuote := false
	var lastQuoteChar rune
	buffer := make([]rune, 0)
	nowP := 0
	for _, elem := range richTexts {
		if elem.Type == "text" {
			// Text needs to be parsed
			for _, r := range elem.Text {
				// state operations
				if lastEscape {
					lastEscape = false
					buffer = append(buffer, r)
					continue
				}
				if r == ESCAPE_CHAR {
					lastEscape = true
					continue
				}
				if lastWhiteChar {
					if p.isWhiteChar(r) {
						continue
					} else {
						lastWhiteChar = false
					}
				}
				if inQuote {
					if r == lastQuoteChar {
						// end of quote
						inQuote = false
					} else {
						buffer = append(buffer, r)
					}
					continue
				}
				// state transfer
				if p.isWhiteChar(r) {
					lastWhiteChar = true
					// append and clear buffer
					if len(buffer) > 0 {
						parsedCommand[nowP] = append(parsedCommand[nowP],
							ubm_api.RichTextElement{
								Type: "text",
								Text: string(buffer),
							})
						buffer = make([]rune, 0)
					}
					// append parsedCommand
					parsedCommand = append(parsedCommand, make([]ubm_api.RichTextElement, 0))
					nowP++
				} else if p.isQuoteChar(r) {
					inQuote = true
					lastQuoteChar = r
				} else {
					// normal char
					buffer = append(buffer, r)
				}
			}
		} else {
			// Other type of message, append buffer
			if len(buffer) > 0 {
				parsedCommand[nowP] = append(parsedCommand[nowP],
					ubm_api.RichTextElement{
						Type: "text",
						Text: string(buffer),
					})
				buffer = make([]rune, 0)
			}
			// and append cur elem to the end
			parsedCommand[nowP] = append(parsedCommand[nowP], elem)
			// and clear some states
			lastEscape = false
			lastWhiteChar = false
		}
	}
	// Append buffer in the end
	if len(buffer) > 0 {
		parsedCommand[nowP] = append(parsedCommand[nowP],
			ubm_api.RichTextElement{
				Type: "text",
				Text: string(buffer),
			})
		buffer = make([]rune, 0)
	}

	// compose response according to config.ResponseMode in bit mask
	c := cmd.Command{}
	c.CmdPrefix = pfx
	if p.config.ResponseMode&RESPONSE_CMD != 0 {
		c.Cmd = parsedCommand[0]
	}
	if p.config.ResponseMode&RESPONSE_CMDSTR != 0 {
		for _, elem := range parsedCommand[0] {
			if elem.Type == "text" {
				c.CmdStr += elem.Text
			}
		}
	}
	if p.config.ResponseMode&RESPONSE_ARGS != 0 {
		c.Args = parsedCommand[1:]
	}
	if p.config.ResponseMode&RESPONSE_ARGSTXT != 0 ||
		p.config.ResponseMode&RESPONSE_ARGSSTR != 0 {
		tmpArgsTxt := make([]string, 0)
		for _, aCmd := range parsedCommand[1:] {
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
		if p.config.ResponseMode&RESPONSE_ARGSTXT != 0 {
			c.ArgsTxt = tmpArgsTxt
		}
		if p.config.ResponseMode&RESPONSE_ARGSSTR != 0 {
			c.ArgsStr = strings.Join(tmpArgsTxt, " ")
		}
	}
	c.Message = req.Message

	return &c
}
