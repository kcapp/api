package models

import (
	"encoding/json"
	"time"

	"github.com/guregu/null"
)

// Badge represents a badge model.
type Badge struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Hidden      bool     `json:"hidden"`
	Secret      bool     `json:"secret"`
	Filename    string   `json:"filename"`
	Levels      null.Int `json:"levels,omitempty"`
}

// BadgeStatistics represents badge statistics.
type BadgeStatistics struct {
	BadgeID         int       `json:"badge_id"`
	Level           null.Int  `json:"level,omitempty"`
	Value           null.Int  `json:"value,omitempty"`
	UnlockedPlayers int       `json:"unlocked_players"`
	UnlockedPercent float32   `json:"unlocked_percent"`
	FirstUnlock     null.Time `json:"first_unlock"`
	Players         []int     `json:"players"`
}

// PlayerBadge represents a Player2Badge model.
type PlayerBadge struct {
	Badge            *Badge    `json:"badge"`
	PlayerID         int       `json:"player_id"`
	Level            null.Int  `json:"level,omitempty"`
	LegID            null.Int  `json:"leg_id,omitempty"`
	Value            null.Int  `json:"value,omitempty"`
	MatchID          null.Int  `json:"match_id,omitempty"`
	OpponentPlayerID null.Int  `json:"opponent_player_id,omitempty"`
	TournamentID     null.Int  `json:"tournament_id,omitempty"`
	VisitID          null.Int  `json:"visit_id,omitempty"`
	Darts            []*Dart   `json:"darts,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// MarshalJSON will marshall the given object to JSON
func (pb PlayerBadge) MarshalJSON() ([]byte, error) {
	// Use a type to get consistent order of JSON key-value pairs.
	type playerBadgeJSON struct {
		Badge            *Badge    `json:"badge"`
		PlayerID         int       `json:"player_id"`
		Level            null.Int  `json:"level,omitempty"`
		LegID            null.Int  `json:"leg_id,omitempty"`
		Value            null.Int  `json:"value,omitempty"`
		MatchID          null.Int  `json:"match_id,omitempty"`
		OpponentPlayerID null.Int  `json:"opponent_player_id,omitempty"`
		TournamentID     null.Int  `json:"tournament_id,omitempty"`
		VisitID          null.Int  `json:"visit_id,omitempty"`
		Darts            []*Dart   `json:"darts,omitempty"`
		DartsString      string    `json:"darts_string,omitempty"`
		CreatedAt        time.Time `json:"created_at"`
	}
	var dartsString string
	if pb.Darts != nil {
		dartsString = pb.Darts[0].String()
		if pb.Darts[1].Value.Valid {
			dartsString += " " + pb.Darts[1].String()
		}
		if pb.Darts[2].Value.Valid {
			dartsString += " " + pb.Darts[2].String()
		}
	}
	return json.Marshal(playerBadgeJSON{
		Badge:            pb.Badge,
		PlayerID:         pb.PlayerID,
		Level:            pb.Level,
		LegID:            pb.LegID,
		Value:            pb.Value,
		MatchID:          pb.MatchID,
		OpponentPlayerID: pb.OpponentPlayerID,
		TournamentID:     pb.TournamentID,
		VisitID:          pb.VisitID,
		Darts:            pb.Darts,
		DartsString:      dartsString,
		CreatedAt:        pb.CreatedAt,
	})
}

// PlayerBadgeStatistics struct used for storing badge statistics
type PlayerBadgeStatistics struct {
	PlayerID      int
	Score100sPlus int
	Score140sPlus int
	Score180s     int
}

type GlobalBadge interface {
	GetID() int
}

type BadgeKcappSupporter struct{ ID int }
type BadgeSayMyName struct{ ID int }
type BadgeItsOfficial struct{ ID int }
type BadgeTournament1st struct{ ID int }
type BadgeTournament2nd struct{ ID int }
type BadgeTournament3rd struct{ ID int }
type BadgeUntouchable struct{ ID int }

func (b BadgeKcappSupporter) GetID() int {
	return 4
}

func (b BadgeSayMyName) GetID() int {
	return 12
}

func (b BadgeItsOfficial) GetID() int {
	return 17
}

func (b BadgeTournament1st) GetID() int {
	return 18
}

func (b BadgeTournament2nd) GetID() int {
	return 19
}

func (b BadgeTournament3rd) GetID() int {
	return 20
}

func (b BadgeUntouchable) GetID() int {
	return 26
}

var MatchBadges = []MatchBadge{}

type MatchBadge interface {
	GetID() int
	Validate(*Match) (bool, *int)
}

var LegBadges = []LegBadge{
	BadgeDoubleDouble{ID: 6},
	BadgeTripleDouble{ID: 7},
	BadgeMadHouse{ID: 8},
	BadgeMerryChristmas{ID: 9},
	BadgeHappyNewYear{ID: 10},
	BadgeBigFish{ID: 11},
	BadgeGettingCrowded{ID: 13},
	BadgeBullseye{ID: 14},
	BadgeEasyAs123{ID: 15},
	BadgeCloseToPerfect{ID: 16},
}

type LegBadge interface {
	GetID() int
	Validate(*Leg) (bool, *int, *int)
}

type BadgeDoubleDouble struct{ ID int }
type BadgeTripleDouble struct{ ID int }
type BadgeMadHouse struct{ ID int }
type BadgeMerryChristmas struct{ ID int }
type BadgeHappyNewYear struct{ ID int }
type BadgeBigFish struct{ ID int }
type BadgeGettingCrowded struct{ ID int }
type BadgeBullseye struct{ ID int }
type BadgeEasyAs123 struct{ ID int }
type BadgeCloseToPerfect struct{ ID int }

func (b BadgeDoubleDouble) GetID() int {
	return b.ID
}
func (b BadgeDoubleDouble) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	doubles := 0
	if visit.ThirdDart.IsDouble() {
		doubles++
	}
	if visit.SecondDart.IsDouble() {
		doubles++
	}
	if visit.FirstDart.IsDouble() {
		doubles++
	}
	return doubles == 2, &visit.PlayerID, &visit.ID
}

func (b BadgeTripleDouble) GetID() int {
	return b.ID
}
func (b BadgeTripleDouble) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	return visit.FirstDart.IsDouble() && visit.SecondDart.IsDouble() && visit.ThirdDart.IsDouble(), &visit.PlayerID, &visit.ID
}

func (b BadgeMadHouse) GetID() int {
	return b.ID
}
func (b BadgeMadHouse) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	last := visit.GetLastDart()
	return last.IsDouble() && last.ValueRaw() == 1, &visit.PlayerID, &visit.ID
}

func (b BadgeMerryChristmas) GetID() int {
	return b.ID
}
func (b BadgeMerryChristmas) Validate(leg *Leg) (bool, *int, *int) {
	d := leg.Endtime.Time
	return d.Day() == 25 && d.Month() == 12, nil, nil
}

func (b BadgeHappyNewYear) GetID() int {
	return b.ID
}
func (b BadgeHappyNewYear) Validate(leg *Leg) (bool, *int, *int) {
	d := leg.Endtime.Time
	return d.Day() == 31 && d.Month() == 12, nil, nil
}

func (b BadgeBigFish) GetID() int {
	return b.ID
}
func (b BadgeBigFish) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	return visit.FirstDart.IsTriple() && visit.FirstDart.ValueRaw() == 20 &&
		visit.SecondDart.IsTriple() && visit.SecondDart.ValueRaw() == 20 &&
		visit.ThirdDart.IsDouble() && visit.ThirdDart.IsBull(), &visit.PlayerID, &visit.ID
}

func (b BadgeGettingCrowded) GetID() int {
	return b.ID
}
func (b BadgeGettingCrowded) Validate(leg *Leg) (bool, *int, *int) {
	return len(leg.Players) > 4, nil, nil
}

func (b BadgeBullseye) GetID() int {
	return b.ID
}
func (b BadgeBullseye) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	last := visit.GetLastDart()
	return last.ValueRaw() == BULLSEYE && last.Multiplier == DOUBLE, &visit.PlayerID, &visit.ID
}

func (b BadgeEasyAs123) GetID() int {
	return b.ID
}
func (b BadgeEasyAs123) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	last := visit.GetLastDart()
	return visit.GetScore() == 123 && last.IsDouble(), &visit.PlayerID, &visit.ID
}

func (b BadgeCloseToPerfect) GetID() int {
	return b.ID
}
func (b BadgeCloseToPerfect) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	return leg.StartingScore == 501 && visit.DartsThrown < 15 && visit.DartsThrown > 9, &visit.PlayerID, nil
}

var LegPlayerBadges = []LegPlayerBadge{
	BadgeImpersonator{ID: 21},
	BadgeBotBeaterEasy{ID: 22},
	BadgeBotBeaterMedium{ID: 23},
	BadgeBotBeaterHard{ID: 24},
}

type LegPlayerBadge interface {
	GetID() int
	Validate(*Leg, []*Player2Leg) (bool, *int)
}

type BadgeImpersonator struct{ ID int }
type BadgeBotBeaterEasy struct{ ID int }
type BadgeBotBeaterMedium struct{ ID int }
type BadgeBotBeaterHard struct{ ID int }

func (b BadgeImpersonator) GetID() int {
	return b.ID
}
func (b BadgeImpersonator) Validate(leg *Leg, players []*Player2Leg) (bool, *int) {
	var bot *Player2Leg
	for _, p2l := range players {
		if p2l.Player.IsBot && p2l.BotConfig.PlayerID.Valid {
			bot = p2l
		}
	}
	winner := int(leg.WinnerPlayerID.Int64)
	if bot != nil && bot.PlayerID != winner {
		return true, &winner
	}
	return false, nil
}

func (b BadgeBotBeaterEasy) GetID() int {
	return b.ID
}
func (b BadgeBotBeaterEasy) Validate(leg *Leg, players []*Player2Leg) (bool, *int) {
	bot := getBot(BOT_EASY, players)
	winner := int(leg.WinnerPlayerID.Int64)
	if bot != nil && bot.PlayerID != winner {
		return true, &winner
	}
	return false, nil
}

func (b BadgeBotBeaterMedium) GetID() int {
	return b.ID
}
func (b BadgeBotBeaterMedium) Validate(leg *Leg, players []*Player2Leg) (bool, *int) {
	bot := getBot(BOT_MEDIUM, players)
	winner := int(leg.WinnerPlayerID.Int64)
	if bot != nil && bot.PlayerID != winner {
		return true, &winner
	}
	return false, nil
}

func (b BadgeBotBeaterHard) GetID() int {
	return b.ID
}
func (b BadgeBotBeaterHard) Validate(leg *Leg, players []*Player2Leg) (bool, *int) {
	bot := getBot(BOT_HARD, players)
	winner := int(leg.WinnerPlayerID.Int64)
	if bot != nil && bot.PlayerID != winner {
		return true, &winner
	}
	return false, nil
}

var VisitBadges = []VisitBadge{
	BadgeHighScore{ID: 1},
	BadgeHigherScore{ID: 2},
	BadgeTheMaximum{ID: 3},
}

type VisitBadge interface {
	GetID() int
	Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int)
	Levels() []int
}

type BadgeHighScore struct{ ID int }
type BadgeHigherScore struct{ ID int }
type BadgeTheMaximum struct{ ID int }

func (b BadgeHighScore) GetID() int {
	return b.ID
}

func (b BadgeHighScore) Levels() []int {
	return []int{1, 10, 100, 1000}
}

func (b BadgeHighScore) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int) {
	count := 0
	for _, visit := range visits {
		if visit.PlayerID != stats.PlayerID {
			continue
		}
		if visit.Score >= 100 && visit.Score < 140 {
			count++
		}
	}
	if count > 0 {
		level := getLevel(stats.Score100sPlus+count, b.Levels())
		return true, &level
	}
	return false, nil
}

func (b BadgeHigherScore) GetID() int {
	return b.ID
}

func (b BadgeHigherScore) Levels() []int {
	return []int{1, 10, 100, 1000}
}

func (b BadgeHigherScore) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int) {
	count := 0
	for _, visit := range visits {
		if visit.PlayerID != stats.PlayerID {
			continue
		}
		if visit.Score >= 140 && visit.Score < 180 {
			count++
		}
	}
	if count > 0 {
		level := getLevel(stats.Score140sPlus+count, b.Levels())
		return true, &level
	}
	return false, nil
}

func (b BadgeTheMaximum) GetID() int {
	return b.ID
}

func (b BadgeTheMaximum) Levels() []int {
	return []int{1, 10, 50, 100}
}

func (b BadgeTheMaximum) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int) {
	count := 0
	for _, visit := range visits {
		if visit.PlayerID != stats.PlayerID {
			continue
		}
		if visit.Score == 180 {
			count++
		}
	}
	if count > 0 {
		level := getLevel(stats.Score180s+count, b.Levels())
		return true, &level
	}
	return false, nil
}

func getLevel(value int, levels []int) int {
	level := 1
	for i, treshold := range levels {
		if value > treshold {
			level = i + 1
		}
	}
	return level
}

func getBot(skill int64, players []*Player2Leg) *Player2Leg {
	for _, p2l := range players {
		if p2l.Player.IsBot && p2l.BotConfig.Skill.Int64 == skill {
			return p2l
		}
	}
	return nil
}
