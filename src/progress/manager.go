package progress

import (
	"fmt"
	"strings"
	"time"
)

type Manager[T int | uint | float32 | float64] struct {
	size  int
	chars string
	add   chan Bar[T]
	done  chan struct{}
}

func New[T int | uint | float32 | float64](size int, chars string) Manager[T] {
	p := Manager[T]{
		size:  size,
		chars: chars,
		add:   make(chan Bar[T]),
		done:  make(chan struct{}),
	}
	go BarCtrl(p)
	return p
}

func (p Manager[T]) Close() {
	close(p.done)
}

func _IsCharsValid(chars string) error {
	if len(chars) == 0 {
		return fmt.Errorf("chars is empty")
	}
	if len(chars) > 3 {
		return fmt.Errorf("chars is too long, more than 3 characters (%d)", len(chars))
	}
	return nil
}

func FormatBar[T int | uint | float32 | float64](size uint, chars string, current T, total T) (string, error) {
	if err := _IsCharsValid(chars); err != nil {
		return "", err
	}

	if float64(total) == 0 {
		return "", fmt.Errorf("total must not be zero")
	}

	coef := float32(float64(current) / float64(total))
	if coef > 1 {
		coef = 1
	} else if coef < 0 {
		coef = 0
	}

	filledLen := int(float32(size) * coef)
	if filledLen > int(size) {
		filledLen = int(size)
	}

	filled := strings.Repeat(string(chars[0]), filledLen)

	remaining := int(size) - filledLen
	emptyChar := chars[len(chars)-1]

	if remaining > 0 {
		if len(chars) > 2 {
			// Place le caractère intermédiaire en premier dans la partie vide
			filled += string(chars[1])
			remaining--
		}
		filled += strings.Repeat(string(emptyChar), remaining)
	}

	return filled, nil
}

func BarCtrl[T int | uint | float32 | float64](p Manager[T]) error {
	var Bars []Bar[T]
	var values []T
	var locked []bool
	prevCount := 0

	for {
		select {
		case <-p.done:
			fmt.Fprintln(output, "done")
			return nil // arrêt propre
		default:
			// suite du traitement normal
		}
		// Récupération des nouveaux transmetteurs
		select {
		case v := <-p.add:
			Bars = append(Bars, v)
			var zero T
			values = append(values, zero)
			locked = append(locked, false)
		default:
			// pas d'ajout
		}

		var bars []string
		for i := range Bars {
			if !locked[i] {
				select {
				case val := <-Bars[i].current:
					values[i] = val
					// Vérifie si la barre est complète
					if float64(val) >= float64(Bars[i].total) {
						locked[i] = true
					}
				default:
					// pas de nouvelle valeur, on garde l'ancienne
				}
			}

			barStr, err := FormatBar(uint(p.size), p.chars, values[i], Bars[i].total)
			if err != nil {
				return err
			}
			if values[i] == Bars[i].total {
				bars = append(bars, fmt.Sprintf("<%s> [%s] done", Bars[i].title, barStr))
			} else {
				bars = append(bars, fmt.Sprintf("<%s> [%s] %v/%v", Bars[i].title, barStr, values[i], Bars[i].total))

			}
		}

		// Effacer les anciennes lignes
		for i := 0; i < prevCount; i++ {
			fmt.Fprint(output, "\x1b[1A")
			fmt.Fprint(output, "\x1b[2K")
		}

		// Afficher les barres
		for _, bar := range bars {
			fmt.Fprintln(output, bar)
		}

		prevCount = len(bars)
		time.Sleep(50 * time.Millisecond)
	}
}

func (p *Manager[T]) Add(channel chan T, total T, title string) {
	var Bar Bar[T] = Bar[T]{
		current: channel,
		total:   total,
		title:   title,
	}
	p.add <- Bar
}
