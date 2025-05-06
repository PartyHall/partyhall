# PartyHall

PartyHall is the appliance software, that does the bulk of the features.

## Getting Started

PartyHall is available through an ansible playbook to be run on a stock debian setup.

To get started, you should visit our [documentation](https://partyhall.github.io/).

Note that the whole software has been rewritten in v0.8 and is no longer compatible at all with previous versions.

## Known issue

- Setup the hotspot doesn't work on real hw (run the script manually and it works)
	- 200 but bad request
	- Golang code only (the script works)
	- "The specified eth iface was not found"
	- Probably will fail too for the wifi one
	- Added logging in latest beta to check this out
	- iface passed to the script is empty
- Spotify no longer works (no clue why, I need to debug this??)

## Architecture todo
Some stuff are truely ugly in the architecture of the app. I know, I need to fix them. Here's a todo list.

- Frontend apps should have a global handling of errors (e.g. 401, 403, 500) through React Router
- All errors should be catched and AT LEAST display a snackbar notification
- The appliance settings / onboarding components should be remade properly, those are ugly and prone to issues I don't like it. (Like they should get all the data they need by themself and not rely on passed props, they should be able to decide whether the user CAN submit or not the form, ...)
- Make an eventbus (e.g. instead of doing mercure sendstate + other stuff by hand just say "bus.trigger("backdrop_updated")" and each packages register their events so that it handles everything by itself)
- Maybe we should try to build a REST-er api but that would need a lot of change including probably moving to a query-builder (or ORM) such as [bun](https://bun.uptrace.dev/) having better collection filtering / pagination (with [this](https://github.com/webstradev/gin-pagination) ?), ... It would also need to check how we would handle PUT method as they should only overwrite updated elements not the whole item (This could be done using a [JSON MERGE PATCH](https://github.com/evanphx/json-patch) library similarly to Api Platform but we need to check out how this library handles null value vs unset)
- The appliance interface should be translated, maybe check out the system language to set the default then allow overriding as UserSettings ?

## Links

- [Website / Docs](https://partyhall.github.io/)
- [Main software](https://github.com/partyhall/partyhall)
- [PartyNexus](https://github.com/partyhall/partynexus)
- [Docs repository](https://github.com/partyhall/docs)
