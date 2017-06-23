# Setting up etcd in a secure way

## etcd2

### Create Certificate Authority

Use following options for the Certificate Authority (CA):

```
$ cat ca-config.json
{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "server": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth"
                ]
            },
            "client": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            }
        }
    }
}
```

The Certificate Signing Request (CSR) should look like:

```
$ cat ca-csr.json
{
    "CN": "The ReShifter CA",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "CA",
            "O": "ReShifter",
            "ST": "San Francisco",
            "OU": "Pier 39"
        }
    ]
}
```

Create CA:

```
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
```

### Server certs

Generate certs for the server:

```
$ cat server.json
{
    "CN": "*",
    "hosts": [
        "127.0.0.1",
        "localhost"
    ],
    "key": {
        "algo": "ecdsa",
        "size": 256
    },
    "names": [
        {
            "C": "US",
            "L": "CA",
            "ST": "San Francisco"
        }
    ]
}
$ cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server.json | cfssljson -bare server
```

### Client certs

Generate certs for the client:

```
$ cat client.json
{
    "CN": "client",
    "hosts": [""],
    "key": {
        "algo": "ecdsa",
        "size": 256
    },
    "names": [
        {
            "C": "US",
            "L": "CA",
            "ST": "San Francisco"
        }
    ]
}
$ cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client.json | cfssljson -bare client
```

### Verify

```
$ openssl x509 -in ca.pem -text -noout
$ openssl x509 -in server.pem -text -noout
$ openssl x509 -in client.pem -text -noout
```
