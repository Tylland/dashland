# Systems Overview

This document lists the game systems, a short description of what each does, and the recommended per-frame execution order.

## Systems (what they do)

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

1. `InputSystem`
2. `InputActionSystem`
3. `GravitySystem`
4. `WallWalkerSystem`
5. `Pusher`
6. `BlockCollisionSystem`
7. `Combat` and `Collector` (event handlers)
8. `BlockMovement`
9. `Renderer`

## Notes

- Movement intent producers (input, AI, gravity, pusher) must run before `BlockCollisionSystem`.
- Event handlers should run immediately after collision detection so removals/changes are applied before movement commit.
- `BlockMovement` must run after collision resolution so only validated moves are committed.

---

File generated from repository analysis on December 25, 2025.
