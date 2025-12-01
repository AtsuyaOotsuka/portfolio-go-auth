package provider

import "gorm.io/gorm"

type Provider struct {
	db *gorm.DB
}

func NewProvider(db *gorm.DB) *Provider {
	return &Provider{
		db: db,
	}
}
