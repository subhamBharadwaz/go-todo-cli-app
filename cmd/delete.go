/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"

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

		id := args[0]
		if err := deleteTask(id); err != nil {
			fmt.Println("Error deleting task:", err)
		} else {
			fmt.Println("Task deleted successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteTask(id string) error {
	file, err := utils.LoadFile("tasks.csv")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer utils.CloseFile(file)

	// Read existing tasks
	tasks, err := utils.ReadTasks(file)
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	// Find and remove the task with the given ID
	var updatedTasks [][]string
	var found bool
	for _, task := range tasks {
		if task[0] == id {
			found = true
			continue
		}
		updatedTasks = append(updatedTasks, task)
	}

	if !found {
		return fmt.Errorf("task with ID %s not found", id)
	}

	// Write updated tasks back to file
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(updatedTasks); err != nil {
		return fmt.Errorf("failed to write tasks: %w", err)
	}

	writer.Flush()
	return nil
}
