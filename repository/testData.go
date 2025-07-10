package repository
import(
	"example.com/project/models"
	"go.uber.org/zap"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
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

	oauth2Config *oauth2.Config

    Verifier *oidc.IDTokenVerifier

    NonceStore = map[string]bool{}
)