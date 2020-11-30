// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package certs is set of default certificates serialized in to string format.
*/
package certs

// Default Client Certificate name
const (
	Client1Key = "client1.key"
	Client1Crt = "client1.crt"
)

/*
All of these are copied over from https://github.com/onosproject/gnxi-simulators/tree/master/pkg/certs
*/

/*
OnfCaCrt is the default CA certificate
Certificate:
    Data:
        Version: 1 (0x0)
        Serial Number:
            de:f7:d7:d2:37:da:b1:49
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = ca.opennetworking.org
        Validity
            Not Before: Apr 11 09:06:13 2019 GMT
            Not After : Apr  8 09:06:13 2029 GMT
        Subject: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = ca.opennetworking.org
*/
const OnfCaCrt = `
-----BEGIN CERTIFICATE-----
MIIDYDCCAkgCCQDe99fSN9qxSTANBgkqhkiG9w0BAQsFADByMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCU1lbmxvUGFyazEMMAoGA1UECgwDT05G
MRQwEgYDVQQLDAtFbmdpbmVlcmluZzEeMBwGA1UEAwwVY2Eub3Blbm5ldHdvcmtp
bmcub3JnMB4XDTE5MDQxMTA5MDYxM1oXDTI5MDQwODA5MDYxM1owcjELMAkGA1UE
BhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlNZW5sb1BhcmsxDDAKBgNVBAoM
A09ORjEUMBIGA1UECwwLRW5naW5lZXJpbmcxHjAcBgNVBAMMFWNhLm9wZW5uZXR3
b3JraW5nLm9yZzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMEg7CZR
X8Y+syKHaQCh6mNIL1D065trwX8RnuKM2kBwSu034zefQAPloWugSoJgJnf5fe0j
nUD8gN3Sm8XRhCkvf67pzfabgw4n8eJmHScyL/ugyExB6Kahwzn37bt3oT3gSqhr
6PUznWJ8fvfVuCHZZkv/HPRp4eyAcGzbJ4TuB0go4s6VE0WU5OCxCSlAiK3lvpVr
3DOLdYLVoCa5q8Ctl3wXDrfTLw5/Bpfrg9fF9ED2/YKIdV8KZ2ki/gwEOQqWcKp8
0LkTlfOWsdGjp4opPuPT7njMBGXMJzJ8/J1e1aJvIsoB7n8XrfvkNiWL5U3fM4N7
UZN9jfcl7ULmm7cCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAIh6FjkQuTfXddmZY
FYpoTen/VD5iu2Xxc1TexwmKeH+YtaKp1Zk8PTgbCtMEwEiyslfeHTMtODfnpUIk
DwvtB4W0PAnreRsqh9MBzdU6YZmzGyZ92vSUB3yukkHaYzyjeKM0AwgVl9yRNEZw
Y/OM070hJXXzJh3eJpLl9dlUbMKzaoAh2bZx6y3ZJIZFs/zrpGfg4lvBAvfO/59i
mxJ9bQBSN3U2Hwp6ioOQzP0LpllfXtx9N5LanWpB0cu/HN9vAgtp3kRTBZD0M1XI
Ctit8bXV7Mz+1iGqoyUhfCYcCSjuWTgAxzir+hrdn7uO67Hv4ndCoSj4SQaGka3W
eEfVeA==
-----END CERTIFICATE-----
`

/*
DefaultClientKey is the default client key
openssl rsa -in client1.key -text -noout
Private-Key: (2048 bit)
*/
const DefaultClientKey = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDK2amlhSGecWBI
9jlj3OvCoGMAlHXffwHAsHskE/bFkyRKeYs0MHlcNBXSjgADtHt7FpGJ7nbZsMdu
0XkNOlN/X0i8FKU7X9DnOHbHhOitg469J89D8tHnN2wrTyGxY2Km95PxeCQ4NH7z
4MSSbdFfFXnw9k+qr+e1owswloJsYo+eyP7hjy7QvR2sAX8CC6D0ZaVuqLGK5kHJ
UyMxifqMa4D4wD8zfzahcTe2QAeOxQwhsaK2pZ/09bI/KCj8vxbxG2nkcwUgjl0L
gh0g1BFNSu39v1DjsiFWrIxnj3CL2Ncu0ps0+dPSdpx90ImvjL4BRPz8V4iLPljm
BmQfU2HpAgMBAAECggEAVEd9Ca03m5nldEsA6zHVrmZu28XS94nQU5u/fezhgZMx
59N597QQKDPnwTSIYwGwsCJfU5yFOssNAUj873cFTA1trd8yC2oy5G58Q0dAWR8o
xgRtRAD2HwfS5GebSxVM3qxMhm3xNnzxJiiD44bHD6dfo7LixLsTHU9hjc1q4NaR
BfJCl01f1EgEtPB4L7l3E1ivgy2OQ91gSJsIneA4RFLn5NbMac7ZWOTMGThjGFAw
qR38Ua7TzQXL+qNzHU0JUOBg4smk5tzPxD8ik2p8VTc8WrPEBX129fn+JumhKoIC
4LysZoazLhKem56LBVSW/TNybXPy9z+nBhbFbnoLAQKBgQDpPfvc83f/f5Iu3PaG
8BsIoOonasqGsFoDB7XN/DZdjk8d9MKF+2tK8tRtwfZ+cWJL7m94ZCudtU72uhwN
0x1E/p+qAUBvE3WR3MxQQb6R5kyzMUlwTjwrM49kMuQdU7CaL2abqRvvekkE6L13
9KesvE0TQSTx+beMbzUUdDyCCQKBgQDepInF/pmwbM0E2zKYny0Ro3caW6HoUEoJ
ZGcZ/NLdTgiQ3GJsVddc0M4jgXMndjOJFZcDDeFPgRAgQVnJHbDDWNygmNDBKavj
LNcTIB2AO4ObWIH2C6x38V6l+62Nc9QPlI52tiekWVO1tAhy3FJXRHAUhP43bi36
VnCg0mJY4QKBgFuSTEnpBJm4+imP8vHzXom6s3OaR70ti4lZA5XFiYqdjo5SQ/Ta
Srt4LtKQrjfiSBdLm1QG7+DRCBlx5AXBduJZnVHff+6cEzKbH1P7G9ioNEC9/vkq
nhDQA2HxYQHqk5FVPtGqSR9yQSy+O3TXBuWYYCJJFzoxMlDecFaBdCgRAoGBAIRc
qZPOQyyB4nkKn8/ggfjEh+BhraXhZcKjsC/hALOU2r7UZqcleX2ynXq6UO2a9hR/
g2HLdLHBdwbWEzzfq+DXCYNolmLgFVJfrBWwuBkuSJWoTssqMYS1OKHROGKqA96n
YPLuZC7u9DdIKuWuWj2LcF6imkf19tunXBogOVvBAoGBAJvIYUGba3WU+/Ugs3yZ
ubNBC52mS2kLrOTZfzQOP6sL3OC5rFgm7A4e/CMCWphns98ftB0c1RaZBCcwskUE
mCxUxvivP7oUNf0CMZd6CW2M0LEWMVrIFJ1o0BtNpUVZ/+nipJblRA1B/JLsgNvn
11mpIWR4KHd8A5wjkkfb2d/e
-----END PRIVATE KEY-----
`

/*
DefaultClientCrt is the default client certificate
    Data:
        Version: 1 (0x0)
        Serial Number:
            4a:c1:1b:3b:17:1e:8d:65:f1:b9:99:46:60:e4:17:e8:76:6e:c7:53
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = ca.opennetworking.org
        Validity
            Not Before: May  5 07:42:13 2020 GMT
            Not After : May  3 07:42:13 2030 GMT
        Subject: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = client1.opennetworking.org
*/
const DefaultClientCrt = `
-----BEGIN CERTIFICATE-----
MIIDcDCCAlgCFErBGzsXHo1l8bmZRmDkF+h2bsdTMA0GCSqGSIb3DQEBCwUAMHIx
CzAJBgNVBAYTAlVTMQswCQYDVQQIDAJDQTESMBAGA1UEBwwJTWVubG9QYXJrMQww
CgYDVQQKDANPTkYxFDASBgNVBAsMC0VuZ2luZWVyaW5nMR4wHAYDVQQDDBVjYS5v
cGVubmV0d29ya2luZy5vcmcwHhcNMjAwNTA1MDc0MjEzWhcNMzAwNTAzMDc0MjEz
WjB3MQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCU1lbmxvUGFy
azEMMAoGA1UECgwDT05GMRQwEgYDVQQLDAtFbmdpbmVlcmluZzEjMCEGA1UEAwwa
Y2xpZW50MS5vcGVubmV0d29ya2luZy5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDK2amlhSGecWBI9jlj3OvCoGMAlHXffwHAsHskE/bFkyRKeYs0
MHlcNBXSjgADtHt7FpGJ7nbZsMdu0XkNOlN/X0i8FKU7X9DnOHbHhOitg469J89D
8tHnN2wrTyGxY2Km95PxeCQ4NH7z4MSSbdFfFXnw9k+qr+e1owswloJsYo+eyP7h
jy7QvR2sAX8CC6D0ZaVuqLGK5kHJUyMxifqMa4D4wD8zfzahcTe2QAeOxQwhsaK2
pZ/09bI/KCj8vxbxG2nkcwUgjl0Lgh0g1BFNSu39v1DjsiFWrIxnj3CL2Ncu0ps0
+dPSdpx90ImvjL4BRPz8V4iLPljmBmQfU2HpAgMBAAEwDQYJKoZIhvcNAQELBQAD
ggEBADERPV5ykcjUaskgPHk081A0LmsKKL9AdqLfREreRIJgPg1T+ppE7PAf+i9m
W6lojz+qCOkjd9d7Cd7wWRvrndRs6xt3HoIu3Z+knds/SC6emz/Rf+ZYsODNC3b3
zv7agEeuD2M/Lq8pLfh9MsCxFmbz4JuiJ3kkT8zhiRAbUW6LkYFzdSeq6uUrhZPv
C6VCZtHRQiEHE84Wwkxcz7shyknPwwDZxLoNeGjQPruESsY9Dyt96dz8yBQ6ukD1
hPpnMDb+sksoX0ZOoEYpNdqDvZWa25n64LBBLuxbCuelSVUpsEdjRvh10znFY36x
op/V1+9Wot9zMohw6bu5Aj9K6s8=
-----END CERTIFICATE-----
`

/*
DefaultLocalhostKey is the default localhost server key
openssl rsa -in client1.key -text -noout
Private-Key: (2048 bit)
*/
const DefaultLocalhostKey = `
-----BEGIN PRIVATE KEY-----
MIIEugIBADANBgkqhkiG9w0BAQEFAASCBKQwggSgAgEAAoIBAQC2F3yVkmiXIdeC
NMYNzltYaRdUPptIRz9nKCFaKHWu+tVB1Pz0nh1xShQ3ObIkhQ8g4v+bms8tEbcZ
XiVeD5hOJJEYse04HZ1EoqpXQB+xl8gFN5UMu0mEGC67tEkXi8narqSEebqUJtPK
MGBM4F9QEre1dOgOiI5+LOZmGgd58DKzF0Tel1SVv7k0ysZ4dzEnsitaulWPBKz7
UzQ4OFyY/boN72tJ9Q3BiqDnUQH/ac+kqg8G1kXUn5EsAH3d6Nv/ku3lW0Gr9uxI
y/jiwvZVTG7l792jPNzRI3LmnffAzg/nGcqTQXZgAvxPz7eAc7c2v4xXrcp9/tWU
NEa9alN5AgMBAAECggEAcrj1ax7k+mL97jDlnwkmD9uWMSOInc8VqR5ldPIMwwOR
nHpeLJf5oMi1V93n2I5ka6nYtOaiJJkGrNrd3BcjNAhhyhc/h51Q2k9J1tK1pSQl
hvPv2iedN7Ysq2H4svcFY9uoFzbCUFjuEnLMGWM7aa2BRLe1BIMQk3oiZq17jFyy
v/IYDI8jNFLGyQFwP3iO5EwCIza1MOeKz1vXKPqnLYUk6VjDaOYH3hPqeaMdfzAp
1NzEZtzfPkJFYu11qmoNwZmTEthhIyJ1pKn82CuE8CiL8jdSQq144k1YR19NmEZX
HU5hZvijU4Iz+1NkAfka6MYuYXEvUKkqkqQzfwnjiQKBgQDfKOAxK1/xalbn9wjo
SkywpmRjbS6M2q7D4fTk7O7PGLTSjroFRcNUe7HI8+jH3GK7Jm1ddViUNaWqQb+8
9Q7SAAro+CakAuPAbf0e4VXzzyHE/lvtruc/SL953wZGIoAZ38TXsi4fm/UTD9Ox
hVqu/N5AvY4P1g6IXfqynPdXjwKBgQDQ43QStDnOjsiwImOOsKjefhYtHhZQvGQf
cXE3d+qarmxtifYkEyPLLYtQkb8biBDVJaWnEVcqg9Kuxbq6yOLROGPRtio2WIns
FiHiu/h6sgdsBfzKIlcn23UU0Fce2nuMLwC8PGoWkdMQPR+6NTKe1qKkkzhDpA3x
TbGJSzVgdwKBgG5ReLMV7DIeDaRSnRaoVE0nlI0KVm7PVIIFW9knv86lOg60/ATL
Pgqvs23SFgtnSW+XSY1gC1AJTUJjinPQ+WibGMmekwuVWh2wwebYInOKu/j0fWF8
i1jfj7ihpipZt9YSpu6yaNa7dGXd9xrU/8VtwDlk+6uceEa1ns9ZhXTFAn8SxFyp
UYfgBvQA3xYSu8xwMOPNKebXWhWkvYxub1ekjgcv0DVNCGsu1eiuVGnXD2Jzw+4e
FHDAYReMnDcqkOHP6kENllA0kb/SdiqVNE4et9/y1JbhkjRCYHUkaZNqMjbnYVGv
l73wSSmtS9CN6jmiC6aRIqjratHV3CUXMKqbAoGAHvRu9AjfHAnURpUk93DBvMoe
soydv29w2p6V5aiEj0+qRcoDZc1iO4QVy463twkKyEduW+pSJPwYI3/KZCf1zPvZ
apAnh710u1lw5DBjSDN9H+ZL2m6myYy2vdsTUh0T8TN3J5i7kUpobCj7fGvng8dt
f/cOAHPSczsq+mGIYxo=
-----END PRIVATE KEY-----
`

/*
DefaultLocalhostCrt is the default localhost server certificate
Certificate:
    Data:
        Version: 1 (0x0)
        Serial Number:
            4a:c1:1b:3b:17:1e:8d:65:f1:b9:99:46:60:e4:17:e8:76:6e:c7:54
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = ca.opennetworking.org
        Validity
            Not Before: May  5 07:43:55 2020 GMT
            Not After : May  3 07:43:55 2030 GMT
        Subject: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = localhost
*/
const DefaultLocalhostCrt = `
-----BEGIN CERTIFICATE-----
MIIDXzCCAkcCFErBGzsXHo1l8bmZRmDkF+h2bsdUMA0GCSqGSIb3DQEBCwUAMHIx
CzAJBgNVBAYTAlVTMQswCQYDVQQIDAJDQTESMBAGA1UEBwwJTWVubG9QYXJrMQww
CgYDVQQKDANPTkYxFDASBgNVBAsMC0VuZ2luZWVyaW5nMR4wHAYDVQQDDBVjYS5v
cGVubmV0d29ya2luZy5vcmcwHhcNMjAwNTA1MDc0MzU1WhcNMzAwNTAzMDc0MzU1
WjBmMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCU1lbmxvUGFy
azEMMAoGA1UECgwDT05GMRQwEgYDVQQLDAtFbmdpbmVlcmluZzESMBAGA1UEAwwJ
bG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAthd8lZJo
lyHXgjTGDc5bWGkXVD6bSEc/ZyghWih1rvrVQdT89J4dcUoUNzmyJIUPIOL/m5rP
LRG3GV4lXg+YTiSRGLHtOB2dRKKqV0AfsZfIBTeVDLtJhBguu7RJF4vJ2q6khHm6
lCbTyjBgTOBfUBK3tXToDoiOfizmZhoHefAysxdE3pdUlb+5NMrGeHcxJ7IrWrpV
jwSs+1M0ODhcmP26De9rSfUNwYqg51EB/2nPpKoPBtZF1J+RLAB93ejb/5Lt5VtB
q/bsSMv44sL2VUxu5e/dozzc0SNy5p33wM4P5xnKk0F2YAL8T8+3gHO3Nr+MV63K
ff7VlDRGvWpTeQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBkaee8haWLCP9cUIsm
TWAZxLFV4jqbEdu0/lJDpe8vPlT8rh+Kh47VS0snTR7aL2SwPz+vrU13cJdNtkCQ
NvgpTUZN6GSPrPCOPgJWDvjQ4eeNd+KmdD/GAxk/cZUbC1gJmXTTPu/ZUN03IRKF
98Tg/oVxq74fKPlJHHAa/eBZBvrjIsFyTXaYMnFMY1vRoWX2bxfFC2vHhJt6xAyV
6s4ymYhfr2zrWwDpLxsAjyBQ7dCi/UEYVYeTANpk0U4CJqBRb9DQuWEm4oc3eAfL
ac3qA/O+DFJMLDAa/ChFI1c08mlB/+gbfYPRPB3xj0zvI02KGkDigqHrVcaPXQzU
0C6a
-----END CERTIFICATE-----
`

/*
DefaultOnosConfigKey is the default onos-config server key
openssl rsa -in onso-config.key -text -noout
Private-Key: (2048 bit)
*/
const DefaultOnosConfigKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCZn2yJTuwzOI08
3/2/OaIwzQzvzuGBnp4ttZGxrSi0KdX7r19Vu2dI5eEDbSNNYa4+9xSwT0lTf9u0
/9i6mJYgejLNJU0CnAZDGxTLuxMHeGjPjeSpdVGuXq4qOj+kXKmdpATXZmyXwZlK
gamh7/u66Wcrgn2y6eK1iCIvK+L0U+vwAa6NokUJwZlFtLBPmE0wwZxnW9LMyWkz
K0KANCsua3uVBSFclREsnM+jw4V2Gk9smXxfRV1HWkP7xsZuGi8NHq8A1/JC1gE5
CNLVsnD/ePHV3vHYFzVTu66Q9hDJLGpyu2cTpjUv44p7/3fQsaz3fYuT31nAj5//
l9ZHJwXRAgMBAAECggEAVEwRCL+QCQNNLUxUNyxu/YxnPugtAi2B6t8pVXAJV+Nl
EjjHfYnaQTwzXufyaTHipZZ7ecvoFrOgYg/KY4n7R1MGsV94hKgNH6Gqpai/5meC
S/I2uW4xJhe6Rl20MoLOaDxqk7AWgqevcBz6cmv3nDcbb9qpExYYWziaWXwhi6Pv
unvmZGzpAu9ehb4BWlQmHVqIDAcLijTKD7uxHl9yiM7zG5KLzpPCmr4drWjyrsE+
oo7fA0GhLpjQCJ6Ye+2hBYNvMJLmvcgex++16vQDszhIy9rTtlFCqbPaADF/9bOn
rjTMPPYWijMFx0lltHhPq09527iL/Yt6szg/PbNOyQKBgQDIZjdNfNXs128Qm7Zq
opo6rQU3j17/QINUBOx/+xCWnNohBbo+pGHP6Z8Ld4Lx6P1LPoWUw5tX4v8fwXgJ
Rn3gnwTAiXx/S/dHEkiDQLMU7nlpkwRNUA00DDqAz8dZIgRZeHAE1PRxBihCjXVf
Z4PKTHsbomxOrVn+NuAV80qr2wKBgQDEPs47MTUI96Xh8BuixVAwi2QVo01QB5Dl
x/ds1v/IDHen9sV1Yjl7gLC9YpO/7kweFDZRt6Zm1IRh4Qy9yqtDUSAQEOtoI6A1
FvHCBC06NzYISBFNdmRnkS3AJEb9piJY62YeP9+vWJqMuVc39igIbYoGR654uszJ
gEH2lgu6wwKBgQDHfCjU882oBBRFPhvqLo7ElfM5iXiRMtEIVBZwl6W9p8njUWZC
cTQE2ZQ+v+sTkFCEFGq42bbLV+WK4PXylb88WE9Msg/CUAaJMwQH0+Hwlis6EuUX
aPabtwiNrUfNzHTz81XfGXVzBSQSi+oo3Exulo99xMN31kxdKJcMgrD0PQKBgEx+
1tDH645lSin5+CvIket6SjcNArPxXw/SlKW+YNHP2kyEqo+JDDMSBNKtvD4SW2VW
J55O4fQvXrLwkJDikUOaOc9JaRmc2XQYT4B7NE3++3ba8LOrNJQSSS0edvWkbrsO
dy3PZBfrh8LW9CKCNzShzi2If3/cALuC3TOLZWMVAoGAOHexq8+7EoO9FXcqyRfN
tgv108T3UD7LNcH2Fj5sx+o/NN/OcGcCCb40I5Jrh305a+GPRG1RKW7fmD8UaMaM
L9Q8NQ7nUfUvl0RR9vwAgcvLy1zD+IgqAzqF0MN/ejIHTACIgL4S0cXlUyCk5gQn
wXqkSZWW53eg8HKb/cVMLoM=
-----END PRIVATE KEY-----
`

/*
DefaultOnosConfigCrt is the default onos-config server certificate
Certificate:
    Data:
        Version: 1 (0x0)
        Serial Number:
            56:56:04:5e:9b:45:15:87:d7:24:3d:2a:22:21:df:87:11:e0:f2:0b
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = ca.opennetworking.org
        Validity
            Not Before: Jul 16 18:38:15 2019 GMT
            Not After : Jul 13 18:38:15 2029 GMT
        Subject: C = US, ST = CA, L = MenloPark, O = ONF, OU = Engineering, CN = onos-config
*/
const DefaultOnosConfigCrt = `
-----BEGIN CERTIFICATE-----
MIIDYTCCAkkCFFZWBF6bRRWH1yQ9KiIh34cR4PILMA0GCSqGSIb3DQEBCwUAMHIx
CzAJBgNVBAYTAlVTMQswCQYDVQQIDAJDQTESMBAGA1UEBwwJTWVubG9QYXJrMQww
CgYDVQQKDANPTkYxFDASBgNVBAsMC0VuZ2luZWVyaW5nMR4wHAYDVQQDDBVjYS5v
cGVubmV0d29ya2luZy5vcmcwHhcNMTkwNzE2MTgzODE1WhcNMjkwNzEzMTgzODE1
WjBoMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCU1lbmxvUGFy
azEMMAoGA1UECgwDT05GMRQwEgYDVQQLDAtFbmdpbmVlcmluZzEUMBIGA1UEAwwL
b25vcy1jb25maWcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCZn2yJ
TuwzOI083/2/OaIwzQzvzuGBnp4ttZGxrSi0KdX7r19Vu2dI5eEDbSNNYa4+9xSw
T0lTf9u0/9i6mJYgejLNJU0CnAZDGxTLuxMHeGjPjeSpdVGuXq4qOj+kXKmdpATX
ZmyXwZlKgamh7/u66Wcrgn2y6eK1iCIvK+L0U+vwAa6NokUJwZlFtLBPmE0wwZxn
W9LMyWkzK0KANCsua3uVBSFclREsnM+jw4V2Gk9smXxfRV1HWkP7xsZuGi8NHq8A
1/JC1gE5CNLVsnD/ePHV3vHYFzVTu66Q9hDJLGpyu2cTpjUv44p7/3fQsaz3fYuT
31nAj5//l9ZHJwXRAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAGZfTp6qKTGZGrEl
bZWWKe8ilMqgZtcz7J93LOYk+l8nTEg5hIQ015mHY1+R+0F2gniciayjpVG2BChD
MOHfes0StdKY0nVHy83TpG0TsY76e//DSmekZwtm+OoxualpEOLW0PgKFEE8+PdJ
b/QlN8AyWJ3cvA7hDGlCrCNontLJS+W0VAPLDrFi/NeK0RpiQ6rI2U4B2jdGIrhw
AJD2FhJHDcdfBeR80KHiSFlhgSSMChKBrlzYw2vdeuSuAuuzTn88CzXKaIki56xQ
xm422D8l3cAo4W+GP6HGtwMl/UcI+WBMpa6yGtkaCaUA9v5EwKk+oceDlSnNnHS9
tWg4SMI=
-----END CERTIFICATE-----
`
