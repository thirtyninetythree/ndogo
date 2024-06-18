package db

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

func (db *Database) CreateCollection(name string, dimension int, distance DistanceMetric) error {
	if _, exists := db.Collections[name]; exists {
		return ErrUniqueViolation
	}

	db.Collections[name] = &Collection{
		Dimension:  dimension,
		Distance:   distance,
		Embeddings: make([]*Embedding, 0),
	}

	return nil
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

	// Normalize the vector if the distance metric is cosine
	if collection.Distance == CosineDistance {
		embedding.Vector = normalize(embedding.Vector)
	}

	collection.Embeddings = append(collection.Embeddings, embedding)
	return nil
}
