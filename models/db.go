package models

import (
	"errors"
	"fmt"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	// postgresql driver
	// _ "github.com/lib/pq"
)

// 	экземпляр БД
var storage *DB

var dbSchema string

//Query type custom Query
type Query string

//Queryx type custom Queryx
type Queryx struct {
	Query  Query
	Params []interface{}
}

//DB type custom DB
type DB struct {
	*sqlx.DB
	schema string
}

//Tx type custom Transacton struct
type Tx struct {
	*sqlx.Tx
}

var (
	// ErrNoGetterFound - у структуры не найден интерфейс Getter
	ErrNoGetterFound = errors.New("No getter found")
	// ErrNoDeleterFound - у структуры не найден интерфейс Deleter
	ErrNoDeleterFound = errors.New("No deleter found")
	// ErrNoSelecterFound - у структуры не найден интерфейс Selecter
	ErrNoSelecterFound = errors.New("No selecter found")
	// ErrNoUpdaterFound - у структуры не найден интерфейс Updater
	ErrNoUpdaterFound = errors.New("No updater found")
	// ErrNoInserterFound - у структуры не найден интерфейс Inserter
	ErrNoInserterFound = errors.New("No inserter found")
)

// Limit функция возвращает тип limitOption, для запроса count строк, начинающихся с offset записи
func Limit(offset, count int) SelectOption {
	return &limitOption{offset, count}
}

type limitOption struct {
	offset int
	count  int
}

func (o *limitOption) Wrap(query string, params []interface{}) (string, []interface{}) {
	query = fmt.Sprintf("SELECT a.* FROM (%s) a LIMIT %d OFFSET %d", query, o.count, o.offset)
	return query, params
}

// SelectOption - интерфейс для обертки запросов (физический смысл не помню)
type SelectOption interface {
	Wrap(string, []interface{}) (string, []interface{})
}

// Selectx - запрос возвращает массив результатов
func (tx *Tx) Selectx(o interface{}, qx Queryx, options ...SelectOption) error {
	q := string(qx.Query)
	params := qx.Params
	log.Warn(q)
	log.Warn(params)
	for _, option := range options {
		q, params = option.Wrap(q, params)
	}

	if u, ok := o.(Selecter); ok {
		return u.Select(tx.Tx, Query(q), params...)
	}
	stmt, err := tx.Preparex(q)
	if err != nil {
		return err
	}
	return stmt.Select(o, params...)
}

// Countx - функция для подсчета количества записей в запросе
func (tx *Tx) Countx(qx Queryx) (int, error) {
	stmt, err := tx.Preparex(fmt.Sprintf("SELECT COUNT(*) FROM (%s) q", string(qx.Query)))
	if err != nil {
		return 0, err
	}
	count := 0
	err = stmt.Get(&count, qx.Params...)
	return count, err
}

// Getx фунция возвращает результат, состоящий из одной записи
func (tx *Tx) Getx(o interface{}, qx Queryx) error {
	log.Infof("Into getx: %+v", qx)
	if u, ok := o.(Getter); ok {
		return u.Get(tx.Tx, qx.Query, qx.Params...)
	}
	stmt, err := tx.Preparex(string(qx.Query))
	if err != nil {
		return err
	}
	return stmt.Get(o, qx.Params...)
}

// Get фунция возвращает результат, состоящий из одной записи
func (tx *Tx) Get(o interface{}, query Query, params ...interface{}) error {
	if u, ok := o.(Getter); ok {
		return u.Get(tx.Tx, query, params...)
	}
	stmt, err := tx.Preparex(string(query))
	if err != nil {
		return err
	}
	return stmt.Get(o, params...)
}

// Update - обобщенная функция обновления данных
func (tx *Tx) Update(o interface{}) error {
	if u, ok := o.(Updater); ok {
		return u.Update(tx.Tx)
	}
	log.Debugf("No updater found for object: %s", reflect.TypeOf(o))
	return ErrNoUpdaterFound
}

// Delete - обобщенная функция удаления данных
func (tx *Tx) Delete(o interface{}) error {
	if u, ok := o.(Deleter); ok {
		return u.Delete(tx.Tx)
	}
	log.Debugf("No deleter found for object: %s", reflect.TypeOf(o))
	return ErrNoDeleterFound
}

// Insert - обобщенная функция вставки данных
func (tx *Tx) Insert(o interface{}) (int64, error) {
	var id int64
	var err error
	if u, ok := o.(Inserter); ok {
		id, err = u.Insert(tx.Tx)
		if err != nil {
			log.Error(err.Error())
		}
		return id, err
	}
	log.Debugf("No inserter found for object: %s", reflect.TypeOf(o))
	return id, ErrNoInserterFound
}

// Begin - функция возвращает транзакцию
func (d *DB) Begin() *Tx {
	tx := d.MustBegin()
	return &Tx{tx}
}

// Updater interface
type Updater interface {
	Update(*sqlx.Tx) error
}

// Inserter interface
type Inserter interface {
	Insert(*sqlx.Tx) (int64, error)
}

// Selecter interface
type Selecter interface {
	Select(*sqlx.Tx, Query, ...interface{}) error
}

//Getter Interface
type Getter interface {
	Get(*sqlx.Tx, Query, ...interface{}) error
}

//Deleter interface
type Deleter interface {
	Delete(*sqlx.Tx) error
}

// InitDB Returns new db object
func InitDB(host string, port int, username string, password string, dbName string, sslMode string, schema string) (*DB, error) {

	d, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		username,
		password,
		dbName,
		sslMode))

	if err != nil {
		log.Fatal("Ошибка в настройках БД: ", err)
		return nil, err
	} else {
		log.Info("Авторизация в базе прошла успешно")
	}
	err = d.Ping()
	if err != nil {
		return nil, err
	}

	_, err = d.Exec(fmt.Sprintf("SET search_path TO %s", schema))
	if err != nil {
		return &DB{}, err
	} else {
		log.Info("Выбрана схема ", schema)
	}
	dbSchema = schema
	database := DB{d, schema}
	return &database, nil
}
