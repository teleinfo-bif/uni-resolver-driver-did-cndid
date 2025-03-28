![DIF Logo](https://raw.githubusercontent.com/decentralized-identity/universal-resolver/master/docs/logo-dif.png)

# Universal Resolver Driver: did:cndid

This is a [Universal Resolver](https://github.com/decentralized-identity/universal-resolver/) driver for Teleinfo **did:cndid** identifiers.

### Specifications

* [Decentralized Identifiers](https://www.w3.org/TR/did-core/)
* [DID Method Specification](https://github.com/teleinfo-bif/cndid/blob/main/doc/en/CNDID%20Protocol%20Specification.md)

### Example DIDs

```
did:cndid:sf24eYrmwXt6nx4fig3XJm7n9UP6PNRJ3
```

### Run with Docker

To build and run the driver as a Docker container:

```
docker build -t teleinfo/driver-did-cndid:v1.0.0 .
docker run -p 8080:8080 teleinfo/driver-did-cndid:v1.0.0
curl -X GET http://localhost:8080/1.0/identifiers/did:cndid:sf24eYrmwXt6nx4fig3XJm7n9UP6PNRJ3
```

### Pull from DockerHub

To pull the Docker image from DockerHub, run:

```
# Pull the image
docker pull teleinfo/driver-did-cndid:v1.0.0
# Run the image, as in the previous section
docker run -p 8080:8080 teleinfo/driver-did-cndid:v1.0.0
curl -X GET http://localhost:8080/1.0/identifiers/did:cndid:sf24eYrmwXt6nx4fig3XJm7n9UP6PNRJ3
```

### Driver Metadata

The driver returns the following metadata in addition to a DID document:

* `proof`: Some proof info about the DID document.
* `created`: The DID create time.
* `updated`: The DID document last update time.