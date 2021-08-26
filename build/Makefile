build:
	env GOOS=linux GOARCH=amd64 go build -o telegramupworkbot github.com/flasherup/telegrum_upwork_bot/cmd


upload:
	pscp telegramupworkbot flasherup@91.238.103.109:/home/flasherup/utils
	pscp config.yml flasherup@91.238.103.109:/home/flasherup/utils
