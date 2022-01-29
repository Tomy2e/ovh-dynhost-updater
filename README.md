# ovh-dynhost-updater

A tool to automatically update an OVH DynHost record with your public IP address.

It can be run as a long running process (daemon) or as a regular cronjob task.

## Usage

### Options

The tool accepts the following command-line options:

| Name      | Description                                                                   | Default value |
| --------- | ----------------------------------------------------------------------------- | ------------- |
| -daemon   | Run the program as a daemon.                                                  | `false`       |
| -delay    | If the program is run as a daemon, number of minutes to wait between updates. | `10`          |
| -hostname | Hostname of the DynHost record to update.                                     |               |

The tool reads the following environment variables:

| Name     | Description                                                                       | Default value |
| -------- | --------------------------------------------------------------------------------- | ------------- |
| LOGIN    | DynHost login. The tool will exit if this environment variable is not defined.    |               |
| PASSWORD | DynHost password. The tool will exit if this environment variable is not defined. |               |

## Installation

### cron job

Run the following commands to download the tool:

```bash
sudo wget https://github.com/Tomy2e/ovh-dynhost-updater/releases/latest/download/ovh-dynhost-updater-linux-amd64 -O /usr/local/bin/ovh-dynhost-updater
sudo chmod +x /usr/local/bin/ovh-dynhost-updater
```

Update the `CHANGEME` values then run the following commands to create the file
that contains your DynHost credentials:

```bash
sudo bash -c "cat >> /usr/local/etc/ovh-dynhost-updater-creds" << EOL
export LOGIN=CHANGEME
export PASSWORD=CHANGEME
EOL
sudo chmod 0700 /usr/local/etc/ovh-dynhost-updater-creds
```

Open the crontab file:

```bash
sudo crontab -e
```

Append the following line to your crontab file, update the `CHANGEME` values,
then save:

```bash
*/10 * * * * (. /usr/local/etc/ovh-dynhost-updater-creds && /usr/local/bin/ovh-dynhost-updater --hostname=CHANGEME) 2>&1 | logger -t ovh-dynhost-updater
```

The tool will run every ten minutes, run the following command to see the logs:

```bash
grep 'ovh-dynhost-updater' /var/log/syslog
```
