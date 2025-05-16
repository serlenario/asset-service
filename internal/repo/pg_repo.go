package repo

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	GetByLogin(ctx context.Context, login string) (id int64, hash string, err error)
}

type SessionRepo interface {
	DeleteByUser(ctx context.Context, uid int64) error
	Create(ctx context.Context, uid int64, ip net.IP, expiresAt time.Time) (string, error)
	Validate(ctx context.Context, token string) (uid int64, expiresAt time.Time, err error)
}

type AssetRepo interface {
	Upsert(ctx context.Context, uid int64, name string, data []byte) error
	Get(ctx context.Context, uid int64, name string) ([]byte, error)
	List(ctx context.Context, uid int64) ([]string, error)
	Delete(ctx context.Context, uid int64, name string) error
	Exists(ctx context.Context, name string) (bool, error)
}

type PG struct{ pool *pgxpool.Pool }

func NewPG(pool *pgxpool.Pool) *PG { return &PG{pool: pool} }

func (p *PG) GetByLogin(ctx context.Context, login string) (int64, string, error) {
	var id int64
	var hash string
	err := p.pool.QueryRow(ctx, `SELECT id, password_hash FROM users WHERE login=$1`, login).Scan(&id, &hash)
	if err != nil {
		log.Printf("repo.GetByLogin error: %v login=%s", err, login)
	}
	return id, hash, err
}

func (p *PG) DeleteByUser(ctx context.Context, uid int64) error {
	_, err := p.pool.Exec(ctx, `DELETE FROM sessions WHERE uid=$1`, uid)
	if err != nil {
		log.Printf("repo.DeleteByUser error: %v uid=%d", err, uid)
	}
	return err
}

func (p *PG) Create(ctx context.Context, uid int64, ip net.IP, expiresAt time.Time) (string, error) {
	var sid string
	err := p.pool.QueryRow(ctx,
		`INSERT INTO sessions(uid, ip_addr, expires_at) VALUES($1, $2, $3) RETURNING id`,
		uid, ip, expiresAt).Scan(&sid)
	if err != nil {
		log.Printf("repo.CreateSession error: %v uid=%d", err, uid)
	}
	return sid, err
}

func (p *PG) Validate(ctx context.Context, token string) (int64, time.Time, error) {
	var uid int64
	var exp time.Time
	err := p.pool.QueryRow(ctx,
		`SELECT uid, expires_at FROM sessions WHERE id=$1`, token).Scan(&uid, &exp)
	if err != nil {
		log.Printf("repo.ValidateSession error: %v token=%s", err, token)
	}
	return uid, exp, err
}

func (p *PG) Upsert(ctx context.Context, uid int64, name string, data []byte) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO assets(name, uid, data) VALUES($1, $2, $3)
         ON CONFLICT(name, uid) DO UPDATE SET data=EXCLUDED.data, created_at=now()`,
		name, uid, data)
	if err != nil {
		log.Printf("repo.UpsertAsset error: %v uid=%d name=%s", err, uid, name)
	}
	return err
}

func (p *PG) Get(ctx context.Context, uid int64, name string) ([]byte, error) {
	var data []byte
	err := p.pool.QueryRow(ctx,
		`SELECT data FROM assets WHERE uid=$1 AND name=$2`, uid, name).Scan(&data)
	if err != nil {
		log.Printf("repo.GetAsset error: %v uid=%d name=%s", err, uid, name)
	}
	return data, err
}

func (p *PG) List(ctx context.Context, uid int64) ([]string, error) {
	rows, err := p.pool.Query(ctx, `SELECT name FROM assets WHERE uid=$1 ORDER BY created_at DESC`, uid)
	if err != nil {
		log.Printf("repo.ListAssets error: %v uid=%d", err, uid)
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			log.Printf("repo.ListAssets scan error: %v", err)
			return nil, err
		}
		names = append(names, n)
	}
	return names, rows.Err()
}

func (p *PG) Delete(ctx context.Context, uid int64, name string) error {
	_, err := p.pool.Exec(ctx, `DELETE FROM assets WHERE uid=$1 AND name=$2`, uid, name)
	if err != nil {
		log.Printf("repo.DeleteAsset error: %v uid=%d name=%s", err, uid, name)
	}
	return err
}

func (p *PG) Exists(ctx context.Context, name string) (bool, error) {
	var ex bool
	err := p.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM assets WHERE name=$1)`, name).Scan(&ex)
	if err != nil {
		log.Printf("repo.ExistsAsset error: %v name=%s", err, name)
	}
	return ex, err
}
