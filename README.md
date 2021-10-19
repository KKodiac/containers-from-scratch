# containers-from-scratch

Creating a docker-like container service from scratch using Go

Go를 사용한 도커와 비슷한 컨테이너 서비스 시작부터 만들기
## 시작하는 방법
해당 코드는 우분투 20.04 lts 에서 확인하였습니다.

```sh
    # 해당 저장소를 클론 합니다.
    git clone https://github.com/KKodiac/containers-from-scratch.git
    # main.go 를 빌드 후 executable 생성합니다.
    go build main.go 
    # 루트 권한으로 실행합니다.
    # 컨테이너 프로세스를 bash쉘로 모니터 합니다.
    sudo ./main run /bin/bash
```

## namespace 
[도커](https://docs.docker.com/get-started/overview/#the-underlying-technology) 등 컨테이너는 `namespace`라는 기법을 활용해서 각 컨테이너에게 독립된 환경설정을 제공합니다. 

자세한 것은 [코드]()를 확인해주세요.

## 컨테이너 확인
위 처럼 시작 하였으면 bash로 컨테이너 루트(`/`)로 설정된 디렉토리에서 시작하게 됩니다.
해당 컨테이너 쉘의 호스트 이름은 `container`로 설정되어 있습니다.
```sh
    # root@container:/#
    ls # '/' 디렉토리 내용 출력
```


## 최대 프로세스 제한 
해당 컨테이너는 생성할 수 있는 최대 프로세스가 `20`으로 제한되어 있습니다.
```sh
    # 컨테이너 프로세스를 실행합니다.
    sudo ./main run /bin/bash
```
```sh
    # Fork Bomb을 생성해서 확인해 봅니다. 무한의 프로세스 포크를 생성합니다.
    # 컨테이너 내부에서 실행해야 됩니다! 호스트 쉘에서 실행하면... ㅎㅇㅌ...
    root@container:/ :(){ :|: & }; :
```

호스트 OS(Linux Ubuntu)에서 확인하면,
```sh
    cd /sys/fs/cgroup/pids/ubuntu 
    cat pids.max # 20
    # 호스트에서 cgroup/pids/ubuntu 를 확인해보면,
    # 컨테이너로 인해 생성된 프로세스가 그렇게 많지 않은 것을 확인할 수 있습니다. (20개)
    ps aux 
```

## /proc 내용 확인
시스템 내부의 프로세스 ID, 자원 등을 확인할 수 있습니다.

호스트에서 `mount` 한 /proc 디렉토리를 확인합니다.
```sh
    ls /proc # 프로세스 ID 등등
    ps # 컨테이너 내부의 프로세스 리스트
```