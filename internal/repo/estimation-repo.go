package repo

import (
	"context"
	"es/internal/model"
	"fmt"
	"github.com/gammazero/deque"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

// ExecBatchItem is used for making the query which is going to add to database
type ExecBatchItem struct {
	Query     string
	Arguments []interface{}
}

type EstimationRepo interface {
	SaveSegmentTagForUser(model.Estimation)
	GetSegmentTagFor14dLastDays(segment string) (uint32, error)
	DbEexecCronTask() error
}

type EstimationRepoImpl struct {
	pool        *pgxpool.Pool
	DbExexQueue chan ExecBatchItem
}

func (e *EstimationRepoImpl) GetSegmentTagFor14dLastDays(segment string) (uint32, error) {
	var count uint32
	err := e.pool.QueryRow(context.Background(), "select count(*) from live_users where created_at<$1 and segment=$2", time.Now().AddDate(0, 0, -14), segment).Scan(&count)
	return count, err
}

// SaveSegmentTagForUser save data if it is not in live_users data otherwise it updates the date of row
func (e *EstimationRepoImpl) SaveSegmentTagForUser(estimation model.Estimation) {

	e.DbExexQueue <- ExecBatchItem{
		Query:     "insert into live_users(user_id,segment,created_at) VALUES($1,$2,now()) ON CONFLICT (user_id,segment) DO UPDATE SET created_at = now()",
		Arguments: []interface{}{estimation.UserId, estimation.Segment},
	}
}

// DbEexecQueueTask its run on the another go routine every insert gets queued  and  on every 50 Millisecond send all the data in queue
//todo move logic to the service liar later
func (e *EstimationRepoImpl) DbEexecQueueTask() {

	TimeTicker := time.NewTicker(50 * time.Millisecond)

	var PgxBatch pgx.Batch
	var queryQ deque.Deque

	for {
		select {

		case ExexQueueItem := <-e.DbExexQueue:

			PgxBatch.Queue(ExexQueueItem.Query, ExexQueueItem.Arguments...)
			queryQ.PushFront(ExexQueueItem)
			if PgxBatch.Len() > 5000 {
				if err := e.pool.SendBatch(context.Background(), &PgxBatch).Close(); err != nil {
					query_q_len := queryQ.Len()
					for i := 0; i < query_q_len; i++ {
						ExecInfo := queryQ.PopBack().(ExecBatchItem)
						if _, err := e.pool.Exec(context.Background(), ExecInfo.Query, ExecInfo.Arguments...); err != nil {
							fmt.Println(err)
						}
					}
				}
				queryQ.Clear()
				PgxBatch = pgx.Batch{}
			}

		case <-TimeTicker.C:
			if PgxBatch.Len() > 0 {
				if err := e.pool.SendBatch(context.Background(), &PgxBatch).Close(); err != nil {
					query_q_len := queryQ.Len()
					for i := 0; i < query_q_len; i++ {
						ExecInfo := queryQ.PopBack().(ExecBatchItem)
						if _, err := e.pool.Exec(context.Background(), ExecInfo.Query, ExecInfo.Arguments...); err != nil {
							fmt.Println(err)
						}
					}
				}
				queryQ.Clear()
				PgxBatch = pgx.Batch{}
			}

		}

	}
}

// DbEexecCronTask transmute data from live table to archive table
func (e *EstimationRepoImpl) DbEexecCronTask() error {

	two_week_ago := time.Now().AddDate(0, 0, -14)
	_, err := e.pool.Exec(context.Background(), "insert into archive_users select * from live_users where created_at < $1", two_week_ago)
	if err != nil {
		return err
	}
	_, err = e.pool.Exec(context.Background(), "delete from live_users where created_at < $1", two_week_ago)
	if err != nil {
		return err
	}

	return nil
}

func NewEstimationRepo(pool *pgxpool.Pool) *EstimationRepoImpl {
	e := &EstimationRepoImpl{pool: pool, DbExexQueue: make(chan ExecBatchItem, 1000000)}
	go e.DbEexecQueueTask()
	return e
}
