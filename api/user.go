package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sensor-data-simulator/global"
	"sensor-data-simulator/tools"

	routing "github.com/julienschmidt/httprouter"
)

/*----------------------*/

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	TokenHash string `json:"tokenHash"`

	// LastLogin time.Time `json:"lastlogin"`
}

/*----------------------*/

/*
* This function implements GET /userDevices
 */
func GetUserDevicesAndSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	userId, err := getAuthorizedUserID(resp, req)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*------*/

	user, err := getUserById(userId)
	if err != nil {
		http.Error(resp, "Something went wrong", http.StatusInternalServerError)
		return
	}

	/*------*/

	listOfDevices, err := getUserDevices(user)
	if err != nil {
		http.Error(resp, "Something went wrong", http.StatusInternalServerError)
		return
	}

	/*------*/

	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(listOfDevices))

}

/*-------------*/

/*
* This function queries Waziup API for user's devices
* and returns all of them up to a 1000 records
 */
func getUserDevices(user User) (string, error) {

	var queryURI = fmt.Sprintf(`devices?q=owner==%s&limit=%d`, user.Username, 1000)

	req, err := http.NewRequest("GET", global.ENV.WAZIUP_API_PATH+queryURI, nil)
	if err != nil {
		log.Printf("[FETCH] could not make the request: %v", err)
		return "", CodeError{500, "Something went wrong!"}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.Token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[FETCH] did not receive a response from Waziup Server: %v", err)
		return "", CodeError{500, "Did not receive a response from Waziup Server!"}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", CodeError{resp.StatusCode, resp.Status}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[AUTH ] Failed to read the response content: %v", err)
		return "", CodeError{500, "Something went wrong!"}
	}

	return string(body), nil
}

/*
* This function implements GET /user
 */
func GetUser(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	userId, err := getAuthorizedUserID(resp, req)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*------*/

	user, err := getUserById(userId)
	if err != nil {
		http.Error(resp, "Something went wrong", http.StatusInternalServerError)
		return
	}

	output := map[string]interface{}{
		"username":  user.Username,
		"tokenHash": user.TokenHash,
	}

	tools.SendJSON(resp, output)
}
