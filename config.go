package main

type Config struct {
	ChannelUUID      string   `toml:"channel_uuid"`
	CommandPrefix    []string `toml:"command_prefix"`
	CommandAliasPath string   `toml:"command_alias_path"`
}
