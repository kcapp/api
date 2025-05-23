package models

import (
	"encoding/json"
	"sort"
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
	Badge            *Badge                 `json:"badge"`
	PlayerID         int                    `json:"player_id"`
	Level            null.Int               `json:"level,omitempty"`
	LegID            null.Int               `json:"leg_id,omitempty"`
	Value            null.Int               `json:"value,omitempty"`
	MatchID          null.Int               `json:"match_id,omitempty"`
	OpponentPlayerID int                    `json:"opponent_player_id,omitempty"`
	TournamentID     null.Int               `json:"tournament_id,omitempty"`
	VisitID          null.Int               `json:"visit_id,omitempty"`
	Darts            []*Dart                `json:"darts,omitempty"`
	Data             null.String            `json:"data,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	Statistics       *PlayerBadgeStatistics `json:"statistics,omitempty"`
}

// MarshalJSON will marshall the given object to JSON
func (pb PlayerBadge) MarshalJSON() ([]byte, error) {
	// Use a type to get consistent order of JSON key-value pairs.
	type playerBadgeJSON struct {
		Badge            *Badge                 `json:"badge"`
		PlayerID         int                    `json:"player_id"`
		Level            null.Int               `json:"level,omitempty"`
		LegID            null.Int               `json:"leg_id,omitempty"`
		Value            null.Int               `json:"value,omitempty"`
		MatchID          null.Int               `json:"match_id,omitempty"`
		OpponentPlayerID null.Int               `json:"opponent_player_id,omitempty"`
		TournamentID     null.Int               `json:"tournament_id,omitempty"`
		VisitID          null.Int               `json:"visit_id,omitempty"`
		Darts            []*Dart                `json:"darts,omitempty"`
		DartsString      string                 `json:"darts_string,omitempty"`
		Data             null.String            `json:"data,omitempty"`
		CreatedAt        time.Time              `json:"created_at"`
		Statistics       *PlayerBadgeStatistics `json:"statistics,omitempty"`
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

	var opponentPlayerID null.Int
	if pb.OpponentPlayerID != 0 {
		// 0 means NULL for unique constraint, so make it null here in the API
		opponentPlayerID = null.IntFrom(int64(pb.OpponentPlayerID))
	}

	return json.Marshal(playerBadgeJSON{
		Badge:            pb.Badge,
		PlayerID:         pb.PlayerID,
		Level:            pb.Level,
		LegID:            pb.LegID,
		Value:            pb.Value,
		MatchID:          pb.MatchID,
		OpponentPlayerID: opponentPlayerID,
		TournamentID:     pb.TournamentID,
		VisitID:          pb.VisitID,
		Darts:            pb.Darts,
		DartsString:      dartsString,
		Data:             pb.Data,
		CreatedAt:        pb.CreatedAt,
		Statistics:       pb.Statistics,
	})
}

// PlayerBadgeStatistics struct used for storing badge statistics
type PlayerBadgeStatistics struct {
	PlayerID      int                 `json:"player_id"`
	Score100sPlus int                 `json:"score_100_plus"`
	Score140sPlus int                 `json:"score_140_plus"`
	Score180s     int                 `json:"score_180s"`
	Shanghais     []int               `json:"shanghais"`
	BadgeMap      map[int]interface{} `json:"values"`
}

// MarshalJSON will marshall the given object to JSON
func (pbs PlayerBadgeStatistics) MarshalJSON() ([]byte, error) {
	// Use a type to get consistent order of JSON key-value pairs.
	type playerBadgeStatisticsJSON struct {
		PlayerID      int                 `json:"player_id"`
		Score100sPlus int                 `json:"score_100_plus"`
		Score140sPlus int                 `json:"score_140_plus"`
		Score180s     int                 `json:"score_180s"`
		Shanghais     []int               `json:"shanghais"`
		BadgeMap      map[int]interface{} `json:"values"`
	}

	pbs.BadgeMap = make(map[int]interface{}, 0)
	pbs.BadgeMap[1] = pbs.Score100sPlus
	pbs.BadgeMap[2] = pbs.Score140sPlus
	pbs.BadgeMap[3] = pbs.Score180s
	sort.Ints(pbs.Shanghais)
	pbs.BadgeMap[46] = pbs.Shanghais

	return json.Marshal(playerBadgeStatisticsJSON{
		PlayerID:      pbs.PlayerID,
		Score100sPlus: pbs.Score100sPlus,
		Score140sPlus: pbs.Score140sPlus,
		Score180s:     pbs.Score180s,
		Shanghais:     pbs.Shanghais,
		BadgeMap:      pbs.BadgeMap,
	})
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
type BadgeByeForNow struct{ ID int }
type BadgeOldTimer struct{ ID int }

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

func (b BadgeByeForNow) GetID() int {
	return 27
}

func (b BadgeOldTimer) GetID() int {
	return 28
}

type GlobalLevelBadge interface {
	GetID() int
	Levels() []int
}

type BadgeVersatilePlayer struct{ ID int }

func (b BadgeVersatilePlayer) GetID() int {
	return 29
}

func (b BadgeVersatilePlayer) Levels() []int {
	return []int{5, 10, 15, 20}
}

var LeaderboardBadges = []LeaderboardBadge{
	BadgeKingslayer{ID: 48},
}

type LeaderboardBadge interface {
	GetID() int
	// Validate returns bool, player.ID
	Validate(match *Match, matchStatistics []*StatisticsX01, leaderboard []*StatisticsX01) (bool, *int, *int)
}

type BadgeKingslayer struct{ ID int }

func (b BadgeKingslayer) GetID() int {
	return b.ID
}
func (b BadgeKingslayer) Validate(match *Match, matchStatistics []*StatisticsX01, leaderboard []*StatisticsX01) (bool, *int, *int) {
	if !match.IsX01() {
		return false, nil, nil
	}

	// Get the current king in the correct office
	var king *int
	var kingOfKings *int
	if len(leaderboard) > 0 {
		// One king to rule them all...
		kingOfKings = &leaderboard[0].PlayerID
		king = kingOfKings
	}
	if match.OfficeID.Valid {
		// Get the local king
		for _, stat := range leaderboard {
			if stat.OfficeID.Int64 == match.OfficeID.Int64 {
				king = &stat.PlayerID
				break
			}
		}
	}

	if king == nil {
		// The king is in a different office
		return false, nil, nil
	}
	if !match.MatchMode.IsChallenge {
		// No usurpers in friendly skirmishes
		return false, nil, nil
	}

	if len(match.Players) != 2 {
		// There can be only one
		return false, nil, nil
	}
	if !containsInt(match.Players, *king) && !containsInt(match.Players, *kingOfKings) {
		// This was not an attempt to slay the king
		return false, nil, nil
	}
	if containsInt(match.Players, *kingOfKings) {
		// It's the one true king
		king = kingOfKings
	}

	if !match.WinnerID.Valid || match.WinnerID.Int64 == int64(*king) {
		// If you come for the king, you best not miss...
		return false, nil, nil
	}

	kingStatistics, _ := findPlayerX01Statistics(matchStatistics, *king)
	slayerStatistics, _ := findPlayerX01Statistics(matchStatistics, int(match.WinnerID.Int64))
	if slayerStatistics.ThreeDartAvg < kingStatistics.ThreeDartAvg {
		// You must best the king, not just wound him.
		return false, nil, nil
	}

	// The king is dead, long live the king!
	kingslayer := int(match.WinnerID.Int64)
	return true, &kingslayer, king
}

var MatchBadges = []MatchBadge{
	BadgeJustAQuickie{ID: 37},
	BadgeAroundTheWorld{ID: 38},
	BadgeOfficiallyGood{ID: 41},
}

type MatchBadge interface {
	GetID() int
	// Validate returns bool, player.ID
	Validate(match *Match) (bool, []int)
}

type BadgeJustAQuickie struct{ ID int }
type BadgeAroundTheWorld struct{ ID int }
type BadgeOfficiallyGood struct{ ID int }

func (b BadgeJustAQuickie) GetID() int {
	return b.ID
}
func (b BadgeJustAQuickie) Validate(match *Match) (bool, []int) {
	if !match.IsX01() {
		return false, nil
	}
	if len(match.Legs) == 3 && len(match.Players) > 1 {
		first := match.Legs[0]
		second := match.Legs[1]
		third := match.Legs[2]
		if first.GetLastVisit().CreatedAt.Sub(first.Visits[0].CreatedAt).Minutes() <= 3 &&
			second.GetLastVisit().CreatedAt.Sub(second.Visits[0].CreatedAt).Minutes() <= 3 &&
			third.GetLastVisit().CreatedAt.Sub(third.Visits[0].CreatedAt).Minutes() <= 3 {
			return true, []int{int(match.WinnerID.Int64)}
		}
	}

	return false, nil
}

func (b BadgeAroundTheWorld) GetID() int {
	return b.ID
}
func (b BadgeAroundTheWorld) Validate(match *Match) (bool, []int) {
	if !match.IsX01() {
		return false, nil
	}

	playerHits := make(map[int][]int)
	for playerID := range match.Players {
		playerHits[playerID] = make([]int, 0)
	}

	for _, leg := range match.Legs {
		for _, visit := range leg.Visits {
			if !visit.IsBust {
				playerHits[visit.PlayerID] = append(playerHits[visit.PlayerID], visit.FirstDart.ValueRaw())
				playerHits[visit.PlayerID] = append(playerHits[visit.PlayerID], visit.SecondDart.ValueRaw())
				playerHits[visit.PlayerID] = append(playerHits[visit.PlayerID], visit.ThirdDart.ValueRaw())
			}
		}
	}

	players := make([]int, 0)
	for playerID, hits := range playerHits {
		allHit := true
		for i := 1; i <= 20; i++ {
			if !containsInt(hits, i) {
				allHit = false
				break
			}
		}
		if !containsInt(hits, 25) {
			allHit = false
		}
		if allHit {
			players = append(players, playerID)
		}
	}
	if len(players) > 0 {
		return true, players
	}
	return false, nil
}

func (b BadgeOfficiallyGood) GetID() int {
	return b.ID
}
func (b BadgeOfficiallyGood) Validate(match *Match) (bool, []int) {
	playerIDs := make([]int, 0)
	if match.TournamentID.Valid {
		for _, leg := range match.Legs {
			for _, visit := range leg.Visits {
				if visit.GetScore() == 180 {
					playerIDs = append(playerIDs, visit.PlayerID)
				}
			}
		}
	}
	if len(playerIDs) > 0 {
		return true, playerIDs
	}
	return false, nil
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
	BadgeLittleFish{ID: 33},
	BadgeShanghaiCheckout{ID: 36},
	BadgeTripleTrouble{ID: 39},
	BadgePerfection{ID: 40},
	BadgeChampagneShot{ID: 42},
	BadgeYin{ID: 44},
	BadgeYang{ID: 45},
	BadgeZebra{ID: 47},
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
type BadgeLittleFish struct{ ID int }
type BadgeShanghaiCheckout struct{ ID int }
type BadgeTripleThreat struct{ ID int }
type BadgeBabyTon struct{ ID int }
type BadgeBullBullBull struct{ ID int }
type BadgeSoClose struct{ ID int }
type BadgeTripleTrouble struct{ ID int }
type BadgePerfection struct{ ID int }
type BadgeChampagneShot struct{ ID int }
type BadgeYin struct{ ID int }
type BadgeYang struct{ ID int }
type BadgeZebra struct{ ID int }

func (b BadgeDoubleDouble) GetID() int {
	return b.ID
}
func (b BadgeDoubleDouble) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
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
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
	visit := leg.GetLastVisit()
	return visit.FirstDart.IsDouble() && visit.SecondDart.IsDouble() && visit.ThirdDart.IsDouble(), &visit.PlayerID, &visit.ID
}

func (b BadgeMadHouse) GetID() int {
	return b.ID
}
func (b BadgeMadHouse) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
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
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
	visit := leg.GetLastVisit()
	last := visit.GetLastDart()
	return last.ValueRaw() == BULLSEYE && last.Multiplier == DOUBLE, &visit.PlayerID, &visit.ID
}

func (b BadgeEasyAs123) GetID() int {
	return b.ID
}
func (b BadgeEasyAs123) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
	visit := leg.GetLastVisit()
	last := visit.GetLastDart()
	return visit.GetScore() == 123 && last.IsDouble(), &visit.PlayerID, &visit.ID
}

func (b BadgeCloseToPerfect) GetID() int {
	return b.ID
}
func (b BadgeCloseToPerfect) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
	visit := leg.GetLastVisit()
	return leg.StartingScore == 501 && visit.DartsThrown < 15 && visit.DartsThrown > 9, &visit.PlayerID, nil
}

func (b BadgeLittleFish) GetID() int {
	return b.ID
}
func (b BadgeLittleFish) Validate(leg *Leg) (bool, *int, *int) {
	visit := leg.GetLastVisit()
	return visit.GetScore() == 130 &&
		visit.FirstDart.ValueRaw() == 20 && (visit.FirstDart.IsSingle() || visit.FirstDart.IsTriple()) &&
		visit.SecondDart.ValueRaw() == 20 && (visit.SecondDart.IsSingle() || visit.SecondDart.IsTriple()) &&
		visit.ThirdDart.IsBull() && visit.ThirdDart.IsDouble(), &visit.PlayerID, &visit.ID
}

func (b BadgeShanghaiCheckout) GetID() int {
	return b.ID
}
func (b BadgeShanghaiCheckout) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
	visit := leg.GetLastVisit()
	return visit.IsShanghai(), &visit.PlayerID, &visit.ID
}

func (b BadgeTripleTrouble) GetID() int {
	return b.ID
}
func (b BadgeTripleTrouble) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}

	visit := leg.GetLastVisit()
	if visit.FirstDart.IsDouble() && visit.SecondDart.IsDouble() && visit.ThirdDart.IsDouble() &&
		visit.FirstDart.ValueRaw() == visit.SecondDart.ValueRaw() && visit.FirstDart.ValueRaw() == visit.ThirdDart.ValueRaw() {
		return true, &visit.PlayerID, &visit.ID

	}
	return false, nil, nil
}

func (b BadgePerfection) GetID() int {
	return b.ID
}
func (b BadgePerfection) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}

	visit := leg.GetLastVisit()
	if leg.StartingScore == 501 && visit.DartsThrown == 9 {
		return true, &visit.PlayerID, nil
	}
	return false, nil, nil
}

func (b BadgeChampagneShot) GetID() int {
	return b.ID
}
func (b BadgeChampagneShot) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}

	visit := leg.GetLastVisit()

	if visit.FirstDart.IsBull() && visit.FirstDart.IsDouble() &&
		visit.SecondDart.IsBull() && visit.SecondDart.IsDouble() &&
		visit.ThirdDart.ValueRaw() == 16 && visit.ThirdDart.IsDouble() {
		return true, &visit.PlayerID, &visit.ID

	}
	return false, nil, nil
}

func (b BadgeYin) GetID() int {
	return b.ID
}
func (b BadgeYin) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}

	completed := true
	winner := leg.GetLastVisit().PlayerID
	for _, visit := range leg.Visits {
		if visit.PlayerID != winner {
			continue
		}
		if !containsInt(NUMS_BLACK, visit.FirstDart.ValueRaw()) ||
			!containsInt(NUMS_BLACK, visit.SecondDart.ValueRaw()) ||
			!containsInt(NUMS_BLACK, visit.ThirdDart.ValueRaw()) {
			completed = false
			break

		}
	}
	if completed {
		return true, &winner, nil
	}
	return false, nil, nil
}

func (b BadgeYang) GetID() int {
	return b.ID
}
func (b BadgeYang) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}

	completed := true
	winner := leg.GetLastVisit().PlayerID
	for _, visit := range leg.Visits {
		if visit.PlayerID != winner {
			continue
		}
		if !containsInt(NUMS_WHITE, visit.FirstDart.ValueRaw()) ||
			!containsInt(NUMS_WHITE, visit.SecondDart.ValueRaw()) ||
			!containsInt(NUMS_WHITE, visit.ThirdDart.ValueRaw()) {
			completed = false
			break

		}
	}
	if completed {
		return true, &winner, nil
	}
	return false, nil, nil
}

func (b BadgeZebra) GetID() int {
	return b.ID
}
func (b BadgeZebra) Validate(leg *Leg) (bool, *int, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil, nil
	}
	winner := leg.GetLastVisit().PlayerID

	COLORS := map[int][]int{
		0: NUMS_BLACK,
		1: NUMS_WHITE,
	}
	completed := true

	first := leg.GetFirstHitDart(winner)
	var color int
	if containsInt(COLORS[0], first.ValueRaw()) {
		color = 0
	} else {
		color = 1
	}
	for _, visit := range leg.Visits {
		if visit.PlayerID != winner {
			continue
		}

		for _, dart := range visit.GetDarts() {
			if dart.IsMiss() {
				continue
			}
			colors := COLORS[color%2]
			if !containsInt(colors, dart.ValueRaw()) {
				completed = false
				break
			}
			color++
		}

		if !completed {
			break
		}
	}
	if completed {
		return true, &winner, nil
	}
	return false, nil, nil
}

var LegPlayerBadges = []LegPlayerBadge{
	BadgeImpersonator{ID: 21},
	BadgeBotBeaterEasy{ID: 22},
	BadgeBotBeaterMedium{ID: 23},
	BadgeBotBeaterHard{ID: 24},
	BadgeBeerGame{ID: 43},
}

type LegPlayerBadge interface {
	GetID() int
	// Validate returns bool, player.ID
	Validate(*Leg, []*Player2Leg) (bool, *int)
}

type BadgeImpersonator struct{ ID int }
type BadgeBotBeaterEasy struct{ ID int }
type BadgeBotBeaterMedium struct{ ID int }
type BadgeBotBeaterHard struct{ ID int }
type BadgeBeerGame struct{ ID int }

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

func (b BadgeBeerGame) GetID() int {
	return b.ID
}
func (b BadgeBeerGame) Validate(leg *Leg, players []*Player2Leg) (bool, *int) {
	if !leg.IsX01() || !leg.IsLegCheckout() {
		return false, nil
	}
	visit := leg.GetLastVisit()

	for _, player := range players {
		if player.PlayerID == visit.PlayerID {
			continue
		}
		if player.CurrentScore >= 200 {
			return true, &visit.PlayerID
		}
	}
	return false, nil
}

var VisitBadgesLevel = []VisitBadgeLevel{
	BadgeHighScore{ID: 1},
	BadgeHigherScore{ID: 2},
	BadgeTheMaximum{ID: 3},
	BadgeMonotonous{ID: 30},
	BadgeShanghai{ID: 46},
}

type VisitBadgeLevel interface {
	GetID() int
	// Validate returns bool, level, visit.ID
	Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int, *int)
	Levels() []int
}

type BadgeHighScore struct{ ID int }
type BadgeHigherScore struct{ ID int }
type BadgeTheMaximum struct{ ID int }
type BadgeMonotonous struct{ ID int }
type BadgeShanghai struct{ ID int }

func (b BadgeHighScore) GetID() int {
	return b.ID
}
func (b BadgeHighScore) Levels() []int {
	return []int{1, 10, 100, 1000}
}
func (b BadgeHighScore) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int, *int) {
	count := 0
	playerVisits := getVisitsForPlayer(visits, stats.PlayerID)
	for _, visit := range playerVisits {
		if visit.Score >= 100 && visit.Score < 140 {
			count++
		}
	}
	if count > 0 {
		level := GetLevel(stats.Score100sPlus+count, b.Levels())
		return true, &level, nil
	}
	return false, nil, nil
}

func (b BadgeHigherScore) GetID() int {
	return b.ID
}
func (b BadgeHigherScore) Levels() []int {
	return []int{1, 10, 100, 1000}
}
func (b BadgeHigherScore) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int, *int) {
	count := 0
	playerVisits := getVisitsForPlayer(visits, stats.PlayerID)
	for _, visit := range playerVisits {
		if visit.PlayerID != stats.PlayerID {
			continue
		}
		if visit.Score >= 140 && visit.Score < 180 {
			count++
		}
	}
	if count > 0 {
		level := GetLevel(stats.Score140sPlus+count, b.Levels())
		return true, &level, nil
	}
	return false, nil, nil
}

func (b BadgeTheMaximum) GetID() int {
	return b.ID
}
func (b BadgeTheMaximum) Levels() []int {
	return []int{1, 10, 50, 100}
}
func (b BadgeTheMaximum) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int, *int) {
	count := 0
	playerVisits := getVisitsForPlayer(visits, stats.PlayerID)
	for _, visit := range playerVisits {
		if visit.PlayerID != stats.PlayerID {
			continue
		}
		if visit.Score == 180 {
			count++
		}
	}
	if count > 0 {
		level := GetLevel(stats.Score180s+count, b.Levels())
		return true, &level, nil
	}
	return false, nil, nil
}

func (b BadgeMonotonous) GetID() int {
	return b.ID
}
func (b BadgeMonotonous) Levels() []int {
	return []int{3, 4, 5, 6}
}
func (b BadgeMonotonous) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int, *int) {
	playerVisits := getVisitsForPlayer(visits, stats.PlayerID)
	for i := len(b.Levels()) - 1; i >= 0; i-- {
		monotonous, visit := hasSameVisitsInARow(playerVisits, b.Levels()[i])
		if monotonous {
			level := i + 1
			return true, &level, &visit.ID
		}
	}
	return false, nil, nil
}

func (b BadgeShanghai) GetID() int {
	return b.ID
}
func (b BadgeShanghai) Levels() []int {
	return []int{1, 5, 10, 15, 20}
}
func (b BadgeShanghai) Validate(stats *PlayerBadgeStatistics, visits []*Visit) (bool, *int, *int) {
	count := 0
	playerVisits := getVisitsForPlayer(visits, stats.PlayerID)
	for _, visit := range playerVisits {
		if visit.IsShanghai() && containsInt(stats.Shanghais, visit.FirstDart.ValueRaw()) {
			count++
		}
	}
	if count > 0 {
		level := GetLevel(len(stats.Shanghais)+count, b.Levels())
		return true, &level, nil
	}
	return false, nil, nil
}

var VisitBadges = []VisitBadge{
	BadgeTripleThreat{ID: 31},
	BadgeBabyTon{ID: 32},
	BadgeBullBullBull{ID: 34},
	BadgeSoClose{ID: 35},
}

type VisitBadge interface {
	GetID() int
	// Validate returns bool, visit.ID
	Validate(playerID int, visits []*Visit) (bool, *int)
}

func (b BadgeTripleThreat) GetID() int {
	return b.ID
}
func (b BadgeTripleThreat) Validate(playerID int, visits []*Visit) (bool, *int) {
	values := []int{20, 19, 18}
	playerVisits := getVisitsForPlayer(visits, playerID)
	for _, visit := range playerVisits {
		if visit.GetScore() == 168 &&
			visit.FirstDart.IsTriple() && visit.SecondDart.IsTriple() && visit.ThirdDart.IsTriple() &&
			visit.FirstDart.IsValue(values) && visit.SecondDart.IsValue(values) && visit.ThirdDart.IsValue(values) {
			return true, &visit.ID
		}
	}
	return false, nil
}

func (b BadgeBabyTon) GetID() int {
	return b.ID
}
func (b BadgeBabyTon) Validate(playerID int, visits []*Visit) (bool, *int) {
	value := []int{19}
	playerVisits := getVisitsForPlayer(visits, playerID)
	for _, visit := range playerVisits {
		if visit.GetScore() == 95 && visit.FirstDart.IsValue(value) && visit.SecondDart.IsValue(value) && visit.ThirdDart.IsValue(value) &&
			// Only allow a Baby Ton to be T19, 19, 19
			!visit.FirstDart.IsDouble() && !visit.SecondDart.IsDouble() && !visit.ThirdDart.IsDouble() {
			return true, &visit.ID
		}
	}
	return false, nil
}

func (b BadgeBullBullBull) GetID() int {
	return b.ID
}
func (b BadgeBullBullBull) Validate(playerID int, visits []*Visit) (bool, *int) {
	playerVisits := getVisitsForPlayer(visits, playerID)
	for _, visit := range playerVisits {
		if visit.FirstDart.IsBull() && visit.FirstDart.IsDouble() &&
			visit.SecondDart.IsBull() && visit.SecondDart.IsDouble() &&
			visit.ThirdDart.IsBull() && visit.ThirdDart.IsDouble() {
			return true, &visit.ID
		}
	}
	return false, nil
}

func (b BadgeSoClose) GetID() int {
	return b.ID
}
func (b BadgeSoClose) Validate(playerID int, visits []*Visit) (bool, *int) {
	value := []int{1}
	playerVisits := getVisitsForPlayer(visits, playerID)
	for _, visit := range playerVisits {
		if visit.FirstDart.IsTriple() && visit.FirstDart.IsValue(value) &&
			visit.SecondDart.IsTriple() && visit.SecondDart.IsValue(value) &&
			visit.ThirdDart.IsTriple() && visit.ThirdDart.IsValue(value) {
			return true, &visit.ID
		}
	}
	return false, nil
}

func GetLevel(value int, levels []int) int {
	level := 1
	for i, treshold := range levels {
		if value >= treshold {
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

func hasSameVisitsInARow(visits []*Visit, numVisits int) (bool, *Visit) {
	if len(visits) < numVisits {
		return false, nil
	}

	for i := numVisits - 1; i < len(visits); i++ {
		sameVisits := true
		var visit *Visit
		for j := 0; j < numVisits-1; j++ {
			visit = visits[i-j]
			if visit.FirstDart.IsMiss() || visit.SecondDart.IsMiss() || visit.ThirdDart.IsMiss() {
				sameVisits = false
				break
			}
			if !visits[i-j].isEqualTo(*visits[i-j-1]) {
				sameVisits = false
				break
			}
		}
		if sameVisits {
			return true, visit
		}
	}
	return false, nil
}

func getVisitsForPlayer(visits []*Visit, playerID int) []*Visit {
	playerVisits := make([]*Visit, 0)
	for _, visit := range visits {
		if visit.PlayerID == playerID {
			playerVisits = append(playerVisits, visit)
		}
	}
	return playerVisits
}

// findPlayerX01Statistics takes a slice and returns the statistics for the given playerID
func findPlayerX01Statistics(stats []*StatisticsX01, playerID int) (*StatisticsX01, bool) {
	for _, stat := range stats {
		if stat.PlayerID == playerID {
			return stat, true
		}
	}
	return nil, false
}
