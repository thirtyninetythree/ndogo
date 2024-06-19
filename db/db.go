package db

import "sort"

type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}

var (
	ErrUniqueViolation   = &Error{"Collection already exists"}
	ErrNotFound          = &Error{"Collection doesn't exist"}
	ErrDimensionMismatch = &Error{"The dimension of the vector doesn't match the dimension of the collection"}
)

type Database struct {
	Collections map[string]*Collection
}

type Collection struct {
	Dimension  int
	Distance   DistanceMetric
	Embeddings []*Embedding
}

type Embedding struct {
	ID       string
	Vector   []float32
	Metadata map[string]string
}

func NewDatabase() *Database {
	return &Database{
		Collections: make(map[string]*Collection),
	}
}

func (db *Database) CreateCollection(name string, dimension int, distance DistanceMetric) (*Collection, error) {
	if _, exists := db.Collections[name]; exists {
		return nil, ErrUniqueViolation
	}

	collection := &Collection{
		Dimension:  dimension,
		Distance:   distance,
		Embeddings: make([]*Embedding, 0),
	}

	db.Collections[name] = collection

	return collection, nil
}

func (db *Database) DeleteCollection(name string) error {
	if _, exists := db.Collections[name]; !exists {
		return ErrNotFound
	}

	delete(db.Collections, name)
	return nil
}

func (db *Database) InsertIntoCollection(name string, embedding *Embedding) error {
	collection, exists := db.Collections[name]
	if !exists {
		return ErrNotFound
	}

	if len(embedding.Vector) != collection.Dimension {
		return ErrDimensionMismatch
	}

	if collection.Distance == Cosine {
		embedding.Vector = normalize(embedding.Vector)
	}

	collection.Embeddings = append(collection.Embeddings, embedding)
	return nil
}

func (db *Database) FindSimilarVectors(collectionName string, queryVector []float32, k int) ([]ScoreIndex, error) {
	collection, exists := db.Collections[collectionName]
	if !exists {
		return nil, ErrNotFound
	}

	if len(queryVector) != collection.Dimension {
		return nil, ErrDimensionMismatch
	}

	normalizedQueryVector := normalize(queryVector)

	distances := make([]ScoreIndex, len(collection.Embeddings))
	distanceFunc := getDistanceFunc(collection.Distance)
	for i, embedding := range collection.Embeddings {
		normalizedDbVector := embedding.Vector
		if collection.Distance == Cosine {
			normalizedDbVector = normalize(embedding.Vector)
		}

		score := distanceFunc(normalizedQueryVector, normalizedDbVector, 0.0)

		distances[i] = ScoreIndex{Score: score, Index: i}
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].Score < distances[j].Score
	})

	return distances[:k], nil
}
