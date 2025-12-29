package repository

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/rinarudhei/pomcli/model"
	"github.com/rinarudhei/pomcli/session"
)

const (
	createTableInterval string = `CREATE TABLE IF NOT EXISTS "pomodoros" (
		"id" INTEGER,
		"start_time" DATETIME NOT NULL,
		"planned_duration" INTEGER DEFAULT 0,
		"actual_duration" INTEGER DEFAULT 0,
		"type" TEXT NOT NULL,
		"state" INTEGER DEFAULT 1,
		PRIMARY KEY("id")
	);`
)

type DbRepo struct {
	db *sql.DB
	sync.RWMutex
}

func NewSQLiteRepo(dbfile string) (*DbRepo, error) {
	db, err := sql.Open("sqlite", dbfile)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableInterval); err != nil {
		return nil, err
	}

	return &DbRepo{
		db: db,
	}, nil
}

func (r *DbRepo) Close() error {
	return r.db.Close()
}

func (r *DbRepo) Add(i model.Pomodoro) error {
	r.Lock()
	defer r.Unlock()

	insStmt, err := r.db.Prepare("INSERT INTO pomodoros VALUES(NULL, ?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer insStmt.Close()

	_, err = insStmt.Exec(i.Starttime, i.PlannedDuration, i.ActualDuration, i.Type, i.State)
	if err != nil {
		return err
	}

	return nil
}

func (r *DbRepo) GetTodayPomodoro() ([]model.Pomodoro, error) {
	r.RLock()
	defer r.RUnlock()

	rows, err := r.db.Query("SELECT * FROM pomodoros where start_time >= ? AND type = ? ORDER BY start_time DESC", time.Now().Format("2006-01-02"), session.PomodoroSession)
	if err != nil {
		return []model.Pomodoro{}, err
	}

	var pomodoros []model.Pomodoro
	for rows.Next() {
		i := model.Pomodoro{}
		if err := rows.Scan(&i.ID, &i.Starttime, &i.PlannedDuration, &i.ActualDuration, &i.Type, &i.State); err != nil {
			return []model.Pomodoro{}, err
		}

		pomodoros = append(pomodoros, i)
	}

	return pomodoros, nil
}
