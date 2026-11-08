package engine

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koohq/forge-lite/internal/model"
	"github.com/koohq/forge-lite/internal/save"
)

func TestEvaluateAutoExploreStopReasons(t *testing.T) {
	// Table-driven check without I/O
	tests := []struct {
		name     string
		params   AutoExploreStopParams
		expected []AutoExploreStopReasonKind
	}{
		{
			name: "10階ごと (エンドレス)",
			params: AutoExploreStopParams{
				IsEndless:         true,
				ClearedFloor:      20,
				StartDeepestFloor: 50,
				StartFloor:        1,
				HasNextFloor:      true,
				NextFloor:         21,
				NextEnemyAtk:      10,
				EstTurns:          1,
				EstDmg:            0,
				CurPlayerHP:       100,
				MaxPlayerHP:       100,
			},
			expected: []AutoExploreStopReasonKind{StopReasonMilestone},
		},
		{
			name: "最高到達階層突破 (エンドレス、出撃時の階)",
			params: AutoExploreStopParams{
				IsEndless:         true,
				ClearedFloor:      15,
				StartDeepestFloor: 15,
				StartFloor:        10,
				HasNextFloor:      true,
				NextFloor:         16,
				NextEnemyAtk:      10,
				EstTurns:          1,
				EstDmg:            0,
				CurPlayerHP:       100,
				MaxPlayerHP:       100,
			},
			expected: []AutoExploreStopReasonKind{StopReasonRecordBreak},
		},
		{
			name: "最高到達が始めた階より浅い場合は記録で止まらない",
			params: AutoExploreStopParams{
				IsEndless:         true,
				ClearedFloor:      1,
				StartDeepestFloor: 0,
				StartFloor:        1,
				HasNextFloor:      true,
				NextFloor:         2,
				NextEnemyAtk:      10,
				EstTurns:          1,
				EstDmg:            0,
				CurPlayerHP:       100,
				MaxPlayerHP:       100,
			},
			expected: nil,
		},
		{
			name: "すでに記録で止まった後は記録で止まらない",
			params: AutoExploreStopParams{
				IsEndless:               true,
				ClearedFloor:            16,
				StartDeepestFloor:       15,
				StartFloor:              10,
				AlreadyStoppedForRecord: true,
				HasNextFloor:            true,
				NextFloor:               17,
				NextEnemyAtk:            10,
				EstTurns:                1,
				EstDmg:                  0,
				CurPlayerHP:             100,
				MaxPlayerHP:             100,
			},
			expected: nil,
		},
		{
			name: "B100F 初クリア (10階ごとと重ねずクリアのみ)",
			params: AutoExploreStopParams{
				IsEndless:         true,
				ClearedFloor:      100,
				JustClearedAbyss:  true,
				StartDeepestFloor: 100,
				StartFloor:        91,
				HasNextFloor:      true,
				NextFloor:         101,
				NextEnemyAtk:      10,
				EstTurns:          1,
				EstDmg:            0,
				CurPlayerHP:       100,
				MaxPlayerHP:       100,
			},
			// StartDeepestFloor==100 also triggers RecordBreak; check the order
			expected: []AutoExploreStopReasonKind{StopReasonClear, StopReasonRecordBreak},
		},
		{
			name: "B100F 初クリア (記録突破なし、クリアのみ)",
			params: AutoExploreStopParams{
				IsEndless:               true,
				ClearedFloor:            100,
				JustClearedAbyss:        true,
				StartDeepestFloor:       100,
				StartFloor:              91,
				AlreadyStoppedForRecord: true,
				HasNextFloor:            true,
				NextFloor:               101,
				NextEnemyAtk:            10,
				EstTurns:                1,
				EstDmg:                  0,
				CurPlayerHP:             100,
				MaxPlayerHP:             100,
			},
			expected: []AutoExploreStopReasonKind{StopReasonClear},
		},
		{
			name: "クリア済み B100F (クリアではなく10階ごとで止まる)",
			params: AutoExploreStopParams{
				IsEndless:               true,
				ClearedFloor:            100,
				JustClearedAbyss:        false,
				StartDeepestFloor:       105,
				StartFloor:              91,
				AlreadyStoppedForRecord: true,
				HasNextFloor:            true,
				NextFloor:               101,
				NextEnemyAtk:            10,
				EstTurns:                1,
				EstDmg:                  0,
				CurPlayerHP:             100,
				MaxPlayerHP:             100,
			},
			expected: []AutoExploreStopReasonKind{StopReasonMilestone},
		},
		{
			name: "安全の線で止まる",
			params: AutoExploreStopParams{
				IsEndless:    false,
				ClearedFloor: 2,
				HasNextFloor: true,
				NextFloor:    3,
				NextEnemyAtk: 153,
				EstTurns:     6,
				EstDmg:       765,
				CurPlayerHP:  900,
				MaxPlayerHP:  1000, // dangerLine = 300. 900 - 765 = 135 <= 300
			},
			expected: []AutoExploreStopReasonKind{StopReasonDangerLine},
		},
		{
			name: "複数理由の重なりと順序 (記録突破 + 10階ごと + 安全の線)",
			params: AutoExploreStopParams{
				IsEndless:         true,
				ClearedFloor:      20,
				StartDeepestFloor: 20,
				StartFloor:        11,
				HasNextFloor:      true,
				NextFloor:         21,
				NextEnemyAtk:      100,
				EstTurns:          5,
				EstDmg:            400,
				CurPlayerHP:       450,
				MaxPlayerHP:       1000, // dangerLine = 300. 450 - 400 = 50 <= 300
			},
			expected: []AutoExploreStopReasonKind{
				StopReasonRecordBreak,
				StopReasonMilestone,
				StopReasonDangerLine,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reasons := EvaluateAutoExploreStopReasons(tc.params)
			if len(reasons) != len(tc.expected) {
				t.Fatalf("expected %d reasons, got %d: %+v", len(tc.expected), len(reasons), reasons)
			}
			for i, exp := range tc.expected {
				if reasons[i].Kind != exp {
					t.Errorf("reason[%d] = %v; want %v", i, reasons[i].Kind, exp)
				}
			}
		})
	}
}

func TestAutoExplore_Endless_ClearB100(t *testing.T) {
	// Acceptance criteria 2, 4, 5.
	// Deepest 91, not cleared, start from 91F (10th option via "10\n").
	// "4" on 91F -> stop for the record after 91F -> "1".
	// "4" on 92F -> advance to 100F and stop after the clear screen -> "2" to retreat.
	savePath := filepath.Join(t.TempDir(), "save.json")

	// Create the initial save data
	maxHP := 1000
	deepest := 91
	cleared := false
	initData := model.SaveData{
		MaxHP:        &maxHP,
		Weapon:       model.Weapon{BaseAtk: 100000, Plus: 0, Slots: []model.OrbType{}},
		StockScrolls: 0,
		StockOrbs:    []model.OrbType{},
		DeepestFloor: &deepest,
		AbyssCleared: &cleared,
	}
	if err := save.SaveSaveData(savePath, &initData); err != nil {
		t.Fatalf("failed to init save: %v", err)
	}

	input := "10\n" + // start from 91F
		"4\n" + // 91F battle: auto explore
		"1\n" + // stop reason (record) -> proceed
		"4\n" + // 92F battle: auto explore
		"2\n" // B100F victory & stop reason (clear) -> retreat

	r := strings.NewReader(input)
	game := NewGame(r)
	game.saveFilePath = savePath
	game.language = model.LanguageJA
	game.player.MaxHP = 1000
	game.player.HP = 1000
	game.player.Weapon.BaseAtk = 100000 // one-shot every enemy
	game.player.Weapon.Slots = nil
	game.deepestFloor = 91
	game.abyssCleared = false

	out := captureStdout(t, func() {
		game.endlessDungeonPhase()
	})

	// 1. No "proceed" prompt between 92F and 99F;
	// it appears only after 91F and after 100F.
	nextFloorPromptCount := strings.Count(out, "1: 次の階層へ進む")
	if nextFloorPromptCount != 2 {
		t.Errorf("expected '1: 次の階層へ進む' exactly 2 times, got %d", nextFloorPromptCount)
	}

	// 2. Ten summary lines for 91F-100F
	for f := 91; f <= 100; f++ {
		prefix := fmt.Sprintf("  B%dF ", f)
		if !strings.Contains(out, prefix) {
			t.Errorf("missing summary line for floor B%dF", f)
		}
	}

	// 3. The clear screen (AbyssCleared) appears once
	clearCount := strings.Count(out, "本編クリア、おめでとうございます！")
	if clearCount != 1 {
		t.Errorf("expected AbyssCleared screen exactly once, got %d", clearCount)
	}

	// 4. Stop reasons
	// Stop after 91F: record passed
	if !strings.Contains(out, "これまでの最高到達階層 B91F を突破しました。") {
		t.Errorf("missing record break stop reason for B91F")
	}
	// Stop after 100F: clear reason only (no milestone reason)
	if !strings.Contains(out, "B100F を踏破しました（本編クリア）。") {
		t.Errorf("missing clear stop reason for B100F")
	}
	if strings.Contains(out, "B100F を突破しました（10階ごとの区切り）。") {
		t.Errorf("B100F should not show milestone reason on first clear")
	}

	// 5. No encounter banner or EndlessRecordUpdated for 93F-100F
	for f := 93; f <= 100; f++ {
		encounterMarker := fmt.Sprintf("B%dF :", f)
		if strings.Contains(out, encounterMarker) {
			t.Errorf("encounter message should not appear during auto explore for B%dF", f)
		}
	}
	// The record line appears once on 92F (entered manually) and not during auto explore (93F-100F)
	recordUpdatedCount := strings.Count(out, "最高到達階層を更新")
	if recordUpdatedCount != 1 {
		t.Errorf("EndlessRecordUpdated should only appear before auto explore on 92F, got %d", recordUpdatedCount)
	}

	// 6. Deepest floor is 100
	if game.deepestFloor != 100 {
		t.Errorf("expected deepestFloor 100, got %d", game.deepestFloor)
	}

	// 7. Saved AbyssCleared is true
	saved, err := save.LoadSaveData(savePath)
	if err != nil {
		t.Fatalf("failed to load save: %v", err)
	}
	if saved.AbyssCleared == nil || !*saved.AbyssCleared {
		t.Errorf("expected saved.AbyssCleared to be true")
	}
	if saved.DeepestFloor == nil || *saved.DeepestFloor != 100 {
		t.Errorf("expected saved.DeepestFloor to be 100, got %v", saved.DeepestFloor)
	}

	// 8. Retreating stores the loot
	if game.stockScrolls <= 0 {
		t.Errorf("expected stockScrolls > 0, got %d", game.stockScrolls)
	}
}

func TestAutoExplore_Endless_ClearedB100_Milestone(t *testing.T) {
	// Acceptance criteria 3, 5.
	// Already cleared, deepest 91, start from 91F. After the record stop, "1" then "4"
	// advances to B100F and stops for the milestone (no clear screen).
	savePath := filepath.Join(t.TempDir(), "save.json")

	input := "10\n" + // start from 91F
		"4\n" + // 91F: auto explore
		"1\n" + // stop at record break -> proceed
		"4\n" + // 92F: auto explore
		"2\n" // stop at B100F (milestone) -> retreat

	r := strings.NewReader(input)
	game := NewGame(r)
	game.saveFilePath = savePath
	game.language = model.LanguageJA
	game.player.MaxHP = 1000
	game.player.HP = 1000
	game.player.Weapon.BaseAtk = 100000
	game.deepestFloor = 91
	game.abyssCleared = true // already cleared

	out := captureStdout(t, func() {
		game.endlessDungeonPhase()
	})

	// No clear screen
	if strings.Contains(out, "本編クリア、おめでとうございます") {
		t.Errorf("cleared abyss should not show AbyssCleared screen")
	}

	// Stops for the milestone
	if !strings.Contains(out, "B100F を突破しました（10階ごとの区切り）。") {
		t.Errorf("missing milestone stop reason for B100F")
	}
}

func TestAutoExplore_DangerZoneStop(t *testing.T) {
	// Acceptance criteria 6, 7.
	// Start from 1F.
	// maxHP 100, curHP 40 -> dangerLine = 30
	// Weapon BaseAtk 72 one-shots 1F (HP 72), so HP stays at 40.
	// 2F enemy: HP 91, Atk 13.
	// Estimate: 2 turns, 13 damage -> 40 - 13 = 27 <= 30.
	// Auto explore stops before 2F because of the danger estimate.
	// After "1", the 2F enemy matches the one used for the estimate (HP/Atk).
	// Then "2" flees from the battle.
	savePath := filepath.Join(t.TempDir(), "save.json")

	input := "4\n" + // 1F battle: auto explore
		"1\n" + // stop reason (danger line) -> proceed to 2F
		"2\n" // 2F battle: retreat from battle

	r := strings.NewReader(input)
	game := NewGame(r)
	game.saveFilePath = savePath
	game.language = model.LanguageJA
	game.player.MaxHP = 100
	game.player.HP = 40
	game.player.Weapon.BaseAtk = 72
	game.deepestFloor = 0

	out := captureStdout(t, func() {
		game.endlessDungeonPhase()
	})

	// The stop reason shows the estimated numbers.
	expectedReason := "次の B2F の敵（攻撃 13）は、倒すまでに約 2 ターン、被害 約 13 の見込みです（現在 HP 40、危険域 30 以下）。"
	if !strings.Contains(out, expectedReason) {
		t.Errorf("expected danger line reason %q in output, got:\n%s", expectedReason, out)
	}

	// The command-3 summary must not appear before the stop
	if strings.Contains(out, "【オート戦闘 終了要約】") {
		t.Errorf("should not contain auto combat summary before stopping")
	}

	// After "1", the 2F battle shows the same enemy (HP 91)
	if !strings.Contains(out, "HP: 91/91") {
		t.Errorf("expected 2F enemy HP: 91/91 in encounter vs display, got:\n%s", out)
	}
}

func TestAutoExplore_FirstFight_DangerZoneStop(t *testing.T) {
	// Acceptance criterion 8.
	// The first fight after choosing "4" hits the danger zone (weak weapon, strong enemy).
	// The command-3 summary appears and the battle menu returns ("2" then flees).
	// No one-line summary is printed.
	savePath := filepath.Join(t.TempDir(), "save.json")

	// 1F: enemy HP 72, Atk 10.
	// Player: HP 35, maxHP 100 (dangerLine = 30).
	// Weapon: BaseAtk 5 (many turns needed).
	// After turn 1 the player has HP 25 (<= 30) and auto combat stops.
	input := "4\n" + // battle: choose 4
		"2\n" // return to command menu, choose 2 to retreat from combat

	r := strings.NewReader(input)
	game := NewGame(r)
	game.saveFilePath = savePath
	game.language = model.LanguageJA
	game.player.MaxHP = 100
	game.player.HP = 35
	game.player.Weapon.BaseAtk = 5
	game.deepestFloor = 0

	out := captureStdout(t, func() {
		game.endlessDungeonPhase()
	})

	// The same summary as command 3 appears
	if !strings.Contains(out, "【オート戦闘 終了要約】") {
		t.Errorf("expected auto combat summary in output")
	}
	if !strings.Contains(out, "危険域 (HP 25/100 <= 30%) に到達したため自動停止") {
		t.Errorf("expected danger zone stop reason in output")
	}

	// No one-line summary
	if strings.Contains(out, "を撃破") {
		t.Errorf("should not contain victory summary line")
	}
}

func TestAutoExplore_StarterDungeon_AllFloors(t *testing.T) {
	// Acceptance criteria 6, 9.
	// Pass a test DungeonDef that counts GetEnemy calls. With a strong weapon, "4" on 1F
	// reaches the last floor, shows DungeonComplete, stores the loot, and never shows the
	// "proceed" prompt. GetEnemy is called exactly once per floor.
	savePath := filepath.Join(t.TempDir(), "save.json")

	getEnemyCalls := make(map[int]int)
	customDungeon := model.DungeonDef{
		Name:   "テストの洞窟",
		Floors: 3,
		GetEnemy: func(floor int) (model.EnemyTemplate, error) {
			getEnemyCalls[floor]++
			return model.EnemyTemplate{
				Name: fmt.Sprintf("テスト敵%d", floor),
				HP:   20,
				Atk:  5,
				GetDrops: func() model.RunInventory {
					return model.RunInventory{Scrolls: 2, Orbs: []model.OrbType{}}
				},
			}, nil
		},
	}

	input := "4\n" // choose 4 on 1F
	r := strings.NewReader(input)
	game := NewGame(r)
	game.saveFilePath = savePath
	game.language = model.LanguageJA
	game.player.MaxHP = 100
	game.player.HP = 100
	game.player.Weapon.BaseAtk = 100 // one-shot every enemy

	out := captureStdout(t, func() {
		game.dungeonPhase(customDungeon)
	})

	// The dungeon complete screen appears
	if !strings.Contains(out, "完全踏破") {
		t.Errorf("expected dungeon complete message in output")
	}

	// The "proceed" prompt never appears
	if strings.Contains(out, "1: 次の階層へ進む") {
		t.Errorf("'1: 次の階層へ進む' should not appear during uninterrupted auto explore")
	}

	// Loot is stored (2 scrolls * 3 floors = 6)
	if game.stockScrolls != 6 {
		t.Errorf("expected 6 scrolls in stock, got %d", game.stockScrolls)
	}

	// GetEnemy is called once per floor
	for f := 1; f <= 3; f++ {
		if getEnemyCalls[f] != 1 {
			t.Errorf("GetEnemy for floor %d called %d times; want 1", f, getEnemyCalls[f])
		}
	}
}

func TestAutoCombat_Command3_BehaviorUnchanged(t *testing.T) {
	// Acceptance criterion 1.
	// Command 3 behaves as before (summary, loot lines, and the "proceed" prompt).
	savePath := filepath.Join(t.TempDir(), "save.json")

	input := "3\n" + // 1F battle: choose 3 (auto combat)
		"2\n" // between floors: retreat

	r := strings.NewReader(input)
	game := NewGame(r)
	game.saveFilePath = savePath
	game.language = model.LanguageJA
	game.player.MaxHP = 100
	game.player.HP = 100
	game.player.Weapon.BaseAtk = 100
	game.deepestFloor = 0

	out := captureStdout(t, func() {
		game.endlessDungeonPhase()
	})

	// The summary appears
	if !strings.Contains(out, "【オート戦闘 終了要約】") {
		t.Errorf("expected auto combat summary report")
	}
	if !strings.Contains(out, "- 終了理由: 敵を撃破！") {
		t.Errorf("expected enemy defeated reason")
	}

	// The loot lines appear
	if !strings.Contains(out, "戦利品獲得:") {
		t.Errorf("expected victory loot display")
	}

	// The "proceed" prompt appears
	if !strings.Contains(out, "1: 次の階層へ進む") {
		t.Errorf("expected '1: 次の階層へ進む' prompt")
	}
}
