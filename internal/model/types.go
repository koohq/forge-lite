package model

type Language = string

const (
	LanguageEN Language = "en"
	LanguageJA Language = "ja"
)

type OrbType string

const (
	OrbMultiHit     OrbType = "MULTI_HIT"
	OrbCritical     OrbType = "CRITICAL"
	OrbVamp         OrbType = "VAMP"
	OrbPoison       OrbType = "POISON"
	OrbMultiHitPlus OrbType = "MULTI_HIT_PLUS"
	OrbCriticalPlus OrbType = "CRITICAL_PLUS"
	OrbVampPlus     OrbType = "VAMP_PLUS"
	OrbPoisonPlus   OrbType = "POISON_PLUS"
)

// Weapon represents player equipment.
type Weapon struct {
	Name    string    `json:"name"`
	BaseAtk int       `json:"baseAtk"`
	Plus    int       `json:"plus"`
	Slots   []OrbType `json:"slots"`
}

// Player represents the player's active state.
type Player struct {
	MaxHP  int    `json:"maxHp"`
	HP     int    `json:"hp"`
	Weapon Weapon `json:"weapon"`
}

// Enemy represents an encounter opponent.
type Enemy struct {
	Name   string `json:"name"`
	HP     int    `json:"hp"`
	MaxHP  int    `json:"maxHp"`
	Atk    int    `json:"atk"`
	Poison int    `json:"poison"`
}

// RunInventory holds temporary drops acquired during an expedition.
type RunInventory struct {
	Scrolls int       `json:"scrolls"`
	Orbs    []OrbType `json:"orbs"`
}

// SaveData matches the JSON persistence format of save.json.
type SaveData struct {
	Weapon       Weapon    `json:"weapon"`
	StockScrolls int       `json:"stockScrolls"`
	StockOrbs    []OrbType `json:"stockOrbs"`
	MaxHP        *int      `json:"maxHp,omitempty"`
	DeepestFloor *int      `json:"deepestFloor,omitempty"`
	Language     *string   `json:"language,omitempty"`
	ThemeID      *string   `json:"themeId,omitempty"`
}

// EnemyTemplate defines an enemy generation archetype.
type EnemyTemplate struct {
	Name     string
	HP       int
	Atk      int
	GetDrops func() RunInventory
}

// DungeonDef defines a dungeon's structure and enemy encounter generator.
type DungeonDef struct {
	Name        string
	Floors      int
	Recommended string
	GetEnemy    func(floor int) (EnemyTemplate, error)
}

// OrbDesc holds a localized name and short description of an orb.
type OrbDesc struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// Dungeons holds localized dungeon names.
type Dungeons struct {
	Starter string `json:"starter"`
	Deep    string `json:"deep"`
	Abyss   string `json:"abyss"`
}

// TargetRank represents a weapon/partner rank threshold and display name.
type TargetRank struct {
	MinPlus int    `json:"minPlus"`
	Name    string `json:"name"`
}

// ThemeLocaleData holds all display vocabulary and configuration for a locale.
type ThemeLocaleData struct {
	Name              string              `json:"name"`
	TargetNameLabel   string              `json:"targetNameLabel"`
	DefaultTargetName string              `json:"defaultTargetName"`
	TargetRanks       []TargetRank        `json:"targetRanks,omitempty"`
	EnhanceVerb       string              `json:"enhanceVerb"`
	ResourceName      string              `json:"resourceName"`
	Orbs              map[OrbType]OrbDesc `json:"orbs"`
	Dungeons          Dungeons            `json:"dungeons"`
	EnemyPrefixes     []string            `json:"enemyPrefixes"`
	EnemyBases        []string            `json:"enemyBases"`
	RetreatMessage    string              `json:"retreatMessage"`
	HubTitle          string              `json:"hubTitle"`
	HPLabel           string              `json:"hpLabel"`
	TargetPrefix      string              `json:"targetPrefix"`
	StatLabel         string              `json:"statLabel"`
	PlusPrefix        string              `json:"plusPrefix"`
	StorageLabel      string              `json:"storageLabel"`
	OrbLabel          string              `json:"orbLabel"`
	EnhanceHPVerb     string              `json:"enhanceHpVerb"`
	InstallVerb       string              `json:"installVerb"`
}

// GetTargetRankName returns the rank name corresponding to the current plus value.
// If plus >= 1000 or exceeds the highest threshold, the highest rank name is maintained.
func (d ThemeLocaleData) GetTargetRankName(plus int) string {
	if len(d.TargetRanks) == 0 {
		return d.DefaultTargetName
	}
	matched := d.DefaultTargetName
	for _, r := range d.TargetRanks {
		if plus >= r.MinPlus {
			matched = r.Name
		}
	}
	return matched
}

// GetTargetPrefix returns TargetPrefix if non-empty, falling back to TargetNameLabel.
func (d ThemeLocaleData) GetTargetPrefix() string {
	if d.TargetPrefix != "" {
		return d.TargetPrefix
	}
	return d.TargetNameLabel
}

// ThemeDefinition represents a bilingual theme dataset.
type ThemeDefinition struct {
	ID string          `json:"id"`
	JA ThemeLocaleData `json:"ja"`
	EN ThemeLocaleData `json:"en"`
}

// Theme represents an active resolved theme.
type Theme struct {
	ID       string
	Language Language
	ThemeLocaleData
}

