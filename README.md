# corona-decoder

This is a super simple CLI application that uses [@stapelberg](https://github.com/stapelberg)'s coronaqr library / CLI
to provide quickly some information about a COVID-19 certificate provided in textual form.

## Usage

```
corona-decoder -f ./samples/cert-1.txt -v
```

Output:
```
VR 0: C=FR,ID=URN:UVCI:01:FR:W7V2BE46QSBJ#L,ISS=CNAM
KID: 53FOjX/4aJs=
Issued At: 2021-10-27 13:20:48 CEST
Signed By: CN=DSC_FR_023,OU=180035024,O=CNAM,C=FR (issued by: CN=CSCA-FRANCE,O=Gouv,C=FR)
Expiration: 2023-10-14 00:00:00 CEST
Personal Name: MICKEY MOUSE
DOB: 2001-12-31
```
