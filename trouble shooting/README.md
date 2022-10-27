## Golang struct variable name

Marshal, Unmarshal 을 하기 위해 struct 에 변수를 만들었다면 첫글자를 반드시 대문자로 해야한다.

```go
type Book struct {
  name string `json:"Name"`
} (X)

type Book struct {
  Name string `json:"Name"`
} (O)
```
## MongoDB 설치 오류

~~해당 문제는 Ubuntu 의 버전에 따라 차이가 있다.~~  
~~기존 Ubuntu 에 설치되어 있는 mongodb 와 설치하려는 mongodb 의 충돌로 인하여 동작이 제대로 되지 않을 때 수행한다.~~  
처음에 받은 이미지에 MongoDB 가 이미 설치되어 있었던 것으로 보인다. 버전 충돌로 인해 삭제 후 재설치했다.

1. MongoDB 삭제
```
$ apt remove mongodb-org*
$ apt-get purge mongodb-org*
$ rm -r /var/log/mongodb
$ rm -r /var/lib/mongodb
$ rm -r /etc/mongodb.conf
```

2. 확인 (+삭제)
   * `$ apt list --installed | grep mongo` 를 통해 mongodb 가 남아있는지 확인한다.
   * `$ dpkg --list` 를 통해 mongodb 관련 패키지를 찾아 제거한다.
   * ex) `$ sudo apt-get --purge remove mongodb*`

3. 재설치

