package data

import (
	"github.com/kcapp/api/models"
)

// GetTournaments will return all tournaments
func GetTournaments() ([]*models.Tournament, error) {
	rows, err := models.DB.Query("SELECT id, name, short_name, start_time, end_time FROM tournament ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := make([]*models.Tournament, 0)
	for rows.Next() {
		tournament := new(models.Tournament)
		err := rows.Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.StartTime, &tournament.EndTime)
		if err != nil {
			return nil, err
		}
		tournaments = append(tournaments, tournament)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tournaments, nil
}

// GetTournamentGroups will return all tournament groups
func GetTournamentGroups() (map[int]*models.TournamentGroup, error) {
	rows, err := models.DB.Query("SELECT id, name, division FROM tournament_group")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make(map[int]*models.TournamentGroup, 0)
	for rows.Next() {
		group := new(models.TournamentGroup)
		err := rows.Scan(&group.ID, &group.Name, &group.Division)
		if err != nil {
			return nil, err
		}
		groups[group.ID] = group
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetTournament will return a given tournament
func GetTournament(id int) (*models.Tournament, error) {
	tournament := new(models.Tournament)
	err := models.DB.QueryRow(`
		SELECT t.id, t.name, t.short_name, t.start_time, t.end_time
		FROM tournament t WHERE t.id = ?`, id).Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.StartTime, &tournament.EndTime)
	if err != nil {
		return nil, err
	}
	return tournament, nil
}

// GetTournamentStatistics will return statistics for a given tournament
func GetTournamentStatistics(id int) (map[int][]*models.TournamentStatistics, error) {
	rows, err := models.DB.Query(`
		SELECT
			t.id, t.name, t.short_name, t.start_time, t.end_time,
			tg.id, tg.name, tg.division,
			p.id as 'player_id',
			COUNT(DISTINCT m.id) AS 'p',
			COUNT(DISTINCT won.id) AS 'w',
			COUNT(DISTINCT draw.id) AS 'd',
			COUNT(DISTINCT lost.id) AS 'l',
			COUNT(DISTINCT legs_for.id) AS 'F',
			COUNT(DISTINCT legs_against.id) AS 'A',
			(COUNT(DISTINCT legs_for.id) - COUNT(DISTINCT legs_against.id)) AS 'diff',
			COUNT(DISTINCT won.id) * 2 + COUNT(DISTINCT draw.id) AS 'pts'
		FROM player2leg p2l
			JOIN matches m ON m.id = p2l.match_id
			JOIN player p ON p.id = p2l.player_id
			LEFT JOIN matches won ON won.id = p2l.match_id AND won.winner_id = p.id
			LEFT JOIN matches lost ON lost.id = p2l.match_id AND lost.winner_id <> p.id
			LEFT JOIN matches draw ON draw.id = p2l.match_id AND draw.winner_id IS NULL
			LEFT JOIN leg legs_for ON legs_for.id = p2l.leg_id AND legs_for.winner_id = p.id
			LEFT JOIN leg legs_against ON legs_against.id = p2l.leg_id AND legs_against.winner_id <> p.id
			JOIN tournament t ON t.id = m.tournament_id
			JOIN player2tournament p2t ON p2t.player_id = p.id AND p2t.tournament_id = t.id
			JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
		WHERE m.tournament_id = ?
		GROUP BY p2l.player_id
		ORDER BY tg.id, pts DESC, diff DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	statistics := make(map[int][]*models.TournamentStatistics, 0)
	for rows.Next() {
		tournament := new(models.Tournament)
		group := new(models.TournamentGroup)
		stats := new(models.TournamentStatistics)
		err := rows.Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.StartTime, &tournament.EndTime, &group.ID, &group.Name, &group.Division,
			&stats.PlayerID, &stats.Played, &stats.MatchesWon, &stats.MatchesDraw, &stats.MatchesLost, &stats.LegsFor, &stats.LegsAgainst, &stats.LegsDifference, &stats.Points)
		if err != nil {
			return nil, err
		}
		stats.Tournament = tournament
		stats.Group = group

		if _, ok := statistics[group.ID]; !ok {
			statistics[group.ID] = make([]*models.TournamentStatistics, 0)
		}
		statistics[group.ID] = append(statistics[group.ID], stats)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return statistics, nil
}
