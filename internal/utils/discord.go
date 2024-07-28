package utils

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func CheckAccountAge(userID string, minDays int) bool {
	createdAt, err := discordgo.SnowflakeTimestamp(userID)
	if err != nil {
		return false
	}
	accountAge := time.Since(createdAt).Hours() / 24
	return accountAge >= float64(minDays)
}