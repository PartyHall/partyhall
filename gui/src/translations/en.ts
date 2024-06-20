export const ENGLISH = {
    general: {
        error_occured: "An error occured",
        something_went_wrong: "Something went wrong",
        cancel: "Cancel",
        dt_not_available: "Datetime not available",
        yes: 'Yes',
        no: 'No',
        version: 'Version',
    },
    osd: {
        no_event: "No event selected !"
    },
    login: {
        name: 'Displayed name',
        username: 'Username',
        password: 'Password',
        bt: 'Login',
        failed: 'Failed to login',
        logout: 'Logout',
    },
    admin_main: {
        partyhall: {
            current_event: "Selected event",
            new_event: "New event",
            logged_in_as: "Logged in as {{name}}"
        },
        mode: "Mode",
        hw_flash: "Hardware flash",
        system_info: "System infos",
        current_time: "Current time",
        set_to_my_time: "Set to my device's time",
        show_debug_info: "Show debug infos (30 sec)",
        shutdown: {
            title: "Shutdown",
            text: "You are trying to shutdown the partyhall. Are you sure ?",
            bt: "Shut down",
        },
        change_event: {
            title: "Change event",
            content: 'You are updating the current event to "{{event}}" (by {{author}})<br />Doing so will make that all new pictures are sent to this event instead of the current one.',
            bt: "Change event"
        },
        settings: "Global",
        photobooth: "Photobooth",
        karaoke: "Karaoke",
        volume: "Main volume",
        device: "Device",
        no_devices: "No audio devices detected"
    },
    exports: {
        last_exports: "Last exports",
        file: "Files",
        date: "Date",
        export_as_zip: "Export as zip",
        export: "Export",
        failed_to_download: "Failed to download the file",
        modal_infos: "You're trying to export the event {{name}}. <br /> This will create a zip with all the pictures and let you download, thus it could take a long time. <br /> Are you sure you want to continue ?",
        started: "Export started",
        completed: "Export completed",
    },
    event: {
        create: "Create an event",
        edit: "Editing {{name}}",
        name: "Title",
        host: "Host",
        date: "Date",
        location: "Location",
        save: "Save",
        saved: "Event saved!",
        failed: "Failed to save event"
    },
    karaoke: {
        no_song_playing: "No music playing",
        now_playing: "Now playing",
        sung_by: "Sung by",
        next_up: "Next up",
        search: "Search",
        queue: "Queue",
        admin: "Admin",
        amt_songs: "Amt songs",
        current: "Current",
        what_to_do: "What to do ?",
        wtd_add_song: "Add a song",
        wtd_import: "Import",
        wtd_settings: "Settings",
        title: "Title",
        artist: "Artist",
        cover_source: "Cover source",
        no_cover: "No cover",
        uploaded: "Uploaded",
        format: "Format",
        upload: "Upload {{ format }}",
        add: "Create",
        created: 'Music created',
        failed: 'Failed to create song',
        waiting_song_upload: "Waiting for a song upload",
        upload_a_song: "Upload a song",
        upload_in_progress: "Uploading...",
        rescan_songs: "Re-scan songs",
        songs_rescanned: "Song scanning completed",
        song_rescan_failed: "Failed to re-scan songs",
        adding_to_queue: 'Added to queue',
        empty_queue: "Empty queue",
        skip: {
            title: "Skip the song?",
            content: "There is currently a song, you might want to add to the queue.",
            skip_and_play: "Skip & play",
            skip: "Skip",
        },
        queue_remove: {
            title: "Remove from the queue",
            content: "Are you sure you want to remove {{name}} from the queue?",
            removed: "Removed from the queue"
        },
        current_remove: {
            title: "Stop the music",
            content: "Are you sure you want to remove the current song ({{name}}) ?",
            remove: "Stop",
        },
        volume: {
            title: "Volume",
            instrumental: "Instrumental",
            vocals: "Voices",
            full: "Full"
        }
    },
    photobooth: {
        remote_take_picture: "Remote picture",
        amt_hand_taken: "Amt taken manually",
        amt_unattended: "Amt unattended",
    },
    disabled: {
        partyhall_disabled: 'PartyHall disabled',
        msg: 'This PartyHall booth is currently disabled by an admin. Sorry for the inconvenience.'
    }
};