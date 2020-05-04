package user

import (
	"net/http"
	"strconv"

	"github.com/cesarnono/bookstore_users-api/domain/users"
	"github.com/cesarnono/bookstore_users-api/services"
	"github.com/cesarnono/bookstore_users-api/utils/errors"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user users.User
	// First way to handle request incomming=====
	// bytes, err := ioutil.ReadAll(c.Request.Body)
	// fmt.Println(string(bytes))
	// if err != nil {
	// 	//TODO: handle error
	// 	return
	// }
	// if err := json.Unmarshal(bytes, &user); err != nil {
	// 	//TODO: handle error json
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// Second way====
	if err := c.ShouldBindJSON(&user); err != nil {
		restError := errors.NewBadRequestError("invalid json body")
		c.JSON(restError.Status, restError)
		return
	}
	result, saveError := services.CreateUser(user)
	if saveError != nil {
		c.JSON(saveError.Status, saveError)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		error := errors.NewBadRequestError("User id should be a number")
		c.JSON(error.Status, error)
		return
	}

	user, getError := services.Get(userId)
	if getError != nil {
		c.JSON(getError.Status, getError)
		return
	}
	c.JSON(http.StatusOK, user)

}
