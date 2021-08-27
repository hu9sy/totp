### Generate Time-based One-time Password

-----
#### How to Use

```shell
# ** Add New Definition
# otpauth://totp/Example:hoge@example.com?secret=NBXWOZLGOVTWC===&issuer=Example
# totp add [issuer] [account name] [secret]
totp add Example hoge@example.com NBXWOZLGOVTWC===
```

```shell
# ** Generate
totp
> No     ISSUER      IDENTIFIER           TOTP      REMAINING
> 0      Example     hoge@example.com     376947    21
```

```shell
# ** Delete Definition
# totp delete [index]
totp delete 0
```
