package engine

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/koohq/forge-lite/internal/model"
	"github.com/koohq/forge-lite/internal/save"
	"github.com/koohq/forge-lite/internal/theme"
)

func TestEngine_LoadGame_DynamicRankResolution(t *testing.T) {
	// Setup a temporary save file with plus: 2188 and old default weapon name
	tmpSave := "test_save_rank.json"
	maxHp := 50
	deepest := 15
	lang := model.LanguageJA
	themeID := "classic_fantasy"

	saveData := &model.SaveData{
		Weapon: model.Weapon{
			Name:    "どうのつるぎ",
			BaseAtk: 10,
			Plus:    2188,
			Slots:   []model.OrbType{},
		},
		StockScrolls: 10,
		StockOrbs:    []model.OrbType{},
		MaxHP:        &maxHp,
		DeepestFloor: &deepest,
		Language:     &lang,
		ThemeID:      &themeID,
	}

	err := save.SaveSaveData(tmpSave, saveData)
	if err != nil {
		t.Fatalf("failed to save test save: %v", err)
	}
	defer os.Remove(tmpSave)

	in := bytes.NewBufferString("1\n0\n")
	game := NewGame(in)
	game.saveFilePath = tmpSave
	if !game.loadGame() {
		t.Fatal("expected loadGame to succeed")
	}

	// Verify weapon name is dynamically upgraded according to plus: 2188
	if game.player.Weapon.Name != "終焉を断つ刃" {
		t.Errorf("expected weapon name '終焉を断つ刃', got '%s'", game.player.Weapon.Name)
	}

	// Verify theme switch updates weapon name and abyss name
	game.theme = theme.GetPresetTheme("partner_sync", game.language)
	game.player.Weapon.Name = game.theme.GetTargetRankName(game.player.Weapon.Plus)
	if game.player.Weapon.Name != "歴戦の相棒アイリス" {
		t.Errorf("expected partner name '歴戦の相棒アイリス', got '%s'", game.player.Weapon.Name)
	}
	if game.theme.Dungeons.Abyss != "未知の最深部" {
		t.Errorf("expected abyss name '未知の最深部', got '%s'", game.theme.Dungeons.Abyss)
	}
	if game.theme.StatLabel != "連携力" {
		t.Errorf("expected stat label '連携力', got '%s'", game.theme.StatLabel)
	}

	// Verify cyberpunk theme switch
	game.theme = theme.GetPresetTheme("cyberpunk", game.language)
	game.player.Weapon.Name = game.theme.GetTargetRankName(game.player.Weapon.Plus)
	if game.player.Weapon.Name != "反物質崩壊刃" {
		t.Errorf("expected weapon name '反物質崩壊刃', got '%s'", game.player.Weapon.Name)
	}
	if game.theme.Dungeons.Abyss != "無限の電脳網" {
		t.Errorf("expected abyss name '無限の電脳網', got '%s'", game.theme.Dungeons.Abyss)
	}
}

func TestEngine_UpgradeWeapon_RankUpgradedFanfare(t *testing.T) {
	in := bytes.NewBufferString("")
	game := NewGame(in)
	game.language = model.LanguageJA
	game.theme = theme.GetPresetTheme("classic_fantasy", model.LanguageJA)
	game.player.Weapon = model.Weapon{
		Name:    "どうのつるぎ",
		BaseAtk: 10,
		Plus:    49,
		Slots:   []model.OrbType{},
	}
	game.stockScrolls = 5

	// Check pre-rank
	prevRank := game.theme.GetTargetRankName(game.player.Weapon.Plus)
	if prevRank != "どうのつるぎ" {
		t.Fatalf("expected initial rank 'どうのつるぎ', got '%s'", prevRank)
	}

	// Enhance by 1 (49 -> 50)
	EnhanceWeapon(&game.player.Weapon, &game.stockScrolls, 1)
	newRank := game.theme.GetTargetRankName(game.player.Weapon.Plus)
	game.player.Weapon.Name = newRank

	if newRank != "はがねのつるぎ" {
		t.Errorf("expected upgraded rank 'はがねのつるぎ', got '%s'", newRank)
	}

	msgJA := game.msg().RankUpgraded(newRank)
	if !strings.Contains(msgJA, "【はがねのつるぎ】に進化した！") {
		t.Errorf("unexpected JA rank upgraded message: %s", msgJA)
	}

	// Test English mode
	game.language = model.LanguageEN
	msgEN := game.msg().RankUpgraded("Steel Sword")
	if !strings.Contains(msgEN, "Rank upgraded! Evolved to 【Steel Sword】!") {
		t.Errorf("unexpected EN rank upgraded message: %s", msgEN)
	}
}
