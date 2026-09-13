package engine

import (
	"testing"

	"forge-lite/internal/i18n"
	"forge-lite/internal/model"
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
