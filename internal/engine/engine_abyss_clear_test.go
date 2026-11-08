package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koohq/forge-lite/internal/model"
	"github.com/koohq/forge-lite/internal/save"
	"github.com/koohq/forge-lite/internal/theme"
)

func TestEndlessAbyssClear_Cases(t *testing.T) {
	// Helper to build 91F start input.
	// 91F selection index is 10 ("10\n").
	buildInput := func(steps []string) string {
		var sb strings.Builder
		sb.WriteString("10\n")
		for _, s := range steps {
			sb.WriteString(s)
		}
		return sb.String()
	}

	// 1: Clear 100F, defeat 101F, then retreat between floors
	case1Steps := make([]string, 0, 22)
	for f := 91; f <= 100; f++ {
		case1Steps = append(case1Steps, "1\n", "1\n") // attack, proceed
	}
	case1Steps = append(case1Steps, "1\n", "2\n") // 101F: attack, retreat

	// 2: Defeat 99F, then retreat between floors
	case2Steps := make([]string, 0, 20)
	for f := 91; f <= 98; f++ {
		case2Steps = append(case2Steps, "1\n", "1\n") // attack, proceed
	}
	case2Steps = append(case2Steps, "1\n", "2\n") // 99F: attack, retreat

	// 3: Retreat during battle at 100F
	case3Steps := make([]string, 0, 20)
	for f := 91; f <= 99; f++ {
		case3Steps = append(case3Steps, "1\n", "1\n") // attack, proceed
	}
	case3Steps = append(case3Steps, "2\n") // 100F: retreat in battle

	// 4: Already cleared, defeat 100F, then retreat between floors
	case4Steps := make([]string, 0, 22)
	for f := 91; f <= 99; f++ {
		case4Steps = append(case4Steps, "1\n", "1\n") // attack, proceed
	}
	case4Steps = append(case4Steps, "1\n", "2\n") // 100F: attack, retreat

	// 5: Defeat 100F, proceed to 101F, retreat during battle at 101F
	case5Steps := make([]string, 0, 22)
	for f := 91; f <= 100; f++ {
		case5Steps = append(case5Steps, "1\n", "1\n") // attack, proceed
	}
	case5Steps = append(case5Steps, "2\n") // 101F: retreat in battle

	tests := []struct {
		name                 string
		initialCleared       bool
		input                string
		wantClearBanner      bool
		wantBannerCount      int
		wantAbyssCleared     bool
		wantSaveAbyssCleared bool
		wantDeepestFloor     int
		checkBannerOrder     bool
		checkStockPreserved  bool
	}{
		{
			name:                 "case 1: clear 100F and 101F then retreat",
			initialCleared:       false,
			input:                buildInput(case1Steps),
			wantClearBanner:      true,
			wantBannerCount:      1,
			wantAbyssCleared:     true,
			wantSaveAbyssCleared: true,
			wantDeepestFloor:     101,
			checkBannerOrder:     true,
		},
		{
			name:                 "case 2: defeat 99F and retreat",
			initialCleared:       false,
			input:                buildInput(case2Steps),
			wantClearBanner:      false,
			wantAbyssCleared:     false,
			wantSaveAbyssCleared: false,
			wantDeepestFloor:     99,
		},
		{
			name:                 "case 3: retreat during battle at 100F",
			initialCleared:       false,
			input:                buildInput(case3Steps),
			wantClearBanner:      false,
			wantAbyssCleared:     false,
			wantSaveAbyssCleared: false,
			wantDeepestFloor:     100,
		},
		{
			name:                 "case 4: already cleared, defeat 100F and retreat",
			initialCleared:       true,
			input:                buildInput(case4Steps),
			wantClearBanner:      false,
			wantAbyssCleared:     true,
			wantSaveAbyssCleared: true,
			wantDeepestFloor:     100,
		},
		{
			name:                 "case 5: defeat 100F then retreat during 101F battle (loot lost)",
			initialCleared:       false,
			input:                buildInput(case5Steps),
			wantClearBanner:      true,
			wantBannerCount:      1,
			wantAbyssCleared:     true,
			wantSaveAbyssCleared: true,
			wantDeepestFloor:     101,
			checkStockPreserved:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			savePath := filepath.Join(t.TempDir(), "save.json")
			game := NewGame(strings.NewReader(tc.input))
			game.language = model.LanguageJA
			game.theme = theme.GetPresetTheme(theme.DefaultThemeID, model.LanguageJA)
			game.saveFilePath = savePath
			game.deepestFloor = 91
			game.abyssCleared = tc.initialCleared
			game.player.MaxHP = 10000
			game.player.HP = 10000
			game.player.Weapon.BaseAtk = 100000
			game.stockScrolls = 5

			output := captureStdout(t, func() {
				game.endlessDungeonPhase()
			})

			clearBannerToken := "B100F 踏破！ 本編クリア、おめでとうございます！"
			if tc.wantClearBanner {
				count := strings.Count(output, clearBannerToken)
				if count != tc.wantBannerCount {
					t.Errorf("banner count mismatch: got %d, want %d", count, tc.wantBannerCount)
				}
			} else {
				if strings.Contains(output, clearBannerToken) {
					t.Errorf("unexpected clear banner in output")
				}
			}

			if tc.checkBannerOrder {
				idxBanner := strings.Index(output, clearBannerToken)
				if idxBanner == -1 {
					t.Errorf("clear banner not found in output")
				} else {
					savedToken := game.msg().SavedSuccess
					idxSaved := strings.Index(output[idxBanner:], savedToken)
					if idxSaved == -1 {
						t.Errorf("save message not found after clear banner")
					} else {
						idxSaved += idxBanner
						nextFloorToken := game.msg().DungeonNextFloor
						idxNext := strings.Index(output[idxSaved:], nextFloorToken)
						if idxNext == -1 {
							t.Errorf("next floor prompt not found after save message")
						}
					}
				}
			}

			if game.abyssCleared != tc.wantAbyssCleared {
				t.Errorf("game.abyssCleared mismatch: got %v, want %v", game.abyssCleared, tc.wantAbyssCleared)
			}
			if game.deepestFloor != tc.wantDeepestFloor {
				t.Errorf("game.deepestFloor mismatch: got %d, want %d", game.deepestFloor, tc.wantDeepestFloor)
			}

			loadedSave, err := save.LoadSaveData(savePath)
			if err != nil {
				t.Fatalf("failed to load saved data: %v", err)
			}
			if loadedSave.AbyssCleared == nil || *loadedSave.AbyssCleared != tc.wantSaveAbyssCleared {
				t.Errorf("loaded save AbyssCleared mismatch: got %+v, want %v", loadedSave.AbyssCleared, tc.wantSaveAbyssCleared)
			}
			if *loadedSave.DeepestFloor != tc.wantDeepestFloor {
				t.Errorf("loaded save DeepestFloor mismatch: got %d, want %d", *loadedSave.DeepestFloor, tc.wantDeepestFloor)
			}

			if tc.checkStockPreserved {
				if loadedSave.StockScrolls != 5 {
					t.Errorf("expected stockScrolls to be preserved as 5, got %d", loadedSave.StockScrolls)
				}
			}
		})
	}
}

func TestHubPhase_AbyssClearedMark(t *testing.T) {
	tests := []struct {
		name         string
		lang         model.Language
		abyssCleared bool
		wantMark     string
		shouldHave   bool
	}{
		{"JA not cleared", model.LanguageJA, false, "★本編クリア", false},
		{"JA cleared", model.LanguageJA, true, "★本編クリア", true},
		{"EN not cleared", model.LanguageEN, false, "★Game Cleared", false},
		{"EN cleared", model.LanguageEN, true, "★Game Cleared", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			game := NewGame(strings.NewReader("0\n"))
			game.language = tc.lang
			game.abyssCleared = tc.abyssCleared
			game.deepestFloor = 123

			output := captureStdout(t, func() {
				game.hubPhase()
			})

			hasMark := strings.Contains(output, tc.wantMark)
			if hasMark != tc.shouldHave {
				t.Errorf("mark %q presence got %v, want %v\nOutput:\n%s", tc.wantMark, hasMark, tc.shouldHave, output)
			}
		})
	}
}

func TestLoadGame_LegacySaveWithoutAbyssCleared(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Legacy save: deepestFloor 150, no abyssCleared field
	legacyJSON := `{
		"weapon": {"name": "銅の剣", "baseAtk": 10, "plus": 0},
		"stockScrolls": 0,
		"stockOrbs": [],
		"maxHp": 50,
		"deepestFloor": 150
	}`
	legacyPath := filepath.Join(tmpDir, "legacy.json")
	if err := os.WriteFile(legacyPath, []byte(legacyJSON), 0644); err != nil {
		t.Fatal(err)
	}

	game := NewGame(strings.NewReader(""))
	game.saveFilePath = legacyPath
	if !game.loadGame() {
		t.Fatalf("expected loadGame to succeed for legacy save")
	}
	if game.abyssCleared {
		t.Errorf("expected abyssCleared to be false for legacy save, got true")
	}
	if game.deepestFloor != 150 {
		t.Errorf("expected deepestFloor 150, got %d", game.deepestFloor)
	}

	// 2. Save with abyssCleared: true
	clearedJSON := `{
		"weapon": {"name": "銅の剣", "baseAtk": 10, "plus": 0},
		"stockScrolls": 0,
		"stockOrbs": [],
		"maxHp": 50,
		"deepestFloor": 150,
		"abyssCleared": true
	}`
	clearedPath := filepath.Join(tmpDir, "cleared.json")
	if err := os.WriteFile(clearedPath, []byte(clearedJSON), 0644); err != nil {
		t.Fatal(err)
	}

	game2 := NewGame(strings.NewReader(""))
	game2.saveFilePath = clearedPath
	if !game2.loadGame() {
		t.Fatalf("expected loadGame to succeed for cleared save")
	}
	if !game2.abyssCleared {
		t.Errorf("expected abyssCleared to be true, got false")
	}
}
