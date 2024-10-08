#!/bin/bash

package_name="smokep"
platforms=("windows/amd64" "linux/amd64" "darwin/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    
    flags=''
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
        flags+='-ldflags -H=windowsgui'
    fi

    set GOOS=$GOOS
    set GOARCH=$GOARCH

    go build $tags $flags -o "bin/"$output_name .
    if [ $? -ne 0 ]; then
        echo 'An error occurred! Aborting the script execution...'
        exit 1
    fi
done