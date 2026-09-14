package theme

import (
	"os"
	"testing"

	"github.com/koohq/forge-lite/internal/model"
)

func TestThemePresets_Registered(t *testing.T) {
	presets := []string{"classic_fantasy", "cyberpunk", "partner_sync"}
	for _, id := range presets {
		theme := GetPresetTheme(id, model.LanguageJA)
		if theme.ID != id {
			t.Errorf("expected theme ID %s, got %s", id, theme.ID)
		}
	}
}

func TestClassicFantasy_Fidelity(t *testing.T) {
	ja := GetPresetTheme("classic_fantasy", model.LanguageJA)
	en := GetPresetTheme("classic_fantasy", model.LanguageEN)

	if ja.DefaultTargetName != "どうのつるぎ" {
		t.Errorf("unexpected ja.DefaultTargetName: %s", ja.DefaultTargetName)
	}
	if en.DefaultTargetName != "Bronze Sword" {
		t.Errorf("unexpected en.DefaultTargetName: %s", en.DefaultTargetName)
	}
	if ja.EnhanceVerb != "鍛冶屋で鍛える" {
		t.Errorf("unexpected ja.EnhanceVerb: %s", ja.EnhanceVerb)
	}
	if en.EnhanceVerb != "Forge at Blacksmith" {
		t.Errorf("unexpected en.EnhanceVerb: %s", en.EnhanceVerb)
	}
	if ja.HubTitle != "【冒険者の酒場】" {
		t.Errorf("unexpected ja.HubTitle: %s", ja.HubTitle)
	}
	if en.HubTitle != "【Adventurer's Guild】" {
		t.Errorf("unexpected en.HubTitle: %s", en.HubTitle)
	}
	if ja.HPLabel != "最大体力" || en.HPLabel != "Max HP" {
		t.Errorf("unexpected hpLabel: ja=%s, en=%s", ja.HPLabel, en.HPLabel)
	}
}

func TestCyberpunk_Fidelity(t *testing.T) {
	ja := GetPresetTheme("cyberpunk", model.LanguageJA)
	en := GetPresetTheme("cyberpunk", model.LanguageEN)

	if ja.DefaultTargetName != "パルスブレード" {
		t.Errorf("unexpected ja.DefaultTargetName: %s", ja.DefaultTargetName)
	}
	if en.DefaultTargetName != "Pulse Blade" {
		t.Errorf("unexpected en.DefaultTargetName: %s", en.DefaultTargetName)
	}
	if ja.HubTitle != "【セーフハウス】" {
		t.Errorf("unexpected ja.HubTitle: %s", ja.HubTitle)
	}
	if en.HubTitle != "【Safehouse】" {
		t.Errorf("unexpected en.HubTitle: %s", en.HubTitle)
	}
	if ja.InstallVerb != "インストール" || en.InstallVerb != "Install" {
		t.Errorf("unexpected installVerb: ja=%s, en=%s", ja.InstallVerb, en.InstallVerb)
	}
}

func TestPartnerSync_Fidelity(t *testing.T) {
	ja := GetPresetTheme("partner_sync", model.LanguageJA)
	en := GetPresetTheme("partner_sync", model.LanguageEN)

	if ja.DefaultTargetName != "戦術アンドロイド「アイリス」" {
		t.Errorf("unexpected ja.DefaultTargetName: %s", ja.DefaultTargetName)
	}
	if en.DefaultTargetName != "Tactical Android \"Iris\"" {
		t.Errorf("unexpected en.DefaultTargetName: %s", en.DefaultTargetName)
	}
	if ja.PlusPrefix != " Sync: +" || en.PlusPrefix != " Sync: +" {
		t.Errorf("unexpected plusPrefix: ja=%s, en=%s", ja.PlusPrefix, en.PlusPrefix)
	}
	if ja.HubTitle != "【作戦拠点】" || en.HubTitle != "【Operation Base】" {
		t.Errorf("unexpected hubTitle: ja=%s, en=%s", ja.HubTitle, en.HubTitle)
	}
	if ja.StorageLabel != "ロッカー" || en.StorageLabel != "Locker" {
		t.Errorf("unexpected storageLabel: ja=%s, en=%s", ja.StorageLabel, en.StorageLabel)
	}
	if ja.GetSlotLabel() != "連携スロット" || en.GetSlotLabel() != "Sync Slots" {
		t.Errorf("unexpected slotLabel: ja=%s, en=%s", ja.GetSlotLabel(), en.GetSlotLabel())
	}
	if ja.GetUninstalledOrbLabel() != "未設定の連携スキル" || en.GetUninstalledOrbLabel() != "Standby Sync Skills" {
		t.Errorf("unexpected uninstalledOrbLabel: ja=%s, en=%s", ja.GetUninstalledOrbLabel(), en.GetUninstalledOrbLabel())
	}
	if ja.Dungeons.Deep != "汚染隔離区" || en.Dungeons.Deep != "Contaminated Quarantine Zone" {
		t.Errorf("unexpected deep dungeon: ja=%s, en=%s", ja.Dungeons.Deep, en.Dungeons.Deep)
	}
	if ja.HPLabel != "継戦力" || en.HPLabel != "Endurance" {
		t.Errorf("unexpected hpLabel: ja=%s, en=%s", ja.HPLabel, en.HPLabel)
	}
	if ja.StatLabel != "連携力" || en.StatLabel != "Sync Power" {
		t.Errorf("unexpected statLabel: ja=%s, en=%s", ja.StatLabel, en.StatLabel)
	}
	if ja.ResourceName != "訓練記録" || en.ResourceName != "Training Log" {
		t.Errorf("unexpected resourceName: ja=%s, en=%s", ja.ResourceName, en.ResourceName)
	}
	if ja.EnhanceVerb != "アイリスと訓練する" || en.EnhanceVerb != "Train with Iris" {
		t.Errorf("unexpected enhanceVerb: ja=%s, en=%s", ja.EnhanceVerb, en.EnhanceVerb)
	}
	if ja.InstallVerb != "連携スキルを設定" || en.InstallVerb != "Set Sync" {
		t.Errorf("unexpected installVerb: ja=%s, en=%s", ja.InstallVerb, en.InstallVerb)
	}
}

func TestTargetRanks_Thresholds(t *testing.T) {
	cfJA := GetPresetTheme("classic_fantasy", model.LanguageJA)
	testsCF := []struct {
		plus     int
		expected string
	}{
		{0, "どうのつるぎ"},
		{10, "どうのつるぎ"},
		{49, "どうのつるぎ"},
		{50, "はがねのつるぎ"},
		{199, "はがねのつるぎ"},
		{200, "勇者のつるぎ"},
		{499, "勇者のつるぎ"},
		{500, "竜殺しの神剣"},
		{999, "竜殺しの神剣"},
		{1000, "終焉を断つ刃"},
		{2188, "終焉を断つ刃"},
	}
	for _, tc := range testsCF {
		got := cfJA.GetTargetRankName(tc.plus)
		if got != tc.expected {
			t.Errorf("classic_fantasy plus=%d: got %s, want %s", tc.plus, got, tc.expected)
		}
	}

	psJA := GetPresetTheme("partner_sync", model.LanguageJA)
	testsPS := []struct {
		plus     int
		expected string
	}{
		{0, "戦術アンドロイド「アイリス」"},
		{50, "共闘に慣れてきたアイリス"},
		{200, "息の合うアイリス"},
		{500, "頼れるアイリス"},
		{1000, "歴戦のアイリス"},
		{5000, "歴戦のアイリス"},
	}
	for _, tc := range testsPS {
		got := psJA.GetTargetRankName(tc.plus)
		if got != tc.expected {
			t.Errorf("partner_sync plus=%d: got %s, want %s", tc.plus, got, tc.expected)
		}
	}

	cpEN := GetPresetTheme("cyberpunk", model.LanguageEN)
	testsCP := []struct {
		plus     int
		expected string
	}{
		{0, "Pulse Blade"},
		{50, "Pulse Blade Mk-II"},
		{200, "High-Frequency Katana"},
		{500, "Plasma Saber"},
		{1000, "Antimatter Void Edge"},
	}
	for _, tc := range testsCP {
		got := cpEN.GetTargetRankName(tc.plus)
		if got != tc.expected {
			t.Errorf("cyberpunk plus=%d: got %s, want %s", tc.plus, got, tc.expected)
		}
	}
}

func TestExportCustomThemeTemplate(t *testing.T) {
	tmpFile := "test_custom_theme_template.json"
	defer os.Remove(tmpFile)

	err := ExportCustomThemeTemplate(tmpFile)
	if err != nil {
		t.Fatalf("failed to export template: %v", err)
	}

	def, err := LoadCustomThemeDefIfExists(tmpFile)
	if err != nil {
		t.Fatalf("failed to load exported template: %v", err)
	}
	if def == nil {
		t.Fatal("expected loaded def to not be nil")
	}
	if def.ID != "custom_scifi" {
		t.Errorf("expected ID custom_scifi, got %s", def.ID)
	}
	if def.JA.ResourceName != "プラズマコア" {
		t.Errorf("unexpected ja resource: %s", def.JA.ResourceName)
	}
	if def.EN.ResourceName != "Plasma Core" {
		t.Errorf("unexpected en resource: %s", def.EN.ResourceName)
	}
	if len(def.JA.TargetRanks) != 5 {
		t.Errorf("expected 5 ja target ranks, got %d", len(def.JA.TargetRanks))
	}
	if def.JA.GetTargetRankName(1500) != "次元崩壊特異点兵装" {
		t.Errorf("unexpected custom rank name: %s", def.JA.GetTargetRankName(1500))
	}
}
