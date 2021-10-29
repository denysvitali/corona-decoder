package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	"github.com/stapelberg/coronaqr"
	"github.com/stapelberg/coronaqr/trustlist/trustlistmirror"
)

var args struct {
	InputFile string `arg:"-f,--input,required"`
	Verify    *bool  `arg:"-v"`
}

func main() {
	arg.MustParse(&args)
	log := logrus.New()
	f, err := os.Open(args.InputFile)
	if err != nil {
		log.Fatalf("unable to open certificate file: %v", err)
	}
	defer f.Close()
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	cert, err := coronaqr.Decode(strings.TrimSpace(string(fileBytes)))
	if err != nil {
		log.Fatalf("unable to decode certificate: %v. Make sure you provided a valid certificate in textual form.", err)
	}

	var decoded *coronaqr.Decoded

	if args.Verify != nil && *args.Verify {
		provider, err := trustlistmirror.NewCertificateProvider(context.Background(), trustlistmirror.TrustlistDE)
		if err != nil {
			log.Fatalf("unable to initialize trustlistmirror (certificate provider): %v", err)
		}
		decoded, err = cert.Verify(provider)
		if err != nil {
			log.Fatalf("unable to verify certificate: %v", err)
		}
	} else {
		decoded = cert.SkipVerification()
	}

	const dateFormat = "2006-01-02 15:04:05 MST"

	for idx, v := range decoded.Cert.VaccineRecords {
		fmt.Printf("VR %d: C=%s,ID=%s,ISS=%s\n", idx, v.Country, v.CertificateID, v.Issuer)
	}
	kid := decoded.Kid
	if len(kid) == 0 && decoded.SignedBy != nil {
		hash := sha256.Sum256(decoded.SignedBy.Raw)
		kid = hash[:8]
	}
	fmt.Printf("KID: %s\n", base64.StdEncoding.EncodeToString(kid))
	fmt.Printf("Issued At: %+v\n", decoded.IssuedAt.Format(dateFormat))
	if decoded.SignedBy != nil {
		fmt.Printf("Signed By: %s (issued by: %s)\n", decoded.SignedBy.Subject, decoded.SignedBy.Issuer)
	}
	fmt.Printf("Expiration: %+v\n", decoded.Expiration.Format(dateFormat))
	fmt.Printf(
		"Personal Name: %s %s\n",
		decoded.Cert.PersonalName.GivenName,
		decoded.Cert.PersonalName.FamilyName)
	fmt.Printf("DOB: %+v\n", decoded.Cert.DateOfBirth)
}
