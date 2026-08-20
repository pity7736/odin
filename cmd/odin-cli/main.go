package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/crypto/argon2"
)

const (
	defaultBaseURL = "http://localhost:8000/api/v1"
	sessionFile    = ".odin-cli-session.json"
	argonTime      = 3
	argonMemory    = 65536
	argonThreads   = 4
	masterKeySize  = 32
	subkeySize     = 32
	saltSize       = 16
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: odin-cli <register|login> [flags]")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "register":
		if err := runRegister(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "login":
		if err := runLogin(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func runRegister(args []string) error {
	flags := flag.NewFlagSet("register", flag.ExitOnError)
	email := flags.String("email", "", "account email")
	password := flags.String("password", "", "account password (never sent to the server)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *email == "" || *password == "" {
		return fmt.Errorf("email and password are required")
	}
	salt, err := randomBytes(saltSize)
	if err != nil {
		return err
	}
	authHash, encryptionKey := deriveKeys(*password, salt)
	encryptedMasterKey, err := generateWrappedMasterKey(encryptionKey)
	if err != nil {
		return err
	}
	response, err := register(defaultBaseURL, *email, authHash, encryptedMasterKey, salt)
	if err != nil {
		return err
	}
	sessions, err := loadSessions()
	if err != nil {
		return err
	}
	sessions[*email] = session{Salt: base64.StdEncoding.EncodeToString(salt)}
	if err := saveSessions(sessions); err != nil {
		return err
	}
	fmt.Printf("registered id=%s email=%s (session saved to %s)\n", response.ID, response.Email, sessionFile)
	return nil
}

func randomBytes(size int) ([]byte, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return nil, err
	}
	return buffer, nil
}

func deriveKeys(password string, salt []byte) (string, []byte) {
	stretched := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, subkeySize*2)
	authHash := base64.StdEncoding.EncodeToString(stretched[:subkeySize])
	encryptionKey := stretched[subkeySize:]
	return authHash, encryptionKey
}

func generateWrappedMasterKey(encryptionKey []byte) (string, error) {
	masterKey, err := randomBytes(masterKeySize)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce, err := randomBytes(gcm.NonceSize())
	if err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, masterKey, nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func register(baseURL, email, authHash, encryptedMasterKey string, salt []byte) (registerResponse, error) {
	body := registerRequest{
		Email:              email,
		AuthHash:           authHash,
		EncryptedMasterKey: encryptedMasterKey,
		KeyParams: keyParams{
			Algorithm:   "argon2id",
			Salt:        base64.StdEncoding.EncodeToString(salt),
			Iterations:  argonTime,
			Memory:      argonMemory,
			Parallelism: argonThreads,
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return registerResponse{}, err
	}
	httpResponse, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader(payload))
	if err != nil {
		return registerResponse{}, err
	}
	defer func() { _ = httpResponse.Body.Close() }()
	raw, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return registerResponse{}, err
	}
	if httpResponse.StatusCode != http.StatusCreated {
		return registerResponse{}, fmt.Errorf("register failed (%d): %s", httpResponse.StatusCode, raw)
	}
	var response registerResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return registerResponse{}, err
	}
	return response, nil
}

func loadSessions() (map[string]session, error) {
	raw, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]session{}, nil
		}
		return nil, err
	}
	sessions := map[string]session{}
	if err := json.Unmarshal(raw, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func saveSessions(sessions map[string]session) error {
	payload, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sessionFile, payload, 0o600)
}

func runLogin(args []string) error {
	flags := flag.NewFlagSet("login", flag.ExitOnError)
	email := flags.String("email", "", "account email")
	password := flags.String("password", "", "account password (never sent to the server)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *email == "" || *password == "" {
		return fmt.Errorf("email and password are required")
	}
	sessions, err := loadSessions()
	if err != nil {
		return err
	}
	current, ok := sessions[*email]
	if !ok {
		return fmt.Errorf("no session for %s (run register first)", *email)
	}
	salt, err := base64.StdEncoding.DecodeString(current.Salt)
	if err != nil {
		return err
	}
	authHash, encryptionKey := deriveKeys(*password, salt)
	response, err := login(defaultBaseURL, *email, authHash)
	if err != nil {
		return err
	}
	masterKey, err := unwrapMasterKey(encryptionKey, response.EncryptedMasterKey)
	if err != nil {
		return fmt.Errorf("master key unwrap failed: %w", err)
	}
	current.Token = response.Token
	current.MasterKey = base64.StdEncoding.EncodeToString(masterKey)
	sessions[*email] = current
	if err := saveSessions(sessions); err != nil {
		return err
	}
	fmt.Printf("login ok for %s (token + master key saved to %s)\n", *email, sessionFile)
	fmt.Printf("key params: algorithm=%s iterations=%d memory=%d parallelism=%d\n", response.KeyParams.Algorithm, response.KeyParams.Iterations, response.KeyParams.Memory, response.KeyParams.Parallelism)
	fmt.Println("master key recovered — crypto round-trip verified")
	return nil
}

func login(baseURL, email, authHash string) (loginResponse, error) {
	payload, err := json.Marshal(loginRequest{Email: email, AuthHash: authHash})
	if err != nil {
		return loginResponse{}, err
	}
	httpResponse, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(payload))
	if err != nil {
		return loginResponse{}, err
	}
	defer func() { _ = httpResponse.Body.Close() }()
	raw, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return loginResponse{}, err
	}
	if httpResponse.StatusCode != http.StatusCreated {
		return loginResponse{}, fmt.Errorf("login failed (%d): %s", httpResponse.StatusCode, raw)
	}
	var response loginResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return loginResponse{}, err
	}
	return response, nil
}

func unwrapMasterKey(encryptionKey []byte, encryptedMasterKey string) ([]byte, error) {
	sealed, err := base64.StdEncoding.DecodeString(encryptedMasterKey)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(sealed) < gcm.NonceSize() {
		return nil, fmt.Errorf("encrypted master key is too short")
	}
	nonce := sealed[:gcm.NonceSize()]
	ciphertext := sealed[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

type registerRequest struct {
	Email              string    `json:"email"`
	AuthHash           string    `json:"auth_hash"`
	EncryptedMasterKey string    `json:"encrypted_master_key"`
	KeyParams          keyParams `json:"key_params"`
}

type keyParams struct {
	Algorithm   string `json:"algorithm"`
	Salt        string `json:"salt"`
	Iterations  int    `json:"iterations"`
	Memory      int    `json:"memory"`
	Parallelism int    `json:"parallelism"`
}

type registerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type loginRequest struct {
	Email    string `json:"email"`
	AuthHash string `json:"auth_hash"`
}

type loginResponse struct {
	Token              string    `json:"token"`
	EncryptedMasterKey string    `json:"encrypted_master_key"`
	KeyParams          keyParams `json:"key_params"`
}

type session struct {
	Salt      string `json:"salt"`
	Token     string `json:"token"`
	MasterKey string `json:"master_key"`
}
