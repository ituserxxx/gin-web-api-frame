package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const Secret = "UTEK123#@"

type claims struct {
	ID    uint
	Phone string
	jwt.StandardClaims
}

// Gin 中间件函数，用于验证 JWT Token
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		splitAuth := strings.Fields(auth)

		if auth == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("invalid user"))
			c.Abort()
			return
		}
		claims, errParse := jwtParse(splitAuth[len(splitAuth)-1])
		if errParse != nil {
			c.AbortWithError(http.StatusUnauthorized, errParse)
			return
		}
		userId := claims.ID

		c.Set("uid", userId)
		if userId <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": "500",
				"msg":  "验证错误",
			})
			c.Abort()
			return

		} else {
			//用户验证--获取到这个id的用户名
			//userName, err := db.GetUserById(userId)

			//if err != nil {
			//	c.JSON(http.StatusUnauthorized, gin.H{
			//		"code": "500",
			//		"msg":  "验证错误",
			//	})
			//	c.Abort()
			//
			//}
			//if userName == "" {
			//	c.JSON(http.StatusUnauthorized, gin.H{
			//		"code": "500",
			//		"msg":  "验证错误",
			//	})
			//	//中止请求处理
			//	c.Abort()
			//}
		}
	}
}

// 生成 JWT Token 的函数
func JwtSign(ID uint, Phone string) (tokenString string, err error) {
	// The token content.
	// iss: （Issuer）签发者
	// iat: （Issued At）签发时间，用Unix时间戳表示
	// exp: （Expiration Time）过期时间，用Unix时间戳表示
	// aud: （Audience）接收该JWT的一方
	// sub: （Subject）该JWT的主题
	// nbf: （Not Before）不要早于这个时间
	// jti: （JWT ID）用于标识JWT的唯一ID
	claims := claims{
		ID,
		Phone,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    "UTEK",
		},
	}
	//fmt.Printf("jwt:%#v", claims)
	//fmt.Println("jwt:", userId)
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(Secret))
	//fmt.Println(tokenString)
	return
}

// 解析 JWT Token 的函数
func jwtParse(tokenString string) (*claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// 跨域
func CrossSite(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注册CORS中间件
		r.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
}
