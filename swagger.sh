swagger generate spec --output=./file.yml --scan-models

swagger serve --no-open -F=swagger --port 36666 file.yml