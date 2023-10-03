# morningjuegos

MorningJuegos is a Discord bot taken too far. What started as a Discord channel for sharing results from games like NYT Connections and Coffee Golf became a full blown tournament.

## What it does
MorningJuegos will watch Discord and try to match any messages shared from Coffee Golf. It will then pull all of the information from the message and store it. Each message is connected to a 10 day tournament. The first shared message of the day from each user is counted torwards their total in the tournament. During the tournament, the bot can be queried to see the overall leaderboard.

### Example Share Message
```
Coffee Golf - Oct 3
10 Strokes - Top 1% 🏆

🟦🟥🟪🟨🟩
1️⃣2️⃣2️⃣2️⃣3️⃣
```

### Example Leaderboard
```
Current Tournament: Sep 30, 2023 - Oct 9, 2023

Leaders
1: @player1 - 35 Total Strokes 🥇 ⬆️ 1 👑
2: @player2 - 36 Total Strokes 🥈 ⬇️ 
3: @player3 - 37 Total Strokes 🥉 ⬆️ 
4: @player4 - 39 Total Strokes  ⬆️ 
5: @player5 - 39 Total Strokes  ⬇️ 
```

### Example Stats Page
```
Current Tournament: Sep 30, 2023 - Oct 9, 2023

Most hole in ones: @player2 with 1 hole in ones
Worst round of the tournament: @player3, 17 strokes 🤡
Most common opening hole: 🟩
Most common finishing hole: 🟪
Hardest hole: 🟨 with an average of 3.07 strokes
[All Time] Most consistent players: @player1 @player5 with a standard deviation of 0.577 strokes
[All Time] Least consistent players: @player3 with a standard deviation of 1.732 strokes
```

## Architecture
The Discord bot uses Postgres to store all data. Previously this was done using SQLite and [Bun](https://bun.uptrace.dev/). Schemas/migrations are managed with [dbmate](https://github.com/amacneil/dbmate). [sqlc](https://sqlc.dev/) is used to turn raw SQL queries in to Go code. Optionally, the leaderboards and stats pages can be cached via Redis. [NATS](https://nats.io/) is used for event driven actions, including adding missing rounds for new players or updating final tournament placements. Everything is hosted on DigitalOcean using their managed Kubernetes and database services. Docker images are built by CircleCI and deployed using FluxCD.