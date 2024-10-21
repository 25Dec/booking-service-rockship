package middleware

import (
	"booking-service/internal/openapis"

	"github.com/gin-gonic/gin"
	"fmt"
)

type LearnerMiddleware struct {
	edtronaut openapis.EdtronautAPI
}

func NewVerifyLearner(edtronaut openapis.EdtronautAPI) *LearnerMiddleware {
	return &LearnerMiddleware{edtronaut: edtronaut}
}

func (v *LearnerMiddleware) VerifyLearner() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get("Authorization")
		if len(authorizationHeader) < 7 || authorizationHeader[7:] == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		//get token form header bearer
		token := authorizationHeader[7:]

		learner, err := v.edtronaut.GetUserByToken(token)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("learner", learner)

		// userCourses, err := v.edtronaut.GetUserCourses(token)
		// if err != nil {
		// 	c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		// 	return
		// }

		// hasCourse := false

		// for _, course := range userCourses {
		// 	if course.AllowBooking {
		// 		hasCourse = true
		// 		break
		// 	}
		// }

		// if !hasCourse {
		// 	c.AbortWithStatusJSON(401, gin.H{"error": "User has no access to this course"})
		// 	return
		// }

		c.Next()
	}
}
