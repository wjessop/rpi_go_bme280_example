all:
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/i2c_temp_humid
	scp bin/i2c_temp_humid pi@<your rpi IP>:~/
