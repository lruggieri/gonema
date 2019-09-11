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

shw_norm "initializing project"

DISTRO_REDHAT=("Red Hat" "CentOS")
DISTRO_UBUNTU=("Ubuntu" "Mint")
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

OS="$(/bin/bash getOS.sh)"


if isDistro "$OS" "${DISTRO_REDHAT[@]}"; then
    shw_info "Red Hat distribution '$OS' detected"

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
elif isDistro "$OS" "${DISTRO_UBUNTU[@]}"; then
    shw_info "Ubuntu distribution '$OS' detected"

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

else
    shw_err "Distribution $OS not recognized. Aborting."
    exit 1
fi

#Google chrome is needed as the main tool used by chromedp to scrape the web
if command -v google-chrome >/dev/null 2>&1  ; then
    shw_norm "Google Chrome already installed"
else
    shw_norm "Installing Google Chrome"
    curl https://intoli.com/install-google-chrome.sh | bash
    shw_info "Google Chrome installed"
fi

#Now, proceed with Go dependencies installation; this step uses dep (https://github.com/golang/dep)
if command -v $GOPATH/bin/dep >/dev/null 2>&1 ; then
    shw_norm "Dep already installed"
else
    shw_info "Dep not installed! Installing..."

    #dep requires folder $GOPATH/bin to exist
    if [[ ! -d "$GOPATH/dep" ]]; then
        sudo mkdir $GOPATH/dep
    fi

    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

    if command -v $GOPATH/bin/dep ; then
        shw_info "Dep installed"
    else
        shw_err "Cannot install dep. Aborting"
        exit 1
    fi
fi

shw_norm "Installing project dependencies"
$GOPATH/bin/dep ensure
shw_info "Project dependencies installed"