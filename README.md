# PartyHall

PartyHall is the appliance software, that does the bulk of the features.

## Getting Started

PartyHall is available through an ansible playbook to be run on a stock debian setup.

To get started, you should visit our [documentation](https://partyhall.github.io/).

Note that the whole software has been rewritten in v0.8 and is no longer compatible at all with previous versions.

## Architecture todo
Some stuff are truely ugly in the architecture of the app. I know, I need to fix them. Here's a todo list.

- Frontend apps should have a global handling of errors (e.g. 401, 403, 500) through React Router
- All errors should be catched and AT LEAST display a snackbar notification
- The appliance settings / onboarding components should be remade properly, those are ugly and prone to issues I don't like it. (Like they should get all the data they need by themself and not rely on passed props, they should be able to decide whether the user CAN submit or not the form, ...)
- Make an eventbus (e.g. instead of doing mercure sendstate + other stuff by hand just say "bus.trigger("backdrop_updated")" and each packages register their events so that it handles everything by itself)

## Links

- [Website / Docs](https://partyhall.github.io/)
- [Main software](https://github.com/partyhall/partyhall)
- [PartyNexus](https://github.com/partyhall/partynexus)
- [Docs repository](https://github.com/partyhall/docs)