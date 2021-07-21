package models

import "github.com/PuerkitoBio/goquery"

type ItemType string

const (
    ItemTypeCollectible ItemType = "collectible"
    ItemTypeActive ItemType = "active"
    ItemTypePassive ItemType = "passive"
    ItemTypeTrinket ItemType = "trinket"
    ItemTypePickup ItemType = "pickup"
    ItemTypeCard ItemType = "card"
    ItemTypeRune ItemType = "rune"
    ItemTypeSoulstone ItemType = "soulstone"
    ItemTypeOther ItemType = "other"
)

type TableRow struct {
    Name *goquery.Selection
    Id *goquery.Selection
    Icon *goquery.Selection
    Quote *goquery.Selection
    Description *goquery.Selection
    Quality *goquery.Selection
    Type ItemType
}

type Item struct {
    Id string
    IconUrl string
    Name string
    PageUrl string
    Charges *string
    Quote string
    Description string
    Type ItemType
    Quality *int
}
