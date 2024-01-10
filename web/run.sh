#!/bin/sh
if [ -f ./cert/certificate.pem ]; then
    echo "File exists, running your command"
    npx serve --ssl-cert cert/certificate.pem --ssl-key cert/private-key.pem -s build
else
    echo "File does not exist"
    npx serve -s build
fi