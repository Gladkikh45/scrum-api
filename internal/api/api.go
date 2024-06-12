package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// "os/user"

	"github.com/jackc/pgx/v5"
	// _ "github.com/jackc/pgx/v5"
)

type GetUsersResponse struct {
	Count int64    `json:"count"` // Зачем нужны эти оранжевые слова?
	Users []DbUser `json:"users"`
}

type GetBoardsResponse struct {
	Count  int64     `json:"count"`
	Boards []DbBoard `json:"boards"`
}

type DbBoard struct {
	Title   string   `json:"title,omitempty"`
	Columns []string `json:"columns,omitempty"`
}

type DbUser struct {
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Login       string `json:"login,omitempty"`
	Password    string `json:"password,omitempty"`
}

func Connection() (*pgx.Conn, error) {
	//urlExample := "postgres://postgres:45863@localhost:5432/test_db"
	urlExample := "postgres://postgres:45863@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn, err
}

const sqlCreateReport = `
    INSERT INTO main.boards
        (title, columns)
    VALUES
    ($1, $2)
    RETURNING id
`

type CreateBoardRequest struct {
	Title   string   `json:"title,omitempty"`
	Columns []string `json:"columns,omitempty"`
}

func СreateBoard(w http.ResponseWriter, r *http.Request) {
	Connection()
	conn, err := Connection()
	defer conn.Close(context.Background())

	decoder := json.NewDecoder(r.Body)
	var message CreateBoardRequest

	err = decoder.Decode(&message)
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	//conn.

	// fmt.Scan(title)
	rows, err := conn.Query(context.Background(), sqlCreateReport, message.Title, message.Columns)
	if err != nil {
		JsonResponse(w, 500, err.Error())
		return
	}

	board := BoardCreated{}

	for rows.Next() {
		err = rows.Scan(
			&board.Id,
		)

		if err != nil {
			JsonResponse(w, 500, err.Error())
		}
	}

	JsonResponse(w, 200, board)
}

type BoardCreated struct {
	Id *string `json:"id,omitempty"`
}

func BoardList(w http.ResponseWriter, r *http.Request) {
	Connection()
	conn, err := Connection()

	rows, err := conn.Query(context.Background(), "SELECT title, columns FROM main.boards")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	boards := GetBoardsResponse{Boards: make([]DbBoard, 0)}
	for rows.Next() {
		board := DbBoard{}
		err := rows.Scan(
			&board.Title,
			&board.Columns,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow asdasd: %v\n", err)
			continue
		}
		boards.Boards = append(boards.Boards, board)
		fmt.Printf("%+v", boards)
	}
	boards.Count = int64(len(boards.Boards))
	defer conn.Close(context.Background())
	JsonResponse(w, 200, boards)
}

const sqlDeleteReport = `
    delete from main.boards
    where title = $1
`

func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	Connection()
	conn, err := Connection()

	defer conn.Close(context.Background())

	title := r.URL.Query().Get("title")

	rows, err := conn.Query(context.Background(), sqlDeleteReport, title) // Как отличать когда надо conn.Query, а когда conn.QueryRow?
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		fmt.Println(rows)
		os.Exit(1)
	}

	JsonResponse(w, 200, rows)
}

func Authentication(login, secret string) {
	Connection()
	conn, err := Connection()

	var name string
	data := conn.QueryRow(context.Background(), "select user_name from users_table where user_name=$1", login).Scan(&name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(data)
	defer conn.Close(context.Background())
}

func UsersList(w http.ResponseWriter, r *http.Request) {
	Connection()
	conn, _ := Connection()

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT id, display_name, login, password FROM main.users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	users := GetUsersResponse{Users: make([]DbUser, 0)}
	for rows.Next() {
		user := DbUser{}
		err := rows.Scan(
			&user.Id,
			&user.DisplayName,
			&user.Login,
			&user.Password,
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow asdasd: %v\n", err)
			continue
		}
		users.Users = append(users.Users, user)
		fmt.Printf("%+v", users)
	}

	users.Count = int64(len(users.Users))

	JsonResponse(w, 200, users)
}

func JsonResponse(w http.ResponseWriter, code int, resp any) {
	if resp == nil {
		w.WriteHeader(code)
		return
	}

	response, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error", err.Error())
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// SELECT * FROM public.users_table

// 	_, err = stmt.Exec()
// INSERT INTO users (age, email, first_name, last_name)
// VALUES (30, 'jon@calhoun.io', 'Jonathan', 'Calhoun');
// insert into users_table (user_name, password) values ('moon', 'spell');

// SELECT * FROM public.users_table
// ORDER BY user_id ASC

// func (c *Client) GetProducts(ctx context.Context, msg dbmsg.GetProducts) ([]dbmsg.Product, error) {
// 	if err := c.conn(); err != nil {
// 		return nil, err
// 	}

// 	rows, err := c.driver.QueryContext(ctx, sqlGetProducts)
// 	if err != nil {
// 		//var dbError sqlite3.Error
// 		//if errors.As(err, &dbError) {
// 		// fmt.Println(dbError)
// 		//}

// 		return nil, err
// 	}

// 	products := make([]dbmsg.Product, 0)

// 	for rows.Next() {
// 		p := dbmsg.Product{}
// 		err := rows.Scan(&p.Id, &p.DisplayName, &p.Description)
// 		if err != nil {
// 			// TODO log
// 			continue
// 		}
// 		products = append(products, p)
// 	}

// 	return products, nil
// }
