package models

import "github.com/guregu/null"

// Tournament struct for storing tournaments
type Tournament struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	ShortName string             `json:"short_name"`
	StartTime null.String        `json:"start_time"`
	EndTime   null.String        `json:"end_time"`
	Groups    []*TournamentGroup `json:"groups,omitempty"`
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
}
