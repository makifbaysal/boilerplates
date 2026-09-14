// Package mongo implements job.Store against mongo-driver/v2. The
// client this type takes as a constructor argument comes from
// pkg/mongo.
package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
)

// jobDoc is this adapter's own wire shape — never store core.Job
// directly, the bson tags belong here, not in the domain.
type jobDoc struct {
	ID          string    `bson:"_id"`
	TaskID      string    `bson:"task_id"`
	Type        string    `bson:"type"`
	Payload     string    `bson:"payload"`
	ReceivedAt  time.Time `bson:"received_at"`
	ProcessedAt time.Time `bson:"processed_at"`
}

func toDoc(j *job.Job) jobDoc {
	return jobDoc{ID: j.ID, TaskID: j.TaskID, Type: j.Type, Payload: j.Payload, ReceivedAt: j.ReceivedAt, ProcessedAt: j.ProcessedAt}
}

func (d jobDoc) toJob() *job.Job {
	return &job.Job{ID: d.ID, TaskID: d.TaskID, Type: d.Type, Payload: d.Payload, ReceivedAt: d.ReceivedAt, ProcessedAt: d.ProcessedAt}
}

type JobStore struct {
	coll *mongo.Collection
}

var _ job.Store = (*JobStore)(nil)

func NewJobStore(client *mongo.Client, database string) *JobStore {
	return &JobStore{coll: client.Database(database).Collection("jobs")}
}

func (s *JobStore) Record(ctx context.Context, j *job.Job) error {
	_, err := s.coll.InsertOne(ctx, toDoc(j))
	return err
}

func (s *JobStore) Get(ctx context.Context, id string) (*job.Job, error) {
	var d jobDoc
	err := s.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, job.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return d.toJob(), nil
}

func (s *JobStore) List(ctx context.Context) ([]*job.Job, error) {
	cur, err := s.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var docs []jobDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	out := make([]*job.Job, 0, len(docs))
	for _, d := range docs {
		out = append(out, d.toJob())
	}
	return out, nil
}
