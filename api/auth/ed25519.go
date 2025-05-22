package auth

func verifySignature(pubKeyStr, message, signatureBase64 string) bool {
	// pubKey, err := base64.StdEncoding.DecodeString(pubKeyStr)
	// if err != nil {
	// 	return false
	// }

	// signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	// if err != nil {
	// 	return false
	// }

	// Ed25519 서명 검증
	// if isValid := ed25519.Verify(pubKey, []byte(message), signature); !isValid {
	// 	return false
	// }

	return true
}
