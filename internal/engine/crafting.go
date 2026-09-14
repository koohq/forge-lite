package engine

import (
	"errors"

	"github.com/koohq/forge-lite/internal/model"
)

var (
	ErrNotEnoughResources = errors.New("not enough resources")
	ErrInvalidCount       = errors.New("invalid count specified")
	ErrSlotOutOfRange     = errors.New("slot index out of range")
	ErrOrbOutOfRange      = errors.New("orb index out of range")
	ErrSameOrbSelected    = errors.New("cannot select the same orb twice")
	ErrNoSynthesisRecipe  = errors.New("no synthesis recipe available")
)

var NormalOrbs = []model.OrbType{
	model.OrbMultiHit,
	model.OrbCritical,
	model.OrbVamp,
	model.OrbPoison,
}

var SynthesisRecipes = map[model.OrbType]model.OrbType{
	model.OrbMultiHit: model.OrbMultiHitPlus,
	model.OrbCritical: model.OrbCriticalPlus,
	model.OrbVamp:     model.OrbVampPlus,
	model.OrbPoison:   model.OrbPoisonPlus,
}

// EnhanceWeapon upgrades the weapon plus value.
func EnhanceWeapon(weapon *model.Weapon, stockScrolls *int, count int) (int, error) {
	if count < 1 || count > *stockScrolls {
		return 0, ErrInvalidCount
	}
	*stockScrolls -= count
	weapon.Plus += count
	return count, nil
}

// EnhanceHP upgrades player max HP. 2 scrolls = +10 HP.
func EnhanceHP(player *model.Player, stockScrolls *int, requestedScrolls int) (usedScrolls, hpGain int, hasRemainder bool, err error) {
	if requestedScrolls < 2 || requestedScrolls > *stockScrolls {
		return 0, 0, false, ErrInvalidCount
	}
	times := requestedScrolls / 2
	usedScrolls = times * 2
	hpGain = times * 10
	hasRemainder = (requestedScrolls % 2 != 0)

	*stockScrolls -= usedScrolls
	player.MaxHP += hpGain
	player.HP = player.MaxHP
	return usedScrolls, hpGain, hasRemainder, nil
}

// AttachOrb equips an orb to weapon. If slots < 3, appends. If replaceIndex >= 0 and slots == 3, replaces and returns replaced orb.
func AttachOrb(weapon *model.Weapon, stockOrbs *[]model.OrbType, orbIndex int, replaceIndex int) (attached model.OrbType, replaced model.OrbType, err error) {
	if orbIndex < 0 || orbIndex >= len(*stockOrbs) {
		return "", "", ErrOrbOutOfRange
	}
	targetOrb := (*stockOrbs)[orbIndex]

	if len(weapon.Slots) < 3 {
		weapon.Slots = append(weapon.Slots, targetOrb)
		*stockOrbs = append((*stockOrbs)[:orbIndex], (*stockOrbs)[orbIndex+1:]...)
		return targetOrb, "", nil
	}

	if replaceIndex < 0 || replaceIndex >= 3 {
		return "", "", ErrSlotOutOfRange
	}

	old := weapon.Slots[replaceIndex]
	weapon.Slots[replaceIndex] = targetOrb
	*stockOrbs = append((*stockOrbs)[:orbIndex], (*stockOrbs)[orbIndex+1:]...)
	return targetOrb, old, nil
}

// DismantleOrbPair removes two orbs from stock and awards 1 scroll.
func DismantleOrbPair(stockOrbs *[]model.OrbType, stockScrolls *int, idx1, idx2 int) (orb1, orb2 model.OrbType, err error) {
	if idx1 == idx2 {
		return "", "", ErrSameOrbSelected
	}
	if idx1 < 0 || idx1 >= len(*stockOrbs) || idx2 < 0 || idx2 >= len(*stockOrbs) {
		return "", "", ErrOrbOutOfRange
	}

	first := idx1
	second := idx2
	if first < second {
		first, second = second, first
	}

	orb1 = (*stockOrbs)[first]
	orb2 = (*stockOrbs)[second]

	*stockOrbs = append((*stockOrbs)[:first], (*stockOrbs)[first+1:]...)
	*stockOrbs = append((*stockOrbs)[:second], (*stockOrbs)[second+1:]...)

	*stockScrolls++
	return orb1, orb2, nil
}

// DismantleAllOrbs disassembles all unassigned orbs into scrolls (2 orbs -> 1 scroll).
func DismantleAllOrbs(stockOrbs *[]model.OrbType, stockScrolls *int) (usedCount, gainedScrolls int, remainingOrb *model.OrbType) {
	total := len(*stockOrbs)
	gainedScrolls = total / 2
	usedCount = gainedScrolls * 2

	if total%2 == 1 {
		last := (*stockOrbs)[total-1]
		remainingOrb = &last
		*stockOrbs = []model.OrbType{last}
	} else {
		*stockOrbs = []model.OrbType{}
	}

	*stockScrolls += gainedScrolls
	return usedCount, gainedScrolls, remainingOrb
}

// SynthesisOption represents an available synthesis recipe based on stock.
type SynthesisOption struct {
	From  model.OrbType
	To    model.OrbType
	Count int
}

// GetAvailableSynthesisRecipes finds all synthesizable normal orbs with count >= 2.
func GetAvailableSynthesisRecipes(stockOrbs []model.OrbType) []SynthesisOption {
	counts := make(map[model.OrbType]int)
	for _, o := range stockOrbs {
		counts[o]++
	}

	var options []SynthesisOption
	for _, normal := range NormalOrbs {
		c := counts[normal]
		if c >= 2 {
			options = append(options, SynthesisOption{
				From:  normal,
				To:    SynthesisRecipes[normal],
				Count: c,
			})
		}
	}
	return options
}

// SynthesizeOrb consumes 2 of targetFrom and appends targetTo to stock.
func SynthesizeOrb(stockOrbs *[]model.OrbType, targetFrom, targetTo model.OrbType) error {
	removed := 0
	list := *stockOrbs
	for i := len(list) - 1; i >= 0 && removed < 2; i-- {
		if list[i] == targetFrom {
			list = append(list[:i], list[i+1:]...)
			removed++
		}
	}
	if removed < 2 {
		return ErrNotEnoughResources
	}

	list = append(list, targetTo)
	*stockOrbs = list
	return nil
}
