package db

import (
	"crypto/md5"
	"fmt"
	"punkpushups/domain"
	"strconv"
)

func AddPushups(userID int, count int) error {
	_, err := DB.Exec("INSERT INTO pushups (user_id, count) VALUES (?, ?)", userID, count)
	return err
}

func GetStats(days int) ([]domain.UserStats, error) {
	timeFilter := ""
	if days > 0 {
		timeFilter = "WHERE p.created_at >= date('now', '-" + strconv.Itoa(days) + " days')"
	}
	query := `
	SELECT
		u.username,
		date(p.created_at) as day,
		SUM(p.count) as total
	FROM users u
	JOIN pushups p ON u.id = p.user_id
	` + timeFilter + `
	GROUP BY u.id, day
	ORDER BY day ASC
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	statsMap := make(map[string]*domain.UserStats)
	for rows.Next() {
		var username, date string
		var total int
		if err := rows.Scan(&username, &date, &total); err != nil {
			return nil, err
		}
		if _, ok := statsMap[username]; !ok {
			statsMap[username] = &domain.UserStats{
				Username: username,
				Color:    "#" + fmt.Sprintf("%x", md5.Sum([]byte(username)))[0:6],
				Points:   []domain.ChartPoint{},
			}
		}
		statsMap[username].Points = append(statsMap[username].Points, domain.ChartPoint{
			Date:  date,
			Value: total,
		})
	}
	var result []domain.UserStats
	for _, val := range statsMap {
		result = append(result, *val)
	}
	return result, nil
}
