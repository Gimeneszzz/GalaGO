package entities

import (
	"testing"

	"galago/internal/config"
)

func TestMeteorStatsBySize(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		height         int
		wantHealth     int
		wantReward     int
		wantEscapeCost int
	}{
		{"pequeno", 28, 28, 1, config.SmallMeteorReward, config.SmallMeteorEscapePenalty},
		{"medio", 45, 40, 2, config.MediumMeteorReward, config.MediumMeteorEscapePenalty},
		{"grande", 120, 98, 3, config.LargeMeteorReward, config.LargeMeteorEscapePenalty},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			health, reward, escapeCost := meteorStats(test.width, test.height)
			if health != test.wantHealth || reward != test.wantReward || escapeCost != test.wantEscapeCost {
				t.Fatalf(
					"meteorStats(%d, %d) = (%d, %d, %d), esperado (%d, %d, %d)",
					test.width,
					test.height,
					health,
					reward,
					escapeCost,
					test.wantHealth,
					test.wantReward,
					test.wantEscapeCost,
				)
			}
		})
	}
}
