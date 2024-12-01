package models

import "github.com/guregu/null"

type TournamentMatchTemplate struct {
	Home int
	Away int
}

var TournamentTemplateSemiFinals = [2]TournamentMatchTemplate{
	{Home: 0, Away: 1},
	{Home: 1, Away: 0},
}
var TournamentTemplateSemiFinalsSingle = [2]TournamentMatchTemplate{
	{Home: 0, Away: 3},
	{Home: 1, Away: 2},
}
var TournamentTemplateQuarterFinals = [4]TournamentMatchTemplate{
	{Home: 0, Away: 3},
	{Home: 2, Away: 1},
	{Home: 1, Away: 2},
	{Home: 3, Away: 0},
}
var TournamentTemplateQuarterFinalsSingle = [4]TournamentMatchTemplate{
	{Home: 0, Away: 7},
	{Home: 3, Away: 4},
	{Home: 1, Away: 6},
	{Home: 2, Away: 5},
}
var TournamentTemplateLast16 = [8]TournamentMatchTemplate{
	{Home: 0, Away: 7},
	{Home: 4, Away: 2},
	{Home: 3, Away: 5},
	{Home: 6, Away: 1},
	{Home: 1, Away: 6},
	{Home: 5, Away: 3},
	{Home: 2, Away: 4},
	{Home: 7, Away: 0},
}

// Tournament struct for storing tournaments
type Tournament struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	ShortName            string                `json:"short_name"`
	IsFinished           bool                  `json:"is_finished"`
	IsPlayoffs           bool                  `json:"is_playoffs"`
	PlayoffsTournamentID null.Int              `json:"playoffs_tournament_id,omitempty"`
	PlayoffsTournament   *Tournament           `json:"playoffs,omitempty"`
	PresetID             null.Int              `json:"preset_id,omitempty"`
	Preset               *TournamentPreset     `json:"preset,omitempty"`
	ManualAdmin          bool                  `json:"manual_admin"`
	OfficeID             int                   `json:"office_id"`
	StartTime            null.Time             `json:"start_time"`
	EndTime              null.Time             `json:"end_time"`
	Groups               []*TournamentGroup    `json:"groups,omitempty"`
	Standings            []*TournamentStanding `json:"standings,omitempty"`
	Players              []*Player2Tournament  `json:"players,omitempty"`
}

// TournamentGroup struct for storing tournament groups
type TournamentGroup struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	IsGenerated bool     `json:"is_generated"`
	IsPlayoffs  bool     `json:"is_playoffs"`
	Division    null.Int `json:"division,omitempty"`
}

// Player2Tournament struct for storing player to tounament links
type Player2Tournament struct {
	PlayerID          int  `json:"player_id"`
	TournamentID      int  `json:"tournament_id"`
	TournamentGroupID int  `json:"tournament_group_id"`
	IsPromoted        bool `json:"is_promoted"`
	IsRelegated       bool `json:"is_relegated"`
	IsWinner          bool `json:"is_winner"`
}

// TournamentStanding struct for stroring final tournament standings
type TournamentStanding struct {
	TournamentID     int    `json:"tournament_id"`
	TournamentName   string `json:"tournament_name"`
	PlayerID         int    `json:"player_id"`
	PlayerName       string `json:"player_name"`
	Rank             int    `json:"rank"`
	Elo              int    `json:"elo"`
	EloPlayed        int    `json:"elo_matches"`
	CurrentElo       int    `json:"current_elo"`
	CurrentEloPlayed int    `json:"current_elo_matches"`
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

type TournamentProbabilities struct {
}

type TournamentPreset struct {
	ID                      int              `json:"id"`
	Name                    string           `json:"name"`
	MatchType               *MatchType       `json:"match_type_id"`
	StartingScore           int              `json:"starting_score"`
	MatchMode               *MatchMode       `json:"match_mode"`
	MatchModeLast16         *MatchMode       `json:"match_mode_last_16"`
	MatchModeQuarterFinal   *MatchMode       `json:"match_mode_quarter_final"`
	MatchModeSemiFinal      *MatchMode       `json:"match_mode_semi_final"`
	MatchModeGrandFinal     *MatchMode       `json:"match_mode_grand_final"`
	PlayoffsTournamentGroup *TournamentGroup `json:"playoffs_tournament_group"`
	Group1TournamentGroup   *TournamentGroup `json:"group1_tournament_group"`
	Group2TournamentGroup   *TournamentGroup `json:"group2_tournament_group"`
	PlayerIDWalkover        int              `json:"player_id_walkover"`
	PlayerIDPlaceholderHome int              `json:"player_id_placeholder_home"`
	PlayerIDPlaceholderAway int              `json:"player_id_placeholder_away"`
	Description             null.String      `json:"description"`
}

// GenerateTournamentInput struct for storing generate tournament inputs
type GenerateTournamentInput struct {
	Name          string               `json:"name"`
	ShortName     string               `json:"short_name"`
	IsPlayoffs    bool                 `json:"is_playoffs"`
	ManualAdmin   bool                 `json:"manual_admin"`
	OfficeID      int                  `json:"office_id"`
	MatchModeID   int                  `json:"match_mode_id"`
	MatchTypeID   int                  `json:"match_type_id"`
	StartingScore int                  `json:"starting_score"`
	MaxRounds     int                  `json:"max_rounds"`
	Players       []*Player2Tournament `json:"players,omitempty"`
}

// GeneratePlayoffsInput struct for storing generate playoffs inputs
type GeneratePlayoffsInput struct {
	MatchModeLast32ID int `json:"match_mode_last32"`
	MatchModeLast16ID int `json:"match_mode_last16"`
	MatchModeQFID     int `json:"match_mode_quarterFinals"`
	MatchModeSFID     int `json:"match_mode_semiFinals"`
	MatchModeGFID     int `json:"match_mode_grandFinals"`
}
