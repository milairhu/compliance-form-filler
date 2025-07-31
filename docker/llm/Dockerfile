FROM ollama/ollama

COPY Modelfile /Modelfile
RUN ollama serve
RUN ollama create compliance-mistral -f /Modelfile

CMD ["run", "compliance-mistral"]
