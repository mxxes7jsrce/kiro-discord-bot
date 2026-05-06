package bot

import (
	"strings"
	"testing"
)

func TestFlipCoin_ValidSide(t *testing.T) {
	for i := 0; i < 50; i++ {
		result := FlipCoin()
		if result.Side != "Heads" && result.Side != "Tails" {
			t.Errorf("unexpected side: %s", result.Side)
		}
	}
}

func TestFlipCoin_ValidEmoji(t *testing.T) {
	for i := 0; i < 50; i++ {
		result := FlipCoin()
		if result.Emoji != "🪙" && result.Emoji != "🌕" {
			t.Errorf("unexpected emoji: %s", result.Emoji)
		}
	}
}

func TestFlipCoin_EmojiMatchesSide(t *testing.T) {
	for i := 0; i < 50; i++ {
		result := FlipCoin()
		if result.Side == "Heads" && result.Emoji != "🪙" {
			t.Errorf("Heads should have 🪙 emoji, got %s", result.Emoji)
		}
		if result.Side == "Tails" && result.Emoji != "🌕" {
			t.Errorf("Tails should have 🌕 emoji, got %s", result.Emoji)
		}
	}
}

func TestIsCoinFlipCommand(t *testing.T) {
	valid := []string{"!coinflip", "!flip", "!coin"}
	for _, cmd := range valid {
		if !isCoinFlipCommand(cmd) {
			t.Errorf("expected %q to be a coinflip command", cmd)
		}
	}

	invalid := []string{"!roll", "coinflip", "!COINFLIP", "!flip extra", ""}
	for _, cmd := range invalid {
		if isCoinFlipCommand(cmd) {
			t.Errorf("expected %q to not be a coinflip command", cmd)
		}
	}
}

func TestHandleCoinFlip_ContainsSideAndEmoji(t *testing.T) {
	for i := 0; i < 20; i++ {
		resp := handleCoinFlip()
		hasHeads := strings.Contains(resp, "Heads")
		hasTails := strings.Contains(resp, "Tails")
		if !hasHeads && !hasTails {
			t.Errorf("response missing side: %s", resp)
		}
	}
}

func TestDispatchCoinFlipCommand_Valid(t *testing.T) {
	resp, handled := dispatchCoinFlipCommand("!coinflip")
	if !handled {
		t.Error("expected command to be handled")
	}
	if resp == "" {
		t.Error("expected non-empty response")
	}
}

func TestDispatchCoinFlipCommand_Invalid(t *testing.T) {
	_, handled := dispatchCoinFlipCommand("!roll 2d6")
	if handled {
		t.Error("expected command to not be handled")
	}
}
