FROM ollama/ollama

COPY Modelfile /Modelfile

RUN ollama serve & sleep 10 && ollama create compliance-model -f /Modelfile

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 11434

ENTRYPOINT ["./entrypoint.sh"]
