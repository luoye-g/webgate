package endecrypt

import (
	"fmt"
	"testing"
)

func TestEndecrypt(t *testing.T) {

	plaintext := "Hello, AES encryption!"

	// Encrypt the plaintext
	encrypted, err := EncryptAES([]byte(plaintext), UserSessionkey)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Println("Encrypted:", encrypted)

	// Decrypt the ciphertext
	decrypted, err := DecryptAES(encrypted, UserSessionkey)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Println("Decrypted:", decrypted)

}
