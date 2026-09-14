package engine

import (
	"math/rand/v2"

	"github.com/koohq/forge-lite/internal/i18n"
	"github.com/koohq/forge-lite/internal/model"
)

// CalcWeaponAtk calculates weapon attack power: baseAtk + (plus * 2).
func CalcWeaponAtk(weapon model.Weapon) int {
	return weapon.BaseAtk + weapon.Plus*2
}

// OrbSynergies caches orb flags on a weapon.
type OrbSynergies struct {
	HasMulti      bool
	HasMultiPlus  bool
	HasCrit       bool
	HasCritPlus   bool
	HasVamp       bool
	HasVampPlus   bool
	HasPoison     bool
	HasPoisonPlus bool
}

// GetOrbSynergies extracts equipped orb effects.
func GetOrbSynergies(slots []model.OrbType) OrbSynergies {
	var syn OrbSynergies
	for _, s := range slots {
		switch s {
		case model.OrbMultiHitPlus:
			syn.HasMulti = true
			syn.HasMultiPlus = true
		case model.OrbMultiHit:
			syn.HasMulti = true
		case model.OrbCriticalPlus:
			syn.HasCrit = true
			syn.HasCritPlus = true
		case model.OrbCritical:
			syn.HasCrit = true
		case model.OrbVampPlus:
			syn.HasVamp = true
			syn.HasVampPlus = true
		case model.OrbVamp:
			syn.HasVamp = true
		case model.OrbPoisonPlus:
			syn.HasPoison = true
			syn.HasPoisonPlus = true
		case model.OrbPoison:
			syn.HasPoison = true
		}
	}
	return syn
}

// HitResult holds outcome of a single attack strike.
type HitResult struct {
	HitIndex    int
	IsCrit      bool
	Damage      int
	Heal        int
	PoisonAdded int
}

// ResolvePlayerAttack performs player attacks for one turn.
func ResolvePlayerAttack(
	weapon model.Weapon,
	playerHP, maxHP int,
	enemy *model.Enemy,
	randFloat func() float64,
) ([]HitResult, int) {
	if randFloat == nil {
		randFloat = rand.Float64
	}

	syn := GetOrbSynergies(weapon.Slots)
	baseDmg := CalcWeaponAtk(weapon)

	hits := 1
	if syn.HasMulti {
		hits = 2
	}

	rate := 1.0
	if syn.HasMultiPlus {
		rate = 0.75
	} else if syn.HasMulti {
		rate = 0.65
	}

	critChance := 0.0
	if syn.HasCritPlus {
		critChance = 0.40
	} else if syn.HasCrit {
		critChance = 0.25
	}

	vampRate := 0.0
	if syn.HasVampPlus {
		vampRate = 0.25
	} else if syn.HasVamp {
		vampRate = 0.15
	}

	poisonAdd := 0
	if syn.HasPoisonPlus {
		poisonAdd = 2
	} else if syn.HasPoison {
		poisonAdd = 1
	}

	results := make([]HitResult, 0, hits)
	currHP := playerHP

	for i := 1; i <= hits; i++ {
		if enemy.HP <= 0 {
			break
		}

		dmg := int(float64(baseDmg) * rate)
		if dmg < 1 {
			dmg = 1
		}

		isCrit := false
		if critChance > 0 && randFloat() < critChance {
			isCrit = true
			dmg *= 2
		}

		enemy.HP -= dmg

		heal := 0
		if vampRate > 0 {
			h := int(float64(dmg) * vampRate)
			if h < 1 {
				h = 1
			}
			actual := maxHP - currHP
			if h < actual {
				actual = h
			}
			if actual > 0 {
				currHP += actual
				heal = actual
			}
		}

		addedPoison := 0
		if poisonAdd > 0 {
			enemy.Poison += poisonAdd
			addedPoison = poisonAdd
		}

		results = append(results, HitResult{
			HitIndex:    i,
			IsCrit:      isCrit,
			Damage:      dmg,
			Heal:        heal,
			PoisonAdded: addedPoison,
		})
	}

	return results, currHP
}

// AutoCombatResult summarizes the automated combat execution.
type AutoCombatResult struct {
	EnemyDefeated     bool
	PlayerFainted     bool
	DangerZoneStopped bool
	Turns             int
	TotalDmgDealt     int
	TotalDmgTaken     int
	TotalHealed       int
	StopReason        string
	FinalPlayerHP     int
	FinalEnemyHP      int
}

// RunAutoCombat simulates combat until enemy defeat, player fainting, or danger zone HP <= 30%.
func RunAutoCombat(
	weapon model.Weapon,
	playerHP, maxHP int,
	enemy *model.Enemy,
	msgs i18n.Messages,
	randFloat func() float64,
) AutoCombatResult {
	if randFloat == nil {
		randFloat = rand.Float64
	}

	dangerHP := int(float64(maxHP) * 0.3)
	currPlayerHP := playerHP

	res := AutoCombatResult{}

	for currPlayerHP > 0 && enemy.HP > 0 {
		res.Turns++

		hits, nextPlayerHP := ResolvePlayerAttack(weapon, currPlayerHP, maxHP, enemy, randFloat)
		for _, h := range hits {
			res.TotalDmgDealt += h.Damage
			res.TotalHealed += h.Heal
		}
		currPlayerHP = nextPlayerHP

		if enemy.HP <= 0 {
			res.EnemyDefeated = true
			res.StopReason = msgs.BattleAutoReasonEnemyDefeated
			break
		}

		// Poison tick at turn end
		if enemy.Poison > 0 {
			poisonDmg := enemy.Poison * 3
			enemy.HP -= poisonDmg
			res.TotalDmgDealt += poisonDmg
			if enemy.HP <= 0 {
				res.EnemyDefeated = true
				res.StopReason = msgs.BattleAutoReasonPoisonDefeated
				break
			}
		}

		// Enemy counterattack
		currPlayerHP -= enemy.Atk
		res.TotalDmgTaken += enemy.Atk

		if currPlayerHP <= 0 {
			res.PlayerFainted = true
			res.StopReason = msgs.BattleAutoReasonFainted
			break
		}

		// Danger threshold check
		if currPlayerHP <= dangerHP {
			res.DangerZoneStopped = true
			res.StopReason = msgs.BattleAutoReasonDangerHp(currPlayerHP, maxHP)
			break
		}
	}

	res.FinalPlayerHP = currPlayerHP
	res.FinalEnemyHP = enemy.HP
	return res
}
