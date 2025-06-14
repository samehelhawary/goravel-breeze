package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	httpContract "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// --- Contracts and Implementations ---

// DecryptException represents a failure during the decryption process.
type DecryptException struct {
	Err error
}

func (e *DecryptException) Error() string {
	return fmt.Sprintf("decryption failed: %v", e.Err)
}

// EncrypterContract defines the interface for an encryption service.
type EncrypterContract interface {
	Encrypt(value string, serialize bool) (string, error)
	Decrypt(payload string, unserialize bool) (string, error)
	GetKey() string
	GetAllKeys() []string
}

// AesEncrypter provides an AES-256 GCM implementation of the EncrypterContract.
type AesEncrypter struct {
	key []byte
}

// NewAesEncrypter creates a new AES encrypter. The key must be 32 bytes.
func NewAesEncrypter(key string) (*AesEncrypter, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}
	return &AesEncrypter{key: []byte(key)}, nil
}

// Encrypt encrypts a string using AES-256 GCM.
func (e *AesEncrypter) Encrypt(value string, serialize bool) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(value), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a string using AES-256 GCM.
func (e *AesEncrypter) Decrypt(payload string, unserialize bool) (string, error) {
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", &DecryptException{Err: err}
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", &DecryptException{Err: err}
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", &DecryptException{Err: err}
	}
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", &DecryptException{Err: errors.New("ciphertext is too short")}
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", &DecryptException{Err: err}
	}
	return string(plaintext), nil
}

func (e *AesEncrypter) GetKey() string {
	return string(e.key)
}

func (e *AesEncrypter) GetAllKeys() []string {
	return []string{string(e.key)}
}

// CookieValuePrefix handles adding and validating a signed HMAC prefix to cookie values.
type CookieValuePrefix struct {
	key string
}

func NewCookieValuePrefix(key string) *CookieValuePrefix {
	return &CookieValuePrefix{key: key}
}

func (c *CookieValuePrefix) Create(name string) string {
	h := hmac.New(sha256.New, []byte(c.key))
	h.Write([]byte(name + "v2"))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)) + "|"
}

func (c *CookieValuePrefix) Validate(name, value string, allKeys []string) (string, error) {
	parts := strings.SplitN(value, "|", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid cookie prefix format")
	}
	prefix, actualValue := parts[0], parts[1]
	for _, key := range allKeys {
		validator := NewCookieValuePrefix(key)
		expectedPrefix := strings.TrimSuffix(validator.Create(name), "|")
		if hmac.Equal([]byte(prefix), []byte(expectedPrefix)) {
			return actualValue, nil
		}
	}
	return "", errors.New("invalid cookie prefix signature")
}

// --- EncryptCookies Middleware ---

// EncryptCookies is middleware for encrypting and decrypting HTTP cookies.
type EncryptCookies struct {
	encrypter EncrypterContract
	except    map[string]bool
}

var (
	neverEncrypt sync.Map
	serialize    bool
)

// NewEncryptCookies creates a new EncryptCookies middleware instance.
// It automatically initializes the encrypter from the application's configuration.
func NewEncryptCookies() *EncryptCookies {
	appKey := facades.Config().GetString("APP_KEY")
	aesEncrypter, err := NewAesEncrypter(appKey)
	if err != nil {
		// Panic because this is a critical configuration error.
		panic(fmt.Sprintf("FATAL: Failed to create cookie encrypter: %v", err))
	}

	return &EncryptCookies{
		encrypter: aesEncrypter,
		except:    make(map[string]bool),
	}
}

// DisableFor adds cookie names to the 'except' list for this middleware instance
// and returns the middleware instance for fluent chaining.
func (m *EncryptCookies) DisableFor(names ...string) *EncryptCookies {
	for _, name := range names {
		m.except[name] = true
	}
	return m
}

// Handle processes the HTTP request and response.
func (m *EncryptCookies) Handle() httpContract.Middleware {
	return func(ctx httpContract.Context) {
		m.decryptRequestCookies(ctx.Request().Origin())
		ctx.Request().Next()
		m.encryptResponseCookies(ctx)
	}
}

func (m *EncryptCookies) decryptRequestCookies(r *http.Request) {
	if len(r.Cookies()) == 0 {
		return
	}
	decryptedCookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		if m.isDisabled(cookie.Name) {
			decryptedCookies[cookie.Name] = cookie.Value
			continue
		}
		decryptedValue, err := m.decryptCookie(cookie.Name, cookie.Value)
		if err != nil {
			fmt.Printf("Warning: could not decrypt cookie '%s': %v\n", cookie.Name, err)
			continue
		}
		prefixer := NewCookieValuePrefix(m.encrypter.GetKey())
		validatedValue, err := prefixer.Validate(cookie.Name, decryptedValue, m.encrypter.GetAllKeys())
		if err != nil {
			fmt.Printf("Warning: invalid prefix for cookie '%s': %v\n", cookie.Name, err)
			continue
		}
		decryptedCookies[cookie.Name] = validatedValue
	}
	r.Header.Del("Cookie")
	var newCookieStrings []string
	for name, value := range decryptedCookies {
		newCookieStrings = append(newCookieStrings, fmt.Sprintf("%s=%s", name, value))
	}
	if len(newCookieStrings) > 0 {
		r.Header.Set("Cookie", strings.Join(newCookieStrings, "; "))
	}
}

func (m *EncryptCookies) decryptCookie(name, value string) (string, error) {
	return m.encrypter.Decrypt(value, m.serialized(name))
}

func (m *EncryptCookies) encryptResponseCookies(ctx httpContract.Context) {
	responseHeaders := ctx.Response().Origin().Header()
	setCookieHeaders := append([]string{}, responseHeaders["Set-Cookie"]...)
	if len(setCookieHeaders) == 0 {
		return
	}
	responseHeaders.Del("Set-Cookie")
	for _, cookieStr := range setCookieHeaders {
		m.handleOutgoingCookie(ctx, cookieStr)
	}
}

func (m *EncryptCookies) handleOutgoingCookie(ctx httpContract.Context, cookieStr string) {
	dummyRes := http.Response{Header: http.Header{"Set-Cookie": {cookieStr}}}
	parsedCookies := dummyRes.Cookies()
	if len(parsedCookies) == 0 {
		return
	}
	cookie := parsedCookies[0]
	if cookie.MaxAge < 0 {
		ctx.Response().Origin().Header().Add("Set-Cookie", cookieStr)
		return
	}
	if m.isDisabled(cookie.Name) {
		ctx.Response().Cookie(toContractCookie(cookie))
		return
	}
	prefixer := NewCookieValuePrefix(m.encrypter.GetKey())
	valueToEncrypt := prefixer.Create(cookie.Name) + cookie.Value
	encryptedValue, err := m.encrypter.Encrypt(valueToEncrypt, m.serialized(cookie.Name))
	if err != nil {
		fmt.Printf("Warning: could not encrypt cookie '%s': %v\n", cookie.Name, err)
		return
	}
	encryptedCookie := toContractCookie(cookie)
	encryptedCookie.Value = encryptedValue
	ctx.Response().Cookie(encryptedCookie)
}

func (m *EncryptCookies) isDisabled(name string) bool {
	if _, ok := m.except[name]; ok {
		return true
	}
	_, ok := neverEncrypt.Load(name)
	return ok
}

func (m *EncryptCookies) serialized(name string) bool {
	return serialize
}

func toContractCookie(c *http.Cookie) httpContract.Cookie {
	return httpContract.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  c.Expires,
		MaxAge:   c.MaxAge,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
	}
}
