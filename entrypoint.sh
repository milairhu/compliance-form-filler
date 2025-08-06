#!/bin/sh
set -e

ARGS="answer"

[ -n "$SOURCE_FILE" ] && ARGS="$ARGS --source-file \"$SOURCE_FILE\""
[ -n "$OUTPUT_FILE" ] && ARGS="$ARGS --output-file \"$OUTPUT_FILE\""
[ -n "$QDRANT_URL" ] && ARGS="$ARGS --qdrant-url \"$QDRANT_URL\""
[ -n "$LLM_URL" ] && ARGS="$ARGS --llm-url \"$LLM_URL\""
[ -n "$EMBEDDING_API_URL" ] && ARGS="$ARGS --embedding-api-url \"$EMBEDDING_API_URL\""

while :; do
  STATUS=$(curl -s http://ingestor:8000/ready | jq -r '.status')
  [ "$STATUS" = "true" ] && break
  sleep 30
done

eval ./compliance-form-filler $ARGS