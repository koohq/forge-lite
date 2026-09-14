# Forge Lite

[English](README.md) | [日本語](README.ja.md)

A lightweight terminal CLI roguelite with zero external dependencies, powered entirely by the Go standard library. Features rapid stat inflation, strategic orb attachments, auto-combat, and customizable themes.

---

## Table of Contents

- [Overview](#overview)
- [Requirements](#requirements)
- [Quick Start](#quick-start)
  - [Run](#run)
  - [Build](#build)
  - [Cross-Compilation](#cross-compilation)
  - [Install via Go](#install-via-go)
- [Game Systems](#game-systems)
  - [Dungeons](#dungeons)
  - [Town Hub & QoL Features](#town-hub--qol-features)
  - [Dynamic Weapon Rank Progression](#dynamic-weapon-rank-progression)
- [Preset Themes](#preset-themes)
- [Custom Themes](#custom-themes)
- [License](#license)

---

## Overview

**Forge Lite** is a pure Go CLI roguelite focused on the satisfying loop of power progression: delving into perilous dungeons, harvesting upgrade scrolls and trait orbs, and reforging your weapon into a godlike instrument of destruction.

- **Zero External Dependencies**: Implemented using strictly Go standard libraries (`fmt`, `bufio`, `embed`, etc.). No CGo, no heavy third-party TUI frameworks.
- **Rapid Inflation Gameplay**: Weapon damage scales exponentially as plus values rise into the hundreds and thousands.
- **Hands-Free Auto Combat**: Clean, fast-paced turn resolution with automated safety stops when player HP falls to 30% or below.
- **Built-in & Custom Themes**: Switch effortlessly between distinct fantasy, sci-fi, and tactical themes, or create your own custom theme via JSON.

---

## Requirements

- **Go**: `1.27.1` or higher
- **Task** ([go-task](https://taskfile.dev)): Optional (used for simplified build automation)

---

## Quick Start

### Run

Launch the game directly from source:

```bash
# Using Task
task run

# Or using Go directly
go run ./cmd/game
```

### Build

Compile the binary for your current operating system into `bin/game` (or `bin/game.exe` on Windows):

```bash
# Using Task
task build

# Or using Go directly
mkdir -p bin
go build -o bin/game ./cmd/game
```

### Cross-Compilation

Build binaries for Windows, Linux, and macOS in one command:

```bash
task build:all
```

Outputs:
- `bin/game-windows-amd64.exe` (Windows amd64)
- `bin/game-linux-amd64` (Linux amd64)
- `bin/game-darwin-arm64` (macOS arm64)

### Install via Go

Install the latest release directly into your `$GOPATH/bin`:

```bash
go install github.com/koohq/forge-lite/cmd/game@latest
```

---

## Game Systems

### Dungeons

1. **Starter Cave (始まりの洞窟)**: 5 Floors. Designed for new runs to gather initial scrolls and establish basic orb loadouts.
2. **Scorching Depths (灼熱の深層)**: 10 Floors. Intermediate challenge with hazardous enemies and richer rewards.
3. **Infinite Abyss (無限の深淵)**: Endless floors. Features procedural enemy scaling and checkpoint skips every 10 floors based on your highest recorded floor.

> **Risk & Return**: Retiring or clearing a floor lets you carry all gathered scrolls and orbs safely back to town. Falling in battle loses that expedition's loot, but all weapon levels, maximum HP, and highest reached depths persist permanently.

### Town Hub & QoL Features

- **Auto-Combat**: Seamlessly battles dungeon enemies turn-by-turn. Automatically yields control back to you when HP drops to 30% or below, or when encountering lethal danger.
- **Bulk & Target Enhancements**:
  - **Single Enhance**: Spend 1 scroll to upgrade weapon (+1).
  - **Target / Count Enhance**: Specify an exact number of scrolls or target plus value.
  - **Max Enhance**: Dump all remaining scrolls into your weapon instantly.
- **Orb Management & Crafting**:
  - **Socketing**: Customize your weapon with up to 3 trait orbs (`MULTI_HIT`, `CRITICAL`, `VAMP`, `POISON`).
  - **Bulk Dismantling**: Break down all unequipped inventory orbs into upgrade scrolls in a single action.
  - **Higher-Tier Synthesis**: Combine 2 identical standard orbs into a potent `+` version (e.g., `CRITICAL` -> `CRITICAL_PLUS`).
- **HP Fortification**: Exchange upgrade scrolls to permanently boost your hero's maximum health.

### Dynamic Weapon Rank Progression

As your weapon's `+` enhancement value reaches specific milestones, it dynamically gains new titles and visual fanfare in town and during combat:

| Enhancement Milestone | Classic Fantasy | Cyberpunk | Partner Sync |
| :--- | :--- | :--- | :--- |
| **+0** | Bronze Sword | Pulse Blade | Tactical Android "Iris" |
| **+50** | Steel Sword | Pulse Blade Mk-II | Iris, Combat Companion |
| **+200** | Hero's Blade | High-Frequency Katana | Iris, In Sync |
| **+500** | Dragon Slayer | Plasma Saber | Iris, Trusted Partner |
| **+1000+** | Blade of Armageddon | Antimatter Void Edge | Iris, Veteran Partner |

---

## Preset Themes

Forge Lite comes bundled with three fully localized embedded themes (supporting English and Japanese):

- 🗡️ **Classic Fantasy**: A traditional fantasy adventure. Visit the village blacksmith, temper your blade, and delve deep into monster-infested dungeons.
- ⚡ **Cyberpunk**: High-tech cyberspace infiltration. Overclock your pulse blade, breach firewall subnets, and hack mainframe nodes.
- 🤖 **Partner Sync**: Sci-fi buddy tactical combat. Synchronize data links and conduct joint combat drills alongside your tactical android partner, "Iris".

Switch between preset themes at any time through the in-game **Settings** menu.

---

## Custom Themes

You can fully customize all displayed names, descriptions, dungeon environments, enemy descriptors, and hub titles using a single JSON file.

### How to use a Custom Theme

1. **Export Template**: Open the in-game **Settings** menu and select **Export Custom Theme Template**. This writes `custom_theme.json` to the current working directory.
   - Alternatively, copy the bundled template:
     ```bash
     cp custom_theme.example.json custom_theme.json
     ```
2. **Edit `custom_theme.json`**: Modify weapon titles, dungeon names, orb effects, and status strings to fit your desired narrative or language.
3. **Automatic Loading**: When `custom_theme.json` is present in the working directory, the game will automatically load it with the highest priority upon launch.

---

## License

This project is open source and available under the [MIT License](LICENSE).
