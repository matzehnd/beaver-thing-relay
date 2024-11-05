package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struktur für die Antwort von der externen API
type TokenResponse struct {
	Valid bool `json:"valid"`
}

// Middleware für den Token-Check
func TokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Holen des Authorization-Headers
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Überprüfen, ob der Header das Bearer-Präfix hat
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Extrahieren des Tokens
		token := authHeader[7:] // Token ohne "Bearer " Präfix

		// Dekodieren des Tokens (Base64)
		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Erstelle die Anfrage an die externe API
		reqBody, _ := json.Marshal(map[string]string{"token": string(decodedToken)})
		req, err := http.NewRequest("POST", "https://external-api.com/validate", bytes.NewBuffer(reqBody))
		if err != nil {
			log.Println("Error creating request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Sende die Anfrage
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error making request to external API:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Lese die Antwort von der externen API
		body, _ := ioutil.ReadAll(resp.Body)
		var tokenResp TokenResponse

		// Überprüfe den Statuscode der Antwort
		if resp.StatusCode == http.StatusOK {
			// Unmarshale den Body nur, wenn der Status OK ist
			if err := json.Unmarshal(body, &tokenResp); err != nil {
				log.Println("Error unmarshaling response:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}

			// Überprüfe, ob der Token gültig ist
			if !tokenResp.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
		} else {
			// Handle den Fall, wenn die API nicht mit 200 antwortet
			log.Println("External API returned non-200 status:", resp.Status)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating token"})
			c.Abort()
			return
		}

		// Token ist gültig, setze die nächste Middleware oder den Handler fort
		c.Next()
	}
}
