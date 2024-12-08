# PartyHall

## WIP
- Config PulseAudio / Pipewire (+ tester avec l'evo4)
- Maitriser le volume des micros dans PH + volume sonore global de Firefox

=> Some work done in the pipewire package
=> Need to add routes to settings to set the audio interface (both input & output)
=> + change the volume of output (microphone) and output (global => Karaoke mic + spotify + karaoke songs)
=> + Install the pipewire config in ansible (the file is already present in ansible/files/pipewire_karaoke.conf)
=> + Configure pipewire properly (add zamgate + zamcomp / amplify properly)

## TODOs
- When logout, mercure shits itself and spam connect attempts
- Fixer la carte mère de merde (Ou la remplacer avec celle d'antoine)

## VSCode config

CA MARCHE PAS

https://dev.to/andreidascalu/setup-go-with-vscode-in-docker-for-debugging-24ch

```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Delve into Docker",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "substitutePath": [
                {
                    "from": "<full absolute path to project>",
                    "to": "/app/",
                },
            ],
            "port": 2345,
            "host": "127.0.0.1",
            "showLog": true,
            "apiVersion": 2,
            "trace": "verbose"
        }
    ]
}
```
