# set GOPATH
uname -a | grep -i cygwin > /dev/null
if [[ $? -eq 0 ]]; then
	export GOPATH=$(cygpath -w $(pwd))
else
	export GOPATH=$(pwd)	
fi


# build and put executable to $GOPATH/bin/
go install s3cl
