# Setting up etcd in a secure way

Use the pre-generated certs and keys in this directory or create your own following below steps.

Pre-generated certs and keys are:

- [ca.pem](ca.pem) … Certificate Authority (CA) file
- [server.pem](server.pem) … server certificate
- [server-key.pem](server-key.pem) … server private key
- [client.p12](client.p12) … client PKCS12 certificate+key with password `reshifter`
- [client-key.pem](client-key.pem) … client PKCS12 private key


## etcd2

To launch a secure etcd2:

```
$ cd certs/
$ docker run -d  -v $(pwd)/:/etc/ssl/certs -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 \
--ca-file /etc/ssl/certs/ca.pem --cert-file /etc/ssl/certs/server.pem --key-file /etc/ssl/certs/server-key.pem \
--advertise-client-urls https://0.0.0.0:2379 --listen-client-urls https://0.0.0.0:2379
```

To query the secure etcd2:

```
$ cd certs/

# using curl, with SSL/TLS verification, using the CA file:
$ curl --cacert $(pwd)/ca.pem --cert $(pwd)/client.p12 --pass reshifter -L https://127.0.0.1:2379/version

# using curl, without SSL/TLS verification:
$ curl --insecure --cert $(pwd)/client.p12 --pass reshifter -L https://127.0.0.1:2379/version

# using http, without SSL/TLS verification (since http doesn't do verification, neither password checks the client key):
$ http --verify=no --cert=./client.pem --cert-key=./client-key.pem  https://127.0.0.1:2379/version
```

## etcd3

To launch a secure etcd3:

```
$ cd certs/
$ docker run -d  -v $(pwd)/:/etc/ssl/certs -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v3.1.0 /usr/local/bin/etcd \
--ca-file /etc/ssl/certs/ca.pem --cert-file /etc/ssl/certs/server.pem --key-file /etc/ssl/certs/server-key.pem \
--advertise-client-urls https://0.0.0.0:2379 --listen-client-urls https://0.0.0.0:2379
```

To query the secure etcd3, use the same commands as above for etcd2.

## Create your own certificates

You will need [openssl](https://www.openssl.org/source/) and [cfssl](https://pkg.cfssl.org/) for the following steps.

### Create CA

Use following options for the CA:

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

### Convert to PKCS12

```
$ cd certs/
$ openssl pkcs12 -export -in ./client.pem -inkey ./client-key.pem -out client.p12 -password pass:reshifter
```
