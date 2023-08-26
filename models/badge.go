package models

import (
	"time"

	"github.com/guregu/null"
)

// Badge represents a badge model.
type Badge struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	Secret      bool   `json:"secret"`
	Filename    string `json:"filename"`
}

// PlayerBadge represents a Player2Badge model.
type PlayerBadge struct {
	Badge     *Badge    `json:"badge"`
	PlayerID  int       `json:"player_id"`
	Level     null.Int  `json:"level,omitempty"`
	LegID     null.Int  `json:"leg_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

var LegBadges = []LegBadge{
	BadgeDoubleDouble{ID: 8},
	BadgeTripleDouble{ID: 9},
	BadgeMadHouse{ID: 10},
	BadgeMerryChristmas{ID: 11},
	BadgeBigFish{ID: 14},
	BadgeGettingCrowded{ID: 17},
	BadgeBullseye{ID: 20},
}

type LegBadge interface {
	GetID() int
	Validate(*Leg) (bool, *int)
}

type BadgeMerryChristmas struct{ ID int }
type BadgeGettingCrowded struct{ ID int }
type BadgeBullseye struct{ ID int }
type BadgeDoubleDouble struct{ ID int }
type BadgeTripleDouble struct{ ID int }
type BadgeMadHouse struct{ ID int }
type BadgeBigFish struct{ ID int }

type BadgeHighScore struct{ ID int }
type BadgeHigherScore struct{ ID int }
type BadgeTheMaximum struct{ ID int }
type BadgeGlobetrotter struct{ ID int }
type BadgePartyOfTwo struct{ ID int }
type BadgeWorkFromHome struct{ ID int }
type BadgeHardlyWorking struct{ ID int }
type BadgeWeeklyPlayer struct{ ID int }

func (b BadgeGettingCrowded) GetID() int {
	return b.ID
}
func (b BadgeGettingCrowded) Validate(leg *Leg) (bool, *int) {
	return len(leg.Players) > 4, nil
}
func (b BadgeMerryChristmas) GetID() int {
	return b.ID
}
func (b BadgeMerryChristmas) Validate(leg *Leg) (bool, *int) {
	d := leg.Endtime.Time
	return d.Day() == 25 && d.Month() == 12, nil
}
func (b BadgeBullseye) GetID() int {
	return b.ID
}
func (b BadgeBullseye) Validate(leg *Leg) (bool, *int) {
	visit := leg.GetLastVisit()
	last := visit.GetLastDart()
	return last.ValueRaw() == BULLSEYE && last.Multiplier == DOUBLE, &visit.PlayerID
}

func (b BadgeDoubleDouble) GetID() int {
	return b.ID
}

func (b BadgeDoubleDouble) Validate(leg *Leg) (bool, *int) {
	visit := leg.GetLastVisit()
	return visit.SecondDart.IsDouble(), &visit.PlayerID
}

func (b BadgeTripleDouble) GetID() int {
	return b.ID
}

func (b BadgeTripleDouble) Validate(leg *Leg) (bool, *int) {
	visit := leg.GetLastVisit()
	return visit.FirstDart.IsDouble() && visit.SecondDart.IsDouble(), &visit.PlayerID
}

func (b BadgeMadHouse) GetID() int {
	return b.ID
}

func (b BadgeMadHouse) Validate(leg *Leg) (bool, *int) {
	visit := leg.GetLastVisit()
	last := visit.GetLastDart()
	return last.IsDouble() && last.ValueRaw() == 1, &visit.PlayerID
}

func (b BadgeBigFish) GetID() int {
	return b.ID
}

func (b BadgeBigFish) Validate(leg *Leg) (bool, *int) {
	visit := leg.GetLastVisit()
	return visit.FirstDart.IsTriple() && visit.FirstDart.ValueRaw() == 20 &&
		visit.SecondDart.IsTriple() && visit.SecondDart.ValueRaw() == 20 &&
		visit.ThirdDart.IsDouble() && visit.ThirdDart.IsBull(), &visit.PlayerID
}
