package signInHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"fullstack-journal/app/config"
	"fullstack-journal/app/entity"
	"fullstack-journal/app/entity/usersEntity"
	"fullstack-journal/app/filters"
	"fullstack-journal/app/helpers/globalFunc"
	"fullstack-journal/app/helpers/validationFunc"
	"fullstack-journal/app/response"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SignIn(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		safeInput, err := vSignIn(c)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": err.Error(),
			})
		}

		// Get User By Email
		users := usersEntity.FindFirstByEmail(db, safeInput["email"].(string))
		if users.ID == uuid.Nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "accountNotFound"),
			})
		}

		// Verify Password
		if err:=globalFunc.VerifyHash(users.Password, safeInput["password"].(string)); err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "accountNotFound"),
			})
		}

		exp := time.Now().Add(time.Hour * 24).Unix() 
		if safeInput["rememberMe"] == 1 {
			exp = time.Now().Add(time.Hour * 48).Unix() 
		}
		
		// payload token
		resUserToken := response.ResUserToken{
			ID: users.ID,
			Name: users.Name,
			Email: users.Email,
			CurrentCompany: response.ResCurrentCompany{},
		}

		resToken := response.ResToken{
			User: resUserToken,
			Exp: exp,
		}

		token, err := globalFunc.GetJwt(resToken)
		if err != nil{
			log.Println("Error generate jwt token", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "networkError"),
			})
		}

		db.Model(&entity.Users{}).Where("id = ?", users.ID).Update("token", token)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"token": token,
				"user": resUserToken,
			},
			"success": true,
		});
	}
}

func SignInGoogle(c echo.Context)error{
	env, err := config.GetEnv()

	if err != nil{
		log.Fatal("Error load env", err)
	}

	googleClientId := env["googleClientId"].(string)
	redirectUrl := env["appUrl"].(string)+"/api/v1/auth/google/callback"
	authUrl := "https://accounts.google.com/o/oauth2/v2/auth"

	params := url.Values{}
	params.Add("client_id", googleClientId)
	params.Add("redirect_uri", redirectUrl)
	params.Add("response_type", "code")
	params.Add("scope", "openid email profile")
	params.Add("prompt", "select_account")
	params.Add("access_type", "offline")

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?%s", authUrl, params.Encode()))
}


func SignGoogleCallback(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		var userJson map[string]interface{}
		var tokenRes map[string]interface{}
		var resJson []byte

		env, err := config.GetEnv()

		if err != nil{
			log.Fatal("Error load env", err)
		}
	
		code := c.QueryParam("code")
		if code == ""{
			resJson, _ = json.Marshal(map[string]interface{}{
				"success": false,
				"error": "Invalid State, Please try again",
			})

			return c.HTML(http.StatusBadRequest, fmt.Sprintf(`
				<script>
					window.opener.postMessage(%s, "*");
					window.close();
				</script>
			`, string(resJson)))
		}
	
		googleClientId := env["googleClientId"].(string)
		googleClientSecrets := env["googleClientSecrets"].(string)
		redirectUrl := env["appUrl"].(string)+"/api/v1/auth/google/callback"
	
		data := url.Values{}
		data.Set("code", code)
		data.Set("client_id", googleClientId)
		data.Set("client_secret", googleClientSecrets)
		data.Set("redirect_uri", redirectUrl)
		data.Set("grant_type", "authorization_code")
	
		resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
		if err != nil {
			return err
		}
	
		defer resp.Body.Close()
	
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &tokenRes)
	
		accessToken := tokenRes["access_token"].(string)
		req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
	
		userResp, _ := http.DefaultClient.Do(req)
		defer userResp.Body.Close()
	
		userBody, _ := io.ReadAll(userResp.Body)
		json.Unmarshal(userBody, &userJson)
	
		users := usersEntity.FindFirstByEmail(db, userJson["email"].(string))

		if users.ID == uuid.Nil {
			resJson, _ = json.Marshal(map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "accountNotFound"),
			})

			return c.HTML(http.StatusBadRequest, fmt.Sprintf(`
				<script>
					window.opener.postMessage(%s, "*");
					window.close();
				</script>
			`, string(resJson)))
		}
	
		exp := time.Now().Add(time.Hour * 24).Unix() 
		// For Token
		// payload token
		resUserToken := response.ResUserToken{
			ID: users.ID,
			Name: users.Name,
			Email: users.Email,
			CurrentCompany: response.ResCurrentCompany{},
		}

		resToken := response.ResToken{
			User: resUserToken,
			Exp: exp,
		}

		token, err := globalFunc.GetJwt(resToken)
		if err != nil{
			resJson, _ = json.Marshal(map[string]interface{}{
				"success": false,
				"error": filters.Translate(c, "networkError"),
			})
		
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf(`
				<script>
					window.opener.postMessage(%s, "*");
					window.close();
				</script>
			`, string(resJson)))
		}

		db.Model(&entity.Users{}).Where("id = ?", users.ID).Update("token", token).Update("google_id", userJson["id"].(string))
		resJson, _ = json.Marshal(map[string]interface{}{
			"data": map[string]interface{}{
				"token": token,
				"user": resUserToken,
			},
			"success": true,
		})
		
		return c.HTML(http.StatusOK, fmt.Sprintf(`
			<script>
			window.opener.postMessage(%s, "*");
			window.close();
			</script>
		`, string(resJson)))
	  
	}
}

func vSignIn(c echo.Context)(map[string]interface{}, error){
	safeInput := make(map[string]interface{})

	email := c.FormValue("email")
	password:=c.FormValue("password")
	rememberMeStr := c.FormValue("rememberMe")

	if err := validationFunc.VEmail(c, email); err != nil{
		return safeInput, err
	}

	if password == ""{
		return safeInput, errors.New(
			filters.Translate(c, "passwordRequired"),
		)
	}

	rememberMe, err := strconv.ParseBool(rememberMeStr)
	if err != nil  {
		return safeInput, errors.New(
			filters.Translate(c, "invalidRememberMe"),
		)
	}

	safeInput = map[string]interface{}{
		"email":      email,
		"password":   password,
		"rememberMe": rememberMe,
	}

	return safeInput, nil

}