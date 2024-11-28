package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
)

const sequenceLength = 10000

func rsaBitGenerator(privateKey *rsa.PrivateKey, bits int) []byte {
	sequence := make([]byte, bits)

	for i := 0; i < bits; i++ {
		randomData := make([]byte, 16)
		_, err := rand.Read(randomData)
		if err != nil {
			fmt.Println("Ошибка генерации случайных данных:", err)
			return nil
		}

		ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, randomData)
		if err != nil {
			fmt.Println("Ошибка шифрования:", err)
			return nil
		}

		deciphered, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
		if err != nil {
			fmt.Println("Ошибка расшифровки:", err)
			return nil
		}

		if len(deciphered) > 0 && deciphered[len(deciphered)-1]%2 == 1 {
			sequence[i] = 1
		} else {
			sequence[i] = 0
		}
	}

	return sequence
}

func saveToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, b := range data {
		_, err := file.WriteString(fmt.Sprintf("%d", b))
		if err != nil {
			return err
		}
	}
	return nil
}

func testGolombFirst(sequence []byte) bool {
	count0, count1 := 0, 0
	for _, bit := range sequence {
		if bit == 0 {
			count0++
		} else {
			count1++
		}
	}
	return abs(count0-count1) <= 100
}

func testGolombSecond(sequence []byte) bool {
	runs := make(map[int]int)
	currentRun := 1
	for i := 1; i < len(sequence); i++ {
		if sequence[i] == sequence[i-1] {
			currentRun++
		} else {
			runs[currentRun]++
			currentRun = 1
		}
	}
	for i := 2; i <= len(runs); i++ {
		if runs[i-1] < runs[i] {
			return false
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Println("Ошибка генерации ключей RSA:", err)
		return
	}

	sequence := rsaBitGenerator(privateKey, sequenceLength)

	err = saveToFile("rsa_sequence.txt", sequence)
	if err != nil {
		fmt.Println("Ошибка при сохранении последовательности в файл:", err)
		return
	}
	fmt.Println("Последовательность сохранена")

	if testGolombFirst(sequence) {
		fmt.Println("Последовательность удовлетворяет первому постулату Голомба")
	} else {
		fmt.Println("Последовательность НЕ удовлетворяет первому постулату Голомба")
	}

	if testGolombSecond(sequence) {
		fmt.Println("Последовательность удовлетворяет второму постулату Голомба")
	} else {
		fmt.Println("Последовательность НЕ удовлетворяет второму постулату Голомба")
	}
}
