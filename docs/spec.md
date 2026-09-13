# Minimal Rogue-lite Specification (v0.1)

## 1. 概要
「手塩にかけて育てた装備で壁を越える」ループを検証するための、画面描画を省いたCUIベースの最小ローグライト。

## 2. コアループ
1. **拠点**:
   - 武器の強化（ステータス底上げ）
   - 印の着脱（シナジー構築）
   - オーブの分解（余剰オーブ2個 -> 強化の書1枚）
   - オーブの合成（同種通常オーブ2個 -> 上位オーブ1個）
   - 体力の強化（強化の書2枚 -> 最大HP+10）
2. **出撃**:
   - 始まりの洞窟（全5階層 / 初級）
   - 灼熱の深層（全10階層 / 上級）
   - 無限の深淵（階層無制限 / エンドレス）
3. **選択**: 各階層突破時、または戦闘中に「前進」か「撤退（戦利品持ち帰り）」を選択。
4. **結果**:
   - 撤退または踏破: 戦利品（強化の書・印）を拠点に持ち帰り。
   - 死亡: その回の戦利品は全ロスト。ただし育てた武器や最大HP、最高到達階層は永続保持。

## 3. データ構造

```typescript
type NormalOrbType = 'MULTI_HIT' | 'CRITICAL' | 'VAMP' | 'POISON';
type PlusOrbType = 'MULTI_HIT_PLUS' | 'CRITICAL_PLUS' | 'VAMP_PLUS' | 'POISON_PLUS';
type OrbType = NormalOrbType | PlusOrbType;

interface Weapon {
  name: string;
  baseAtk: number; // 初期値: 10
  plus: number;    // 強化値 (初期値: 0, 1ごとに攻撃力+2)
  slots: OrbType[];// 最大3枠
}

interface Player {
  maxHp: number;   // 50
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
  maxHp?: number;        // 初期値: 50
  deepestFloor?: number; // 無限の深淵の最高到達階層 (初期値: 0)
}
```

## 4. 印（オーブ）の定義
- `MULTI_HIT` (連撃): 2回攻撃を行う。1撃あたりのダメージは通常の65%。
- `CRITICAL` (会心): 各撃ごとに25%の確率でダメージが2倍。
- `VAMP` (吸血): 与えた実ダメージの15%分、自身のHPを回復（最大HP上限）。
- `POISON` (猛毒): 攻撃ヒット時、敵に毒カウンター+1。ターン終了時に [毒カウンター × 3] の固定ダメージ。

### 上位印（合成限定）
同種の通常オーブ2個を合成することで獲得可能。
- `MULTI_HIT_PLUS` (連撃+): 2回攻撃を行う。1撃あたりのダメージは通常の75%。
- `CRITICAL_PLUS` (会心+): 各撃ごとに40%の確率でダメージが2倍。
- `VAMP_PLUS` (吸血+): 与えた実ダメージの25%分、自身のHPを回復（最大HP上限）。
- `POISON_PLUS` (猛毒+): 攻撃ヒット時、敵に毒カウンター+2。

## 5. 戦闘計算ルール
プレイヤー先攻のターン制。
1. **プレイヤー攻撃フェーズ**:
   - 攻撃力 = `baseAtk + (plus * 2)`
   - 攻撃回数: `MULTI_HIT_PLUS` または `MULTI_HIT` あり ? 2回 : 1回
   - 威力倍率: `MULTI_HIT_PLUS` あり ? 0.75 : (`MULTI_HIT` あり ? 0.65 : 1.0)
   - 各撃ごとに判定:
     - 会心判定:
       - `CRITICAL_PLUS` あり -> 40%でダメージ2倍
       - `CRITICAL` あり（通常のみ） -> 25%でダメージ2倍
     - 敵HP減算
     - 吸血判定:
       - `VAMP_PLUS` あり -> 与ダメの25%回復
       - `VAMP` あり（通常のみ） -> 与ダメの15%回復
     - 毒付与:
       - `POISON_PLUS` あり -> 毒+2
       - `POISON` あり（通常のみ） -> 毒+1
     - 敵HP <= 0 なら即座に勝利
2. **ターン終了フェーズ**:
   - 敵の毒ダメージ処理 (`poison * 3`)
   - 敵HP <= 0 なら勝利
3. **敵反撃フェーズ**:
   - プレイヤーHP -= 敵攻撃力
   - プレイヤーHP <= 0 なら死亡処理

## 6. エネミー & ドロップテーブル
### 始まりの洞窟
- B1F: スライム (HP 15, ATK 3) -> 強化の書 ×1 (100%)
- B2F: コボルト (HP 28, ATK 6) -> 強化の書 ×1 (100%)
- B3F: オーク (HP 45, ATK 10) -> オーブ ×1 (4種から均等抽選)
- B4F: ゴーレム (HP 75, ATK 14) -> 強化の書 ×2 (100%)
- B5F: レッドドラゴン (BOSS) (HP 130, ATK 20) -> 強化の書 ×3 + オーブ ×1

## 7. エンドレスダンジョン「無限の深淵」
階層上限のない高難度チャレンジダンジョン。
- **階層**: B1F〜無限
- **敵生成**:
  - 名前: ランダムな接頭辞（凶暴な/深淵の/古代の/紅蓮の/漆黒の/彷徨える/狂気の/奈落の）＋敵名（スライム/コボルト/オーク/ゴーレム/ドラゴン/ワイバーン/デーモン/死霊騎士）
  - 5階層ごと（B5F, B10F...）は「中ボス」扱い
  - HP: `Math.floor(60 + Math.pow(floor, 1.4) * 12)`
  - ATK: `Math.floor(8 + floor * 2.5)`
- **ドロップ**:
  - 各階層突破時: 強化の書 `Math.min(10, Math.floor(1 + floor / 2))` 枚
  - 5階層ごと（中ボス）: 上位オーブ（〜_PLUS）×1 確定ドロップ
- **ハイスコア**:
  - 到達階層を `deepestFloor` としてセーブデータに自動記録。
