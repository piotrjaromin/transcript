#!/bin/sh

if [ -n "$WHISPER_MODEL" ]; then
    if [ ! -f "/app/models/$WHISPER_MODEL" ]; then
        echo "Downloading custom model: $WHISPER_MODEL"dockerfle
        wget -O "/app/models/$WHISPER_MODEL" "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/$WHISPER_MODEL"
    fi
    MODEL_PATH="/app/models/$WHISPER_MODEL"
else
    MODEL_PATH="/app/models/${DEFAULT_MODEL}"
fi

/app/transcript --model "$MODEL_PATH" "$@"
