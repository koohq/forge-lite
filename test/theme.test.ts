import assert from "node:assert/strict";
import { existsSync, unlinkSync } from "node:fs";
import test from "node:test";
import {
	exportCustomThemeTemplate,
	getPresetTheme,
	loadCustomThemeDefIfExists,
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
});
