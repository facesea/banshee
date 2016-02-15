# coding=utf8

"""
Fabfile (http://www.fabfile.org/) to deploy banshee to remote host.

Requirements: fabric

Example usage:

    $ python deploy.py -u hit9 -H remote-host:22 --remote-path "/srv/banshee"
      --remote-user root:root --refresh

This script will do the following jobs:

    1. Install static depdencies.
    2. Build static files.
    3. Build banshee binary.
    4. Rsync the static files and binary to remote host.
    5. Restart banshee service via supervisorctl.

Note:

    1. If the remote host is linux, this script should be also called
    on linux to build the right binary.
    2. The banshee service should be maintained in supervisor. You need
    to create a new service named banshee in supervisor.
"""

import os
import argparse

from fabric.api import (
    abort,
    env,
    execute,
    local,
    sudo,
    task,
    warn,
)
from fabric.contrib.files import exists
from fabric.contrib.project import rsync_project

###
# Global
###

LOCAL_DIR = "deploy-local-tmp"
LOCAL_STATIC_DIR = os.path.join(LOCAL_DIR, "static")
BINARY_NAME = "banshee"
STATIC_DIR = os.path.join("static", "dist")
SERVICE_NAME = "banshee"

###
# Local
###


def record_commit():
    """Record current commit.
    """
    local("git rev-parse HEAD > commit")


def install_static_deps():
    """Install local static dependencies.
    """
    local("cd static && npm install -q")
    local("cd static/public && npm install -q")


def build_static_files():
    """Build static files via gulp.
    """
    local("cd static && rm -rf dist/")
    local("cd static && gulp build")


def build_binary():
    """Build banshee binary via godep.
    """
    local("godep go build")


def make_local_dir():
    """Make local temporary directory.

        deploy-local-tmp/
            |- commit
            |- binary
            |- static/
                |- dist/
                    |- css/
                    |- js/
                    ...
    """
    local("mkdir -p {}".format(LOCAL_DIR))
    local("cp {0} {1}".format(BINARY_NAME, LOCAL_DIR))
    local("cp -r {0} {1}".format(STATIC_DIR, LOCAL_STATIC_DIR))
    local("mv commit {}".format(LOCAL_DIR))


def remove_local_dir():
    """Remove local temporary directory.
    """
    local("rm -rf {}".format(LOCAL_DIR))


###
# Remote
###


def upload():
    """Upload local directory to remote directory.
    """
    if not exists(env.remote_path):
        sudo("mkdir -p {}".format(env.remote_path))
    sudo("chmod 755 {}".format(env.remote_path))
    sudo("chown -R {0} {1}".format(env.user, env.remote_path))
    rsync_project(env.remote_path, LOCAL_DIR + '/')
    sudo("chown -R {0} {1}".format(env.remote_user, env.remote_path))


def refresh():
    """Refresh service via supervisor.
    """
    sudo("supervisorctl restart {}".format(SERVICE_NAME))


@task
def deploy():
    """Deploy banshee.
    """
    try:
        record_commit()
        install_static_deps()
        build_static_files()
        build_binary()
        make_local_dir()
        upload()
        if env.refresh:
            refresh()
    finally:
        remove_local_dir()


def main(host=None, user=None):
    parser = argparse.ArgumentParser()
    parser.add_argument('-u', '--user', help="user to connect")
    parser.add_argument('-H', '--host', help="host to deploy", required=True)
    parser.add_argument('--refresh', help="whether to refresh service",
                        action='store_true', default=False)
    parser.add_argument('--remote-path', help="remote service path",
                        required=True)
    parser.add_argument("--remote-user", help="remote service user",
                        default="root:root")
    args = parser.parse_args()

    if not args.user:
        warn("Using default user: {}".format(env.user))
    else:
        env.user = args.user

    if not args.remote_user:
        env.remote_user = "root:root"
        warn("Using default remote user: root:root")
    else:
        env.remote_user = args.remote_user

    env.remote_path = args.remote_path
    env.refresh = args.refresh
    env.use_ssh_config = True
    execute(deploy, hosts=[args.host])


if __name__ == '__main__':
    main()
