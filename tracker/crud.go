package tracker

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	db      = "tracker"
	bookCol = "book"
	noteCol = "note"
)

// Book
func listBook(query map[string]string) (books []Book, err error) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database(db).Collection(bookCol)

	listAll := true
	for _, v := range query {
		if v != "" {
			listAll = false
		}
	}
	if listAll == true {
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println("error 1")
			log.Fatal(err)
			return books, err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var book Book
			if err = cursor.Decode(&book); err != nil {
				log.Fatal(err)
				return books, err
			} else {
				books = append(books, book)
			}
		}
		return books, nil
	} else {
		var filter bson.D
		for k, v := range query {
			if v != "" {
				fmt.Println(v)
				if err != nil {
					log.Fatal(err)
				}
				filter = append(filter, bson.E{k, v})
			}
		}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			fmt.Println("error 1")
			log.Fatal(err)
			return books, err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var book Book
			if err = cursor.Decode(&book); err != nil {
				log.Fatal(err)
				return books, err
			} else {
				books = append(books, book)
			}
		}
		return books, nil
	}
}

func getBook(id string) (book Book, err error) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database(db).Collection(bookCol)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
		return book, err
	}

	res := collection.FindOne(ctx, bson.M{"id": oid})
	res.Decode(&book)

	return book, nil
}

func addBook(book *Book) (primitive.ObjectID, error) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	book.ID = primitive.NewObjectID()

	collection := client.Database(db).Collection(bookCol)

	res, err := collection.InsertOne(ctx, book)
	if err != nil {
		log.Printf("Could not create Book: %v", err)
		return primitive.NilObjectID, err
	}
	oid := res.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func deleteBook(ids []string) (int, error) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database(db).Collection(bookCol)

	var oids bson.D
	for _, id := range ids {
		if id != "" {
			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				log.Println("Invalid id")
				return 0, err
			}
			oids = append(oids, bson.E{"id", oid})
		}
	}

	res, err := collection.DeleteMany(ctx, oids)
	if err != nil {
		log.Fatal(err)
		return int(res.DeletedCount), err
	}
	return int(res.DeletedCount), nil
}

func editBook(fields map[string]interface{}) (int, error) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database(db).Collection(bookCol)

	var updateFields bson.D
	for k, v := range fields {
		if v != "" {
			updateFields = append(updateFields, bson.E{k, v})
		}
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"id": fields["id"]},
		bson.D{
			{"$set", updateFields},
		},
	)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return int(result.ModifiedCount), nil
}

// Note
