package models

import "errors"

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrAdvertiserNotFound = errors.New("advertiser not found")
	ErrInvalidStartDate   = errors.New("invalid start date")
	ErrCampaignNotFound   = errors.New("campaign not found")
	ErrCantUpdateCampaign = errors.New("can`t update campaign")
	ErrNoAdsForClient     = errors.New("no ads for client")
	ErrAlreadyImpressed   = errors.New("already impressed")
	ErrAlreadyClicked     = errors.New("already clicked")
	ErrNotImpressed       = errors.New("not impressed")
	ErrStaticNotFound     = errors.New("static not found")
)
