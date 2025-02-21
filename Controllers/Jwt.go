package controller

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

const (
	privateKeyPem = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC31WQlvnveBvgJ
g2XrxpQMtBeJok5YbV7A8Ilvgenmt3JIH/dJNTINoDDTlqAn/A8UfhUWftYpZup1
Gjr/8+ikH8Xao+rEKCZ3f94pln8uQJhGElUs0Ge2kUTSHRZ4/1L3uOd+oEPRuBtW
gsHkijolj33Hry1A6op4HwsrMQGpNEh54sue3KdxlI93yaGEGV1eGIBovm36A03k
Ol4zCOht3UTmWpmgD1pyi8gB6icEVGIHPwCd344AnfkJIQ7W7ap4kWf3hvbb2Vwc
yt9APVALzgz5dG+u27bu1/plgiP156ZcRllDxg0eFako9d70VMvKIg29Iafy0GBb
2TSjroA5AgMBAAECggEAGyIS4wNHcxDiQT00qOUpauqV4smi+KhD6QRXtK6fIF+J
LZ4SOKryVVKEgmZkAyLP8v1dDXHxGDFJf7k8ZhTRDJBn+ophF0y5yL+Fweuln+UG
1KjWC4RDGo48cyq562f8DfYrrOPovqaG2nD4P0wroumX7gYsDr3PbEVgt6JHFXsS
1YvRib/9/J4Y4MXSZfRU5h8DW/U9vVuGStClDTbKrghr+MqOFT6sN8cCcxlxSCS3
AHVKk32aVY4KfSkjLxJqiIHquSkb9j14uyUrGvTxblZiUHLjFMg8XDw+fHobfULX
t8FebDhJSSlwT8CNpq7klYeF6GHAb3epdbFS+BrclQKBgQDpafTmo7siIsgEACge
bw3d9JU+6LCIh5ixCyWj24SZFsQhfnM0QOb44vAqonsOe/N6ISqibL+YrrXBASnT
vhCKijMJQUJb7Xg5Eo+mu2bCMW7B2xQShldGI2l4oVZsPUKMNwtJJpxpZNrSGaSh
2QL+0A1aktOvflFSd836J9bR1wKBgQDJnz+rnYkB6tY+lrWRfD2giO3zo8RHCRqC
2ZXh9t6IiWE64AQ4W/0EvB3fHClm/Lh54J6zp42HndpZjrYJkyUD6AFI9JJvz0wt
UZWzDOwIpSe5q0yo264KXGtxvFJAkG9ng7XAfYfsJcWV8yllxx5fDghvT3iVrilw
yF3jj1ocbwKBgGlSilNYJiStFRvZBkFVUyiIKKAOVzoEFX4tzXo2n4qEn1ONv2Yg
sxgzLrPORUCv5ZmCRb6s23eFvjWs1Lba2JPq8ESI0eyxJsJ6AZ/2h9Owgo2u0Uva
mp7nc7we8OQ+cDzcyZbkeUeFXsfXElaFmbhpIN1xy1sw4HkZ3jO2rlRZAoGBAJUv
DAGDsxiEFrqA8SAQ+diK+OZJyrV0+vTO4qQr8kS8wgC6OOUqy3BxcPjg8ZGOdUFY
/pSX6ZTrK+EQQ4maIs1dIXZF4QRyMj9mGoo9iXhsG5S6NyLKSWDJOYcSfRngxU2m
mxkuyR/mYuis33i7eLGExKD4AJVgJLGa0D3MmDRpAoGBAKCBDb4XdPdBAcz5Bnf1
7pj9GtPy6tmmNW/R/XyOMwQe5iyzNdtaPmr/tYtrWLBcjPKS2LCH1Qs6dok2JDUG
iQxNfkaekiJ+PFbYrOFM0Mj+z4/hXz5MkM8WZeaft2ILzl4W8VAnpxzWY6zwSXuU
StDlwga4pZiicEeFA5Nizqrt
-----END PRIVATE KEY-----
`
)

func GenerateToken(claims jwt.Claims) (string,error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))
	if err != nil {
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "",err
	}
	return tokenString,nil
}

