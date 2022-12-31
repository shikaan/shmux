test:
  go test github.com/shikaan/shmux/pkg/...

build:
  go build -o shmux-dev .
  mv shmux-dev /usr/local/bin/shmux-dev

greet:
  echo "Hello $1, my old friend"