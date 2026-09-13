import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";
import { argv, stdin as input, stdout as output } from "node:process";
import * as readline from "node:readline/promises";
import { fileURLToPath } from "node:url";

const SAVE_FILE_PATH = "save.json";
const CUSTOM_THEME_FILE_PATH = "custom_theme.json";
const CUSTOM_THEME_EXAMPLE_PATH = "custom_theme.example.json";

// ============================================================================
// Types
// ============================================================================

export type Language = "en" | "ja";

export type NormalOrbType = "MULTI_HIT" | "CRITICAL" | "VAMP" | "POISON";
export type PlusOrbType =
	| "MULTI_HIT_PLUS"
	| "CRITICAL_PLUS"
	| "VAMP_PLUS"
	| "POISON_PLUS";
export type OrbType = NormalOrbType | PlusOrbType;

export interface Weapon {
	name: string;
	baseAtk: number;
	plus: number;
	slots: OrbType[];
}

export interface Player {
	maxHp: number;
	hp: number;
	weapon: Weapon;
}

export interface Enemy {
	name: string;
	hp: number;
	maxHp: number;
	atk: number;
	poison: number;
}

export interface RunInventory {
	scrolls: number;
	orbs: OrbType[];
}

export interface SaveData {
	weapon: Weapon;
	stockScrolls: number;
	stockOrbs: OrbType[];
	maxHp?: number | undefined;
	deepestFloor?: number | undefined;
	language?: Language | undefined;
	themeId?: string | undefined;
}

export interface EnemyTemplate {
	name: string;
	hp: number;
	atk: number;
	getDrops: () => RunInventory;
}

export interface DungeonDef {
	name: string;
	floors: number;
	recommended?: string | undefined;
	getEnemy: (floor: number) => EnemyTemplate;
}

export interface ThemeLocaleData {
	name: string;
	targetNameLabel: string;
	defaultTargetName: string;
	enhanceVerb: string;
	resourceName: string;
	orbs: Record<OrbType, { name: string; desc: string }>;
	dungeons: { starter: string; deep: string; abyss: string };
	enemyPrefixes: string[];
	enemyBases: string[];
	retreatMessage: string;
}

export interface ThemeDefinition {
	id: string;
	ja: ThemeLocaleData;
	en: ThemeLocaleData;
}

export interface Theme extends ThemeLocaleData {
	id: string;
	language: Language;
}

export function resolveTheme(def: ThemeDefinition, lang: Language): Theme {
	const locale = def[lang] ?? def.ja ?? def.en;
	return {
		id: def.id,
		language: lang,
		...locale,
	};
}

// ============================================================================
// System Message Layer (i18n)
// ============================================================================

export interface I18nMessages {
	title: string;
	startPrompt: string;
	startContinue: string;
	startNew: string;
	invalidChoice1or2: string;
	saveLoaded: string;
	saveStats: (hp: number, deepest: number) => string;
	untested: string;
	saveNotFoundOrCorrupted: string;
	newGameStarted: string;
	savedSuccess: string;
	saveFailed: string;
	loadFailed: string;
	farewell: string;

	// Hub
	hubHeader: string;
	hubStats: (hp: number, deepest: number) => string;
	hubEquip: (label: string, name: string, plus: number, atk: number) => string;
	hubSlots: (count: number, slots: string) => string;
	hubEmptySlot: string;
	hubStorage: (resource: string, scrolls: number, orbs: string) => string;
	hubNone: string;
	hubMenu1Dungeon: string;
	hubMenu2Enhance: (verb: string) => string;
	hubMenu3AttachOrb: (label: string) => string;
	hubMenu4Disassemble: (resource: string) => string;
	hubMenu5Synthesize: string;
	hubMenu6EnhanceHp: (resource: string) => string;
	hubMenu7Settings: string;
	hubMenu0Exit: string;
	chooseAction: string;
	invalidChoice: string;

	// Dungeon Select
	dungeonSelectTitle: string;
	dungeonOptionStarter: (name: string, floors: number) => string;
	dungeonOptionDeep: (
		name: string,
		floors: number,
		recommended: string,
	) => string;
	dungeonOptionAbyss: (name: string) => string;
	dungeonCancel: string;
	dungeonPrompt: string;
	dungeonCanceled: string;
	recommendedDeep: string;

	// Weapon / Target Upgrade
	enhanceTitle: (label: string) => string;
	enhanceCurrentRes: (resource: string, count: number) => string;
	enhanceOpt1: string;
	enhanceOpt2: (resource: string) => string;
	enhanceOpt3: (resource: string, count: number) => string;
	enhanceOpt0: string;
	enhanceNoResource: (resource: string) => string;
	enhancePromptCount: (resource: string, max: number) => string;
	enhanceInvalidCount: string;
	enhanceSuccessSingle: (
		verb: string,
		label: string,
		name: string,
		plus: number,
		atk: number,
	) => string;
	enhanceSuccessMultiple: (
		count: number,
		resource: string,
		verb: string,
		name: string,
		plus: number,
		atk: number,
	) => string;
	enhanceSuccessAll: (
		count: number,
		resource: string,
		verb: string,
		name: string,
		plus: number,
		atk: number,
	) => string;
	enhanceCanceled: string;

	// HP Upgrade
	hpUpgradeTitle: string;
	hpUpgradeCurrentHp: (hp: number) => string;
	hpUpgradeCurrentRes: (resource: string, count: number) => string;
	hpUpgradeNotEnough: (resource: string, count: number) => string;
	hpUpgradeOpt1: (resource: string) => string;
	hpUpgradeOpt2: (resource: string, max: number) => string;
	hpUpgradeOpt3: (resource: string, used: number, gain: number) => string;
	hpUpgradeOpt0: string;
	hpUpgradePromptCount: (resource: string, max: number) => string;
	hpUpgradeInvalidCount: string;
	hpUpgradeSuccessSingle: (prev: number, next: number) => string;
	hpUpgradeSuccessMultiple: (
		used: number,
		resource: string,
		prev: number,
		next: number,
	) => string;
	hpUpgradeOddRemainder: string;
	hpUpgradeSuccessAll: (
		used: number,
		resource: string,
		prev: number,
		next: number,
	) => string;
	hpUpgradeCanceled: string;

	// Orb Attach
	orbAttachNoOrbs: string;
	orbAttachTitle: string;
	orbAttachPrompt: string;
	orbAttachSuccess: (orb: string) => string;
	orbAttachFullTitle: string;
	orbAttachReplacePrompt: string;
	orbAttachReplaced: (removed: string, added: string) => string;

	// Orb Disassemble
	orbDisassembleNotEnough: string;
	orbDisassembleTitle: (resource: string) => string;
	orbDisassembleCount: (count: number) => string;
	orbDisassembleOpt1: string;
	orbDisassembleOpt2: (resource: string, scrolls: number) => string;
	orbDisassembleOpt0: string;
	orbDisassembleAllSuccess: (
		usedCount: number,
		resource: string,
		scrolls: number,
		totalScrolls: number,
	) => string;
	orbDisassembleRemainder: (orb: string) => string;
	orbDisassembleIndivTitle: string;
	orbDisassemblePick1: string;
	orbDisassembleSelected1: (orb: string) => string;
	orbDisassemblePick2: string;
	orbDisassembleSameError: string;
	orbDisassembleSuccess: (
		orb1: string,
		orb2: string,
		resource: string,
		totalScrolls: number,
	) => string;
	orbDisassembleCanceled: string;

	// Orb Synthesize
	orbSynthNoRecipes: string;
	orbSynthTitle: string;
	orbSynthRecipeItem: (
		index: number,
		from: string,
		count: number,
		to: string,
		toDesc: string,
	) => string;
	orbSynthPrompt: string;
	orbSynthSuccess: (from: string, to: string) => string;

	// Dungeon Exploration
	dungeonEnter: (name: string) => string;
	dungeonEncounter: (floor: number, total: number, name: string) => string;
	dungeonDefeatPrompt: string;
	dungeonItemsLost: (resource: string, scrolls: number, orbs: number) => string;
	dungeonCarriedBack: string;
	dungeonVictory: (enemy: string) => string;
	dungeonDropScroll: (resource: string, count: number) => string;
	dungeonDropOrb: (orb: string, name: string) => string;
	dungeonComplete: (name: string) => string;
	dungeonCurrentStatus: (
		hp: number,
		maxHp: number,
		resource: string,
		scrolls: number,
		orbs: number,
	) => string;
	dungeonNextFloor: string;
	dungeonRetreat: string;
	dungeonPromptAction: string;
	dungeonRetreatSuccess: (
		resource: string,
		scrolls: number,
		orbs: number,
	) => string;

	// Endless
	endlessEnter: (name: string) => string;
	endlessIntro: string;
	endlessStartFloorTitle: string;
	endlessHighestRecord: (floor: number) => string;
	endlessStartOption: (floor: number) => string;
	endlessCancelOption: string;
	endlessStartFrom: (floor: number) => string;
	endlessFloorEncounter: (floor: number, name: string) => string;
	endlessRecordUpdated: (floor: number) => string;

	// Battle
	battleVs: (
		playerHp: number,
		maxHp: number,
		enemyName: string,
		enemyHp: number,
		enemyMaxHp: number,
		poison: number,
	) => string;
	battleCmdAttack: string;
	battleCmdRetreat: string;
	battleCmdAuto: string;
	battleCmdPrompt: string;
	battleRetreatCombat: string;
	battleAutoStart: string;
	battleAutoSummaryTitle: string;
	battleAutoReason: (reason: string) => string;
	battleAutoTurns: (turns: number) => string;
	battleAutoDamageDealt: (dmg: number) => string;
	battleAutoDamageTaken: (dmg: number) => string;
	battleAutoHealed: (heal: number) => string;
	battleAutoRemainingHp: (
		playerHp: number,
		maxHp: number,
		enemyName: string,
		enemyHp: number,
		enemyMaxHp: number,
	) => string;
	battleAutoReasonEnemyDefeated: string;
	battleAutoReasonPoisonDefeated: string;
	battleAutoReasonFainted: string;
	battleAutoReasonDangerHp: (hp: number, maxHp: number) => string;
	battleInvalidCmd: string;
	battlePlayerAttack: (
		hitIndex: string,
		isCrit: boolean,
		enemyName: string,
		dmg: number,
	) => string;
	battleCritLabel: string;
	battleHitNumber: (hit: number) => string;
	battleVampHeal: (heal: number, hp: number) => string;
	battlePoisonInflict: (
		enemyName: string,
		added: number,
		total: number,
	) => string;
	battlePoisonTick: (enemyName: string, dmg: number) => string;
	battleEnemyCounter: (enemyName: string, dmg: number) => string;

	// Settings
	settingsTitle: string;
	settingsCurrent: (lang: string, theme: string) => string;
	settingsMenu1Lang: string;
	settingsMenu2Theme: string;
	settingsMenu3ExportTemplate: string;
	settingsMenu0Back: string;
	settingsLangChanged: (lang: string) => string;
	settingsThemeSelectTitle: string;
	settingsThemeChanged: (theme: string) => string;
	settingsTemplateExported: (path: string) => string;
	settingsTemplateExportFailed: string;
}

export const MESSAGES: Record<Language, I18nMessages> = {
	ja: {
		title:
			"==============================================\n   Minimal Rogue-lite Prototype (CUI Ver)   \n==============================================",
		startPrompt: "\n選択してください: ",
		startContinue: "1: つづきから (save.json を読み込んで開始)",
		startNew: "2: はじめから (初期状態で開始)",
		invalidChoice1or2: ">> 無効な選択です。1 または 2 を入力してください。",
		saveLoaded: ">> セーブデータを読み込みました！",
		saveStats: (hp, deepest) =>
			`>> 最大体力: HP ${hp} | 最高到達階層: ${deepest > 0 ? `B${deepest}F` : "未挑戦"}`,
		untested: "未挑戦",
		saveNotFoundOrCorrupted:
			">> save.json が見つからないか破損しています。新規データで開始します。",
		newGameStarted: ">> はじめからゲームを開始します。",
		savedSuccess: ">> セーブデータを保存しました。(save.json)",
		saveFailed: ">> セーブデータの保存に失敗しました:",
		loadFailed: ">> セーブデータの読み込みに失敗しました:",
		farewell: "お疲れ様でした。",

		// Hub
		hubHeader: "【拠点】",
		hubStats: (hp, deepest) =>
			`最大体力: HP ${hp} | 無限の深淵 最高到達: ${deepest > 0 ? `B${deepest}F` : "未挑戦"}`,
		hubEquip: (label, name, plus, atk) =>
			`${label}: ${name}+${plus} (攻撃力: ${atk})`,
		hubSlots: (count, slots) => `装着スロット [${count}/3]: ${slots}`,
		hubEmptySlot: "(空き)",
		hubStorage: (resource, scrolls, orbs) =>
			`倉庫: ${resource} x${scrolls} | 未装着オーブ: ${orbs}`,
		hubNone: "(なし)",
		hubMenu1Dungeon: "1: ダンジョンへ出撃",
		hubMenu2Enhance: (verb) => `2: ${verb}`,
		hubMenu3AttachOrb: (label) => `3: オーブを${label}に装着`,
		hubMenu4Disassemble: (res) =>
			`4: オーブを分解 (任意のオーブ2個 -> ${res}1枚)`,
		hubMenu5Synthesize: "5: オーブを合成 (同種オーブ2個 -> 上位オーブ)",
		hubMenu6EnhanceHp: (res) => `6: 体力を強化する (${res}2枚 -> 最大HP+10)`,
		hubMenu7Settings: "7: 設定 (Settings)",
		hubMenu0Exit: "0: ゲーム終了",
		chooseAction: "\n行動を選択してください: ",
		invalidChoice: "無効な選択です。",

		// Dungeon Select
		dungeonSelectTitle: "\n--- 出撃ダンジョン選択 ---",
		dungeonOptionStarter: (name, floors) => `1: ${name} (全${floors}階 / 初級)`,
		dungeonOptionDeep: (name, floors, rec) =>
			`2: ${name} (全${floors}階 / ${rec})`,
		dungeonOptionAbyss: (name) => `3: ${name} (エンドレス / 階層無制限)`,
		dungeonCancel: "0: キャンセル (拠点に戻る)",
		dungeonPrompt: "\nダンジョンを選択してください: ",
		dungeonCanceled: ">> 出撃を取りやめました。",
		recommendedDeep: "推奨+25以上の上級ダンジョン",

		// Weapon / Target Upgrade
		enhanceTitle: (label) => `\n--- ${label}の強化 ---`,
		enhanceCurrentRes: (res, count) => `現在の所持${res}: ${count}個`,
		enhanceOpt1: "1: 1回鍛える (+1)",
		enhanceOpt2: (res) => `2: 指定した回数分鍛える（消費する${res}を直接入力）`,
		enhanceOpt3: (res, count) =>
			`3: 所持している${res}ですべて鍛える (+${count})`,
		enhanceOpt0: "0: キャンセル",
		enhanceNoResource: (res) => `>> ${res}がありません！`,
		enhancePromptCount: (res, max) =>
			`消費する${res}の数を入力してください (1〜${max}): `,
		enhanceInvalidCount:
			">> 無効な数値です。1以上の所持数以内の数値を入力してください。",
		enhanceSuccessSingle: (verb, label, name, plus, atk) =>
			`>> 【${label}】${verb}を行いました！ ${name}+${plus} (攻撃力: ${atk})`,
		enhanceSuccessMultiple: (count, res, verb, name, plus, atk) =>
			`>> ${res}を ${count}個 消費して${verb}を行いました！ ${name}+${plus} (攻撃力: ${atk})`,
		enhanceSuccessAll: (count, res, verb, name, plus, atk) =>
			`>> ${res}をすべて(${count}個)消費して一括で${verb}を行いました！ ${name}+${plus} (攻撃力: ${atk})`,
		enhanceCanceled: ">> 強化をキャンセルしました。",

		// HP Upgrade
		hpUpgradeTitle: "\n--- 体力の強化 ---",
		hpUpgradeCurrentHp: (hp) => `現在の最大HP: ${hp}`,
		hpUpgradeCurrentRes: (res, count) =>
			`所持${res}: ${count}個 (2個消費で最大HP+10)`,
		hpUpgradeNotEnough: (res, count) =>
			`>> ${res}が足りません！ (必要: 2個 / 所持: ${count}個)`,
		hpUpgradeOpt1: (res) => `1: 1回強化する (${res}2個消費 -> 最大HP+10)`,
		hpUpgradeOpt2: (res, max) =>
			`2: 指定した回数分強化する（消費する${res}を直接入力、最大: ${max}個）`,
		hpUpgradeOpt3: (res, used, gain) =>
			`3: 所持している${res}ですべて強化する (${used}個消費 -> 最大HP+${gain})`,
		hpUpgradeOpt0: "0: キャンセル",
		hpUpgradePromptCount: (res, max) =>
			`消費する${res}の数を入力してください (2〜${max}): `,
		hpUpgradeInvalidCount:
			">> 無効な数値です。2個以上、所持数以内の数値を入力してください。",
		hpUpgradeSuccessSingle: (prev, next) =>
			`>> 体力を強化しました！ 最大HP: ${prev} -> ${next}`,
		hpUpgradeSuccessMultiple: (used, res, prev, next) =>
			`>> ${res}を ${used}個 消費して強化しました！ 最大HP: ${prev} -> ${next}`,
		hpUpgradeOddRemainder:
			">> ※ 2個単位での消費のため、端数の1個は温存されました。",
		hpUpgradeSuccessAll: (used, res, prev, next) =>
			`>> ${res}をすべて(${used}個)消費して一括で強化しました！ 最大HP: ${prev} -> ${next}`,
		hpUpgradeCanceled: ">> 強化をキャンセルしました。",

		// Orb Attach
		orbAttachNoOrbs: ">> 装着できるオーブが倉庫にありません！",
		orbAttachTitle: "\n--- 倉庫のオーブ一覧 ---",
		orbAttachPrompt: "装着するオーブの番号を選択: ",
		orbAttachSuccess: (orb) => `>> スロットに [${orb}] を装着しました！`,
		orbAttachFullTitle: "\nスロットが満杯です。上書きする枠を選んでください:",
		orbAttachReplacePrompt: "番号を選択: ",
		orbAttachReplaced: (removed, added) =>
			`>> [${removed}] を破棄し、[${added}] を装着しました！`,

		// Orb Disassemble
		orbDisassembleNotEnough: ">> 分解には倉庫にオーブが2個以上必要です！",
		orbDisassembleTitle: (res) =>
			`\n--- オーブの分解 (任意のオーブ2個 -> ${res}1個) ---`,
		orbDisassembleCount: (count) => `倉庫の未装着オーブ: ${count}個`,
		orbDisassembleOpt1: "1: オーブを個別に選んで分解",
		orbDisassembleOpt2: (res, scrolls) =>
			`2: 倉庫にある未装着オーブをすべて分解する (一括分解: ${res} x${scrolls})`,
		orbDisassembleOpt0: "0: キャンセル",
		orbDisassembleAllSuccess: (usedCount, res, scrolls, total) =>
			`>> 倉庫の未装着オーブ ${usedCount}個 をすべて分解し、${res} x${scrolls} を獲得しました！ (所持: ${res} x${total})`,
		orbDisassembleRemainder: (orb) =>
			`>> ※ 分解できなかった余りのオーブ [${orb}] 1個は倉庫に残りました。`,
		orbDisassembleIndivTitle: "\n--- オーブの個別分解 ---",
		orbDisassemblePick1: "1つ目に分解するオーブの番号を選択: ",
		orbDisassembleSelected1: (orb) => `>> 1つ目: [${orb}] を選択しました。`,
		orbDisassemblePick2: "2つ目に分解するオーブの番号を選択: ",
		orbDisassembleSameError: ">> 1つ目と同じオーブは選択できません！",
		orbDisassembleSuccess: (orb1, orb2, res, total) =>
			`>> [${orb1}] と [${orb2}] を分解し、${res} x1 を獲得しました！ (所持: ${res} x${total})`,
		orbDisassembleCanceled: ">> 分解をキャンセルしました。",

		// Orb Synthesize
		orbSynthNoRecipes:
			">> 合成可能なオーブ（同種の通常オーブ2個以上）が倉庫にありません！",
		orbSynthTitle:
			"\n--- オーブの合成 (同種の通常オーブ2個 -> 上位オーブ1個) ---",
		orbSynthRecipeItem: (idx, from, count, to, toDesc) =>
			`${idx}: [${from}] (所持: ${count}個) -> [${to}] ${toDesc} を合成`,
		orbSynthPrompt: "合成するオーブの番号を選択: ",
		orbSynthSuccess: (from, to) =>
			`>> [${from}] を2個消費し、上位オーブ [${to}] を合成しました！`,

		// Dungeon Exploration
		dungeonEnter: (name) => `\n>>> 【${name}】に突入しました！ <<<`,
		dungeonEncounter: (floor, total, name) =>
			`\n==============================================\n   B${floor}F / B${total}F : ${name} が現れた！\n==============================================`,
		dungeonDefeatPrompt: "\n[!] あなたは力尽きた...",
		dungeonItemsLost: (res, scrolls, orbs) =>
			`[!] 今回獲得したアイテム（${res}: ${scrolls}個, オーブ: ${orbs}個）は失われました。`,
		dungeonCarriedBack:
			"[!] 命からがら拠点へ運ばれました。（所持装備は無事です）",
		dungeonVictory: (enemy) => `\n>> ${enemy} を討伐！`,
		dungeonDropScroll: (res, count) => `   戦利品獲得: ${res} x${count}`,
		dungeonDropOrb: (orb, name) => `   戦利品獲得: オーブ [${orb}] (${name})`,
		dungeonComplete: (name) =>
			`\n**********************************************\n   【${name}】完全踏破！おめでとうございます！   \n**********************************************`,
		dungeonCurrentStatus: (hp, maxHp, res, scrolls, orbs) =>
			`\n現在HP: ${hp}/${maxHp}\n現在の未確定戦利品: ${res} x${scrolls}, オーブ x${orbs}`,
		dungeonNextFloor: "1: 次の階層へ進む",
		dungeonRetreat: "2: 撤退する (戦利品を持ち帰って拠点に戻る)",
		dungeonPromptAction: "行動を選択: ",
		dungeonRetreatSuccess: (res, scrolls, orbs) =>
			`>> 戦利品（${res}: ${scrolls}個, オーブ: ${orbs}個）を持ち帰りました！`,

		// Endless
		endlessEnter: (name) => `\n>>> 【${name}】に突入しました！ <<<`,
		endlessIntro:
			"※ 階層上限のないエンドレスモードです。どこまで潜れるか挑戦しましょう！",
		endlessStartFloorTitle: "\n--- スタート階層選択 ---",
		endlessHighestRecord: (floor) => `最高到達階層: B${floor}F`,
		endlessStartOption: (floor) => `B${floor}F からスタート`,
		endlessCancelOption: "0: 出撃をキャンセル (拠点に戻る)",
		endlessStartFrom: (floor) => `\n>> B${floor}F から探索を開始します！`,
		endlessFloorEncounter: (floor, name) =>
			`\n==============================================\n   B${floor}F : ${name} が現れた！\n==============================================`,
		endlessRecordUpdated: (floor) => `★ 最高到達階層を更新！ (B${floor}F)`,

		// Battle
		battleVs: (php, pmax, ename, ehp, emax, poison) =>
			`\n[YOU] HP: ${php}/${pmax}  vs  [${ename}] HP: ${ehp}/${emax}${
				poison > 0 ? ` (毒:${poison})` : ""
			}`,
		battleCmdAttack: "1: 攻撃する",
		battleCmdRetreat: "2: 撤退する (戦闘から逃げて拠点へ)",
		battleCmdAuto: "3: オート戦闘",
		battleCmdPrompt: "コマンド: ",
		battleRetreatCombat:
			"\n>> 戦闘から離脱し、命からがら帰還しました。（戦利品は持ち帰れません）",
		battleAutoStart: "\n>> [オート戦闘開始] 高速で戦闘を進行します...",
		battleAutoSummaryTitle:
			"\n==============================================\n             【オート戦闘 終了要約】            \n==============================================",
		battleAutoReason: (reason) => `- 終了理由: ${reason}`,
		battleAutoTurns: (turns) => `- 経過ターン数: ${turns} ターン`,
		battleAutoDamageDealt: (dmg) => `- 与えた総ダメージ: ${dmg}`,
		battleAutoDamageTaken: (dmg) => `- 受けた総ダメージ: ${dmg}`,
		battleAutoHealed: (heal) => `- 吸血総回復量: ${heal}`,
		battleAutoRemainingHp: (php, pmax, ename, ehp, emax) =>
			`- 残りHP: あなた ${Math.max(0, php)}/${pmax} | ${ename} ${Math.max(0, ehp)}/${emax}`,
		battleAutoReasonEnemyDefeated: "敵を撃破！",
		battleAutoReasonPoisonDefeated: "毒ダメージにより敵を撃破！",
		battleAutoReasonFainted: "力尽きました...",
		battleAutoReasonDangerHp: (hp, maxHp) =>
			`危険域 (HP ${hp}/${maxHp} <= 30%) に到達したため自動停止`,
		battleInvalidCmd:
			">> 無効なコマンドです。1, 2, 3 のいずれかを入力してください。",
		battlePlayerAttack: (hitIndex, isCrit, ename, dmg) =>
			`>> あなたの${hitIndex}攻撃！ ${isCrit ? "【会心の一撃！】 " : ""}${ename}に ${dmg} ダメージ！`,
		battleCritLabel: "【会心の一撃！】 ",
		battleHitNumber: (hit) => `${hit}撃目の`,
		battleVampHeal: (heal, hp) =>
			`   [吸血] HPが ${heal} 回復した！ (現在HP: ${hp})`,
		battlePoisonInflict: (ename, added, total) =>
			`   [猛毒] ${ename}に毒を付与！ (+${added} / 毒カウント: ${total})`,
		battlePoisonTick: (ename, dmg) =>
			`>> [毒効果] ${ename}は毒で ${dmg} のダメージを受けた！`,
		battleEnemyCounter: (ename, dmg) =>
			`>> ${ename}の反撃！ あなたは ${dmg} のダメージを受けた！`,

		// Settings
		settingsTitle: "\n--- 設定 (Settings) ---",
		settingsCurrent: (lang, theme) =>
			`現在の設定: 言語: [${lang === "ja" ? "日本語 (ja)" : "English (en)"}] | テーマ: [${theme}]`,
		settingsMenu1Lang: "1: 言語切り替え (Toggle Language: en / ja)",
		settingsMenu2Theme: "2: テーマ切り替え (Switch Theme)",
		settingsMenu3ExportTemplate:
			"3: カスタムテーマ雛形を出力 (Export custom_theme.example.json)",
		settingsMenu0Back: "0: 戻る (Back)",
		settingsLangChanged: (lang) =>
			`>> 言語を「${lang === "ja" ? "日本語 (ja)" : "English (en)"}」に切り替えました。`,
		settingsThemeSelectTitle: "\n--- テーマ選択 ---",
		settingsThemeChanged: (theme) =>
			`>> テーマを「${theme}」に切り替えました。`,
		settingsTemplateExported: (path) =>
			`>> カスタムテーマ雛形を出力しました: ${path}`,
		settingsTemplateExportFailed: ">> 雛形の出力に失敗しました:",
	},
	en: {
		title:
			"==============================================\n   Minimal Rogue-lite Prototype (CUI Ver)   \n==============================================",
		startPrompt: "\nPlease choose an option: ",
		startContinue: "1: Continue (Load save.json)",
		startNew: "2: New Game (Start with initial stats)",
		invalidChoice1or2: ">> Invalid choice. Please enter 1 or 2.",
		saveLoaded: ">> Save data loaded successfully!",
		saveStats: (hp, deepest) =>
			`>> Max HP: ${hp} | Deepest Floor: ${deepest > 0 ? `B${deepest}F` : "None"}`,
		untested: "None",
		saveNotFoundOrCorrupted:
			">> save.json was not found or corrupted. Starting a new game.",
		newGameStarted: ">> Starting a new game.",
		savedSuccess: ">> Game saved successfully (save.json).",
		saveFailed: ">> Failed to save game:",
		loadFailed: ">> Failed to load save data:",
		farewell: "Thank you for playing!",

		// Hub
		hubHeader: "【Base Camp】",
		hubStats: (hp, deepest) =>
			`Max HP: ${hp} | Infinite Abyss Record: ${deepest > 0 ? `B${deepest}F` : "None"}`,
		hubEquip: (label, name, plus, atk) =>
			`${label}: ${name}+${plus} (ATK: ${atk})`,
		hubSlots: (count, slots) => `Equipped Slots [${count}/3]: ${slots}`,
		hubEmptySlot: "(Empty)",
		hubStorage: (resource, scrolls, orbs) =>
			`Storage: ${resource} x${scrolls} | Stock Orbs: ${orbs}`,
		hubNone: "(None)",
		hubMenu1Dungeon: "1: Embark to Dungeon",
		hubMenu2Enhance: (verb) => `2: ${verb}`,
		hubMenu3AttachOrb: (label) => `3: Attach Orb to ${label}`,
		hubMenu4Disassemble: (res) =>
			`4: Disassemble Orbs (Any 2 Orbs -> 1 ${res})`,
		hubMenu5Synthesize: "5: Synthesize Orbs (2 Same Orbs -> Plus Orb)",
		hubMenu6EnhanceHp: (res) => `6: Upgrade Max HP (2 ${res} -> +10 Max HP)`,
		hubMenu7Settings: "7: Settings",
		hubMenu0Exit: "0: Quit Game",
		chooseAction: "\nChoose an action: ",
		invalidChoice: "Invalid choice.",

		// Dungeon Select
		dungeonSelectTitle: "\n--- Select Dungeon ---",
		dungeonOptionStarter: (name, floors) =>
			`1: ${name} (${floors} Floors / Beginner)`,
		dungeonOptionDeep: (name, floors, rec) =>
			`2: ${name} (${floors} Floors / ${rec})`,
		dungeonOptionAbyss: (name) => `3: ${name} (Endless / Unlimited Floors)`,
		dungeonCancel: "0: Cancel (Return to Base)",
		dungeonPrompt: "\nSelect a dungeon: ",
		dungeonCanceled: ">> Canceled departure.",
		recommendedDeep: "Recommended +25 or higher",

		// Weapon / Target Upgrade
		enhanceTitle: (label) => `\n--- Enhance ${label} ---`,
		enhanceCurrentRes: (res, count) => `Current ${res}: ${count}`,
		enhanceOpt1: "1: Upgrade once (+1)",
		enhanceOpt2: (res) => `2: Upgrade multiple times (Input amount of ${res})`,
		enhanceOpt3: (res, count) =>
			`3: Upgrade with all available ${res} (+${count})`,
		enhanceOpt0: "0: Cancel",
		enhanceNoResource: (res) => `>> Not enough ${res}!`,
		enhancePromptCount: (res, max) =>
			`Enter amount of ${res} to consume (1 to ${max}): `,
		enhanceInvalidCount:
			">> Invalid amount. Please enter a valid number within your inventory.",
		enhanceSuccessSingle: (verb, label, name, plus, atk) =>
			`>> [${label}] Performed ${verb}! ${name}+${plus} (ATK: ${atk})`,
		enhanceSuccessMultiple: (count, res, verb, name, plus, atk) =>
			`>> Consumed ${count} ${res} for ${verb}! ${name}+${plus} (ATK: ${atk})`,
		enhanceSuccessAll: (count, res, verb, name, plus, atk) =>
			`>> Consumed all (${count}) ${res} for ${verb}! ${name}+${plus} (ATK: ${atk})`,
		enhanceCanceled: ">> Enhancement canceled.",

		// HP Upgrade
		hpUpgradeTitle: "\n--- Upgrade Max HP ---",
		hpUpgradeCurrentHp: (hp) => `Current Max HP: ${hp}`,
		hpUpgradeCurrentRes: (res, count) =>
			`Current ${res}: ${count} (2 ${res} -> +10 Max HP)`,
		hpUpgradeNotEnough: (res, count) =>
			`>> Not enough ${res}! (Required: 2 / Owned: ${count})`,
		hpUpgradeOpt1: (res) => `1: Upgrade once (Consume 2 ${res} -> +10 Max HP)`,
		hpUpgradeOpt2: (res, max) =>
			`2: Upgrade multiple times (Input amount of ${res}, max: ${max})`,
		hpUpgradeOpt3: (res, used, gain) =>
			`3: Upgrade with all available ${res} (Consume ${used} ${res} -> +${gain} Max HP)`,
		hpUpgradeOpt0: "0: Cancel",
		hpUpgradePromptCount: (res, max) =>
			`Enter amount of ${res} to consume (2 to ${max}): `,
		hpUpgradeInvalidCount:
			">> Invalid amount. Please enter a number >= 2 within your inventory.",
		hpUpgradeSuccessSingle: (prev, next) =>
			`>> Max HP upgraded! ${prev} -> ${next}`,
		hpUpgradeSuccessMultiple: (used, res, prev, next) =>
			`>> Consumed ${used} ${res} to upgrade Max HP! ${prev} -> ${next}`,
		hpUpgradeOddRemainder:
			">> * Consumed in units of 2; 1 remaining was retained.",
		hpUpgradeSuccessAll: (used, res, prev, next) =>
			`>> Consumed all (${used}) ${res} to upgrade Max HP! ${prev} -> ${next}`,
		hpUpgradeCanceled: ">> HP upgrade canceled.",

		// Orb Attach
		orbAttachNoOrbs: ">> No orbs available in storage!",
		orbAttachTitle: "\n--- Stored Orbs ---",
		orbAttachPrompt: "Select an orb to attach: ",
		orbAttachSuccess: (orb) => `>> Attached [${orb}] to weapon slot!`,
		orbAttachFullTitle: "\nSlots are full. Choose a slot to overwrite:",
		orbAttachReplacePrompt: "Select slot number: ",
		orbAttachReplaced: (removed, added) =>
			`>> Discarded [${removed}] and equipped [${added}]!`,

		// Orb Disassemble
		orbDisassembleNotEnough:
			">> You need at least 2 orbs in storage to disassemble!",
		orbDisassembleTitle: (res) =>
			`\n--- Disassemble Orbs (Any 2 Orbs -> 1 ${res}) ---`,
		orbDisassembleCount: (count) => `Unassigned Orbs in Storage: ${count}`,
		orbDisassembleOpt1: "1: Select individual orbs to disassemble",
		orbDisassembleOpt2: (res, scrolls) =>
			`2: Disassemble all unassigned orbs (Batch: ${res} x${scrolls})`,
		orbDisassembleOpt0: "0: Cancel",
		orbDisassembleAllSuccess: (usedCount, res, scrolls, total) =>
			`>> Disassembled ${usedCount} orbs and obtained ${res} x${scrolls}! (Total: ${res} x${total})`,
		orbDisassembleRemainder: (orb) =>
			`>> * 1 remaining orb [${orb}] could not be paired and remains in storage.`,
		orbDisassembleIndivTitle: "\n--- Individual Orb Disassembly ---",
		orbDisassemblePick1: "Select first orb to disassemble: ",
		orbDisassembleSelected1: (orb) => `>> First orb selected: [${orb}]`,
		orbDisassemblePick2: "Select second orb to disassemble: ",
		orbDisassembleSameError: ">> Cannot choose the same orb twice!",
		orbDisassembleSuccess: (orb1, orb2, res, total) =>
			`>> Disassembled [${orb1}] and [${orb2}], obtained ${res} x1! (Total: ${res} x${total})`,
		orbDisassembleCanceled: ">> Disassembly canceled.",

		// Orb Synthesize
		orbSynthNoRecipes:
			">> No synthesizable orbs (need at least 2 of the same normal orb) in storage!",
		orbSynthTitle:
			"\n--- Synthesize Orbs (2 Same Normal Orbs -> 1 Plus Orb) ---",
		orbSynthRecipeItem: (idx, from, count, to, toDesc) =>
			`${idx}: [${from}] (Owned: ${count}) -> Synthesize [${to}] ${toDesc}`,
		orbSynthPrompt: "Select recipe number: ",
		orbSynthSuccess: (from, to) =>
			`>> Consumed 2x [${from}] and synthesized [${to}]!`,

		// Dungeon Exploration
		dungeonEnter: (name) => `\n>>> Entered 【${name}】! <<<`,
		dungeonEncounter: (floor, total, name) =>
			`\n==============================================\n   B${floor}F / B${total}F : ${name} appeared!\n==============================================`,
		dungeonDefeatPrompt: "\n[!] You have fallen...",
		dungeonItemsLost: (res, scrolls, orbs) =>
			`[!] Loot gathered on this expedition (${res}: ${scrolls}, Orbs: ${orbs}) was lost.`,
		dungeonCarriedBack:
			"[!] Carried back to base safely. (Equipped items preserved)",
		dungeonVictory: (enemy) => `\n>> Defeated ${enemy}!`,
		dungeonDropScroll: (res, count) => `   Loot obtained: ${res} x${count}`,
		dungeonDropOrb: (orb, name) => `   Loot obtained: Orb [${orb}] (${name})`,
		dungeonComplete: (name) =>
			`\n**********************************************\n   【${name}】 Cleared! Congratulations!   \n**********************************************`,
		dungeonCurrentStatus: (hp, maxHp, res, scrolls, orbs) =>
			`\nCurrent HP: ${hp}/${maxHp}\nPending Loot: ${res} x${scrolls}, Orbs x${orbs}`,
		dungeonNextFloor: "1: Proceed to next floor",
		dungeonRetreat: "2: Retreat (Return to base with collected loot)",
		dungeonPromptAction: "Choose action: ",
		dungeonRetreatSuccess: (res, scrolls, orbs) =>
			`>> Brought back loot (${res}: ${scrolls}, Orbs: ${orbs}) safely!`,

		// Endless
		endlessEnter: (name) => `\n>>> Entered 【${name}】! <<<`,
		endlessIntro:
			"* Endless dungeon with no floor limit. Test how deep you can survive!",
		endlessStartFloorTitle: "\n--- Select Start Floor ---",
		endlessHighestRecord: (floor) => `Highest Record: B${floor}F`,
		endlessStartOption: (floor) => `Start from B${floor}F`,
		endlessCancelOption: "0: Cancel departure (Return to base)",
		endlessStartFrom: (floor) => `\n>> Starting expedition from B${floor}F!`,
		endlessFloorEncounter: (floor, name) =>
			`\n==============================================\n   B${floor}F : ${name} appeared!\n==============================================`,
		endlessRecordUpdated: (floor) =>
			`★ New highest record reached! (B${floor}F)`,

		// Battle
		battleVs: (php, pmax, ename, ehp, emax, poison) =>
			`\n[YOU] HP: ${php}/${pmax}  vs  [${ename}] HP: ${ehp}/${emax}${
				poison > 0 ? ` (Poison:${poison})` : ""
			}`,
		battleCmdAttack: "1: Attack",
		battleCmdRetreat: "2: Retreat (Flee to base)",
		battleCmdAuto: "3: Auto Combat",
		battleCmdPrompt: "Command: ",
		battleRetreatCombat:
			"\n>> Fled from battle and escaped back to base. (Loot lost)",
		battleAutoStart:
			"\n>> [Auto Combat Started] Resolving combat at high speed...",
		battleAutoSummaryTitle:
			"\n==============================================\n          【Auto Combat Summary Report】         \n==============================================",
		battleAutoReason: (reason) => `- End Reason: ${reason}`,
		battleAutoTurns: (turns) => `- Elapsed Turns: ${turns} turns`,
		battleAutoDamageDealt: (dmg) => `- Total Damage Dealt: ${dmg}`,
		battleAutoDamageTaken: (dmg) => `- Total Damage Taken: ${dmg}`,
		battleAutoHealed: (heal) => `- Total HP Restored via Vampirism: ${heal}`,
		battleAutoRemainingHp: (php, pmax, ename, ehp, emax) =>
			`- Remaining HP: You ${Math.max(0, php)}/${pmax} | ${ename} ${Math.max(0, ehp)}/${emax}`,
		battleAutoReasonEnemyDefeated: "Enemy Defeated!",
		battleAutoReasonPoisonDefeated: "Enemy slain by poison damage!",
		battleAutoReasonFainted: "You were knocked out...",
		battleAutoReasonDangerHp: (hp, maxHp) =>
			`Auto stopped upon reaching danger zone (HP ${hp}/${maxHp} <= 30%)`,
		battleInvalidCmd: ">> Invalid command. Please enter 1, 2, or 3.",
		battlePlayerAttack: (hitIndex, isCrit, ename, dmg) =>
			`>> Your ${hitIndex}attack! ${isCrit ? "【Critical Strike!】 " : ""}${ename} took ${dmg} damage!`,
		battleCritLabel: "【Critical Strike!】 ",
		battleHitNumber: (hit) => `strike #${hit} `,
		battleVampHeal: (heal, hp) =>
			`   [Vampiric] Restored ${heal} HP! (Current HP: ${hp})`,
		battlePoisonInflict: (ename, added, total) =>
			`   [Poison] Inflicted poison to ${ename}! (+${added} / Poison count: ${total})`,
		battlePoisonTick: (ename, dmg) =>
			`>> [Poison] ${ename} took ${dmg} poison damage!`,
		battleEnemyCounter: (ename, dmg) =>
			`>> ${ename}'s counterattack! You took ${dmg} damage!`,

		// Settings
		settingsTitle: "\n--- Settings ---",
		settingsCurrent: (lang, theme) =>
			`Current Settings: Language: [${lang === "ja" ? "日本語 (ja)" : "English (en)"}] | Theme: [${theme}]`,
		settingsMenu1Lang: "1: Toggle Language (en / ja)",
		settingsMenu2Theme: "2: Switch Theme",
		settingsMenu3ExportTemplate:
			"3: Export Custom Theme Template (custom_theme.example.json)",
		settingsMenu0Back: "0: Back",
		settingsLangChanged: (lang) =>
			`>> Language switched to "${lang === "ja" ? "日本語 (ja)" : "English (en)"}".`,
		settingsThemeSelectTitle: "\n--- Select Theme ---",
		settingsThemeChanged: (theme) => `>> Switched theme to "${theme}".`,
		settingsTemplateExported: (path) =>
			`>> Exported custom theme template: ${path}`,
		settingsTemplateExportFailed: ">> Failed to export template:",
	},
};

export function detectDefaultLanguage(): Language {
	const env = process.env as {
		LANG?: string | undefined;
		LC_ALL?: string | undefined;
		LC_MESSAGES?: string | undefined;
	};
	const envLang = (
		env.LANG ||
		env.LC_ALL ||
		env.LC_MESSAGES ||
		""
	).toLowerCase();
	if (envLang.startsWith("ja") || envLang.includes("ja_jp")) {
		return "ja";
	}
	try {
		const locale = Intl.DateTimeFormat().resolvedOptions().locale.toLowerCase();
		if (locale.startsWith("ja")) {
			return "ja";
		}
	} catch {
		// ignore
	}
	return "en";
}

// ============================================================================
// Themes & Built-in Presets
// ============================================================================

export const PRESET_CLASSIC_FANTASY: ThemeDefinition = {
	id: "classic_fantasy",
	ja: {
		name: "王道ファンタジー (Classic Fantasy)",
		targetNameLabel: "所持装備",
		defaultTargetName: "どうのつるぎ",
		enhanceVerb: "鍛冶屋で鍛える",
		resourceName: "強化の書",
		orbs: {
			MULTI_HIT: { name: "連撃の印", desc: "2回攻撃/威力65%" },
			CRITICAL: { name: "会心の印", desc: "25%で2倍" },
			VAMP: { name: "吸血の印", desc: "与ダメの15%回復" },
			POISON: { name: "猛毒の印", desc: "攻撃時毒+1/ターン末毒x3ダメ" },
			MULTI_HIT_PLUS: { name: "連撃の印+", desc: "2回攻撃/威力75%" },
			CRITICAL_PLUS: { name: "会心の印+", desc: "40%で2倍" },
			VAMP_PLUS: { name: "吸血の印+", desc: "与ダメの25%回復" },
			POISON_PLUS: { name: "猛毒の印+", desc: "攻撃時毒+2/ターン末毒x3ダメ" },
		},
		dungeons: {
			starter: "始まりの洞窟",
			deep: "灼熱の深層",
			abyss: "無限の深淵",
		},
		enemyPrefixes: [
			"凶暴な",
			"深淵の",
			"古代の",
			"紅蓮の",
			"漆黒の",
			"彷徨える",
			"狂気の",
			"奈落の",
		],
		enemyBases: [
			"スライム",
			"コボルト",
			"オーク",
			"ゴーレム",
			"ドラゴン",
			"ワイバーン",
			"デーモン",
			"死霊騎士",
		],
		retreatMessage: "慎重に撤退を選択し、拠点へ帰還した。",
	},
	en: {
		name: "Classic Fantasy",
		targetNameLabel: "Equipment",
		defaultTargetName: "Bronze Sword",
		enhanceVerb: "Forge at Blacksmith",
		resourceName: "Upgrade Scroll",
		orbs: {
			MULTI_HIT: { name: "Twin Strike", desc: "2 hits at 65% power each" },
			CRITICAL: { name: "Critical Strike", desc: "25% chance for 2x damage" },
			VAMP: { name: "Vampiric Drain", desc: "Heal for 15% of damage dealt" },
			POISON: {
				name: "Deadly Poison",
				desc: "+1 poison on hit, 3x dmg per turn",
			},
			MULTI_HIT_PLUS: {
				name: "Twin Strike+",
				desc: "2 hits at 75% power each",
			},
			CRITICAL_PLUS: {
				name: "Critical Strike+",
				desc: "40% chance for 2x damage",
			},
			VAMP_PLUS: {
				name: "Vampiric Drain+",
				desc: "Heal for 25% of damage dealt",
			},
			POISON_PLUS: {
				name: "Deadly Poison+",
				desc: "+2 poison on hit, 3x dmg per turn",
			},
		},
		dungeons: {
			starter: "Cave of Beginnings",
			deep: "Scorching Depths",
			abyss: "Infinite Abyss",
		},
		enemyPrefixes: [
			"Fierce",
			"Abyssal",
			"Ancient",
			"Crimson",
			"Dark",
			"Wandering",
			"Mad",
			"Infernal",
		],
		enemyBases: [
			"Slime",
			"Kobold",
			"Orc",
			"Golem",
			"Dragon",
			"Wyvern",
			"Demon",
			"Death Knight",
		],
		retreatMessage: "Carefully chose to retreat and returned to base.",
	},
};

export const PRESET_CYBERPUNK: ThemeDefinition = {
	id: "cyberpunk",
	ja: {
		name: "サイバーパンク (Cyberpunk)",
		targetNameLabel: "サイバー兵装",
		defaultTargetName: "パルスブレード",
		enhanceVerb: "システムオーバークロック",
		resourceName: "ナノチップ",
		orbs: {
			MULTI_HIT: { name: "多段バースト", desc: "2連撃/出力65%×2" },
			CRITICAL: { name: "クリティカル注入", desc: "25%で2倍電圧" },
			VAMP: {
				name: "ナノ修復ドレイン",
				desc: "与ダメの15%で機体修復",
			},
			POISON: {
				name: "神経汚染",
				desc: "攻撃時毒+1/ターン末毒x3ダメ",
			},
			MULTI_HIT_PLUS: {
				name: "多段バースト+",
				desc: "2連撃/出力75%×2",
			},
			CRITICAL_PLUS: {
				name: "クリティカル注入+",
				desc: "40%で2倍電圧",
			},
			VAMP_PLUS: {
				name: "ナノ修復ドレイン+",
				desc: "与ダメの25%で機体修復",
			},
			POISON_PLUS: {
				name: "神経汚染+",
				desc: "攻撃時毒+2/ターン末毒x3ダメ",
			},
		},
		dungeons: {
			starter: "閉鎖ネットワーク",
			deep: "企業中枢サーバー",
			abyss: "無限の電脳網",
		},
		enemyPrefixes: [
			"グリッチの",
			"暴走した",
			"オーバークロックの",
			"侵食された",
			"強化型の",
			"ステルス",
			"試作型",
			"自律型の",
		],
		enemyBases: [
			"セキュリティドローン",
			"偵察ボット",
			"戦闘アンドロイド",
			"重装歩行メカ",
			"中枢メインフレーム",
			"電脳ドラゴン",
			"ウイルスロード",
			"暴走AI",
		],
		retreatMessage:
			"機体の致命的なシャットダウンを防ぐため、戦闘から離脱した。",
	},
	en: {
		name: "Cyberpunk Protocol",
		targetNameLabel: "Cyber Weapon",
		defaultTargetName: "Pulse Blade",
		enhanceVerb: "Overclock System",
		resourceName: "Nanite Chip",
		orbs: {
			MULTI_HIT: { name: "Multi-Burst", desc: "2 strikes at 65% output each" },
			CRITICAL: {
				name: "Critical Injection",
				desc: "25% chance for 2x voltage",
			},
			VAMP: {
				name: "Nano-Drain",
				desc: "Siphon 15% damage as chassis repair",
			},
			POISON: {
				name: "Neuro-Toxin",
				desc: "Inject +1 corrosive payload on hit",
			},
			MULTI_HIT_PLUS: {
				name: "Multi-Burst+",
				desc: "2 strikes at 75% output each",
			},
			CRITICAL_PLUS: {
				name: "Critical Injection+",
				desc: "40% chance for 2x voltage",
			},
			VAMP_PLUS: {
				name: "Nano-Drain+",
				desc: "Siphon 25% damage as chassis repair",
			},
			POISON_PLUS: {
				name: "Neuro-Toxin+",
				desc: "Inject +2 corrosive payload on hit",
			},
		},
		dungeons: {
			starter: "Subnet Alpha",
			deep: "Corporate Core",
			abyss: "Infinite Cyberspace",
		},
		enemyPrefixes: [
			"Glitch",
			"Rogue",
			"Overclocked",
			"Corrupted",
			"Augmented",
			"Stealth",
			"Prototype",
			"Autonomous",
		],
		enemyBases: [
			"Security Drone",
			"Recon Bot",
			"Combat Android",
			"Heavy Mech",
			"Mainframe Core",
			"Cyber Dragon",
			"Virus Lord",
			"Rogue AI",
		],
		retreatMessage:
			"Disengaged from combat to prevent terminal hardware shutdown.",
	},
};

export const PRESET_PARTNER_SYNC: ThemeDefinition = {
	id: "partner_sync",
	ja: {
		name: "パートナー・シンクロ (相棒育成)",
		targetNameLabel: "相棒",
		defaultTargetName: "戦術アンドロイド「アイリス」",
		enhanceVerb: "同期率を向上させる",
		resourceName: "メモリコア",
		orbs: {
			MULTI_HIT: {
				name: "連携戦術",
				desc: "息の合った2連撃を行う (威力65%×2)",
			},
			CRITICAL: {
				name: "弱点看破",
				desc: "隙を突いて致命打を与える (25%で2倍)",
			},
			VAMP: {
				name: "自己修復",
				desc: "戦闘データから装甲をナノ修復 (与ダメの15%回復)",
			},
			POISON: {
				name: "浸食ノイズ",
				desc: "敵システムに持続ダメージ (毒+1/ターン末毒x3)",
			},
			MULTI_HIT_PLUS: {
				name: "連携戦術・極",
				desc: "息の合った2連撃を行う (威力75%×2)",
			},
			CRITICAL_PLUS: {
				name: "弱点看破・極",
				desc: "隙を突いて致命打を与える (40%で2倍)",
			},
			VAMP_PLUS: {
				name: "自己修復・極",
				desc: "戦闘データから装甲をナノ修復 (与ダメの25%回復)",
			},
			POISON_PLUS: {
				name: "浸食ノイズ・極",
				desc: "敵システムに持続ダメージ (毒+2/ターン末毒x3)",
			},
		},
		dungeons: {
			starter: "廃墟区域",
			deep: "汚染中枢",
			abyss: "未知の最深部",
		},
		enemyPrefixes: [
			"暴走した",
			"変異型",
			"侵略型",
			"警戒態勢の",
			"重装甲",
			"高機動",
			"汚染された",
			"古代遺産の",
		],
		enemyBases: [
			"暴走ドローン",
			"変異体",
			"侵略尖兵",
			"防衛要塞",
			"掃討機兵",
			"自律兵器",
			"殲滅ユニット",
			"支配者コア",
		],
		retreatMessage: "アイリスの負荷を考慮し、一時帰還した。",
	},
	en: {
		name: "Partner Sync",
		targetNameLabel: "Partner",
		defaultTargetName: 'Tactical Android "Iris"',
		enhanceVerb: "Deepen Sync",
		resourceName: "Memory Core",
		orbs: {
			MULTI_HIT: {
				name: "Tandem Tactics",
				desc: "2 coordinated strikes at 65% power each",
			},
			CRITICAL: {
				name: "Exploit Weakness",
				desc: "25% chance for 2x fatal damage",
			},
			VAMP: {
				name: "Self-Repair",
				desc: "Restore 15% damage dealt as nanite repairs",
			},
			POISON: {
				name: "Corrosive Noise",
				desc: "+1 corrosive noise on hit, 3x dmg per turn",
			},
			MULTI_HIT_PLUS: {
				name: "Tandem Tactics+",
				desc: "2 coordinated strikes at 75% power each",
			},
			CRITICAL_PLUS: {
				name: "Exploit Weakness+",
				desc: "40% chance for 2x fatal damage",
			},
			VAMP_PLUS: {
				name: "Self-Repair+",
				desc: "Restore 25% damage dealt as nanite repairs",
			},
			POISON_PLUS: {
				name: "Corrosive Noise+",
				desc: "+2 corrosive noise on hit, 3x dmg per turn",
			},
		},
		dungeons: {
			starter: "Ruined Sector",
			deep: "Contaminated Core",
			abyss: "Unknown Depths",
		},
		enemyPrefixes: [
			"Rampaging",
			"Mutant",
			"Invading",
			"Alert",
			"Armored",
			"High-Mobility",
			"Contaminated",
			"Relic",
		],
		enemyBases: [
			"Rogue Drone",
			"Mutant Beast",
			"Vanguard Scout",
			"Defense Bastion",
			"Sweeper Mech",
			"Autonomous Weapon",
			"Annihilator Unit",
			"Overlord Core",
		],
		retreatMessage: "Retreated to base to reduce load on Iris.",
	},
};

export const PRESETS: Record<string, ThemeDefinition> = {
	classic_fantasy: PRESET_CLASSIC_FANTASY,
	cyberpunk: PRESET_CYBERPUNK,
	partner_sync: PRESET_PARTNER_SYNC,
};

export const PRESET_CLASSIC_FANTASY_JA: Theme = resolveTheme(
	PRESET_CLASSIC_FANTASY,
	"ja",
);
export const PRESET_CLASSIC_FANTASY_EN: Theme = resolveTheme(
	PRESET_CLASSIC_FANTASY,
	"en",
);

export function getPresetTheme(themeId: string, lang: Language): Theme {
	const preset = PRESETS[themeId] ?? PRESET_CLASSIC_FANTASY;
	return resolveTheme(preset, lang);
}

export function validateThemeLocaleData(data: unknown): ThemeLocaleData | null {
	if (!data || typeof data !== "object") return null;
	const obj = data as Partial<ThemeLocaleData>;
	if (
		typeof obj.name !== "string" ||
		typeof obj.targetNameLabel !== "string" ||
		typeof obj.defaultTargetName !== "string" ||
		typeof obj.enhanceVerb !== "string" ||
		typeof obj.resourceName !== "string" ||
		!obj.orbs ||
		typeof obj.orbs !== "object" ||
		!obj.dungeons ||
		typeof obj.dungeons !== "object" ||
		!Array.isArray(obj.enemyPrefixes) ||
		!Array.isArray(obj.enemyBases) ||
		typeof obj.retreatMessage !== "string"
	) {
		return null;
	}

	const requiredOrbs: OrbType[] = [
		"MULTI_HIT",
		"CRITICAL",
		"VAMP",
		"POISON",
		"MULTI_HIT_PLUS",
		"CRITICAL_PLUS",
		"VAMP_PLUS",
		"POISON_PLUS",
	];
	for (const orb of requiredOrbs) {
		const orbDef = (obj.orbs as Record<string, unknown>)[orb];
		if (
			!orbDef ||
			typeof orbDef !== "object" ||
			typeof (orbDef as { name?: unknown }).name !== "string" ||
			typeof (orbDef as { desc?: unknown }).desc !== "string"
		) {
			return null;
		}
	}

	const dungeons = obj.dungeons as {
		starter?: unknown;
		deep?: unknown;
		abyss?: unknown;
	};
	if (
		typeof dungeons.starter !== "string" ||
		typeof dungeons.deep !== "string" ||
		typeof dungeons.abyss !== "string"
	) {
		return null;
	}

	return obj as ThemeLocaleData;
}

interface RawThemeDef {
	id?: unknown;
	ja?: unknown;
	en?: unknown;
}

export function validateThemeDefinition(data: unknown): ThemeDefinition | null {
	if (!data || typeof data !== "object") return null;
	const obj = data as RawThemeDef;
	if (typeof obj.id !== "string") return null;

	const jaLocale = validateThemeLocaleData(obj.ja);
	const enLocale = validateThemeLocaleData(obj.en);
	if (jaLocale && enLocale) {
		return {
			id: obj.id,
			ja: jaLocale,
			en: enLocale,
		};
	}

	// Fallback for legacy flat format
	const legacyLocale = validateThemeLocaleData(obj);
	if (legacyLocale) {
		return {
			id: obj.id,
			ja: jaLocale ?? legacyLocale,
			en: enLocale ?? legacyLocale,
		};
	}

	return null;
}

export function validateTheme(
	data: unknown,
	lang: Language = "ja",
): Theme | null {
	const def = validateThemeDefinition(data);
	if (!def) return null;
	return resolveTheme(def, lang);
}

export function loadCustomThemeDefIfExists(
	filePath = CUSTOM_THEME_FILE_PATH,
): ThemeDefinition | null {
	try {
		if (!existsSync(filePath)) {
			return null;
		}
		const raw = readFileSync(filePath, "utf-8");
		const parsed = JSON.parse(raw);
		const valid = validateThemeDefinition(parsed);
		if (valid) {
			return valid;
		}
		console.warn(
			`>> ${filePath} exists but is missing required theme properties.`,
		);
		return null;
	} catch (error) {
		console.warn(`>> Failed to read ${filePath}:`, error);
		return null;
	}
}

export function loadCustomThemeIfExists(
	filePath = CUSTOM_THEME_FILE_PATH,
	lang: Language = "ja",
): Theme | null {
	const def = loadCustomThemeDefIfExists(filePath);
	if (!def) return null;
	return resolveTheme(def, lang);
}

export function exportCustomThemeTemplate(
	filePath = CUSTOM_THEME_EXAMPLE_PATH,
): boolean {
	try {
		const template: ThemeDefinition = {
			id: "custom_scifi",
			ja: {
				name: "SF星間探査 (Sci-Fi Star Explorer)",
				targetNameLabel: "旗艦武装",
				defaultTargetName: "フォトンランス",
				enhanceVerb: "リアクターを調整する",
				resourceName: "プラズマコア",
				orbs: {
					MULTI_HIT: {
						name: "ツインビーム",
						desc: "2連照射/威力65%",
					},
					CRITICAL: {
						name: "クリティカルパルス",
						desc: "25%で2倍",
					},
					VAMP: {
						name: "シールドドレイン",
						desc: "与ダメの15%回復",
					},
					POISON: {
						name: "アシッド腐食",
						desc: "攻撃時毒+1/ターン末毒x3ダメ",
					},
					MULTI_HIT_PLUS: {
						name: "ツインビーム+",
						desc: "2連照射/威力75%",
					},
					CRITICAL_PLUS: {
						name: "クリティカルパルス+",
						desc: "40%で2倍",
					},
					VAMP_PLUS: {
						name: "シールドドレイン+",
						desc: "与ダメの25%回復",
					},
					POISON_PLUS: {
						name: "アシッド腐食+",
						desc: "攻撃時毒+2/ターン末毒x3ダメ",
					},
				},
				dungeons: {
					starter: "小惑星前哨基地",
					deep: "遺棄された弩級艦",
					abyss: "事象の地平線",
				},
				enemyPrefixes: [
					"敵対的な",
					"過負荷の",
					"変異した",
					"未知生命の",
					"サイバネの",
					"凶暴な",
					"古代の",
					"コズミック",
				],
				enemyBases: [
					"偵察ドローン",
					"バイオスウォーム",
					"宇宙海賊",
					"戦闘メカ",
					"虚空の怪異",
					"星喰らい",
					"ナノ集合体",
					"タイタン弩級艦",
				],
				retreatMessage:
					"ワープドライブを緊急起動し、軌道ステーションへ退避した。",
			},
			en: {
				name: "Sci-Fi Star Explorer",
				targetNameLabel: "Flagship Weapon",
				defaultTargetName: "Photon Lance",
				enhanceVerb: "Calibrate Reactor",
				resourceName: "Plasma Core",
				orbs: {
					MULTI_HIT: {
						name: "Twin Beam",
						desc: "Fire dual beams at 65% energy each",
					},
					CRITICAL: {
						name: "Critical Pulse",
						desc: "25% chance for 2x focused damage",
					},
					VAMP: {
						name: "Shield Siphon",
						desc: "Absorb 15% damage dealt as shield energy",
					},
					POISON: {
						name: "Corrosive Acid",
						desc: "+1 chemical payload on hit, deals 3x dmg per tick",
					},
					MULTI_HIT_PLUS: {
						name: "Twin Beam+",
						desc: "Fire dual beams at 75% energy each",
					},
					CRITICAL_PLUS: {
						name: "Critical Pulse+",
						desc: "40% chance for 2x focused damage",
					},
					VAMP_PLUS: {
						name: "Shield Siphon+",
						desc: "Absorb 25% damage dealt as shield energy",
					},
					POISON_PLUS: {
						name: "Corrosive Acid+",
						desc: "+2 chemical payload on hit, deals 3x dmg per tick",
					},
				},
				dungeons: {
					starter: "Asteroid Outpost",
					deep: "Derelict Dreadnought",
					abyss: "Event Horizon Void",
				},
				enemyPrefixes: [
					"Hostile",
					"Overcharged",
					"Mutated",
					"Alien",
					"Cybernetic",
					"Vicious",
					"Ancient",
					"Cosmic",
				],
				enemyBases: [
					"Scout Drone",
					"Bio-Swarm",
					"Pirate Raider",
					"War Mech",
					"Void Beast",
					"Star Devourer",
					"Nanite Hive",
					"Titan Dreadnought",
				],
				retreatMessage:
					"Warp drive engaged; executed emergency tactical jump back to orbital station.",
			},
		};
		writeFileSync(filePath, JSON.stringify(template, null, 2), "utf-8");
		return true;
	} catch (error) {
		console.error(">> Failed to write theme template:", error);
		return false;
	}
}

// ============================================================================
// Combat Constants & Drop Tables
// ============================================================================

const SYNTHESIS_RECIPES: Record<NormalOrbType, PlusOrbType> = {
	MULTI_HIT: "MULTI_HIT_PLUS",
	CRITICAL: "CRITICAL_PLUS",
	VAMP: "VAMP_PLUS",
	POISON: "POISON_PLUS",
};

const DROP_ORBS: readonly NormalOrbType[] = [
	"MULTI_HIT",
	"CRITICAL",
	"VAMP",
	"POISON",
];

const PLUS_ORBS: readonly PlusOrbType[] = [
	"MULTI_HIT_PLUS",
	"CRITICAL_PLUS",
	"VAMP_PLUS",
	"POISON_PLUS",
];

function getRandomOrb(): NormalOrbType {
	const orb = DROP_ORBS[Math.floor(Math.random() * DROP_ORBS.length)];
	return orb ?? "MULTI_HIT";
}

function getRandomPlusOrb(): PlusOrbType {
	const orb = PLUS_ORBS[Math.floor(Math.random() * PLUS_ORBS.length)];
	return orb ?? "MULTI_HIT_PLUS";
}

export function createStarterDungeon(theme: Theme): DungeonDef {
	const bases = theme.enemyBases;
	const b1 = bases[0] ?? "Slime";
	const b2 = bases[1] ?? bases[0] ?? "Kobold";
	const b3 = bases[2] ?? bases[0] ?? "Orc";
	const b4 = bases[3] ?? bases[0] ?? "Golem";
	const b5 = bases[4] ?? bases[0] ?? "Dragon";

	const enemies: EnemyTemplate[] = [
		{
			name: b1,
			hp: 15,
			atk: 3,
			getDrops: () => ({ scrolls: 1, orbs: [] }),
		},
		{
			name: b2,
			hp: 28,
			atk: 6,
			getDrops: () => ({ scrolls: 1, orbs: [] }),
		},
		{
			name: b3,
			hp: 45,
			atk: 10,
			getDrops: () => ({ scrolls: 0, orbs: [getRandomOrb()] }),
		},
		{
			name: b4,
			hp: 75,
			atk: 14,
			getDrops: () => ({ scrolls: 2, orbs: [] }),
		},
		{
			name: `${b5} (BOSS)`,
			hp: 130,
			atk: 20,
			getDrops: () => ({ scrolls: 3, orbs: [getRandomOrb()] }),
		},
	];

	return {
		name: theme.dungeons.starter,
		floors: 5,
		getEnemy: (floor: number): EnemyTemplate => {
			const template = enemies[floor - 1];
			if (!template) throw new Error(`Invalid floor: ${floor}`);
			return template;
		},
	};
}

export function createDeepDungeon(
	theme: Theme,
	recommended: string,
): DungeonDef {
	const prefixes = theme.enemyPrefixes;
	const bases = theme.enemyBases;

	const getRandomItem = <T>(items: readonly T[]): T => {
		const item = items[Math.floor(Math.random() * items.length)];
		if (!item) throw new Error("Empty array");
		return item;
	};

	const getPrefix = (idx: number) => prefixes[idx % prefixes.length] ?? "";
	const getBase = (idx: number) => bases[idx % bases.length] ?? "Enemy";

	return {
		name: theme.dungeons.deep,
		floors: 10,
		recommended,
		getEnemy: (floor: number): EnemyTemplate => {
			const getDeepFloorScrolls = () => Math.floor(Math.random() * 4) + 2;

			if (floor >= 1 && floor <= 3) {
				const candidates = [
					{
						name: `${getPrefix(0)} ${getBase(2)}`.trim(),
						hp: 80,
						atk: 12,
					},
					{
						name: `${getPrefix(1)} ${getBase(4)}`.trim(),
						hp: 120,
						atk: 16,
					},
				];
				const picked = getRandomItem(candidates);
				return {
					name: picked.name,
					hp: picked.hp,
					atk: picked.atk,
					getDrops: () => ({ scrolls: getDeepFloorScrolls(), orbs: [] }),
				};
			}

			if (floor >= 4 && floor <= 6) {
				const candidates = [
					{
						name: `${getPrefix(2)} ${getBase(3)}`.trim(),
						hp: 180,
						atk: 20,
					},
					{
						name: `${getPrefix(3)} ${getBase(7)}`.trim(),
						hp: 220,
						atk: 24,
					},
				];
				const picked = getRandomItem(candidates);
				return {
					name: picked.name,
					hp: picked.hp,
					atk: picked.atk,
					getDrops: () => ({ scrolls: getDeepFloorScrolls(), orbs: [] }),
				};
			}

			if (floor >= 7 && floor <= 9) {
				const candidates = [
					{
						name: `${getPrefix(4)} ${getBase(5)}`.trim(),
						hp: 280,
						atk: 28,
					},
					{
						name: `${getPrefix(5)} ${getBase(6)}`.trim(),
						hp: 350,
						atk: 32,
					},
				];
				const picked = getRandomItem(candidates);
				return {
					name: picked.name,
					hp: picked.hp,
					atk: picked.atk,
					getDrops: () => ({ scrolls: getDeepFloorScrolls(), orbs: [] }),
				};
			}

			if (floor === 10) {
				return {
					name: `${getPrefix(6)} ${getBase(7)} (BOSS)`.trim(),
					hp: 500,
					atk: 36,
					getDrops: () => ({
						scrolls: 10,
						orbs: [getRandomPlusOrb()],
					}),
				};
			}

			throw new Error(`Invalid floor: ${floor}`);
		},
	};
}

export function generateEndlessEnemy(
	floor: number,
	theme: Theme,
): EnemyTemplate {
	const prefixes = theme.enemyPrefixes;
	const enemyNames = theme.enemyBases;

	const prefix =
		prefixes[Math.floor(Math.random() * prefixes.length)] ?? "Dark";
	const baseName =
		enemyNames[Math.floor(Math.random() * enemyNames.length)] ?? "Beast";
	const isBoss = floor % 5 === 0;
	const bossTag = theme.language === "ja" ? "(中ボス)" : "(Mini-Boss)";
	const separator = theme.language === "ja" ? "" : " ";
	const name = isBoss
		? `${prefix}${separator}${baseName} ${bossTag}`
		: `${prefix}${separator}${baseName}`;

	const hp = Math.floor(60 + floor ** 1.4 * 12);
	const atk = Math.floor(8 + floor * 2.5);

	const scrolls = Math.min(10, Math.floor(1 + floor / 2));
	const orbs: OrbType[] = isBoss ? [getRandomPlusOrb()] : [];

	return {
		name,
		hp,
		atk,
		getDrops: () => ({ scrolls, orbs }),
	};
}

// ============================================================================
// Main Game Engine
// ============================================================================

export class Game {
	private readonly rl: readline.Interface;
	private language: Language;
	private theme: Theme;
	private customThemeDef: ThemeDefinition | null = null;
	private readonly player: Player;
	private stockScrolls = 0;
	private stockOrbs: OrbType[] = [];
	private deepestFloor = 0;
	private isRunning = true;

	constructor() {
		this.rl = readline.createInterface({ input, output });
		this.language = detectDefaultLanguage();
		this.customThemeDef = loadCustomThemeDefIfExists();

		if (this.customThemeDef) {
			this.theme = resolveTheme(this.customThemeDef, this.language);
		} else {
			this.theme = getPresetTheme("classic_fantasy", this.language);
		}

		this.player = {
			maxHp: 50,
			hp: 50,
			weapon: {
				name: this.theme.defaultTargetName,
				baseAtk: 10,
				plus: 0,
				slots: [],
			},
		};
	}

	private get msg(): I18nMessages {
		return MESSAGES[this.language];
	}

	private getOrbDesc(orb: OrbType): string {
		const orbDef = this.theme.orbs[orb];
		if (!orbDef) return orb;
		return `${orbDef.name} (${orbDef.desc})`;
	}

	private getWeaponAtk(): number {
		return this.player.weapon.baseAtk + this.player.weapon.plus * 2;
	}

	private saveGame(): void {
		try {
			const data: SaveData = {
				weapon: this.player.weapon,
				stockScrolls: this.stockScrolls,
				stockOrbs: this.stockOrbs,
				maxHp: this.player.maxHp,
				deepestFloor: this.deepestFloor,
				language: this.language,
				themeId: this.theme.id,
			};
			writeFileSync(SAVE_FILE_PATH, JSON.stringify(data, null, 2), "utf-8");
			console.log(this.msg.savedSuccess);
		} catch (error) {
			console.error(this.msg.saveFailed, error);
		}
	}

	private loadGame(): boolean {
		try {
			if (!existsSync(SAVE_FILE_PATH)) {
				return false;
			}
			const raw = readFileSync(SAVE_FILE_PATH, "utf-8");
			const data = JSON.parse(raw) as SaveData;
			if (
				data?.weapon &&
				typeof data.stockScrolls === "number" &&
				Array.isArray(data.stockOrbs)
			) {
				if (data.language === "en" || data.language === "ja") {
					this.language = data.language;
				}

				if (this.customThemeDef) {
					this.theme = resolveTheme(this.customThemeDef, this.language);
				} else if (data.themeId) {
					this.theme = getPresetTheme(data.themeId, this.language);
				} else {
					this.theme = getPresetTheme("classic_fantasy", this.language);
				}

				this.player.weapon = {
					name: data.weapon.name,
					baseAtk: data.weapon.baseAtk,
					plus: data.weapon.plus,
					slots: Array.isArray(data.weapon.slots) ? [...data.weapon.slots] : [],
				};
				this.stockScrolls = data.stockScrolls;
				this.stockOrbs = [...data.stockOrbs];
				this.player.maxHp =
					typeof data.maxHp === "number" &&
					Number.isFinite(data.maxHp) &&
					data.maxHp > 0
						? data.maxHp
						: 50;
				this.player.hp = this.player.maxHp;
				this.deepestFloor =
					typeof data.deepestFloor === "number" &&
					Number.isFinite(data.deepestFloor) &&
					data.deepestFloor >= 0
						? data.deepestFloor
						: 0;
				return true;
			}
			return false;
		} catch (error) {
			console.error(this.msg.loadFailed, error);
			return false;
		}
	}

	async start(): Promise<void> {
		console.log(this.msg.title);
		console.log(`\n${this.msg.startContinue}`);
		console.log(this.msg.startNew);

		while (true) {
			const choice = await this.rl.question(this.msg.startPrompt);
			const trimmed = choice.trim();
			if (trimmed === "1") {
				if (this.loadGame()) {
					console.log(this.msg.saveLoaded);
					console.log(this.msg.saveStats(this.player.maxHp, this.deepestFloor));
				} else {
					console.log(this.msg.saveNotFoundOrCorrupted);
				}
				break;
			}
			if (trimmed === "2") {
				console.log(this.msg.newGameStarted);
				break;
			}
			console.log(this.msg.invalidChoice1or2);
		}

		while (this.isRunning) {
			await this.hubPhase();
		}
	}

	private async hubPhase(): Promise<void> {
		this.player.hp = this.player.maxHp;
		console.log("\n----------------------------------------------");
		console.log(this.msg.hubHeader);
		console.log(this.msg.hubStats(this.player.maxHp, this.deepestFloor));
		console.log(
			this.msg.hubEquip(
				this.theme.targetNameLabel,
				this.player.weapon.name,
				this.player.weapon.plus,
				this.getWeaponAtk(),
			),
		);
		console.log(
			this.msg.hubSlots(
				this.player.weapon.slots.length,
				this.player.weapon.slots.map((s) => `[${s}]`).join(" ") ||
					this.msg.hubEmptySlot,
			),
		);
		console.log(
			this.msg.hubStorage(
				this.theme.resourceName,
				this.stockScrolls,
				this.stockOrbs.map((o) => `[${o}]`).join(" ") || this.msg.hubNone,
			),
		);
		console.log("----------------------------------------------");
		console.log(this.msg.hubMenu1Dungeon);
		console.log(this.msg.hubMenu2Enhance(this.theme.enhanceVerb));
		console.log(this.msg.hubMenu3AttachOrb(this.theme.targetNameLabel));
		console.log(this.msg.hubMenu4Disassemble(this.theme.resourceName));
		console.log(this.msg.hubMenu5Synthesize);
		console.log(this.msg.hubMenu6EnhanceHp(this.theme.resourceName));
		console.log(this.msg.hubMenu7Settings);
		console.log(this.msg.hubMenu0Exit);

		const choice = await this.rl.question(this.msg.chooseAction);

		switch (choice.trim()) {
			case "1":
				await this.chooseDungeonPhase();
				break;
			case "2":
				await this.upgradeWeaponPhase();
				break;
			case "3":
				await this.attachOrbPhase();
				break;
			case "4":
				await this.disassembleOrbPhase();
				break;
			case "5":
				await this.synthesizeOrbPhase();
				break;
			case "6":
				await this.upgradeHpPhase();
				break;
			case "7":
				await this.settingsPhase();
				break;
			case "0":
				console.log(this.msg.farewell);
				this.isRunning = false;
				this.rl.close();
				return;
			default:
				console.log(this.msg.invalidChoice);
		}
	}

	private async settingsPhase(): Promise<void> {
		while (true) {
			console.log(this.msg.settingsTitle);
			console.log(this.msg.settingsCurrent(this.language, this.theme.name));
			console.log(this.msg.settingsMenu1Lang);
			console.log(this.msg.settingsMenu2Theme);
			console.log(this.msg.settingsMenu3ExportTemplate);
			console.log(this.msg.settingsMenu0Back);

			const choice = await this.rl.question(this.msg.chooseAction);
			const trimmed = choice.trim();

			if (trimmed === "0") {
				return;
			}

			if (trimmed === "1") {
				const nextLang: Language = this.language === "ja" ? "en" : "ja";
				const oldDefaultName = this.theme.defaultTargetName;
				this.language = nextLang;
				if (this.customThemeDef && this.theme.id === this.customThemeDef.id) {
					this.theme = resolveTheme(this.customThemeDef, this.language);
				} else {
					this.theme = getPresetTheme(this.theme.id, this.language);
				}
				if (this.player.weapon.name === oldDefaultName) {
					this.player.weapon.name = this.theme.defaultTargetName;
				}
				this.saveGame();
				console.log(this.msg.settingsLangChanged(this.language));
				continue;
			}

			if (trimmed === "2") {
				await this.switchThemeSubmenu();
				continue;
			}

			if (trimmed === "3") {
				const success = exportCustomThemeTemplate();
				if (success) {
					console.log(
						this.msg.settingsTemplateExported(CUSTOM_THEME_EXAMPLE_PATH),
					);
				} else {
					console.log(this.msg.settingsTemplateExportFailed);
				}
				continue;
			}

			console.log(this.msg.invalidChoice);
		}
	}

	private async switchThemeSubmenu(): Promise<void> {
		console.log(this.msg.settingsThemeSelectTitle);

		const themeOptions: Theme[] = [];

		if (this.customThemeDef) {
			themeOptions.push(resolveTheme(this.customThemeDef, this.language));
		} else {
			const checkedCustom = loadCustomThemeDefIfExists();
			if (checkedCustom) {
				this.customThemeDef = checkedCustom;
				themeOptions.push(resolveTheme(checkedCustom, this.language));
			}
		}

		themeOptions.push(
			getPresetTheme("classic_fantasy", this.language),
			getPresetTheme("cyberpunk", this.language),
			getPresetTheme("partner_sync", this.language),
		);

		for (let i = 0; i < themeOptions.length; i++) {
			const t = themeOptions[i];
			if (t) {
				console.log(`${i + 1}: ${t.name} (${t.language}) [${t.id}]`);
			}
		}
		console.log("0: Cancel");

		const choice = await this.rl.question(this.msg.chooseAction);
		const idx = Number.parseInt(choice.trim(), 10) - 1;

		if (idx >= 0 && idx < themeOptions.length) {
			const selected = themeOptions[idx];
			if (selected) {
				const oldDefault = this.theme.defaultTargetName;
				this.theme = selected;
				if (this.player.weapon.name === oldDefault) {
					this.player.weapon.name = selected.defaultTargetName;
				}
				this.saveGame();
				console.log(this.msg.settingsThemeChanged(selected.name));
			}
		}
	}

	private async chooseDungeonPhase(): Promise<void> {
		const starterDungeon = createStarterDungeon(this.theme);
		const deepDungeon = createDeepDungeon(this.theme, this.msg.recommendedDeep);

		console.log(this.msg.dungeonSelectTitle);
		console.log(
			this.msg.dungeonOptionStarter(starterDungeon.name, starterDungeon.floors),
		);
		console.log(
			this.msg.dungeonOptionDeep(
				deepDungeon.name,
				deepDungeon.floors,
				deepDungeon.recommended ?? "",
			),
		);
		console.log(this.msg.dungeonOptionAbyss(this.theme.dungeons.abyss));
		console.log(this.msg.dungeonCancel);

		const choice = await this.rl.question(this.msg.dungeonPrompt);
		switch (choice.trim()) {
			case "1":
				await this.dungeonPhase(starterDungeon);
				break;
			case "2":
				await this.dungeonPhase(deepDungeon);
				break;
			case "3":
				await this.endlessDungeonPhase();
				break;
			case "0":
				console.log(this.msg.dungeonCanceled);
				break;
			default:
				console.log(this.msg.invalidChoice);
				break;
		}
	}

	private async upgradeWeaponPhase(): Promise<void> {
		if (this.stockScrolls <= 0) {
			console.log(this.msg.enhanceNoResource(this.theme.resourceName));
			return;
		}

		console.log(this.msg.enhanceTitle(this.theme.targetNameLabel));
		console.log(
			this.msg.enhanceCurrentRes(this.theme.resourceName, this.stockScrolls),
		);
		console.log(this.msg.enhanceOpt1);
		console.log(this.msg.enhanceOpt2(this.theme.resourceName));
		console.log(
			this.msg.enhanceOpt3(this.theme.resourceName, this.stockScrolls),
		);
		console.log(this.msg.enhanceOpt0);

		const choice = await this.rl.question(this.msg.chooseAction);
		switch (choice.trim()) {
			case "1": {
				this.stockScrolls--;
				this.player.weapon.plus++;
				console.log(
					this.msg.enhanceSuccessSingle(
						this.theme.enhanceVerb,
						this.theme.targetNameLabel,
						this.player.weapon.name,
						this.player.weapon.plus,
						this.getWeaponAtk(),
					),
				);
				this.saveGame();
				break;
			}
			case "2": {
				const inputStr = await this.rl.question(
					this.msg.enhancePromptCount(
						this.theme.resourceName,
						this.stockScrolls,
					),
				);
				const count = Number.parseInt(inputStr.trim(), 10);
				if (Number.isNaN(count) || count < 1 || count > this.stockScrolls) {
					console.log(this.msg.enhanceInvalidCount);
					break;
				}
				this.stockScrolls -= count;
				this.player.weapon.plus += count;
				console.log(
					this.msg.enhanceSuccessMultiple(
						count,
						this.theme.resourceName,
						this.theme.enhanceVerb,
						this.player.weapon.name,
						this.player.weapon.plus,
						this.getWeaponAtk(),
					),
				);
				this.saveGame();
				break;
			}
			case "3": {
				const count = this.stockScrolls;
				this.player.weapon.plus += count;
				this.stockScrolls = 0;
				console.log(
					this.msg.enhanceSuccessAll(
						count,
						this.theme.resourceName,
						this.theme.enhanceVerb,
						this.player.weapon.name,
						this.player.weapon.plus,
						this.getWeaponAtk(),
					),
				);
				this.saveGame();
				break;
			}
			case "0":
				console.log(this.msg.enhanceCanceled);
				break;
			default:
				console.log(this.msg.invalidChoice);
				break;
		}
	}

	private async upgradeHpPhase(): Promise<void> {
		if (this.stockScrolls < 2) {
			console.log(
				this.msg.hpUpgradeNotEnough(this.theme.resourceName, this.stockScrolls),
			);
			return;
		}

		const maxCount = Math.floor(this.stockScrolls / 2);
		console.log(this.msg.hpUpgradeTitle);
		console.log(this.msg.hpUpgradeCurrentHp(this.player.maxHp));
		console.log(
			this.msg.hpUpgradeCurrentRes(this.theme.resourceName, this.stockScrolls),
		);
		console.log(this.msg.hpUpgradeOpt1(this.theme.resourceName));
		console.log(this.msg.hpUpgradeOpt2(this.theme.resourceName, maxCount * 2));
		console.log(
			this.msg.hpUpgradeOpt3(
				this.theme.resourceName,
				maxCount * 2,
				maxCount * 10,
			),
		);
		console.log(this.msg.hpUpgradeOpt0);

		const choice = await this.rl.question(this.msg.chooseAction);
		switch (choice.trim()) {
			case "1": {
				this.stockScrolls -= 2;
				const prevHp = this.player.maxHp;
				this.player.maxHp += 10;
				this.player.hp = this.player.maxHp;
				console.log(this.msg.hpUpgradeSuccessSingle(prevHp, this.player.maxHp));
				this.saveGame();
				break;
			}
			case "2": {
				const inputStr = await this.rl.question(
					this.msg.hpUpgradePromptCount(this.theme.resourceName, maxCount * 2),
				);
				const inputNum = Number.parseInt(inputStr.trim(), 10);
				if (
					Number.isNaN(inputNum) ||
					inputNum < 2 ||
					inputNum > this.stockScrolls
				) {
					console.log(this.msg.hpUpgradeInvalidCount);
					break;
				}
				const times = Math.floor(inputNum / 2);
				const usedScrolls = times * 2;
				const hpGain = times * 10;
				this.stockScrolls -= usedScrolls;
				const prevHp = this.player.maxHp;
				this.player.maxHp += hpGain;
				this.player.hp = this.player.maxHp;
				console.log(
					this.msg.hpUpgradeSuccessMultiple(
						usedScrolls,
						this.theme.resourceName,
						prevHp,
						this.player.maxHp,
					),
				);
				if (inputNum % 2 !== 0) {
					console.log(this.msg.hpUpgradeOddRemainder);
				}
				this.saveGame();
				break;
			}
			case "3": {
				const usedScrolls = maxCount * 2;
				const hpGain = maxCount * 10;
				this.stockScrolls -= usedScrolls;
				const prevHp = this.player.maxHp;
				this.player.maxHp += hpGain;
				this.player.hp = this.player.maxHp;
				console.log(
					this.msg.hpUpgradeSuccessAll(
						usedScrolls,
						this.theme.resourceName,
						prevHp,
						this.player.maxHp,
					),
				);
				this.saveGame();
				break;
			}
			case "0":
				console.log(this.msg.hpUpgradeCanceled);
				break;
			default:
				console.log(this.msg.invalidChoice);
				break;
		}
	}

	private async attachOrbPhase(): Promise<void> {
		if (this.stockOrbs.length === 0) {
			console.log(this.msg.orbAttachNoOrbs);
			return;
		}

		console.log(this.msg.orbAttachTitle);
		for (let i = 0; i < this.stockOrbs.length; i++) {
			const orb = this.stockOrbs[i];
			if (orb) {
				console.log(`${i + 1}: [${orb}] ${this.getOrbDesc(orb)}`);
			}
		}
		console.log("0: Cancel");

		const select = await this.rl.question(this.msg.orbAttachPrompt);
		const idx = Number.parseInt(select.trim(), 10) - 1;

		if (Number.isNaN(idx) || idx < 0 || idx >= this.stockOrbs.length) {
			return;
		}

		const targetOrb = this.stockOrbs[idx];
		if (!targetOrb) {
			return;
		}

		if (this.player.weapon.slots.length < 3) {
			this.player.weapon.slots.push(targetOrb);
			this.stockOrbs.splice(idx, 1);
			console.log(this.msg.orbAttachSuccess(targetOrb));
			this.saveGame();
		} else {
			console.log(this.msg.orbAttachFullTitle);
			for (let sIdx = 0; sIdx < this.player.weapon.slots.length; sIdx++) {
				const s = this.player.weapon.slots[sIdx];
				if (s) {
					console.log(`${sIdx + 1}: [${s}]`);
				}
			}
			console.log("0: Cancel");
			const replaceSelect = await this.rl.question(
				this.msg.orbAttachReplacePrompt,
			);
			const rIdx = Number.parseInt(replaceSelect.trim(), 10) - 1;
			if (!Number.isNaN(rIdx) && rIdx >= 0 && rIdx < 3) {
				const removed = this.player.weapon.slots[rIdx];
				if (removed) {
					this.player.weapon.slots[rIdx] = targetOrb;
					this.stockOrbs.splice(idx, 1);
					console.log(this.msg.orbAttachReplaced(removed, targetOrb));
					this.saveGame();
				}
			}
		}
	}

	private async disassembleOrbPhase(): Promise<void> {
		if (this.stockOrbs.length < 2) {
			console.log(this.msg.orbDisassembleNotEnough);
			return;
		}

		console.log(this.msg.orbDisassembleTitle(this.theme.resourceName));
		console.log(this.msg.orbDisassembleCount(this.stockOrbs.length));
		console.log(this.msg.orbDisassembleOpt1);
		console.log(
			this.msg.orbDisassembleOpt2(
				this.theme.resourceName,
				Math.floor(this.stockOrbs.length / 2),
			),
		);
		console.log(this.msg.orbDisassembleOpt0);

		const menuChoice = await this.rl.question(this.msg.chooseAction);
		if (menuChoice.trim() === "0") {
			console.log(this.msg.orbDisassembleCanceled);
			return;
		}
		if (menuChoice.trim() === "2") {
			const totalOrbs = this.stockOrbs.length;
			const gainedScrolls = Math.floor(totalOrbs / 2);
			const remainder = totalOrbs % 2;
			const remainingOrb =
				remainder === 1 ? this.stockOrbs[totalOrbs - 1] : undefined;

			this.stockScrolls += gainedScrolls;
			this.stockOrbs = remainingOrb ? [remainingOrb] : [];

			console.log(
				this.msg.orbDisassembleAllSuccess(
					totalOrbs - remainder,
					this.theme.resourceName,
					gainedScrolls,
					this.stockScrolls,
				),
			);
			if (remainingOrb) {
				console.log(this.msg.orbDisassembleRemainder(remainingOrb));
			}
			this.saveGame();
			return;
		}
		if (menuChoice.trim() !== "1") {
			console.log(this.msg.invalidChoice);
			return;
		}

		console.log(this.msg.orbDisassembleIndivTitle);
		for (let i = 0; i < this.stockOrbs.length; i++) {
			const orb = this.stockOrbs[i];
			if (orb) {
				console.log(`${i + 1}: [${orb}] ${this.getOrbDesc(orb)}`);
			}
		}
		console.log("0: Cancel");

		const select1 = await this.rl.question(this.msg.orbDisassemblePick1);
		const idx1 = Number.parseInt(select1.trim(), 10) - 1;
		if (Number.isNaN(idx1) || idx1 < 0 || idx1 >= this.stockOrbs.length) {
			return;
		}
		const orb1 = this.stockOrbs[idx1];
		if (!orb1) return;

		console.log(this.msg.orbDisassembleSelected1(orb1));

		const select2 = await this.rl.question(this.msg.orbDisassemblePick2);
		const idx2 = Number.parseInt(select2.trim(), 10) - 1;
		if (Number.isNaN(idx2) || idx2 < 0 || idx2 >= this.stockOrbs.length) {
			return;
		}
		if (idx1 === idx2) {
			console.log(this.msg.orbDisassembleSameError);
			return;
		}
		const orb2 = this.stockOrbs[idx2];
		if (!orb2) return;

		const firstRemoveIdx = Math.max(idx1, idx2);
		const secondRemoveIdx = Math.min(idx1, idx2);
		this.stockOrbs.splice(firstRemoveIdx, 1);
		this.stockOrbs.splice(secondRemoveIdx, 1);

		this.stockScrolls++;
		console.log(
			this.msg.orbDisassembleSuccess(
				orb1,
				orb2,
				this.theme.resourceName,
				this.stockScrolls,
			),
		);
		this.saveGame();
	}

	private async synthesizeOrbPhase(): Promise<void> {
		const normalOrbs: NormalOrbType[] = [
			"MULTI_HIT",
			"CRITICAL",
			"VAMP",
			"POISON",
		];

		const orbCounts = new Map<NormalOrbType, number>();
		for (const orb of this.stockOrbs) {
			if (normalOrbs.includes(orb as NormalOrbType)) {
				const nOrb = orb as NormalOrbType;
				orbCounts.set(nOrb, (orbCounts.get(nOrb) ?? 0) + 1);
			}
		}

		const availableRecipes: Array<{
			from: NormalOrbType;
			to: PlusOrbType;
			count: number;
		}> = [];

		for (const orb of normalOrbs) {
			const count = orbCounts.get(orb) ?? 0;
			if (count >= 2) {
				availableRecipes.push({
					from: orb,
					to: SYNTHESIS_RECIPES[orb],
					count,
				});
			}
		}

		if (availableRecipes.length === 0) {
			console.log(this.msg.orbSynthNoRecipes);
			return;
		}

		console.log(this.msg.orbSynthTitle);
		for (let i = 0; i < availableRecipes.length; i++) {
			const recipe = availableRecipes[i];
			if (recipe) {
				console.log(
					this.msg.orbSynthRecipeItem(
						i + 1,
						recipe.from,
						recipe.count,
						recipe.to,
						this.getOrbDesc(recipe.to),
					),
				);
			}
		}
		console.log("0: Cancel");

		const select = await this.rl.question(this.msg.orbSynthPrompt);
		const idx = Number.parseInt(select.trim(), 10) - 1;
		if (Number.isNaN(idx) || idx < 0 || idx >= availableRecipes.length) {
			return;
		}

		const targetRecipe = availableRecipes[idx];
		if (!targetRecipe) return;

		let removedCount = 0;
		for (let i = this.stockOrbs.length - 1; i >= 0 && removedCount < 2; i--) {
			if (this.stockOrbs[i] === targetRecipe.from) {
				this.stockOrbs.splice(i, 1);
				removedCount++;
			}
		}

		this.stockOrbs.push(targetRecipe.to);
		console.log(this.msg.orbSynthSuccess(targetRecipe.from, targetRecipe.to));
		this.saveGame();
	}

	private async dungeonPhase(dungeon: DungeonDef): Promise<void> {
		console.log(this.msg.dungeonEnter(dungeon.name));
		const runInv: RunInventory = { scrolls: 0, orbs: [] };

		for (let floor = 1; floor <= dungeon.floors; floor++) {
			const template = dungeon.getEnemy(floor);

			const enemy: Enemy = {
				name: template.name,
				hp: template.hp,
				maxHp: template.hp,
				atk: template.atk,
				poison: 0,
			};

			console.log(this.msg.dungeonEncounter(floor, dungeon.floors, enemy.name));

			const survived = await this.battle(enemy);

			if (!survived) {
				console.log(this.msg.dungeonDefeatPrompt);
				console.log(
					this.msg.dungeonItemsLost(
						this.theme.resourceName,
						runInv.scrolls,
						runInv.orbs.length,
					),
				);
				console.log(this.msg.dungeonCarriedBack);
				this.saveGame();
				return;
			}

			// Drops
			const drops = template.getDrops();
			runInv.scrolls += drops.scrolls;
			runInv.orbs.push(...drops.orbs);

			console.log(this.msg.dungeonVictory(enemy.name));
			if (drops.scrolls > 0) {
				console.log(
					this.msg.dungeonDropScroll(this.theme.resourceName, drops.scrolls),
				);
			}
			for (const o of drops.orbs) {
				console.log(this.msg.dungeonDropOrb(o, this.theme.orbs[o]?.name ?? o));
			}

			if (floor === dungeon.floors) {
				console.log(this.msg.dungeonComplete(dungeon.name));
				this.stockScrolls += runInv.scrolls;
				this.stockOrbs.push(...runInv.orbs);
				this.saveGame();
				return;
			}

			console.log(
				this.msg.dungeonCurrentStatus(
					this.player.hp,
					this.player.maxHp,
					this.theme.resourceName,
					runInv.scrolls,
					runInv.orbs.length,
				),
			);
			console.log(this.msg.dungeonNextFloor);
			console.log(this.msg.dungeonRetreat);

			const nextAction = await this.rl.question(this.msg.dungeonPromptAction);
			if (nextAction.trim() === "2") {
				console.log(`\n>> ${this.theme.retreatMessage}`);
				this.stockScrolls += runInv.scrolls;
				this.stockOrbs.push(...runInv.orbs);
				console.log(
					this.msg.dungeonRetreatSuccess(
						this.theme.resourceName,
						runInv.scrolls,
						runInv.orbs.length,
					),
				);
				this.saveGame();
				return;
			}
		}
	}

	private async endlessDungeonPhase(): Promise<void> {
		console.log(this.msg.endlessEnter(this.theme.dungeons.abyss));
		console.log(this.msg.endlessIntro);

		let startFloor = 1;
		if (this.deepestFloor >= 11) {
			const maxStartFloor = Math.floor((this.deepestFloor - 1) / 10) * 10 + 1;
			const availableFloors: number[] = [];
			for (let f = 1; f <= maxStartFloor; f += 10) {
				availableFloors.push(f);
			}

			console.log(this.msg.endlessStartFloorTitle);
			console.log(this.msg.endlessHighestRecord(this.deepestFloor));
			for (let i = 0; i < availableFloors.length; i++) {
				const f = availableFloors[i];
				if (f !== undefined) {
					console.log(`${i + 1}: ${this.msg.endlessStartOption(f)}`);
				}
			}
			console.log(this.msg.endlessCancelOption);

			while (true) {
				const choice = await this.rl.question(this.msg.startPrompt);
				const trimmed = choice.trim();
				if (trimmed === "0") {
					console.log(this.msg.dungeonCanceled);
					return;
				}
				const idx = Number.parseInt(trimmed, 10) - 1;
				if (!Number.isNaN(idx) && idx >= 0 && idx < availableFloors.length) {
					const chosen = availableFloors[idx];
					if (chosen !== undefined) {
						startFloor = chosen;
						break;
					}
				}
				console.log(this.msg.invalidChoice);
			}
		}

		console.log(this.msg.endlessStartFrom(startFloor));
		const runInv: RunInventory = { scrolls: 0, orbs: [] };
		let floor = startFloor;

		while (true) {
			if (floor > this.deepestFloor) {
				this.deepestFloor = floor;
				console.log(this.msg.endlessRecordUpdated(floor));
			}

			const template = generateEndlessEnemy(floor, this.theme);
			const enemy: Enemy = {
				name: template.name,
				hp: template.hp,
				maxHp: template.hp,
				atk: template.atk,
				poison: 0,
			};

			console.log(this.msg.endlessFloorEncounter(floor, enemy.name));

			const survived = await this.battle(enemy);

			if (!survived) {
				console.log(this.msg.dungeonDefeatPrompt);
				console.log(
					this.msg.dungeonItemsLost(
						this.theme.resourceName,
						runInv.scrolls,
						runInv.orbs.length,
					),
				);
				console.log(this.msg.dungeonCarriedBack);
				this.saveGame();
				return;
			}

			const drops = template.getDrops();
			runInv.scrolls += drops.scrolls;
			runInv.orbs.push(...drops.orbs);

			console.log(this.msg.dungeonVictory(enemy.name));
			if (drops.scrolls > 0) {
				console.log(
					this.msg.dungeonDropScroll(this.theme.resourceName, drops.scrolls),
				);
			}
			for (const o of drops.orbs) {
				console.log(this.msg.dungeonDropOrb(o, this.theme.orbs[o]?.name ?? o));
			}

			console.log(
				this.msg.dungeonCurrentStatus(
					this.player.hp,
					this.player.maxHp,
					this.theme.resourceName,
					runInv.scrolls,
					runInv.orbs.length,
				),
			);
			console.log(this.msg.dungeonNextFloor);
			console.log(this.msg.dungeonRetreat);

			const nextAction = await this.rl.question(this.msg.dungeonPromptAction);
			if (nextAction.trim() === "2") {
				console.log(`\n>> ${this.theme.retreatMessage}`);
				this.stockScrolls += runInv.scrolls;
				this.stockOrbs.push(...runInv.orbs);
				console.log(
					this.msg.dungeonRetreatSuccess(
						this.theme.resourceName,
						runInv.scrolls,
						runInv.orbs.length,
					),
				);
				this.saveGame();
				return;
			}

			floor++;
		}
	}

	private async battle(enemy: Enemy): Promise<boolean> {
		const slots = this.player.weapon.slots;
		const hasMultiPlus = slots.includes("MULTI_HIT_PLUS");
		const hasMulti = hasMultiPlus || slots.includes("MULTI_HIT");

		const hasCritPlus = slots.includes("CRITICAL_PLUS");
		const hasCrit = hasCritPlus || slots.includes("CRITICAL");

		const hasVampPlus = slots.includes("VAMP_PLUS");
		const hasVamp = hasVampPlus || slots.includes("VAMP");

		const hasPoisonPlus = slots.includes("POISON_PLUS");
		const hasPoison = hasPoisonPlus || slots.includes("POISON");

		while (this.player.hp > 0 && enemy.hp > 0) {
			console.log(
				this.msg.battleVs(
					this.player.hp,
					this.player.maxHp,
					enemy.name,
					enemy.hp,
					enemy.maxHp,
					enemy.poison,
				),
			);
			console.log(this.msg.battleCmdAttack);
			console.log(this.msg.battleCmdRetreat);
			console.log(this.msg.battleCmdAuto);

			const act = await this.rl.question(this.msg.battleCmdPrompt);
			if (act.trim() === "2") {
				console.log(this.msg.battleRetreatCombat);
				return false;
			}

			if (act.trim() === "3") {
				console.log(this.msg.battleAutoStart);
				let autoTurns = 0;
				let totalDmgDealt = 0;
				let totalDmgTaken = 0;
				let totalHealed = 0;
				let stopReason = "";
				const dangerHp = Math.floor(this.player.maxHp * 0.3);

				while (this.player.hp > 0 && enemy.hp > 0) {
					autoTurns++;

					const baseDmg = this.getWeaponAtk();
					const hits = hasMulti ? 2 : 1;
					const rate = hasMultiPlus ? 0.75 : hasMulti ? 0.65 : 1.0;
					const critChance = hasCritPlus ? 0.4 : hasCrit ? 0.25 : 0;
					const vampRate = hasVampPlus ? 0.25 : hasVamp ? 0.15 : 0;
					const poisonAdd = hasPoisonPlus ? 2 : hasPoison ? 1 : 0;

					for (let i = 1; i <= hits; i++) {
						if (enemy.hp <= 0) break;
						let dmg = Math.max(1, Math.floor(baseDmg * rate));
						if (critChance > 0 && Math.random() < critChance) {
							dmg *= 2;
						}
						enemy.hp -= dmg;
						totalDmgDealt += dmg;

						if (vampRate > 0) {
							const heal = Math.max(1, Math.floor(dmg * vampRate));
							const actualHeal = Math.min(
								this.player.maxHp - this.player.hp,
								heal,
							);
							this.player.hp += actualHeal;
							totalHealed += actualHeal;
						}

						if (poisonAdd > 0) {
							enemy.poison += poisonAdd;
						}
					}

					if (enemy.hp <= 0) {
						stopReason = this.msg.battleAutoReasonEnemyDefeated;
						break;
					}

					// Poison turn end
					if (enemy.poison > 0) {
						const poisonDmg = enemy.poison * 3;
						enemy.hp -= poisonDmg;
						totalDmgDealt += poisonDmg;
						if (enemy.hp <= 0) {
							stopReason = this.msg.battleAutoReasonPoisonDefeated;
							break;
						}
					}

					// Enemy counter
					this.player.hp -= enemy.atk;
					totalDmgTaken += enemy.atk;

					if (this.player.hp <= 0) {
						stopReason = this.msg.battleAutoReasonFainted;
						break;
					}

					// Danger check
					if (this.player.hp <= dangerHp) {
						stopReason = this.msg.battleAutoReasonDangerHp(
							this.player.hp,
							this.player.maxHp,
						);
						break;
					}
				}

				console.log(this.msg.battleAutoSummaryTitle);
				console.log(this.msg.battleAutoReason(stopReason));
				console.log(this.msg.battleAutoTurns(autoTurns));
				console.log(this.msg.battleAutoDamageDealt(totalDmgDealt));
				console.log(this.msg.battleAutoDamageTaken(totalDmgTaken));
				if (totalHealed > 0) {
					console.log(this.msg.battleAutoHealed(totalHealed));
				}
				console.log(
					this.msg.battleAutoRemainingHp(
						this.player.hp,
						this.player.maxHp,
						enemy.name,
						enemy.hp,
						enemy.maxHp,
					),
				);
				console.log("==============================================");

				if (enemy.hp <= 0) {
					return true;
				}
				if (this.player.hp <= 0) {
					return false;
				}
				continue;
			}

			if (act.trim() !== "1") {
				console.log(this.msg.battleInvalidCmd);
				continue;
			}

			// Attack
			const baseDmg = this.getWeaponAtk();
			const hits = hasMulti ? 2 : 1;
			const rate = hasMultiPlus ? 0.75 : hasMulti ? 0.65 : 1.0;
			const critChance = hasCritPlus ? 0.4 : hasCrit ? 0.25 : 0;
			const vampRate = hasVampPlus ? 0.25 : hasVamp ? 0.15 : 0;
			const poisonAdd = hasPoisonPlus ? 2 : hasPoison ? 1 : 0;

			for (let i = 1; i <= hits; i++) {
				if (enemy.hp <= 0) break;
				let isCrit = false;
				let dmg = Math.max(1, Math.floor(baseDmg * rate));

				if (critChance > 0 && Math.random() < critChance) {
					isCrit = true;
					dmg *= 2;
				}

				enemy.hp -= dmg;
				const hitIndexStr = hits > 1 ? this.msg.battleHitNumber(i) : "";
				console.log(
					this.msg.battlePlayerAttack(hitIndexStr, isCrit, enemy.name, dmg),
				);

				if (vampRate > 0) {
					const heal = Math.max(1, Math.floor(dmg * vampRate));
					this.player.hp = Math.min(this.player.maxHp, this.player.hp + heal);
					console.log(this.msg.battleVampHeal(heal, this.player.hp));
				}

				if (poisonAdd > 0) {
					enemy.poison += poisonAdd;
					console.log(
						this.msg.battlePoisonInflict(enemy.name, poisonAdd, enemy.poison),
					);
				}
			}

			if (enemy.hp <= 0) {
				return true;
			}

			// Poison tick
			if (enemy.poison > 0) {
				const poisonDmg = enemy.poison * 3;
				enemy.hp -= poisonDmg;
				console.log(this.msg.battlePoisonTick(enemy.name, poisonDmg));
				if (enemy.hp <= 0) return true;
			}

			// Enemy counter
			console.log(this.msg.battleEnemyCounter(enemy.name, enemy.atk));
			this.player.hp -= enemy.atk;

			if (this.player.hp <= 0) {
				return false;
			}
		}

		return this.player.hp > 0;
	}
}

const isMain = argv[1] && resolve(argv[1]) === fileURLToPath(import.meta.url);
if (isMain) {
	new Game().start().catch(console.error);
}
