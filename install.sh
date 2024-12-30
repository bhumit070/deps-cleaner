#!/bin/bash

shellFile=""
if [ -n "$($SHELL -c 'echo $ZSH_VERSION')" ]; then
	shellFile=".zshrc"
elif [ -n "$($SHELL -c 'echo $BASH_VERSION')" ]; then
	shellFile=".bashrc"
else
	echo "no installation script found for this, please install manually"
	exit 0
fi

platform=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
	platform="linux-amd64"
elif [[ "$OSTYPE" == "darwin"* ]]; then
	if [[ $(sysctl -n machdep.cpu.brand_string) == *"Apple"* ]]; then
		platform="darwin-arm64"
	else
		platform="darwin-amd64"
	fi
else
	echo "platform not supported"
	exit 0
fi

platform="deps-cleaner-$platform"
tagName=$(curl -s https://api.github.com/repos/bhumit070/deps-cleaner/releases/192679296 | awk -F'"tag_name":' '{print $2}' | awk -F'"' '{print $2}' | xargs)
downloadableUrl="https://github.com/bhumit070/deps-cleaner/releases/download/$tagName/$platform"

echo "downloading..."
installDir="$HOME/.deps-cleaner"
destinationPath="$installDir/deps-cleaner"
mkdir -p "$installDir"
curl -# -L "$downloadableUrl" -o "$destinationPath" && chmod +x "$destinationPath"

command="export PATH=\$PATH:\$HOME/.deps-cleaner"
if grep -q "$command" "$HOME/$shellFile"; then
	echo "already installed."
	exit 0
fi

echo $command >>"$HOME/$shellFile"

echo "installation completed."
