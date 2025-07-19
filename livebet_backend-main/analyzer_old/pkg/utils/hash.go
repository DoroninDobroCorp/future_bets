package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
)

// Генерация ключа матча
func GenerateMatchKey(home, away string) string {
	const emptyHash = "da39a3ee5e6b4b0d3255bfef95601890afd80709"

	h1 := sha1.Sum([]byte(home))
	h2 := sha1.Sum([]byte(away))

	hexH1 := hex.EncodeToString(h1[:])
	hexH2 := hex.EncodeToString(h2[:])

	if hexH1 == emptyHash || hexH2 == emptyHash {
		log.Printf("[WARNING] One of the teams has an empty hash: home hash = %s, away hash = %s", hexH1, hexH2)
	}

	return hexH1 + hexH2
}

func GenerateFullMatchKey(book1, book2, matchID1, matchID2, sport, outcome string) string {
	var str []string
	str = append(str, book1)
	str = append(str, book2)
	str = append(str, matchID1)
	str = append(str, matchID2)
	str = append(str, sport)
	str = append(str, outcome)

	sort.Strings(str)

	hash := md5.Sum([]byte(fmt.Sprintf("%s%s%s%s%s%s", str[0], str[1], str[2], str[3], str[4], str[5])))
	return hex.EncodeToString(hash[:])
}
