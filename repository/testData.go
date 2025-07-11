package repository
import(
	"example.com/project/models"
	"go.uber.org/zap"
	"github.com/coreos/go-oidc"
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
			{
				Username: "Ivars",
				Password: "123490",
				Email:    "ivars@example.com",
			},

		},
	}

	Logger   *zap.Logger

    Verifier *oidc.IDTokenVerifier

    NonceStore = map[string]bool{}
)