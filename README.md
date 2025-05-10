# League Rank

League ranker calculates the rankings of teams.

# Requirements

POSIX environment and Go 1.24.

# How to run

Input is taken from `STDIN`

`make run < input.txt` where `input.txt` is for example:

```
Lions 3, Snakes 3
Tarantulas 1, FC Awesome 0
Lions 1, FC Awesome 1
Tarantulas 3, Snakes 1
Lions 4, Grouches 0
```

Expected output: 

```
1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts
```