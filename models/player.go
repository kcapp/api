package models

import (
	"encoding/json"
	"strings"

	"github.com/guregu/null"
)

// Player struct used for storing players
type Player struct {
	ID            int         `json:"id"`
	FirstName     string      `json:"first_name"`
	LastName      null.String `json:"last_name"`
	VocalName     null.String `json:"vocal_name,omitempty"`
	Nickname      null.String `json:"nickname,omitempty"`
	SlackHandle   null.String `json:"slack_handle,omitempty"`
	MatchesPlayed int         `json:"matches_played"`
	MatchesWon    int         `json:"matches_won"`
	LegsPlayed    int         `json:"legs_played"`
	LegsWon       int         `json:"legs_won"`
	Color         null.String `json:"color,omitempty"`
	ProfilePicURL null.String `json:"profile_pic_url,omitempty"`
	OfficeID      null.Int    `json:"office_id,omitempty"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at,omitempty"`
	TournamentElo int         `json:"tournament_elo,omitempty"`
	CurrentElo    int         `json:"current_elo,omitempty"`
}

// MarshalJSON will marshall the given object to JSON
func (player Player) MarshalJSON() ([]byte, error) {
	// Use a type to get consistnt order of JSON key-value pairs.
	type playerJSON struct {
		ID            int         `json:"id"`
		Name          string      `json:"name"`
		FirstName     string      `json:"first_name"`
		LastName      null.String `json:"last_name"`
		VocalName     null.String `json:"vocal_name,omitempty"`
		Nickname      null.String `json:"nickname,omitempty"`
		SlackHandle   null.String `json:"slack_handle,omitempty"`
		MatchesPlayed int         `json:"matches_played"`
		MatchesWon    int         `json:"matches_won"`
		LegsPlayed    int         `json:"legs_played"`
		LegsWon       int         `json:"legs_won"`
		Color         null.String `json:"color,omitempty"`
		ProfilePicURL null.String `json:"profile_pic_url,omitempty"`
		OfficeID      null.Int    `json:"office_id,omitempty"`
		CreatedAt     string      `json:"created_at"`
		UpdatedAt     string      `json:"updated_at,omitempty"`
		TournamentElo int         `json:"tournament_elo,omitempty"`
		CurrentElo    int         `json:"current_elo,omitempty"`
	}

	return json.Marshal(playerJSON{
		ID:            player.ID,
		FirstName:     player.FirstName,
		LastName:      player.LastName,
		VocalName:     player.VocalName,
		Nickname:      player.Nickname,
		SlackHandle:   player.SlackHandle,
		MatchesPlayed: player.MatchesPlayed,
		MatchesWon:    player.MatchesWon,
		LegsPlayed:    player.LegsPlayed,
		LegsWon:       player.LegsWon,
		Color:         player.Color,
		ProfilePicURL: player.ProfilePicURL,
		OfficeID:      player.OfficeID,
		CreatedAt:     player.CreatedAt,
		UpdatedAt:     player.UpdatedAt,
		TournamentElo: player.TournamentElo,
		CurrentElo:    player.CurrentElo,
		Name:          strings.Trim((player.FirstName + " " + player.LastName.ValueOrZero()), " "),
	})
}
