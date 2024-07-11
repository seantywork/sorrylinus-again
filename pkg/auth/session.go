package auth

import (
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/seantywork/sorrylinus-again/pkg/dbquery"
)

func WhoAmI(c *gin.Context) (string, string, string) {

	var your_key string = ""

	var your_type string = ""

	var your_id string = ""

	session := sessions.Default(c)

	var session_id string

	v := session.Get("SOLIAGAIN")

	if v == nil {

		fmt.Printf("you: %s\n", "nobody-yet")
		return "", "", ""

	} else {

		session_id = v.(string)

	}

	ss, err := dbquery.GetBySessionKeyFromSession(session_id)

	if err != nil {

		fmt.Printf("you: %s\n", "nodody-no-session")

		return "", "", ""
	}

	your_key = session_id
	your_type = ss.Type
	your_id = ss.Id

	return your_key, your_type, your_id
}
