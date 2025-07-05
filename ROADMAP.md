# Roadmap for release 0.8 (final)

- Fully working onboarding process (and the setting pages)
- Fix hotspot
- Fix Spotify
- Changing Nexus HUB at runtime doesn't work (UI shows its ok but still tries to push to old URL)

# Long-term (to do before 1.0)

- Frontend apps should have a global handling of errors (e.g. 401, 403, 500) through Tanstack Router
- All errors should be catched and AT LEAST display a snackbar notification
- The appliance interface should be translated, maybe check out the system language to set the default then allow overriding as UserSettings ?

# Some day

- Make an eventbus (e.g. instead of doing mercure sendstate + other stuff by hand just say "bus.trigger("backdrop_updated")" and each packages register their events so that it handles everything by itself)
- Maybe we should try to build a REST-er api but that would need a lot of change including probably moving to a query-builder (or ORM) such as [bun](https://bun.uptrace.dev/) having better collection filtering / pagination (with [this](https://github.com/webstradev/gin-pagination) ?), ... It would also need to check how we would handle PUT method as they should only overwrite updated elements not the whole item (This could be done using a [JSON MERGE PATCH](https://github.com/evanphx/json-patch) library similarly to Api Platform but we need to check out how this library handles null value vs unset)