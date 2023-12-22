test:
  go test github.com/shikaan/shmux/pkg/...

build:
  go build -o shmux-dev .
  mv shmux-dev ~/.local/bin/shmux-dev