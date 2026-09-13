import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { stdin as input, stdout as output } from "node:process";
import * as readline from "node:readline/promises";

const SAVE_FILE_PATH = "save.json";

type NormalOrbType = "MULTI_HIT" | "CRITICAL" | "VAMP" | "POISON";
type PlusOrbType =
	| "MULTI_HIT_PLUS"
	| "CRITICAL_PLUS"
	| "VAMP_PLUS"
	| "POISON_PLUS";
type OrbType = NormalOrbType | PlusOrbType;

interface Weapon {
	name: string;
	baseAtk: number;
	plus: number;
	slots: OrbType[];
}

interface Player {
	maxHp: number;
	hp: number;
	weapon: Weapon;
}

interface Enemy {
	name: string;
	hp: number;
	maxHp: number;
	atk: number;
	poison: number;
}

interface RunInventory {
	scrolls: number;
	orbs: OrbType[];
}

interface SaveData {
	weapon: Weapon;
	stockScrolls: number;
	stockOrbs: OrbType[];
}

interface EnemyTemplate {
	name: string;
	hp: number;
	atk: number;
	getDrops: () => RunInventory;
}

interface DungeonDef {
	name: string;
	floors: number;
	recommended?: string | undefined;
	getEnemy: (floor: number) => EnemyTemplate;
}

const ORB_NAMES: Record<OrbType, string> = {
	MULTI_HIT: "連撃の印 (2回攻撃/威力65%)",
	CRITICAL: "会心の印 (25%で2倍)",
	VAMP: "吸血の印 (与ダメの15%回復)",
	POISON: "猛毒の印 (攻撃時毒+1/ターン末毒x3ダメ)",
	MULTI_HIT_PLUS: "連撃の印+ (2回攻撃/威力75%)",
	CRITICAL_PLUS: "会心の印+ (40%で2倍)",
	VAMP_PLUS: "吸血の印+ (与ダメの25%回復)",
	POISON_PLUS: "猛毒の印+ (攻撃時毒+2/ターン末毒x3ダメ)",
};

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

const FLOOR_ENEMIES_CAVE: readonly EnemyTemplate[] = [
	{
		name: "スライム",
		hp: 15,
		atk: 3,
		getDrops: () => ({ scrolls: 1, orbs: [] }),
	},
	{
		name: "コボルト",
		hp: 28,
		atk: 6,
		getDrops: () => ({ scrolls: 1, orbs: [] }),
	},
	{
		name: "オーク",
		hp: 45,
		atk: 10,
		getDrops: () => ({ scrolls: 0, orbs: [getRandomOrb()] }),
	},
	{
		name: "ゴーレム",
		hp: 75,
		atk: 14,
		getDrops: () => ({ scrolls: 2, orbs: [] }),
	},
	{
		name: "レッドドラゴン (BOSS)",
		hp: 130,
		atk: 20,
		getDrops: () => ({ scrolls: 3, orbs: [getRandomOrb()] }),
	},
];

const CAVE_DUNGEON: DungeonDef = {
	name: "始まりの洞窟",
	floors: 5,
	getEnemy: (floor: number): EnemyTemplate => {
		const template = FLOOR_ENEMIES_CAVE[floor - 1];
		if (!template) {
			throw new Error(`Invalid floor: ${floor}`);
		}
		return template;
	},
};

const DEEP_DUNGEON: DungeonDef = {
	name: "灼熱の深層",
	floors: 10,
	recommended: "推奨+25以上の上級ダンジョン",
	getEnemy: (floor: number): EnemyTemplate => {
		const getRandomItem = <T>(items: readonly T[]): T => {
			const item = items[Math.floor(Math.random() * items.length)];
			if (!item) {
				throw new Error("Empty array");
			}
			return item;
		};

		const getDeepFloorScrolls = () => Math.floor(Math.random() * 4) + 2; // 2〜5枚

		if (floor >= 1 && floor <= 3) {
			const candidates = [
				{ name: "狂暴なオーク", hp: 80, atk: 12 },
				{ name: "サンダードラゴン", hp: 120, atk: 16 },
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
				{ name: "アイアンゴーレム", hp: 180, atk: 20 },
				{ name: "デスナイト", hp: 220, atk: 24 },
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
				{ name: "エンシェントワイバーン", hp: 280, atk: 28 },
				{ name: "ベヒモス", hp: 350, atk: 32 },
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
				name: "冥王 (BOSS)",
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

class Game {
	private readonly rl: readline.Interface;
	private readonly player: Player;
	private stockScrolls = 0;
	private stockOrbs: OrbType[] = [];
	private isRunning = true;

	constructor() {
		this.rl = readline.createInterface({ input, output });
		this.player = {
			maxHp: 50,
			hp: 50,
			weapon: {
				name: "どうのつるぎ",
				baseAtk: 10,
				plus: 0,
				slots: [],
			},
		};
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
			};
			writeFileSync(SAVE_FILE_PATH, JSON.stringify(data, null, 2), "utf-8");
			console.log(">> セーブデータを保存しました。(save.json)");
		} catch (error) {
			console.error(">> セーブデータの保存に失敗しました:", error);
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
				this.player.weapon = {
					name: data.weapon.name,
					baseAtk: data.weapon.baseAtk,
					plus: data.weapon.plus,
					slots: Array.isArray(data.weapon.slots) ? [...data.weapon.slots] : [],
				};
				this.stockScrolls = data.stockScrolls;
				this.stockOrbs = [...data.stockOrbs];
				return true;
			}
			return false;
		} catch (error) {
			console.error(">> セーブデータの読み込みに失敗しました:", error);
			return false;
		}
	}

	async start(): Promise<void> {
		console.log("==============================================");
		console.log("   Minimal Rogue-lite Prototype (CUI Ver)   ");
		console.log("==============================================");

		console.log("\n1: つづきから (save.json を読み込んで開始)");
		console.log("2: はじめから (初期状態で開始)");

		while (true) {
			const choice = await this.rl.question("\n選択してください: ");
			const trimmed = choice.trim();
			if (trimmed === "1") {
				if (this.loadGame()) {
					console.log(">> セーブデータを読み込みました！");
				} else {
					console.log(
						">> save.json が見つからないか破損しています。新規データで開始します。",
					);
				}
				break;
			}
			if (trimmed === "2") {
				console.log(">> はじめからゲームを開始します。");
				break;
			}
			console.log(">> 無効な選択です。1 または 2 を入力してください。");
		}

		while (this.isRunning) {
			await this.hubPhase();
		}
	}

	private async hubPhase(): Promise<void> {
		this.player.hp = this.player.maxHp;
		console.log("\n----------------------------------------------");
		console.log("【拠点】");
		console.log(
			`所持装備: ${this.player.weapon.name}+${this.player.weapon.plus} (攻撃力: ${this.getWeaponAtk()})`,
		);
		console.log(
			`装着スロット [${this.player.weapon.slots.length}/3]: ${
				this.player.weapon.slots.map((s) => `[${s}]`).join(" ") || "(空き)"
			}`,
		);
		console.log(
			`倉庫: 強化の書 x${this.stockScrolls} | 未装着オーブ: ${
				this.stockOrbs.map((o) => `[${o}]`).join(" ") || "(なし)"
			}`,
		);
		console.log("----------------------------------------------");
		console.log("1: ダンジョンへ出撃");
		console.log("2: 武器を鍛える");
		console.log("3: オーブを武器に装着");
		console.log("4: オーブを分解 (任意のオーブ2個 -> 強化の書1枚)");
		console.log("5: オーブを合成 (同種オーブ2個 -> 上位オーブ)");
		console.log("0: ゲーム終了");

		const choice = await this.rl.question("\n行動を選択してください: ");

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
			case "0":
				console.log("お疲れ様でした。");
				this.isRunning = false;
				this.rl.close();
				return;
			default:
				console.log("無効な選択です。");
		}
	}

	private async chooseDungeonPhase(): Promise<void> {
		console.log("\n--- 出撃ダンジョン選択 ---");
		console.log(`1: ${CAVE_DUNGEON.name} (全${CAVE_DUNGEON.floors}階 / 初級)`);
		console.log(
			`2: ${DEEP_DUNGEON.name} (全${DEEP_DUNGEON.floors}階 / ${DEEP_DUNGEON.recommended})`,
		);
		console.log("0: キャンセル (拠点に戻る)");

		const choice = await this.rl.question("\nダンジョンを選択してください: ");
		switch (choice.trim()) {
			case "1":
				await this.dungeonPhase(CAVE_DUNGEON);
				break;
			case "2":
				await this.dungeonPhase(DEEP_DUNGEON);
				break;
			case "0":
				console.log(">> 出撃を取りやめました。");
				break;
			default:
				console.log(">> 無効な選択です。");
				break;
		}
	}

	private async upgradeWeaponPhase(): Promise<void> {
		if (this.stockScrolls <= 0) {
			console.log(">> 強化の書がありません！");
			return;
		}

		console.log("\n--- 武器の強化 ---");
		console.log(`現在の所持強化の書: ${this.stockScrolls}枚`);
		console.log("1: 1枚だけ消費して鍛える (+1)");
		console.log(
			`2: 所持している強化の書をすべて消費して一括で鍛える (+${this.stockScrolls})`,
		);
		console.log("0: キャンセル");

		const choice = await this.rl.question("選択してください: ");
		switch (choice.trim()) {
			case "1": {
				this.stockScrolls--;
				this.player.weapon.plus++;
				console.log(
					`>> 武器を鍛えました！ ${this.player.weapon.name}+${this.player.weapon.plus} (攻撃力: ${this.getWeaponAtk()})`,
				);
				this.saveGame();
				break;
			}
			case "2": {
				const count = this.stockScrolls;
				this.player.weapon.plus += count;
				this.stockScrolls = 0;
				console.log(
					`>> 強化の書をすべて(${count}枚)消費して一括で鍛えました！ ${this.player.weapon.name}+${this.player.weapon.plus} (攻撃力: ${this.getWeaponAtk()})`,
				);
				this.saveGame();
				break;
			}
			case "0":
				console.log(">> 強化をキャンセルしました。");
				break;
			default:
				console.log(">> 無効な選択です。");
				break;
		}
	}

	private async attachOrbPhase(): Promise<void> {
		if (this.stockOrbs.length === 0) {
			console.log(">> 装着できるオーブが倉庫にありません！");
			return;
		}

		console.log("\n--- 倉庫のオーブ一覧 ---");
		for (let i = 0; i < this.stockOrbs.length; i++) {
			const orb = this.stockOrbs[i];
			if (orb) {
				console.log(`${i + 1}: [${orb}] ${ORB_NAMES[orb]}`);
			}
		}
		console.log("0: キャンセル");

		const select = await this.rl.question("装着するオーブの番号を選択: ");
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
			console.log(`>> スロットに [${targetOrb}] を装着しました！`);
			this.saveGame();
		} else {
			console.log("\nスロットが満杯です。上書きする枠を選んでください:");
			for (let sIdx = 0; sIdx < this.player.weapon.slots.length; sIdx++) {
				const s = this.player.weapon.slots[sIdx];
				if (s) {
					console.log(`${sIdx + 1}: 上書き対象 -> [${s}]`);
				}
			}
			console.log("0: キャンセル");
			const replaceSelect = await this.rl.question("番号を選択: ");
			const rIdx = Number.parseInt(replaceSelect.trim(), 10) - 1;
			if (!Number.isNaN(rIdx) && rIdx >= 0 && rIdx < 3) {
				const removed = this.player.weapon.slots[rIdx];
				if (removed) {
					this.player.weapon.slots[rIdx] = targetOrb;
					this.stockOrbs.splice(idx, 1);
					console.log(
						`>> [${removed}] を破棄し、[${targetOrb}] を装着しました！`,
					);
					this.saveGame();
				}
			}
		}
	}

	private async disassembleOrbPhase(): Promise<void> {
		if (this.stockOrbs.length < 2) {
			console.log(">> 分解には倉庫にオーブが2個以上必要です！");
			return;
		}

		console.log("\n--- オーブの分解 (任意のオーブ2個 -> 強化の書1枚) ---");
		for (let i = 0; i < this.stockOrbs.length; i++) {
			const orb = this.stockOrbs[i];
			if (orb) {
				console.log(`${i + 1}: [${orb}] ${ORB_NAMES[orb]}`);
			}
		}
		console.log("0: キャンセル");

		const select1 = await this.rl.question(
			"1つ目に分解するオーブの番号を選択: ",
		);
		const idx1 = Number.parseInt(select1.trim(), 10) - 1;
		if (Number.isNaN(idx1) || idx1 < 0 || idx1 >= this.stockOrbs.length) {
			return;
		}
		const orb1 = this.stockOrbs[idx1];
		if (!orb1) return;

		console.log(`>> 1つ目: [${orb1}] を選択しました。`);

		const select2 = await this.rl.question(
			"2つ目に分解するオーブの番号を選択: ",
		);
		const idx2 = Number.parseInt(select2.trim(), 10) - 1;
		if (Number.isNaN(idx2) || idx2 < 0 || idx2 >= this.stockOrbs.length) {
			return;
		}
		if (idx1 === idx2) {
			console.log(">> 1つ目と同じオーブは選択できません！");
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
			`>> [${orb1}] と [${orb2}] を分解し、強化の書 x1 を獲得しました！ (所持: 強化の書 x${this.stockScrolls})`,
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
			console.log(
				">> 合成可能なオーブ（同種の通常オーブ2個以上）が倉庫にありません！",
			);
			return;
		}

		console.log(
			"\n--- オーブの合成 (同種の通常オーブ2個 -> 上位オーブ1個) ---",
		);
		for (let i = 0; i < availableRecipes.length; i++) {
			const recipe = availableRecipes[i];
			if (recipe) {
				console.log(
					`${i + 1}: [${recipe.from}] (所持: ${recipe.count}個) -> [${recipe.to}] ${ORB_NAMES[recipe.to]} を合成`,
				);
			}
		}
		console.log("0: キャンセル");

		const select = await this.rl.question("合成するオーブの番号を選択: ");
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
		console.log(
			`>> [${targetRecipe.from}] を2個消費し、上位オーブ [${targetRecipe.to}] を合成しました！`,
		);
		this.saveGame();
	}

	private async dungeonPhase(dungeon: DungeonDef): Promise<void> {
		console.log(`\n>>> 【${dungeon.name}】に突入しました！ <<<`);
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

			console.log("\n==============================================");
			console.log(
				`   B${floor}F / B${dungeon.floors}F : ${enemy.name} が現れた！`,
			);
			console.log("==============================================");

			const survived = await this.battle(enemy);

			if (!survived) {
				console.log("\n[!] あなたは力尽きた...");
				console.log(
					`[!] 今回獲得したアイテム（書: ${runInv.scrolls}枚, オーブ: ${runInv.orbs.length}個）は失われました。`,
				);
				console.log("[!] 命からがら拠点へ運ばれました。（所持武器は無事です）");
				this.saveGame();
				return;
			}

			// ドロップ獲得
			const drops = template.getDrops();
			runInv.scrolls += drops.scrolls;
			runInv.orbs.push(...drops.orbs);

			console.log(`\n>> ${enemy.name} を討伐！`);
			if (drops.scrolls > 0)
				console.log(`   戦利品獲得: 強化の書 x${drops.scrolls}`);
			for (const o of drops.orbs) {
				console.log(`   戦利品獲得: オーブ [${o}] (${ORB_NAMES[o]})`);
			}

			if (floor === dungeon.floors) {
				console.log("\n**********************************************");
				console.log(
					`   【${dungeon.name}】完全踏破！おめでとうございます！   `,
				);
				console.log("**********************************************");
				this.stockScrolls += runInv.scrolls;
				this.stockOrbs.push(...runInv.orbs);
				this.saveGame();
				return;
			}

			console.log(`\n現在HP: ${this.player.hp}/${this.player.maxHp}`);
			console.log(
				`現在の未確定戦利品: 書 x${runInv.scrolls}, オーブ x${runInv.orbs.length}`,
			);
			console.log("1: 次の階層へ進む");
			console.log("2: 撤退する (戦利品を持ち帰って拠点に戻る)");

			const nextAction = await this.rl.question("行動を選択: ");
			if (nextAction.trim() === "2") {
				console.log("\n>> 慎重に撤退を選択しました。");
				this.stockScrolls += runInv.scrolls;
				this.stockOrbs.push(...runInv.orbs);
				console.log(
					`>> 戦利品（書: ${runInv.scrolls}枚, オーブ: ${runInv.orbs.length}個）を持ち帰りました！`,
				);
				this.saveGame();
				return;
			}
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
				`\n[YOU] HP: ${this.player.hp}/${this.player.maxHp}  vs  [${enemy.name}] HP: ${enemy.hp}/${enemy.maxHp}${
					enemy.poison > 0 ? ` (毒:${enemy.poison})` : ""
				}`,
			);
			console.log("1: 攻撃する");
			console.log("2: 撤退する (戦闘から逃げて拠点へ)");

			const act = await this.rl.question("コマンド: ");
			if (act.trim() === "2") {
				console.log(
					"\n>> 戦闘から離脱し、命からがら帰還しました。（戦利品は持ち帰れません）",
				);
				return false;
			}

			// プレイヤー攻撃
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
				console.log(
					`>> あなたの${hits > 1 ? `${i}撃目の` : ""}攻撃！ ${
						isCrit ? "【会心の一撃！】 " : ""
					}${enemy.name}に ${dmg} ダメージ！`,
				);

				if (vampRate > 0) {
					const heal = Math.max(1, Math.floor(dmg * vampRate));
					this.player.hp = Math.min(this.player.maxHp, this.player.hp + heal);
					console.log(
						`   [吸血] HPが ${heal} 回復した！ (現在HP: ${this.player.hp})`,
					);
				}

				if (poisonAdd > 0) {
					enemy.poison += poisonAdd;
					console.log(
						`   [猛毒] ${enemy.name}に毒を付与！ (+${poisonAdd} / 毒カウント: ${enemy.poison})`,
					);
				}
			}

			if (enemy.hp <= 0) {
				return true;
			}

			// ターン終了時：毒ダメージ
			if (enemy.poison > 0) {
				const poisonDmg = enemy.poison * 3;
				enemy.hp -= poisonDmg;
				console.log(
					`>> [毒効果] ${enemy.name}は毒で ${poisonDmg} のダメージを受けた！`,
				);
				if (enemy.hp <= 0) return true;
			}

			// 敵の攻撃
			console.log(
				`>> ${enemy.name}の反撃！ あなたは ${enemy.atk} のダメージを受けた！`,
			);
			this.player.hp -= enemy.atk;

			if (this.player.hp <= 0) {
				return false;
			}
		}

		return this.player.hp > 0;
	}
}

new Game().start().catch(console.error);
