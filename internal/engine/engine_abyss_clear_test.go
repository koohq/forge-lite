package engine

import (
	"fmt"
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
	intPtr := func(v int) *int { return &v }
	boolPtr := func(v bool) *bool { return &v }

	makeSaveJSON := func(deepestFloor *int, abyssCleared *bool) string {
		parts := []string{
			`"weapon": {"name": "銅の剣", "baseAtk": 10, "plus": 0}`,
			`"stockScrolls": 0`,
			`"stockOrbs": []`,
			`"maxHp": 50`,
		}
		if deepestFloor != nil {
			parts = append(parts, fmt.Sprintf(`"deepestFloor": %d`, *deepestFloor))
		}
		if abyssCleared != nil {
			parts = append(parts, fmt.Sprintf(`"abyssCleared": %t`, *abyssCleared))
		}
		return "{\n  " + strings.Join(parts, ",\n  ") + "\n}"
	}

	tests := []struct {
		name             string
		deepestFloor     *int
		abyssCleared     *bool
		wantAbyssCleared bool
		wantDeepestFloor int
	}{
		{
			name:             "deepestFloor omitted, abyssCleared omitted -> false",
			deepestFloor:     nil,
			abyssCleared:     nil,
			wantAbyssCleared: false,
			wantDeepestFloor: 0,
		},
		{
			name:             "deepestFloor 0, abyssCleared omitted -> false",
			deepestFloor:     intPtr(0),
			abyssCleared:     nil,
			wantAbyssCleared: false,
			wantDeepestFloor: 0,
		},
		{
			name:             "deepestFloor 100, abyssCleared omitted -> false",
			deepestFloor:     intPtr(100),
			abyssCleared:     nil,
			wantAbyssCleared: false,
			wantDeepestFloor: 100,
		},
		{
			name:             "deepestFloor 100, abyssCleared false -> false",
			deepestFloor:     intPtr(100),
			abyssCleared:     boolPtr(false),
			wantAbyssCleared: false,
			wantDeepestFloor: 100,
		},
		{
			name:             "deepestFloor 100, abyssCleared true -> true",
			deepestFloor:     intPtr(100),
			abyssCleared:     boolPtr(true),
			wantAbyssCleared: true,
			wantDeepestFloor: 100,
		},
		{
			name:             "deepestFloor 101, abyssCleared omitted -> true",
			deepestFloor:     intPtr(101),
			abyssCleared:     nil,
			wantAbyssCleared: true,
			wantDeepestFloor: 101,
		},
		{
			name:             "deepestFloor 101, abyssCleared false -> true",
			deepestFloor:     intPtr(101),
			abyssCleared:     boolPtr(false),
			wantAbyssCleared: true,
			wantDeepestFloor: 101,
		},
		{
			name:             "deepestFloor 150, abyssCleared omitted -> true",
			deepestFloor:     intPtr(150),
			abyssCleared:     nil,
			wantAbyssCleared: true,
			wantDeepestFloor: 150,
		},
		{
			name:             "deepestFloor 150, abyssCleared true -> true",
			deepestFloor:     intPtr(150),
			abyssCleared:     boolPtr(true),
			wantAbyssCleared: true,
			wantDeepestFloor: 150,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rawJSON := makeSaveJSON(tc.deepestFloor, tc.abyssCleared)
			savePath := filepath.Join(t.TempDir(), "save.json")
			if err := os.WriteFile(savePath, []byte(rawJSON), 0644); err != nil {
				t.Fatal(err)
			}

			game := NewGame(strings.NewReader(""))
			game.saveFilePath = savePath
			if !game.loadGame() {
				t.Fatalf("expected loadGame to succeed")
			}

			if game.abyssCleared != tc.wantAbyssCleared {
				t.Errorf("abyssCleared mismatch: got %v, want %v", game.abyssCleared, tc.wantAbyssCleared)
			}
			if game.deepestFloor != tc.wantDeepestFloor {
				t.Errorf("deepestFloor mismatch: got %d, want %d", game.deepestFloor, tc.wantDeepestFloor)
			}

			// Acceptance criteria 4: Loading save data must not modify the file on disk.
			afterBytes, err := os.ReadFile(savePath)
			if err != nil {
				t.Fatalf("failed to read back save file: %v", err)
			}
			if string(afterBytes) != rawJSON {
				t.Errorf("save file was modified upon loadGame()")
			}
		})
	}

	// Acceptance criteria 3: Saving after loading a legacy save (deepestFloor 101, no abyssCleared)
	// writes "abyssCleared": true to the save file.
	t.Run("saveGame writes abyssCleared true after loading legacy 101F save", func(t *testing.T) {
		savePath := filepath.Join(t.TempDir(), "save.json")
		rawJSON := makeSaveJSON(intPtr(101), nil)
		if err := os.WriteFile(savePath, []byte(rawJSON), 0644); err != nil {
			t.Fatal(err)
		}

		game := NewGame(strings.NewReader(""))
		game.saveFilePath = savePath
		if !game.loadGame() {
			t.Fatalf("expected loadGame to succeed")
		}
		if !game.abyssCleared {
			t.Fatalf("expected abyssCleared to be true for 101F save, got false")
		}

		captureStdout(t, func() {
			game.saveGame()
		})

		savedData, err := save.LoadSaveData(savePath)
		if err != nil {
			t.Fatalf("failed to reload saved data: %v", err)
		}
		if savedData.AbyssCleared == nil || !*savedData.AbyssCleared {
			t.Errorf("expected saved AbyssCleared to be true, got %v", savedData.AbyssCleared)
		}
	})
}
