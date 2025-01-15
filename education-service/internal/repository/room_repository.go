package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"fmt"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) CreateRoom(companyId string, title string, capacity int32) error {
	query := "INSERT INTO rooms (title, capacity , company_id) VALUES ($1 , $2 ,$3)"
	_, err := r.db.Exec(query, title, capacity, companyId)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

func (r *RoomRepository) UpdateRoom(companyId string, id, title *string, capacity *int32) error {
	query := "UPDATE rooms SET title = $1 , capacity=$2 WHERE id = $3 and company_id=$4"
	_, err := r.db.Exec(query, title, capacity, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}
	return nil
}

func (r *RoomRepository) DeleteRoom(companyId string, id string) error {
	query := "DELETE FROM rooms WHERE id = $1 and company_id=$3"
	_, err := r.db.Exec(query, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

func (r *RoomRepository) GetRoom(companyId string) (*pb.GetUpdateRoomAbs, error) {
	query := "SELECT id , title, capacity FROM rooms where company_id=$1"
	rows, err := r.db.Query(query, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result pb.GetUpdateRoomAbs
	var rooms []*pb.AbsRoom
	for rows.Next() {
		var res pb.AbsRoom
		err := rows.Scan(&res.Id, &res.Name, &res.Capacity)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &res)
	}
	result.Rooms = rooms
	return &result, nil
}
