package theme

import (
	"os"
	"testing"

	"forge-lite/internal/model"
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
	if ja.HubTitle != "【セーフハウス・端末】" {
		t.Errorf("unexpected ja.HubTitle: %s", ja.HubTitle)
	}
	if en.HubTitle != "【Safehouse Terminal】" {
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
	if ja.PlusPrefix != " Sync:+" || en.PlusPrefix != " Sync:+" {
		t.Errorf("unexpected plusPrefix: ja=%s, en=%s", ja.PlusPrefix, en.PlusPrefix)
	}
	if ja.HubTitle != "【司令室・ドック】" || en.HubTitle != "【Command Dock】" {
		t.Errorf("unexpected hubTitle: ja=%s, en=%s", ja.HubTitle, en.HubTitle)
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
}
