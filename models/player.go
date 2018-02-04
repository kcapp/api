package models

import (
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

// Player struct used for storing players
type Player struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Nickname     null.String `json:"nickname,omitempty"`
	GamesPlayed  int         `json:"games_played"`
	GamesWon     int         `json:"games_won"`
	PPD          float32     `json:"ppd,omitempty"`
	FirstNinePPD float32     `json:"first_nine_ppd,omitempty"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
}

// GetPlayers returns a map of all players
func GetPlayers() (map[int]*Player, error) {
	rows, err := db.Query(`SELECT p.id, p.name, p.nickname, p.games_played, p.games_won, p.created_at FROM player p`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*Player)
	for rows.Next() {
		p := new(Player)
		err := rows.Scan(&p.ID, &p.Name, &p.Nickname, &p.GamesPlayed, &p.GamesWon, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		players[p.ID] = p
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

// AddPlayer will add a new player to the database
func AddPlayer(player Player) error {
	// Prepare statement for inserting data
	stmt, err := db.Prepare("INSERT INTO player (name, nickname) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.Name, player.Nickname)
	return err // Will be nil if no error occured
}

// GetPlayerStatistics will get statistics about the given player id
func GetPlayerStatistics(id int) (*StatisticsX01, error) {
	s := new(StatisticsX01)
	err := db.QueryRow(`
		SELECT
			p.id,
			SUM(s.ppd) / p.games_played,
			SUM(s.first_nine_ppd) / p.games_played,
			SUM(s.60s_plus),
			SUM(s.100s_plus),
			SUM(s.140s_plus),
			SUM(s.180s),
			SUM(accuracy_20) / COUNT(accuracy_20),
			SUM(accuracy_19) / COUNT(accuracy_19),
			SUM(overall_accuracy) / COUNT(overall_accuracy),
			SUM(checkout_percentage) / COUNT(checkout_percentage)
		FROM statistics_x01 s
		JOIN player p ON p.id = s.player_id
		JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE s.player_id = ?
		GROUP BY s.player_id`, id).Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus,
		&s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
	if err != nil {
		return nil, err
	}

	visits, err := GetPlayerVisits(id)
	if err != nil {
		return nil, err
	}
	s.Hits, s.DartsThrown = GetHitsMap(visits)

	return s, nil
}

// GetPlayersStatistics will get statistics about all the the given player IDs
func GetPlayersStatistics(ids []int) ([]*StatisticsX01, error) {
	q, args, err := sqlx.In(`
		SELECT
			p.id,
			SUM(s.ppd) / p.games_played,
			SUM(s.first_nine_ppd) / p.games_played,
			SUM(s.60s_plus),
			SUM(s.100s_plus),
			SUM(s.140s_plus),
			SUM(s.180s),
			SUM(accuracy_20) / COUNT(accuracy_20),
			SUM(accuracy_19) / COUNT(accuracy_19),
			SUM(overall_accuracy) / COUNT(overall_accuracy),
			SUM(checkout_percentage) / COUNT(checkout_percentage)
		FROM statistics_x01 s
		JOIN player p ON p.id = s.player_id
		JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE s.player_id IN (?)
		GROUP BY s.player_id`, ids)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(q, args...)
	defer rows.Close()

	statisticsMap := make(map[int]*StatisticsX01)
	for rows.Next() {
		s := new(StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus,
			&s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		statisticsMap[s.PlayerID] = s
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate Best PPD, Best First 9, Best 301 and Best 501
	q, args, err = sqlx.In(`
		SELECT
			p.id,
			s.ppd,
			s.first_nine_ppd,
			s.checkout_percentage,
			s.darts_thrown,
			m.starting_score
		FROM statistics_x01 s
		JOIN player p ON p.id = s.player_id
		JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE s.player_id IN (?)`, ids)
	if err != nil {
		return nil, err
	}
	rows, err = db.Query(q, args...)
	defer rows.Close()

	rawStatistics := make([]*StatisticsX01, 0)
	for rows.Next() {
		s := new(StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.CheckoutPercentage, &s.DartsThrown, &s.StartingScore)
		if err != nil {
			return nil, err
		}
		rawStatistics = append(rawStatistics, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, stat := range rawStatistics {
		real := statisticsMap[stat.PlayerID]
		if stat.StartingScore.Int64 == 301 && (stat.DartsThrown < real.Best301 || real.Best301 == 0) {
			real.Best301 = stat.DartsThrown
		}
		if stat.StartingScore.Int64 == 501 && (stat.DartsThrown < real.Best501 || real.Best501 == 0) {
			real.Best501 = stat.DartsThrown
		}
		if stat.PPD > real.BestPPD {
			real.BestPPD = stat.PPD
		}
		if stat.FirstNinePPD > real.BestFirstNinePPD {
			real.BestFirstNinePPD = stat.FirstNinePPD
		}
	}

	q, args, err = sqlx.In(`
		SELECT
			s.player_id,
			MAX(IFNULL(s.first_dart * s.first_dart_multiplier, 0) +
			IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
			IFNULL(s.third_dart * s.third_dart_multiplier, 0)) AS 'highest_checkout'
		FROM score s
		JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE m.winner_id = s.player_id
			AND s.player_id IN (?)
			AND s.id IN (SELECT MAX(s.id) FROM score s JOIN `+"`match`"+`m ON m.id = s.match_id WHERE m.winner_id = s.player_id GROUP BY match_id)
		GROUP BY player_id
		ORDER BY highest_checkout DESC`, ids)
	if err != nil {
		return nil, err
	}
	rows, err = db.Query(q, args...)
	defer rows.Close()

	for rows.Next() {
		var playerID int
		var checkout int
		err := rows.Scan(&playerID, &checkout)
		if err != nil {
			return nil, err
		}
		statisticsMap[playerID].HighestCheckout = checkout
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	statistics := make([]*StatisticsX01, 0)
	for _, s := range statisticsMap {
		statistics = append(statistics, s)
	}
	return statistics, nil
}
