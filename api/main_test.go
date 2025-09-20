package api

import (
	"os"
	"testing"
	"time"

	db "github.com/ayushrakesh/gopay/db/sqlc"
	"github.com/ayushrakesh/gopay/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func testServer(t *testing.T, store db.Store) *Server {

	config := util.Config{
		SymmetricKey:      util.RandomString(32),
		AccessTokenExpiry: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
