package repository

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/djcrock/fwip/internal/model"
	"io"
	"io/fs"
	"log"
	"math"
	"strconv"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

//go:embed migrations/*.sql
var migrations embed.FS

var (
	ErrNoSuchTitle   = errors.New("title does not exist")
	ErrNoSuchService = errors.New("service does not exist")
	ErrNoSuchUser    = errors.New("user does not exist")
)

// A Repository is a Repository stores persisted state in an SQLite database.
type Repository struct {
	conn *sqlite.Conn
	// TODO: should I get rid of this?
	isClosed bool
}

// A Pool is a Pool that contains Repository instances of Repository.
type Pool struct {
	pool *sqlitex.Pool
}

// TODO: decide on a range of valid IDs
const (
	minId = 1
	maxId = math.MaxInt64
)

func NewPool(pool *sqlitex.Pool) *Pool {
	// Initialize the database
	p := &Pool{
		pool: pool,
	}
	repo := p.GetRepository(context.Background())
	defer p.PutRepository(repo)

	stmt := repo.conn.Prep("PRAGMA foreign_keys = ON;")
	stmt.Step()
	stmt.Finalize()

	err := repo.applyMigrations()
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func (p *Pool) GetRepository(ctx context.Context) *Repository {
	// TODO: this CAN return nil (e.g. if ctx expires), so figure out a way to handle it
	conn := p.pool.Get(ctx)
	return &Repository{conn: conn, isClosed: false}
}

func (p *Pool) PutRepository(repo *Repository) {
	repo.isClosed = true
	p.pool.Put(repo.conn)
}

func (r *Repository) Transact() (completeFn func(*error)) {
	return sqlitex.Save(r.conn)
}

func (r *Repository) GetTitles() ([]*model.Title, error) {
	stmt := r.conn.Prep(`
SELECT id, imdb_id, type, name, year, release_date, runtime
FROM title
;`,
	)
	defer stmt.Reset()

	titles := make([]*model.Title, 0)
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to retrieve titles: %w", err)
		} else if !hasRow {
			break
		}
		titles = append(titles, &model.Title{
			Id:          stmt.GetInt64("id"),
			ImdbId:      stmt.GetText("imdb_id"),
			Type:        stmt.GetText("type"),
			Name:        stmt.GetText("name"),
			Year:        stmt.GetInt64("year"),
			ReleaseDate: stmt.GetText("release_date"),
			Runtime:     stmt.GetInt64("runtime"),
		})
	}

	return titles, nil
}

func (r *Repository) GetTitlesByService(serviceId int64) ([]*model.Title, error) {
	stmt := r.conn.Prep(`
SELECT t.id, t.imdb_id, t.type, t.name, t.year, t.release_date, t.runtime
FROM title t
INNER JOIN main.service_title st on t.id = st.title_id
WHERE st.service_id = $serviceId
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$serviceId", serviceId)

	titles := make([]*model.Title, 0)
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to retrieve titles: %w", err)
		} else if !hasRow {
			break
		}
		titles = append(titles, &model.Title{
			Id:          stmt.GetInt64("id"),
			ImdbId:      stmt.GetText("imdb_id"),
			Type:        stmt.GetText("type"),
			Name:        stmt.GetText("name"),
			Year:        stmt.GetInt64("year"),
			ReleaseDate: stmt.GetText("release_date"),
			Runtime:     stmt.GetInt64("runtime"),
		})
	}

	return titles, nil
}

func (r *Repository) GetTitle(titleId int64) (*model.Title, error) {
	stmt := r.conn.Prep(`
SELECT id, imdb_id, type, name, year, release_date, runtime
FROM title
WHERE id = $id
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$id", titleId)
	if hasRow, err := stmt.Step(); err != nil {
		return nil, fmt.Errorf("failed to retrieve title %d: %w", titleId, err)
	} else if !hasRow {
		return nil, ErrNoSuchTitle
	}

	title := &model.Title{
		Id:          stmt.GetInt64("id"),
		ImdbId:      stmt.GetText("imdb_id"),
		Type:        stmt.GetText("type"),
		Name:        stmt.GetText("name"),
		Year:        stmt.GetInt64("year"),
		ReleaseDate: stmt.GetText("release_date"),
		Runtime:     stmt.GetInt64("runtime"),
	}

	return title, nil
}

func (r *Repository) PutTitle(title *model.Title) (titleId int64, err error) {
	defer sqlitex.Save(r.conn)(&err)
	if title.Id == model.NoId {
		titleId, err = r.insertTitle(title)
		return
	}
	titleId = title.Id
	err = r.updateTitle(title)
	return
}

func (r *Repository) insertTitle(title *model.Title) (int64, error) {
	stmt := r.conn.Prep(`
INSERT INTO title (
	id,
	imdb_id,
	type,
    name,
    year,
	release_date,
	runtime
)
VALUES (
	$id,
    $imdbId,
	$type,
    $name,
    $year,
	$releaseDate,
	$runtime
)
;`,
	)
	defer stmt.Reset()
	stmt.SetText("$imdbId", title.ImdbId)
	stmt.SetText("$type", title.Type)
	stmt.SetText("$name", title.Name)
	stmt.SetInt64("$year", title.Year)
	stmt.SetText("$releaseDate", title.ReleaseDate)
	stmt.SetInt64("$runtime", title.Runtime)
	id, err := sqlitex.InsertRandID(stmt, "$id", minId, maxId)

	if err != nil {
		return model.NoId, err
	}
	title.Id = id

	return id, nil
}

func (r *Repository) updateTitle(title *model.Title) error {
	stmt := r.conn.Prep(`
UPDATE title
SET name = $name
WHERE id = $id
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$id", title.Id)
	stmt.SetText("$name", title.Name)
	_, err := stmt.Step()
	return err
}

func (r *Repository) GetServices() ([]*model.Service, error) {
	stmt := r.conn.Prep(`
SELECT id, name
FROM service
;`,
	)
	defer stmt.Reset()

	titles := make([]*model.Service, 0)
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to retrieve services: %w", err)
		} else if !hasRow {
			break
		}
		titles = append(titles, &model.Service{
			Id:   stmt.GetInt64("id"),
			Name: stmt.GetText("name"),
		})
	}

	return titles, nil
}

func (r *Repository) GetService(id int64) (*model.Service, error) {
	stmt := r.conn.Prep(`
SELECT id, name
FROM service
WHERE id = $id
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$id", id)
	if hasRow, err := stmt.Step(); err != nil {
		return nil, fmt.Errorf("failed to retrieve service %d: %w", id, err)
	} else if !hasRow {
		return nil, ErrNoSuchService
	}

	service := &model.Service{
		Id:   stmt.GetInt64("id"),
		Name: stmt.GetText("name"),
	}

	return service, nil
}

func (r *Repository) GetUsers() ([]*model.User, error) {
	stmt := r.conn.Prep(`
SELECT id, username
FROM user
;`,
	)
	defer stmt.Reset()

	titles := make([]*model.User, 0)
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to retrieve users: %w", err)
		} else if !hasRow {
			break
		}
		titles = append(titles, &model.User{
			Id:       stmt.GetInt64("id"),
			Username: stmt.GetText("username"),
		})
	}

	return titles, nil
}

func (r *Repository) GetUser(id int64) (*model.User, error) {
	stmt := r.conn.Prep(`
SELECT id, username
FROM user
WHERE id = $id
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$id", id)
	if hasRow, err := stmt.Step(); err != nil {
		return nil, fmt.Errorf("failed to retrieve user %d: %w", id, err)
	} else if !hasRow {
		return nil, ErrNoSuchUser
	}

	user := &model.User{
		Id:       stmt.GetInt64("id"),
		Username: stmt.GetText("username"),
	}

	return user, nil
}

func (r *Repository) PutUser(user *model.User) (userId int64, err error) {
	defer sqlitex.Save(r.conn)(&err)
	if user.Id == model.NoId {
		userId, err = r.insertUser(user)
		return
	}
	userId = user.Id
	err = r.updateUser(user)
	return
}

func (r *Repository) insertUser(user *model.User) (int64, error) {
	stmt := r.conn.Prep(`
INSERT INTO user (
	id,
    username
)
VALUES (
	$id,
    $username
)
;`,
	)
	defer stmt.Reset()
	stmt.SetText("$username", user.Username)
	id, err := sqlitex.InsertRandID(stmt, "$id", minId, maxId)

	if err != nil {
		return model.NoId, err
	}
	user.Id = id

	return id, nil
}

func (r *Repository) updateUser(user *model.User) error {
	stmt := r.conn.Prep(`
UPDATE user
SET username = $username
WHERE id = $id
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$id", user.Id)
	stmt.SetText("$username", user.Username)
	_, err := stmt.Step()
	return err
}

func (r *Repository) GetUserWatchHistory(userId int64) ([]*model.WatchHistory, error) {
	stmt := r.conn.Prep(`
SELECT user_id, title_id, watched, want_to_watch
FROM watch_history
WHERE user_id = $userId
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$userId", userId)

	watchHistory := make([]*model.WatchHistory, 0)
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to retrieve watch history: %w", err)
		} else if !hasRow {
			break
		}
		watchHistory = append(watchHistory, &model.WatchHistory{
			UserId:      stmt.GetInt64("user_id"),
			TitleId:     stmt.GetInt64("title_id"),
			Watched:     stmt.GetBool("watched"),
			WantToWatch: stmt.GetInt64("want_to_watch"),
		})
	}

	return watchHistory, nil
}

func (r *Repository) PutWatchHistory(watchHistory *model.WatchHistory) error {
	stmt := r.conn.Prep(`
INSERT INTO watch_history (
	user_id,
    title_id,
	watched,
	want_to_watch
)
VALUES (
	$userId,
    $titleId,
	$watched,
	$wantToWatch
)
ON CONFLICT DO UPDATE SET
	watched = excluded.watched,
	want_to_watch = excluded.want_to_watch
;`,
	)
	defer stmt.Reset()
	stmt.SetInt64("$userId", watchHistory.UserId)
	stmt.SetInt64("$titleId", watchHistory.TitleId)
	stmt.SetBool("$watched", watchHistory.Watched)
	stmt.SetInt64("$wantToWatch", watchHistory.WantToWatch)

	_, err := stmt.Step()

	return err
}

// applyMigrations walks through the scripts in the migration directory,
// running any that haven't yet been applied.
func (r *Repository) applyMigrations() (err error) {
	defer sqlitex.Save(r.conn)(&err)
	var version int64
	version, err = getSchemaVersion(r.conn)
	if err != nil {
		return
	}
	log.Println("current schema version is", version)
	err = fs.WalkDir(
		migrations,
		"migrations",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			// Migrations have a filename like 0001-initial-schema.sql
			migrationId, err := strconv.ParseInt(d.Name()[:4], 10, 64)
			if err != nil {
				return err
			}
			if migrationId > version {
				log.Println("applying schema migration", path)
				return applyMigration(r.conn, path, migrationId)
			}
			return nil
		},
	)
	log.Println("finished applying schema migrations")

	return
}

// getSchemaVersion gets the ID of the most recent migration to have run on the database.
func getSchemaVersion(conn *sqlite.Conn) (version int64, err error) {
	stmt, _, err := conn.PrepareTransient("PRAGMA user_version;")
	if err != nil {
		return
	}
	version, err = sqlitex.ResultInt64(stmt)
	err = stmt.Finalize()
	return
}

// applyMigration runs the given migration script on the database.
func applyMigration(conn *sqlite.Conn, path string, migrationId int64) (err error) {
	migrationFile, err := migrations.Open(path)
	if err != nil {
		return
	}
	defer migrationFile.Close()
	migrationScript, err := io.ReadAll(migrationFile)
	if err != nil {
		return
	}
	setVersion := "\n\nPRAGMA user_version=" + strconv.FormatInt(migrationId, 10) + ";"
	err = sqlitex.ExecScript(conn, string(migrationScript)+setVersion)
	if err != nil {
		return
	}
	// Keep track of which migrations have been applied. Can't use parameters for PRAGMA statements.
	return
}
