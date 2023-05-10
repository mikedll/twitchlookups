
Looking for twitch videos at specific timestamps.

### Building

    > make
    go build -o bin/web_server cmd/web_server/main.go
    go build -o bin/cli cmd/cli/main.go

### Running

    > ./bin/cli pokimane "12:32 PM PDT Apr 11, 2023"
    Using timestamp of: Tue Apr 11, 2023 at 12:32pm PDT
    Found 20 possible videos
    Video URL: https://www.twitch.tv/videos/1791015541?t=44s

### References

[SO: Using internal packages](https://stackoverflow.com/questions/33351387/how-to-use-internal-packages)
