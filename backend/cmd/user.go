package cmd

import (
	"os"
	"strconv"

	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/services"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User related commands",
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a user, user create [username] [password] [name]",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			log.LOG.Error(err)
			os.Exit(1)
		}

		dal.DB = services.DB

		username := args[0]
		password, err := services.GetArgon().Hash(args[1])
		if err != nil {
			log.LOG.Errorw("Failed to hash password", "err", err)
			os.Exit(1)
		}

		user := models.User{
			Name:     args[2],
			Username: username,
			Password: password,
			Roles:    models.Roles([]string{models.ROLE_USER}),
		}

		dbUser, err := dal.USERS.Create(user)
		if err != nil {
			log.LOG.Errorw("Failed to create user", "err", err)
			os.Exit(1)
		}

		log.LOG.Infof("User %v (%v) created", dbUser.Id, dbUser.Username)
	},
}

var getUserCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a user by id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			log.LOG.Error(err)
			os.Exit(1)
		}

		userId, err := strconv.Atoi(args[0])
		if err != nil {
			log.LOG.Errorw("Failed to parse id", "err", err)
			os.Exit(1)
		}

		dbUser, err := dal.USERS.Get(userId)
		if err != nil {
			log.LOG.Errorw("Failed to get user", "err", err)
			os.Exit(1)
		}

		log.LOG.Infof("User %v (%v) has the following roles: %v", dbUser.Id, dbUser.Username, dbUser.Roles)
	},
}

var createAdminUserCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "Creates an admin user, idempotent so creating an already existing one will simply update its password / display name",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		fullname, _ := cmd.Flags().GetString("name")

		if err := services.Load(); err != nil {
			log.LOG.Error(err)
			os.Exit(1)
		}

		dal.DB = services.DB

		user, err := dal.USERS.FindByUsername(username)
		if err != nil {
			log.LOG.Errorw("Failed to create admin", "err", err)

			os.Exit(1)
		}

		if user == nil {
			user = &models.User{
				Name:     fullname,
				Username: username,
				Roles:    models.Roles([]string{models.ROLE_USER, models.ROLE_ADMIN}),
			}
		}

		password, err = services.GetArgon().Hash(password)
		if err != nil {
			log.LOG.Errorw("Failed to hash password", "err", err)
			os.Exit(1)
		}

		user.Password = password

		dbUser, err := dal.USERS.Upsert(user)
		if err != nil {
			log.LOG.Errorw("Failed to create user", "err", err)
			os.Exit(1)
		}

		action := "created"
		if user.Id > 0 {
			action = "updated"
		}

		log.LOG.Infof("User %v (%v) logs", dbUser.Id, dbUser.Username, action)
	},
}

func getInitializeUserCmd() *cobra.Command {
	createAdminUserCmd.Flags().String("username", "", "The username for the admin")
	createAdminUserCmd.Flags().String("password", "", "The password for the admin")
	createAdminUserCmd.Flags().String("name", "", "Their full name")

	createAdminUserCmd.MarkFlagsOneRequired("username")
	createAdminUserCmd.MarkFlagsOneRequired("password")
	createAdminUserCmd.MarkFlagsOneRequired("name")

	return createAdminUserCmd
}
