FROM ollama/ollama

COPY Modelfile /Modelfile

RUN ollama serve & sleep 10 && ollama create compliance-model -f /Modelfile

CMD ["run", "compliance-model"]
