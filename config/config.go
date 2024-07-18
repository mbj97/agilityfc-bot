package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token                          string `json:"token"`
	RunnerRequestChannelID         string `json:"runner_request_channel_id"`
	RunnerRequestSpecificMessageID string `json:"runner_request_specific_message_id"`
	AntiPkRequestMessageID         string `json:"anti_pk_request_message_id"`
	AntiPkResponseChannelID        string `json:"anti_pk_response_channel_id"`
	CheckMarkEmoji                 string `json:"check_mark_emoji"`
	RedXEmoji                      string `json:"red_x_emoji"`
	RunnerRoleID                   string `json:"runner_role_id"`
	MemberRoleID                   string `json:"member_role_id"`
	NonMemberRoleID                string `json:"non_member_role_id"`
	DataFilePath                   string `json:"data_file_path"`
	AntiPkResponseDataFilePath     string `json:"anti_pk_response_data_file_path"`
}

func LoadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
