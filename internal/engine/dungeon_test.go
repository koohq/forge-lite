package engine

import (
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/koohq/forge-lite/internal/model"
	"github.com/koohq/forge-lite/internal/theme"
)

func TestCreateStarterDungeon(t *testing.T) {
	th := theme.GetPresetTheme("classic_fantasy", model.LanguageJA)
	dungeon := CreateStarterDungeon(th, func(n int) int { return 0 })

	if dungeon.Floors != 5 {
		t.Fatalf("expected 5 floors, got %d", dungeon.Floors)
	}

	boss, err := dungeon.GetEnemy(5)
	if err != nil {
		t.Fatalf("failed to get boss: %v", err)
	}
	if boss.HP != 130 || boss.Atk != 20 {
		t.Errorf("expected boss stats HP 130, Atk 20; got %d, %d", boss.HP, boss.Atk)
	}
	drops := boss.GetDrops()
	if drops.Scrolls != 3 || len(drops.Orbs) != 1 {
		t.Errorf("unexpected boss drops: %+v", drops)
	}
}

func TestGenerateEndlessEnemy_BossAndStats(t *testing.T) {
	th := theme.GetPresetTheme("cyberpunk", model.LanguageEN)

	// Floor 5 should be a mini-boss with plus orb
	enemy := GenerateEndlessEnemy(5, th, func(n int) int { return 0 })
	if enemy.HP <= 60 {
		t.Errorf("expected scaled HP, got %d", enemy.HP)
	}
	drops := enemy.GetDrops()
	if len(drops.Orbs) != 1 || drops.Orbs[0] != model.OrbMultiHitPlus {
		t.Errorf("expected 1 plus orb, got %+v", drops.Orbs)
	}

	// Floor 1 should not be boss and have 0 orbs
	normalEnemy := GenerateEndlessEnemy(1, th, func(n int) int { return 0 })
	normalDrops := normalEnemy.GetDrops()
	if len(normalDrops.Orbs) != 0 {
		t.Errorf("expected 0 orbs on floor 1, got %+v", normalDrops.Orbs)
	}
}

func TestPartnerSyncEnemyNames(t *testing.T) {
	th := theme.GetPresetTheme("partner_sync", model.LanguageJA)

	expectedPrefixes := []string{
		"強襲型",
		"重装型",
		"高機動型",
		"近接特化型",
		"深層配備の",
		"封鎖区画の",
		"廃墟に残る",
		"異常活性化した",
	}

	if len(th.EnemyPrefixes) != len(expectedPrefixes) {
		t.Fatalf("expected %d prefixes, got %d", len(expectedPrefixes), len(th.EnemyPrefixes))
	}
	for i, p := range expectedPrefixes {
		if th.EnemyPrefixes[i] != p {
			t.Errorf("prefix %d: expected %s, got %s", i, p, th.EnemyPrefixes[i])
		}
	}

	// Verify 10 random sample generations
	r := rand.New(rand.NewPCG(1234, 5678))
	for i := 1; i <= 10; i++ {
		floor := i*3 + 1
		enemy := GenerateEndlessEnemy(floor, th, func(n int) int { return r.IntN(n) })
		t.Logf("Sample %2d (B%dF): %s", i, floor, enemy.Name)
	}
}

func assertValidEnemyName(t *testing.T, name string) {
	t.Helper()
	if strings.Contains(name, "  ") {
		t.Errorf("name contains double spaces: %q", name)
	}
	if strings.Contains(name, "- ") {
		t.Errorf("name contains hyphen followed by space: %q", name)
	}
	if strings.TrimSpace(name) != name {
		t.Errorf("name has leading/trailing spaces: %q", name)
	}
}

// pickRandInt returns a randInt stub that yields picks in order, then 0.
func pickRandInt(picks ...int) func(n int) int {
	call := 0
	return func(_ int) int {
		call++
		if call <= len(picks) {
			return picks[call-1]
		}
		return 0
	}
}

// deepVariantRandInt returns a randInt stub that picks the given variant
// for the two-way enemy choice of the deep dungeon, and 0 otherwise.
func deepVariantRandInt(variant int) func(n int) int {
	return func(n int) int {
		if n == 2 {
			return variant
		}
		return 0
	}
}

func TestEnemyNames_AllCombinations(t *testing.T) {
	themeIDs := []string{"classic_fantasy", "cyberpunk", "partner_sync"}
	languages := []model.Language{model.LanguageJA, model.LanguageEN}

	for _, themeID := range themeIDs {
		for _, lang := range languages {
			th := theme.GetPresetTheme(themeID, lang)

			// 1. Endless Dungeon: all combinations of prefix and base
			for i, prefix := range th.EnemyPrefixes {
				for j, base := range th.EnemyBases {
					for _, floor := range []int{1, 5} {
						enemy := GenerateEndlessEnemy(floor, th, pickRandInt(i, j))

						assertValidEnemyName(t, enemy.Name)

						if lang == model.LanguageJA {
							expected := prefix + base
							if floor == 5 {
								expected += " (中ボス)"
							}
							if enemy.Name != expected {
								t.Errorf("[%s/JA Endless] expected %q, got %q", themeID, expected, enemy.Name)
							}
						}
					}
				}
			}

			// 2. Deep Dungeon: floors 1 to 10
			for _, variant := range []int{0, 1} {
				dungeon := CreateDeepDungeon(th, "Normal", deepVariantRandInt(variant))
				for floor := 1; floor <= 10; floor++ {
					enemy, err := dungeon.GetEnemy(floor)
					if err != nil {
						t.Fatalf("[%s/%s Deep B%dF] unexpected error: %v", themeID, lang, floor, err)
					}
					assertValidEnemyName(t, enemy.Name)
				}
			}
		}
	}
}

func TestEnemyNames_EnglishExamples(t *testing.T) {
	endlessName := func(th model.Theme, floor, prefixIdx, baseIdx int) string {
		return GenerateEndlessEnemy(floor, th, pickRandInt(prefixIdx, baseIdx)).Name
	}

	// Cyberpunk Endless
	cpEN := theme.GetPresetTheme("cyberpunk", model.LanguageEN)
	// "Rogue " (idx 1), "Firewall Node" (idx 3)
	if name := endlessName(cpEN, 1, 1, 3); name != "Rogue Firewall Node" {
		t.Errorf("expected 'Rogue Firewall Node', got %q", name)
	}
	// "Glitch-" (idx 0), "AI Core" (idx 7)
	if name := endlessName(cpEN, 1, 0, 7); name != "Glitch-AI Core" {
		t.Errorf("expected 'Glitch-AI Core', got %q", name)
	}
	// "Glitch-" (idx 0), "Security Drone" (idx 0), floor 5 (Mini-Boss)
	if name := endlessName(cpEN, 5, 0, 0); name != "Glitch-Security Drone (Mini-Boss)" {
		t.Errorf("expected 'Glitch-Security Drone (Mini-Boss)', got %q", name)
	}

	// Partner Sync Endless
	psEN := theme.GetPresetTheme("partner_sync", model.LanguageEN)
	// "Quarantine-Zone " (idx 5), "Mutant Beast" (idx 1)
	if name := endlessName(psEN, 1, 5, 1); name != "Quarantine-Zone Mutant Beast" {
		t.Errorf("expected 'Quarantine-Zone Mutant Beast', got %q", name)
	}
	// "Heavy Armored " (idx 1), "Drone" (idx 0)
	if name := endlessName(psEN, 1, 1, 0); name != "Heavy Armored Drone" {
		t.Errorf("expected 'Heavy Armored Drone', got %q", name)
	}

	// Classic Fantasy Endless
	cfEN := theme.GetPresetTheme("classic_fantasy", model.LanguageEN)
	// "Fierce" (idx 0), "Slime" (idx 0)
	if name := endlessName(cfEN, 1, 0, 0); name != "Fierce Slime" {
		t.Errorf("expected 'Fierce Slime', got %q", name)
	}

	// Deep Dungeon Examples
	cpDeep := CreateDeepDungeon(cpEN, "Normal", deepVariantRandInt(0))
	if e, _ := cpDeep.GetEnemy(1); e.Name != "Glitch-Combat Avatar" {
		t.Errorf("expected 'Glitch-Combat Avatar', got %q", e.Name)
	}
	cpDeep2 := CreateDeepDungeon(cpEN, "Normal", deepVariantRandInt(1))
	if e, _ := cpDeep2.GetEnemy(1); e.Name != "Rogue Mainframe" {
		t.Errorf("expected 'Rogue Mainframe', got %q", e.Name)
	}

	psDeep := CreateDeepDungeon(psEN, "Normal", deepVariantRandInt(0))
	if e, _ := psDeep.GetEnemy(10); e.Name != "Ruins-Lurking Command Core (BOSS)" {
		t.Errorf("expected 'Ruins-Lurking Command Core (BOSS)', got %q", e.Name)
	}

	cfDeep := CreateDeepDungeon(cfEN, "Normal", deepVariantRandInt(0))
	if e, _ := cfDeep.GetEnemy(10); e.Name != "Mad Death Knight (BOSS)" {
		t.Errorf("expected 'Mad Death Knight (BOSS)', got %q", e.Name)
	}
}

func TestEnemyNames_JapaneseExamples(t *testing.T) {
	cpJA := theme.GetPresetTheme("cyberpunk", model.LanguageJA)
	cpDeep1 := CreateDeepDungeon(cpJA, "Normal", deepVariantRandInt(0))
	if e, _ := cpDeep1.GetEnemy(1); e.Name != "グリッチ・戦闘アバター" {
		t.Errorf("expected 'グリッチ・戦闘アバター', got %q", e.Name)
	}
	cpDeep2 := CreateDeepDungeon(cpJA, "Normal", deepVariantRandInt(1))
	if e, _ := cpDeep2.GetEnemy(1); e.Name != "ローグ・メインフレーム" {
		t.Errorf("expected 'ローグ・メインフレーム', got %q", e.Name)
	}

	psJA := theme.GetPresetTheme("partner_sync", model.LanguageJA)
	psDeep2 := CreateDeepDungeon(psJA, "Normal", deepVariantRandInt(1))
	if e, _ := psDeep2.GetEnemy(1); e.Name != "重装型掃討機兵" {
		t.Errorf("expected '重装型掃討機兵', got %q", e.Name)
	}

	cfJA := theme.GetPresetTheme("classic_fantasy", model.LanguageJA)
	cfDeep := CreateDeepDungeon(cfJA, "Normal", deepVariantRandInt(0))
	if e, _ := cfDeep.GetEnemy(10); e.Name != "狂気の死霊騎士 (BOSS)" {
		t.Errorf("expected '狂気の死霊騎士 (BOSS)', got %q", e.Name)
	}
}

func TestJoinEnemyName(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		base     string
		lang     model.Language
		expected string
	}{
		// Japanese cases
		{
			name:     "JA standard concatenation without separator",
			prefix:   "凶暴な",
			base:     "スライム",
			lang:     model.LanguageJA,
			expected: "凶暴なスライム",
		},
		{
			name:     "JA with hyphen-like katakana middle dot",
			prefix:   "グリッチ・",
			base:     "防壁ノード",
			lang:     model.LanguageJA,
			expected: "グリッチ・防壁ノード",
		},
		{
			name:     "JA empty prefix returns base",
			prefix:   "",
			base:     "スライム",
			lang:     model.LanguageJA,
			expected: "スライム",
		},
		{
			name:     "JA empty base returns prefix",
			prefix:   "凶暴な",
			base:     "",
			lang:     model.LanguageJA,
			expected: "凶暴な",
		},
		{
			name:     "JA both empty",
			prefix:   "",
			base:     "",
			lang:     model.LanguageJA,
			expected: "",
		},
		{
			name:     "JA trims outer spaces",
			prefix:   " 凶暴な",
			base:     "スライム ",
			lang:     model.LanguageJA,
			expected: "凶暴なスライム",
		},

		// English cases
		{
			name:     "EN no trailing space in prefix (Classic Fantasy style)",
			prefix:   "Fierce",
			base:     "Slime",
			lang:     model.LanguageEN,
			expected: "Fierce Slime",
		},
		{
			name:     "EN trailing space in prefix (Cyberpunk Rogue style)",
			prefix:   "Rogue ",
			base:     "Firewall Node",
			lang:     model.LanguageEN,
			expected: "Rogue Firewall Node",
		},
		{
			name:     "EN ending with hyphen (Cyberpunk Glitch- style)",
			prefix:   "Glitch-",
			base:     "AI Core",
			lang:     model.LanguageEN,
			expected: "Glitch-AI Core",
		},
		{
			name:     "EN hyphen ending with trailing space",
			prefix:   "Glitch- ",
			base:     "AI Core",
			lang:     model.LanguageEN,
			expected: "Glitch-AI Core",
		},
		{
			name:     "EN multiple trailing spaces in prefix",
			prefix:   "Overclocked   ",
			base:     "Combat Avatar",
			lang:     model.LanguageEN,
			expected: "Overclocked Combat Avatar",
		},
		{
			name:     "EN leading and trailing spaces on both prefix and base",
			prefix:   "  Assault  ",
			base:     "  Drone  ",
			lang:     model.LanguageEN,
			expected: "Assault Drone",
		},
		{
			name:     "EN internal space preserved in base",
			prefix:   "Quarantine-Zone ",
			base:     "Mutant Beast",
			lang:     model.LanguageEN,
			expected: "Quarantine-Zone Mutant Beast",
		},
		{
			name:     "EN internal space preserved in prefix",
			prefix:   "Heavy Armored ",
			base:     "Sweeper Mech",
			lang:     model.LanguageEN,
			expected: "Heavy Armored Sweeper Mech",
		},
		{
			name:     "EN acceptance criterion 4: Void without trailing space",
			prefix:   "Void",
			base:     "Beast",
			lang:     model.LanguageEN,
			expected: "Void Beast",
		},
		{
			name:     "EN acceptance criterion 4: Void with trailing space",
			prefix:   "Void ",
			base:     "Beast",
			lang:     model.LanguageEN,
			expected: "Void Beast",
		},
		{
			name:     "EN acceptance criterion 4: empty prefix returns base only",
			prefix:   "",
			base:     "Enemy",
			lang:     model.LanguageEN,
			expected: "Enemy",
		},
		{
			name:     "EN whitespace-only prefix returns base only",
			prefix:   "   ",
			base:     "Enemy",
			lang:     model.LanguageEN,
			expected: "Enemy",
		},
		{
			name:     "EN empty base returns prefix only",
			prefix:   "Dark",
			base:     "",
			lang:     model.LanguageEN,
			expected: "Dark",
		},
		{
			name:     "EN whitespace-only base returns prefix only",
			prefix:   "Dark",
			base:     "   ",
			lang:     model.LanguageEN,
			expected: "Dark",
		},
		{
			name:     "EN both empty",
			prefix:   "",
			base:     "",
			lang:     model.LanguageEN,
			expected: "",
		},
		{
			name:     "EN both whitespace only",
			prefix:   "  ",
			base:     "   ",
			lang:     model.LanguageEN,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := joinEnemyName(tc.prefix, tc.base, tc.lang)
			if got != tc.expected {
				t.Errorf("joinEnemyName(%q, %q, %s) = %q; want %q", tc.prefix, tc.base, tc.lang, got, tc.expected)
			}
		})
	}
}
