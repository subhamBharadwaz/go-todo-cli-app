/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/subhamBharadwaz/go-todo-cli-app/utils"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a todo item by ID from the CSV file",
	Long: `The delete command removes a todo item from the CSV file 'tasks.csv' based on its ID. You need to provide the ID of the todo item you want to delete. If the item exists, it will be removed from the file, and the updated list of tasks will be displayed.

Example usage:
  $ tasks delete 3

This command searches for the todo item with the specified ID, deletes it if found, and updates the 'tasks.csv' file. If the ID does not exist, an error message will be shown. After successful deletion, the updated list of tasks is displayed to reflect the change.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide the ID of the task to delete.")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID", err)
			return
		}
		err = deleteTask(id)
		if err != nil {
			fmt.Println("Error deleting task", err)
			return
		}
		fmt.Printf("Task with ID %d has been deleted\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteTask(id int) error {
	_, err := utils.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
