import { stdin as input, stdout as output } from "node:process";
import * as readline from "node:readline/promises";

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

function getRandomOrb(): NormalOrbType {
	const orb = DROP_ORBS[Math.floor(Math.random() * DROP_ORBS.length)];
	return orb ?? "MULTI_HIT";
}

const FLOOR_ENEMIES: Array<{
	name: string;
	hp: number;
	atk: number;
	getDrops: () => RunInventory;
}> = [
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

class Game {
	private readonly rl: readline.Interface;
	private readonly player: Player;
	private stockScrolls = 0;
	private readonly stockOrbs: OrbType[] = [];
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

	async start(): Promise<void> {
		console.log("==============================================");
		console.log("   Minimal Rogue-lite Prototype (CUI Ver)   ");
		console.log("==============================================");

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
		console.log("1: ダンジョンへ出撃 (全5階)");
		console.log("2: 武器を鍛える (強化の書を1枚消費)");
		console.log("3: オーブを武器に装着");
		console.log("4: オーブを分解 (任意のオーブ2個 -> 強化の書1枚)");
		console.log("5: オーブを合成 (同種オーブ2個 -> 上位オーブ)");
		console.log("0: ゲーム終了");

		const choice = await this.rl.question("\n行動を選択してください: ");

		switch (choice.trim()) {
			case "1":
				await this.dungeonPhase();
				break;
			case "2":
				this.upgradeWeapon();
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

	private upgradeWeapon(): void {
		if (this.stockScrolls <= 0) {
			console.log(">> 強化の書がありません！");
			return;
		}
		this.stockScrolls--;
		this.player.weapon.plus++;
		console.log(
			`>> 武器を鍛えました！ ${this.player.weapon.name}+${this.player.weapon.plus} (攻撃力: ${this.getWeaponAtk()})`,
		);
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
	}

	private async dungeonPhase(): Promise<void> {
		console.log("\n>>> ダンジョンに突入しました！ <<<");
		const runInv: RunInventory = { scrolls: 0, orbs: [] };

		for (let floor = 1; floor <= 5; floor++) {
			const template = FLOOR_ENEMIES[floor - 1];
			if (!template) {
				break;
			}

			const enemy: Enemy = {
				name: template.name,
				hp: template.hp,
				maxHp: template.hp,
				atk: template.atk,
				poison: 0,
			};

			console.log("\n==============================================");
			console.log(`   B${floor}F : ${enemy.name} が現れた！`);
			console.log("==============================================");

			const survived = await this.battle(enemy);

			if (!survived) {
				console.log("\n[!] あなたは力尽きた...");
				console.log(
					`[!] 今回獲得したアイテム（書: ${runInv.scrolls}枚, オーブ: ${runInv.orbs.length}個）は失われました。`,
				);
				console.log("[!] 命からがら拠点へ運ばれました。（所持武器は無事です）");
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
				console.log(`   戦利品獲得: オーブ [${o}]`);
			}

			if (floor === 5) {
				console.log("\n**********************************************");
				console.log("   ダンジョン完全踏破！おめでとうございます！   ");
				console.log("**********************************************");
				this.stockScrolls += runInv.scrolls;
				this.stockOrbs.push(...runInv.orbs);
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
