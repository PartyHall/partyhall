package cmd

import (
	"os"
	"strconv"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/orm"
	"github.com/partyhall/partyhall/services"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User related commands",
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a user, user-create [username] [password] [name]",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			logs.Error(err)
			os.Exit(1)
		}

		username := args[0]
		password, err := services.GetArgon().Hash(args[1])
		if err != nil {
			logs.Errorf("Failed to hash password: ", err)
			os.Exit(1)
		}

		row := orm.GET.DB.QueryRow("SELECT COUNT(*) FROM ph_user")
		if row.Err() != nil {
			logs.Errorf("Failed to create user: %v", row.Err)
			os.Exit(1)
		}

		var amtUsers int
		err = row.Scan(&amtUsers)
		if err != nil {
			logs.Errorf("Failed to create user: %v", err)
			os.Exit(1)
		}

		roles := []string{models.ROLE_USER}
		if amtUsers == 0 {
			roles = append(roles, models.ROLE_ADMIN)
		}

		user := models.User{
			Name:     args[2],
			Username: username,
			Password: password,
			Roles:    models.Roles(roles),
		}

		dbUser, err := orm.GET.Users.Create(user)
		if err != nil {
			logs.Errorf("Failed to create user: ", err)
			os.Exit(1)
		}

		logs.Infof("User %v (%v) created", dbUser.Id, dbUser.Username)
	},
}

var getUserCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a user by id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			logs.Error(err)
			os.Exit(1)
		}

		userId, err := strconv.Atoi(args[0])
		if err != nil {
			logs.Errorf("Failed to parse id: ", err)
			os.Exit(1)
		}

		dbUser, err := orm.GET.Users.Get(userId)
		if err != nil {
			logs.Errorf("Failed to create user: ", err)
			os.Exit(1)
		}

		logs.Infof("User %v (%v) has the following roles: %v", dbUser.Id, dbUser.Username, dbUser.Roles)
	},
}
