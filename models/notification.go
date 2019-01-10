package models

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

// 	запросы
var (
	notificationTableName         = "notifications"
	queryAllNotifications         = "SELECT id, date_create, date_read, opened, message, read, subject, user_from, user_to FROM %s"
	queryCountNotificationsByName = "SELECT count(*) as count FROM %s WHERE user_to=$1 AND read=false"
	queryNotificationsByUserName  = "SELECT id, date_create, date_read, opened, message, read, subject, user_from, user_to FROM %s WHERE user_to=$1 AND read=$2"
	queryNotificationsByID        = "SELECT id, date_create, date_read, opened, message, read, subject, user_from, user_to FROM %s WHERE id=$1"
	queryInsertNotification       = "INSERT INTO %s (id, message, subject, user_from, user_to)  VALUES (:id, :message, :subject, :user_from, :user_to) RETURNING id"
	queryUpdateNotification       = "UPDATE %s SET date_create=:date_create , date_read=:date_read, opened=:opened, message=:message, read=:read, subject=:subject, user_from=:user_from, user_to=:user_to WHERE id=:id"
	querySwitchOpenedNotification = "UPDATE %s SET opened=not (SELECT opened FROM notification_service.notifications where id=:id) WHERE id=:id"
	querySetReadNotification      = "UPDATE %s SET date_read=current_timestamp, read=true WHERE id=:id"
	queryDeleteNotification       = "DELETE FROM %s WHERE id=:id"
)

// 	инициируем
func init() {
	// 	регистрируем таблицу в схеме
	RegisterTable(GetNotificationsSchema())
}

// 	инициируем
func InitNotificationsSchema() {
	// 	регистрируем таблицу в схеме
	RegisterTable(GetNotificationsSchema())
}

// queryAllNotifications - возврашает список всех сообщений
func QueryAllNotifications() Queryx {
	log.Info("QueryAllNotifications")
	return Queryx{
		Query:  getSchemedQuery(queryAllNotifications, notificationTableName),
		Params: []interface{}{},
	}
}

// queryNotificationsByUserName - возврашает список сообщений конкретного пользователя
func QueryCountNotificationsByName(user string) Queryx {
	log.Info("QueryNotificationsByUserName")
	return Queryx{
		Query: getSchemedQuery(queryCountNotificationsByName, notificationTableName),
		Params: []interface{}{
			user,
		},
	}
}

// queryNotificationsByUserName - возврашает список сообщений конкретного пользователя
func QueryNotificationsByUserName(user_to string) Queryx {
	log.Info("QueryNotificationsByUserName")
	return Queryx{
		Query: getSchemedQuery(queryNotificationsByUserName, notificationTableName),
		Params: []interface{}{
			user_to,
		},
	}
}

// queryNotificationsByUserName - возврашает список сообщений конкретного пользователя
func QueryNotificationsByID(id string) Queryx {
	log.Info("QueryNotificationsByUserName")
	return Queryx{
		Query: getSchemedQuery(queryNotificationsByID, notificationTableName),
		Params: []interface{}{
			id,
		},
	}
}

// QuerySetReadNotification - возврашает список сообщений конкретного пользователя
func (noticeRequest *NoticeRequest) UpdateSetReadNotification(tx *sqlx.Tx) error {
	_, err := tx.NamedExec(string(getSchemedQuery(querySetReadNotification, notificationTableName)), noticeRequest)
	return err
}

// QuerySwitchOpenedNotification - возврашает список сообщений конкретного пользователя
func (noticeRequest *NoticeRequest) UpdateSwitchOpenedNotification(tx *sqlx.Tx) error {
	log.Warnf("%+v", noticeRequest)
	_, err := tx.NamedExec(string(getSchemedQuery(querySwitchOpenedNotification, notificationTableName)), noticeRequest)
	return err
}

//	добавление нового сообщения
func (noticeRequest *NoticeRequest) Insert(tx *sqlx.Tx) (int64, error) {
	// var lastID int64
	r, err := tx.NamedQuery(string(getSchemedQuery(queryInsertNotification, notificationTableName)), noticeRequest)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	return 0, nil
}

//	изменение в таблице сообщений
func (noticeRequest *NoticeRequest) Update(tx *sqlx.Tx) error {
	_, err := tx.NamedExec(string(getSchemedQuery(queryUpdateNotification, notificationTableName)), noticeRequest)
	return err
}

// 	удаление сообщения из таблицы сообщений
func (notification *Notification) Delete(tx *sqlx.Tx) error {
	_, err := tx.NamedExec(string(getSchemedQuery(queryDeleteNotification, notificationTableName)), notification)
	return err
}

// GetNotificationsSchema - возвращает структуру для создания таблицы с сообщениями
func GetNotificationsSchema() Table {
	return Table{
		Tablename: string(getSchemedQuery("%s", notificationTableName)),
		Schema: []Column{
			{
				Colname:    "id",
				Coltype:    "varchar",
				PrimaryKey: true,
			},
			{
				Colname: "date_create",
				Coltype: "timestamp",
				Default: "current_timestamp",
			},
			{
				Colname:    "date_read",
				Coltype:    "timestamp",
				IsNullable: true,
			},
			{
				Colname: "opened",
				Coltype: "boolean",
				Default: "false",
			},
			{
				Colname: "message",
				Coltype: "varchar",
			},
			{
				Colname: "read",
				Coltype: "boolean",
				Default: "false",
			},
			{
				Colname: "subject",
				Coltype: "varchar",
			},
			{
				Colname: "user_from",
				Coltype: "varchar",
			},
			{
				Colname: "user_to",
				Coltype: "varchar",
			},
		},
	}
}
