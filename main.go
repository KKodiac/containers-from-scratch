package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	fmt.Printf("Running %v \n", os.Args[2:])

	// 자기 자신의 자식 프로세스 생성
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	// 쉘 기본 IO 추가
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Cloneflags 는 오직 리눅스에서 가능합니다
	// 아래는 컨테이너 namespace를 설정해 줍니다.
	// CLONE_NEWUTS 독립된 hostname용 namespace
	// CLONE_NEWPID 독립된 process용 namespace
	// CLONE_NEWNS 독립된 mounts용 namespace
	// 자세한 내용은 아래 확인
	// https://pkg.go.dev/syscall#pkg-constants
	// https://pkg.go.dev/syscall#SysProcAttr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as Child\n", os.Args[2:])
	cg()
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	// 쉘 기본 IO 설정
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 컨테이너 호스트 이름을 'container' 로 설정
	must(syscall.Sethostname([]byte("container")))
	// 루트 디렉토리 설정을 호스트의 /home/ubuntu/ubuntufs 에서 가져와 설정
	must(syscall.Chroot("/home/ubuntu/ubuntufs"))
	// 컨테이너 쉘 시작을 '/' 로 설정
	must(os.Chdir("/"))
	// ps 명령어가 컨테이너에서 돌아가도록 /proc 마운트
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	// 임시 파일시스템 마운트
	// must(syscall.Mount("thing", "mytemp", "tmpfs", 0, "")) // file not found error

	must(cmd.Run())

	// 마운트 클린업
	must(syscall.Unmount("proc", 0))
	must(syscall.Unmount("thing", 0))
}

// 호스트의 control group 패스에 설정 초기값 설정
func cg() {
	cgroups := "/sys/fs/cgroup/" // 리눅스 cgroup 폴더 위치
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "ubuntu"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "ubuntu/pids.max"), []byte("20"), 0700))                          // 컨테이너 내 최대 지정 가능 PID를 20개로 설정
	must(ioutil.WriteFile(filepath.Join(pids, "ubuntu/notify_on_release"), []byte("1"), 0700))                  // cgroup 미사용 시 클린업
	must(ioutil.WriteFile(filepath.Join(pids, "ubuntu/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)) // cgroup 내 있는 모든 프로세스에 적용
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
