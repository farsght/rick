/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rick",
	Short: "Rick is your friendly AI assistant",
	Long: `
	██████╗ ██╗ ██████╗██╗  ██╗
	██╔══██╗██║██╔════╝██║ ██╔╝
	██████╔╝██║██║     █████╔╝ 
	██╔══██╗██║██║     ██╔═██╗ 
	██║  ██║██║╚██████╗██║  ██╗
	╚═╝  ╚═╝╚═╝ ╚═════╝╚═╝  ╚═╝
	
	Rick is your new best friend
	Use 'rick chat' to start a conversation.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initializeConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.llm_cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initializeConfig() {
	if err := InitConfig(); err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		os.Exit(1)
	}
}
