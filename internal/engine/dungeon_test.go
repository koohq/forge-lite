package engine

import (
	"math/rand/v2"
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
