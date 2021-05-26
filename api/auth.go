package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sensor-data-simulator/database"
	"sensor-data-simulator/global"
	"sensor-data-simulator/tools"
	"strings"
	"time"

	routing "github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

/*----------------*/

const loginSessionExpTimeMinutes int = 20 // in minutes

const (
	SameSiteDefaultMode http.SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

/*----------------------*/

// PostAuth implements POST /auth
// It returns and sets a tokenHash which is used for authorization in the App
func PostAuth(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	body, err := tools.ReadAll(req.Body)
	if err != nil {
		log.Printf("[ERR  ] PostAuth: %s", err.Error())
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	var inputUser User

	err = json.Unmarshal(body, &inputUser)
	if err != nil {
		log.Printf("[ERR  ] PostAuth: %s", err.Error())
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	// log.Printf("Input User: %q", inputUser)

	token, err := CheckUserCredentials(inputUser.Username, inputUser.Password)

	if err != nil {
		log.Printf("[ERR  ] PostAuth: %s", err.Error())
		http.Error(resp, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	//Login success.

	tokenHash, err := GenerateTokenHash(token)
	if err != nil {
		log.Printf("[ERR  ] PostAuth: %s", err.Error())
		http.Error(resp, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	/*---------*/

	// Save the user creds. in the DB
	err = SaveUserInfo(User{
		Username:  inputUser.Username,
		Password:  inputUser.Password,
		Token:     token,
		TokenHash: tokenHash,
	})

	if err != nil {
		log.Printf("[ERR  ] PostAuth: %s", err.Error())
		http.Error(resp, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	/*---------*/

	// Set Cookie, it is just an extra feature that makes the life easier on the UI part
	expiration := time.Now().Add(time.Minute * time.Duration(loginSessionExpTimeMinutes))
	setAuthCookie(resp, tokenHash, expiration)

	/*---------*/

	resp.Write([]byte(tokenHash))
}

/*---------------------*/

func setAuthCookie(resp http.ResponseWriter, tokenHash string, expirationTime time.Time) {

	cookie := http.Cookie{
		Name:     "TokenHash",
		Value:    string(tokenHash),
		Path:     "/",
		Expires:  expirationTime,
		HttpOnly: true,
		MaxAge:   60 * loginSessionExpTimeMinutes,
		// Secure:     true,
		SameSite: SameSiteStrictMode,
	}
	http.SetCookie(resp, &cookie)

}

/*---------------------*/

// This function checks the user's credential from the Waziup cloud API
// and if the user does not exist in the local database, it adds it.
func CheckUserCredentials(username string, password string) (string, error) {

	var token string

	var postBody = []byte(fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password))

	req, err := http.NewRequest("POST", global.ENV.WAZIUP_API_PATH+`auth/token`, bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("[AUTH ] could not make the request: %v", err)
		return token, CodeError{500, "Something went wrong!"}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[AUTH ] did not receive a response from Waziup Server: %v", err)
		return token, CodeError{500, "Did not receive a response from Waziup Server!"}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return token, CodeError{resp.StatusCode, resp.Status}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[AUTH ] Failed to read the response content: %v", err)
		return token, CodeError{500, "Something went wrong!"}
	}

	token = string(body)

	return token, nil
}

/*---------------------*/

// This function receives an authorized user info, stores it in the database

func SaveUserInfo(user User) error {

	existingUser, err := GetUserByUsername(user.Username)
	if err != nil {
		// The user not found, so let's add it
		row := database.RowType{
			"username":  user.Username,
			"password":  user.Password,
			"token":     user.Token,
			"tokenHash": user.TokenHash,
		}
		insRes, err := global.DB.Insert("users", row)
		if err != nil {
			// Let's ignore duplicate key as some user's have used same sensor name twice or more
			log.Printf("\nError in `users` insertion: %v", err)
			return err
		}
		if insRes.RowsAffected > 0 {
			return nil // Everything is alright
		}
	}

	// Let's update the current user as it exists
	row := database.RowType{
		"password":  user.Password,
		"token":     user.Token,
		"tokenHash": user.TokenHash,
	}
	_, err = global.DB.Update("users", row, database.RowType{"id": existingUser.ID})
	if err != nil {
		log.Printf("\nError in `users` update: %v", err)
		return err
	}

	return nil // Done without error

}

/*---------------------*/

func getAuthorizedUserID(resp http.ResponseWriter, req *http.Request) (int64, error) {

	reqTokenHash := ""

	if req.Header["Authorization"] != nil && len(req.Header["Authorization"][0]) > 0 {

		bearToken := req.Header["Authorization"][0]
		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			reqTokenHash = strArr[1]
		}

	} else {

		c, err := req.Cookie("TokenHash")
		if err != nil {
			// log.Printf("[ERR  ] Auth reading cookie: %s", err.Error())
		} else {
			reqTokenHash = c.Value
		}
	}

	/*---------*/

	if len(reqTokenHash) == 0 {

		return 0, fmt.Errorf("not authorized")
	}

	userId, err := GetUserIdByTokenHash(reqTokenHash)
	if err != nil {
		return 0, err
	}

	/*---------*/

	// Refresh the cookies: this keeps the user logged in while she/he using the service
	expiration := time.Now().Add(time.Minute * time.Duration(loginSessionExpTimeMinutes))
	setAuthCookie(resp, reqTokenHash, expiration)

	/*---------*/

	return userId, nil
}

/*---------------------*/

func GetUserIdByTokenHash(tokenHash string) (int64, error) {

	SQL := `SELECT "id" FROM "users" WHERE "tokenHash" = $1`

	rows, err := global.DB.Query(SQL, database.QueryParams{tokenHash})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		return 0, err
	}

	if rows == nil || len(rows) == 0 {
		return 0, fmt.Errorf("User not found!")
	}

	return rows[0]["id"].(int64), err
}

/*---------------------*/

func GetUserById(userId int64) (User, error) {

	var user User

	SQL := `SELECT * FROM "users" WHERE "id" = $1`

	rows, err := global.DB.Query(SQL, database.QueryParams{userId})
	if err != nil {
		log.Printf("Error in db query: %v , SQL: %s", err, SQL)
		return user, err
	}

	if rows == nil || len(rows) == 0 {
		return user, fmt.Errorf("User not found!")
	}

	user.ID = rows[0]["id"].(int64)
	user.Username = rows[0]["username"].(string)
	user.Password = rows[0]["password"].(string)
	user.Token = rows[0]["token"].(string)
	user.TokenHash = rows[0]["tokenHash"].(string)

	return user, nil
}

/*---------------------*/

func GetUserByUsername(username string) (User, error) {

	var user User

	SQL := `SELECT * FROM "users" WHERE "username" = $1`

	rows, err := global.DB.Query(SQL, database.QueryParams{username})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		return user, err
	}

	if rows == nil || len(rows) == 0 {
		return user, fmt.Errorf("User not found")
	}

	user.ID = rows[0]["id"].(int64)
	user.Username = rows[0]["username"].(string)
	user.Password = rows[0]["password"].(string)
	user.Token = rows[0]["token"].(string)
	user.TokenHash = rows[0]["tokenHash"].(string)

	return user, nil
}

/*---------------------*/

func GenerateTokenHash(token string) (string, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

/*---------------------*/

// PostLogout implements GET /auth/logout
func PostLogout(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	// c := http.Cookie{
	// 	Name:     "TokenHash",
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	// Secure:     true,
	// 	SameSite: SameSiteStrictMode,
	// 	MaxAge:   -1}
	// http.SetCookie(resp, &c)
	expiration := time.Now().Add(time.Minute * time.Duration(-1))
	setAuthCookie(resp, "", expiration)

	//TODO: Other actions that we may need to do in future
	tools.SendJSON(resp, "Logged out.")
}

/*---------------------*/
