package i18n

import (
	"fmt"
	"os"
	"strings"

	"github.com/koohq/forge-lite/internal/model"
)

// Messages defines all UI text and string templates.
type Messages struct {
	Title                   string
	StartPrompt             string
	StartContinue           string
	StartNew                string
	InvalidChoice1or2       string
	SaveLoaded              string
	SaveStats               func(hp, deepest int) string
	Untested                string
	SaveNotFoundOrCorrupted string
	NewGameStarted          string
	SavedSuccess            string
	SaveFailed              string
	LoadFailed              string
	Farewell                string

	// Hub
	HubHeader           func(title string) string
	HubStats            func(hpLabel string, hp int, abyssName string, deepest int) string
	HubEquip            func(label, name, plusPrefix string, plus int, statLabel string, atk int) string
	HubSlots            func(verb string, count int, slots string) string
	HubEmptySlot        string
	HubStorage          func(storage, resource string, scrolls int, orb, orbs string) string
	HubNone             string
	HubMenu1Dungeon     string
	HubMenu2Enhance     func(verb, statLabel string) string
	HubMenu3AttachOrb   func(verb string) string
	HubMenu4Disassemble func(orb string) string
	HubMenu5Synthesize  func(orb string) string
	HubMenu6EnhanceHp   func(verb, hpLabel string) string
	HubMenu7Settings    string
	HubMenu0Exit        string
	ChooseAction        string
	InvalidChoice       string

	// Dungeon Select
	DungeonSelectTitle   string
	DungeonOptionStarter func(name string, floors int) string
	DungeonOptionDeep    func(name string, floors int, rec string) string
	DungeonOptionAbyss   func(name string) string
	DungeonCancel        string
	DungeonPrompt        string
	DungeonCanceled      string
	RecommendedDeep      string

	// Weapon / Target Upgrade
	EnhanceTitle           func(label string) string
	EnhanceCurrentRes      func(res string, count int) string
	EnhanceOpt1            string
	EnhanceOpt2            func(res string) string
	EnhanceOpt3            func(res string, count int) string
	EnhanceOpt0            string
	EnhanceNoResource      func(res string) string
	EnhancePromptCount     func(res string, max int) string
	EnhanceInvalidCount    string
	EnhanceSuccessSingle   func(verb, label, name string, plus, atk int) string
	EnhanceSuccessMultiple func(count int, res, verb, name string, plus, atk int) string
	EnhanceSuccessAll      func(count int, res, verb, name string, plus, atk int) string
	EnhanceCanceled        string
	RankUpgraded           func(name string) string

	// HP Upgrade
	HPUpgradeTitle           string
	HPUpgradeCurrentHp       func(hp int) string
	HPUpgradeCurrentRes      func(res string, count int) string
	HPUpgradeNotEnough       func(res string, count int) string
	HPUpgradeOpt1            func(res string) string
	HPUpgradeOpt2            func(res string, max int) string
	HPUpgradeOpt3            func(res string, used, gain int) string
	HPUpgradeOpt0            string
	HPUpgradePromptCount     func(res string, max int) string
	HPUpgradeInvalidCount    string
	HPUpgradeSuccessSingle   func(prev, next int) string
	HPUpgradeSuccessMultiple func(used int, res string, prev, next int) string
	HPUpgradeOddRemainder    string
	HPUpgradeSuccessAll      func(used int, res string, prev, next int) string
	HPUpgradeCanceled        string

	// Orb Attach
	OrbAttachNoOrbs        string
	OrbAttachTitle         string
	OrbAttachPrompt        string
	OrbAttachSuccess       func(orb string) string
	OrbAttachFullTitle     string
	OrbAttachReplacePrompt string
	OrbAttachReplaced      func(removed, added string) string

	// Orb Disassemble
	OrbDisassembleNotEnough  string
	OrbDisassembleTitle      func(res string) string
	OrbDisassembleCount      func(count int) string
	OrbDisassembleOpt1       string
	OrbDisassembleOpt2       func(res string, scrolls int) string
	OrbDisassembleOpt0       string
	OrbDisassembleAllSuccess func(usedCount int, res string, scrolls, total int) string
	OrbDisassembleRemainder  func(orb string) string
	OrbDisassembleIndivTitle string
	OrbDisassemblePick1      string
	OrbDisassembleSelected1  func(orb string) string
	OrbDisassemblePick2      string
	OrbDisassembleSameError  string
	OrbDisassembleSuccess    func(orb1, orb2, res string, total int) string
	OrbDisassembleCanceled   string

	// Orb Synthesize
	OrbSynthNoRecipes  string
	OrbSynthTitle      string
	OrbSynthRecipeItem func(idx int, from string, count int, to, toDesc string) string
	OrbSynthPrompt     string
	OrbSynthSuccess    func(from, to string) string

	// Dungeon Exploration
	DungeonEnter          func(name string) string
	DungeonEncounter      func(floor, total int, name string) string
	DungeonDefeatPrompt   string
	DungeonItemsLost      func(res string, scrolls, orbs int) string
	DungeonCarriedBack    string
	DungeonVictory        func(enemy string) string
	DungeonDropScroll     func(res string, count int) string
	DungeonDropOrb        func(orb, name string) string
	DungeonComplete       func(name string) string
	DungeonCurrentStatus  func(hp, maxHp int, res string, scrolls, orbs int) string
	DungeonNextFloor      string
	DungeonRetreat        string
	DungeonPromptAction   string
	DungeonRetreatSuccess func(res string, scrolls, orbs int) string

	// Endless
	EndlessEnter           func(name string) string
	EndlessIntro           string
	EndlessStartFloorTitle string
	EndlessHighestRecord   func(floor int) string
	EndlessStartOption     func(floor int) string
	EndlessCancelOption    string
	EndlessStartFrom       func(floor int) string
	EndlessFloorEncounter  func(floor int, name string) string
	EndlessRecordUpdated   func(floor int) string

	// Battle
	BattleVs                      func(php, pmax int, ename string, ehp, emax, poison int) string
	BattleCmdAttack               string
	BattleCmdRetreat              string
	BattleCmdAuto                 string
	BattleCmdPrompt               string
	BattleRetreatCombat           string
	BattleAutoStart               string
	BattleAutoSummaryTitle        string
	BattleAutoReason              func(reason string) string
	BattleAutoTurns               func(turns int) string
	BattleAutoDamageDealt         func(dmg int) string
	BattleAutoDamageTaken         func(dmg int) string
	BattleAutoHealed              func(heal int) string
	BattleAutoRemainingHp         func(php, pmax int, ename string, ehp, emax int) string
	BattleAutoReasonEnemyDefeated string
	BattleAutoReasonPoisonDefeated string
	BattleAutoReasonFainted       string
	BattleAutoReasonDangerHp      func(hp, maxHp int) string
	BattleInvalidCmd              string
	BattlePlayerAttack            func(hitIndex string, isCrit bool, ename string, dmg int) string
	BattleCritLabel               string
	BattleHitNumber               func(hit int) string
	BattleVampHeal                func(heal, hp int) string
	BattlePoisonInflict           func(ename string, added, total int) string
	BattlePoisonTick              func(ename string, dmg int) string
	BattleEnemyCounter            func(ename string, dmg int) string

	// Settings
	SettingsTitle                string
	SettingsCurrent              func(lang, theme string) string
	SettingsMenu1Lang            string
	SettingsMenu2Theme           string
	SettingsMenu3ExportTemplate  string
	SettingsMenu0Back            string
	SettingsLangChanged          func(lang string) string
	SettingsThemeSelectTitle     string
	SettingsThemeChanged         func(theme string) string
	SettingsTemplateExported     func(path string) string
	SettingsTemplateExportFailed string
}

var (
	MessagesJA = Messages{
		Title: "==============================================\n   Minimal Rogue-lite Prototype (CUI Ver)   \n==============================================",
		StartPrompt: "\n選択してください: ",
		StartContinue: "1: つづきから (save.json を読み込んで開始)",
		StartNew: "2: はじめから (初期状態で開始)",
		InvalidChoice1or2: ">> 無効な選択です。1 または 2 を入力してください。",
		SaveLoaded: ">> セーブデータを読み込みました！",
		SaveStats: func(hp, deepest int) string {
			rec := "未挑戦"
			if deepest > 0 {
				rec = fmt.Sprintf("B%dF", deepest)
			}
			return fmt.Sprintf(">> 最大体力: HP %d | 最高到達階層: %s", hp, rec)
		},
		Untested: "未挑戦",
		SaveNotFoundOrCorrupted: ">> save.json が見つからないか破損しています。新規データで開始します。",
		NewGameStarted: ">> はじめからゲームを開始します。",
		SavedSuccess: ">> セーブデータを保存しました。(save.json)",
		SaveFailed: ">> セーブデータの保存に失敗しました:",
		LoadFailed: ">> セーブデータの読み込みに失敗しました:",
		Farewell: "お疲れ様でした。",

		// Hub
		HubHeader: func(title string) string { return title },
		HubStats: func(hpLabel string, hp int, abyssName string, deepest int) string {
			rec := "未挑戦"
			if deepest > 0 {
				rec = fmt.Sprintf("B%dF", deepest)
			}
			return fmt.Sprintf("%s: HP %d | %s 最高到達: %s", hpLabel, hp, abyssName, rec)
		},
		HubEquip: func(label, name, plusPrefix string, plus int, statLabel string, atk int) string {
			return fmt.Sprintf("%s: %s%s%d (%s: %d)", label, name, plusPrefix, plus, statLabel, atk)
		},
		HubSlots: func(verb string, count int, slots string) string {
			return fmt.Sprintf("%sスロット [%d/3]: %s", verb, count, slots)
		},
		HubEmptySlot: "(空き)",
		HubStorage: func(storage, resource string, scrolls int, orb, orbs string) string {
			return fmt.Sprintf("%s: %s x%d | 未装着%s: %s", storage, resource, scrolls, orb, orbs)
		},
		HubNone: "(なし)",
		HubMenu1Dungeon: "1: ダンジョンへ出撃 (探索開始)",
		HubMenu2Enhance: func(verb, statLabel string) string { return fmt.Sprintf("2: %s (%sの強化)", verb, statLabel) },
		HubMenu3AttachOrb: func(verb string) string {
			return fmt.Sprintf("3: %s (パッシブ効果の装着)", verb)
		},
		HubMenu4Disassemble: func(orb string) string {
			return fmt.Sprintf("4: %sを分解 (素材への還元)", orb)
		},
		HubMenu5Synthesize: func(orb string) string {
			return fmt.Sprintf("5: %sを合成 (上位性能への強化)", orb)
		},
		HubMenu6EnhanceHp: func(verb, hpLabel string) string {
			return fmt.Sprintf("6: %s (%s +10)", verb, hpLabel)
		},
		HubMenu7Settings: "7: 設定 / Settings",
		HubMenu0Exit: "0: ゲーム終了",
		ChooseAction: "\n行動を選択してください: ",
		InvalidChoice: "無効な選択です。",

		// Dungeon Select
		DungeonSelectTitle: "\n--- 出撃ダンジョン選択 ---",
		DungeonOptionStarter: func(name string, floors int) string {
			return fmt.Sprintf("1: %s (全%d階 / 初級)", name, floors)
		},
		DungeonOptionDeep: func(name string, floors int, rec string) string {
			return fmt.Sprintf("2: %s (全%d階 / %s)", name, floors, rec)
		},
		DungeonOptionAbyss: func(name string) string {
			return fmt.Sprintf("3: %s (エンドレス / 階層無制限)", name)
		},
		DungeonCancel: "0: キャンセル (拠点に戻る)",
		DungeonPrompt: "\nダンジョンを選択してください: ",
		DungeonCanceled: ">> 出撃を取りやめました。",
		RecommendedDeep: "推奨+25以上の上級ダンジョン",

		// Weapon / Target Upgrade
		EnhanceTitle: func(label string) string {
			return fmt.Sprintf("\n--- %sの強化 ---", label)
		},
		EnhanceCurrentRes: func(res string, count int) string {
			return fmt.Sprintf("現在の所持%s: %d個", res, count)
		},
		EnhanceOpt1: "1: 1回鍛える (+1)",
		EnhanceOpt2: func(res string) string {
			return fmt.Sprintf("2: 指定した回数分鍛える（消費する%sを直接入力）", res)
		},
		EnhanceOpt3: func(res string, count int) string {
			return fmt.Sprintf("3: 所持している%sですべて鍛える (+%d)", res, count)
		},
		EnhanceOpt0: "0: キャンセル",
		EnhanceNoResource: func(res string) string {
			return fmt.Sprintf(">> %sがありません！", res)
		},
		EnhancePromptCount: func(res string, max int) string {
			return fmt.Sprintf("消費する%sの数を入力してください (1〜%d): ", res, max)
		},
		EnhanceInvalidCount: ">> 無効な数値です。1以上の所持数以内の数値を入力してください。",
		EnhanceSuccessSingle: func(verb, label, name string, plus, atk int) string {
			return fmt.Sprintf(">> 【%s】%sを行いました！ %s+%d (攻撃力: %d)", label, verb, name, plus, atk)
		},
		EnhanceSuccessMultiple: func(count int, res, verb, name string, plus, atk int) string {
			return fmt.Sprintf(">> %sを %d個 消費して%sを行いました！ %s+%d (攻撃力: %d)", res, count, verb, name, plus, atk)
		},
		EnhanceSuccessAll: func(count int, res, verb, name string, plus, atk int) string {
			return fmt.Sprintf(">> %sをすべて(%d個)消費して一括で%sを行いました！ %s+%d (攻撃力: %d)", res, count, verb, name, plus, atk)
		},
		EnhanceCanceled: ">> 強化をキャンセルしました。",
		RankUpgraded: func(name string) string {
			return fmt.Sprintf(">> ランクが上がった！ 【%s】に進化した！", name)
		},

		// HP Upgrade
		HPUpgradeTitle: "\n--- 体力の強化 ---",
		HPUpgradeCurrentHp: func(hp int) string {
			return fmt.Sprintf("現在の最大HP: %d", hp)
		},
		HPUpgradeCurrentRes: func(res string, count int) string {
			return fmt.Sprintf("所持%s: %d個 (2個消費で最大HP+10)", res, count)
		},
		HPUpgradeNotEnough: func(res string, count int) string {
			return fmt.Sprintf(">> %sが足りません！ (必要: 2個 / 所持: %d個)", res, count)
		},
		HPUpgradeOpt1: func(res string) string {
			return fmt.Sprintf("1: 1回強化する (%s2個消費 -> 最大HP+10)", res)
		},
		HPUpgradeOpt2: func(res string, max int) string {
			return fmt.Sprintf("2: 指定した回数分強化する（消費する%sを直接入力、最大: %d個）", res, max)
		},
		HPUpgradeOpt3: func(res string, used, gain int) string {
			return fmt.Sprintf("3: 所持している%sですべて強化する (%d個消費 -> 最大HP+%d)", res, used, gain)
		},
		HPUpgradeOpt0: "0: キャンセル",
		HPUpgradePromptCount: func(res string, max int) string {
			return fmt.Sprintf("消費する%sの数を入力してください (2〜%d): ", res, max)
		},
		HPUpgradeInvalidCount: ">> 無効な数値です。2個以上、所持数以内の数値を入力してください。",
		HPUpgradeSuccessSingle: func(prev, next int) string {
			return fmt.Sprintf(">> 体力を強化しました！ 最大HP: %d -> %d", prev, next)
		},
		HPUpgradeSuccessMultiple: func(used int, res string, prev, next int) string {
			return fmt.Sprintf(">> %sを %d個 消費して強化しました！ 最大HP: %d -> %d", res, used, prev, next)
		},
		HPUpgradeOddRemainder: ">> ※ 2個単位での消費のため、端数の1個は温存されました。",
		HPUpgradeSuccessAll: func(used int, res string, prev, next int) string {
			return fmt.Sprintf(">> %sをすべて(%d個)消費して一括で強化しました！ 最大HP: %d -> %d", res, used, prev, next)
		},
		HPUpgradeCanceled: ">> 強化をキャンセルしました。",

		// Orb Attach
		OrbAttachNoOrbs: ">> 装着できるオーブが倉庫にありません！",
		OrbAttachTitle: "\n--- 倉庫のオーブ一覧 ---",
		OrbAttachPrompt: "装着するオーブの番号を選択: ",
		OrbAttachSuccess: func(orb string) string {
			return fmt.Sprintf(">> スロットに [%s] を装着しました！", orb)
		},
		OrbAttachFullTitle: "\nスロットが満杯です。上書きする枠を選んでください:",
		OrbAttachReplacePrompt: "番号を選択: ",
		OrbAttachReplaced: func(removed, added string) string {
			return fmt.Sprintf(">> [%s] を破棄し、[%s] を装着しました！", removed, added)
		},

		// Orb Disassemble
		OrbDisassembleNotEnough: ">> 分解には倉庫にオーブが2個以上必要です！",
		OrbDisassembleTitle: func(res string) string {
			return fmt.Sprintf("\n--- オーブの分解 (任意のオーブ2個 -> %s1個) ---", res)
		},
		OrbDisassembleCount: func(count int) string {
			return fmt.Sprintf("倉庫の未装着オーブ: %d個", count)
		},
		OrbDisassembleOpt1: "1: オーブを個別に選んで分解",
		OrbDisassembleOpt2: func(res string, scrolls int) string {
			return fmt.Sprintf("2: 倉庫にある未装着オーブをすべて分解する (一括分解: %s x%d)", res, scrolls)
		},
		OrbDisassembleOpt0: "0: キャンセル",
		OrbDisassembleAllSuccess: func(usedCount int, res string, scrolls, total int) string {
			return fmt.Sprintf(">> 倉庫の未装着オーブ %d個 をすべて分解し、%s x%d を獲得しました！ (所持: %s x%d)", usedCount, res, scrolls, res, total)
		},
		OrbDisassembleRemainder: func(orb string) string {
			return fmt.Sprintf(">> ※ 分解できなかった余りのオーブ [%s] 1個は倉庫に残りました。", orb)
		},
		OrbDisassembleIndivTitle: "\n--- オーブの個別分解 ---",
		OrbDisassemblePick1: "1つ目に分解するオーブの番号を選択: ",
		OrbDisassembleSelected1: func(orb string) string {
			return fmt.Sprintf(">> 1つ目: [%s] を選択しました。", orb)
		},
		OrbDisassemblePick2: "2つ目に分解するオーブの番号を選択: ",
		OrbDisassembleSameError: ">> 1つ目と同じオーブは選択できません！",
		OrbDisassembleSuccess: func(orb1, orb2, res string, total int) string {
			return fmt.Sprintf(">> [%s] と [%s] を分解し、%s x1 を獲得しました！ (所持: %s x%d)", orb1, orb2, res, res, total)
		},
		OrbDisassembleCanceled: ">> 分解をキャンセルしました。",

		// Orb Synthesize
		OrbSynthNoRecipes: ">> 合成可能なオーブ（同種の通常オーブ2個以上）が倉庫にありません！",
		OrbSynthTitle: "\n--- オーブの合成 (同種の通常オーブ2個 -> 上位オーブ1個) ---",
		OrbSynthRecipeItem: func(idx int, from string, count int, to, toDesc string) string {
			return fmt.Sprintf("%d: [%s] (所持: %d個) -> [%s] %s を合成", idx, from, count, to, toDesc)
		},
		OrbSynthPrompt: "合成するオーブの番号を選択: ",
		OrbSynthSuccess: func(from, to string) string {
			return fmt.Sprintf(">> [%s] を2個消費し、上位オーブ [%s] を合成しました！", from, to)
		},

		// Dungeon Exploration
		DungeonEnter: func(name string) string {
			return fmt.Sprintf("\n>>> 【%s】に突入しました！ <<<", name)
		},
		DungeonEncounter: func(floor, total int, name string) string {
			return fmt.Sprintf("\n==============================================\n   B%dF / B%dF : %s が現れた！\n==============================================", floor, total, name)
		},
		DungeonDefeatPrompt: "\n[!] あなたは力尽きた...",
		DungeonItemsLost: func(res string, scrolls, orbs int) string {
			return fmt.Sprintf("[!] 今回獲得したアイテム（%s: %d個, オーブ: %d個）は失われました。", res, scrolls, orbs)
		},
		DungeonCarriedBack: "[!] 命からがら拠点へ運ばれました。（所持装備は無事です）",
		DungeonVictory: func(enemy string) string {
			return fmt.Sprintf("\n>> %s を討伐！", enemy)
		},
		DungeonDropScroll: func(res string, count int) string {
			return fmt.Sprintf("   戦利品獲得: %s x%d", res, count)
		},
		DungeonDropOrb: func(orb, name string) string {
			return fmt.Sprintf("   戦利品獲得: オーブ [%s] (%s)", orb, name)
		},
		DungeonComplete: func(name string) string {
			return fmt.Sprintf("\n**********************************************\n   【%s】完全踏破！おめでとうございます！   \n**********************************************", name)
		},
		DungeonCurrentStatus: func(hp, maxHp int, res string, scrolls, orbs int) string {
			return fmt.Sprintf("\n現在HP: %d/%d\n現在の未確定戦利品: %s x%d, オーブ x%d", hp, maxHp, res, scrolls, orbs)
		},
		DungeonNextFloor: "1: 次の階層へ進む",
		DungeonRetreat: "2: 撤退する (戦利品を持ち帰って拠点に戻る)",
		DungeonPromptAction: "行動を選択: ",
		DungeonRetreatSuccess: func(res string, scrolls, orbs int) string {
			return fmt.Sprintf(">> 戦利品（%s: %d個, オーブ: %d個）を持ち帰りました！", res, scrolls, orbs)
		},

		// Endless
		EndlessEnter: func(name string) string {
			return fmt.Sprintf("\n>>> 【%s】に突入しました！ <<<", name)
		},
		EndlessIntro: "※ 階層上限のないエンドレスモードです。どこまで潜れるか挑戦しましょう！",
		EndlessStartFloorTitle: "\n--- スタート階層選択 ---",
		EndlessHighestRecord: func(floor int) string {
			return fmt.Sprintf("最高到達階層: B%dF", floor)
		},
		EndlessStartOption: func(floor int) string {
			return fmt.Sprintf("B%dF からスタート", floor)
		},
		EndlessCancelOption: "0: 出撃をキャンセル (拠点に戻る)",
		EndlessStartFrom: func(floor int) string {
			return fmt.Sprintf("\n>> B%dF から探索を開始します！", floor)
		},
		EndlessFloorEncounter: func(floor int, name string) string {
			return fmt.Sprintf("\n==============================================\n   B%dF : %s が現れた！\n==============================================", floor, name)
		},
		EndlessRecordUpdated: func(floor int) string {
			return fmt.Sprintf("★ 最高到達階層を更新！ (B%dF)", floor)
		},

		// Battle
		BattleVs: func(php, pmax int, ename string, ehp, emax, poison int) string {
			poisonStr := ""
			if poison > 0 {
				poisonStr = fmt.Sprintf(" (毒:%d)", poison)
			}
			return fmt.Sprintf("\n[YOU] HP: %d/%d  vs  [%s] HP: %d/%d%s", php, pmax, ename, ehp, emax, poisonStr)
		},
		BattleCmdAttack: "1: 攻撃する",
		BattleCmdRetreat: "2: 撤退する (戦闘から逃げて拠点へ)",
		BattleCmdAuto: "3: オート戦闘",
		BattleCmdPrompt: "コマンド: ",
		BattleRetreatCombat: "\n>> 戦闘から離脱し、命からがら帰還しました。（戦利品は持ち帰れません）",
		BattleAutoStart: "\n>> [オート戦闘開始] 高速で戦闘を進行します...",
		BattleAutoSummaryTitle: "\n==============================================\n             【オート戦闘 終了要約】            \n==============================================",
		BattleAutoReason: func(reason string) string {
			return fmt.Sprintf("- 終了理由: %s", reason)
		},
		BattleAutoTurns: func(turns int) string {
			return fmt.Sprintf("- 経過ターン数: %d ターン", turns)
		},
		BattleAutoDamageDealt: func(dmg int) string {
			return fmt.Sprintf("- 与えた総ダメージ: %d", dmg)
		},
		BattleAutoDamageTaken: func(dmg int) string {
			return fmt.Sprintf("- 受けた総ダメージ: %d", dmg)
		},
		BattleAutoHealed: func(heal int) string {
			return fmt.Sprintf("- 吸血総回復量: %d", heal)
		},
		BattleAutoRemainingHp: func(php, pmax int, ename string, ehp, emax int) string {
			remPlayer := php
			if remPlayer < 0 {
				remPlayer = 0
			}
			remEnemy := ehp
			if remEnemy < 0 {
				remEnemy = 0
			}
			return fmt.Sprintf("- 残りHP: あなた %d/%d | %s %d/%d", remPlayer, pmax, ename, remEnemy, emax)
		},
		BattleAutoReasonEnemyDefeated: "敵を撃破！",
		BattleAutoReasonPoisonDefeated: "毒ダメージにより敵を撃破！",
		BattleAutoReasonFainted: "力尽きました...",
		BattleAutoReasonDangerHp: func(hp, maxHp int) string {
			return fmt.Sprintf("危険域 (HP %d/%d <= 30%%) に到達したため自動停止", hp, maxHp)
		},
		BattleInvalidCmd: ">> 無効なコマンドです。1, 2, 3 のいずれかを入力してください。",
		BattlePlayerAttack: func(hitIndex string, isCrit bool, ename string, dmg int) string {
			critStr := ""
			if isCrit {
				critStr = "【会心の一撃！】 "
			}
			return fmt.Sprintf(">> あなたの%s攻撃！ %s%sに %d ダメージ！", hitIndex, critStr, ename, dmg)
		},
		BattleCritLabel: "【会心の一撃！】 ",
		BattleHitNumber: func(hit int) string {
			return fmt.Sprintf("%d撃目の", hit)
		},
		BattleVampHeal: func(heal, hp int) string {
			return fmt.Sprintf("   [吸血] HPが %d 回復した！ (現在HP: %d)", heal, hp)
		},
		BattlePoisonInflict: func(ename string, added, total int) string {
			return fmt.Sprintf("   [猛毒] %sに毒を付与！ (+%d / 毒カウント: %d)", ename, added, total)
		},
		BattlePoisonTick: func(ename string, dmg int) string {
			return fmt.Sprintf(">> [毒効果] %sは毒で %d のダメージを受けた！", ename, dmg)
		},
		BattleEnemyCounter: func(ename string, dmg int) string {
			return fmt.Sprintf(">> %sの反撃！ あなたは %d のダメージを受けた！", ename, dmg)
		},

		// Settings
		SettingsTitle: "\n--- 設定 (Settings) ---",
		SettingsCurrent: func(lang, theme string) string {
			langLabel := "English (en)"
			if lang == "ja" {
				langLabel = "日本語 (ja)"
			}
			return fmt.Sprintf("現在の設定: 言語: [%s] | テーマ: [%s]", langLabel, theme)
		},
		SettingsMenu1Lang: "1: 言語切り替え (Toggle Language: en / ja)",
		SettingsMenu2Theme: "2: テーマ切り替え (Switch Theme)",
		SettingsMenu3ExportTemplate: "3: カスタムテーマ雛形を出力 (Export custom_theme.example.json)",
		SettingsMenu0Back: "0: 戻る (Back)",
		SettingsLangChanged: func(lang string) string {
			langLabel := "English (en)"
			if lang == "ja" {
				langLabel = "日本語 (ja)"
			}
			return fmt.Sprintf(">> 言語を「%s」に切り替えました。", langLabel)
		},
		SettingsThemeSelectTitle: "\n--- テーマ選択 ---",
		SettingsThemeChanged: func(theme string) string {
			return fmt.Sprintf(">> テーマを「%s」に切り替えました。", theme)
		},
		SettingsTemplateExported: func(path string) string {
			return fmt.Sprintf(">> カスタムテーマ雛形を出力しました: %s", path)
		},
		SettingsTemplateExportFailed: ">> 雛形の出力に失敗しました:",
	}

	MessagesEN = Messages{
		Title: "==============================================\n   Minimal Rogue-lite Prototype (CUI Ver)   \n==============================================",
		StartPrompt: "\nPlease choose an option: ",
		StartContinue: "1: Continue (Load save.json)",
		StartNew: "2: New Game (Start with initial stats)",
		InvalidChoice1or2: ">> Invalid choice. Please enter 1 or 2.",
		SaveLoaded: ">> Save data loaded successfully!",
		SaveStats: func(hp, deepest int) string {
			rec := "None"
			if deepest > 0 {
				rec = fmt.Sprintf("B%dF", deepest)
			}
			return fmt.Sprintf(">> Max HP: %d | Deepest Floor: %s", hp, rec)
		},
		Untested: "None",
		SaveNotFoundOrCorrupted: ">> save.json was not found or corrupted. Starting a new game.",
		NewGameStarted: ">> Starting a new game.",
		SavedSuccess: ">> Game saved successfully (save.json).",
		SaveFailed: ">> Failed to save game:",
		LoadFailed: ">> Failed to load save data:",
		Farewell: "Thank you for playing!",

		// Hub
		HubHeader: func(title string) string { return title },
		HubStats: func(hpLabel string, hp int, abyssName string, deepest int) string {
			rec := "None"
			if deepest > 0 {
				rec = fmt.Sprintf("B%dF", deepest)
			}
			return fmt.Sprintf("%s: %d | %s Record: %s", hpLabel, hp, abyssName, rec)
		},
		HubEquip: func(label, name, plusPrefix string, plus int, statLabel string, atk int) string {
			return fmt.Sprintf("%s: %s%s%d (%s: %d)", label, name, plusPrefix, plus, statLabel, atk)
		},
		HubSlots: func(verb string, count int, slots string) string {
			return fmt.Sprintf("%s Slots [%d/3]: %s", verb, count, slots)
		},
		HubEmptySlot: "(Empty)",
		HubStorage: func(storage, resource string, scrolls int, orb, orbs string) string {
			return fmt.Sprintf("%s: %s x%d | Stock %ss: %s", storage, resource, scrolls, orb, orbs)
		},
		HubNone: "(None)",
		HubMenu1Dungeon: "1: Embark to Dungeon (Start Run)",
		HubMenu2Enhance: func(verb, statLabel string) string { return fmt.Sprintf("2: %s (Upgrade %s)", verb, statLabel) },
		HubMenu3AttachOrb: func(verb string) string {
			return fmt.Sprintf("3: %s (Equip Passives)", verb)
		},
		HubMenu4Disassemble: func(orb string) string {
			return fmt.Sprintf("4: Dismantle %s (Convert to Materials)", orb)
		},
		HubMenu5Synthesize: func(orb string) string {
			return fmt.Sprintf("5: Synthesize %s (Upgrade to Plus)", orb)
		},
		HubMenu6EnhanceHp: func(verb, hpLabel string) string {
			return fmt.Sprintf("6: %s (%s +10)", verb, hpLabel)
		},
		HubMenu7Settings: "7: Settings",
		HubMenu0Exit: "0: Quit Game",
		ChooseAction: "\nChoose an action: ",
		InvalidChoice: "Invalid choice.",

		// Dungeon Select
		DungeonSelectTitle: "\n--- Select Dungeon ---",
		DungeonOptionStarter: func(name string, floors int) string {
			return fmt.Sprintf("1: %s (%d Floors / Beginner)", name, floors)
		},
		DungeonOptionDeep: func(name string, floors int, rec string) string {
			return fmt.Sprintf("2: %s (%d Floors / %s)", name, floors, rec)
		},
		DungeonOptionAbyss: func(name string) string {
			return fmt.Sprintf("3: %s (Endless / Unlimited Floors)", name)
		},
		DungeonCancel: "0: Cancel (Return to Base)",
		DungeonPrompt: "\nSelect a dungeon: ",
		DungeonCanceled: ">> Canceled departure.",
		RecommendedDeep: "Recommended +25 or higher",

		// Weapon / Target Upgrade
		EnhanceTitle: func(label string) string {
			return fmt.Sprintf("\n--- Enhance %s ---", label)
		},
		EnhanceCurrentRes: func(res string, count int) string {
			return fmt.Sprintf("Current %s: %d", res, count)
		},
		EnhanceOpt1: "1: Upgrade once (+1)",
		EnhanceOpt2: func(res string) string {
			return fmt.Sprintf("2: Upgrade multiple times (Input amount of %s)", res)
		},
		EnhanceOpt3: func(res string, count int) string {
			return fmt.Sprintf("3: Upgrade with all available %s (+%d)", res, count)
		},
		EnhanceOpt0: "0: Cancel",
		EnhanceNoResource: func(res string) string {
			return fmt.Sprintf(">> Not enough %s!", res)
		},
		EnhancePromptCount: func(res string, max int) string {
			return fmt.Sprintf("Enter amount of %s to consume (1 to %d): ", res, max)
		},
		EnhanceInvalidCount: ">> Invalid amount. Please enter a valid number within your inventory.",
		EnhanceSuccessSingle: func(verb, label, name string, plus, atk int) string {
			return fmt.Sprintf(">> [%s] Performed %s! %s+%d (ATK: %d)", label, verb, name, plus, atk)
		},
		EnhanceSuccessMultiple: func(count int, res, verb, name string, plus, atk int) string {
			return fmt.Sprintf(">> Consumed %d %s for %s! %s+%d (ATK: %d)", count, res, verb, name, plus, atk)
		},
		EnhanceSuccessAll: func(count int, res, verb, name string, plus, atk int) string {
			return fmt.Sprintf(">> Consumed all (%d) %s for %s! %s+%d (ATK: %d)", count, res, verb, name, plus, atk)
		},
		EnhanceCanceled: ">> Enhancement canceled.",
		RankUpgraded: func(name string) string {
			return fmt.Sprintf(">> Rank upgraded! Evolved to 【%s】!", name)
		},

		// HP Upgrade
		HPUpgradeTitle: "\n--- Upgrade Max HP ---",
		HPUpgradeCurrentHp: func(hp int) string {
			return fmt.Sprintf("Current Max HP: %d", hp)
		},
		HPUpgradeCurrentRes: func(res string, count int) string {
			return fmt.Sprintf("Current %s: %d (2 %s -> +10 Max HP)", res, count, res)
		},
		HPUpgradeNotEnough: func(res string, count int) string {
			return fmt.Sprintf(">> Not enough %s! (Required: 2 / Owned: %d)", res, count)
		},
		HPUpgradeOpt1: func(res string) string {
			return fmt.Sprintf("1: Upgrade once (Consume 2 %s -> +10 Max HP)", res)
		},
		HPUpgradeOpt2: func(res string, max int) string {
			return fmt.Sprintf("2: Upgrade multiple times (Input amount of %s, max: %d)", res, max)
		},
		HPUpgradeOpt3: func(res string, used, gain int) string {
			return fmt.Sprintf("3: Upgrade with all available %s (Consume %d %s -> +%d Max HP)", res, used, res, gain)
		},
		HPUpgradeOpt0: "0: Cancel",
		HPUpgradePromptCount: func(res string, max int) string {
			return fmt.Sprintf("Enter amount of %s to consume (2 to %d): ", res, max)
		},
		HPUpgradeInvalidCount: ">> Invalid amount. Please enter a number >= 2 within your inventory.",
		HPUpgradeSuccessSingle: func(prev, next int) string {
			return fmt.Sprintf(">> Max HP upgraded! %d -> %d", prev, next)
		},
		HPUpgradeSuccessMultiple: func(used int, res string, prev, next int) string {
			return fmt.Sprintf(">> Consumed %d %s to upgrade Max HP! %d -> %d", used, res, prev, next)
		},
		HPUpgradeOddRemainder: ">> * Consumed in units of 2; 1 remaining was retained.",
		HPUpgradeSuccessAll: func(used int, res string, prev, next int) string {
			return fmt.Sprintf(">> Consumed all (%d) %s to upgrade Max HP! %d -> %d", used, res, prev, next)
		},
		HPUpgradeCanceled: ">> HP upgrade canceled.",

		// Orb Attach
		OrbAttachNoOrbs: ">> No orbs available in storage!",
		OrbAttachTitle: "\n--- Stored Orbs ---",
		OrbAttachPrompt: "Select an orb to attach: ",
		OrbAttachSuccess: func(orb string) string {
			return fmt.Sprintf(">> Attached [%s] to weapon slot!", orb)
		},
		OrbAttachFullTitle: "\nSlots are full. Choose a slot to overwrite:",
		OrbAttachReplacePrompt: "Select slot number: ",
		OrbAttachReplaced: func(removed, added string) string {
			return fmt.Sprintf(">> Discarded [%s] and equipped [%s]!", removed, added)
		},

		// Orb Disassemble
		OrbDisassembleNotEnough: ">> You need at least 2 orbs in storage to disassemble!",
		OrbDisassembleTitle: func(res string) string {
			return fmt.Sprintf("\n--- Disassemble Orbs (Any 2 Orbs -> 1 %s) ---", res)
		},
		OrbDisassembleCount: func(count int) string {
			return fmt.Sprintf("Unassigned Orbs in Storage: %d", count)
		},
		OrbDisassembleOpt1: "1: Select individual orbs to disassemble",
		OrbDisassembleOpt2: func(res string, scrolls int) string {
			return fmt.Sprintf("2: Disassemble all unassigned orbs (Batch: %s x%d)", res, scrolls)
		},
		OrbDisassembleOpt0: "0: Cancel",
		OrbDisassembleAllSuccess: func(usedCount int, res string, scrolls, total int) string {
			return fmt.Sprintf(">> Disassembled %d orbs and obtained %s x%d! (Total: %s x%d)", usedCount, res, scrolls, res, total)
		},
		OrbDisassembleRemainder: func(orb string) string {
			return fmt.Sprintf(">> * 1 remaining orb [%s] could not be paired and remains in storage.", orb)
		},
		OrbDisassembleIndivTitle: "\n--- Individual Orb Disassembly ---",
		OrbDisassemblePick1: "Select first orb to disassemble: ",
		OrbDisassembleSelected1: func(orb string) string {
			return fmt.Sprintf(">> First orb selected: [%s]", orb)
		},
		OrbDisassemblePick2: "Select second orb to disassemble: ",
		OrbDisassembleSameError: ">> Cannot choose the same orb twice!",
		OrbDisassembleSuccess: func(orb1, orb2, res string, total int) string {
			return fmt.Sprintf(">> Disassembled [%s] and [%s], obtained %s x1! (Total: %s x%d)", orb1, orb2, res, res, total)
		},
		OrbDisassembleCanceled: ">> Disassembly canceled.",

		// Orb Synthesize
		OrbSynthNoRecipes: ">> No synthesizable orbs (need at least 2 of the same normal orb) in storage!",
		OrbSynthTitle: "\n--- Synthesize Orbs (2 Same Normal Orbs -> 1 Plus Orb) ---",
		OrbSynthRecipeItem: func(idx int, from string, count int, to, toDesc string) string {
			return fmt.Sprintf("%d: [%s] (Owned: %d) -> Synthesize [%s] %s", idx, from, count, to, toDesc)
		},
		OrbSynthPrompt: "Select recipe number: ",
		OrbSynthSuccess: func(from, to string) string {
			return fmt.Sprintf(">> Consumed 2x [%s] and synthesized [%s]!", from, to)
		},

		// Dungeon Exploration
		DungeonEnter: func(name string) string {
			return fmt.Sprintf("\n>>> Entered 【%s】! <<<", name)
		},
		DungeonEncounter: func(floor, total int, name string) string {
			return fmt.Sprintf("\n==============================================\n   B%dF / B%dF : %s appeared!\n==============================================", floor, total, name)
		},
		DungeonDefeatPrompt: "\n[!] You have fallen...",
		DungeonItemsLost: func(res string, scrolls, orbs int) string {
			return fmt.Sprintf("[!] Loot gathered on this expedition (%s: %d, Orbs: %d) was lost.", res, scrolls, orbs)
		},
		DungeonCarriedBack: "[!] Carried back to base safely. (Equipped items preserved)",
		DungeonVictory: func(enemy string) string {
			return fmt.Sprintf("\n>> Defeated %s!", enemy)
		},
		DungeonDropScroll: func(res string, count int) string {
			return fmt.Sprintf("   Loot obtained: %s x%d", res, count)
		},
		DungeonDropOrb: func(orb, name string) string {
			return fmt.Sprintf("   Loot obtained: Orb [%s] (%s)", orb, name)
		},
		DungeonComplete: func(name string) string {
			return fmt.Sprintf("\n**********************************************\n   【%s】 Cleared! Congratulations!   \n**********************************************", name)
		},
		DungeonCurrentStatus: func(hp, maxHp int, res string, scrolls, orbs int) string {
			return fmt.Sprintf("\nCurrent HP: %d/%d\nPending Loot: %s x%d, Orbs x%d", hp, maxHp, res, scrolls, orbs)
		},
		DungeonNextFloor: "1: Proceed to next floor",
		DungeonRetreat: "2: Retreat (Return to base with collected loot)",
		DungeonPromptAction: "Choose action: ",
		DungeonRetreatSuccess: func(res string, scrolls, orbs int) string {
			return fmt.Sprintf(">> Brought back loot (%s: %d, Orbs: %d) safely!", res, scrolls, orbs)
		},

		// Endless
		EndlessEnter: func(name string) string {
			return fmt.Sprintf("\n>>> Entered 【%s】! <<<", name)
		},
		EndlessIntro: "* Endless dungeon with no floor limit. Test how deep you can survive!",
		EndlessStartFloorTitle: "\n--- Select Start Floor ---",
		EndlessHighestRecord: func(floor int) string {
			return fmt.Sprintf("Highest Record: B%dF", floor)
		},
		EndlessStartOption: func(floor int) string {
			return fmt.Sprintf("Start from B%dF", floor)
		},
		EndlessCancelOption: "0: Cancel departure (Return to base)",
		EndlessStartFrom: func(floor int) string {
			return fmt.Sprintf("\n>> Starting expedition from B%dF!", floor)
		},
		EndlessFloorEncounter: func(floor int, name string) string {
			return fmt.Sprintf("\n==============================================\n   B%dF : %s appeared!\n==============================================", floor, name)
		},
		EndlessRecordUpdated: func(floor int) string {
			return fmt.Sprintf("★ New highest record reached! (B%dF)", floor)
		},

		// Battle
		BattleVs: func(php, pmax int, ename string, ehp, emax, poison int) string {
			poisonStr := ""
			if poison > 0 {
				poisonStr = fmt.Sprintf(" (Poison:%d)", poison)
			}
			return fmt.Sprintf("\n[YOU] HP: %d/%d  vs  [%s] HP: %d/%d%s", php, pmax, ename, ehp, emax, poisonStr)
		},
		BattleCmdAttack: "1: Attack",
		BattleCmdRetreat: "2: Retreat (Flee to base)",
		BattleCmdAuto: "3: Auto Combat",
		BattleCmdPrompt: "Command: ",
		BattleRetreatCombat: "\n>> Fled from battle and escaped back to base. (Loot lost)",
		BattleAutoStart: "\n>> [Auto Combat Started] Resolving combat at high speed...",
		BattleAutoSummaryTitle: "\n==============================================\n          【Auto Combat Summary Report】         \n==============================================",
		BattleAutoReason: func(reason string) string {
			return fmt.Sprintf("- End Reason: %s", reason)
		},
		BattleAutoTurns: func(turns int) string {
			return fmt.Sprintf("- Elapsed Turns: %d turns", turns)
		},
		BattleAutoDamageDealt: func(dmg int) string {
			return fmt.Sprintf("- Total Damage Dealt: %d", dmg)
		},
		BattleAutoDamageTaken: func(dmg int) string {
			return fmt.Sprintf("- Total Damage Taken: %d", dmg)
		},
		BattleAutoHealed: func(heal int) string {
			return fmt.Sprintf("- Total HP Restored via Vampirism: %d", heal)
		},
		BattleAutoRemainingHp: func(php, pmax int, ename string, ehp, emax int) string {
			remPlayer := php
			if remPlayer < 0 {
				remPlayer = 0
			}
			remEnemy := ehp
			if remEnemy < 0 {
				remEnemy = 0
			}
			return fmt.Sprintf("- Remaining HP: You %d/%d | %s %d/%d", remPlayer, pmax, ename, remEnemy, emax)
		},
		BattleAutoReasonEnemyDefeated: "Enemy Defeated!",
		BattleAutoReasonPoisonDefeated: "Enemy slain by poison damage!",
		BattleAutoReasonFainted: "You were knocked out...",
		BattleAutoReasonDangerHp: func(hp, maxHp int) string {
			return fmt.Sprintf("Auto stopped upon reaching danger zone (HP %d/%d <= 30%%)", hp, maxHp)
		},
		BattleInvalidCmd: ">> Invalid command. Please enter 1, 2, or 3.",
		BattlePlayerAttack: func(hitIndex string, isCrit bool, ename string, dmg int) string {
			critStr := ""
			if isCrit {
				critStr = "【Critical Strike!】 "
			}
			return fmt.Sprintf(">> Your %sattack! %s%s took %d damage!", hitIndex, critStr, ename, dmg)
		},
		BattleCritLabel: "【Critical Strike!】 ",
		BattleHitNumber: func(hit int) string {
			return fmt.Sprintf("strike #%d ", hit)
		},
		BattleVampHeal: func(heal, hp int) string {
			return fmt.Sprintf("   [Vampiric] Restored %d HP! (Current HP: %d)", heal, hp)
		},
		BattlePoisonInflict: func(ename string, added, total int) string {
			return fmt.Sprintf("   [Poison] Inflicted poison to %s! (+%d / Poison count: %d)", ename, added, total)
		},
		BattlePoisonTick: func(ename string, dmg int) string {
			return fmt.Sprintf(">> [Poison] %s took %d poison damage!", ename, dmg)
		},
		BattleEnemyCounter: func(ename string, dmg int) string {
			return fmt.Sprintf(">> %s's counterattack! You took %d damage!", ename, dmg)
		},

		// Settings
		SettingsTitle: "\n--- Settings ---",
		SettingsCurrent: func(lang, theme string) string {
			langLabel := "English (en)"
			if lang == "ja" {
				langLabel = "日本語 (ja)"
			}
			return fmt.Sprintf("Current Settings: Language: [%s] | Theme: [%s]", langLabel, theme)
		},
		SettingsMenu1Lang: "1: Toggle Language (en / ja)",
		SettingsMenu2Theme: "2: Switch Theme",
		SettingsMenu3ExportTemplate: "3: Export Custom Theme Template (custom_theme.example.json)",
		SettingsMenu0Back: "0: Back",
		SettingsLangChanged: func(lang string) string {
			langLabel := "English (en)"
			if lang == "ja" {
				langLabel = "日本語 (ja)"
			}
			return fmt.Sprintf(">> Language switched to \"%s\".", langLabel)
		},
		SettingsThemeSelectTitle: "\n--- Select Theme ---",
		SettingsThemeChanged: func(theme string) string {
			return fmt.Sprintf(">> Switched theme to \"%s\".", theme)
		},
		SettingsTemplateExported: func(path string) string {
			return fmt.Sprintf(">> Exported custom theme template: %s", path)
		},
		SettingsTemplateExportFailed: ">> Failed to export template:",
	}
)

// GetMessages returns the dictionary for the specified language.
func GetMessages(lang model.Language) Messages {
	if lang == model.LanguageJA {
		return MessagesJA
	}
	return MessagesEN
}

// DetectDefaultLanguage checks environment variables for Japanese locale.
func DetectDefaultLanguage() model.Language {
	for _, key := range []string{"LANG", "LC_ALL", "LC_MESSAGES"} {
		val := strings.ToLower(os.Getenv(key))
		if strings.HasPrefix(val, "ja") || strings.Contains(val, "ja_jp") {
			return model.LanguageJA
		}
	}
	return model.LanguageEN
}
