package engine

import (
	"testing"

	"github.com/koohq/forge-lite/internal/i18n"
	"github.com/koohq/forge-lite/internal/model"
)

func TestCalcWeaponAtk(t *testing.T) {
	w := model.Weapon{
		BaseAtk: 10,
		Plus:    5,
	}
	atk := CalcWeaponAtk(w)
	if atk != 20 {
		t.Errorf("expected 20, got %d", atk)
	}
}

func TestResolvePlayerAttack_Standard(t *testing.T) {
	w := model.Weapon{
		BaseAtk: 10,
		Plus:    0,
	}
	enemy := &model.Enemy{HP: 100, MaxHP: 100, Atk: 5}

	hits, nextHP := ResolvePlayerAttack(w, 50, 50, enemy, func() float64 { return 1.0 })
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(hits))
	}
	if hits[0].Damage != 10 {
		t.Errorf("expected 10 damage, got %d", hits[0].Damage)
	}
	if enemy.HP != 90 {
		t.Errorf("expected 90 enemy HP, got %d", enemy.HP)
	}
	if nextHP != 50 {
		t.Errorf("expected 50 player HP, got %d", nextHP)
	}
}

func TestResolvePlayerAttack_MultiHitAndCrit(t *testing.T) {
	w := model.Weapon{
		BaseAtk: 10,
		Plus:    0,
		Slots:   []model.OrbType{model.OrbMultiHitPlus, model.OrbCriticalPlus},
	}
	enemy := &model.Enemy{HP: 100, MaxHP: 100, Atk: 5}

	// randFloat returning 0.1 will trigger crit (0.1 < 0.4)
	hits, _ := ResolvePlayerAttack(w, 50, 50, enemy, func() float64 { return 0.1 })
	if len(hits) != 2 {
		t.Fatalf("expected 2 hits, got %d", len(hits))
	}
	// MultiHitPlus rate is 0.75 -> floor(10 * 0.75) = 7. Crit doubles to 14.
	if hits[0].Damage != 14 || !hits[0].IsCrit {
		t.Errorf("expected 14 crit damage, got %d (crit=%v)", hits[0].Damage, hits[0].IsCrit)
	}
	if hits[1].Damage != 14 || !hits[1].IsCrit {
		t.Errorf("expected 14 crit damage, got %d (crit=%v)", hits[1].Damage, hits[1].IsCrit)
	}
	if enemy.HP != 72 {
		t.Errorf("expected 72 enemy HP, got %d", enemy.HP)
	}
}

func TestResolvePlayerAttack_VampAndPoison(t *testing.T) {
	w := model.Weapon{
		BaseAtk: 100,
		Plus:    0,
		Slots:   []model.OrbType{model.OrbVampPlus, model.OrbPoisonPlus},
	}
	enemy := &model.Enemy{HP: 500, MaxHP: 500, Atk: 5}

	hits, nextHP := ResolvePlayerAttack(w, 10, 100, enemy, func() float64 { return 0.9 })
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(hits))
	}
	// VampPlus is 25% of 100 dmg = 25 heal. Player HP: 10 + 25 = 35.
	if hits[0].Heal != 25 {
		t.Errorf("expected 25 heal, got %d", hits[0].Heal)
	}
	if nextHP != 35 {
		t.Errorf("expected 35 player HP, got %d", nextHP)
	}
	// PoisonPlus adds 2 poison
	if enemy.Poison != 2 {
		t.Errorf("expected 2 poison, got %d", enemy.Poison)
	}
}

func TestRunAutoCombat_DangerThreshold(t *testing.T) {
	// Player maxHP 100, danger threshold is <= 30.
	// Player starts with 35 HP. Enemy attacks for 10.
	// After 1 turn, player HP becomes 25 (<= 30), so auto combat should stop.
	w := model.Weapon{BaseAtk: 5, Plus: 0}
	enemy := &model.Enemy{HP: 100, MaxHP: 100, Atk: 10}
	msgs := i18n.GetMessages(model.LanguageJA)

	res := RunAutoCombat(w, 35, 100, enemy, msgs, func() float64 { return 0.9 })
	if !res.DangerZoneStopped {
		t.Errorf("expected danger zone stop, got result %+v", res)
	}
	if res.FinalPlayerHP != 25 {
		t.Errorf("expected 25 final player HP, got %d", res.FinalPlayerHP)
	}
}

func TestRunAutoCombat_EnemyDefeated(t *testing.T) {
	w := model.Weapon{BaseAtk: 50, Plus: 0}
	enemy := &model.Enemy{HP: 40, MaxHP: 40, Atk: 5}
	msgs := i18n.GetMessages(model.LanguageEN)

	res := RunAutoCombat(w, 50, 50, enemy, msgs, func() float64 { return 0.9 })
	if !res.EnemyDefeated {
		t.Errorf("expected enemy defeat, got result %+v", res)
	}
	if res.Turns != 1 {
		t.Errorf("expected 1 turn, got %d", res.Turns)
	}
}

func TestEstimateAutoFight(t *testing.T) {
	// 1. No multi-hit: atk 10, enemy HP 25, atk 7 -> 3 turns, 14 damage
	wStandard := model.Weapon{BaseAtk: 10}
	turns, dmg := EstimateAutoFight(wStandard, 25, 7)
	if turns != 3 || dmg != 14 {
		t.Errorf("expected turns=3, dmg=14, got turns=%d, dmg=%d", turns, dmg)
	}

	// 2. Multi-hit: int(10*0.65)=6 per hit, 12 per turn. Enemy HP 25 -> 3 turns, 14 damage
	wMulti := model.Weapon{BaseAtk: 10, Slots: []model.OrbType{model.OrbMultiHit}}
	turns, dmg = EstimateAutoFight(wMulti, 25, 7)
	if turns != 3 || dmg != 14 {
		t.Errorf("multi hit: expected turns=3, dmg=14, got turns=%d, dmg=%d", turns, dmg)
	}
	// Enemy HP 12 -> 1 turn, 0 damage
	turns, dmg = EstimateAutoFight(wMulti, 12, 7)
	if turns != 1 || dmg != 0 {
		t.Errorf("multi hit 1-turn: expected turns=1, dmg=0, got turns=%d, dmg=%d", turns, dmg)
	}

	// 3. Multi-hit+: int(10*0.75)=7 per hit, 14 per turn. Enemy HP 25 -> 2 turns, 7 damage
	wMultiPlus := model.Weapon{BaseAtk: 10, Slots: []model.OrbType{model.OrbMultiHitPlus}}
	turns, dmg = EstimateAutoFight(wMultiPlus, 25, 7)
	if turns != 2 || dmg != 7 {
		t.Errorf("multi hit+: expected turns=2, dmg=7, got turns=%d, dmg=%d", turns, dmg)
	}

	// 4. Weak weapon below 1 per hit (minimum 1): BaseAtk=0, rate=0.65 -> dmgPerHit=1, 2 hits -> 2 dmg/turn
	wWeak := model.Weapon{BaseAtk: 0, Slots: []model.OrbType{model.OrbMultiHit}}
	turns, dmg = EstimateAutoFight(wWeak, 5, 3)
	if turns != 3 || dmg != 6 {
		t.Errorf("weak weapon: expected turns=3, dmg=6, got turns=%d, dmg=%d", turns, dmg)
	}

	// 5. One-shot enemy (1 turn, 0 damage)
	turns, dmg = EstimateAutoFight(wStandard, 10, 50)
	if turns != 1 || dmg != 0 {
		t.Errorf("1-shot: expected turns=1, dmg=0, got turns=%d, dmg=%d", turns, dmg)
	}

	// 6. Critical, poison and vamp orbs do not change the estimate
	wAllOtherOrbs := model.Weapon{
		BaseAtk: 10,
		Slots: []model.OrbType{
			model.OrbCriticalPlus,
			model.OrbPoisonPlus,
			model.OrbVampPlus,
		},
	}
	turnsOther, dmgOther := EstimateAutoFight(wAllOtherOrbs, 25, 7)
	if turnsOther != 3 || dmgOther != 14 {
		t.Errorf("other orbs: expected turns=3, dmg=14, got turns=%d, dmg=%d", turnsOther, dmgOther)
	}
}

func TestDangerHPLine(t *testing.T) {
	tests := []struct {
		maxHP    int
		expected int
	}{
		{100, 30},
		{50, 15},
		{33, 9},
		{1000, 300},
	}
	for _, tc := range tests {
		got := DangerHPLine(tc.maxHP)
		if got != tc.expected {
			t.Errorf("DangerHPLine(%d) = %d; want %d", tc.maxHP, got, tc.expected)
		}
	}
}

func TestEstimateAutoFight_ConservativeProperty(t *testing.T) {
	// The estimate must be pessimistic (turns >= actual turns, damage >= actual damage)
	testCases := []struct {
		weapon   model.Weapon
		enemyHP  int
		enemyAtk int
	}{
		{
			weapon: model.Weapon{
				BaseAtk: 20,
				Slots:   []model.OrbType{model.OrbCriticalPlus, model.OrbMultiHitPlus, model.OrbVampPlus},
			},
			enemyHP:  100,
			enemyAtk: 15,
		},
		{
			weapon: model.Weapon{
				BaseAtk: 15,
				Slots:   []model.OrbType{model.OrbPoisonPlus, model.OrbCritical},
			},
			enemyHP:  80,
			enemyAtk: 10,
		},
		{
			weapon: model.Weapon{
				BaseAtk: 30,
				Slots:   []model.OrbType{model.OrbMultiHit, model.OrbPoison},
			},
			enemyHP:  120,
			enemyAtk: 20,
		},
	}

	msgs := i18n.GetMessages(model.LanguageJA)
	for idx, tc := range testCases {
		estTurns, estDmg := EstimateAutoFight(tc.weapon, tc.enemyHP, tc.enemyAtk)

		// Run the random simulation several times
		for seed := 0; seed < 10; seed++ {
			enemy := &model.Enemy{HP: tc.enemyHP, MaxHP: tc.enemyHP, Atk: tc.enemyAtk}
			res := RunAutoCombat(tc.weapon, 1000, 1000, enemy, msgs, nil)

			if !res.EnemyDefeated {
				t.Fatalf("case %d seed %d: expected enemy defeat", idx, seed)
			}
			if res.Turns > estTurns {
				t.Errorf("case %d seed %d: actual turns %d > estimated turns %d", idx, seed, res.Turns, estTurns)
			}
			if res.TotalDmgTaken > estDmg {
				t.Errorf("case %d seed %d: actual dmg %d > estimated dmg %d", idx, seed, res.TotalDmgTaken, estDmg)
			}
		}
	}
}
