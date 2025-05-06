package models

type OracleCardImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	PNG        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

type OracleCardLegalities struct {
	Standard        string `json:"standard"`
	Future          string `json:"future"`
	Historic        string `json:"historic"`
	Timeless        string `json:"timeless"`
	Gladiator       string `json:"gladiator"`
	Pioneer         string `json:"pioneer"`
	Explorer        string `json:"explorer"`
	Modern          string `json:"modern"`
	Legacy          string `json:"legacy"`
	Pauper          string `json:"pauper"`
	Vintage         string `json:"vintage"`
	Penny           string `json:"penny"`
	Commander       string `json:"commander"`
	Oathbreaker     string `json:"oathbreaker"`
	StandardBrawl   string `json:"standardbrawl"`
	Brawl           string `json:"brawl"`
	Alchemy         string `json:"alchemy"`
	PauperCommander string `json:"paupercommander"`
	Duel            string `json:"duel"`
	Oldschool       string `json:"oldschool"`
	Premodern       string `json:"premodern"`
	Predh           string `json:"predh"`
}

type OracleCardPrices struct {
	USD       string `json:"usd"`
	USDFoil   string `json:"usd_foil"`
	USDEtched string `json:"usd_etched"`
	EUR       string `json:"eur"`
	EURFoil   string `json:"eur_foil"`
	Tix       string `json:"tix"`
}

type OracleCard struct {
	ID              string               `json:"id"`
	OracleID        string               `json:"oracle_id"`
	MultiverseIDs   []int                `json:"multiverse_ids"`
	MTGOID          int                  `json:"mtgo_id"`
	MTGOFoilID      int                  `json:"mtgo_foil_id"`
	TCGPlayerID     int                  `json:"tcgplayer_id"`
	CardMarketID    int                  `json:"cardmarket_id"`
	Name            string               `json:"name"`
	Lang            string               `json:"lang"`
	ReleasedAt      string               `json:"released_at"`
	URI             string               `json:"uri"`
	ScryfallURI     string               `json:"scryfall_uri"`
	Layout          string               `json:"layout"`
	HighResImage    bool                 `json:"highres_image"`
	ImageStatus     string               `json:"image_status"`
	ImageURIs       OracleCardImageURIs  `json:"image_uris"`
	ManaCost        string               `json:"mana_cost"`
	CMC             float64              `json:"cmc"`
	TypeLine        string               `json:"type_line"`
	OracleText      string               `json:"oracle_text"`
	Colors          []string             `json:"colors"`
	ColorIdentity   []string             `json:"color_identity"`
	Keywords        []string             `json:"keywords"`
	Legalities      OracleCardLegalities `json:"legalities"`
	Games           []string             `json:"games"`
	Reserved        bool                 `json:"reserved"`
	GameChanger     bool                 `json:"game_changer"`
	Foil            bool                 `json:"foil"`
	NonFoil         bool                 `json:"nonfoil"`
	Finishes        []string             `json:"finishes"`
	Oversized       bool                 `json:"oversized"`
	Promo           bool                 `json:"promo"`
	Reprint         bool                 `json:"reprint"`
	Variation       bool                 `json:"variation"`
	SetID           string               `json:"set_id"`
	Set             string               `json:"set"`
	SetName         string               `json:"set_name"`
	SetType         string               `json:"set_type"`
	SetURI          string               `json:"set_uri"`
	SetSearchURI    string               `json:"set_search_uri"`
	ScryfallSetURI  string               `json:"scryfall_set_uri"`
	RulingsURI      string               `json:"rulings_uri"`
	PrintsSearchURI string               `json:"prints_search_uri"`
	CollectorNumber string               `json:"collector_number"`
	Digital         bool                 `json:"digital"`
	Rarity          string               `json:"rarity"`
	FlavorText      string               `json:"flavor_text"`
	CardBackID      string               `json:"card_back_id"`
	Artist          string               `json:"artist"`
	ArtistIDs       []string             `json:"artist_ids"`
	IllustrationID  string               `json:"illustration_id"`
	BorderColor     string               `json:"border_color"`
	Frame           string               `json:"frame"`
	FullArt         bool                 `json:"full_art"`
	Textless        bool                 `json:"textless"`
	Booster         bool                 `json:"booster"`
	StorySpotlight  bool                 `json:"story_spotlight"`
	EDHRecRank      int                  `json:"edhrec_rank"`
	Prices          OracleCardPrices     `json:"prices"`
}
