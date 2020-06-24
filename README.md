# BME280 climate data logging

Get data from a BME280 sensor, put it in an Influx database and optionally log it to an LCD screen.

## Starting a systemd service

Put in /lib/systemd/system/i2c_temp_humid.service

```
[Unit]
Description=i2c temperature, humidity and pressure monitor
ConditionPathExists=/usr/bin/i2c_temp_humid
After=network.target

[Service]
Type=simple
User=root
Group=adm
LimitNOFILE=1024

Restart=on-failure
RestartSec=10

ExecStart=/usr/bin/i2c_temp_humid -name study -influxhost 192.168.1.30

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/i2c_temp_humid
ExecStartPre=/bin/chown root:adm /var/log/i2c_temp_humid
ExecStartPre=/bin/chmod 755 /var/log/i2c_temp_humid
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=i2c_temp_humid

[Install]
WantedBy=multi-user.target
```

Then:

```
sudo systemctl edit myservice
```

Make the file look like this:

```
[Service]
Environment="INFLUX_DB_SECRET=<influx DB password goes here>"
```

Then enable and start the service.

```
sudo systemctl unmask i2c_temp_humid.service
sudo systemctl start i2c_temp_humid
sudo systemctl enable i2c_temp_humid.service
```
