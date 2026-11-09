# Systems Overview

This document lists the game systems, a short description of what each does, and the recommended per-frame execution order.

## Systems (what they do)

- **StageTimerSystem**: counts down the configured stage time limit, pauses gameplay at expiry, and triggers player death after a three-second grace period.
- **InputSystem**: reads keyboard state into InputComponent booleans.
- **InputActionSystem**: converts InputComponent into movement intent (`BlockStep.Move`) for entities (player).
- **GravitySystem**: applies gravity to entities with `CanFall`, creating downward/diagonal `BlockStep` moves when unsupported.
- **WallWalkerSystem**: AI for enemies (`WallWalkerComponent`); decides directions and sets enemy `BlockStep` moves.
- **Pusher**: handles player push interactions; when player moves into a pushable entity it sets that entity's `BlockStep`.
- **BlockCollisionSystem**: inspects each entity's intended `BlockStep` target, cancels blocked moves, and emits events (`blockcollision`, `collect`, `entitycollision`, `damage`).
- **Combat**: consumes collision/damage/block events, removes or damages entities, spawns diamonds, and triggers stage changes (game rules for collisions/damage).
- **Collector**: consumes `collect` events, removes collected entities, and plays collection sound.
- **BlockMovement**: attempts and commits moves (`Stage.TryMoveEntity`), interpolates `Position.Vector2` toward step targets, and cancels/commits steps.
- **Renderer**: updates stage rendering and draws entities (sprite frame updates and texture draws).

## Recommended Execution Order (per frame)

1. `StageTimerSystem`
2. `InputSystem`
3. `InputActionSystem`
4. `GravitySystem`
5. `WallWalkerSystem`
6. `Pusher`
7. `BlockCollisionSystem`
8. `Combat` and `Collector` (event handlers)
9. `BlockMovement`
10. `Renderer`

## Notes

- Movement intent producers (input, AI, gravity, pusher) must run before `BlockCollisionSystem`.
- Event handlers should run immediately after collision detection so removals/changes are applied before movement commit.
- `BlockMovement` must run after collision resolution so only validated moves are committed.
- Stages set `TimeLimit` in seconds; `0` disables the countdown. On expiry the world pauses except for the timer, gameplay-event handling, rendering, and cleanup.
- The timer kills the player through the existing damage flow after a three-second timeout grace period.

---

File generated from repository analysis on December 25, 2025.
