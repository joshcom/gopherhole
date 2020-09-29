# Installing and Running as a systemd service

## WIP

0. build/install gopherhole.  Place binary in /usr/local/bin, or elsewhere.
1. Create a gopher user and group, nologin.
2. Create a configuration file, place it in /etc/gopherhole/, and make sure it is readable.
3. Create your gopherhole directory.  /var/gopherhole.  Make sure gopher.gopher can read it.
4. Set up systemd.

```
sudo touch /etc/systemd/system/gopherhole.service
sudo chmod 644 /etc/systemd/system/gopherhole.service
```

```
[Unit]
Description=Gopherhole Server

[Service]
Type=simple
User=gopher
Group=gopher
ExecStart=/usr/local/bin/gopherhole -config="/path/to/config.json"

[Install]
WantedBy=multi-user.target 
```

```
sudo systemctl daemon-reload
sudo systemctl enable gopherhole.service
```
