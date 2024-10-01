package repository

import (
	"database/sql"
	"log"
	"net/http"
	"user-service/proto/pb"
)

type UserRepository interface {
	CreateUser(name, surname, phoneNumber, role, chatId string) (*pb.AbsResponse, error)
	GetUserByChatIdOrPhone(chatId, phoneNumber, id *string) (*pb.User, error)
	GetAllUsers(page, size int32) (*pb.GetAllUserResponse, error)
	UpdateNameOrSurname(name, surname, userId string) (*pb.AbsResponse, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(name, surname, phoneNumber, role, chatId string) (*pb.AbsResponse, error) {
	_, err := r.db.Exec("INSERT INTO users (name, surname, phone_number, role , chat_id) VALUES ($1, $2, $3, $4 , $5);", name, surname, phoneNumber, role, chatId)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return &pb.AbsResponse{
			Status:  http.StatusConflict,
			Message: err.Error(),
		}, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "user created successfully",
	}, nil
}

func (r *PostgresUserRepository) GetUserByChatIdOrPhone(chatId, phoneNumber, id *string) (*pb.User, error) {
	var user pb.User
	if chatId != nil {
		err := r.db.QueryRow(`SELECT id , name , surname, phone_number , chat_id, role FROM users where chat_id=$1`, chatId).Scan(&user.Id, &user.Name, &user.Surname, &user.PhoneNumber, &user.ChatId, &user.Role)
		if err != nil {
			return nil, err
		}
	} else if phoneNumber != nil {
		err := r.db.QueryRow(`SELECT id , name , surname, phone_number , chat_id, role FROM users where phone_number=$1`, phoneNumber).Scan(&user.Id, &user.Name, &user.Surname, &user.PhoneNumber, &user.ChatId, &user.Role)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.QueryRow(`SELECT id , name , surname , phone_number , chat_id , role FROM users where id=$1`, id).Scan(&user.Id, &user.Name, &user.Surname, &user.PhoneNumber, &user.ChatId, &user.Role)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetAllUsers(page, size int32) (*pb.GetAllUserResponse, error) {
	return nil, nil
}

func (r *PostgresUserRepository) UpdateNameOrSurname(name, surname, userId string) (*pb.AbsResponse, error) {
	_, err := r.db.Exec(`UPDATE users SET name=$1 , surname=$2 where id=$3`, name, surname, userId)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "updated",
	}, nil
}
