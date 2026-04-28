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

func GetStatsByUserAndFriends(userID int, days int) ([]domain.UserStats, error) {
	friends, err := GetFriends(userID)
	if err != nil {
		return nil, err
	}
	var userIDs []int
	userIDs = append(userIDs, userID)
	for _, friend := range friends {
		userIDs = append(userIDs, friend.FriendID)
	}
	if len(userIDs) == 0 {
		return []domain.UserStats{}, nil
	}
	placeholders := ""
	var args []interface{}
	for i, id := range userIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		args = append(args, id)
	}
	dateFilter := ""
	if days > 0 {
		dateFilter = " AND p.created_at >= date('now', '-" + strconv.Itoa(days) + " days') "
	}
	query := `
	SELECT
		u.username,
		date(p.created_at) as day,
		SUM(p.count) as total
	FROM users u
	JOIN pushups p ON u.id = p.user_id
	WHERE u.id IN (` + placeholders + `) AND u.is_approved = 1` + dateFilter + `
	GROUP BY u.id, day
	ORDER BY day ASC
	`
	rows, err := DB.Query(query, args...)
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

/*
// public graph, unsafe
func GetStats(days int) ([]domain.UserStats, error) {
	whereFilter := " WHERE u.is_approved = 1 "
	if days > 0 {
		whereFilter += " AND p.created_at >= date('now', '-" + strconv.Itoa(days) + " days') "
	}
	query := `
	SELECT
		u.username,
		date(p.created_at) as day,
		SUM(p.count) as total
	FROM users u
	JOIN pushups p ON u.id = p.user_id
	` + whereFilter + `
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
*/
