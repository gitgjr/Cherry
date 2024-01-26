@echo off
SETLOCAL ENABLEDELAYEDEXPANSION

REM Path to the directory
SET "directory=E:\go\Distributed-meeting-recording-storage\data\serverWork"

REM Loop through all files in the directory
FOR %%F IN ("%directory%\*") DO (
    SET "filename=%%~nxF"
    
    REM Check if the file is not one of the exceptions
    IF "!filename!" NEQ "1080left.mp4" (
    IF "!filename!" NEQ "4k30left.mp4" (
    IF "!filename!" NEQ "left.mp4" (
    IF "!filename!" NEQ "1080right.mp4" (
    IF "!filename!" NEQ "4k30right.mp4" (
    IF "!filename!" NEQ "right.mp4" (
        REM Delete the file
        ECHO Deleting: %%F
        DEL "%%F"
    ))))))
)

ECHO File deletion process is complete.
PAUSE
