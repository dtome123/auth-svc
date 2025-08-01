package models

import "time"

type Session struct {
	UserID           string    `bson:"user_id" json:"user_id"`
	DeviceID         string    `bson:"device_id" json:"device_id"`
	Type             string    `bson:"type" json:"type"`
	AccessTokenHash  string    `bson:"access_token_hash" json:"access_token_hash"`
	RefreshTokenHash string    `bson:"refresh_token_hash" json:"refresh_token_hash"`
	TTL              time.Time `bson:"ttl" json:"ttl"`
}
