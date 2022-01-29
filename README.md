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
