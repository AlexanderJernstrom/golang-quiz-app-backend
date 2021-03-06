package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	Quiz "server/models/QuizModel"
	User "server/models/UserModel"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)


type DB struct{
	Client *mongo.Client
	
}

type Claims struct{
	ID *primitive.ObjectID `json:"id"`
	jwt.StandardClaims
}

type QuestionAnswer struct{
	QuestionID *primitive.ObjectID `json:"id"`
	Answer *string `json:"answer"`
}

type SaveResultBody struct{
	RightAmount *int `json:"rightAmount"`
	QuizId *string `json:"_id"` 
}

var mongoString string = os.Getenv("MONGO_URI")
var Key string = os.Getenv("key")

func ConnectToDatabase() *DB {


	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoString))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		Client: client,
		
	}
}



func (db *DB) GetAllQuizzes(c *fiber.Ctx){
	quizCollection := db.Client.Database("test").Collection("quizzes")

	id  := c.Locals("id")
	fmt.Println(id)
	cursor, err := quizCollection.Find(context.Background(), bson.M{})
	defer cursor.Close(context.Background())
	if err != nil{
		log.Fatal(err)
	}
	var quizzes []Quiz.Quiz

	if err = cursor.All(context.Background(), &quizzes); err != nil{
		
		log.Fatal(err) 
	}

	
	c.JSON(quizzes)
}

func (db *DB) GetQuiz(c *fiber.Ctx){
	quizCollection := db.Client.Database("test").Collection("quizzes")

	quizId := c.Params("id")
	formattedId, _ := primitive.ObjectIDFromHex(quizId)

	var quiz Quiz.Quiz

	quizCollection.FindOne(context.Background(), bson.M{"_id": formattedId}).Decode(&quiz)

	c.JSON(quiz)
}

func (db *DB) CreateQuiz(c *fiber.Ctx){
	userId := c.Locals("id")

	quizCollection := db.Client.Database("test").Collection("quizzes")
	userCollection := db.Client.Database("test").Collection("users")
	newQuiz := new(Quiz.Quiz)


	if err := c.BodyParser(newQuiz); err != nil{
		log.Fatal(err)
		return
	}
	
	

	result, err := quizCollection.InsertOne(context.Background(), newQuiz)

	userCollection.UpdateOne( context.Background(),
		bson.M{"_id": userId},
		bson.D{
			{"$push", bson.D{{"createdQuizzes", result.InsertedID}}},
		},)
		
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(result)

	c.JSON(newQuiz)

}

func (db *DB) DeleteQuiz(c *fiber.Ctx){
	quizCollection := db.Client.Database("test").Collection("quizzes")

	id := c.Params("id")
	formattedId, _ := primitive.ObjectIDFromHex(id)

	_, err := quizCollection.DeleteOne(context.Background(), bson.M{"_id": formattedId})

	if err != nil{
		log.Fatal(err)
	}

	c.Send("Quiz was succesfully deleted")
}

func (db *DB) GetUser(c *fiber.Ctx){
	userCollection := db.Client.Database("test").Collection("users")

	id := c.Params("id")
	formattedId, _ := primitive.ObjectIDFromHex(id)

	var user *User.User

	if err := userCollection.FindOne(context.Background(), bson.M{"_id": formattedId}).Decode(&user); err != nil{
		log.Fatal(err)
	}

	c.JSON(user)

}


func (db *DB) AswerQuestion(c *fiber.Ctx){
	quizCollection := db.Client.Database("test").Collection("quizzes")
	quiz := new(Quiz.Quiz)
	quizId, _ := primitive.ObjectIDFromHex(c.Params("quizId"))
	answerBody := new(QuestionAnswer) 

	quizCollection.FindOne(context.Background(),bson.M{"_id": quizId}).Decode(&quiz)

	if err := c.BodyParser(answerBody); err != nil {
		log.Fatal(err)
	}

	questions := quiz.Questions
	var rightQuestion *Quiz.Question
	for i := range(questions){
		if questions[i].ID == *answerBody.QuestionID{
			rightQuestion = questions[i]
		}
	}
	fmt.Println(reflect.TypeOf(&rightQuestion.Answer) == reflect.TypeOf(answerBody.Answer))
	if rightQuestion.Answer == *answerBody.Answer{
		responseMap := map[string]bool{rightQuestion.Text: true}
		
		c.JSON(responseMap)
	} else {
		responseMap := map[string]bool{rightQuestion.Text: false}

		c.JSON(responseMap)
	}
}

func (db *DB) SaveResult(c *fiber.Ctx){
	id := c.Locals("id")
	fmt.Println(id)
	quizCollection := db.Client.Database("test").Collection("quizzes")
	userCollection := db.Client.Database("test").Collection("users")

	user := new(User.User)
	var requestBody *SaveResultBody = new(SaveResultBody)

	userCollection.FindOne( context.Background(), bson.M{"_id": id}).Decode(&user)
	fmt.Println(requestBody)
	if err := c.BodyParser(requestBody); err != nil {
		log.Fatal(err)
	}
	var insertResultVar *Quiz.PlayedBy = &Quiz.PlayedBy{Name: *user.Name, RightAmount: *requestBody.RightAmount, Date: time.Now()}
	
	quizCollection.FindOneAndUpdate( context.Background(), bson.M{"_id": requestBody.QuizId}, bson.D{
		{"$push", bson.D{{"playedBy", insertResultVar}}},
	})
	userCollection.UpdateOne( context.Background(), bson.M{"_id": id}, bson.D{
		{"$push", bson.D{{"playedQuizzes", requestBody.QuizId}}},
	})

	c.Send("This was done")
}

func (db *DB) Login(c *fiber.Ctx){
	userCollection := db.Client.Database("test").Collection("users")
	user := new(User.User)	
	notLoggedInUser := new(User.User)

	if err := c.BodyParser(notLoggedInUser); err != nil{
		log.Fatal(err)
	}

	if err := userCollection.FindOne(context.Background(), bson.M{"email": notLoggedInUser.Email}).Decode(user); err != nil{
		log.Fatal(err)
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*notLoggedInUser.Password)); err != nil{
		log.Fatal(err)
	}

	

	claims := jwt.MapClaims{}
	claims["_id"] = user.ID
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encodedToken, err := token.SignedString([]byte(Key))
	if err != nil{
		log.Fatal(err)
	}

	c.Send(encodedToken)
}

func (db *DB) Register(c *fiber.Ctx){
	userCollection := db.Client.Database("test").Collection("users")
	
	emptyArray := make([]primitive.ObjectID, 0, 0)

	newUser := new(User.User)
	newUser.CreatedQuizzes = &emptyArray
	newUser.PlayedQuizzes = &emptyArray

	if err := c.BodyParser(newUser); err != nil{
		log.Fatal(err)
	}

	formattedPassword := []byte(*newUser.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(formattedPassword, bcrypt.DefaultCost)

	if err != nil{
		log.Fatal(err)
	}
	passwordHashString := string(hashedPassword)

	newUser.Password = &passwordHashString

	userCollection.InsertOne(context.Background(), newUser)

	c.JSON(newUser)
}