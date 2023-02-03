package queries

func QueryBestStatistics() string {
	return `
		SELECT
			p.id,
			l.winner_id,
			l.id,
			(s.ppd_score * 3) / s.darts_thrown,
			((s.first_nine_ppd_score) * 3 / if(s.darts_thrown < 9, s.darts_thrown, 9)),
			s.checkout_percentage,
			s.darts_thrown,
			l.starting_score
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
		WHERE s.player_id IN (?)
			AND l.starting_score IN (?)`
}

func QueryHighestCheckout() string {
	return `
		SELECT
			player_id,
			leg_id,
			MAX(checkout)
		FROM (SELECT
				s.player_id,
				s.leg_id,
				IFNULL(s.first_dart * s.first_dart_multiplier, 0) +
					IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
					IFNULL(s.third_dart * s.third_dart_multiplier, 0) AS 'checkout'
			FROM score s
			JOIN leg l ON l.id = s.leg_id
			WHERE l.winner_id = s.player_id
				AND s.player_id IN (?)
				AND s.id IN (SELECT MAX(s.id) FROM score s JOIN leg l ON l.id = s.leg_id WHERE l.winner_id = s.player_id GROUP BY leg_id)
				AND l.starting_score IN (?)
			GROUP BY s.player_id, s.id
			ORDER BY checkout DESC) checkouts
		GROUP BY player_id`
}

func baseQueryPlayersX01Statistics() string {
	return `SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(p.id) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
			SUM(s.first_nine_ppd) / COUNT(p.id) * 3 AS 'first_nine_three_dart_avg',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
			SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
			SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage'
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l2.match_id AND l2.winner_id = p.id
		WHERE s.player_id IN (?)
			AND l.starting_score IN (301, 501, 701)
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND m.match_type_id = 1`
}

func QueryPlayerCurrentX01Statistics() string {
	return baseQueryPlayersX01Statistics()
}

func QueryPlayerPreviousX01Statistics() string {
	return baseQueryPlayersX01Statistics() +
		` AND m.updated_at < (CURRENT_DATE - INTERVAL WEEKDAY(CURRENT_DATE) DAY)`
}

func QueryPlayersX01Statistics() string {
	return baseQueryPlayersX01Statistics() +
		` UNION ALL ` + baseQueryPlayersX01Statistics() +
		` AND m.updated_at < (CURRENT_DATE - INTERVAL WEEKDAY(CURRENT_DATE) DAY)`
}

func QueryPreviousPlayersX01Statistics() string {
	return `
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(p.id) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
			SUM(s.first_nine_ppd) / COUNT(p.id) * 3 AS 'first_nine_three_dart_avg',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
			SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
			SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage'
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l2.match_id AND l2.winner_id = p.id
		WHERE s.player_id IN (?)
			AND l.starting_score IN (?)
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND m.match_type_id = 1
			-- Exclude all matches played this week
			AND m.updated_at < (CURRENT_DATE - INTERVAL WEEKDAY(CURRENT_DATE) DAY)
		GROUP BY s.player_id
		ORDER BY p.id`
}
