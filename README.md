# ec2-metadata-environment

A small utility exports EC2 metadata as environment variables. It's intended to be used as a systemd service, and expose the information to be used for other systemd services.

## Example

A systemd unit override, which writes the `EC2_SSH_PUBLIC_KEYS` variable to the authorized keys file before the ssh daemon starts.

`/etc/systemd/system/sshd.service.d/10-update-ec2-public-keys.conf`:

```
[Unit]
Wants=ec2-metadata-environment.service
After=ec2-metadata-environment.service

[Service]
EnvironmentFile=/etc/ec2/environment
ExecStartPre=/bin/sh -c "/bin/echo ${EC2_SSH_PUBLIC_KEYS} > /etc/ssh/authorized_keys"
```

## Variables

The following variables are exported:

* `EC2_ACCOUNT_ID`
* `EC2_ARCHITECTURE`
* `EC2_AVAILABILITY_ZONE`
* `EC2_IMAGE_ID`
* `EC2_INSTANCE_ID`
* `EC2_INSTANCE_TYPE`
* `EC2_KERNEL_ID`
* `EC2_LOCAL_HOSTNAME`
* `EC2_LOCAL_IPV4`
* `EC2_MAC_ADDRESS`
* `EC2_PENDING_TIME`
* `EC2_PUBLIC_HOSTNAME`
* `EC2_PUBLIC_IPV4`
* `EC2_RAMDISK_ID`
* `EC2_REGION`
* `EC2_RESERVATION_ID`
* `EC2_SECURITY_GROUPS` - comma-separated list of security groups
* `EC2_SSH_PUBLIC_KEYS` - newline-separated list of public SSH keys
* `EC2_VPC_DNS_SERVER_ADDRESS` - the IP address of the Amazon DNS server (the second IP address in the VPC CIDR block)
* `EC2_VPC_ID`
* `EC2_VPC_IPV4_CIDR_BLOCK`

