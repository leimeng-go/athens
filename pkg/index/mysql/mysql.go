package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/leimeng-go/athens/pkg/config"
	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/index"
)

//Mysql 创建go模块的的索引，这个是表格初始化
// New returns a new Indexer with a MySQL implementation.
// It attempts to connect to the DB and create the index table
// if it doesn ot already exist.
func New(cfg *config.MySQL) (index.Indexer, error) {
	dataSource := getMySQLSource(cfg)
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	return &indexer{db}, nil
}

const schema = `
	CREATE TABLE IF NOT EXISTS indexes(
	id INT
		AUTO_INCREMENT
		PRIMARY KEY
		COMMENT 'Unique identifier for a module line',

	path VARCHAR(255)
		NOT NULL
		COMMENT 'Import path of the module',

	version VARCHAR(255)
		NOT NULL
		COMMENT 'Module version',

	timestamp TIMESTAMP(6)
		COMMENT 'Date and time when the module was first created',

	INDEX (timestamp),
	UNIQUE INDEX idx_module_version (path, version)
	) CHARACTER SET utf8;
`

type indexer struct {
	db *sql.DB
}

//Index 创建go模块的索引
func (i *indexer) Index(ctx context.Context, mod, ver string) error {
	const op errors.Op = "mysql.Index"
	_, err := i.db.ExecContext(
		ctx,
		`INSERT INTO indexes (path, version, timestamp) VALUES (?, ?, ?)`,
		mod,
		ver,
		time.Now().Format("2006-01-02 15:04:05.000"),
	)
	if err != nil {
		fmt.Printf("sql: %s\n",fmt.Sprintf("INSERT INTO indexes (path, version, timestamp) VALUES (%s, %s, %s)", mod,ver,time.Now().Format(time.RFC3339)))
		return errors.E(op, err, getKind(err))
	}
	return nil
}

//Lines 根据时间和limit获取索引列表
func (i *indexer) Lines(ctx context.Context, since time.Time, limit int) ([]*index.Line, error) {
	const op errors.Op = "mysql.Lines"
	if since.IsZero() {
		since = time.Unix(0, 0)
	}
	sinceStr := since.Format("2006-01-02 15:04:05.000")
	rows, err := i.db.QueryContext(ctx, `SELECT path, version, timestamp FROM indexes WHERE timestamp >= ? LIMIT ?`, sinceStr, limit)
	if err != nil {
		return nil, errors.E(op, err)
	}
	defer func() { _ = rows.Close() }()
	var lines []*index.Line
	for rows.Next() {
		var line index.Line
		err = rows.Scan(&line.Path, &line.Version, &line.Timestamp)
		if err != nil {
			return nil, errors.E(op, err)
		}
		lines = append(lines, &line)
	}
	return lines, nil
}

func getMySQLSource(cfg *config.MySQL) string {
	c := mysql.NewConfig()
	c.Net = cfg.Protocol
	c.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	c.User = cfg.User
	c.Passwd = cfg.Password
	c.DBName = cfg.Database
	c.Params = cfg.Params
	return c.FormatDSN()
}

func getKind(err error) int {
	mysqlErr := &mysql.MySQLError{}
	if !errors.AsErr(err, &mysqlErr) {
		return errors.KindUnexpected
	}
	switch mysqlErr.Number {
	case 1062:
		return errors.KindAlreadyExists
	default:
		return errors.KindUnexpected
	}
}
