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

func (r *RoomRepository) CreateRoom(title string, capacity int32) error {
	query := "INSERT INTO rooms (title, capacity) VALUES ($1 , $2)"
	_, err := r.db.Exec(query, title, capacity)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

func (r *RoomRepository) UpdateRoom(id, title *string, capacity *int32) error {
	query := "UPDATE rooms SET title = $1 , capacity=$2 WHERE id = $3"
	_, err := r.db.Exec(query, title, capacity, id)
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}
	return nil
}

func (r *RoomRepository) DeleteRoom(id string) error {
	query := "DELETE FROM rooms WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

func (r *RoomRepository) GetRoom() (*pb.GetUpdateRoomAbs, error) {
	query := "SELECT id , title, capacity FROM rooms"
	rows, err := r.db.Query(query)
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
