package repository
import(
	"example.com/project/models"
	"go.uber.org/zap"
)

var(
	 Users = models.UserList{
		Users: []models.User{
			{
				Username: "Oleg",
				Password: "123456",
				Email:    "olegs@gmail.com",
			},
			{
				Username: "Anna",
				Password: "654321",
				Email:    "anna@example.com",
			},
		},
	}

		Logger   *zap.Logger

)