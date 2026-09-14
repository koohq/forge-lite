package i18n

import (
	"testing"

	"forge-lite/internal/model"
	"forge-lite/internal/theme"
)

func TestHubMenu_ClassicFantasy_JA(t *testing.T) {
	th := theme.GetPresetTheme("classic_fantasy", model.LanguageJA)
	msg := GetMessages(model.LanguageJA)

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: ダンジョンへ出撃 (探索開始)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb), "2: 鍛冶屋で鍛える (攻撃力・出力の強化)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: 装着 (パッシブ効果の装着)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: オーブを分解 (素材への還元)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: オーブを合成 (上位性能への強化)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: 体力を強化する (最大体力 +10)"},
		{"7: Settings", msg.HubMenu7Settings, "7: 設定 / Settings"},
		{"0: Exit", msg.HubMenu0Exit, "0: ゲーム終了"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.actual != tc.expected {
				t.Errorf("\ngot : %s\nwant: %s", tc.actual, tc.expected)
			}
		})
	}
}

func TestHubMenu_Cyberpunk_JA(t *testing.T) {
	th := theme.GetPresetTheme("cyberpunk", model.LanguageJA)
	msg := GetMessages(model.LanguageJA)

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: ダンジョンへ出撃 (探索開始)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb), "2: システムオーバークロック (攻撃力・出力の強化)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: インストール (パッシブ効果の装着)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: モジュールを分解 (素材への還元)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: モジュールを合成 (上位性能への強化)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: 生体フレーム拡張 (最大耐久 +10)"},
		{"7: Settings", msg.HubMenu7Settings, "7: 設定 / Settings"},
		{"0: Exit", msg.HubMenu0Exit, "0: ゲーム終了"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.actual != tc.expected {
				t.Errorf("\ngot : %s\nwant: %s", tc.actual, tc.expected)
			}
		})
	}
}

func TestHubMenu_PartnerSync_JA(t *testing.T) {
	th := theme.GetPresetTheme("partner_sync", model.LanguageJA)
	msg := GetMessages(model.LanguageJA)

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: ダンジョンへ出撃 (探索開始)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb), "2: 同期率を向上させる (攻撃力・出力の強化)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: セット (パッシブ効果の装着)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: プロトコルを分解 (素材への還元)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: プロトコルを合成 (上位性能への強化)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: 防護プロトコル強化 (生存限界 +10)"},
		{"7: Settings", msg.HubMenu7Settings, "7: 設定 / Settings"},
		{"0: Exit", msg.HubMenu0Exit, "0: ゲーム終了"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.actual != tc.expected {
				t.Errorf("\ngot : %s\nwant: %s", tc.actual, tc.expected)
			}
		})
	}
}

func TestHubMenu_ClassicFantasy_EN(t *testing.T) {
	th := theme.GetPresetTheme("classic_fantasy", model.LanguageEN)
	msg := GetMessages(model.LanguageEN)

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: Embark to Dungeon (Start run)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb), "2: Forge at Blacksmith (Upgrade ATK / Output)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: Equip (Equip passive effects)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: Dismantle Orb (Convert to materials)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: Synthesize Orb (Upgrade to Plus version)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: Fortify Vitality (Max HP +10)"},
		{"7: Settings", msg.HubMenu7Settings, "7: Settings"},
		{"0: Exit", msg.HubMenu0Exit, "0: Quit Game"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.actual != tc.expected {
				t.Errorf("\ngot : %s\nwant: %s", tc.actual, tc.expected)
			}
		})
	}
}

func TestHubMenu_Cyberpunk_EN(t *testing.T) {
	th := theme.GetPresetTheme("cyberpunk", model.LanguageEN)
	msg := GetMessages(model.LanguageEN)

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: Embark to Dungeon (Start run)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb), "2: Overclock System (Upgrade ATK / Output)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: Install (Equip passive effects)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: Dismantle Module (Convert to materials)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: Synthesize Module (Upgrade to Plus version)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: Upgrade Chassis (Max Hull +10)"},
		{"7: Settings", msg.HubMenu7Settings, "7: Settings"},
		{"0: Exit", msg.HubMenu0Exit, "0: Quit Game"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.actual != tc.expected {
				t.Errorf("\ngot : %s\nwant: %s", tc.actual, tc.expected)
			}
		})
	}
}

func TestHubMenu_PartnerSync_EN(t *testing.T) {
	th := theme.GetPresetTheme("partner_sync", model.LanguageEN)
	msg := GetMessages(model.LanguageEN)

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: Embark to Dungeon (Start run)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb), "2: Deepen Sync (Upgrade ATK / Output)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: Set (Equip passive effects)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: Dismantle Protocol (Convert to materials)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: Synthesize Protocol (Upgrade to Plus version)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: Reinforce Shields (Vitality +10)"},
		{"7: Settings", msg.HubMenu7Settings, "7: Settings"},
		{"0: Exit", msg.HubMenu0Exit, "0: Quit Game"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.actual != tc.expected {
				t.Errorf("\ngot : %s\nwant: %s", tc.actual, tc.expected)
			}
		})
	}
}
