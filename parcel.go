package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

// Создаётся экземпляр структуры ParcelStore в которой хранится объект БД'ых
func NewParcelStore(dbase *sql.DB) ParcelStore {
	return ParcelStore{db: dbase}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление информации в таблицу parcel, используйте данные из переменной p, являющейся экземпляром структуры Parсel
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number).Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	p.Number = number

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Parcel
		_ = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		res = append(res, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :stat WHERE number = :num",
		sql.Named("stat", status),
		sql.Named("num", number))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :adrs WHERE number = :num AND status = :stat",
		sql.Named("adrs", address),
		sql.Named("num", number),
		sql.Named("stat", ParcelStatusRegistered))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	if err != nil {
		return err
	}
	return nil
}
