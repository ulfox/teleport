---
authors: Alex McGrath (alex.mcgrath@goteleport.com)
state: draft
---

# RFD 57 - Automatic user and sudoers provisioning

## What

Automatically create non-existing users and optionally add them to
`sudoers` on Teleport nodes. Users will be removed after all sessions
have logged out.

## Why

Currently, when logging into an SSH node, the user must be
pre-created. Adding automatic user and `sudoer` provisioning would
make it so that any Teleport user would be able to login and have the
account created automatically without manual intervention.

## Details

The following are required for this feature:

- Ability to automatically provision a Linux user if it's not present
  on the node.
- Ability to automatically provision a Linux group if it's not present
  on the node.
- Ability to add the provisioned user to existing Linux groups defined
  in the user traits/role.
- Ability to add the provisioned user to sudoers.
- Clean up the provisioned user / sudoers changes upon logout (being
  careful not to remove pre-existing users).

### Config/Role Changes

Several new fields will need to be added to to the role `options` and
`allow` sections:

```yaml
kind: role
version: v5
metadata:
  name: example
spec:
  options:
    # Controls whether this role supports auto provisioning of users.
    create_host_user: true
  allow:
    # New field listing Linux groups to assign a provisioned user to.
    # Should support user and identity provider traits like other fields (e.g. "logins")
    host_groups: [ubuntu, "{{internal.groups}}", "{{external.xxx}}"]
    # host_sudoers is a list of entries to be included in a users sudoers file
    host_sudoers: ["{{internal.logins}} ALL=(ALL) ALL", ...]
```

An individual `ssh_service` can be configured disable auto user
creation with the below config:

```yaml
ssh_service:
    # when disabled, takes precedence over the role setting
    disable_create_host_user: true
```

### User creation

In order to create users `useradd` will be executed from teleport
after a user has tried to access a Teleport SSH node.

#### User Groups

When a user is created they will be added to the specified groups from
the `host_groups` field in the role. In addition the user will be
added to a special `teleport-system` group which can be used to
indicate that the user was created by teleport and that its safe for
it to be deleted. The groups will be created via `groupadd` at startup
if they do not already exist and users will be added to groups via
`usermod -aG <list of groups> <username>`

#### Valid user/group names

The set of valid names that are valid on Linux varies between distros
and are generally more restrictive than the allowed usernames in
Teleport. This will require that names containing invalid characters
have those characters removed/replaced. Information on the valid
characters between Linux distros is available [here](https://systemd.io/USER_NAMES/).
The common core of valid characters is `^[a-z][a-z0-9-]{0,30}$`.

#### Adding and removing users from sudoers

Each user with entries in `host_sudoers` will have a file created in
`/etc/sudoers.d`, with one entry per line.

If a user is in multiple rules that specify `host_sudoers` they will
be all be concatenated together.

##### sudoers file syntax validation

If a system has `visudo` present, validation could be performed by
executing `visudo -c -f path/to/sudoersfile`, where if it fails to
validate, the user fails to have the shell start and the error is
reported.

##### sudoers security considerations

In order to stop users from being able to edit the sudoers file a
command allow list must be used, as or equivalent to below:

```
${USERNAME} ALL = (${USER TO RUN AS}) NOPASSWD: /bin/cmd1 args, /bin/cmd2 args
```

Should a user be given `root` access to all commands, they will be
able to modify any file, including sudoers files.


### User and group deletion

After all of a users sessions are logged out the created user and any
`sudoers` files that were created for that user will be deleted if
that user is also a member of the `teleport-system` group.

Users can not be deleted while they have running processes so each
time a session ends, an attempt to delete the user can happen, if it
succeeds the sudoers file can also be removed.

If it does not succeed a cleanup process will run every 5 minutes, that
will attempt to delete users if they no longer have running processes.
This clean up process will also ensure that users with running
sessions during a restart will be cleaned up appropriately.

Groups will not be cleaned up and will be created once and be reused
this is to avoid files created with specified groups will remain
accessible between sessions to users in those groups.

### Multiple matching roles 

Automatic user provisioning will require that all roles matching a
node via `labels` have `create_host_user=true`

## UX Examples

### Teleport admin wants each user to have a dedicated host user defined by their Okta attributes
```yaml
kind: role
version: v5
metadata:
  name: auto-user-groups
spec:
  options:
    # allow auto provisioning of users.
    create_host_user: true
  allow:
    # username from external okta attribute
    logins: [ "{{external.username}}" ]
```

### Teleport admin wants to define which Linux groups each auto-created user will be added to

```yaml
kind: role
version: v5
metadata:
  name: auto-user-groups
spec:
  options:
    # allow auto provisioning of users.
    create_host_user: true
  allow:
    # List of each group the user will be added to
    host_groups: [ubuntu, docker, ...]
    # username from external okta attribute
    logins: [ "{{external.username}}" ]	
```

### Teleport admin wants to make each auto-created user a sudoer

```yaml
kind: role
version: v5
metadata:
  name: users-as-sudoers
spec:
  options:
    # allow auto provisioning of users.
    create_host_user: true
  allow:
    # add users to the wheel group
    host_groups: [wheel]
    # make it so users in the wheel group will be able to execute sudoers commands without a password
    host_sudoers: ["%wheel ALL=(ALL) NOPASSWD: ALL"]
```

### Teleport admin wants to define particular commands user will be able to run as root
```yaml
kind: role
version: v5
metadata:
  name: specify-commands-as-sudoers
spec:
  options:
    # allow auto provisioning of users.
    create_host_user: true
  allow:
    # make it so this specific user can execute `systemctl restart nginx.service `
    host_sudoers: ["{{internal.logins}} ALL = (root) NOPASSWD: /usr/bin/systemctl restart nginx.service"]
```

### Teleport admin wants to prohibit some nodes from auto-creating users

Include the below config for the Teleport node that should not allow automatic user creation:

```yaml
ssh_service:
  enabled: "yes"
  # stops a specific node from auto-creating users
  disable_create_host_user: true
```

Nodes where `diable_create_host_user` is `false` will still be able to
have users be automatically created.

### Teleport user has multiple roles but not all of them enable `create_host_user`

In the situtation where a user has roles as below, the user would not
be able to make use of automatically provisioning users as both roles
do not enable `create_host_user`.

```yaml
kind: role
version: v5
metadata:
  name: allow-access-and-auto-create
spec:
  options:
    # allow auto provisioning of users.
    create_host_user: true
    node_labels:
      - 'env': 'example'	
```

```yaml
kind: role
version: v5
metadata:
  name: specify-commands-as-sudoers
spec:
  options:
    node_labels:
      - 'env': 'example'
```
