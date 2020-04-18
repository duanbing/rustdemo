package jwtutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/big"
)

const (
	//ES256 keys
	ECDSAKeyD = "CCFDFDC9C2572D15C639D07E3C6C8804A1E941B13F5D10C7297A2DFAA70E6393"
	ECDSAKeyX = "EE4C3E11EB1BF081CFD4B5CCC482E069BFBECA07D566238F29191716319B809E"
	ECDSAKeyY = "A40CCD993EC355326588E2A9E202C24A2D5D1BE5128B19885FD9F2C4155C3EF1"

	//HS256 signed key
	SIGNED_KEY = "nihaoduanbing"
)

//获取签名算法为ES256的token
//该token的内容只有Redis的key,用于保存用户的登录状态
func GetEStoken(redisValue string) string {
	keyD := new(big.Int)
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyD.SetString(ECDSAKeyD, 16)
	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)

	claims := jwt.MapClaims{
		"redisValue": redisValue,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	privateKey := ecdsa.PrivateKey{D: keyD, PublicKey: publicKey}

	ss, err := token.SignedString(&privateKey)
	if err != nil {
		fmt.Println("ES256的token生成签名错误,err=%v", err)
		return ""
	}
	return ss
}

func GetPublicKey() *ecdsa.PublicKey {
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}
	return publicKey
}

func GetPrivateKey() []byte {
	keyD := new(big.Int)
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyD.SetString(ECDSAKeyD, 16)
	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	key := ecdsa.PrivateKey{D: keyD, PublicKey: publicKey}
	keyDer, err := x509.MarshalECPrivateKey(&key)
	if err != nil {
		panic(err)
	}

	keyBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyDer,
	}
	return pem.EncodeToMemory(&keyBlock)
}

//获取签名算法为HS256的token
func getHStoken(tokenFirst string, id, username string) string {
	claims := jwt.MapClaims{
		"tokenES": tokenFirst,
		//解析时，该变量的类型被转换成float64
		"uid":      id,
		"username": username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密算法是HS256时，这里的SignedString必须是[]byte（）类型
	ss, err := token.SignedString([]byte(SIGNED_KEY))
	if err != nil {
		fmt.Println("token生成签名错误,err=%v", err)
		return ""
	}
	return ss
}

//解析签名算法为ES256的token
func parseEStoken(tokenES string) string {
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	token, err := jwt.Parse(tokenES, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return &publicKey, nil
	})
	if err != nil {
		fmt.Println("ES256的token解析错误,err=%v", err)
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["redisValue"].(string)
	}

	fmt.Println("ParseEStoken:Claims类型转换失败")
	return ""
}

//解析签名算法为HS256的token
func parseHStoken(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNED_KEY), nil
	})
	if err != nil {
		fmt.Println("HS256的token解析错误，err:", err)
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("ParseHStoken:claims类型转换失败")
		return nil
	}
	return claims
}
