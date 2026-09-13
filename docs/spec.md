# Minimal Rogue-lite Specification (v0.1)

## 1. 概要
「手塩にかけて育てた装備で壁を越える」ループを検証するための、画面描画を省いたCUIベースの最小ローグライト。

## 2. コアループ
1. **拠点**: 武器の強化（ステータス底上げ）および印の着脱（シナジー構築）。
2. **出撃**: 全5階層のダンジョンへ挑む。
3. **選択**: 各階層突破時、または戦闘中に「前進」か「撤退（戦利品持ち帰り）」を選択。
4. **結果**:
   - 撤退または踏破: 戦利品（強化の書・印）を拠点に持ち帰り。
   - 死亡: その回の戦利品は全ロスト。ただし育てた武器は永続保持。

## 3. データ構造

```typescript
type OrbType = 'MULTI_HIT' | 'CRITICAL' | 'VAMP' | 'POISON';

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
```

## 4. 印（オーブ）の定義
- `MULTI_HIT` (連撃): 2回攻撃を行う。1撃あたりのダメージは通常の65%。
- `CRITICAL` (会心): 各撃ごとに25%の確率でダメージが2倍。
- `VAMP` (吸血): 与えた実ダメージの15%分、自身のHPを回復（最大HP上限）。
- `POISON` (猛毒): 攻撃ヒット時、敵に毒カウンター+1。ターン終了時に [毒カウンター × 3] の固定ダメージ。

## 5. 戦闘計算ルール
プレイヤー先攻のターン制。
1. **プレイヤー攻撃フェーズ**:
   - 攻撃力 = `baseAtk + (plus * 2)`
   - 攻撃回数 = `MULTI_HIT` あり ? 2回 (倍率0.65) : 1回 (倍率1.0)
   - 各撃ごとに判定:
     - 会心判定 (`CRITICAL` かつ 乱数 < 0.25) -> ダメージ2倍
     - 敵HP減算
     - 吸血判定 (`VAMP`) -> 与ダメの15%回復
     - 毒付与 (`POISON`) -> 毒+1
     - 敵HP <= 0 なら即座に勝利
2. **ターン終了フェーズ**:
   - 敵の毒ダメージ処理 (`poison * 3`)
   - 敵HP <= 0 なら勝利
3. **敵反撃フェーズ**:
   - プレイヤーHP -= 敵攻撃力
   - プレイヤーHP <= 0 なら死亡処理

## 6. エネミー & ドロップテーブル
- B1F: スライム (HP 15, ATK 3) -> 強化の書 ×1 (100%)
- B2F: コボルト (HP 28, ATK 6) -> 強化の書 ×1 (100%)
- B3F: オーク (HP 45, ATK 10) -> オーブ ×1 (4種から均等抽選)
- B4F: ゴーレム (HP 75, ATK 14) -> 強化の書 ×2 (100%)
- B5F: レッドドラゴン (BOSS) (HP 130, ATK 20) -> 強化の書 ×3 + オーブ ×1
