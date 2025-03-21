# Transcript project

Project purpose is to take in as input audio files and as output provide transcript.
To make task simpler optional parameter can be provided with expected audio language.

Application should be able to run in 3 modes:

1. Http server with endpoint where audio data is provided and as result transcript is returned
2. Standalone cli tool with option to provide path to file and transcript is returned
3. Standalone cli tool with recording option. Cli tool after entering "ENTER" starts recording audio, after another "ENTER" input stops recodring and provides transcript of recorded audio.

Cli tools should take in arguments which will switch between used modes.

## Goals

- Take in audio file and produced speech to text as result
- Option to run as http server
- Option to run as cli tool
- Option to include speech to text model inside resulting binary
- Option to provide path to external speech to text model 
- Dockerized environemt
- as minimum project should understand polish and english language
- README.md file should be updated with usage examples, and all relevant information on how to run and build project.
- .gitignore should be updated for files which should not be commitd to git repository
- Github actions pipelinie should be created for pushing docker images into docker hub repository

## Tools and libraries

- application should be written in golang
- module responsible for speech to text is https://github.com/ggerganov/whisper.cpp and its golang binding
- 
