package engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"forge-lite/internal/i18n"
	"forge-lite/internal/model"
	"forge-lite/internal/save"
	"forge-lite/internal/theme"
)

type Game struct {
	reader         *bufio.Reader
	saveFilePath   string
	language       model.Language
	theme          model.Theme
	customThemeDef *model.ThemeDefinition
	player         model.Player
	stockScrolls   int
	stockOrbs      []model.OrbType
	deepestFloor   int
	isRunning      bool
}

func NewGame(r io.Reader) *Game {
	if r == nil {
		r = os.Stdin
	}

	lang := i18n.DetectDefaultLanguage()
	customDef, _ := theme.LoadCustomThemeDefIfExists("")

	var th model.Theme
	if customDef != nil {
		th = theme.ResolveTheme(*customDef, lang)
	} else {
		th = theme.GetPresetTheme(theme.DefaultThemeID, lang)
	}

	return &Game{
		reader:         bufio.NewReader(r),
		language:       lang,
		theme:          th,
		customThemeDef: customDef,
		player: model.Player{
			MaxHP: 50,
			HP:    50,
			Weapon: model.Weapon{
				Name:    th.GetTargetRankName(0),
				BaseAtk: 10,
				Plus:    0,
				Slots:   []model.OrbType{},
			},
		},
		stockScrolls: 0,
		stockOrbs:    []model.OrbType{},
		deepestFloor: 0,
		isRunning:    true,
	}
}

func (g *Game) msg() i18n.Messages {
	return i18n.GetMessages(g.language)
}

func (g *Game) readLine(prompt string) string {
	if prompt != "" {
		fmt.Print(prompt)
	}
	line, err := g.reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return ""
	}
	return strings.TrimSpace(line)
}

func (g *Game) getOrbDesc(orb model.OrbType) string {
	desc, ok := g.theme.Orbs[orb]
	if !ok {
		return string(orb)
	}
	return fmt.Sprintf("%s (%s)", desc.Name, desc.Desc)
}

func (g *Game) getWeaponAtk() int {
	return CalcWeaponAtk(g.player.Weapon)
}

func (g *Game) saveGame() {
	deepest := g.deepestFloor
	maxHp := g.player.MaxHP
	lang := g.language
	themeID := g.theme.ID

	data := &model.SaveData{
		Weapon:       g.player.Weapon,
		StockScrolls: g.stockScrolls,
		StockOrbs:    g.stockOrbs,
		MaxHP:        &maxHp,
		DeepestFloor: &deepest,
		Language:     &lang,
		ThemeID:      &themeID,
	}

	err := save.SaveSaveData(g.saveFilePath, data)
	if err != nil {
		fmt.Printf("%s %v\n", g.msg().SaveFailed, err)
	} else {
		fmt.Println(g.msg().SavedSuccess)
	}
}

func (g *Game) loadGame() bool {
	data, err := save.LoadSaveData(g.saveFilePath)
	if err != nil {
		return false
	}

	if data.Language != nil && (*data.Language == model.LanguageEN || *data.Language == model.LanguageJA) {
		g.language = *data.Language
	}

	if g.customThemeDef != nil {
		g.theme = theme.ResolveTheme(*g.customThemeDef, g.language)
	} else if data.ThemeID != nil {
		g.theme = theme.GetPresetTheme(*data.ThemeID, g.language)
	} else {
		g.theme = theme.GetPresetTheme(theme.DefaultThemeID, g.language)
	}

	slots := make([]model.OrbType, len(data.Weapon.Slots))
	copy(slots, data.Weapon.Slots)

	g.player.Weapon = model.Weapon{
		Name:    g.theme.GetTargetRankName(data.Weapon.Plus),
		BaseAtk: data.Weapon.BaseAtk,
		Plus:    data.Weapon.Plus,
		Slots:   slots,
	}
	g.stockScrolls = data.StockScrolls
	g.stockOrbs = make([]model.OrbType, len(data.StockOrbs))
	copy(g.stockOrbs, data.StockOrbs)

	if data.MaxHP != nil && *data.MaxHP > 0 {
		g.player.MaxHP = *data.MaxHP
	} else {
		g.player.MaxHP = 50
	}
	g.player.HP = g.player.MaxHP

	if data.DeepestFloor != nil && *data.DeepestFloor >= 0 {
		g.deepestFloor = *data.DeepestFloor
	} else {
		g.deepestFloor = 0
	}

	return true
}

func (g *Game) Start() {
	fmt.Println(g.msg().Title)
	fmt.Printf("\n%s\n", g.msg().StartContinue)
	fmt.Println(g.msg().StartNew)

	for {
		choice := g.readLine(g.msg().StartPrompt)
		if choice == "1" {
			if g.loadGame() {
				fmt.Println(g.msg().SaveLoaded)
				fmt.Println(g.msg().SaveStats(g.player.MaxHP, g.deepestFloor))
			} else {
				fmt.Println(g.msg().SaveNotFoundOrCorrupted)
			}
			break
		}
		if choice == "2" {
			fmt.Println(g.msg().NewGameStarted)
			break
		}
		fmt.Println(g.msg().InvalidChoice1or2)
	}

	for g.isRunning {
		g.hubPhase()
	}
}

func (g *Game) hubPhase() {
	g.player.HP = g.player.MaxHP
	fmt.Println("\n----------------------------------------------")
	fmt.Println(g.msg().HubHeader(g.theme.HubTitle))
	fmt.Println(g.msg().HubStats(g.theme.HPLabel, g.player.MaxHP, g.theme.Dungeons.Abyss, g.deepestFloor))
	fmt.Println(g.msg().HubEquip(g.theme.GetTargetPrefix(), g.player.Weapon.Name, g.theme.PlusPrefix, g.player.Weapon.Plus, g.theme.StatLabel, g.getWeaponAtk()))

	slotsStr := g.msg().HubEmptySlot
	if len(g.player.Weapon.Slots) > 0 {
		var b strings.Builder
		for i, s := range g.player.Weapon.Slots {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(fmt.Sprintf("[%s]", s))
		}
		slotsStr = b.String()
	}
	fmt.Println(g.msg().HubSlots(g.theme.InstallVerb, len(g.player.Weapon.Slots), slotsStr))

	orbsStr := g.msg().HubNone
	if len(g.stockOrbs) > 0 {
		var b strings.Builder
		for i, o := range g.stockOrbs {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(fmt.Sprintf("[%s]", o))
		}
		orbsStr = b.String()
	}
	fmt.Println(g.msg().HubStorage(g.theme.StorageLabel, g.theme.ResourceName, g.stockScrolls, g.theme.OrbLabel, orbsStr))
	fmt.Println("----------------------------------------------")

	fmt.Println(g.msg().HubMenu1Dungeon)
	fmt.Println(g.msg().HubMenu2Enhance(g.theme.EnhanceVerb, g.theme.StatLabel))
	fmt.Println(g.msg().HubMenu3AttachOrb(g.theme.InstallVerb))
	fmt.Println(g.msg().HubMenu4Disassemble(g.theme.OrbLabel))
	fmt.Println(g.msg().HubMenu5Synthesize(g.theme.OrbLabel))
	fmt.Println(g.msg().HubMenu6EnhanceHp(g.theme.EnhanceHPVerb, g.theme.HPLabel))
	fmt.Println(g.msg().HubMenu7Settings)
	fmt.Println(g.msg().HubMenu0Exit)

	choice := g.readLine(g.msg().ChooseAction)
	switch choice {
	case "1":
		g.chooseDungeonPhase()
	case "2":
		g.upgradeWeaponPhase()
	case "3":
		g.attachOrbPhase()
	case "4":
		g.disassembleOrbPhase()
	case "5":
		g.synthesizeOrbPhase()
	case "6":
		g.upgradeHpPhase()
	case "7":
		g.settingsPhase()
	case "0":
		fmt.Println(g.msg().Farewell)
		g.isRunning = false
	default:
		fmt.Println(g.msg().InvalidChoice)
	}
}

func (g *Game) settingsPhase() {
	for {
		fmt.Println(g.msg().SettingsTitle)
		fmt.Println(g.msg().SettingsCurrent(g.language, g.theme.Name))
		fmt.Println(g.msg().SettingsMenu1Lang)
		fmt.Println(g.msg().SettingsMenu2Theme)
		fmt.Println(g.msg().SettingsMenu3ExportTemplate)
		fmt.Println(g.msg().SettingsMenu0Back)

		choice := g.readLine(g.msg().ChooseAction)
		switch choice {
		case "0":
			return
		case "1":
			nextLang := model.LanguageEN
			if g.language == model.LanguageEN {
				nextLang = model.LanguageJA
			}
			g.language = nextLang
			if g.customThemeDef != nil && g.theme.ID == g.customThemeDef.ID {
				g.theme = theme.ResolveTheme(*g.customThemeDef, g.language)
			} else {
				g.theme = theme.GetPresetTheme(g.theme.ID, g.language)
			}
			g.player.Weapon.Name = g.theme.GetTargetRankName(g.player.Weapon.Plus)
			g.saveGame()
			fmt.Println(g.msg().SettingsLangChanged(g.language))
		case "2":
			g.switchThemeSubmenu()
		case "3":
			err := theme.ExportCustomThemeTemplate("")
			if err == nil {
				fmt.Println(g.msg().SettingsTemplateExported(theme.CustomThemeExamplePath))
			} else {
				fmt.Printf("%s %v\n", g.msg().SettingsTemplateExportFailed, err)
			}
		default:
			fmt.Println(g.msg().InvalidChoice)
		}
	}
}

func (g *Game) switchThemeSubmenu() {
	fmt.Println(g.msg().SettingsThemeSelectTitle)

	var themeOptions []model.Theme
	if g.customThemeDef != nil {
		themeOptions = append(themeOptions, theme.ResolveTheme(*g.customThemeDef, g.language))
	} else {
		checked, _ := theme.LoadCustomThemeDefIfExists("")
		if checked != nil {
			g.customThemeDef = checked
			themeOptions = append(themeOptions, theme.ResolveTheme(*checked, g.language))
		}
	}

	themeOptions = append(themeOptions,
		theme.GetPresetTheme("classic_fantasy", g.language),
		theme.GetPresetTheme("cyberpunk", g.language),
		theme.GetPresetTheme("partner_sync", g.language),
	)

	for i, t := range themeOptions {
		fmt.Printf("%d: %s (%s) [%s]\n", i+1, t.Name, t.Language, t.ID)
	}
	fmt.Println("0: Cancel")

	choice := g.readLine(g.msg().ChooseAction)
	idx, err := strconv.Atoi(choice)
	if err == nil && idx >= 1 && idx <= len(themeOptions) {
		selected := themeOptions[idx-1]
		g.theme = selected
		g.player.Weapon.Name = selected.GetTargetRankName(g.player.Weapon.Plus)
		g.saveGame()
		fmt.Println(g.msg().SettingsThemeChanged(selected.Name))
	}
}

func (g *Game) chooseDungeonPhase() {
	starter := CreateStarterDungeon(g.theme, nil)
	deep := CreateDeepDungeon(g.theme, g.msg().RecommendedDeep, nil)

	fmt.Println(g.msg().DungeonSelectTitle)
	fmt.Println(g.msg().DungeonOptionStarter(starter.Name, starter.Floors))
	fmt.Println(g.msg().DungeonOptionDeep(deep.Name, deep.Floors, deep.Recommended))
	fmt.Println(g.msg().DungeonOptionAbyss(g.theme.Dungeons.Abyss))
	fmt.Println(g.msg().DungeonCancel)

	choice := g.readLine(g.msg().DungeonPrompt)
	switch choice {
	case "1":
		g.dungeonPhase(starter)
	case "2":
		g.dungeonPhase(deep)
	case "3":
		g.endlessDungeonPhase()
	case "0":
		fmt.Println(g.msg().DungeonCanceled)
	default:
		fmt.Println(g.msg().InvalidChoice)
	}
}

func (g *Game) upgradeWeaponPhase() {
	if g.stockScrolls <= 0 {
		fmt.Println(g.msg().EnhanceNoResource(g.theme.ResourceName))
		return
	}

	fmt.Println(g.msg().EnhanceTitle(g.theme.GetTargetPrefix()))
	fmt.Println(g.msg().EnhanceCurrentRes(g.theme.ResourceName, g.stockScrolls))
	fmt.Println(g.msg().EnhanceOpt1)
	fmt.Println(g.msg().EnhanceOpt2(g.theme.ResourceName))
	fmt.Println(g.msg().EnhanceOpt3(g.theme.ResourceName, g.stockScrolls))
	fmt.Println(g.msg().EnhanceOpt0)

	choice := g.readLine(g.msg().ChooseAction)
	switch choice {
	case "1":
		prevRank := g.theme.GetTargetRankName(g.player.Weapon.Plus)
		EnhanceWeapon(&g.player.Weapon, &g.stockScrolls, 1)
		newRank := g.theme.GetTargetRankName(g.player.Weapon.Plus)
		g.player.Weapon.Name = newRank
		if prevRank != newRank {
			fmt.Println(g.msg().RankUpgraded(newRank))
		}
		fmt.Println(g.msg().EnhanceSuccessSingle(g.theme.EnhanceVerb, g.theme.GetTargetPrefix(), g.player.Weapon.Name, g.player.Weapon.Plus, g.getWeaponAtk()))
		g.saveGame()
	case "2":
		inputStr := g.readLine(g.msg().EnhancePromptCount(g.theme.ResourceName, g.stockScrolls))
		count, err := strconv.Atoi(inputStr)
		if err != nil || count < 1 || count > g.stockScrolls {
			fmt.Println(g.msg().EnhanceInvalidCount)
			return
		}
		prevRank := g.theme.GetTargetRankName(g.player.Weapon.Plus)
		EnhanceWeapon(&g.player.Weapon, &g.stockScrolls, count)
		newRank := g.theme.GetTargetRankName(g.player.Weapon.Plus)
		g.player.Weapon.Name = newRank
		if prevRank != newRank {
			fmt.Println(g.msg().RankUpgraded(newRank))
		}
		fmt.Println(g.msg().EnhanceSuccessMultiple(count, g.theme.ResourceName, g.theme.EnhanceVerb, g.player.Weapon.Name, g.player.Weapon.Plus, g.getWeaponAtk()))
		g.saveGame()
	case "3":
		count := g.stockScrolls
		prevRank := g.theme.GetTargetRankName(g.player.Weapon.Plus)
		EnhanceWeapon(&g.player.Weapon, &g.stockScrolls, count)
		newRank := g.theme.GetTargetRankName(g.player.Weapon.Plus)
		g.player.Weapon.Name = newRank
		if prevRank != newRank {
			fmt.Println(g.msg().RankUpgraded(newRank))
		}
		fmt.Println(g.msg().EnhanceSuccessAll(count, g.theme.ResourceName, g.theme.EnhanceVerb, g.player.Weapon.Name, g.player.Weapon.Plus, g.getWeaponAtk()))
		g.saveGame()
	case "0":
		fmt.Println(g.msg().EnhanceCanceled)
	default:
		fmt.Println(g.msg().InvalidChoice)
	}
}

func (g *Game) upgradeHpPhase() {
	if g.stockScrolls < 2 {
		fmt.Println(g.msg().HPUpgradeNotEnough(g.theme.ResourceName, g.stockScrolls))
		return
	}

	maxCount := g.stockScrolls / 2
	fmt.Println(g.msg().HPUpgradeTitle)
	fmt.Println(g.msg().HPUpgradeCurrentHp(g.player.MaxHP))
	fmt.Println(g.msg().HPUpgradeCurrentRes(g.theme.ResourceName, g.stockScrolls))
	fmt.Println(g.msg().HPUpgradeOpt1(g.theme.ResourceName))
	fmt.Println(g.msg().HPUpgradeOpt2(g.theme.ResourceName, maxCount*2))
	fmt.Println(g.msg().HPUpgradeOpt3(g.theme.ResourceName, maxCount*2, maxCount*10))
	fmt.Println(g.msg().HPUpgradeOpt0)

	choice := g.readLine(g.msg().ChooseAction)
	switch choice {
	case "1":
		prev := g.player.MaxHP
		EnhanceHP(&g.player, &g.stockScrolls, 2)
		fmt.Println(g.msg().HPUpgradeSuccessSingle(prev, g.player.MaxHP))
		g.saveGame()
	case "2":
		inputStr := g.readLine(g.msg().HPUpgradePromptCount(g.theme.ResourceName, maxCount*2))
		val, err := strconv.Atoi(inputStr)
		if err != nil || val < 2 || val > g.stockScrolls {
			fmt.Println(g.msg().HPUpgradeInvalidCount)
			return
		}
		prev := g.player.MaxHP
		used, _, rem, _ := EnhanceHP(&g.player, &g.stockScrolls, val)
		fmt.Println(g.msg().HPUpgradeSuccessMultiple(used, g.theme.ResourceName, prev, g.player.MaxHP))
		if rem {
			fmt.Println(g.msg().HPUpgradeOddRemainder)
		}
		g.saveGame()
	case "3":
		prev := g.player.MaxHP
		used, _, _, _ := EnhanceHP(&g.player, &g.stockScrolls, maxCount*2)
		fmt.Println(g.msg().HPUpgradeSuccessAll(used, g.theme.ResourceName, prev, g.player.MaxHP))
		g.saveGame()
	case "0":
		fmt.Println(g.msg().HPUpgradeCanceled)
	default:
		fmt.Println(g.msg().InvalidChoice)
	}
}

func (g *Game) attachOrbPhase() {
	if len(g.stockOrbs) == 0 {
		fmt.Println(g.msg().OrbAttachNoOrbs)
		return
	}

	fmt.Println(g.msg().OrbAttachTitle)
	for i, orb := range g.stockOrbs {
		fmt.Printf("%d: [%s] %s\n", i+1, orb, g.getOrbDesc(orb))
	}
	fmt.Println("0: Cancel")

	selectStr := g.readLine(g.msg().OrbAttachPrompt)
	idx, err := strconv.Atoi(selectStr)
	if err != nil || idx < 1 || idx > len(g.stockOrbs) {
		return
	}
	orbIdx := idx - 1

	if len(g.player.Weapon.Slots) < 3 {
		attached, _, _ := AttachOrb(&g.player.Weapon, &g.stockOrbs, orbIdx, -1)
		fmt.Println(g.msg().OrbAttachSuccess(string(attached)))
		g.saveGame()
	} else {
		fmt.Println(g.msg().OrbAttachFullTitle)
		for sIdx, s := range g.player.Weapon.Slots {
			fmt.Printf("%d: [%s]\n", sIdx+1, s)
		}
		fmt.Println("0: Cancel")

		rStr := g.readLine(g.msg().OrbAttachReplacePrompt)
		rIdx, rErr := strconv.Atoi(rStr)
		if rErr != nil || rIdx < 1 || rIdx > 3 {
			return
		}
		attached, removed, _ := AttachOrb(&g.player.Weapon, &g.stockOrbs, orbIdx, rIdx-1)
		fmt.Println(g.msg().OrbAttachReplaced(string(removed), string(attached)))
		g.saveGame()
	}
}

func (g *Game) disassembleOrbPhase() {
	if len(g.stockOrbs) < 2 {
		fmt.Println(g.msg().OrbDisassembleNotEnough)
		return
	}

	fmt.Println(g.msg().OrbDisassembleTitle(g.theme.ResourceName))
	fmt.Println(g.msg().OrbDisassembleCount(len(g.stockOrbs)))
	fmt.Println(g.msg().OrbDisassembleOpt1)
	fmt.Println(g.msg().OrbDisassembleOpt2(g.theme.ResourceName, len(g.stockOrbs)/2))
	fmt.Println(g.msg().OrbDisassembleOpt0)

	choice := g.readLine(g.msg().ChooseAction)
	if choice == "0" {
		fmt.Println(g.msg().OrbDisassembleCanceled)
		return
	}
	if choice == "2" {
		used, scrolls, rem := DismantleAllOrbs(&g.stockOrbs, &g.stockScrolls)
		fmt.Println(g.msg().OrbDisassembleAllSuccess(used, g.theme.ResourceName, scrolls, g.stockScrolls))
		if rem != nil {
			fmt.Println(g.msg().OrbDisassembleRemainder(string(*rem)))
		}
		g.saveGame()
		return
	}
	if choice != "1" {
		fmt.Println(g.msg().InvalidChoice)
		return
	}

	fmt.Println(g.msg().OrbDisassembleIndivTitle)
	for i, orb := range g.stockOrbs {
		fmt.Printf("%d: [%s] %s\n", i+1, orb, g.getOrbDesc(orb))
	}
	fmt.Println("0: Cancel")

	s1 := g.readLine(g.msg().OrbDisassemblePick1)
	idx1, err1 := strconv.Atoi(s1)
	if err1 != nil || idx1 < 1 || idx1 > len(g.stockOrbs) {
		return
	}
	orb1 := g.stockOrbs[idx1-1]
	fmt.Println(g.msg().OrbDisassembleSelected1(string(orb1)))

	s2 := g.readLine(g.msg().OrbDisassemblePick2)
	idx2, err2 := strconv.Atoi(s2)
	if err2 != nil || idx2 < 1 || idx2 > len(g.stockOrbs) {
		return
	}
	if idx1 == idx2 {
		fmt.Println(g.msg().OrbDisassembleSameError)
		return
	}

	o1, o2, _ := DismantleOrbPair(&g.stockOrbs, &g.stockScrolls, idx1-1, idx2-1)
	fmt.Println(g.msg().OrbDisassembleSuccess(string(o1), string(o2), g.theme.ResourceName, g.stockScrolls))
	g.saveGame()
}

func (g *Game) synthesizeOrbPhase() {
	recipes := GetAvailableSynthesisRecipes(g.stockOrbs)
	if len(recipes) == 0 {
		fmt.Println(g.msg().OrbSynthNoRecipes)
		return
	}

	fmt.Println(g.msg().OrbSynthTitle)
	for i, rec := range recipes {
		fmt.Println(g.msg().OrbSynthRecipeItem(i+1, string(rec.From), rec.Count, string(rec.To), g.getOrbDesc(rec.To)))
	}
	fmt.Println("0: Cancel")

	choice := g.readLine(g.msg().OrbSynthPrompt)
	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(recipes) {
		return
	}

	target := recipes[idx-1]
	_ = SynthesizeOrb(&g.stockOrbs, target.From, target.To)
	fmt.Println(g.msg().OrbSynthSuccess(string(target.From), string(target.To)))
	g.saveGame()
}

func (g *Game) dungeonPhase(dungeon model.DungeonDef) {
	fmt.Println(g.msg().DungeonEnter(dungeon.Name))
	runInv := model.RunInventory{Scrolls: 0, Orbs: []model.OrbType{}}

	for floor := 1; floor <= dungeon.Floors; floor++ {
		template, err := dungeon.GetEnemy(floor)
		if err != nil {
			break
		}

		enemy := model.Enemy{
			Name:   template.Name,
			HP:     template.HP,
			MaxHP:  template.HP,
			Atk:    template.Atk,
			Poison: 0,
		}

		fmt.Println(g.msg().DungeonEncounter(floor, dungeon.Floors, enemy.Name))
		survived := g.battle(&enemy)
		if !survived {
			fmt.Println(g.msg().DungeonDefeatPrompt)
			fmt.Println(g.msg().DungeonItemsLost(g.theme.ResourceName, runInv.Scrolls, len(runInv.Orbs)))
			fmt.Println(g.msg().DungeonCarriedBack)
			g.saveGame()
			return
		}

		drops := template.GetDrops()
		runInv.Scrolls += drops.Scrolls
		runInv.Orbs = append(runInv.Orbs, drops.Orbs...)

		fmt.Println(g.msg().DungeonVictory(enemy.Name))
		if drops.Scrolls > 0 {
			fmt.Println(g.msg().DungeonDropScroll(g.theme.ResourceName, drops.Scrolls))
		}
		for _, o := range drops.Orbs {
			name := string(o)
			if d, ok := g.theme.Orbs[o]; ok {
				name = d.Name
			}
			fmt.Println(g.msg().DungeonDropOrb(string(o), name))
		}

		if floor == dungeon.Floors {
			fmt.Println(g.msg().DungeonComplete(dungeon.Name))
			g.stockScrolls += runInv.Scrolls
			g.stockOrbs = append(g.stockOrbs, runInv.Orbs...)
			g.saveGame()
			return
		}

		fmt.Println(g.msg().DungeonCurrentStatus(g.player.HP, g.player.MaxHP, g.theme.ResourceName, runInv.Scrolls, len(runInv.Orbs)))
		fmt.Println(g.msg().DungeonNextFloor)
		fmt.Println(g.msg().DungeonRetreat)

		nextAct := g.readLine(g.msg().DungeonPromptAction)
		if nextAct == "2" {
			fmt.Printf("\n>> %s\n", g.theme.RetreatMessage)
			g.stockScrolls += runInv.Scrolls
			g.stockOrbs = append(g.stockOrbs, runInv.Orbs...)
			fmt.Println(g.msg().DungeonRetreatSuccess(g.theme.ResourceName, runInv.Scrolls, len(runInv.Orbs)))
			g.saveGame()
			return
		}
	}
}

func (g *Game) endlessDungeonPhase() {
	fmt.Println(g.msg().EndlessEnter(g.theme.Dungeons.Abyss))
	fmt.Println(g.msg().EndlessIntro)

	startFloor := 1
	if g.deepestFloor >= 11 {
		maxStart := ((g.deepestFloor-1)/10)*10 + 1
		var availableFloors []int
		for f := 1; f <= maxStart; f += 10 {
			availableFloors = append(availableFloors, f)
		}

		fmt.Println(g.msg().EndlessStartFloorTitle)
		fmt.Println(g.msg().EndlessHighestRecord(g.deepestFloor))
		for i, f := range availableFloors {
			fmt.Printf("%d: %s\n", i+1, g.msg().EndlessStartOption(f))
		}
		fmt.Println(g.msg().EndlessCancelOption)

		for {
			choice := g.readLine(g.msg().StartPrompt)
			if choice == "0" {
				fmt.Println(g.msg().DungeonCanceled)
				return
			}
			idx, err := strconv.Atoi(choice)
			if err == nil && idx >= 1 && idx <= len(availableFloors) {
				startFloor = availableFloors[idx-1]
				break
			}
			fmt.Println(g.msg().InvalidChoice)
		}
	}

	fmt.Println(g.msg().EndlessStartFrom(startFloor))
	runInv := model.RunInventory{Scrolls: 0, Orbs: []model.OrbType{}}
	floor := startFloor

	for {
		if floor > g.deepestFloor {
			g.deepestFloor = floor
			fmt.Println(g.msg().EndlessRecordUpdated(floor))
		}

		template := GenerateEndlessEnemy(floor, g.theme, nil)
		enemy := model.Enemy{
			Name:   template.Name,
			HP:     template.HP,
			MaxHP:  template.HP,
			Atk:    template.Atk,
			Poison: 0,
		}

		fmt.Println(g.msg().EndlessFloorEncounter(floor, enemy.Name))
		survived := g.battle(&enemy)
		if !survived {
			fmt.Println(g.msg().DungeonDefeatPrompt)
			fmt.Println(g.msg().DungeonItemsLost(g.theme.ResourceName, runInv.Scrolls, len(runInv.Orbs)))
			fmt.Println(g.msg().DungeonCarriedBack)
			g.saveGame()
			return
		}

		drops := template.GetDrops()
		runInv.Scrolls += drops.Scrolls
		runInv.Orbs = append(runInv.Orbs, drops.Orbs...)

		fmt.Println(g.msg().DungeonVictory(enemy.Name))
		if drops.Scrolls > 0 {
			fmt.Println(g.msg().DungeonDropScroll(g.theme.ResourceName, drops.Scrolls))
		}
		for _, o := range drops.Orbs {
			name := string(o)
			if d, ok := g.theme.Orbs[o]; ok {
				name = d.Name
			}
			fmt.Println(g.msg().DungeonDropOrb(string(o), name))
		}

		fmt.Println(g.msg().DungeonCurrentStatus(g.player.HP, g.player.MaxHP, g.theme.ResourceName, runInv.Scrolls, len(runInv.Orbs)))
		fmt.Println(g.msg().DungeonNextFloor)
		fmt.Println(g.msg().DungeonRetreat)

		nextAct := g.readLine(g.msg().DungeonPromptAction)
		if nextAct == "2" {
			fmt.Printf("\n>> %s\n", g.theme.RetreatMessage)
			g.stockScrolls += runInv.Scrolls
			g.stockOrbs = append(g.stockOrbs, runInv.Orbs...)
			fmt.Println(g.msg().DungeonRetreatSuccess(g.theme.ResourceName, runInv.Scrolls, len(runInv.Orbs)))
			g.saveGame()
			return
		}

		floor++
	}
}

func (g *Game) battle(enemy *model.Enemy) bool {
	for g.player.HP > 0 && enemy.HP > 0 {
		fmt.Println(g.msg().BattleVs(g.player.HP, g.player.MaxHP, enemy.Name, enemy.HP, enemy.MaxHP, enemy.Poison))
		fmt.Println(g.msg().BattleCmdAttack)
		fmt.Println(g.msg().BattleCmdRetreat)
		fmt.Println(g.msg().BattleCmdAuto)

		act := g.readLine(g.msg().BattleCmdPrompt)
		switch act {
		case "2":
			fmt.Println(g.msg().BattleRetreatCombat)
			return false
		case "3":
			fmt.Println(g.msg().BattleAutoStart)
			res := RunAutoCombat(g.player.Weapon, g.player.HP, g.player.MaxHP, enemy, g.msg(), nil)
			g.player.HP = res.FinalPlayerHP

			fmt.Println(g.msg().BattleAutoSummaryTitle)
			fmt.Println(g.msg().BattleAutoReason(res.StopReason))
			fmt.Println(g.msg().BattleAutoTurns(res.Turns))
			fmt.Println(g.msg().BattleAutoDamageDealt(res.TotalDmgDealt))
			fmt.Println(g.msg().BattleAutoDamageTaken(res.TotalDmgTaken))
			if res.TotalHealed > 0 {
				fmt.Println(g.msg().BattleAutoHealed(res.TotalHealed))
			}
			fmt.Println(g.msg().BattleAutoRemainingHp(g.player.HP, g.player.MaxHP, enemy.Name, enemy.HP, enemy.MaxHP))
			fmt.Println("==============================================")

			if enemy.HP <= 0 {
				return true
			}
			if g.player.HP <= 0 {
				return false
			}
			continue
		case "1":
			syn := GetOrbSynergies(g.player.Weapon.Slots)
			hits, nextHP := ResolvePlayerAttack(g.player.Weapon, g.player.HP, g.player.MaxHP, enemy, nil)
			g.player.HP = nextHP

			for _, h := range hits {
				hitIndexStr := ""
				if syn.HasMulti {
					hitIndexStr = g.msg().BattleHitNumber(h.HitIndex)
				}
				fmt.Println(g.msg().BattlePlayerAttack(hitIndexStr, h.IsCrit, enemy.Name, h.Damage))
				if h.Heal > 0 {
					fmt.Println(g.msg().BattleVampHeal(h.Heal, g.player.HP))
				}
				if h.PoisonAdded > 0 {
					fmt.Println(g.msg().BattlePoisonInflict(enemy.Name, h.PoisonAdded, enemy.Poison))
				}
			}

			if enemy.HP <= 0 {
				return true
			}

			// Poison tick
			if enemy.Poison > 0 {
				poisonDmg := enemy.Poison * 3
				enemy.HP -= poisonDmg
				fmt.Println(g.msg().BattlePoisonTick(enemy.Name, poisonDmg))
				if enemy.HP <= 0 {
					return true
				}
			}

			// Enemy counterattack
			fmt.Println(g.msg().BattleEnemyCounter(enemy.Name, enemy.Atk))
			g.player.HP -= enemy.Atk
			if g.player.HP <= 0 {
				return false
			}
		default:
			fmt.Println(g.msg().BattleInvalidCmd)
		}
	}

	return g.player.HP > 0
}
