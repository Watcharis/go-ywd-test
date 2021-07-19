package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/users"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
	Email  string `json:"email"`
	*jwt.StandardClaims
}

type Usersrepository struct {
	repository users.UserRepository
	caching    *redis.Client
}

func NewUserService(repository users.UserRepository, caching *redis.Client) *Usersrepository {
	return &Usersrepository{
		repository: repository,
		caching:    caching,
	}
}

func (r *Usersrepository) Register(ctx echo.Context) error {

	bodyUser := users.UserRequestRegister{}

	if reqbody := ctx.Bind(&bodyUser); reqbody != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: reqbody.Error(), Status: "fall", Data: ""})
	}
	if err := ctx.Validate(&bodyUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: err.Error(), Status: "fall", Data: ""})
	}
	fmt.Println("bodyUser :", bodyUser)
	checkQuery, err := r.repository.ValidateDataUserByEmail(bodyUser.Email)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}

	if len(checkQuery) == 0 {
		getRoleId, err := r.repository.GetRoleIdByRoleName(bodyUser.RoleName)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		dataGenUser := model.UsersCreate{
			Email:       bodyUser.Email,
			Password:    bodyUser.Password,
			Phone:       bodyUser.Phone,
			DisplayName: bodyUser.DisplayName,
			RoleId:      getRoleId,
		}

		register, err := r.repository.RegisterUser(dataGenUser)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: register, Status: "success", Data: ""})
	}
	return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: "User is exists", Status: "fail", Data: ""})
}

func (r *Usersrepository) Login(ctx echo.Context) error {
	bodyUser := users.UserRequestLogin{}
	if reqbody := ctx.Bind(&bodyUser); reqbody != nil {
		return reqbody
	}
	if err := ctx.Validate(&bodyUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: err.Error(), Status: "fall", Data: ""})
	}
	// fmt.Println("bodyUser :", bodyUser)
	data := users.UserRequestLogin{
		Email:    bodyUser.Email,
		Password: bodyUser.Password,
	}
	valiDateUser, err := r.repository.ValidateUserLogin(data)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}
	if len(valiDateUser) != 0 {

		// fmt.Println("now :", time.Now())
		// fmt.Println("after 1 hour :", time.Now().Add(time.Hour*1))

		//TODO get client from caching
		keys, err := r.caching.Keys("token:" + valiDateUser[0].UserId + ":" + valiDateUser[0].Email).Result()
		if err != nil {
			logrus.Error("Error redis error ->", err)
		}
		// s := strings.Split(keys[0], ":")
		if len(keys) > 0 {

			cachingToken, err := r.caching.Get(keys[0]).Result()
			if err != nil {
				logrus.Error("Error redis error ->", err)
				return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
			}

			//TODO check token expire
			validateToken, err := jwt.Parse(cachingToken, func(*jwt.Token) (interface{}, error) { return []byte("goywdtest"), nil })
			if err != nil {
				logrus.Errorln("err validate Token ->", err.Error())
				r.caching.Del("token:" + valiDateUser[0].UserId)
				return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
			}

			token := validateToken.Claims
			convertToByteArray, _ := json.Marshal(token)
			var dataStructToken Claims
			if err := json.Unmarshal(convertToByteArray, &dataStructToken); err != nil {
				return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
			}

			//ยัง bug อยู่
			if dataStructToken.ExpiresAt > time.Now().Unix() {
				claims := &Claims{
					valiDateUser[0].UserId,
					valiDateUser[0].RoleId,
					valiDateUser[0].Email,
					&jwt.StandardClaims{
						ExpiresAt: dataStructToken.ExpiresAt + (int64(time.Hour) * 1),
						IssuedAt:  time.Now().Unix(),
					},
				}
				//TODO ENCODE PAYLOAD TOKEN
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				//TODO GEN REFRESH_TOKEN
				t, err := token.SignedString([]byte("goywdtest"))
				if err != nil {
					return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
				}

				//SET KEYS REDIS
				if err = r.caching.Set("token:"+valiDateUser[0].UserId+":"+valiDateUser[0].Email, t, time.Duration(int64(time.Minute)*5)).Err(); err != nil {
					logrus.Errorln("redis Error ->", err)
					return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
				}
				cachingNewRefreshToken, err := r.caching.Get("token:" + valiDateUser[0].UserId + ":" + valiDateUser[0].Email).Result()
				if err != nil {
					logrus.Error("Error redis error ->", err)
					return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
				}

				if t == cachingNewRefreshToken {
					return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "login success", Status: "success", Data: users.LoginResponse{Accesstoken: t}})
				}
				return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: "caching unauthorize", Status: "fail", Data: ""})
			}
			return ctx.JSON(http.StatusUnauthorized, model.JsonResponse{Message: "TOKEN is expire", Status: "fail", Data: ""})
		}

		// gen token JWT
		claims := &Claims{
			valiDateUser[0].UserId,
			valiDateUser[0].RoleId,
			valiDateUser[0].Email,
			&jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// fmt.Println("token :", token)
		t, err := token.SignedString([]byte("goywdtest"))
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		//SET KEYS REDIS
		if err = r.caching.Set("token:"+valiDateUser[0].UserId+":"+valiDateUser[0].Email, t, time.Duration(int64(time.Minute)*5)).Err(); err != nil {
			logrus.Errorln("redis Error ->", err)
			return ctx.JSON(http.StatusBadGateway, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		//TODO get DATA Redis
		// var getAllKeysInCaching []string
		// getAllKeysInCaching = TestGetRedisClient(r.caching)
		// fmt.Println("getAllKeysInCaching :", getAllKeysInCaching)

		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "login success", Status: "success", Data: users.LoginResponse{Accesstoken: t}})
	} else {
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: "Not Found Users", Status: "fail", Data: ""})
	}
}

func TestGetRedisClient(rdb *redis.Client) []string {
	keysAll, err := rdb.Keys("token:*").Result()
	if err != nil {
		logrus.Errorln("Error redis command KEYS ->", err)
	}

	for _, v := range keysAll {
		checkToken, err := rdb.Get(v).Result()
		if err != nil {
			logrus.Errorln("Error: redis command get ->", err)
		}
		fmt.Println("checkToken :", checkToken)
	}
	return keysAll
}
