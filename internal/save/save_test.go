package save

import (
	"os"
	"testing"

	"forge-lite/internal/model"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	tmpFile := "test_roundtrip_save.json"
	defer os.Remove(tmpFile)

	maxHp := 120
	deepest := 15
	lang := model.LanguageJA
	themeID := "cyberpunk"

	original := &model.SaveData{
		Weapon: model.Weapon{
			Name:    "パルスブレード",
			BaseAtk: 10,
			Plus:    5,
			Slots:   []model.OrbType{model.OrbMultiHitPlus, model.OrbVamp},
		},
		StockScrolls: 7,
		StockOrbs:    []model.OrbType{model.OrbCritical, model.OrbPoisonPlus},
		MaxHP:        &maxHp,
		DeepestFloor: &deepest,
		Language:     &lang,
		ThemeID:      &themeID,
	}

	err := SaveSaveData(tmpFile, original)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	loaded, err := LoadSaveData(tmpFile)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.Weapon.Name != original.Weapon.Name || loaded.Weapon.Plus != original.Weapon.Plus {
		t.Errorf("weapon mismatch: got %+v, want %+v", loaded.Weapon, original.Weapon)
	}
	if len(loaded.Weapon.Slots) != 2 || loaded.Weapon.Slots[0] != model.OrbMultiHitPlus {
		t.Errorf("weapon slots mismatch: %+v", loaded.Weapon.Slots)
	}
	if loaded.StockScrolls != 7 || len(loaded.StockOrbs) != 2 {
		t.Errorf("inventory mismatch: scrolls=%d, orbs=%+v", loaded.StockScrolls, loaded.StockOrbs)
	}
	if *loaded.MaxHP != 120 || *loaded.DeepestFloor != 15 {
		t.Errorf("stats mismatch: maxHp=%d, deepest=%d", *loaded.MaxHP, *loaded.DeepestFloor)
	}
	if *loaded.Language != "ja" || *loaded.ThemeID != "cyberpunk" {
		t.Errorf("meta mismatch: lang=%s, theme=%s", *loaded.Language, *loaded.ThemeID)
	}
}

func TestLoad_DefaultsWhenMissingFields(t *testing.T) {
	tmpFile := "test_defaults_save.json"
	defer os.Remove(tmpFile)

	minimalJSON := []byte(`{
		"weapon": {
			"name": "どうのつるぎ",
			"baseAtk": 10,
			"plus": 0
		},
		"stockScrolls": 0,
		"stockOrbs": []
	}`)
	if err := os.WriteFile(tmpFile, minimalJSON, 0644); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadSaveData(tmpFile)
	if err != nil {
		t.Fatalf("failed to load minimal save: %v", err)
	}

	if *loaded.MaxHP != 50 {
		t.Errorf("expected default MaxHP 50, got %d", *loaded.MaxHP)
	}
	if *loaded.DeepestFloor != 0 {
		t.Errorf("expected default DeepestFloor 0, got %d", *loaded.DeepestFloor)
	}
	if *loaded.ThemeID != "classic_fantasy" {
		t.Errorf("expected default ThemeID classic_fantasy, got %s", *loaded.ThemeID)
	}
	if loaded.Weapon.Slots == nil || len(loaded.Weapon.Slots) != 0 {
		t.Errorf("expected empty non-nil slots")
	}
}
