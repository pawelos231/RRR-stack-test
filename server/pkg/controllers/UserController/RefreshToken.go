package UserController

import (
	"BackendGo/pkg/auth"
	"BackendGo/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func RefreshTokenHandler(col *mongo.Collection, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		user := &models.User{}
		json.NewDecoder(req.Body).Decode(&user)

		tokenCookie2, errCookie := req.Cookie("jid")
		fmt.Println(tokenCookie2.Value, "tokenCookie2")

		value := tokenCookie2.Value
		if errCookie != nil {
			fmt.Println(errCookie)
			json.NewEncoder(w).Encode("coś poszło nie tak")
		}

		tkClaims := jwt.MapClaims{}
		refreshToken, errParsed := jwt.ParseWithClaims(value, tkClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
		})
		//Pass it to utils later beacuse it duplicates
		if errParsed != nil && refreshToken == nil {
			if errParsed == jwt.ErrSignatureInvalid {
				fmt.Println(errParsed)
				json.NewEncoder(w).Encode("invalid token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(errParsed)
			json.NewEncoder(w).Encode("zły token")
			println("blad")
			return
		}

		_, _, User := auth.FindUserByEmail(col, user, ctx)

		RefreshTokenString, _ := auth.CreateRefreshToken(User)
		auth.SendRefreshToken(w, RefreshTokenString)
		tokenString, _ := auth.CreateAccessToken(User)

		var UserInfo = map[string]interface{}{}
		UserInfo["token"] = tokenString
		UserInfo["email"] = User.Email
		UserInfo["name"] = User.Name
		var reponse = map[string]interface{}{"UserInfo": UserInfo}
		reponse["text"] = "Udało się zalogować !"
		reponse["status"] = 1
		fmt.Println(reponse)
		json.NewEncoder(w).Encode(reponse)
	}
}
