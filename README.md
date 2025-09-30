# Compliance Auto-Filler

## Description
The **Compliance Auto-Filler** is a tool that automates the completion of compliance questionnaires (RFPs, audits, vendor forms) using an organization's internal policies and documentation.

### How it works
It leverages:

- Qdrant as a vector database to store and query compliance policies and documents.

- A Python API embedding a model to handle document ingestion, chunking, and vectorization.

- A locally hosted LLM deployed with Ollama, preprompted to generate short, ready-to-use answers.

- An Go CLI (this project) to read questions, feed policies and documents into the Qdrant database, ask the LLM questions with relevant context, and output responses.

### Input
So far, the CLI is fed with a text file containing questions.

### Output
The output is a csv file with questions and answers.

## Architecture

![Filler architecture diagram](./img/architecture.png "Compliance Auto-Filler atchitecture")

## Getting Started

The prerequisites are:
- have the compliance corpus ready in a **./corpus** folder
- have the questions ready in a **./questions.txt** file

```bash
docker compose up --build -d
```

This will start:
- Qdrant on http://localhost:6333
- The Ingestor on http://localhost:8000
- The Compliance Auto-Filler CLI
- A local LLM onn http://localhost:11434

The output will be in **./answers.csv**.
