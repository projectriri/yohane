package main

import (
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
)

var cmdAliasMap = make(map[string]map[string][][]ubm_api.RichTextElement)

func (p *CorePlugin) setAlias(alias []ubm_api.RichTextElement, target [][]ubm_api.RichTextElement, aliasMap map[string][][]ubm_api.RichTextElement)  {
	data, _ := json.Marshal(alias)
	hasher := md5.New()
	hasher.Write(data)
	aliasMap[hex.EncodeToString(hasher.Sum(nil))] = target
}

func (p *CorePlugin) getAlias(alias []ubm_api.RichTextElement, aliasMap map[string][][]ubm_api.RichTextElement) [][]ubm_api.RichTextElement {
	data, _ := json.Marshal(alias)
	hasher := md5.New()
	hasher.Write(data)
	target, ok := aliasMap[hex.EncodeToString(hasher.Sum(nil))]
	if ok {
		return target
	}
	return nil
}

func (p *CorePlugin) handleSetAlias() {

}

func (p *CorePlugin) loadAliasMap() {

}
