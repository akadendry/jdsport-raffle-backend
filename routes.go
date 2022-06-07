package jdsport_raffle_backend

import (

	// "github.com/akadendry/jdsport-raffle-backend/v1middlewares"

	"github.com/akadendry/jdsport-raffle-backend/v1/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/loginmember", controllers.LoginMember)

	app.Post("/api/logout", controllers.Logout)
	app.Post("/api/sendResetPasswordAdmin", controllers.ResetPasswordAdmin)

	app.Get("/api/raffles", controllers.AllRaffles)
	// app.Get("/api/raffles/:id", controllers.GetRaffle)

	app.Post("/api/getRafflesById", controllers.GetRaffleById)
	app.Post("/api/getRafflesBySlug", controllers.GetRaffleBySlug)
	app.Post("/api/getRaffleProductByRaffleId", controllers.GetRaffleProductByRaffleId)
	app.Post("/api/getRaffleProductById", controllers.GetRaffleProductById)

	// app.Post("/api/users", controllers.CreateUser)

	app.Post("/api/users", controllers.CreateUserNew)

	app.Post("/api/getProductSizeStock", controllers.GetProductSizeStock)

	app.Post("/api/finalRegistration", controllers.CreateParticipant)

	app.Post("/api/getUserParticipantYet", controllers.CekUserParticipant)

	// app.Use(middlewares.IsAuthenticated)
	// app.Post("/api/raffles", controllers.CreateRaffle)

	app.Post("/api/cms/createUserAdmin", controllers.CreateUserAdmin)
	app.Post("/api/cms/getUserByToken", controllers.GetUserByToken)
	app.Post("/api/cms/updateUserAdmin", controllers.UpdatePasswordAdmin)

	app.Post("/api/cms/addRaffle", controllers.AddRaffle)
	app.Post("/api/cms/editRaffle", controllers.EditRaffle)
	app.Post("/api/cms/getRaffleByFilter", controllers.GetRaffleByFilter)
	app.Delete("/api/cms/deleteRaffles/:id", controllers.DeleteRaffle)
	app.Post("/api/cms/getTotalWinerParticipanAndTotalStock", controllers.GetTotalWinerParticipanAndTotalStock)

	app.Get("/api/cms/raffle/getSelectOptionRaffle", controllers.GetSelectOptionRaffle)
	app.Post("/api/cms/RaffleParticipant", controllers.RaffleParticipant)
	app.Post("/api/cms/RaffleParticipant", controllers.RaffleParticipant)
	app.Post("/api/cms/raffleSearchParticipant", controllers.RaffleSearchParticipant)
	app.Post("/api/cms/getAllModuleSearchParticipant", controllers.GetAllModuleSearchParticipant)
	app.Post("/api/cms/generateWinerParticipant", controllers.GenerateWinerParticipant)
	app.Post("/api/cms/raffle/participantWinner", controllers.RaffleParticipantWinner)
	app.Post("/api/cms/raffle/participantUpdateStatus", controllers.RaffleParticipantUpdateStatus)
	app.Post("/api/cms/setDatePay", controllers.SetDatePay)
	app.Post("/api/cms/emailBlastWinner", controllers.EmailBlastWinner)
	app.Post("/api/cms/sendEmailPerWinner", controllers.SendEmailPerWinner)

	app.Get("/api/cms/getAllAdminUsers", controllers.AllAdminUsers)
	app.Get("/api/cms/getAllCustomerUsers", controllers.AllCustomerUsers)

	app.Get("/api/cms/roles", controllers.AllRoles)
	app.Post("/api/cms/roles", controllers.CreateRole)
	app.Get("/api/cms/roles/:id", controllers.GetRole)
	app.Put("/api/cms/roles/:id", controllers.UpdateRole)
	app.Delete("/api/cms/roles/:id", controllers.DeleteRole)

	// app.Put("/api/users/info", controllers.UpdateInfo)
	// app.Put("/api/users/password", controllers.UpdatePassword)

	// app.Get("/api/user", controllers.User)
	// app.Post("/api/logout", controllers.Logout)

	// app.Get("/api/users", controllers.AllUsers)

	// app.Get("/api/users/:id", controllers.GetUser)
	// app.Put("/api/users/:id", controllers.UpdateUser)
	// app.Delete("/api/users/:id", controllers.DeleteUser)

	// app.Get("/api/permissions", controllers.AllPermissions)

	// app.Put("/api/raffles/:id", controllers.UpdateRaffle)
	// app.Delete("/api/raffles/:id", controllers.DeleteRaffle)

	// app.Get("/api/products", controllers.AllProducts)
	// app.Post("/api/products", controllers.CreateProduct)
	// app.Get("/api/products/:id", controllers.GetProduct)
	// app.Put("/api/products/:id", controllers.UpdateProduct)
	// app.Delete("/api/products/:id", controllers.DeleteProduct)

	app.Post("/api/cms/upload", controllers.Upload)
	app.Post("/api/cms/uploadMobile", controllers.UploadMobile)
	app.Static("/asset/products", "./asset/products")
	app.Static("/asset/products/mobile", "./asset/products/mobile")

	//CMS
	//Raflle
	app.Post("/api/cms/raffle/getDetailById", controllers.GetDetailRaffleById)
	app.Post("/api/cms/raffle/GetDetailRaffleByIdNew", controllers.GetDetailRaffleByIdNew)
	//Participant
	app.Get("/api/cms/Participant/ListAll", controllers.GetAllParticipant)
	app.Post("/api/cms/Participant/UpdateStatus", controllers.UpdateStatusParticipant)

	//untuk diconsume team backend jd
	app.Get("/api/raffle/ListRaffleActive", controllers.GetRaffleByDateNow)
	app.Post("/api/raffle/CheckUserWinner", controllers.CheckUserWinner)
	app.Post("/api/raffle/UpdateTransactionStatus", controllers.UpdateTransactionStatus)

	// app.Get("/api/orders", controllers.AllOrders)
	// app.Post("/api/export", controllers.Export)
	// app.Get("/api/chart", controllers.Chart)
}
