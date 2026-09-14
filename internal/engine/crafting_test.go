package engine

import (
	"testing"

	"github.com/koohq/forge-lite/internal/model"
)

func TestEnhanceWeapon(t *testing.T) {
	w := model.Weapon{Name: "Sword", BaseAtk: 10, Plus: 0}
	scrolls := 10

	count, err := EnhanceWeapon(&w, &scrolls, 3)
	if err != nil || count != 3 {
		t.Fatalf("enhance failed: %v", err)
	}
	if w.Plus != 3 || scrolls != 7 {
		t.Errorf("expected plus 3, scrolls 7; got plus %d, scrolls %d", w.Plus, scrolls)
	}

	_, err = EnhanceWeapon(&w, &scrolls, 10)
	if err != ErrInvalidCount {
		t.Errorf("expected ErrInvalidCount, got %v", err)
	}
}

func TestEnhanceHP(t *testing.T) {
	p := model.Player{MaxHP: 50, HP: 50}
	scrolls := 5

	used, gain, remainder, err := EnhanceHP(&p, &scrolls, 5)
	if err != nil {
		t.Fatalf("enhance HP failed: %v", err)
	}
	if used != 4 || gain != 20 || !remainder {
		t.Errorf("expected used 4, gain 20, remainder true; got used=%d, gain=%d, rem=%v", used, gain, remainder)
	}
	if p.MaxHP != 70 || scrolls != 1 {
		t.Errorf("expected maxHP 70, scrolls 1; got maxHP=%d, scrolls=%d", p.MaxHP, scrolls)
	}
}

func TestAttachOrb(t *testing.T) {
	w := model.Weapon{Slots: []model.OrbType{model.OrbMultiHit, model.OrbCritical}}
	stock := []model.OrbType{model.OrbVamp, model.OrbPoison}

	// 3rd slot append
	attached, replaced, err := AttachOrb(&w, &stock, 0, -1)
	if err != nil {
		t.Fatalf("attach failed: %v", err)
	}
	if attached != model.OrbVamp || replaced != "" {
		t.Errorf("expected attached Vamp, replaced empty; got %s, %s", attached, replaced)
	}
	if len(w.Slots) != 3 || len(stock) != 1 {
		t.Errorf("slots: %d, stock: %d", len(w.Slots), len(stock))
	}

	// Overwrite slot 1
	attached, replaced, err = AttachOrb(&w, &stock, 0, 1)
	if err != nil {
		t.Fatalf("overwrite failed: %v", err)
	}
	if attached != model.OrbPoison || replaced != model.OrbCritical {
		t.Errorf("expected attached Poison, replaced Critical; got %s, %s", attached, replaced)
	}
	if w.Slots[1] != model.OrbPoison || len(stock) != 0 {
		t.Errorf("unexpected slot content or stock size")
	}
}

func TestDismantleAllOrbs(t *testing.T) {
	stock := []model.OrbType{model.OrbMultiHit, model.OrbCritical, model.OrbVamp}
	scrolls := 0

	used, gained, rem := DismantleAllOrbs(&stock, &scrolls)
	if used != 2 || gained != 1 || rem == nil || *rem != model.OrbVamp {
		t.Errorf("unexpected dismantle all result: used=%d, gained=%d, rem=%v", used, gained, rem)
	}
	if scrolls != 1 || len(stock) != 1 {
		t.Errorf("scrolls=%d, stock len=%d", scrolls, len(stock))
	}
}

func TestSynthesizeOrb(t *testing.T) {
	stock := []model.OrbType{
		model.OrbMultiHit,
		model.OrbPoison,
		model.OrbMultiHit,
	}

	recipes := GetAvailableSynthesisRecipes(stock)
	if len(recipes) != 1 || recipes[0].From != model.OrbMultiHit || recipes[0].To != model.OrbMultiHitPlus {
		t.Fatalf("unexpected recipes: %+v", recipes)
	}

	err := SynthesizeOrb(&stock, recipes[0].From, recipes[0].To)
	if err != nil {
		t.Fatalf("synthesize failed: %v", err)
	}
	if len(stock) != 2 {
		t.Fatalf("expected 2 remaining orbs, got %d", len(stock))
	}
	if stock[0] != model.OrbPoison || stock[1] != model.OrbMultiHitPlus {
		t.Errorf("unexpected stock contents: %+v", stock)
	}
}
