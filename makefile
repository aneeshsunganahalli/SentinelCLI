build:
	go build -ldflags="-X 'github.com/aneeshsunganahalli/SentinelCLI/cmd.Version=1.0.0'" -o sentinel main.go

install:
	go install