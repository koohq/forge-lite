package i18n

import (
	"strings"
	"testing"

	"github.com/koohq/forge-lite/internal/model"
	"github.com/koohq/forge-lite/internal/theme"
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
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb, th.StatLabel), "2: 鍛冶屋で鍛える (攻撃力の強化)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: 装着 (パッシブ効果の装着)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: オーブを分解 (素材を回収)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: オーブを合成 (上位版へ強化)"},
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
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb, th.StatLabel), "2: システムオーバークロック (出力の強化)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: インストール (パッシブ効果の装着)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: モジュールを分解 (素材を回収)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: モジュールを合成 (上位版へ強化)"},
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
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb, th.StatLabel), "2: アイリスと訓練する (連携力の強化)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: 連携スキルを設定 (パッシブ効果の装着)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: 連携スキルを分解 (素材を回収)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: 連携スキルを合成 (上位版へ強化)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: 継戦力を高める (継戦力 +10)"},
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
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: Embark to Dungeon (Start Run)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb, th.StatLabel), "2: Forge at Blacksmith (Upgrade ATK)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: Equip (Equip Passives)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: Dismantle Orb (Recover Materials)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: Synthesize Orb (Upgrade to Plus)"},
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
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: Embark to Dungeon (Start Run)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb, th.StatLabel), "2: Overclock System (Upgrade Output)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: Install (Equip Passives)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: Dismantle Module (Recover Materials)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: Synthesize Module (Upgrade to Plus)"},
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
		{"1: Dungeon", msg.HubMenu1Dungeon, "1: Embark to Dungeon (Start Run)"},
		{"2: Enhance", msg.HubMenu2Enhance(th.EnhanceVerb, th.StatLabel), "2: Train with Iris (Upgrade Sync Power)"},
		{"3: AttachOrb", msg.HubMenu3AttachOrb(th.InstallVerb), "3: Set Sync (Equip Passives)"},
		{"4: Disassemble", msg.HubMenu4Disassemble(th.OrbLabel), "4: Dismantle Sync Skill (Recover Materials)"},
		{"5: Synthesize", msg.HubMenu5Synthesize(th.OrbLabel), "5: Synthesize Sync Skill (Upgrade to Plus)"},
		{"6: EnhanceHp", msg.HubMenu6EnhanceHp(th.EnhanceHPVerb, th.HPLabel), "6: Train Endurance (Endurance +10)"},
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

func TestHubStats_DynamicAbyss(t *testing.T) {
	jaMsg := GetMessages(model.LanguageJA)
	enMsg := GetMessages(model.LanguageEN)

	jaRes := jaMsg.HubStats("最大体力", 100, "事象の地平線", 42, false)
	expectedJA := "最大体力: 100 | 事象の地平線 最高到達: B42F"
	if jaRes != expectedJA {
		t.Errorf("got: %s, want: %s", jaRes, expectedJA)
	}

	enRes := enMsg.HubStats("Max Shields", 150, "Event Horizon Void", 0, false)
	expectedEN := "Max Shields: 150 | Event Horizon Void Record: None"
	if enRes != expectedEN {
		t.Errorf("got: %s, want: %s", enRes, expectedEN)
	}

	jaCleared := jaMsg.HubStats("最大体力", 100, "事象の地平線", 42, true)
	expectedJACleared := "最大体力: 100 | 事象の地平線 最高到達: B42F ★本編クリア"
	if jaCleared != expectedJACleared {
		t.Errorf("got: %s, want: %s", jaCleared, expectedJACleared)
	}

	enCleared := enMsg.HubStats("Max Shields", 150, "Event Horizon Void", 0, true)
	expectedENCleared := "Max Shields: 150 | Event Horizon Void Record: None ★Game Cleared"
	if enCleared != expectedENCleared {
		t.Errorf("got: %s, want: %s", enCleared, expectedENCleared)
	}
}

func TestAbyssCleared_Message(t *testing.T) {
	jaMsg := GetMessages(model.LanguageJA)
	enMsg := GetMessages(model.LanguageEN)

	tests := []struct {
		themeID    string
		expectedJA string
		expectedEN string
	}{
		{
			themeID: "classic_fantasy",
			expectedJA: "\n**********************************************\n" +
				"   【無限の深淵】B100F 踏破！ 本編クリア、おめでとうございます！\n" +
				"   ここから先も潜り続けて、最高記録に挑戦できます。\n" +
				"**********************************************",
			expectedEN: "\n**********************************************\n" +
				"   【Infinite Abyss】 B100F conquered! You've cleared the game. Congratulations!\n" +
				"   The depths go on. Keep descending to push your record.\n" +
				"**********************************************",
		},
		{
			themeID: "cyberpunk",
			expectedJA: "\n**********************************************\n" +
				"   【無限の電脳網】B100F 踏破！ 本編クリア、おめでとうございます！\n" +
				"   ここから先も潜り続けて、最高記録に挑戦できます。\n" +
				"**********************************************",
			expectedEN: "\n**********************************************\n" +
				"   【Infinite Cyberspace】 B100F conquered! You've cleared the game. Congratulations!\n" +
				"   The depths go on. Keep descending to push your record.\n" +
				"**********************************************",
		},
		{
			themeID: "partner_sync",
			expectedJA: "\n**********************************************\n" +
				"   【未知の最深部】B100F 踏破！ 本編クリア、おめでとうございます！\n" +
				"   ここから先も潜り続けて、最高記録に挑戦できます。\n" +
				"**********************************************",
			expectedEN: "\n**********************************************\n" +
				"   【Unknown Depths】 B100F conquered! You've cleared the game. Congratulations!\n" +
				"   The depths go on. Keep descending to push your record.\n" +
				"**********************************************",
		},
	}

	for _, tc := range tests {
		t.Run(tc.themeID, func(t *testing.T) {
			thJA := theme.GetPresetTheme(tc.themeID, model.LanguageJA)
			thEN := theme.GetPresetTheme(tc.themeID, model.LanguageEN)

			jaRes := jaMsg.AbyssCleared(thJA.Dungeons.Abyss, 100)
			if jaRes != tc.expectedJA {
				t.Errorf("\ngot : %s\nwant: %s", jaRes, tc.expectedJA)
			}

			enRes := enMsg.AbyssCleared(thEN.Dungeons.Abyss, 100)
			if enRes != tc.expectedEN {
				t.Errorf("\ngot : %s\nwant: %s", enRes, tc.expectedEN)
			}
		})
	}
}

func TestFarewell_And_HubLabels(t *testing.T) {
	jaMsg := GetMessages(model.LanguageJA)
	enMsg := GetMessages(model.LanguageEN)

	if jaMsg.Farewell != "ゲームを終了します。" {
		t.Errorf("unexpected ja Farewell: %s", jaMsg.Farewell)
	}
	if enMsg.Farewell != "Exiting the game." {
		t.Errorf("unexpected en Farewell: %s", enMsg.Farewell)
	}

	slotsJA := jaMsg.HubSlots("連携スロット", 3, "[追撃連携+] [弱点看破+] [立て直し+]")
	if slotsJA != "連携スロット [3/3]: [追撃連携+] [弱点看破+] [立て直し+]" {
		t.Errorf("unexpected slotsJA: %s", slotsJA)
	}

	storageJA := jaMsg.HubStorage("ロッカー", "訓練記録", 10, "未設定の連携スキル", "(なし)")
	if storageJA != "ロッカー: 訓練記録 x10 | 未設定の連携スキル: (なし)" {
		t.Errorf("unexpected storageJA: %s", storageJA)
	}
}

func TestRankUpgraded_Message(t *testing.T) {
	jaMsg := GetMessages(model.LanguageJA)
	enMsg := GetMessages(model.LanguageEN)

	jaRes := jaMsg.RankUpgraded("終焉を断つ刃")
	expectedJA := ">> ランクが上がった！ 【終焉を断つ刃】に進化した！"
	if jaRes != expectedJA {
		t.Errorf("got: %s, want: %s", jaRes, expectedJA)
	}

	enRes := enMsg.RankUpgraded("Blade of Armageddon")
	expectedEN := ">> Rank upgraded! Evolved to 【Blade of Armageddon】!"
	if enRes != expectedEN {
		t.Errorf("got: %s, want: %s", enRes, expectedEN)
	}
}

func TestTitle_ForgeLite(t *testing.T) {
	tests := []struct {
		name string
		lang model.Language
	}{
		{"JA", model.LanguageJA},
		{"EN", model.LanguageEN},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			title := GetMessages(tc.lang).Title
			if !strings.Contains(title, "Forge Lite") {
				t.Errorf("expected Title to contain 'Forge Lite', got %q", title)
			}
			if strings.Contains(title, "Prototype") {
				t.Errorf("expected Title to not contain 'Prototype', got %q", title)
			}
		})
	}
}

func TestAutoExploreMessages_JA(t *testing.T) {
	msg := GetMessages(model.LanguageJA)

	if msg.BattleCmdAutoExplore != "4: オート探索 (止まる所まで進む)" {
		t.Errorf("unexpected BattleCmdAutoExplore: %q", msg.BattleCmdAutoExplore)
	}

	if msg.BattleAutoExploreStart != "\n>> [オート探索開始] 止まる所まで、勝つたびに次の階層へ進みます..." {
		t.Errorf("unexpected BattleAutoExploreStart: %q", msg.BattleAutoExploreStart)
	}

	if msg.AutoExploreStoppedTitle != "\n>> オート探索を止めました：" {
		t.Errorf("unexpected AutoExploreStoppedTitle: %q", msg.AutoExploreStoppedTitle)
	}

	if got := msg.AutoExploreStopReasonClear(100); got != "   B100F を踏破しました（本編クリア）。" {
		t.Errorf("unexpected AutoExploreStopReasonClear: %q", got)
	}

	if got := msg.AutoExploreStopReasonRecord(150); got != "   これまでの最高到達階層 B150F を突破しました。ここから先は未踏の階層です。" {
		t.Errorf("unexpected AutoExploreStopReasonRecord: %q", got)
	}

	if got := msg.AutoExploreStopReasonMilestone(20); got != "   B20F を突破しました（10階ごとの区切り）。" {
		t.Errorf("unexpected AutoExploreStopReasonMilestone: %q", got)
	}

	expectedDanger := "   次の B58F の敵（攻撃 153）は、倒すまでに約 6 ターン、被害 約 765 の見込みです（現在 HP 900、危険域 300 以下）。"
	if got := msg.AutoExploreStopReasonDangerLine(58, 153, 6, 765, 900, 300); got != expectedDanger {
		t.Errorf("unexpected AutoExploreStopReasonDangerLine: \ngot : %q\nwant: %q", got, expectedDanger)
	}

	if msg.BattleInvalidCmd != ">> 無効なコマンドです。1, 2, 3, 4 のいずれかを入力してください。" {
		t.Errorf("unexpected BattleInvalidCmd: %q", msg.BattleInvalidCmd)
	}

	// Per-fight summary (JA)
	// With resources and one orb
	s1 := msg.AutoExploreBattleSummary(11, "闇のスライム", 2, 960, 1000, "強化の書", 6, []string{"猛毒の印"})
	exp1 := "  B11F 闇のスライム を撃破（2ターン） HP 960/1000 ／ 強化の書 +6 ／ オーブ 猛毒の印"
	if s1 != exp1 {
		t.Errorf("summary 1:\ngot : %q\nwant: %q", s1, exp1)
	}

	// No resources (0) and no orbs
	s2 := msg.AutoExploreBattleSummary(3, "コボルド", 1, 50, 50, "強化の書", 0, nil)
	exp2 := "  B3F コボルド を撃破（1ターン） HP 50/50"
	if s2 != exp2 {
		t.Errorf("summary 2:\ngot : %q\nwant: %q", s2, exp2)
	}

	// With resources and two orbs
	s3 := msg.AutoExploreBattleSummary(20, "氷のゴーレム", 3, 610, 1000, "強化の書", 10, []string{"猛毒の印", "吸血の印+"})
	exp3 := "  B20F 氷のゴーレム を撃破（3ターン） HP 610/1000 ／ 強化の書 +10 ／ オーブ 猛毒の印 ／ オーブ 吸血の印+"
	if s3 != exp3 {
		t.Errorf("summary 3:\ngot : %q\nwant: %q", s3, exp3)
	}
}

func TestAutoExploreMessages_EN(t *testing.T) {
	msg := GetMessages(model.LanguageEN)

	if msg.BattleCmdAutoExplore != "4: Auto Explore (keep going until a stop point)" {
		t.Errorf("unexpected BattleCmdAutoExplore: %q", msg.BattleCmdAutoExplore)
	}

	if msg.BattleAutoExploreStart != "\n>> [Auto Explore Started] Advancing to the next floor after each win until a stop point..." {
		t.Errorf("unexpected BattleAutoExploreStart: %q", msg.BattleAutoExploreStart)
	}

	if msg.AutoExploreStoppedTitle != "\n>> Auto Explore stopped:" {
		t.Errorf("unexpected AutoExploreStoppedTitle: %q", msg.AutoExploreStoppedTitle)
	}

	if got := msg.AutoExploreStopReasonClear(100); got != "   You conquered B100F (game cleared)." {
		t.Errorf("unexpected AutoExploreStopReasonClear: %q", got)
	}

	if got := msg.AutoExploreStopReasonRecord(150); got != "   You passed your previous record, B150F. The floors ahead are uncharted." {
		t.Errorf("unexpected AutoExploreStopReasonRecord: %q", got)
	}

	if got := msg.AutoExploreStopReasonMilestone(20); got != "   You cleared B20F (every 10 floors)." {
		t.Errorf("unexpected AutoExploreStopReasonMilestone: %q", got)
	}

	expectedDanger := "   The next enemy on B58F (ATK 153) would take about 6 turns to defeat, dealing about 765 damage (current HP 900, danger zone at 300 or below)."
	if got := msg.AutoExploreStopReasonDangerLine(58, 153, 6, 765, 900, 300); got != expectedDanger {
		t.Errorf("unexpected AutoExploreStopReasonDangerLine: \ngot : %q\nwant: %q", got, expectedDanger)
	}

	if msg.BattleInvalidCmd != ">> Invalid command. Please enter 1, 2, 3, or 4." {
		t.Errorf("unexpected BattleInvalidCmd: %q", msg.BattleInvalidCmd)
	}

	// Per-fight summary (EN)
	// 2 turns, resource + orb
	s1 := msg.AutoExploreBattleSummary(11, "Dark Slime", 2, 960, 1000, "Upgrade Scroll", 6, []string{"Deadly Poison"})
	exp1 := "  B11F Dark Slime defeated (2 turns) HP 960/1000 / Upgrade Scroll +6 / Orb: Deadly Poison"
	if s1 != exp1 {
		t.Errorf("EN summary 1:\ngot : %q\nwant: %q", s1, exp1)
	}

	// 1 turn, no resource, no orb
	s2 := msg.AutoExploreBattleSummary(1, "Goblin", 1, 50, 50, "Upgrade Scroll", 0, nil)
	exp2 := "  B1F Goblin defeated (1 turn) HP 50/50"
	if s2 != exp2 {
		t.Errorf("EN summary 2:\ngot : %q\nwant: %q", s2, exp2)
	}

	// 3 turns, resource + 2 orbs
	s3 := msg.AutoExploreBattleSummary(20, "Ice Golem", 3, 610, 1000, "Upgrade Scroll", 10, []string{"Deadly Poison", "Vampiric+"})
	exp3 := "  B20F Ice Golem defeated (3 turns) HP 610/1000 / Upgrade Scroll +10 / Orb: Deadly Poison / Orb: Vampiric+"
	if s3 != exp3 {
		t.Errorf("EN summary 3:\ngot : %q\nwant: %q", s3, exp3)
	}
}

func TestAutoExploreMessages_AllThemes(t *testing.T) {
	themes := []string{"classic_fantasy", "cyberpunk", "partner_sync"}
	for _, thID := range themes {
		t.Run(thID+"_JA", func(t *testing.T) {
			th := theme.GetPresetTheme(thID, model.LanguageJA)
			msg := GetMessages(model.LanguageJA)
			orbName := th.Orbs[model.OrbMultiHit].Name
			summary := msg.AutoExploreBattleSummary(15, "テスト敵", 2, 100, 100, th.ResourceName, 5, []string{orbName})
			if !strings.HasPrefix(summary, "  B15F テスト敵 を撃破（2ターン）") {
				t.Errorf("malformed prefix: %q", summary)
			}
			if !strings.Contains(summary, "／ "+th.ResourceName+" +5") {
				t.Errorf("missing resource %s: %q", th.ResourceName, summary)
			}
			if !strings.Contains(summary, "／ オーブ "+orbName) {
				t.Errorf("missing orb %s: %q", orbName, summary)
			}
		})

		t.Run(thID+"_EN", func(t *testing.T) {
			th := theme.GetPresetTheme(thID, model.LanguageEN)
			msg := GetMessages(model.LanguageEN)
			orbName := th.Orbs[model.OrbMultiHit].Name
			summary := msg.AutoExploreBattleSummary(15, "Test Enemy", 1, 100, 100, th.ResourceName, 5, []string{orbName})
			if !strings.HasPrefix(summary, "  B15F Test Enemy defeated (1 turn)") {
				t.Errorf("malformed prefix: %q", summary)
			}
			if !strings.Contains(summary, "/ "+th.ResourceName+" +5") {
				t.Errorf("missing resource %s: %q", th.ResourceName, summary)
			}
			if !strings.Contains(summary, "/ Orb: "+orbName) {
				t.Errorf("missing orb %s: %q", orbName, summary)
			}
		})
	}
}
