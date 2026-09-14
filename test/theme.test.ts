import assert from "node:assert/strict";
import { existsSync, unlinkSync } from "node:fs";
import test from "node:test";
import {
	exportCustomThemeTemplate,
	getPresetTheme,
	loadCustomThemeDefIfExists,
	MESSAGES,
	PRESET_CLASSIC_FANTASY,
	PRESET_CYBERPUNK,
	PRESET_PARTNER_SYNC,
	PRESETS,
	resolveTheme,
	validateThemeDefinition,
} from "../src/game.ts";

test("Theme Presets: All 3 presets are registered", () => {
	assert.ok(PRESETS.classic_fantasy);
	assert.ok(PRESETS.cyberpunk);
	assert.ok(PRESETS.partner_sync);
});

test("Classic Fantasy: Japanese and English strings match requirements", () => {
	const ja = PRESET_CLASSIC_FANTASY.ja;
	const en = PRESET_CLASSIC_FANTASY.en;

	// Target
	assert.equal(ja.defaultTargetName, "どうのつるぎ");
	assert.equal(en.defaultTargetName, "Bronze Sword");

	// Verb
	assert.equal(ja.enhanceVerb, "鍛冶屋で鍛える");
	assert.equal(en.enhanceVerb, "Forge at Blacksmith");

	// Resource
	assert.equal(ja.resourceName, "強化の書");
	assert.equal(en.resourceName, "Upgrade Scroll");

	// Orbs
	assert.equal(ja.orbs.MULTI_HIT.name, "連撃の印");
	assert.equal(en.orbs.MULTI_HIT.name, "Twin Strike");
	assert.equal(ja.orbs.CRITICAL.name, "会心の印");
	assert.equal(en.orbs.CRITICAL.name, "Critical Strike");
	assert.equal(ja.orbs.VAMP.name, "吸血の印");
	assert.equal(en.orbs.VAMP.name, "Vampiric Drain");
	assert.equal(ja.orbs.POISON.name, "猛毒の印");
	assert.equal(en.orbs.POISON.name, "Deadly Poison");

	// Dungeons
	assert.equal(ja.dungeons.starter, "始まりの洞窟");
	assert.equal(en.dungeons.starter, "Cave of Beginnings");
	assert.equal(ja.dungeons.deep, "灼熱の深層");
	assert.equal(en.dungeons.deep, "Scorching Depths");
	assert.equal(ja.dungeons.abyss, "無限の深淵");
	assert.equal(en.dungeons.abyss, "Infinite Abyss");

	// Vocabulary
	assert.equal(ja.hubTitle, "【冒険者の酒場】");
	assert.equal(en.hubTitle, "【Adventurer's Guild】");
	assert.equal(ja.hpLabel, "最大体力");
	assert.equal(en.hpLabel, "Max HP");
	assert.equal(ja.targetPrefix, "所持装備");
	assert.equal(en.targetPrefix, "Equipment");
	assert.equal(ja.statLabel, "攻撃力");
	assert.equal(en.statLabel, "ATK");
	assert.equal(ja.plusPrefix, "+");
	assert.equal(en.plusPrefix, "+");
	assert.equal(ja.storageLabel, "倉庫");
	assert.equal(en.storageLabel, "Storage");
	assert.equal(ja.orbLabel, "オーブ");
	assert.equal(en.orbLabel, "Orb");
	assert.equal(ja.enhanceHpVerb, "体力を強化する");
	assert.equal(en.enhanceHpVerb, "Fortify Vitality");
	assert.equal(ja.installVerb, "装着");
	assert.equal(en.installVerb, "Equip");
});

test("Cyberpunk: Japanese and English strings match requirements", () => {
	const ja = PRESET_CYBERPUNK.ja;
	const en = PRESET_CYBERPUNK.en;

	// Target
	assert.equal(ja.defaultTargetName, "パルスブレード");
	assert.equal(en.defaultTargetName, "Pulse Blade");

	// Verb
	assert.equal(ja.enhanceVerb, "システムオーバークロック");
	assert.equal(en.enhanceVerb, "Overclock System");

	// Resource
	assert.equal(ja.resourceName, "ナノチップ");
	assert.equal(en.resourceName, "Nanite Chip");

	// Orbs
	assert.equal(ja.orbs.MULTI_HIT.name, "多段バースト");
	assert.equal(en.orbs.MULTI_HIT.name, "Multi-Burst");
	assert.equal(ja.orbs.CRITICAL.name, "クリティカル注入");
	assert.equal(en.orbs.CRITICAL.name, "Critical Injection");
	assert.equal(ja.orbs.VAMP.name, "ナノ修復ドレイン");
	assert.equal(en.orbs.VAMP.name, "Nano-Drain");
	assert.equal(ja.orbs.POISON.name, "神経汚染");
	assert.equal(en.orbs.POISON.name, "Neuro-Toxin");

	// Dungeons
	assert.equal(ja.dungeons.starter, "閉鎖ネットワーク");
	assert.equal(en.dungeons.starter, "Subnet Alpha");
	assert.equal(ja.dungeons.deep, "企業中枢サーバー");
	assert.equal(en.dungeons.deep, "Corporate Core");
	assert.equal(ja.dungeons.abyss, "無限の電脳網");
	assert.equal(en.dungeons.abyss, "Infinite Cyberspace");

	// Vocabulary
	assert.equal(ja.hubTitle, "【セーフハウス・端末】");
	assert.equal(en.hubTitle, "【Safehouse Terminal】");
	assert.equal(ja.hpLabel, "最大耐久");
	assert.equal(en.hpLabel, "Max Hull");
	assert.equal(ja.targetPrefix, "主兵装");
	assert.equal(en.targetPrefix, "Main Armament");
	assert.equal(ja.statLabel, "出力");
	assert.equal(en.statLabel, "Output");
	assert.equal(ja.plusPrefix, "+");
	assert.equal(en.plusPrefix, "+");
	assert.equal(ja.storageLabel, "ストレージ");
	assert.equal(en.storageLabel, "Storage");
	assert.equal(ja.orbLabel, "モジュール");
	assert.equal(en.orbLabel, "Module");
	assert.equal(ja.enhanceHpVerb, "生体フレーム拡張");
	assert.equal(en.enhanceHpVerb, "Upgrade Chassis");
	assert.equal(ja.installVerb, "インストール");
	assert.equal(en.installVerb, "Install");
});

test("Partner Sync: Japanese and English strings match requirements", () => {
	const ja = PRESET_PARTNER_SYNC.ja;
	const en = PRESET_PARTNER_SYNC.en;

	// Target
	assert.equal(ja.defaultTargetName, "戦術アンドロイド「アイリス」");
	assert.equal(en.defaultTargetName, 'Tactical Android "Iris"');

	// Verb
	assert.equal(ja.enhanceVerb, "同期率を向上させる");
	assert.equal(en.enhanceVerb, "Deepen Sync");

	// Resource
	assert.equal(ja.resourceName, "メモリコア");
	assert.equal(en.resourceName, "Memory Core");

	// Orbs
	assert.equal(ja.orbs.MULTI_HIT.name, "連携戦術");
	assert.equal(en.orbs.MULTI_HIT.name, "Tandem Tactics");
	assert.equal(ja.orbs.CRITICAL.name, "弱点看破");
	assert.equal(en.orbs.CRITICAL.name, "Exploit Weakness");
	assert.equal(ja.orbs.VAMP.name, "自己修復");
	assert.equal(en.orbs.VAMP.name, "Self-Repair");
	assert.equal(ja.orbs.POISON.name, "浸食ノイズ");
	assert.equal(en.orbs.POISON.name, "Corrosive Noise");

	// Dungeons
	assert.equal(ja.dungeons.starter, "廃墟区域");
	assert.equal(en.dungeons.starter, "Ruined Sector");
	assert.equal(ja.dungeons.deep, "汚染中枢");
	assert.equal(en.dungeons.deep, "Contaminated Core");
	assert.equal(ja.dungeons.abyss, "未知の最深部");
	assert.equal(en.dungeons.abyss, "Unknown Depths");

	// Retreat
	assert.equal(ja.retreatMessage, "アイリスの負荷を考慮し、一時帰還した。");
	assert.equal(en.retreatMessage, "Retreated to base to reduce load on Iris.");

	// Vocabulary
	assert.equal(ja.hubTitle, "【司令室・ドック】");
	assert.equal(en.hubTitle, "【Command Dock】");
	assert.equal(ja.hpLabel, "生存限界");
	assert.equal(en.hpLabel, "Vitality");
	assert.equal(ja.targetPrefix, "相棒");
	assert.equal(en.targetPrefix, "Partner");
	assert.equal(ja.statLabel, "戦闘能力");
	assert.equal(en.statLabel, "Combat Power");
	assert.equal(ja.plusPrefix, " Sync:+");
	assert.equal(en.plusPrefix, " Sync:+");
	assert.equal(ja.storageLabel, "データバンク");
	assert.equal(en.storageLabel, "Data Bank");
	assert.equal(ja.orbLabel, "プロトコル");
	assert.equal(en.orbLabel, "Protocol");
	assert.equal(ja.enhanceHpVerb, "防護プロトコル強化");
	assert.equal(en.enhanceHpVerb, "Reinforce Shields");
	assert.equal(ja.installVerb, "セット");
	assert.equal(en.installVerb, "Set");
});

test("Theme Resolution: resolveTheme switches language instantly", () => {
	const cyberJa = resolveTheme(PRESET_CYBERPUNK, "ja");
	const cyberEn = resolveTheme(PRESET_CYBERPUNK, "en");

	assert.equal(cyberJa.language, "ja");
	assert.equal(cyberJa.defaultTargetName, "パルスブレード");
	assert.equal(cyberJa.enhanceVerb, "システムオーバークロック");

	assert.equal(cyberEn.language, "en");
	assert.equal(cyberEn.defaultTargetName, "Pulse Blade");
	assert.equal(cyberEn.enhanceVerb, "Overclock System");

	const partnerJa = getPresetTheme("partner_sync", "ja");
	const partnerEn = getPresetTheme("partner_sync", "en");

	assert.equal(partnerJa.resourceName, "メモリコア");
	assert.equal(partnerEn.resourceName, "Memory Core");
});

test("Custom Theme Template: Exports valid bilingual template", () => {
	const testPath = "test_custom_theme.example.json";
	if (existsSync(testPath)) unlinkSync(testPath);

	const success = exportCustomThemeTemplate(testPath);
	assert.ok(success);
	assert.ok(existsSync(testPath));

	const loaded = loadCustomThemeDefIfExists(testPath);
	assert.ok(loaded);
	assert.equal(loaded.id, "custom_scifi");
	assert.equal(loaded.ja.resourceName, "プラズマコア");
	assert.equal(loaded.en.resourceName, "Plasma Core");

	unlinkSync(testPath);
});

test("Validation: Backward compatibility with legacy flat format", () => {
	const legacyTheme = {
		id: "legacy_theme",
		name: "Legacy Theme",
		language: "en",
		targetNameLabel: "Old Weapon",
		defaultTargetName: "Old Stick",
		enhanceVerb: "Old Polish",
		resourceName: "Old Dust",
		orbs: {
			MULTI_HIT: { name: "M", desc: "m" },
			CRITICAL: { name: "C", desc: "c" },
			VAMP: { name: "V", desc: "v" },
			POISON: { name: "P", desc: "p" },
			MULTI_HIT_PLUS: { name: "M+", desc: "m+" },
			CRITICAL_PLUS: { name: "C+", desc: "c+" },
			VAMP_PLUS: { name: "V+", desc: "v+" },
			POISON_PLUS: { name: "P+", desc: "p+" },
		},
		dungeons: { starter: "S", deep: "D", abyss: "A" },
		enemyPrefixes: ["E1"],
		enemyBases: ["B1"],
		retreatMessage: "Retreat",
	};

	const validated = validateThemeDefinition(legacyTheme);
	assert.ok(validated);
	assert.equal(validated.id, "legacy_theme");
	assert.equal(validated.ja.defaultTargetName, "Old Stick");
	assert.equal(validated.en.defaultTargetName, "Old Stick");

	// Fallback check for new vocabulary fields
	assert.equal(validated.ja.hubTitle, "【拠点】");
	assert.equal(validated.ja.targetPrefix, "Old Weapon");
	assert.equal(validated.ja.orbLabel, "オーブ");
	assert.equal(validated.ja.hpLabel, "最大体力");
	assert.equal(validated.ja.enhanceHpVerb, "体力を強化する");
	assert.equal(validated.ja.installVerb, "装着");
});

test("Hub Display & Menu: Dynamic vocabulary correctly delegates per theme", () => {
	const fantasyJa = resolveTheme(PRESET_CLASSIC_FANTASY, "ja");
	assert.equal(MESSAGES.ja.hubHeader(fantasyJa.hubTitle), "【冒険者の酒場】");
	assert.equal(
		MESSAGES.ja.hubStats(fantasyJa.hpLabel, 100, 5),
		"最大体力: HP 100 | 無限の深淵 最高到達: B5F",
	);
	assert.equal(
		MESSAGES.ja.hubEquip(
			fantasyJa.targetPrefix,
			"どうのつるぎ",
			fantasyJa.plusPrefix,
			10,
			fantasyJa.statLabel,
			20,
		),
		"所持装備: どうのつるぎ+10 (攻撃力: 20)",
	);
	assert.equal(
		MESSAGES.ja.hubSlots(fantasyJa.installVerb, 1, "[連撃の印]"),
		"装着スロット [1/3]: [連撃の印]",
	);
	assert.equal(
		MESSAGES.ja.hubStorage(
			fantasyJa.storageLabel,
			fantasyJa.resourceName,
			3,
			fantasyJa.orbLabel,
			"[連撃の印]",
		),
		"倉庫: 強化の書 x3 | 未装着オーブ: [連撃の印]",
	);
	assert.equal(MESSAGES.ja.hubMenu1Dungeon, "1: ダンジョンへ出撃 (探索開始)");
	assert.equal(
		MESSAGES.ja.hubMenu2Enhance(fantasyJa.enhanceVerb),
		"2: 鍛冶屋で鍛える (攻撃力・出力の強化)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu3AttachOrb(fantasyJa.installVerb),
		"3: 装着 (パッシブ効果の装着)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu4Disassemble(fantasyJa.orbLabel),
		"4: オーブを分解 (素材への還元)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu5Synthesize(fantasyJa.orbLabel),
		"5: オーブを合成 (上位性能への強化)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu6EnhanceHp(fantasyJa.enhanceHpVerb, fantasyJa.hpLabel),
		"6: 体力を強化する (最大体力 +10)",
	);
	assert.equal(MESSAGES.ja.hubMenu7Settings, "7: 設定 / Settings");
	assert.equal(MESSAGES.ja.hubMenu0Exit, "0: ゲーム終了");

	const cyberJa = resolveTheme(PRESET_CYBERPUNK, "ja");
	assert.equal(
		MESSAGES.ja.hubHeader(cyberJa.hubTitle),
		"【セーフハウス・端末】",
	);
	assert.equal(
		MESSAGES.ja.hubStats(cyberJa.hpLabel, 150, 12),
		"最大耐久: HP 150 | 無限の深淵 最高到達: B12F",
	);
	assert.equal(
		MESSAGES.ja.hubEquip(
			cyberJa.targetPrefix,
			"パルスブレード",
			cyberJa.plusPrefix,
			5,
			cyberJa.statLabel,
			35,
		),
		"主兵装: パルスブレード+5 (出力: 35)",
	);
	assert.equal(
		MESSAGES.ja.hubSlots(
			cyberJa.installVerb,
			2,
			"[多段バースト] [クリティカル注入]",
		),
		"インストールスロット [2/3]: [多段バースト] [クリティカル注入]",
	);
	assert.equal(
		MESSAGES.ja.hubStorage(
			cyberJa.storageLabel,
			cyberJa.resourceName,
			8,
			cyberJa.orbLabel,
			"[多段バースト]",
		),
		"ストレージ: ナノチップ x8 | 未装着モジュール: [多段バースト]",
	);
	assert.equal(
		MESSAGES.ja.hubMenu2Enhance(cyberJa.enhanceVerb),
		"2: システムオーバークロック (攻撃力・出力の強化)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu3AttachOrb(cyberJa.installVerb),
		"3: インストール (パッシブ効果の装着)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu4Disassemble(cyberJa.orbLabel),
		"4: モジュールを分解 (素材への還元)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu5Synthesize(cyberJa.orbLabel),
		"5: モジュールを合成 (上位性能への強化)",
	);
	assert.equal(
		MESSAGES.ja.hubMenu6EnhanceHp(cyberJa.enhanceHpVerb, cyberJa.hpLabel),
		"6: 生体フレーム拡張 (最大耐久 +10)",
	);

	const partnerEn = resolveTheme(PRESET_PARTNER_SYNC, "en");
	assert.equal(MESSAGES.en.hubHeader(partnerEn.hubTitle), "【Command Dock】");
	assert.equal(
		MESSAGES.en.hubStats(partnerEn.hpLabel, 200, 20),
		"Vitality: 200 | Infinite Abyss Record: B20F",
	);
	assert.equal(
		MESSAGES.en.hubEquip(
			partnerEn.targetPrefix,
			'Tactical Android "Iris"',
			partnerEn.plusPrefix,
			15,
			partnerEn.statLabel,
			50,
		),
		'Partner: Tactical Android "Iris" Sync:+15 (Combat Power: 50)',
	);
	assert.equal(
		MESSAGES.en.hubSlots(partnerEn.installVerb, 1, "[Tandem Tactics]"),
		"Set Slots [1/3]: [Tandem Tactics]",
	);
	assert.equal(
		MESSAGES.en.hubStorage(
			partnerEn.storageLabel,
			partnerEn.resourceName,
			12,
			partnerEn.orbLabel,
			"[Tandem Tactics]",
		),
		"Data Bank: Memory Core x12 | Stock Protocols: [Tandem Tactics]",
	);
	assert.equal(MESSAGES.en.hubMenu1Dungeon, "1: Embark to Dungeon (Start run)");
	assert.equal(
		MESSAGES.en.hubMenu2Enhance(partnerEn.enhanceVerb),
		"2: Deepen Sync (Upgrade ATK / Output)",
	);
	assert.equal(
		MESSAGES.en.hubMenu3AttachOrb(partnerEn.installVerb),
		"3: Set (Equip passive effects)",
	);
	assert.equal(
		MESSAGES.en.hubMenu4Disassemble(partnerEn.orbLabel),
		"4: Dismantle Protocol (Convert to materials)",
	);
	assert.equal(
		MESSAGES.en.hubMenu5Synthesize(partnerEn.orbLabel),
		"5: Synthesize Protocol (Upgrade to Plus version)",
	);
	assert.equal(
		MESSAGES.en.hubMenu6EnhanceHp(partnerEn.enhanceHpVerb, partnerEn.hpLabel),
		"6: Reinforce Shields (Vitality +10)",
	);
	assert.equal(MESSAGES.en.hubMenu7Settings, "7: Settings");
	assert.equal(MESSAGES.en.hubMenu0Exit, "0: Quit Game");
});
