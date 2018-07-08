package models

import "github.com/guregu/null"

// Tournament struct for storing tournaments
type Tournament struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	ShortName            string                `json:"short_name"`
	IsFinished           bool                  `json:"is_finished"`
	IsPlayoffs           bool                  `json:"is_playoffs"`
	PlayoffsTournamentID null.Int              `json:"playoffs_tournament_id,omitempty"`
	PlayoffsTournament   *Tournament           `json:"playoffs,omitempty"`
	StartTime            null.String           `json:"start_time"`
	EndTime              null.String           `json:"end_time"`
	Groups               []*TournamentGroup    `json:"groups,omitempty"`
	Standings            []*TournamentStanding `json:"standings,omitempty"`
}

// TournamentGroup struct for storing tournament groups
type TournamentGroup struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Division null.Int `json:"division,omitempty"`
}

// Player2Tournament struct for storing player to tounament links
type Player2Tournament struct {
	PlayerID          int
	TournamentID      int
	TournamentGroupID int
	IsPromoted        bool
	IsRelegated       bool
	IsWinner          bool
}

// TournamentStanding struct for stroring final tournament standings
type TournamentStanding struct {
	TournamentID   int    `json:"tournament_id"`
	TournamentName string `json:"tournament_name"`
	PlayerID       int    `json:"player_id"`
	PlayerName     string `json:"player_name"`
	Rank           int    `json:"rank"`
	Elo            int    `json:"elo"`
	EloPlayed      int    `json:"elo_matches"`
}

// PlayerTournamentStanding struct for storing player tournament standings
type PlayerTournamentStanding struct {
	PlayerID        int              `json:"player_id"`
	FinalStanding   int              `json:"final_standing"`
	TotalPlayers    int              `json:"total_players"`
	Tournament      *Tournament      `json:"tournament"`
	TournamentGroup *TournamentGroup `json:"tournament_group"`
	Elo             int              `json:"elo,omitempty"`
}
