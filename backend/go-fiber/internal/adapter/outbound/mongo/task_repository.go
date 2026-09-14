// Package mongo implements task.Repository against mongo-driver/v2. The
// client this type takes as a constructor argument comes from
// pkg/mongo.
package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

// taskDoc is this adapter's own wire shape — never store core.Task
// directly, the bson tags belong here, not in the domain.
type taskDoc struct {
	ID        string    `bson:"_id"`
	Title     string    `bson:"title"`
	Done      bool      `bson:"done"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func toDoc(t *task.Task) taskDoc {
	return taskDoc{ID: t.ID, Title: t.Title, Done: t.Done, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
}

func (d taskDoc) toTask() *task.Task {
	return &task.Task{ID: d.ID, Title: d.Title, Done: d.Done, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}
}

type TaskRepository struct {
	coll *mongo.Collection
}

var _ task.Repository = (*TaskRepository)(nil)

func NewTaskRepository(client *mongo.Client, database string) *TaskRepository {
	return &TaskRepository{coll: client.Database(database).Collection("tasks")}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	_, err := r.coll.InsertOne(ctx, toDoc(t))
	return err
}

func (r *TaskRepository) Get(ctx context.Context, id string) (*task.Task, error) {
	var d taskDoc
	err := r.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, task.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return d.toTask(), nil
}

func (r *TaskRepository) List(ctx context.Context) ([]*task.Task, error) {
	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var docs []taskDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	out := make([]*task.Task, 0, len(docs))
	for _, d := range docs {
		out = append(out, d.toTask())
	}
	return out, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	res, err := r.coll.UpdateOne(ctx,
		bson.D{{Key: "_id", Value: t.ID}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "title", Value: t.Title},
			{Key: "done", Value: t.Done},
			{Key: "updated_at", Value: t.UpdatedAt},
		}}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return task.ErrNotFound
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	res, err := r.coll.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return task.ErrNotFound
	}
	return nil
}
