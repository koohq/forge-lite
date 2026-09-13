package engine

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"

	"forge-lite/internal/model"
)

var DropOrbs = []model.OrbType{
	model.OrbMultiHit,
	model.OrbCritical,
	model.OrbVamp,
	model.OrbPoison,
}

var PlusOrbs = []model.OrbType{
	model.OrbMultiHitPlus,
	model.OrbCriticalPlus,
	model.OrbVampPlus,
	model.OrbPoisonPlus,
}

func defaultRandInt(n int) int {
	if n <= 0 {
		return 0
	}
	return rand.IntN(n)
}

func GetRandomOrb(randInt func(n int) int) model.OrbType {
	if randInt == nil {
		randInt = defaultRandInt
	}
	return DropOrbs[randInt(len(DropOrbs))]
}

func GetRandomPlusOrb(randInt func(n int) int) model.OrbType {
	if randInt == nil {
		randInt = defaultRandInt
	}
	return PlusOrbs[randInt(len(PlusOrbs))]
}

// CreateStarterDungeon generates the 5-floor starter dungeon definition.
func CreateStarterDungeon(theme model.Theme, randInt func(n int) int) model.DungeonDef {
	if randInt == nil {
		randInt = defaultRandInt
	}

	bases := theme.EnemyBases
	getBase := func(idx int, fallback string) string {
		if idx < len(bases) && bases[idx] != "" {
			return bases[idx]
		}
		if len(bases) > 0 && bases[0] != "" {
			return bases[0]
		}
		return fallback
	}

	b1 := getBase(0, "Slime")
	b2 := getBase(1, "Kobold")
	b3 := getBase(2, "Orc")
	b4 := getBase(3, "Golem")
	b5 := getBase(4, "Dragon")

	enemies := []model.EnemyTemplate{
		{
			Name: b1,
			HP:   15,
			Atk:  3,
			GetDrops: func() model.RunInventory {
				return model.RunInventory{Scrolls: 1, Orbs: []model.OrbType{}}
			},
		},
		{
			Name: b2,
			HP:   28,
			Atk:  6,
			GetDrops: func() model.RunInventory {
				return model.RunInventory{Scrolls: 1, Orbs: []model.OrbType{}}
			},
		},
		{
			Name: b3,
			HP:   45,
			Atk:  10,
			GetDrops: func() model.RunInventory {
				return model.RunInventory{Scrolls: 0, Orbs: []model.OrbType{GetRandomOrb(randInt)}}
			},
		},
		{
			Name: b4,
			HP:   75,
			Atk:  14,
			GetDrops: func() model.RunInventory {
				return model.RunInventory{Scrolls: 2, Orbs: []model.OrbType{}}
			},
		},
		{
			Name: fmt.Sprintf("%s (BOSS)", b5),
			HP:   130,
			Atk:  20,
			GetDrops: func() model.RunInventory {
				return model.RunInventory{Scrolls: 3, Orbs: []model.OrbType{GetRandomOrb(randInt)}}
			},
		},
	}

	return model.DungeonDef{
		Name:   theme.Dungeons.Starter,
		Floors: 5,
		GetEnemy: func(floor int) (model.EnemyTemplate, error) {
			if floor < 1 || floor > len(enemies) {
				return model.EnemyTemplate{}, fmt.Errorf("invalid floor: %d", floor)
			}
			return enemies[floor-1], nil
		},
	}
}

// CreateDeepDungeon generates the 10-floor deep dungeon definition.
func CreateDeepDungeon(theme model.Theme, recommended string, randInt func(n int) int) model.DungeonDef {
	if randInt == nil {
		randInt = defaultRandInt
	}

	prefixes := theme.EnemyPrefixes
	bases := theme.EnemyBases

	getPrefix := func(idx int) string {
		if len(prefixes) == 0 {
			return ""
		}
		return prefixes[idx%len(prefixes)]
	}
	getBase := func(idx int) string {
		if len(bases) == 0 {
			return "Enemy"
		}
		return bases[idx%len(bases)]
	}

	return model.DungeonDef{
		Name:        theme.Dungeons.Deep,
		Floors:      10,
		Recommended: recommended,
		GetEnemy: func(floor int) (model.EnemyTemplate, error) {
			getDeepFloorScrolls := func() int {
				return randInt(4) + 2
			}

			if floor >= 1 && floor <= 3 {
				c1 := model.EnemyTemplate{
					Name: strings.TrimSpace(getPrefix(0) + " " + getBase(2)),
					HP:   80,
					Atk:  12,
					GetDrops: func() model.RunInventory {
						return model.RunInventory{Scrolls: getDeepFloorScrolls(), Orbs: []model.OrbType{}}
					},
				}
				c2 := model.EnemyTemplate{
					Name: strings.TrimSpace(getPrefix(1) + " " + getBase(4)),
					HP:   120,
					Atk:  16,
					GetDrops: func() model.RunInventory {
						return model.RunInventory{Scrolls: getDeepFloorScrolls(), Orbs: []model.OrbType{}}
					},
				}
				if randInt(2) == 0 {
					return c1, nil
				}
				return c2, nil
			}

			if floor >= 4 && floor <= 6 {
				c1 := model.EnemyTemplate{
					Name: strings.TrimSpace(getPrefix(2) + " " + getBase(3)),
					HP:   180,
					Atk:  20,
					GetDrops: func() model.RunInventory {
						return model.RunInventory{Scrolls: getDeepFloorScrolls(), Orbs: []model.OrbType{}}
					},
				}
				c2 := model.EnemyTemplate{
					Name: strings.TrimSpace(getPrefix(3) + " " + getBase(7)),
					HP:   220,
					Atk:  24,
					GetDrops: func() model.RunInventory {
						return model.RunInventory{Scrolls: getDeepFloorScrolls(), Orbs: []model.OrbType{}}
					},
				}
				if randInt(2) == 0 {
					return c1, nil
				}
				return c2, nil
			}

			if floor >= 7 && floor <= 9 {
				c1 := model.EnemyTemplate{
					Name: strings.TrimSpace(getPrefix(4) + " " + getBase(5)),
					HP:   280,
					Atk:  28,
					GetDrops: func() model.RunInventory {
						return model.RunInventory{Scrolls: getDeepFloorScrolls(), Orbs: []model.OrbType{}}
					},
				}
				c2 := model.EnemyTemplate{
					Name: strings.TrimSpace(getPrefix(5) + " " + getBase(6)),
					HP:   350,
					Atk:  32,
					GetDrops: func() model.RunInventory {
						return model.RunInventory{Scrolls: getDeepFloorScrolls(), Orbs: []model.OrbType{}}
					},
				}
				if randInt(2) == 0 {
					return c1, nil
				}
				return c2, nil
			}

			if floor == 10 {
				return model.EnemyTemplate{
					Name: strings.TrimSpace(fmt.Sprintf("%s %s (BOSS)", getPrefix(6), getBase(7))),
					HP:   500,
					Atk:  36,
					GetDrops: func() model.RunInventory {
						return model.RunInventory{
							Scrolls: 10,
							Orbs:    []model.OrbType{GetRandomPlusOrb(randInt)},
						}
					},
				}, nil
			}

			return model.EnemyTemplate{}, fmt.Errorf("invalid floor: %d", floor)
		},
	}
}

// GenerateEndlessEnemy generates an enemy for the specified floor of Infinite Abyss.
func GenerateEndlessEnemy(floor int, theme model.Theme, randInt func(n int) int) model.EnemyTemplate {
	if randInt == nil {
		randInt = defaultRandInt
	}

	prefixes := theme.EnemyPrefixes
	enemyNames := theme.EnemyBases

	prefix := "Dark"
	if len(prefixes) > 0 {
		prefix = prefixes[randInt(len(prefixes))]
	}

	baseName := "Beast"
	if len(enemyNames) > 0 {
		baseName = enemyNames[randInt(len(enemyNames))]
	}

	isBoss := (floor%5 == 0)
	bossTag := "(Mini-Boss)"
	separator := " "
	if theme.Language == model.LanguageJA {
		bossTag = "(中ボス)"
		separator = ""
	}

	var name string
	if isBoss {
		name = fmt.Sprintf("%s%s%s %s", prefix, separator, baseName, bossTag)
	} else {
		name = fmt.Sprintf("%s%s%s", prefix, separator, baseName)
	}

	hp := int(60 + math.Pow(float64(floor), 1.4)*12)
	atk := int(8 + float64(floor)*2.5)

	scrolls := 1 + floor/2
	if scrolls > 10 {
		scrolls = 10
	}

	var orbs []model.OrbType
	if isBoss {
		orbs = []model.OrbType{GetRandomPlusOrb(randInt)}
	}

	return model.EnemyTemplate{
		Name: name,
		HP:   hp,
		Atk:  atk,
		GetDrops: func() model.RunInventory {
			return model.RunInventory{
				Scrolls: scrolls,
				Orbs:    orbs,
			}
		},
	}
}
