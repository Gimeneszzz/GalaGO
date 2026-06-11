package game

import "testing"

func TestCirclesOverlap(t *testing.T) {
	tests := []struct {
		name string
		ax   float64
		ay   float64
		ar   float64
		bx   float64
		by   float64
		br   float64
		want bool
	}{
		{"sobrepostos", 0, 0, 10, 12, 0, 5, true},
		{"tocando", 0, 0, 10, 15, 0, 5, true},
		{"separados", 0, 0, 10, 16, 0, 5, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := circlesOverlap(test.ax, test.ay, test.ar, test.bx, test.by, test.br)
			if got != test.want {
				t.Fatalf("circlesOverlap() = %v, esperado %v", got, test.want)
			}
		})
	}
}

func TestLifeLabel(t *testing.T) {
	if got := lifeLabel(1); got != "vida" {
		t.Fatalf("lifeLabel(1) = %q, esperado %q", got, "vida")
	}
	if got := lifeLabel(3); got != "vidas" {
		t.Fatalf("lifeLabel(3) = %q, esperado %q", got, "vidas")
	}
}
