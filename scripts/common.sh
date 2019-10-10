#!/bin/bash

shw_grey () {
    echo $(tput bold)$(tput setaf 0) $@ $(tput sgr 0)
}

shw_norm () {
    echo $(tput bold)$(tput setaf 9) $@ $(tput sgr 0)
}

shw_info () {
    echo $(tput bold)$(tput setaf 4) $@ $(tput sgr 0)
}

shw_warn () {
    echo $(tput bold)$(tput setaf 2) $@ $(tput sgr 0)
}
shw_err ()  {
    echo $(tput bold)$(tput setaf 1) $@ $(tput sgr 0)
}


getOS(){
    OS=""
    VER=""
    if [[ -f /etc/os-release ]] ; then
        # freedesktop.org and systemd
        . /etc/os-release
        OS=$NAME
        VER=$VERSION_ID
    elif type lsb_release >/dev/null 2>&1; then
        # linuxbase.org
        OS=$(lsb_release -si)
        VER=$(lsb_release -sr)
    elif [[ -f /etc/lsb-release ]]; then
        # For some versions of Debian/Ubuntu without lsb_release command
        . /etc/lsb-release
        OS=$DISTRIB_ID
        VER=$DISTRIB_RELEASE
    elif [[ -f /etc/debian_version ]]; then
        # Older Debian/Ubuntu/etc.
        OS=Debian
        VER=$(cat /etc/debian_version)
    elif [[ -f /etc/SuSe-release ]]; then
        # Older SuSE/etc.
        ...
    elif [[ -f /etc/redhat-release ]]; then
        # Older Red Hat, CentOS, etc.
        ...
    else
        # Fall back to uname, e.g. "Linux <version>", also works for BSD, etc.
        OS=$(uname -s)
        VER=$(uname -r)
    fi

    echo ${OS}
}