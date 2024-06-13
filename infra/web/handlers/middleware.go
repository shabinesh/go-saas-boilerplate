package handlers

import "github.com/gin-gonic/gin"

// gin auth middleware
func (h handlers) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, _ := h.sessionStore.Get(c.Request, sessionKey)
		if v, ok := sess.Values["authenticated"]; !ok || v != true {
			c.Abort()
			c.Redirect(302, "/login")
		} else {
			c.Next()
			return
		}

		c.Abort()
		c.Redirect(302, "/login")
	}
}
