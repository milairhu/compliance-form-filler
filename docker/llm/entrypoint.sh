#!bin/bash


ollama serve & sleep 5 && \
ollama run compliance-model < /dev/null
wait $!