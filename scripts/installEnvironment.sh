#!/bin/bash

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/common.sh"

shw_norm "initializing project"

DISTRO_REDHAT=("Red Hat" "CentOS")
DISTRO_UBUNTU=("Ubuntu" "Mint" "Debian")
isDistro () {
    local os=$1
    shift
    local distros=("$@")
    

    for dist in "${distros[@]}";
    do
        if [[ "${os}" == *"${dist}"* ]]; then
            return 0
        fi
    done

    return 1
}
installGO () {
    # check and eventually install go (does not check for the version...)
    if command -v go >/dev/null 2>&1  ; then
        shw_norm "go already installed"
    else

        if [[ -d "/usr/local/go" ]]
        then
            shw_norm "go directory already present in /usr/local/go; please export your GOROOT/bin"
            return 0
        else

            sudo wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz
            sudo tar -xzf go1.12.7.linux-amd64.tar.gz
            sudo mv go /usr/local # error in case of directory not empty, which means go is already there
            sudo echo 'export GOROOT=/usr/local/go' >> ~/bashrc # PERMANENT
            sudo echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/bashrc # PERMANENT
            source ~/bashrc # reload bashrc in order for last exports to take effect

            sudo rm -rf go1.12.7.linux-amd64.tar.gz* go #removing go dir just in case we could not move it for some reason

            if command -v go >/dev/null 2>&1  ; then
                shw_info "go installed"
                return 0
            else
                return 1
            fi
        fi
    fi

    return 0
}

OS=$(getOS)

if isDistro "$OS" "${DISTRO_REDHAT[@]}"; then
    shw_info "Red Hat distribution '$OS' detected"


    if ! command -v sudo >/dev/bull2>&1 ; then
        yum install -y sudo

        if command -v sudo >/dev/null 2>&1  ; then
            shw_info "sudo installed"
        else
            shw_err "cannot install sudo"
            exit 1
        fi
    fi

    if ! command -v curl >/dev/bull2>&1 ; then
        yum install -y curl

        if command -v curl >/dev/null 2>&1  ; then
            shw_info "curl installed"
        else
            shw_err "cannot install curl"
            exit 1
        fi
    fi

    # check and eventually install git
    if command -v git >/dev/null 2>&1  ; then
     shw_norm "git already installed"
    else
        sudo yum install -y git
        if command -v git >/dev/null 2>&1  ; then
            shw_info "git installed"
        else
            shw_err "cannot install git"
            exit 1
        fi
    fi

    # check and eventually install wget
    if command -v wget >/dev/null 2>&1  ; then
     shw_norm "wget already installed"
    else
        sudo yum install -y wget
        if command -v wget >/dev/null 2>&1  ; then
            shw_info "wget installed"
        else
            shw_err "cannot install wget"
            exit 1
        fi
    fi

    #gcc and g++ are necessary for gosseract and underling libraries
    if command -v gcc >/dev/null 2>&1  ; then
     shw_norm "gcc already installed"
    else
        sudo yum install -y gcc
    fi
    if command -v g++ >/dev/null 2>&1  ; then
     shw_norm "g++ already installed"
    else
        sudo yum install -y gcc-c++
    fi


    if yum list installed "epel-release" >/dev/null 2>&1; then
        shw_norm "epel-release already installed"
    else
        shw_norm "Installing epel-release..."
        sudo yum install -y epel-release
        shw_info "epel-release installed"
    fi


    if command -v tesseract >/dev/null 2>&1  ; then
        shw_norm "tesseract already installed"
    else
        shw_norm "Installing tesseract (with leptonica)..."
        sudo yum install -y tesseract-devel leptonica-devel
        shw_info "tesseract with leptonica installed"
    fi

    #Google chrome is needed as the main tool used by chromedp to scrape the web
    if command -v google-chrome >/dev/null 2>&1  ; then
        shw_norm "Google Chrome already installed"
    else
        shw_norm "Installing Google Chrome"
        curl https://intoli.com/install-google-chrome.sh | bash
        shw_info "Google Chrome installed"

    fi

elif isDistro "$OS" "${DISTRO_UBUNTU[@]}"; then
   shw_info "Ubuntu distribution '$OS' detected"


   if ! command -v sudo >/dev/bull2>&1 ; then
        apt-get update ; apt-get -y install sudo

        if command -v sudo >/dev/null 2>&1  ; then
            shw_info "sudo installed"
        else
            shw_err "cannot install sudo"
            exit 1
        fi
   fi

   if ! command -v curl >/dev/bull2>&1 ; then
        apt-get update ; apt-get -y install curl

        if command -v curl >/dev/null 2>&1  ; then
            shw_info "curl installed"
        else
            shw_err "cannot install curl"
            exit 1
        fi
   fi

    # check and eventually install git
    if command -v git >/dev/null 2>&1  ; then
        shw_norm "git already installed"
    else
        sudo apt-get install -y git

        if command -v git >/dev/null 2>&1  ; then
            shw_info "git installed"
        else
            shw_err "cannot install git"
            exit 1
        fi
    fi

    # check and eventually install wget
    if command -v wget >/dev/null 2>&1  ; then
        shw_norm "wget already installed"
    else
        sudo apt-get install -y wget
        if command -v wget >/dev/null 2>&1  ; then
            shw_info "wget installed"
        else
            shw_err "cannot install wget"
            exit 1
        fi
    fi

    #gcc and g++ are necessary for gosseract and underling libraries
    #on Ubuntu, installing build-essential suffice
    if [ $(dpkg-query -W -f='${Status}' build-essential 2>/dev/null | grep -c "ok installed") -eq 0 ];
    then
        shw_info "Installing build-essential"
        sudo apt install build-essential -y
    else
        shw_norm "build-essential already installed"
    fi


    if command -v tesseract >/dev/null 2>&1  ; then
        shw_norm "tesseract already installed"
    else
        shw_norm "Installing tesseract with leptonica..."
        sudo apt-get install -y libtesseract-dev libleptonica-dev tesseract-ocr-eng
        shw_info "tesseract with leptonica installed"
    fi

    #Google chrome is needed as the main tool used by chromedp to scrape the web
    if command -v google-chrome >/dev/null 2>&1  ; then
        shw_norm "Google Chrome already installed"
    else

        shw_norm "Installing Google Chrome"
        wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
        dpkg -i google-chrome-stable_current_amd64.deb; apt-get -fy install
        shw_info "Google Chrome installed"

    fi

else
    shw_err "Distribution $OS not recognized. Aborting."
    exit 1
fi

if ! installGO ; then
    exit 1
fi