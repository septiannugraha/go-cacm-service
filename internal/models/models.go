package models

import "time"

// CacheData represents the main financial data structure
type CacheData struct {
	ID              int
	Tahun           int32
	KodeDesa        string
	KodeKegiatan    string
	KodePaket       string
	KodeRekening    string
	KodeSumber      string
	Tagging         string
	Anggaran1       float64
	Anggaran2       float64
	Real1           float64
	Real2           float64
	Real3           float64
	Real4           float64
	Real5           float64
	Real6           float64
	Real7           float64
	Real8           float64
	Real9           float64
	Real10          float64
	Real11          float64
	Real12          float64
	TotalReal       float64
	KodePemda       string
	NamaPemda       string
	NamaRekening    string
	NamaSumber      string
	NamaDesa        string
	NamaKegiatan    string
	NamaPaket       *string // nullable
	IDTipologi      *string // nullable
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// BelanjaPerBidangPerSumber represents expense by field and source
type BelanjaPerBidangPerSumber struct {
	Tahun       int32
	KodeProv    string
	NamaProv    string
	KodePemda   string
	NamaPemda   string
	KodeKec     string
	NamaKec     string
	KodeDesa    string
	NamaDesa    string
	KodePosting int32
	SumberDana  string
	AnggBid01   float64
	RealBid01   float64
	AnggBid02   float64
	RealBid02   float64
	AnggBid03   float64
	RealBid03   float64
	AnggBid04   float64
	RealBid04   float64
	AnggBid05   float64
	RealBid05   float64
	CurrentDate time.Time
}

// Add other model structs as needed...