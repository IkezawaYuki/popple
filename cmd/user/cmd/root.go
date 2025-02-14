/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/IkezawaYuki/popple/di"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/service"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("作成するユーザーの種類は入力してください：customer or admin")
		var category string
		_, _ = fmt.Scan(&category)
		if category == "customer" {
			var customer model.Customer
			fmt.Println("名前を入力してください：")
			var name string
			_, _ = fmt.Scan(&name)
			customer.Name = name

			fmt.Println("メールアドレスを入力してください")
			var email string
			_, _ = fmt.Scan(&email)
			customer.Email = email

			fmt.Println("パスワードを入力してください")
			var password string
			_, _ = fmt.Scan(&password)
			customer.Password = password

			fmt.Println("WordpressURLを入力してください")
			var wordpressURL string
			_, _ = fmt.Scan(&wordpressURL)
			customer.WordpressURL = wordpressURL

			fmt.Println("FacebookTokenを入力してください")
			var facebookToken string
			_, _ = fmt.Scan(&facebookToken)
			customer.FacebookToken = &facebookToken

			err := srv.Create(cmd.Context(), &customer)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(customer.ID)
			fmt.Println("success!")
		} else if category == "admin" {

		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var srv service.CustomerService

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.user.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	db := infrastructure.GetMysqlConnection()
	srv = di.NewCustomerService(db)
}
