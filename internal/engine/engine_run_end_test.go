package engine

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koohq/forge-lite/internal/model"
	"github.com/koohq/forge-lite/internal/save"
	"github.com/koohq/forge-lite/internal/theme"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = oldStdout
	out := <-outChan
	_ = r.Close()
	return out
}

func runDungeonRunEnd(t *testing.T, isEndless bool, isRetreat bool) (*Game, string) {
	t.Helper()

	var input string
	var hp int
	var baseAtk int
	if isRetreat {
		// 1F: attack and defeat enemy ("1\n")
		// Between floors: proceed to next floor ("1\n")
		// 2F: retreat during battle ("2\n")
		input = "1\n1\n2\n"
		hp = 100
		baseAtk = 100
	} else {
		// 1F: attack with 1 HP and weak weapon ("1\n")
		input = "1\n"
		hp = 1
		baseAtk = 1
	}

	game := NewGame(strings.NewReader(input))
	game.language = model.LanguageJA
	game.theme = theme.GetPresetTheme("classic_fantasy", model.LanguageJA)
	game.saveFilePath = filepath.Join(t.TempDir(), "save.json")
	game.stockScrolls = 5
	game.player.HP = hp
	game.player.MaxHP = hp
	game.player.Weapon.BaseAtk = baseAtk

	out := captureStdout(t, func() {
		if isEndless {
			game.endlessDungeonPhase()
		} else {
			dungeon := CreateStarterDungeon(game.theme, func(n int) int { return 0 })
			game.dungeonPhase(dungeon)
		}
	})

	return game, out
}

func TestDungeonRunEnd(t *testing.T) {
	tests := []struct {
		name            string
		isEndless       bool
		isRetreat       bool
		expectedDeepest int
	}{
		{
			name:      "Retreat_StarterDungeon",
			isEndless: false,
			isRetreat: true,
		},
		{
			name:            "Retreat_EndlessDungeon",
			isEndless:       true,
			isRetreat:       true,
			expectedDeepest: 2,
		},
		{
			name:      "Defeat_StarterDungeon",
			isEndless: false,
			isRetreat: false,
		},
		{
			name:      "Defeat_EndlessDungeon",
			isEndless: true,
			isRetreat: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			game, output := runDungeonRunEnd(t, tc.isEndless, tc.isRetreat)

			// In-memory warehouse scrolls check
			if game.stockScrolls != 5 {
				t.Errorf("expected stockScrolls 5, got %d", game.stockScrolls)
			}

			// Save file check
			saveData, err := save.LoadSaveData(game.saveFilePath)
			if err != nil {
				t.Fatalf("failed to load save data: %v", err)
			}
			if saveData.StockScrolls != 5 {
				t.Errorf("expected saved StockScrolls 5, got %d", saveData.StockScrolls)
			}
			if tc.expectedDeepest > 0 {
				if saveData.DeepestFloor == nil || *saveData.DeepestFloor != tc.expectedDeepest {
					t.Errorf("expected saved DeepestFloor %d, got %v", tc.expectedDeepest, saveData.DeepestFloor)
				}
			}

			var mustContain []string
			var mustNotContain []string
			if tc.isRetreat {
				mustContain = []string{
					game.msg().BattleRetreatCombat,
					game.msg().DungeonItemsLost(game.theme.ResourceName, 1, 0),
				}
				mustNotContain = []string{
					game.msg().DungeonDefeatPrompt,
					game.msg().DungeonCarriedBack,
				}
			} else {
				mustContain = []string{
					game.msg().DungeonDefeatPrompt,
					game.msg().DungeonItemsLost(game.theme.ResourceName, 0, 0),
					game.msg().DungeonCarriedBack,
				}
				mustNotContain = []string{
					game.msg().BattleRetreatCombat,
				}
			}

			for _, s := range mustContain {
				if !strings.Contains(output, s) {
					t.Errorf("expected output to contain %q, but it did not", s)
				}
			}
			for _, s := range mustNotContain {
				if strings.Contains(output, s) {
					t.Errorf("expected output to not contain %q, but it did", s)
				}
			}
		})
	}
}
