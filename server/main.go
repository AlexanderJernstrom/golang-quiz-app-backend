package main

import (
	database "server/db"
	"server/middlewares"

	"github.com/gofiber/fiber"
)





func main(){
	dbVar := database.ConnectToDatabase()


	app := fiber.New()
	app.Use("/createQuiz", middlewares.IsAuth())
	//Quiz routes
	app.Get("/quizzes", dbVar.GetAllQuizzes)
	app.Get("/quiz/:id", dbVar.GetQuiz)
	app.Post("/createQuiz", dbVar.CreateQuiz)
	app.Delete("/deleteQuiz/:id", dbVar.DeleteQuiz)
	app.Post("/answer/question/:quizId", dbVar.AswerQuestion)
	//User routes
	app.Get("/user/:id", dbVar.GetUser)
	app.Post("/register", dbVar.Register)
	app.Post("/login", dbVar.Login)
	
	app.Listen(":4000")
	
}