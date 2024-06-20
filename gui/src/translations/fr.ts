export const FRENCH = {
    general: {
        error_occured: "Une erreur est survenue",
        something_went_wrong: "Quelque chose s'est mal passé",
        cancel: "Annuler",
        dt_not_available: "Heure non disponible",
        yes: 'Oui',
        no: 'Non',
        version: 'Version',
    },
    osd: {
        no_event: "Pas d'évènement sélectionné !"
    },
    login: {
        name: 'Nom affiché',
        username: 'Pseudo',
        password: 'Mot de passe',
        bt: 'Connexion',
        failed: 'Échec de la connexion',
        logout: 'Déconnexion',
    },
    admin_main: {
        partyhall: {
            current_event: "Event sélectionné",
            new_event: "Nouvel évènement",
            logged_in_as: "Connecté en tant que {{name}}"
        },
        mode: "Mode",
        hw_flash: "Flash matériel",
        system_info: "Info système",
        current_time: "Heure serveur",
        set_to_my_time: "Configurer à mon heure actuelle",
        show_debug_info: "Afficher les info de debug (30s)",
        shutdown: {
            title: "Éteindre",
            text: "Vous tentez d'éteindre la machine, êtes-vous sûr ?",
            bt: 'Éteindre',
        },
        change_event: {
            title: "Changer d'évènement",
            content: 'Vous allez changer l\'évènement pour "{{event}}" (Par {{author}})<br />Toutes les nouvelles photos iront sur cet évènement.',
            bt: "Changer"
        },
        settings: "Global",
        photobooth: "Photomaton",
        karaoke: "Karaoké",
        volume: "Volume principal",
        device: "Carte son",
        no_devices: "Aucune carte son détectée"
    },
    exports: {
        last_exports: "Derniers exports",
        file: "Fichier",
        date: "Date",
        export_as_zip: "Exporter en zip",
        export: "Exporter",
        failed_to_download: "Échec du téléchargement du fichier",
        modal_infos: "Vous allez exporter l'évènement {{name}}. <br /> Cela va générer un zip avec toutes les photos, donc cela peut prendre du temps. <br /> Êtes-vous sûr ?"
    },
    event: {
        create: "Création d'un évènement",
        edit: "Édition de {{name}}",
        name: "Titre",
        host: "Hôte",
        date: "Date",
        location: "Lieux",
        save: "Sauvegarder",
        saved: "Évènement sauvegardé !",
        failed: "Échec de la sauvegarde"
    },
    karaoke: {
        no_song_playing: "Pas de musique en cours",
        now_playing: "La suite",
        sung_by: "Chantée par",
        next_up: "Prochaine musique",
        search: "Rechercher",
        queue: "Queue",
        admin: "Admin",
        amt_songs: "Nb musiques",
        current: "En cours",
        what_to_do: "Que faire ?",
        wtd_add_song: "Ajouter une musique",
        wtd_import: "Importer",
        wtd_settings: "Paramètres",
        title: "Titre",
        artist: "Artiste",
        cover_source: "Source de cover",
        no_cover: "Pas de cover",
        uploaded: "Téléversée",
        format: "Format",
        upload: "Téléverser {{ format }}",
        add: "Créer",
        created: "Musique créée",
        failed: "Échec de la création",
        waiting_song_upload: "En attente de l'envoi d'une musique",
        upload_a_song: "Téléverser une musique",
        upload_in_progress: "En cours...",
        rescan_songs: "Re-scanner les musiques",
        songs_rescanned: "Scan des musiques terminé",
        song_rescan_failed: "Échec du scan des musiques",
        adding_to_queue: 'Ajouté à la queue',
        empty_queue: "Queue vide",
        skip: {
            title: "Passer la musique?",
            content: "Il y a une musique en cours, vous voulez peut-être plutôt la mettre en queue.",
            skip_and_play: "Passer & lancer",
            skip: "Passer",
        },
        queue_remove: {
            title: "Retirer de la queue",
            content: "Êtes-vous sur de vouloir retirer {{name}} de la queue ?",
            removed: "Suppression de la queue"
        },
        current_remove: {
            title: "Stopper la musique",
            content: "Êtes-vous sur de vouloir stopper la musique en cours ({{name}}) ?",
            remove: "Stopper",
        },
        volume: {
            title: "Volume",
            instrumental: "Instrumental",
            vocals: "Voix",
            full: "Complet"
        }
    },
    photobooth: {
        remote_take_picture: "Photo à distance",
        amt_hand_taken: "Prises manuellement",
        amt_unattended: "Prises automatiquement"
    },
    disabled: {
        partyhall_disabled: 'PartyHall désactivé',
        msg: 'Ce photomaton est désactivé par un admin. Désolé pour le dérangement.'
    }
};