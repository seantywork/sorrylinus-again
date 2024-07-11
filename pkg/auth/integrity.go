package auth

import "unicode"

func SanitizePlainNameValue(raw string) string {

	var result string = ""

	for _, c := range raw {

		if unicode.IsLetter(c) {

			result = result + string(c)

		} else if unicode.IsDigit(c) {

			result = result + string(c)

		} else {

			result = result + "-"

		}

	}

	return result

}

func VerifyCodeNameValue(raw string) bool {

	for _, c := range raw {

		if unicode.IsLetter(c) {

			continue

		} else if unicode.IsDigit(c) {

			continue

		} else {

			return false
		}

	}

	return true
}
