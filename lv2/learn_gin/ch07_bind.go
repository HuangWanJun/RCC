package main

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"reflect"
	"strings"
)

var trans ut.Translator

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		fmt.Println(err)
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func InitFrans(locale string) (err error) {
	//注册gin 框架中 validator 引擎属性
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.getTranslator(%s)", locale)
		}
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
		default:
			en_translations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// godoc.org/github.com/go-playground/validator
type LoginForm struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required,min=3,max=10"`
	Password string `from:"password" json:"password" xml:"password" binding:"required"`
}

type SignUpForm struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=40"`
	Name       string `json:"name" binding:"required,min=3"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"repassword" binding:"eqfield=Password"`
}

func main() {
	if err := InitFrans("zh"); err != nil {
		fmt.Println("Ttranslate init has some error")
		return
	}
	router := gin.Default()
	router.POST("/loginJson", func(c *gin.Context) {
		var loginForm LoginForm
		if err := c.ShouldBind(&loginForm); err != nil {
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"msg": err.Error(),
				})
			}
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": removeTopStruct(errs.Translate(trans)),
			})
			return

		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "login success",
		})
	})

	router.POST("/signUP", func(c *gin.Context) {
		var signForm SignUpForm
		if error := c.ShouldBind(&signForm); error != nil {

			_, ok := error.(validator.ValidationErrors)
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"msg": error.Error(),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{})

			fmt.Println("err" + error.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})

	router.Run(":8083")
}
