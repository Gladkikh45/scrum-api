package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	//"os/user"
	//"github.com/jackc/pgx/v5"
	//_ "github.com/jackc/pgx/v5"
)

type GetUsersResponse struct {
	Count int64    `json:"count"` // Зачем нужны эти оранжевые(зелёные) слова?
	Users []DbUser `json:"users"`
}

type GetBoardsResponse struct {
	Count  int64     `json:"count"`
	Boards []DbBoard `json:"boards"`
}

type GetCardsResponse struct {
	Count int64    `json:"count"`
	Cards []DbCard `json:"cards"`
}

type DbBoard struct {
	Id        *string    `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Title     *string    `json:"title,omitempty"`
	Columns   []string   `json:"columns,omitempty"`
}

type DbUser struct {
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Login       string `json:"login,omitempty"`
	Password    string `json:"password,omitempty"`
}

type DbCard struct {
	Id          *string    `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Board       *string    `json:"board,omitempty"`
	BoardId     *string    `json:"BoardId,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Description *string    `json:"description,omitempty"`
	Assignee    *string    `json:"assignee,omitempty"`
	Estimation  *string    `json:"estimation,omitempty"`
}

func Connection() (*pgx.Conn, error) {
	//urlExample := "postgres://postgres:45863@localhost:5432/test_db"
	urlExample := "postgres://postgres:45863@localhost:5432/postgres"
	//urlExample := "postgres://ngfw@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn, err
}

func CreqteConn() (*pgxpool.Pool, error) {
	connStr := "postgres://postgres@localhost:5432/ngfw"
	pollConf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pollConf.LazyConnect = true

	pool, err := pgxpool.ConnectConfig(context.Background(), pollConf)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

//func New(config *configuration.Db) (*Client, error) {
//	pollConf, err := pgxpool.ParseConfig(config.ConnStr)
//	if err != nil {
//		return nil, err
//	}
//
//	pollConf.LazyConnect = true
//
//	pool, err := pgxpool.ConnectConfig(context.Background(), pollConf)
//	if err != nil {
//		return nil, err
//	}
//
//	c := &Client{
//		driver: pool,
//	}
//
//	return c, nil
//}

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
	conn, err := Connection()
	//defer conn.Close(context.Background())

	decoder := json.NewDecoder(r.Body)
	var message CreateBoardRequest

	err = decoder.Decode(&message)
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	//conn.
	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}
	defer tx.Rollback(context.TODO())

	rows, err := tx.Query(context.Background(), sqlCreateReport, message.Title, message.Columns)
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
	_ = tx.Commit(context.Background())
	JsonResponse(w, 200, board)
}

type BoardCreated struct {
	Id *string `json:"id,omitempty"`
}

func BoardList(w http.ResponseWriter, r *http.Request) {

	conn, err := Connection()

	rows, err := conn.Query(context.Background(), "SELECT id, title, columns FROM main.boards")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	boards := GetBoardsResponse{Boards: make([]DbBoard, 0)}
	for rows.Next() {
		board := DbBoard{}
		err := rows.Scan(
			&board.Id,
			//&board.CreatedAt,
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
    where id = $1
`

type DeleteBoardRequest struct {
	Id *string `json:"id,omitempty"`
}

func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := Connection()

	defer conn.Close(ctx)

	decoder := json.NewDecoder(r.Body)
	var message DeleteBoardRequest

	err = decoder.Decode(&message)
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, sqlDeleteReport, message.Id)
	if err != nil {
		JsonResponse(w, 500, err.Error())
		return
	}

	rows.Next()

	if !rows.CommandTag().Delete() {
		JsonResponse(w, 400, "object does not have delete statement")
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	JsonResponse(w, 200, rows)
}

//func Authentication(login, secret string) {
//	Connection()2
//	conn, err := Connection()
//
//	var name string
//	data := conn.QueryRow(context.Background(), "select user_name from  where user_name=$1", login).Scan(&name)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
//		os.Exit(1)
//	}
//	fmt.Println(data)
//	defer conn.Close(context.Background())
//}

func UsersList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, _ := Connection()

	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT id, display_name, login, password FROM main.users")
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

const sqlCreateCard = `
	INSERT INTO main.cards
	(title, board, board_id, status, description, assignee, estimation)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
`

type CreateCardRequest struct {
	Title       string `json:"title,omitempty"`
	Board       string `json:"board,omitempty"`
	BoardID     string `json:"boardID,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
	Assignee    string `json:"assignee,omitempty"`
	Estimation  string `json:"estimation,omitempty"`
}

type CardCreated struct {
	Id *string `json:"id,omitempty"`
}

func СreateCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := Connection()
	defer conn.Close(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	decoder := json.NewDecoder(r.Body)
	var message CreateCardRequest

	err = decoder.Decode(&message)
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}

	rows, err := tx.Query(ctx, sqlCreateCard, message.Title, message.Board, message.BoardID, message.Status, message.Description, message.Assignee, message.Estimation)
	if err != nil {
		JsonResponse(w, 500, err.Error())
		return
	}

	type CardCreated struct {
		Id *string `json:"id,omitempty"`
	}

	card := CardCreated{}

	for rows.Next() {
		err = rows.Scan(
			&card.Id,
		)

		if err != nil {
			JsonResponse(w, 500, err.Error())
		}
	}
	_ = tx.Commit(context.Background())
	JsonResponse(w, 200, card)
}

const sqlUpdateCard = `
	update main.cards
   set title = coalesce($2, title), 
       status = coalesce($3, status), 
       description = coalesce($4, description), 
       assignee = coalesce($5, assignee), 
       estimation = coalesce($6, estimation), 
       updated_at = now()
   where id = $1

	RETURNING id
`

// TODO: sql запрос

type UpdateCardRequest struct {
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
	Assignee    string `json:"assignee,omitempty"`
	Estimation  string `json:"estimation,omitempty"`
}

// TODO: здесь запрашивается то, что будет нужно для обновления
//type CardUpdate struct {
//	Id *string `json:"id,omitempty"`
//}

func UpdateCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := Connection()
	defer conn.Close(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	decoder := json.NewDecoder(r.Body)
	var message UpdateCardRequest

	err = decoder.Decode(&message)
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}

	rows, err := tx.Query(ctx, sqlUpdateCard, message.Id, message.Title, message.Status, message.Description, message.Assignee, message.Estimation)
	if err != nil {
		JsonResponse(w, 500, err.Error())
		return
	}

	type CardUpdated struct {
		Id *string `json:"id,omitempty"`
	}

	card := CardUpdated{}

	for rows.Next() {
		err = rows.Scan(
			&card.Id,
		)

		if err != nil {
			JsonResponse(w, 500, err.Error())
		}
	}
	_ = tx.Commit(context.Background())
	JsonResponse(w, 200, card)
}

const sqlDeleteCard = `
	delete from main.cards
	where id = $1
	
	returning id
`

type DeleteCardRequest struct {
	Id string `json:"id,omitempty"`
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := Connection()
	defer conn.Close(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}
	decoder := json.NewDecoder(r.Body)
	var message DeleteCardRequest

	err = decoder.Decode(&message)
	if err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}

	rows, err := tx.Query(ctx, sqlDeleteCard, message.Id)
	if err != nil {
		JsonResponse(w, 500, err.Error())
		return
	}

	type CardDeleted struct {
		Id *string `json:"id,omitempty"`
	}

	card := CardDeleted{}

	for rows.Next() {
		err = rows.Scan(
			&card.Id,
		)

		if err != nil {
			JsonResponse(w, 500, err.Error())
		}
	}
	_ = tx.Commit(context.Background())
	JsonResponse(w, 200, card)
}

type ReportRequest struct {
	Board    string `json:"board,omitempty"`
	Status   string `json:"status,omitempty"`
	Assignee string `json:"assignee,omitempty"`
}

type ReportCardsCreated struct {
	Title       string     `json:"title,omitempty"`
	Board       string     `json:"board,omitempty"`
	Status      string     `json:"status,omitempty"`
	Assignee    string     `json:"assignee,omitempty"`
	Estimation  string     `json:"estimation,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

const sqlReportCardsRequest = `
	select cards.title,
	       cards.board,
	       cards.status,
		   cards.assignee,
		   cards.estimation,
	       cards.description,
	       cards.created_at
from main.cards
	where cards.board = $1 and cards.status = $2 and cards.assignee = $3
`

func ReportCards(ctx context.Context, board, status, assignee string) ([]ReportCardsCreated, error) {
	conn, err := Connection()
	defer conn.Close(ctx)

	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	rows, err := conn.Query(ctx, sqlReportCardsRequest, board, status, assignee)
	if err != nil {
		return nil, err
	}

	reportcards := make([]ReportCardsCreated, 0)

	for rows.Next() {
		report := ReportCardsCreated{}
		err = rows.Scan(
			&report.Title,
			&report.Board,
			&report.Status,
			&report.Assignee,
			&report.Estimation,
			&report.Description,
			&report.CreatedAt)
		if err != nil {
			return nil, err
		}
		reportcards = append(reportcards, report)
	}

	return reportcards, err
}

const sqlReport = `
	select cards.board,
       cards.status,
       cards.assignee,
       cards.estimation,
       cards.description
from main.boards inner join main.cards
                            on main.boards.id=main.cards.board_id
	where cards.board = $1 and cards.status = $2 and cards.assignee = $3
`

type ReportCreated struct {
	Board       string               `json:"board,omitempty"`
	Status      string               `json:"status,omitempty"`
	Assignee    string               `json:"assignee,omitempty"`
	Estimation  string               `json:"estimation,omitempty"`
	Description string               `json:"description,omitempty"`
	Cards       []ReportCardsCreated `json:"cards,omitempty"`
}

//func CreateCardsReport(ctx context.Context, board, status, assignee string) {
//	conn, err := Connection()
//	defer conn.Close(ctx)
//
//	if err != nil {
//		fmt.Println("Error", err.Error())
//		os.Exit(1)
//	}
//
//rows, err := conn.Query(ctx, sqlReportCardsRequest, board, status, assignee)
//if err != nil {
//	fmt.Println("Error", err.Error())
//}
//}

// TODO: Почему используется именно констнанта?

func Report(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := Connection()
	defer conn.Close(ctx)

	// TODO: Зачем здесь context.Background?

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	decoder := json.NewDecoder(r.Body)
	var message ReportRequest

	if err = decoder.Decode(&message); err != nil {
		JsonResponse(w, 400, err.Error())
		return
	}

	// TODO: Пользоваться конструкцией сверху с if если функция возвращает только ошибку

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}

	row := tx.QueryRow(ctx, sqlReport, message.Board, message.Status, message.Assignee)
	//
	report := ReportCreated{}

	err = row.Scan(
		&report.Board,
		&report.Status,
		&report.Assignee,
		&report.Estimation,
		&report.Description,
	)

	// TODO: У сards статус varchar, а у досок и репортов масив!

	if err != nil {
		JsonResponse(w, 500, err.Error())
	}
	cards, err := ReportCards(ctx, message.Board, message.Status, message.Assignee)
	if err != nil {
		JsonResponse(w, 500, err.Error())
	}
	report.Cards = cards

	//_ = tx.Commit(context.Background())
	//JsonResponse(w, 200, report)

	estimation, err := Estimation(ctx, message.Board, message.Status, message.Assignee)
	if err != nil {
		JsonResponse(w, 500, err.Error())
	}
	fmt.Println(estimation)

	_ = tx.Commit(context.Background())
	JsonResponse(w, 200, report)
}

type EstimationCreated struct {
	Estimation string `json:"estimation,omitempty"`
}

const sqlEstimationRequest = `
	select cards.estimation
from main.cards
	where cards.board = $1 and cards.status = $2 and cards.assignee = $3
`

func Estimation(ctx context.Context, board, status, assignee string) ([]EstimationCreated, error) {
	conn, err := Connection()
	defer conn.Close(ctx)

	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}
	rows, err := conn.Query(ctx, sqlEstimationRequest, board, status, assignee)
	if err != nil {
		return nil, err
	}

	estimation := make([]EstimationCreated, 0)

	for rows.Next() {
		report := EstimationCreated{}
		err = rows.Scan(
			&report.Estimation)
		if err != nil {
			return nil, err
		}
		estimation = append(estimation, report)
	}
	estimationHours(estimation)

	return estimation, err
}

func estimationHours(est []EstimationCreated) (a int) {
	res := est[0]
	fmt.Printf("T%", res)
	return 1
}

//func test
//
//SELECT * FROM public.users_table
//
//	_, err = stmt.Exec()
//INSERT INTO users (age, email, first_name, last_name)
//VALUES (30, 'jon@calhoun.io', 'Jonathan', 'Calhoun');
//insert into users_table (user_name, password) values ('moon', 'spell');
//
//SELECT * FROM public.users_table
//ORDER BY user_id ASC
//
//func (c *Client) GetProducts(ctx context.Context, msg dbmsg.GetProducts) ([]dbmsg.Product, error) {
//	if err := c.conn(); err != nil {
//		return nil, err
//	}
//
//	rows, err := c.driver.QueryContext(ctx, sqlGetProducts)
//	if err != nil {
//		//var dbError sqlite3.Error
//		//if errors.As(err, &dbError) {
//		// fmt.Println(dbError)
//		//}
//
//		return nil, err
//	}
//
//	products := make([]dbmsg.Product, 0)
//
//	for rows.Next() {
//		p := dbmsg.Product{}
//		err := rows.Scan(&p.Id, &p.DisplayName, &p.Description)
//		if err != nil {
//			// TODO log
//			continue
//		}
//		products = append(products, p)
//	}
//
//	return products, nil
//}
