package theme

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/koohq/forge-lite/internal/model"
)

//go:embed embedded/*.json
var embeddedFS embed.FS

const (
	DefaultThemeID         = "classic_fantasy"
	CustomThemeFilePath    = "custom_theme.json"
	CustomThemeExamplePath = "custom_theme.example.json"
)

var (
	presetCache = make(map[string]model.ThemeDefinition)
)

func init() {
	presets := []string{"classic_fantasy", "cyberpunk", "partner_sync"}
	for _, id := range presets {
		data, err := embeddedFS.ReadFile("embedded/" + id + ".json")
		if err != nil {
			panic(fmt.Sprintf("failed to read embedded theme %s: %v", id, err))
		}
		var def model.ThemeDefinition
		if err := json.Unmarshal(data, &def); err != nil {
			panic(fmt.Sprintf("failed to parse embedded theme %s: %v", id, err))
		}
		presetCache[id] = def
	}
}

// ResolveTheme resolves a ThemeDefinition for the requested language.
func ResolveTheme(def model.ThemeDefinition, lang model.Language) model.Theme {
	var locale model.ThemeLocaleData
	if lang == model.LanguageJA {
		if def.JA.Name != "" {
			locale = def.JA
		} else {
			locale = def.EN
		}
	} else {
		if def.EN.Name != "" {
			locale = def.EN
		} else {
			locale = def.JA
		}
	}

	return model.Theme{
		ID:              def.ID,
		Language:        lang,
		ThemeLocaleData: locale,
	}
}

// GetPresetTheme returns a resolved preset theme by ID and language.
func GetPresetTheme(themeID string, lang model.Language) model.Theme {
	def, ok := presetCache[themeID]
	if !ok {
		def = presetCache[DefaultThemeID]
	}
	return ResolveTheme(def, lang)
}

// GetAllPresetDefinitions returns all embedded preset definitions.
func GetAllPresetDefinitions() []model.ThemeDefinition {
	return []model.ThemeDefinition{
		presetCache["classic_fantasy"],
		presetCache["cyberpunk"],
		presetCache["partner_sync"],
	}
}

// LoadCustomThemeDefIfExists attempts to load and validate a custom theme definition.
func LoadCustomThemeDefIfExists(filePath string) (*model.ThemeDefinition, error) {
	if filePath == "" {
		filePath = CustomThemeFilePath
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	// Try standard ThemeDefinition format first
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid json in %s: %w", filePath, err)
	}

	idRaw, hasID := raw["id"]
	if !hasID {
		return nil, errors.New("theme missing required 'id' property")
	}
	var id string
	if err := json.Unmarshal(idRaw, &id); err != nil || id == "" {
		return nil, errors.New("theme 'id' must be a non-empty string")
	}

	var jaLocale, enLocale *model.ThemeLocaleData
	if jaRaw, hasJA := raw["ja"]; hasJA {
		var l model.ThemeLocaleData
		if err := json.Unmarshal(jaRaw, &l); err == nil && validateLocaleData(&l) {
			jaLocale = &l
		}
	}
	if enRaw, hasEN := raw["en"]; hasEN {
		var l model.ThemeLocaleData
		if err := json.Unmarshal(enRaw, &l); err == nil && validateLocaleData(&l) {
			enLocale = &l
		}
	}

	if jaLocale != nil && enLocale != nil {
		fillLocaleDefaults(jaLocale)
		fillLocaleDefaults(enLocale)
		return &model.ThemeDefinition{
			ID: id,
			JA: *jaLocale,
			EN: *enLocale,
		}, nil
	}

	// Fallback to legacy flat format
	var legacy model.ThemeLocaleData
	if err := json.Unmarshal(data, &legacy); err == nil && validateLocaleData(&legacy) {
		fillLocaleDefaults(&legacy)
		finalJA := legacy
		finalEN := legacy
		if jaLocale != nil {
			fillLocaleDefaults(jaLocale)
			finalJA = *jaLocale
		}
		if enLocale != nil {
			fillLocaleDefaults(enLocale)
			finalEN = *enLocale
		}
		return &model.ThemeDefinition{
			ID: id,
			JA: finalJA,
			EN: finalEN,
		}, nil
	}

	return nil, errors.New("theme missing required locale properties or orbs")
}

func validateLocaleData(data *model.ThemeLocaleData) bool {
	targetPrefix := data.GetTargetPrefix()
	if data.Name == "" || targetPrefix == "" || (data.DefaultTargetName == "" && len(data.TargetRanks) == 0) ||
		data.EnhanceVerb == "" || data.ResourceName == "" || data.RetreatMessage == "" {
		return false
	}
	if data.Dungeons.Starter == "" || data.Dungeons.Deep == "" || data.Dungeons.Abyss == "" {
		return false
	}
	if len(data.EnemyPrefixes) == 0 || len(data.EnemyBases) == 0 {
		return false
	}

	requiredOrbs := []model.OrbType{
		model.OrbMultiHit,
		model.OrbCritical,
		model.OrbVamp,
		model.OrbPoison,
		model.OrbMultiHitPlus,
		model.OrbCriticalPlus,
		model.OrbVampPlus,
		model.OrbPoisonPlus,
	}
	for _, orb := range requiredOrbs {
		desc, ok := data.Orbs[orb]
		if !ok || desc.Name == "" || desc.Desc == "" {
			return false
		}
	}
	return true
}

func fillLocaleDefaults(data *model.ThemeLocaleData) {
	if data.DefaultTargetName == "" && len(data.TargetRanks) > 0 {
		data.DefaultTargetName = data.TargetRanks[0].Name
	}
	if data.HubTitle == "" {
		data.HubTitle = "【拠点】"
	}
	if data.HPLabel == "" {
		data.HPLabel = "最大体力"
	}
	if data.TargetPrefix == "" {
		data.TargetPrefix = data.TargetNameLabel
	}
	if data.TargetNameLabel == "" {
		data.TargetNameLabel = data.TargetPrefix
	}
	if data.StatLabel == "" {
		data.StatLabel = "攻撃力"
	}
	if data.PlusPrefix == "" {
		data.PlusPrefix = "+"
	}
	if data.StorageLabel == "" {
		data.StorageLabel = "倉庫"
	}
	if data.OrbLabel == "" {
		data.OrbLabel = "オーブ"
	}
	if data.EnhanceHPVerb == "" {
		data.EnhanceHPVerb = "体力を強化する"
	}
	if data.InstallVerb == "" {
		data.InstallVerb = "装着"
	}
}

// LoadCustomThemeIfExists loads and resolves custom theme if exists.
func LoadCustomThemeIfExists(filePath string, lang model.Language) (*model.Theme, error) {
	def, err := LoadCustomThemeDefIfExists(filePath)
	if err != nil || def == nil {
		return nil, err
	}
	resolved := ResolveTheme(*def, lang)
	return &resolved, nil
}

// ExportCustomThemeTemplate exports the template for custom theme.
func ExportCustomThemeTemplate(filePath string) error {
	if filePath == "" {
		filePath = CustomThemeExamplePath
	}

	template := model.ThemeDefinition{
		ID: "custom_scifi",
		JA: model.ThemeLocaleData{
			Name:              "SF星間探査 (Sci-Fi Star Explorer)",
			TargetNameLabel:   "旗艦武装",
			DefaultTargetName: "フォトンランス",
			TargetRanks: []model.TargetRank{
				{MinPlus: 0, Name: "フォトンランス"},
				{MinPlus: 50, Name: "フォトンランス Mk-II"},
				{MinPlus: 200, Name: "ハイパーレーザー砲"},
				{MinPlus: 500, Name: "超時空波動砲"},
				{MinPlus: 1000, Name: "次元崩壊特異点兵装"},
			},
			EnhanceVerb:  "リアクターを調整する",
			ResourceName: "プラズマコア",
			Orbs: map[model.OrbType]model.OrbDesc{
				model.OrbMultiHit:     {Name: "ツインビーム", Desc: "2連照射/威力65%"},
				model.OrbCritical:     {Name: "クリティカルパルス", Desc: "25%で2倍"},
				model.OrbVamp:         {Name: "シールドドレイン", Desc: "与ダメの15%回復"},
				model.OrbPoison:       {Name: "アシッド腐食", Desc: "攻撃時毒+1/ターン末毒x3ダメ"},
				model.OrbMultiHitPlus: {Name: "ツインビーム+", Desc: "2連照射/威力75%"},
				model.OrbCriticalPlus: {Name: "クリティカルパルス+", Desc: "40%で2倍"},
				model.OrbVampPlus:     {Name: "シールドドレイン+", Desc: "与ダメの25%回復"},
				model.OrbPoisonPlus:   {Name: "アシッド腐食+", Desc: "攻撃時毒+2/ターン末毒x3ダメ"},
			},
			Dungeons: model.Dungeons{
				Starter: "小惑星前哨基地",
				Deep:    "遺棄された弩級艦",
				Abyss:   "事象の地平線",
			},
			EnemyPrefixes: []string{
				"敵対的な",
				"過負荷の",
				"変異した",
				"未知生命の",
				"サイバネの",
				"凶暴な",
				"古代の",
				"コズミック",
			},
			EnemyBases: []string{
				"偵察ドローン",
				"バイオスウォーム",
				"宇宙海賊",
				"戦闘メカ",
				"虚空の怪異",
				"星喰らい",
				"ナノ集合体",
				"タイタン弩級艦",
			},
			RetreatMessage: "ワープドライブを緊急起動し、軌道ステーションへ退避した。",
			HubTitle:       "【旗艦ドック・管制室】",
			HPLabel:        "シールド最大値",
			TargetPrefix:   "旗艦武装",
			StatLabel:      "兵装出力",
			PlusPrefix:     "+",
			StorageLabel:   "カーゴ",
			OrbLabel:       "デバイス",
			EnhanceHPVerb:  "シールド容量を拡張",
			InstallVerb:    "換装",
		},
		EN: model.ThemeLocaleData{
			Name:              "Sci-Fi Star Explorer",
			TargetNameLabel:   "Flagship Weapon",
			DefaultTargetName: "Photon Lance",
			TargetRanks: []model.TargetRank{
				{MinPlus: 0, Name: "Photon Lance"},
				{MinPlus: 50, Name: "Photon Lance Mk-II"},
				{MinPlus: 200, Name: "Hyper Laser Cannon"},
				{MinPlus: 500, Name: "Spatiotemporal Wave Cannon"},
				{MinPlus: 1000, Name: "Dimensional Singularity Armament"},
			},
			EnhanceVerb:  "Calibrate Reactor",
			ResourceName: "Plasma Core",
			Orbs: map[model.OrbType]model.OrbDesc{
				model.OrbMultiHit:     {Name: "Twin Beam", Desc: "Fire dual beams at 65% energy each"},
				model.OrbCritical:     {Name: "Critical Pulse", Desc: "25% chance for 2x focused damage"},
				model.OrbVamp:         {Name: "Shield Siphon", Desc: "Absorb 15% damage dealt as shield energy"},
				model.OrbPoison:       {Name: "Corrosive Acid", Desc: "+1 chemical payload on hit, deals 3x dmg per tick"},
				model.OrbMultiHitPlus: {Name: "Twin Beam+", Desc: "Fire dual beams at 75% energy each"},
				model.OrbCriticalPlus: {Name: "Critical Pulse+", Desc: "40% chance for 2x focused damage"},
				model.OrbVampPlus:     {Name: "Shield Siphon+", Desc: "Absorb 25% damage dealt as shield energy"},
				model.OrbPoisonPlus:   {Name: "Corrosive Acid+", Desc: "+2 chemical payload on hit, deals 3x dmg per tick"},
			},
			Dungeons: model.Dungeons{
				Starter: "Asteroid Outpost",
				Deep:    "Derelict Dreadnought",
				Abyss:   "Event Horizon Void",
			},
			EnemyPrefixes: []string{
				"Cosmic",
				"Void",
				"Cybernetic",
				"Plasma",
				"Starlight",
				"Hyper-Drive",
				"Gravity-Distorted",
				"Supernova",
			},
			EnemyBases: []string{
				"Scout Drone",
				"Bio-Swarm",
				"Pirate Raider",
				"War Mech",
				"Void Beast",
				"Star Devourer",
				"Nanite Hive",
				"Titan Dreadnought",
			},
			RetreatMessage: "Warp drive engaged; executed emergency tactical jump back to orbital station.",
			HubTitle:       "【Flagship Command Dock】",
			HPLabel:        "Max Shields",
			TargetPrefix:   "Flagship Weapon",
			StatLabel:      "Weapon Output",
			PlusPrefix:     "+",
			StorageLabel:   "Cargo Bay",
			OrbLabel:       "Device",
			EnhanceHPVerb:  "Upgrade Shields",
			InstallVerb:    "Fit",
		},
	}

	bytes, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bytes, 0644)
}
