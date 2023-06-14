# go-api-encrypt
Go Language API Encryption Middleware, supporting Gin and Go-Zero frameworks

This project provides a demonstration of encryption concepts and methods, and the specific encryption approach can be implemented by yourself.

Go version >= 1.16

## Use Cases
```
When you need to encrypt and verify the request parameters of an api, you can use this middleware.

The main purpose is to prevent the tampering of request parameters and basic web scraping, ensuring api security.

However, frontend encryption is not completely secure; it only increases the difficulty of cracking. Therefore, this middleware only provides an approach, and the specific encryption method can be implemented by yourself.

The encryption method used is AES-CBC-Hex, which can be modified as needed.
```

## Usage Limitations
```
Currently, only GET and POST requests are supported. Other request methods will be directly passed through.

For GET requests, the parameters in the URL are verified, while for POST requests, the JSON parameters in the body are verified.

Complex JSON parameters have not been tested.
```

## Framework Support
[Gin](https://github.com/gin-gonic/gin)
The Gin framework provides a complete example that can be run directly.

[Go-Zero](https://github.com/zeromicro/go-zero)
For Go-Zero, since code generation is required, only the code for the encryption middleware is provided, and you need to add it to your project yourself.


## 目录结构
```
├── gin     gin framework middleware example
├── go-zero go-zero framework middleware example
├── test    go test
├── utils   util functions
```


## Encryption Process
```
1. Get the request parameters.

2. Sort the request parameters in ascending order by key.

3. Concatenate the sorted request parameters into a string, e.g., age13namejack.

4. Append the timestamp and app secret to the string.

5. Encrypt the string using AES.

6. Set the request headers: timestamp (10 digits), signature, app ID.

```


## Verification Process
```
1. Verify the request headers.

2. Check if the timestamp has expired.

3. Sort the request parameters in the same way as the previous step, by sorting them in ascending order by key and concatenating them. Then, append the timestamp and app secret.

4. Encrypt the string and compare it with the received signature.

```

## Future Plans
- [ ] Support more encryption methods.
- [ ] Implement api replay protection.
- [ ] Support more request methods.