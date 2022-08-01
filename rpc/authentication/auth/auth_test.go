package auth

import (
	"fmt"
	"llfile/rpc"
	"testing"
)

func TestAuth(t *testing.T) {
	authClient := rpc.NewAuth()
	token, err := authClient.CallCreateToken(111)
	fmt.Println(token, err)

}

func TestClient_CallAuthToken(t *testing.T) {
	authClient2 := rpc.NewAuth()
	authToken, err := authClient2.CallAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMTEsInVzZXJfaWTnmoRBRVPmt7fmt4YiOiJTVjUvamdXUlYxdlBKR0cvd3RMU3R3PT0iLCJleHAiOjE2NTgxMzcyMDR9.0gIqiEVcVfv9AEJC2sLAmqEOYQH2mFi8CDbpFzIxONo", 111)
	fmt.Println(authToken, err)
}
