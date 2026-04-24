package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// DiceResult holds the result of a dice roll
type DiceResult struct {
	Dice  int
	Sides int
	Rolls []int
	Total int
}

// RollDice rolls `dice` number of `sides`-sided dice
func RollDice(dice, sides int) (*DiceResult, error) {
	if dice < 1 || dice > 20 {
		return nil, fmt.Errorf("number of dice must be between 1 and 20")
	}
	if sides < 2 || sides > 100 {
		return nil, fmt.Errorf("number of sides must be between 2 and 100")
	}

	rolls := make([]int, dice)
	total := 0
	for i := 0; i < dice; i++ {
		r := rand.Intn(sides) + 1
		rolls[i] = r
		total += r
	}
	return &DiceResult{Dice: dice, Sides: sides, Rolls: rolls, Total: total}, nil
}

// ParseDiceNotation parses a string like "2d6" into dice and sides
func ParseDiceNotation(notation string) (int, int, error) {
	parts := strings.SplitN(strings.ToLower(notation), "d", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid dice notation: use NdM format (e.g. 2d6)")
	}

	dice, err := strconv.Atoi(parts[0])
	if err != nil || parts[0] == "" {
		dice = 1
		if parts[0] != "" {
			return 0, 0, fmt.Errorf("invalid number of dice: %s", parts[0])
		}
	}

	sides, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid number of sides: %s", parts[1])
	}

	return dice, sides, nil
}

// FormatDiceResult formats a DiceResult into a human-readable string
func FormatDiceResult(r *DiceResult) string {
	if r.Dice == 1 {
		return fmt.Sprintf("🎲 Rolled **d%d**: **%d**", r.Sides, r.Total)
	}
	rollStrs := make([]string, len(r.Rolls))
	for i, v := range r.Rolls {
		rollStrs[i] = strconv.Itoa(v)
	}
	return fmt.Sprintf("🎲 Rolled **%dd%d**: [%s] = **%d**", r.Dice, r.Sides, strings.Join(rollStrs, ", "), r.Total)
}
