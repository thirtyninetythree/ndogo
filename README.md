# ndogo: A Tiny Vector Database

ndogo is a lightweight, in-memory vector database designed for fast and efficient similarity search. A one to one port of Miguel's Tiny, a vector databse in rust. 


## TODO
- Save to disk
- Add an API and package with docker

# Usage
## Create a Database:
```
db := db.NewDatabase()
```

## Create a Collection
```
collection, err := db.CreateCollection("products", 128, db.CosineDistance)
if err != nil {
    fmt.Println("Error creating collection:", err)
    return
}
```

## Insert Embeddings
```
embedding := db.Embedding{
    ID:       "product_1",
    Vector:   []float32{1.0, 2.0, 3.0, /* ... 128 elements */}, 
    Metadata: map[string]string{"category": "electronics"},
}

err = db.InsertIntoCollection("products", &embedding)
if err != nil {
    fmt.Println("Error inserting embedding:", err)
    return
}
```

## Find Similar Vectors:
```
queryVector := []float32{1.1, 2.1, 3.1, /* ... 128 elements */ } 
similarVectors, err := db.FindSimilarVectors("products", queryVector, 10)
if err != nil {
    fmt.Println("Error finding similar vectors:", err)
    return
}

fmt.Println("Similar Vectors:")
for _, si := range similarVectors {
    fmt.Printf("Vector at index %d: Score = %f\n", si.Index, si.Score)
}
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.
## Acknowledgements
This project is inspired by and builds upon Tiny, a vector DB written in Rust.
## Support
If you have any questions or issues, please feel free to open an issue on GitHub.
