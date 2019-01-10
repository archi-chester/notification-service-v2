package models

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var tableList []Table

// Column - структура с описанием таблицы
type Column struct {
	Colname              string
	Coltype              string
	IsNullable           bool
	PrimaryKey           bool
	ForeignKeyReferences string
	MaximumLength        int
	Unique               bool
	Default              string
}

// Table is for Table Schema structures
type Table struct {
	Tablename string
	Schema    []Column
}

// 	Создаем таблицу
func (d *DB) createTable(table Table) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	var tableColumnsDefList []string
	for _, col := range table.Schema {
		c := fmt.Sprintf("%s %s", col.Colname, col.Coltype)
		if col.MaximumLength != 0 {
			c += fmt.Sprintf("(%d)", col.MaximumLength)
		}
		if !col.IsNullable {
			c += " NOT NULL"
		}
		if col.PrimaryKey {
			c += " PRIMARY KEY"
		}
		if col.Unique {
			c += " UNIQUE"
		}
		if col.ForeignKeyReferences != "" {
			c += fmt.Sprintf(" REFERENCES %s", col.ForeignKeyReferences)
		}
		if col.Default != "" {
			c += fmt.Sprintf(" DEFAULT %s", col.Default)
		}
		tableColumnsDefList = append(tableColumnsDefList, c)
	}

	log.Warn("Пробую создать таблицу: ", strings.Join(tableColumnsDefList, ","))
	stmtStr := fmt.Sprintf("CREATE TABLE %s (%s)", table.Tablename, strings.Join(tableColumnsDefList, ","))
	_, err = d.Exec(stmtStr)
	if err != nil {
		return err
	}
	log.Infof("Таблица %s создана", table.Tablename)

	return nil

}

// CreateSchema - создает схему в БД
func (d *DB) CreateSchema(tables []Table) error {
	err := d.Ping()
	if err != nil {
		return err
	}

	// 	создаем новую схему
	if d.schema != "" {
		log.Infof("Удаляю схему %s", d.schema)
		_, err = d.Exec("DROP SCHEMA IF EXISTS " + d.schema + " CASCADE")
		if err != nil {
			return err
		}
		log.Infof("Создаю схему %s", d.schema)
		_, err = d.Exec("CREATE SCHEMA " + d.schema)
		db.Exec("SET search_path TO testdb1"
		if err != nil {
			return err
		}
		log.Infof("Схема %s создана", d.schema)
	}

	for _, table := range tables {
		err := d.createTable(table)
		if err != nil {
			log.Errorf("Невозможно создать таблицу %s: %s", table.Tablename, err)
		}
	}

	return nil
}

func getSchemedQuery(qString string, tableName ...interface{}) Query {
	if dbSchema != "" {
		for i := range tableName {
			tableName[i] = dbSchema + "." + tableName[i].(string)
		}
	}
	log.Warn("Запрос селекта", fmt.Sprintf(qString, tableName...))
	return Query(fmt.Sprintf(qString, tableName...))
}

// GetAllRegisteredTables - функция возвращает зарегистрированные таблицы
func GetAllRegisteredTables() []Table {
	return tableList
}

// RegisterTable - функция регистрирует таблицу в общем списке таблиц
func RegisterTable(table Table) {
	tableList = append(tableList, table)
}
