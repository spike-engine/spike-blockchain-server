package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestConfInit(t *testing.T) {
	fmt.Println(os.Getenv("MORALIS_KEY"))
}
