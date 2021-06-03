package tools

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var alog map[string]string = map[string]string{
	"MD2-RSA":       "MD2WithRSA",
	"MD5-RSA":       "MD5WithRSA",
	"SHA1-RSA":      "SHA1WithRSA",
	"SHA256-RSA":    "SHA256WithRSA",
	"SHA384-RSA":    "SHA384WithRSA",
	"SHA512-RSA":    "SHA512WithRSA",
	"SHA256-RSAPSS": "SHA256WithRSAPSS",
	"SHA384-RSAPSS": "SHA384WithRSAPSS",
	"SHA512-RSAPSS": "SHA512WithRSAPSS",
}


func Sha256WithRsa(source string, privateKey *rsa.PrivateKey) (signature string, err error) {
	if privateKey == nil {
		return "", fmt.Errorf("private key should not be nil")
	}
	h := crypto.Hash.New(crypto.SHA256)
	_, err = h.Write([]byte(source))
	if err != nil {
		return "", nil
	}
	hashed := h.Sum(nil)
	signatureByte, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signatureByte), nil
}


func Signer(formdata map[string]string ,privateKey *rsa.PrivateKey) (string,error) {
	Newformdata:=formatBizQueryParaMap(formdata)
	fmt.Println(formdata)
	sign:=""
	for _,v:=range Newformdata{
		sign+=v+"&"
	}
	sign=strings.TrimRight(sign,"&")

	fmt.Println("代签名字符串：",sign)
	return Sha256WithRsa(sign,privateKey)
}
func Signdata(formdata map[string]string)string{
	Newformdata:=formatBizQueryParaMap(formdata)

	sign:=""
	for _,v:=range Newformdata{
		sign+=v+"&"
	}
	sign=strings.TrimRight(sign,"&")

	return sign
}
func formatBizQueryParaMap(formdata map[string]string) []string{

	var strs []string
	for k := range formdata {
		strs = append(strs, k)
	}
	sort.Strings(strs)
	NewMap:=make([]string,0)
	for _, k := range strs {
		NewMap=append(NewMap,k+"="+formdata[k])
	}
	return NewMap

}

// LoadPrivateKeyWithPath 通过私钥的文件路径内容加载私钥
func LoadPrivateKeyWithPath(path string) (privateKey *rsa.PrivateKey, err error) {
	privateKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private pem file err:%s", err.Error())
	}
	return LoadPrivateKey(string(privateKeyBytes))
}
// LoadPrivateKey 通过私钥的文本内容加载私钥
func LoadPrivateKey(privateKeyStr string) (privateKey *rsa.PrivateKey, err error) {
	//fmt.
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("decode private key err")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key err:%s", err.Error())
	}
	//privateKey:= key
	/*if !ok {
		return nil, fmt.Errorf("%s is not rsa private key", privateKeyStr)
	}*/
	return key, nil
}

func LoadPublicKey(privateKeyStr string) (privateKey *rsa.PublicKey, err error) {
	//fmt.
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("decode private key err")
	}
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key err:%s", err.Error())
	}
	//privateKey:= key
	/*if !ok {
		return nil, fmt.Errorf("%s is not rsa private key", privateKeyStr)
	}*/
	return key, nil
}

func GetCertRootSn(certPath string) (string, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", err
	}
	strs := strings.Split(string(certData), "-----END CERTIFICATE-----")

	var cert bytes.Buffer
	for i := 0; i < len(strs); i++ {
		if strs[i] == "" {
			continue
		}
		if blo, _ := pem.Decode([]byte(strs[i] + "-----END CERTIFICATE-----")); blo != nil {
			c, err := x509.ParseCertificate(blo.Bytes)
			if err != nil {
				continue
			}
			if _, ok := alog[c.SignatureAlgorithm.String()]; !ok {
				continue
			}
			si := c.Issuer.String() + c.SerialNumber.String()
			if cert.String() == "" {
				cert.WriteString(Md5(si))
			} else {
				cert.WriteString("_")
				cert.WriteString(Md5(si))
			}
		}

	}
	return cert.String(), nil
}