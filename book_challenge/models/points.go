package models

import (
	"fmt"
	"time"
)

type Points struct {
	Amount         int
	Streak         int
	LastDayReading string
	Multiplier     int
	Prizes         []string
}

func (p *Points) EarnPoints(pointsEarned int) {
	today := time.Now()
	if today.AddDate(0, 0, -1).Format("2006-01-02") == p.LastDayReading {
		p.Streak++
	} else {
		p.Streak = 1
	}
	p.LastDayReading = today.Format("2006-01-02")

	if p.Streak >= 56 {
		p.Multiplier = 2
	}
	if p.Streak > 0 {
		if (p.Streak % 7) == 0 {
			p.Amount = p.Amount + (p.Multiplier * 100) //finish a book on the seventh day of reading a page a day you get an exta 100 points
		}
		if (p.Streak % 28) == 0 {
			p.Amount = p.Amount + (p.Multiplier * 500)
		}
	}
	p.Amount = p.Amount + (p.Multiplier * pointsEarned)
}

func (p *Points) SpendPoints(key string) {
	switch key {
	case "Slurpee":
		p.Amount = p.Amount - 100
	case "Movie Night":
		p.Amount = p.Amount - 200
	case "Small Toy":
		p.Amount = p.Amount - 1000
	case "Switch Game":
		p.Amount = p.Amount - 5000
	}
	fmt.Println(p.Amount)
	p.Prizes = append(p.Prizes, key)
}
