# Rogue — Go + tcell

[Русский](README.md) | [English](README.en.md)

Rogue is a turn-based console game about exploring randomly generated dungeons.
Reach the exit on each of 21 floors while fighting enemies and collecting items.
The game is written in Go. The terminal interface uses
[`github.com/gdamore/tcell/v2`](https://github.com/gdamore/tcell/tree/v2);
saves use the standard `encoding/json` library. No C, cgo, or ncurses is required.

## Running the game

You need Go 1.25 or newer and an interactive terminal, such as Windows Terminal.
Run these commands from the directory containing this README:

```powershell
go mod download
go run ./rogue/cmd
```

Build and run on Windows:

```powershell
go build -o bin/rogue.exe ./rogue/cmd
./bin/rogue.exe
```

Use `go run ./rogue/cmd -seed 42` for reproducible dungeon generation.
The `-save-dir` flag sets the save directory. The default save file is
`saves/rogue.json`, relative to the working directory. Use the same working
directory between launches to find the same save.
The minimum terminal size is 60×20; the map scrolls with the player.
At 80×41, the entire map and interface fit on screen.

## Controls

| Key | Action |
| --- | --- |
| N / C | Start a new game / continue from the main menu |
| Enter | Continue a game, or start a new one if there is no save |
| WASD or arrow keys | Move; moving into an enemy attacks it |
| H | Choose a weapon with 1–9; use 0 to store the equipped weapon |
| J / K / E | Choose food / an elixir / a scroll, then press 1–9 |
| Space or period | Wait one turn |
| R | View the leaderboard |
| Esc | Return from the backpack to the game; resume an active game from the main menu |
| M | Open the main menu during a game |
| Q or Ctrl+C | Save and quit |

Russian keyboard layouts work on the same physical keys: N/Т, C/С,
WASD/ЦФЫВ, H/Р, J/О, K/Л, E/У, M/Ь, R/К, Q/Й.

If the window opens but seems unresponsive, press Enter or N/Т.
Movement becomes available after starting or continuing a game from the main menu.
If you see `Resize terminal`, enlarge the window to at least 60×20.
On Windows, the game uses tcell's native console backend. Launch it from
Windows Terminal or double-click `bin/rogue.exe`.

The player is `@`, stairs are `>`, a door is `+`, and a door with its key found is `/`.
Keys (`k`) and locked doors are red, blue, or yellow.
Items: `%` food, `!` elixir, `?` scroll, `)` weapon.
Enemies: `z` zombie, `v` vampire, `g` ghost, `O` ogre, `s` serpent mage,
and `m` revealed mimic. Before combat, a mimic looks like food.

### Weapons

Press H. The `In hand` row shows your equipped weapon and damage range;
weapons in the backpack are listed separately below it. Pressing 1–9 **equips**
that item: it moves from the backpack to your hand and is used automatically
when you attack by moving into an enemy. It therefore disappears from the
backpack list but remains visible under `In hand` and on the game status bar.
Equipping another weapon drops the previous one onto an adjacent tile.
Walk onto that tile to pick it up again.

Pressing 0 stores the equipped weapon **in the backpack**, keeping it.
Esc closes the backpack without changing equipment or spending a turn;
pressing Esc again during gameplay does nothing. M opens the main menu.
Equipping and unequipping a weapon each spend one turn.

## Features

- 21 floors, each with nine random rooms in a 3×3 grid connected using DFS.
- BFS validation of room, key, and exit reachability with lock states included.
- Turn-based combat, six enemy types with distinct abilities, and BFS pursuit.
- Up to nine backpack items per type, consumables, and weapon switching.
- Temporary elixirs, permanent scrolls, and maximum-health upgrades.
- Fog of war, Bresenham visibility rays, remembered room walls and corridors.
- Difficulty adjustments on subsequent floors based on damage taken and player health.
- Full JSON saves after every action and floor transition: map, actors, inventory,
  effects, keys, statistics, and random generator state.
- Victory and defeat records sorted by treasure. Finished games cannot be resumed;
  recording a result and clearing the active save happen together.

Treasure is awarded only for defeating enemies. Menus, terminal resizing,
invalid item selections, and walking into walls do not spend turns.
N in the main menu replaces the current save with a new game.

## Architecture

- `rogue/domain`: entities, generation, combat, AI, visibility, and repository interface.
  Does not import tcell, presentation, or datalayer.
- `rogue/presentation`: tcell screen, keyboard mapping, and frame rendering.
- `rogue/presenter`: game loop, menus, and state-to-frame mapping.
- `rogue/dataLayer`: JSON persistence, state validation, and atomic file replacement.
- `rogue/cmd`: application entry point and dependency wiring.

## Checks

```powershell
go test ./...
go vet ./...
```

Tests cover generation across floors and seeds, completing all 21 floors without
enemies, inventory limits, combat abilities, elixir expiration, exact save and
random sequence restoration, records, preserving corrupt save files, rendering
through `tcell.SimulationScreen`, keyboard layouts, and menu navigation.
