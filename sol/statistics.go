package sol

import (
	"fmt"

	"oddstream.games/gomps5/util"
)

// Statistics is a container for the statistics for all variants
type Statistics struct {
	// PascalCase for JSON
	StatsMap map[string]*VariantStatistics
}

// VariantStatistics holds the statistics for one variant
type VariantStatistics struct {
	// PascalCase for JSON
	Won, Lost, CurrStreak, BestStreak, WorstStreak int   `json:",omitempty"`
	Percents                                       []int `json:",omitempty"`
	// Won is number of games with 100%
	// Lost is number of games with % less than 100
	// Won + Lost is total number of games played (won or abandoned)
	// []Percents is a record of games where % < 100
	// average % is (sum of Percents) + (100 * Won) / (Won+Lost)
}

func (stats *VariantStatistics) averagePercent() int {
	played := stats.Won + stats.Lost
	if played == 0 {
		return 0
	}
	var av int
	for _, percent := range stats.Percents {
		av += percent
	}
	av += stats.Won * 100
	return av / played
}

func (stats *VariantStatistics) bestPercent() int {
	var best = 0
	for _, percent := range stats.Percents {
		if percent > best {
			best = percent
		}
	}
	return best
}

func (stats *VariantStatistics) endOfGameToasts() {
	toasts := []string{}
	toasts = append(toasts,
		fmt.Sprintf("You have won %s, and lost %s (%d%%)",
			util.Pluralize("game", stats.Won),
			util.Pluralize("game", stats.Lost),
			((stats.Won*100)/(stats.Won+stats.Lost))))
	toasts = append(toasts, fmt.Sprintf("Your average score is %d%%", stats.averagePercent()))

	if stats.CurrStreak > 0 {
		toasts = append(toasts, fmt.Sprintf("You are on a winning streak of %s", util.Pluralize("game", stats.CurrStreak)))
	}
	if stats.CurrStreak < 0 {
		toasts = append(toasts, fmt.Sprintf("You are on a losing streak of %s", util.Pluralize("game", util.Abs(stats.CurrStreak))))
	}

	for _, t := range toasts {
		TheUI.Toast(t)
	}
}

// NewStatistics creates a new Statistics object
func NewStatistics() *Statistics {
	s := &Statistics{StatsMap: make(map[string]*VariantStatistics)}
	s.Load()
	return s
}

func (s *Statistics) findVariant() *VariantStatistics {
	stats, ok := s.StatsMap[ThePreferences.Variant]
	if !ok {
		stats = &VariantStatistics{} // everything 0
		s.StatsMap[ThePreferences.Variant] = stats
		println("statistics has encountered a new variant", ThePreferences.Variant)
	}
	return stats
}

func (s *Statistics) RecordWonGame() {

	TheUI.Toast(fmt.Sprintf("Recording completed game of %s", ThePreferences.Variant))

	stats := s.findVariant()

	stats.Won = stats.Won + 1

	if stats.CurrStreak < 0 {
		stats.CurrStreak = 1
	} else {
		stats.CurrStreak = stats.CurrStreak + 1
	}
	if stats.CurrStreak > stats.BestStreak {
		stats.BestStreak = stats.CurrStreak
	}

	stats.endOfGameToasts()

	s.Save()
}

func (s *Statistics) RecordLostGame() {

	TheUI.Toast(fmt.Sprintf("Recording lost game of %s", ThePreferences.Variant))

	percent := TheBaize.PercentComplete()
	if percent == 100 {
		println("*** That's odd, here is a lost game that is 100% complete ***")
	}

	stats := s.findVariant()

	stats.Lost = stats.Lost + 1
	// don't see that currStreak can ever be zero
	if stats.CurrStreak > 0 {
		stats.CurrStreak = -1
	} else {
		stats.CurrStreak = stats.CurrStreak - 1
	}
	if stats.CurrStreak < stats.WorstStreak {
		stats.WorstStreak = stats.CurrStreak
	}

	stats.Percents = append(stats.Percents, percent)

	s.Save()
}

func (s *Statistics) WelcomeToast() {
	toasts := []string{}

	stats, ok := s.StatsMap[ThePreferences.Variant]
	if !ok || stats.Won+stats.Lost == 0 {
		toasts = append(toasts, fmt.Sprintf("You have not played %s before", ThePreferences.Variant))
	} else {
		avpc := stats.averagePercent()
		bpc := stats.bestPercent()

		if stats.Won == 0 {
			toasts = append(toasts, fmt.Sprintf("You have yet to win a game of %s in %s", ThePreferences.Variant, util.Pluralize("attempt", stats.Lost)))
			if bpc > 0 && bpc != avpc {
				toasts = append(toasts, fmt.Sprintf("Your best score is %d%%, your average score is %d%%", bpc, avpc))
			}
		} else {
			toasts = append(toasts,
				fmt.Sprintf("You have won %s, and lost %s (%d%%)",
					util.Pluralize("game", stats.Won),
					util.Pluralize("game", stats.Lost),
					((stats.Won*100)/(stats.Won+stats.Lost))))

			toasts = append(toasts, fmt.Sprintf("Your average score is %d%%", avpc))

			if stats.CurrStreak > 0 {
				toasts = append(toasts, fmt.Sprintf("You are on a winning streak of %s", util.Pluralize("game", stats.CurrStreak)))
			}
			if stats.CurrStreak < 0 {
				toasts = append(toasts, fmt.Sprintf("You are on a losing streak of %s", util.Pluralize("game", util.Abs(stats.CurrStreak))))
			}
		}
	}

	for _, t := range toasts {
		TheUI.Toast(t)
	}
}
