package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES ( :client, :status, :address, :createdAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt))
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("scan id failed: %w", err)
	}

	return int(id), nil

}

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, fmt.Errorf("select error: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(`SELECT number, client, status, address, created_at FROM parcel WHERE client = :client`, sql.Named("client", client))
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {

			return res, fmt.Errorf("scan failed: %w", err)
		}

		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows unpacking failed: %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (s ParcelStore) GetStatus(number int) (string, error) {
	// не хочу повторять код в двух других функциях, поэтому отдельно реализовала получение статуса
	p := Parcel{}
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Status)
	if err != nil {
		return "", fmt.Errorf("select failed: %w", err)
	}

	return p.Status, nil

}
func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	status, err := s.GetStatus(number)
	if err != nil {
		return fmt.Errorf("getStatus failed: %w", err)
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("address shouldn't be updated") // не знаю что возвратить
	}
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil

}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	status, err := s.GetStatus(number)
	if err != nil {
		return err
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("parcel shouldn't be deleted")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
	return fmt.Errorf("delete failed: %w", err)

}
